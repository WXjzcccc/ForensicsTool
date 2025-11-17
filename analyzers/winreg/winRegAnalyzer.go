package winreg

import (
	"ForensicsTool/analyzers/winreg/structs"
	"ForensicsTool/utils"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/WXjzcccc/registry"
	"github.com/donnie4w/go-logger/logger"
	"github.com/iancoleman/orderedmap"
)

type WinReg struct {
	systemReg   registry.Registry
	samReg      registry.Registry
	softwareReg registry.Registry
	ntRegList   []registry.Registry
}

type Reg struct {
	ctx context.Context
}

func NewReg() *Reg {
	return &Reg{}
}

func (r *Reg) InitCtx(ctx context.Context) {
	r.ctx = ctx
}

func checkSubKeyExist(key registry.Key, keyName string) bool {
	/*
		@param	key:		注册表的键
		@param	keyName:	子键名
		@return:			是否存在
	*/
	names, err := key.ReadSubKeyNames(-1)
	if err != nil {
		logger.Error("读取子键名称失败:", err)
		return false
	}
	for _, name := range names {
		if name == keyName {
			return true
		}
	}
	return false
}

func getStringValue(key registry.Key, valueName string) string {
	result, _, err := key.GetStringValue(valueName)
	if err != nil {
		logger.Error("获取字符串值失败:", err)
		return ""
	}
	return result
}

func getStringsValue(key registry.Key, valueName string) []string {
	result, _, err := key.GetStringsValue(valueName)
	if err != nil || len(result) == 0 {
		logger.Errorf("获取字符串数组值<%s>失败:%v", valueName, err)
		return []string{""}
	}
	return result
}

func getBinaryValue(key registry.Key, valueName string) []byte {
	result, _, err := key.GetBinaryValue(valueName)
	if err != nil {
		logger.Error("获取二进制值失败:", err)
		return nil
	}
	return result
}

func getIntValue(key registry.Key, valueName string) uint64 {
	result, _, err := key.GetIntegerValue(valueName)
	if err != nil {
		logger.Error("获取整数值失败:", err)
		return 0
	}
	return result
}

func timestampByts(timeBytes []byte) string {
	timestamp := binary.LittleEndian.Uint64(timeBytes)
	unixSeconds := int64(timestamp)/1e7 - 11644473600
	//t := time.Unix(unixSeconds, 0)
	//formattedDate := t.Format("2006-01-02 15:04:05")
	formattedDate := utils.DefaultTimestampToDatetime(unixSeconds, "", "")
	return formattedDate
}

func timestamp(ts uint64) string {
	//t := time.Unix(int64(ts), 0)
	//formattedDate := t.Format("2006-01-02 15:04:05")
	formattedDate := utils.DefaultTimestampToDatetime(int64(ts), "", "")
	return formattedDate
}

func byte2mac(macBytes []byte) string {
	var result string
	str := hex.EncodeToString(macBytes)
	for idx, s := range str {
		result = result + string(s)
		if idx%2 == 1 && idx != len(str)-1 {
			result = result + ":"
		}
	}
	return result
}

