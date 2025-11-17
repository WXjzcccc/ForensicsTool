package tool

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/donnie4w/go-logger/logger"
	"github.com/iancoleman/orderedmap"
	"gopkg.in/ini.v1"
)

func decryptXShellStr(sid, encPwd string) (string, error) {
	/*
		低版本解密密文
		@param sid:用户的SID
		@param encPwd:	密文
		@return 解密后的密码、错误
	*/
	decoded, err := base64.StdEncoding.DecodeString(encPwd)
	if err != nil {
		logger.Errorf("failed to decode base64 string: %v", err)
		return "", err
	}
	if len(decoded) < 32 {
		logger.Errorf("invalid encrypted password length: %d", len(decoded))
		return "", errors.New("invalid encrypted password length")
	}
	data := decoded[:len(decoded)-32]
	checksum := decoded[len(decoded)-32:]

	// 使用 SHA256 生成密钥
	hasher := sha256.New()
	hasher.Write([]byte(sid))
	keyBytes := hasher.Sum(nil)
	// RC4 解密
	decrypted := crypto.FromBytes(data).
		WithKey(keyBytes).
		RC4().
		Decrypt().
		ToBytes()

	// 验证校验和
	hasher = sha256.New()
	hasher.Write(decrypted)
	calculatedChecksum := hasher.Sum(nil)
	if hex.EncodeToString(calculatedChecksum) != hex.EncodeToString(checksum) {
		logger.Errorf("checksum verification failed: expected %s, got %s", hex.EncodeToString(checksum), hex.EncodeToString(calculatedChecksum))
		return "", errors.New("checksum verification failed")
	}

	return string(decrypted), nil
}

// 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func decryptXShellStrNew(sid string, encPwd string) (string, error) {
	/*
		7.1以后解密密文
		@param sid:用户的SID
		@param encPwd:	密文
		@return 解密后的密码、错误
	*/
	decoded, err := base64.StdEncoding.DecodeString(encPwd)
	if err != nil {
		logger.Errorf("failed to decode base64 string: %v", err)
		return "", err
	}

	if len(decoded) < 32 {
		logger.Errorf("invalid encrypted password length: %d", len(decoded))
		return "", errors.New("invalid encrypted password length")
	}

	data := decoded[:len(decoded)-32]
	checksum := decoded[len(decoded)-32:]

	// 处理 SID
	index := strings.Index(sid, "S-1-5")
	if index == -1 {
		logger.Errorf("invalid SID format: %s", sid)
		return "", errors.New("invalid SID format")
	}

	username := sid[:index]
	ssid := sid[index:]
	key := reverseString(reverseString(username) + ssid)

	// 使用 SHA256 生成密钥
	hasher := sha256.New()
	hasher.Write([]byte(key))
	keyBytes := hasher.Sum(nil)

	// RC4 解密
	decrypted := crypto.FromBytes(data).
		WithKey(keyBytes).
		RC4().
		Decrypt().
		ToBytes()
	// 验证校验和
	hasher = sha256.New()
	hasher.Write(decrypted)
	calculatedChecksum := hasher.Sum(nil)

	if hex.EncodeToString(calculatedChecksum) != hex.EncodeToString(checksum) {
		logger.Errorf("checksum verification failed: expected %s, got %s", hex.EncodeToString(checksum), hex.EncodeToString(calculatedChecksum))
		return "", errors.New("checksum verification failed")
	}

	return string(decrypted), nil
}

// 预处理文件
func preDeal(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		logger.Errorf("failed to read file: %v", err)
		return "", err
	}

	// 替换特殊字符
	data = bytes.ReplaceAll(data, []byte{0x00}, []byte{})
	data = bytes.ReplaceAll(data, []byte{0xff}, []byte{})
	data = bytes.ReplaceAll(data, []byte{0xfe}, []byte{})

	dealFile := file + ".deal"
	err = os.WriteFile(dealFile, data, 0644)
	if err != nil {
		logger.Errorf("failed to write file: %v", err)
		return "", err
	}

	return dealFile, nil
}

