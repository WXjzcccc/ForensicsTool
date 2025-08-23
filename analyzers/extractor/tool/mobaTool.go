package tool

import (
	"crypto/aes"
	"crypto/sha512"
	"fmt"
	"github.com/WXjzcccc/registry"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/iancoleman/orderedmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"strings"
)

func decryptMoba(ciphertext, masterPasswd string) string {
	hasher := sha512.New()
	hasher.Write([]byte(masterPasswd))
	pwd := hasher.Sum(nil)[:32]
	iv := crypto.FromHexString(strings.Repeat("00", aes.BlockSize)).WithKey(pwd).Aes().ECB().Encrypt().ToBytes()
	return crypto.FromBase64String(ciphertext).WithKey(pwd).WithIv(iv).Aes().CFB8().Decrypt().ToString()
}

func getEncFromIni(file string) ([]string, []string, error) {
	/*
		@param file	:	MobaXterm.ini配置文件路径
		@return		:	密码列表、凭据列表、错误
	*/
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read ini file: %v", err)
	}

	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Content, err := decoder.Bytes(content)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert encoding: %v", err)
	}

	lines := strings.Split(string(utf8Content), "\n")

	var passwords []string
	var credentials []string

	inPasswords := false
	inCredentials := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch line {
		case "[Passwords]":
			inPasswords = true
			inCredentials = false
			continue
		case "[Credentials]":
			inCredentials = true
			inPasswords = false
			continue
		case "":
			inPasswords = false
			inCredentials = false
			continue
		}
		if inPasswords {
			passwords = append(passwords, line)
		}
		if inCredentials {
			credentials = append(credentials, line)
		}
	}
	return passwords, credentials, nil
}

func getPwdOrCredential(reg *registry.Key, path string) ([]string, error) {
	/*
		@param reg	:	注册表实例
		@param path	:	key路径
		@return		:	提取的数据列表，错误
	*/
	var data []string
	nameKey, err := reg.OpenSubKey(path)
	if err != nil {
		return nil, fmt.Errorf("未找到保存的凭据或是密码: %v", err)
	}
	defer nameKey.Close()
	names, err := nameKey.ReadValueNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry key: %v", err)
	}
	for _, name := range names {
		value, _, err := nameKey.GetStringValue(name)
		if err != nil {
			continue
		}
		data = append(data, fmt.Sprintf("%s=%s", name, value))
	}
	return data, nil
}

func getEncFromRegistry(file string) ([]string, []string, error) {
	/*
		@param file	:	NTUSER.DAT配置文件路径
		@return		:	密码列表、凭据列表、错误
	*/
	reg, err := registry.Open(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open registry file: %v", err)
	}
	defer reg.Close()
	rootKey, err := reg.OpenKey("Software\\Mobatek\\MobaXterm")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open registry key: %v", err)
	}
	defer rootKey.Close()
	passwords, _ := getPwdOrCredential(&rootKey, "P") //这俩的错误不处理
	credentials, _ := getPwdOrCredential(&rootKey, "C")
	return passwords, credentials, nil
}

func AnalyzeMobaXterm(file, masterPasswd string) (*orderedmap.OrderedMap, error) {
	/*
		@param	file:	配置文件或注册表文件
		@return:		解析结果、错误
	*/
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a file", file)
	}
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	var passwords []string
	var credentials []string
	if strings.HasSuffix(file, ".ini") {
		passwords, credentials, err = getEncFromIni(file)
		if err != nil {
			return nil, err
		}
	} else {
		passwords, credentials, err = getEncFromRegistry(file)
		if err != nil {
			return nil, err
		}
	}
	var passwordList []*orderedmap.OrderedMap
	for _, password := range passwords {
		info := orderedmap.New()
		info.SetEscapeHTML(false)
		first := ""
		if strings.Contains(password, ":") {
			firstSp := strings.Split(password, ":")
			first = firstSp[0]
		}
		idx := strings.Index(password, "=")
		server := password[len(first)+1 : idx]
		pwd := decryptMoba(password[idx+1:], masterPasswd)
		info.Set("协议与端口", first)
		info.Set("连接账户与地址", server)
		info.Set("保存的密码", pwd)
		passwordList = append(passwordList, info)
	}
	var credentialList []*orderedmap.OrderedMap
	for _, credential := range credentials {
		info := orderedmap.New()
		info.SetEscapeHTML(false)
		user := ""
		encPwd := ""
		idx := strings.Index(credential, "=")
		server := credential[:idx]
		if strings.Contains(credential[idx+1:], ":") {
			user = strings.Split(credential[idx+1:], ":")[0]
			encPwd = strings.Split(credential[idx+1:], ":")[1]
		}
		pwd := decryptMoba(encPwd, masterPasswd)
		info.Set("凭据名", server)
		info.Set("用户名", user)
		info.Set("保存的密码", pwd)
		credentialList = append(credentialList, info)
	}
	result.Set("密码信息", passwordList)
	result.Set("凭据信息", credentialList)
	return result, nil
}