func (w *WinReg) getBootKey() []byte {
	key, err := w.systemReg.OpenKey(fmt.Sprintf("%s\\Control\\LSA", w.getControlSet()))
	if err != nil {
		logger.Error("打开LSA注册表键失败:", err)
		return nil
	}
	defer key.Close()
	var bootKeyObf []byte
	jd, err := key.OpenSubKey("JD")

	if err != nil {
		logger.Error("打开JD子键失败:", err)
		return nil
	}
	skew, err := key.OpenSubKey("Skew1")
	if err != nil {
		logger.Error("打开Skew1子键失败:", err)
		return nil
	}
	gbg, err := key.OpenSubKey("GBG")
	if err != nil {
		logger.Error("打开GBG子键失败:", err)
		return nil
	}
	data, err := key.OpenSubKey("Data")
	if err != nil {
		logger.Error("打开Data子键失败:", err)
		return nil
	}
	jdClassName, err := hex.DecodeString(jd.GetClassName())
	if err != nil {
		logger.Error("解码JD类名失败:", err)
		return nil
	}
	skewClassName, err := hex.DecodeString(skew.GetClassName())
	if err != nil {
		logger.Error("解码Skew1类名失败:", err)
		return nil
	}
	gbgClassName, err := hex.DecodeString(gbg.GetClassName())
	if err != nil {
		logger.Error("解码GBG类名失败:", err)
		return nil
	}
	dataClassName, err := hex.DecodeString(data.GetClassName())
	if err != nil {
		logger.Error("解码Data类名失败:", err)
		return nil
	}
	bootKeyObf = append(bootKeyObf, jdClassName...)
	bootKeyObf = append(bootKeyObf, skewClassName...)
	bootKeyObf = append(bootKeyObf, gbgClassName...)
	bootKeyObf = append(bootKeyObf, dataClassName...)
	transforms := []int{8, 5, 4, 2, 11, 9, 13, 3, 0, 6, 1, 12, 14, 10, 15, 7}
	var bootKey []byte
	for i := 0; i < len(bootKeyObf); i++ {
		bootKey = append(bootKey, bootKeyObf[transforms[i]:transforms[i]+1]...)
	}
	return bootKey
}

func (w *WinReg) getControlSet() string {
	key, err := w.systemReg.OpenKey("select")
	if err != nil {
		logger.Error("打开select注册表键失败:", err)
		return "ControlSet001"
	}
	defer key.Close()
	set := getIntValue(key, "Current")
	return fmt.Sprintf("ControlSet%03d", set)
}

func (w *WinReg) getTimeZone() string {
	key, err := w.systemReg.OpenKey(fmt.Sprintf("%s\\Control\\TimeZoneInformation", w.getControlSet()))
	if err != nil {
		logger.Error("打开TimeZoneInformation注册表键失败:", err)
		return ""
	}
	defer key.Close()
	timeZoneName := getStringValue(key, "TimeZoneKeyName")
	if timeZoneName == "" {
		return ""
	}
	tKey, err := w.softwareReg.OpenKey(fmt.Sprintf("Microsoft\\Windows NT\\CurrentVersion\\Time Zones\\%s", timeZoneName))
	if err != nil {
		logger.Error("打开时区注册表键失败:", err)
		return ""
	}
	defer tKey.Close()
	timeZoneDisplayName := getStringValue(tKey, "Display")
	return timeZoneDisplayName
}

func (w *WinReg) getComputerName() string {
	systemKey, err := w.systemReg.OpenKey(fmt.Sprintf("%s\\Control\\ComputerName\\ComputerName", w.getControlSet()))
	if err != nil {
		logger.Error("打开ComputerName注册表键失败:", err)
		return ""
	}
	defer systemKey.Close()
	return getStringValue(systemKey, "ComputerName")
}

func (w *WinReg) getLastShutdownTime() string {
	systemKey, err := w.systemReg.OpenKey(fmt.Sprintf("%s\\Control\\Windows", w.getControlSet()))
	if err != nil {
		logger.Error("打开Windows注册表键失败:", err)
		return ""
	}
	defer systemKey.Close()
	hexTime := getBinaryValue(systemKey, "ShutdownTime")
	return timestampByts(hexTime)
}

func (w *WinReg) getLastLoginUser() string {
	softwareKey, err := w.softwareReg.OpenKey("Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI")
	if err != nil {
		logger.Error("打开LogonUI注册表键失败:", err)
		return ""
	}
	defer softwareKey.Close()
	return getStringValue(softwareKey, "LastLoggedOnUser")
}

func (w *WinReg) getOrderMap(key, value string) *orderedmap.OrderedMap {
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	result.Set("键", key)
	result.Set("值", value)
	return result
}