// 删除临时文件
func delDeal(file string) error {
	return os.Remove(file)
}

// 分析 XShell/XFtp 配置文件
func AnalyzeXshell(folder string, sid string) (*orderedmap.OrderedMap, error) {
	/*
		@param folder:	配置文件目录
		@param sid:		用户的用户名+sid
		@return:		解析结果
	*/
	conns := orderedmap.New()
	conns.SetEscapeHTML(false)
	var xsh []*orderedmap.OrderedMap
	var xft []*orderedmap.OrderedMap
	fileInfo, err := os.Stat(folder)
	if err != nil {
		logger.Errorf("failed to stat folder: %v", err)
		return nil, err
	}
	if !fileInfo.IsDir() {
		logger.Errorf("%s is not a folder", folder)
		return nil, fmt.Errorf("%s is not a folder", folder)
	}
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Errorf("failed to walk file: %v", err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := info.Name()
		if strings.HasSuffix(filename, ".xsh") {
			originFile := path
			dealFile, err := preDeal(originFile)
			if err != nil {
				logger.Errorf("failed to preDeal file: %v", err)
				return err
			}
			defer delDeal(dealFile)

			cfg, err := ini.LoadSources(ini.LoadOptions{
				AllowBooleanKeys: true,
			}, dealFile)
			if err != nil {
				logger.Errorf("failed to load ini file: %v", err)
				return err
			}

			conn := orderedmap.New()
			conn.SetEscapeHTML(false)
			conn.Set("连接名", strings.TrimSuffix(filename, ".xsh"))
			conn.Set("地址", cfg.Section("CONNECTION").Key("Host").String())
			conn.Set("端口", cfg.Section("CONNECTION").Key("Port").String())

			version := cfg.Section("SessionInfo").Key("Version").String()
			encPwd := cfg.Section("CONNECTION:AUTHENTICATION").Key("Password").String()

			var pwd string
			if version != "7.0" && strings.HasPrefix(version, "7.") {
				pwd, err = decryptXShellStrNew(sid, encPwd)
			} else {
				pwd, err = decryptXShellStr(sid, encPwd)
			}

			if err == nil {
				conn.Set("密码", pwd)
			} else {
				conn.Set("密码", "解密失败: "+err.Error())
			}

			conn.Set("描述", cfg.Section("CONNECTION").Key("Description").String())
			xsh = append(xsh, conn)
		}

		if strings.HasSuffix(filename, ".xfp") {
			originFile := path
			dealFile, err := preDeal(originFile)
			if err != nil {
				logger.Errorf("failed to preDeal file: %v", err)
				return err
			}
			defer delDeal(dealFile)

			cfg, err := ini.LoadSources(ini.LoadOptions{
				AllowBooleanKeys: true,
			}, dealFile)
			if err != nil {
				logger.Errorf("failed to load ini file: %v", err)
				return err
			}

			conn := orderedmap.New()
			conn.Set("连接名", strings.TrimSuffix(filename, ".xfp"))
			conn.Set("地址", cfg.Section("Connection").Key("Host").String())
			conn.Set("端口", cfg.Section("Connection").Key("Port").String())

			version := cfg.Section("SessionInfo").Key("Version").String()
			encPwd := cfg.Section("Connection").Key("Password").String()

			var pwd string
			if version != "7.0" && strings.HasPrefix(version, "7.") {
				pwd, err = decryptXShellStrNew(sid, encPwd)
			} else {
				pwd, err = decryptXShellStr(sid, encPwd)
			}

			if err == nil {
				conn.Set("密码", pwd)
			} else {
				conn.Set("密码", "解密失败: "+err.Error())
			}

			conn.Set("描述", cfg.Section("Connection").Key("Description").String())
			xft = append(xft, conn)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	conns.Set("XShell", xsh)
	conns.Set("XFtp", xft)

	return conns, nil
}
