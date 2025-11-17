package tool

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/WXjzcccc/registry"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/donnie4w/go-logger/logger"
	"github.com/iancoleman/orderedmap"
	"golang.org/x/crypto/blowfish"
)

const (
	Navicat11KEY = "3DC5CA39"
	Navicat12KEY = "libcckeylibcckey"
	Navicat12IV  = "libcciv libcciv "
)

// Navicat11Cipher 用于 Navicat 11 版本的解密
type Navicat11Cipher struct {
	key []byte
	iv  []byte
}

func NewNavicat11Cipher(userKey string) (*Navicat11Cipher, error) {
	if userKey == "" {
		userKey = Navicat11KEY
	}

	// 初始化密钥
	hasher := sha1.New()
	hasher.Write([]byte(userKey))
	key := hasher.Sum(nil) // Blowfish 需要 4-56 字节的密钥

	// 初始化 IV
	iv, err := hex.DecodeString("FFFFFFFFFFFFFFFF")
	if err != nil {
		return nil, err
	}
	encryptedIV := crypto.FromBytes(iv).WithKey(key).Blowfish().ECB().Encrypt().ToBytes()

	return &Navicat11Cipher{
		key: key,
		iv:  encryptedIV,
	}, nil
}

func (c *Navicat11Cipher) xorBytes(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func (c *Navicat11Cipher) Decrypt(data []byte) []byte {
	outData := make([]byte, len(data))
	cv := make([]byte, len(c.iv))
	copy(cv, c.iv)
	blockSize := blowfish.BlockSize
	blocksLen := len(data) / blockSize
	leftLen := len(data) % blockSize
	for i := 0; i < blocksLen; i++ {
		start := i * blockSize
		end := start + blockSize
		temp := data[start:end]
		decrypted := crypto.FromBytes(temp).WithKey(c.key).Blowfish().ECB().Decrypt().ToBytes()
		decrypted = c.xorBytes(decrypted, cv)
		copy(outData[start:end], decrypted)
		for j := 0; j < len(cv); j++ {
			cv[j] ^= data[start+j]
		}
	}
	if leftLen != 0 {
		encryptedCV := crypto.FromBytes(cv).WithKey(c.key).Blowfish().ECB().Encrypt().ToBytes()
		start := blocksLen * blockSize
		temp := data[start:]
		temp = c.xorBytes(temp, encryptedCV[:leftLen])
		copy(outData[start:], temp)
	}

	return outData
}

func (c *Navicat11Cipher) DecryptString(hexStr string) (string, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	decrypted := c.Decrypt(data)
	return string(decrypted), nil
}

// Navicat12Cipher 用于 Navicat 12+ 版本的解密
type Navicat12Cipher struct {
	aesKey []byte
	aesIV  []byte
}

func NewNavicat12Cipher() *Navicat12Cipher {
	return &Navicat12Cipher{
		aesKey: []byte(Navicat12KEY),
		aesIV:  []byte(Navicat12IV),
	}
}

func (c *Navicat12Cipher) DecryptString(ciphertext string) string {
	return crypto.FromString(ciphertext).WithKey(c.aesKey).WithIv(c.aesIV).Aes().CBC().Decrypt().ToString()
}

// 解密密码封装函数
func decryptNavicat(pwd string) string {
	if pwd == "" {
		return ""
	}

	// 先尝试 Navicat 12+ 解密
	dec12 := NewNavicat12Cipher()
	result := dec12.DecryptString(pwd)
	if utf8.ValidString(result) && result != "" {
		return result
	}
	// 再尝试 Navicat 11 解密
	dec11, err := NewNavicat11Cipher("")
	if err != nil {
		logger.Errorf("failed to create navicat 11 cipher: %v", err)
		return ""
	}
	result, err = dec11.DecryptString(pwd)
	if err == nil {
		return result
	}
	logger.Errorf("failed to decrypt navicat 11 password: %v", err)
	return ""
}

func isValueExist(arr []string, value string) bool {
	valueMap := make(map[string]bool)
	for _, v := range arr {
		valueMap[v] = true
	}
	_, exist := valueMap[value]
	return exist
}

// 获取 Navicat 连接信息
func getNavicatConnections(reg *registry.Registry) (map[string]*registry.Key, error) {

	rootKey, err := reg.OpenKey("Software\\PremiumSoft")
	if err != nil {
		logger.Errorf("failed to open registry key: %v", err)
		return nil, fmt.Errorf("failed to open registry key: %v", err)
	}

	subkeyNames, err := rootKey.ReadSubKeyNames(-1)
	if err != nil {
		logger.Errorf("failed to read subkey names: %v", err)
		return nil, fmt.Errorf("failed to read subkey names: %v", err)
	}

	connections := make(map[string]*registry.Key)
	for _, subkeyName := range subkeyNames {
		subKey, err := rootKey.OpenSubKey(subkeyName)
		if err != nil {
			continue
		}
		names, err := subKey.ReadSubKeyNames(-1)
		if err != nil {
			continue
		}
		if !isValueExist(names, "Servers") {
			continue
		}
		serversKey, err := subKey.OpenSubKey("Servers")
		if err != nil {
			continue
		}

		connections[subkeyName] = &serversKey
	}

	return connections, nil
}

// 获取 MySQL 连接信息
func getNormalDBInfo(serverKeys *registry.Key) (map[string]*orderedmap.OrderedMap, error) {
	info := make(map[string]*orderedmap.OrderedMap)
	connections, err := serverKeys.ReadSubKeyNames(-1)
	if err != nil {
		logger.Errorf("failed to read subkey names: %v", err)
		return nil, fmt.Errorf("failed to read subkey names: %v", err)
	}

	for _, connection := range connections {
		basicInfo := orderedmap.New()

		basicInfo.Set("连接名", connection)
		key, err := serverKeys.OpenSubKey(connection)
		if err != nil {
			logger.Errorf("failed to open subkey: %v", err)
			continue
		}
		if host, _, err := key.GetStringValue("Host"); err == nil {
			basicInfo.Set("地址", host)
		}
		if port, _, err := key.GetIntegerValue("Port"); err == nil {
			basicInfo.Set("端口", fmt.Sprintf("%d", port))
		}
		if username, _, err := key.GetStringValue("UserName"); err == nil {
			basicInfo.Set("用户名", username)
		}
		if pwd, _, err := key.GetStringValue("Pwd"); err == nil {
			basicInfo.Set("密码", decryptNavicat(pwd))
		}
		if path, _, err := key.GetStringValue("QuerySavePath"); err == nil {
			basicInfo.Set("缓存路径", path)
		}
		info[connection] = basicInfo
	}

	return info, nil
}

// 获取 SQLServer 连接信息
func getMSSQLInfo(serverKeys *registry.Key) (map[string]*orderedmap.OrderedMap, error) {
	info := make(map[string]*orderedmap.OrderedMap)
	connections, err := serverKeys.ReadSubKeyNames(-1)
	if err != nil {
		logger.Errorf("failed to read subkey names: %v", err)
		return nil, fmt.Errorf("failed to read subkey names: %v", err)
	}

	for _, connection := range connections {
		basicInfo := orderedmap.New()
		basicInfo.SetEscapeHTML(false)
		basicInfo.Set("连接名", connection)
		key, err := serverKeys.OpenSubKey(connection)
		if err != nil {
			logger.Errorf("failed to open subkey: %v", err)
			continue
		}
		if host, _, err := key.GetStringValue("Host"); err == nil {
			basicInfo.Set("地址", host)
		}
		if port, _, err := key.GetIntegerValue("Port"); err == nil {
			basicInfo.Set("端口", fmt.Sprintf("%d", port))
		}
		if username, _, err := key.GetStringValue("UserName"); err == nil {
			basicInfo.Set("用户名", username)
		}
		if pwd, _, err := key.GetStringValue("Pwd"); err == nil {
			basicInfo.Set("密码", decryptNavicat(pwd))
		}
		if path, _, err := key.GetStringValue("QuerySavePath"); err == nil {
			basicInfo.Set("缓存路径", path)
		}
		if db, _, err := key.GetStringValue("InitialDatabase"); err == nil {
			basicInfo.Set("默认数据库", db)
		}
		if auth, _, err := key.GetStringValue("MSSQLAuthenMode"); err == nil {
			basicInfo.Set("认证模式", auth)
		}
		info[connection] = basicInfo
	}
	return info, nil
}

// 获取SQLite数据
func getSQLiteInfo(serverKeys *registry.Key) (map[string]*orderedmap.OrderedMap, error) {
	info := make(map[string]*orderedmap.OrderedMap)
	connections, err := serverKeys.ReadSubKeyNames(-1)
	if err != nil {
		logger.Errorf("failed to read subkey names: %v", err)
		return nil, fmt.Errorf("failed to read subkey names: %v", err)
	}

	for _, connection := range connections {
		basicInfo := orderedmap.New()
		basicInfo.SetEscapeHTML(false)
		basicInfo.Set("连接名", connection)
		key, err := serverKeys.OpenSubKey(connection)
		if err != nil {
			logger.Errorf("failed to open subkey: %v", err)
			continue
		}
		if name, _, err := key.GetStringValue("DatabaseFileName"); err == nil {
			basicInfo.Set("文件名", name)
		}
		if encrypted, _, err := key.GetIntegerValue("SQLiteEncrypted"); err == nil {
			if encrypted == 1 {
				basicInfo.Set("是否加密", "是")
			} else {
				basicInfo.Set("是否加密", "否")
			}
		}
		if saved, _, err := key.GetStringValue("EncryptionSavePassword"); err == nil {
			if saved == "true" {
				basicInfo.Set("是否保存密码", "是")
			} else {
				basicInfo.Set("是否保存密码", "否")
			}
		}
		if pwd, _, err := key.GetStringValue("EncryptionPassword"); err == nil {
			basicInfo.Set("密码", decryptNavicat(pwd))
		}
		if path, _, err := key.GetStringValue("QuerySavePath"); err == nil {
			basicInfo.Set("缓存路径", path)
		}
		if db, _, err := key.GetStringValue("InitialDatabase"); err == nil {
			basicInfo.Set("默认数据库", db)
		}
		info[connection] = basicInfo
	}
	return info, nil
}

func addRecords(result *orderedmap.OrderedMap, tType, tName string, connections map[string]*registry.Key) error {
	var infoFunc func(serverKeys *registry.Key) (map[string]*orderedmap.OrderedMap, error)
	if keys, ok := connections[tType]; ok {
		switch tType {
		case "Navicat", "NavicatMARIADB", "NavicatMongoDB", "NavicatOra", "NavicatPG":
			infoFunc = getNormalDBInfo
		case "NavicatMSSQL":
			infoFunc = getMSSQLInfo
		case "NavicatSQLite":
			infoFunc = getSQLiteInfo
		}
		info, err := infoFunc(keys)
		if err != nil {
			logger.Errorf("failed to get %s info: %v", tName, err)
			return err
		}
		var records []*orderedmap.OrderedMap
		for _, i := range info {
			records = append(records, i)
		}
		if len(records) > 0 {
			result.Set(tName, records)
		}
	}
	return nil
}

// 分析 Navicat 连接信息
func AnalyzeNavicat(regPath string) (*orderedmap.OrderedMap, error) {
	/*
		@param	refPath:	NTUSER.DAT注册表文件路径
		@return:			解析结果、错误
	*/
	fileInfo, err := os.Stat(regPath)
	if err != nil {
		logger.Errorf("failed to stat file: %v", err)
		return nil, err
	}
	if fileInfo.IsDir() {
		logger.Errorf("%s is not a file", regPath)
		return nil, fmt.Errorf("%s is not a file", regPath)
	}
	reg, err := registry.Open(regPath)
	if err != nil {
		logger.Errorf("failed to open registry file: %v", err)
		return nil, err
	}
	defer reg.Close()
	connections, err := getNavicatConnections(&reg)
	if err != nil {
		return nil, err
	}
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	err = addRecords(result, "Navicat", "MySQL连接信息", connections)
	if err != nil {
		return nil, err
	}
	err = addRecords(result, "NavicatMARIADB", "MariaDB连接信息", connections)
	if err != nil {
		return nil, err
	}
	err = addRecords(result, "NavicatMongoDB", "MongoDB连接信息", connections)
	if err != nil {
		return nil, err
	}
	err = addRecords(result, "NavicatOra", "Oracle连接信息", connections)
	if err != nil {
		return nil, err
	}
	err = addRecords(result, "NavicatPG", "PostgreSQL连接信息", connections)
	if err != nil {
		return nil, err
	}
	err = addRecords(result, "NavicatMSSQL", "SQLServer连接信息", connections)
	if err != nil {
		return nil, err
	}
	err = addRecords(result, "NavicatSQLite", "SQLite连接信息", connections)
	if err != nil {
		return nil, err
	}
	err = reg.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}