func (w *WinReg) getSystemInfo() ([]*orderedmap.OrderedMap, error) {
	/*
		获取系统信息
	*/
	var result []*orderedmap.OrderedMap
	softwareKey, err := w.softwareReg.OpenKey("Microsoft\\Windows NT\\CurrentVersion")
	if err != nil {
		logger.Error("打开CurrentVersion注册表键失败:", err)
		return nil, err
	}
	defer softwareKey.Close()
	result = append(result, w.getOrderMap("Build信息", getStringValue(softwareKey, "BuildLabEx")))
	result = append(result, w.getOrderMap("Build版本", getStringValue(softwareKey, "CurrentBuildNumber")))
	result = append(result, w.getOrderMap("计算机名称", w.getComputerName()))
	result = append(result, w.getOrderMap("版本信息", getStringValue(softwareKey, "EditionID")))
	result = append(result, w.getOrderMap("安装时间(本地时区)", timestamp(getIntValue(softwareKey, "InstallDate"))))
	result = append(result, w.getOrderMap("系统名称", getStringValue(softwareKey, "ProductName")))
	result = append(result, w.getOrderMap("发行ID", getStringValue(softwareKey, "ReleaseId")))
	result = append(result, w.getOrderMap("产品ID", getStringValue(softwareKey, "ProductId")))
	result = append(result, w.getOrderMap("注册所有者", getStringValue(softwareKey, "RegisteredOwner")))
	result = append(result, w.getOrderMap("注册组织", getStringValue(softwareKey, "RegisteredOrganization")))
	result = append(result, w.getOrderMap("注册所有者", getStringValue(softwareKey, "RegisteredOwner")))
	result = append(result, w.getOrderMap("系统时区", w.getTimeZone()))
	result = append(result, w.getOrderMap("最后一次正常关机时间(本地时区)", w.getLastShutdownTime()))
	if checkSubKeyExist(softwareKey, "SoftwareProtectionPlatform") {
		productKey, err := softwareKey.OpenSubKey("SoftwareProtectionPlatform")
		if err != nil {
			logger.Error("打开SoftwareProtectionPlatform子键失败:", err)
			return result, err
		}
		result = append(result, w.getOrderMap("产品密钥备份(非当前密钥)", getStringValue(productKey, "BackupProductKeyDefault")))
	}
	result = append(result, w.getOrderMap("上次登录的用户", w.getLastLoginUser()))
	return result, nil
}

func (w *WinReg) getNetInfo() ([]*orderedmap.OrderedMap, error) {
	/*
		获取网卡信息
	*/
	var result []*orderedmap.OrderedMap
	interfaceKey, err := w.systemReg.OpenKey(fmt.Sprintf("%s\\Services\\Tcpip\\Parameters\\Interfaces", w.getControlSet()))
	if err != nil {
		logger.Error("打开Interfaces注册表键失败:", err)
		return nil, err
	}
	defer interfaceKey.Close()
	interfaceSubKeyNames, err := interfaceKey.ReadSubKeyNames(-1)
	if err != nil {
		logger.Error("读取接口子键名称失败:", err)
		return nil, err
	}
	for _, interfaceSubKeyName := range interfaceSubKeyNames {
		info := orderedmap.New()
		info.SetEscapeHTML(false)
		deviceKey, err := w.systemReg.OpenKey(fmt.Sprintf("%s\\Control\\Network\\{4D36E972-E325-11CE-BFC1-08002BE10318}\\%s\\connection", w.getControlSet(), strings.ToUpper(interfaceSubKeyName)))
		if err != nil {
			logger.Error("打开网络设备注册表键失败:", err)
			info.Set("名称", "")
		} else {
			info.Set("名称", getStringValue(deviceKey, "Name"))
		}
		defer deviceKey.Close()
		macKey, err := w.systemReg.OpenKey(fmt.Sprintf("%s\\Control\\NetworkSetup2\\Interfaces\\%s\\Kernel", w.getControlSet(), strings.ToUpper(interfaceSubKeyName)))
		if err != nil {
			logger.Error("打开MAC地址注册表键失败:", err)
			info.Set("当前MAC地址", "")
			info.Set("物理MAC地址", "")
		} else {
			info.Set("当前MAC地址", byte2mac(getBinaryValue(macKey, "CurrentAddress")))
			info.Set("物理MAC地址", byte2mac(getBinaryValue(macKey, "PermanentAddress")))
		}
		defer macKey.Close()
		interfaceSubKey, err := interfaceKey.OpenSubKey(interfaceSubKeyName)
		if err != nil {
			logger.Error("打开接口子键失败:", err)
			return nil, err
		}
		info.Set("DHCP网络地址", getStringValue(interfaceSubKey, "DhcpIPAddress"))
		info.Set("DHCP网关", getStringsValue(interfaceSubKey, "DhcpDefaultGateway")[0])
		info.Set("DHCP服务地址", getStringValue(interfaceSubKey, "DhcpServer"))
		info.Set("租赁时间", timestamp(getIntValue(interfaceSubKey, "LeaseObtainedTime")))
		info.Set("过期时间", timestamp(getIntValue(interfaceSubKey, "LeaseTerminatesTime")))
		info.Set("IP地址", getStringsValue(interfaceSubKey, "IPAddress")[0])
		info.Set("子网掩码", getStringsValue(interfaceSubKey, "SubnetMask")[0])
		info.Set("网关", getStringsValue(interfaceSubKey, "DefaultGateway")[0])
		result = append(result, info)
	}
	return result, nil
}

func (w *WinReg) getUserInfo() ([]*orderedmap.OrderedMap, error) {
	/*
		获取用户信息
	*/
	var result []*orderedmap.OrderedMap
	domainAccountKey, err := w.samReg.OpenKey("SAM\\Domains\\Account")
	if err != nil {
		logger.Error("打开SAM Domains Account注册表键失败:", err)
		return nil, err
	}
	defer domainAccountKey.Close()
	domainAccountVData := getBinaryValue(domainAccountKey, "V")
	domainAccountFData := getBinaryValue(domainAccountKey, "F")
	t := structs.GetDomainAccountF(domainAccountFData)
	hashedBootKey := structs.GetHBootKey(t.Key0, w.getBootKey())
	machineSid := structs.GetMachineSid(domainAccountVData)
	userAccountKey, err := domainAccountKey.OpenSubKey("Users")
	if err != nil {
		logger.Error("打开Users子键失败:", err)
		return nil, err
	}
	defer userAccountKey.Close()
	userNameKey, err := userAccountKey.OpenSubKey("Names")
	if err != nil {
		logger.Error("打开Names子键失败:", err)
		return nil, err
	}
	defer userNameKey.Close()
	userNames, err := userNameKey.ReadSubKeyNames(-1)
	if err != nil {
		logger.Error("读取用户名失败:", err)
		return nil, err
	}
	for _, userName := range userNames {
		info := orderedmap.New()
		info.SetEscapeHTML(false)
		userKey, err := userNameKey.OpenSubKey(userName)
		if err != nil {
			logger.Error("打开用户子键失败:", err)
			continue
		}
		defer userKey.Close()
		_, valType, err := userKey.GetValue("(default)", []byte{})
		if err != nil {
			logger.Error("获取用户默认值失败:", err)
			continue
		}
		userSid := fmt.Sprintf("%s-%v", machineSid, valType)
		rid := fmt.Sprintf("%08x", valType)
		accountKey, err := userAccountKey.OpenSubKey(rid)
		if err != nil {
			logger.Error("打开用户账户子键失败:", err)
			continue
		}
		accountKey.Close()
		userFData := getBinaryValue(accountKey, "F")
		userVData := getBinaryValue(accountKey, "V")
		userV := structs.GetUserV(userVData)
		userF := structs.GetUserF(userFData)

		info.Set("用户名", userName)
		info.Set("SID", userSid)
		info.Set("密码哈希", structs.DecryptHash(hashedBootKey, userF.RID, userV.NTHash, structs.NTHASH))
		info.Set("上次登录时间", userF.LastLoginTime)
		info.Set("上次密码修改时间", userF.LastPwdChangeTime)
		info.Set("上次登录失败时间", userF.LastFailedLoginTime)
		info.Set("用户属性", strconv.Itoa(int(userF.UserAttribute)))
		info.Set("登录次数", strconv.Itoa(int(userF.LogonCount)))
		info.Set("登录失败次数", strconv.Itoa(int(userF.InValidLoginCount)))
		info.Set("用户全名", userV.FullName)
		info.Set("描述信息", userV.Comment)
		result = append(result, info)
	}
	return result, nil
}

//TODO 默认浏览器，最近访问文档等

func (w *WinReg) getDefaultBrowser() ([]*orderedmap.OrderedMap, error) {
	var result []*orderedmap.OrderedMap
	for idx, nt := range w.ntRegList {
		info := orderedmap.New()
		info.SetEscapeHTML(false)
		rootKey, err := nt.OpenKey("SOFTWARE\\Microsoft\\Windows\\Shell\\Associations\\UrlAssociations\\http\\UserChoice")
		if err != nil {
			logger.Error("打开默认浏览器注册表键失败:", err)
			return nil, err
		}
		defer rootKey.Close()
		browser := getStringValue(rootKey, "ProgId")
		result = append(result, w.getOrderMap(strconv.Itoa(idx), browser))
	}
	return result, nil
}

type RegResult struct {
	Data *orderedmap.OrderedMap `json:"data"`
	Err  string                 `json:"err"`
}

func (r *Reg) AnalyzeWinReg(folder string) *RegResult {
	/*
		@param	folder:	包含了SAM、SYSTEM、SOFTWARE、NTUSER.DAT注册表的文件夹，NTUSER.DAT可以存在多个，只需要后缀名是DAT即可
		@return:		解析结果、错误
	*/
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	var systemReg registry.Registry
	var samReg registry.Registry
	var softwareReg registry.Registry
	var ntRegList []registry.Registry
	fileInfo, err := os.Stat(folder)
	if err != nil {
		logger.Error("检查文件夹状态失败:", err)
		return &RegResult{nil, err.Error()}
	}
	if !fileInfo.IsDir() {
		return &RegResult{nil, fmt.Sprintf("%s is not a folder", folder)}
	}
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		filename := info.Name()
		switch filename {
		case "SYSTEM":
			systemReg, err = registry.Open(path)
		case "SOFTWARE":
			softwareReg, err = registry.Open(path)
			if err != nil {
				logger.Error("打开SOFTWARE注册表文件失败:", err)
				return err
			}
		case "SAM":
			samReg, err = registry.Open(path)
			if err != nil {
				logger.Error("打开SAM注册表文件失败:", err)
				return err
			}
		}
		if strings.HasSuffix(filename, ".DAT") {
			ntReg, err := registry.Open(path)
			if err != nil {
				logger.Error("打开NTUSER.DAT注册表文件失败:", err)
				return err
			}
			ntRegList = append(ntRegList, ntReg)
		}
		return nil
	})
	if err != nil {
		logger.Error("遍历文件夹失败:", err)
		return &RegResult{nil, err.Error()}
	}
	winReg := &WinReg{systemReg, samReg, softwareReg, ntRegList}
	sysInfo, err := winReg.getSystemInfo()
	if err != nil {
		logger.Error("获取系统信息失败:", err)
	} else {
		result.Set("系统信息", sysInfo)
	}
	netInfo, err := winReg.getNetInfo()
	if err != nil {
		logger.Error("获取网卡信息失败:", err)
	} else {
		result.Set("网卡信息", netInfo)
	}
	userInfo, err := winReg.getUserInfo()
	if err != nil {
		logger.Error("获取用户信息失败:", err)
	} else {
		result.Set("用户信息", userInfo)
	}
	defaultBrowser, err := winReg.getDefaultBrowser()
	if err != nil {
		logger.Error("获取默认浏览器信息失败:", err)
	} else {
		result.Set("默认浏览器", defaultBrowser)
	}
	return &RegResult{result, ""}
}
