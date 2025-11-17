package tool

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/donnie4w/go-logger/logger"
	"github.com/iancoleman/orderedmap"
	"github.com/tidwall/gjson"
)

const (
	dbeaverDecryptKey = "babb4a9f774ab853c96c2d653dfe544a"
	dbeaverDecryptIV  = "00000000000000000000000000000000"
)

func decryptDbeaver(encFileData []byte) []byte {
	/*
		@param	encFileData	:	credentials-config.json文件的二进制数据
		@return				:	解密后的结果
	*/
	key, err := hex.DecodeString(dbeaverDecryptKey)
	if err != nil {
		return nil
	}
	iv, err := hex.DecodeString(dbeaverDecryptIV)
	if err != nil {
		return nil
	}
	return crypto.FromBytes(encFileData).WithKey(key).WithIv(iv).Aes().CBC().PKCS7Padding().Decrypt().ToBytes()[16:]
}

func AnalyzeDbeaver(folder string) (*orderedmap.OrderedMap, error) {
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	var tmp []*orderedmap.OrderedMap
	var pwdJson gjson.Result
	var conJson gjson.Result
	fileInfo, err := os.Stat(folder)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		logger.Errorf("%s is not a folder", folder)
		return nil, fmt.Errorf("%s is not a folder", folder)
	}
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filename := info.Name()
		if filename == "credentials-config.json" {
			encData, err := os.ReadFile(path)
			if err != nil {
				logger.Errorf("read %s failed: %v", path, err)
				return err
			}
			decData := decryptDbeaver(encData)
			pwdJson = gjson.Parse(string(decData))
		}
		if filename == "data-sources.json" {
			conData, err := os.ReadFile(path)
			if err != nil {
				logger.Errorf("read %s failed: %v", path, err)
				return err
			}
			conJson = gjson.Parse(string(conData))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	children := conJson.Get("connections")
	if !children.Exists() {
		logger.Errorf("data-sources.json is empty")
		return nil, fmt.Errorf("data-sources.json is empty")
	}
	children.ForEach(func(key, value gjson.Result) bool {
		info := orderedmap.New()
		info.SetEscapeHTML(false)
		info.Set("数据库类型", value.Get("provider").String())
		info.Set("连接名", value.Get("name").String())
		info.Set("地址", value.Get("configuration").Get("host").String())
		info.Set("端口", value.Get("configuration").Get("port").String())
		info.Set("数据库名", value.Get("configuration").Get("database").String())
		info.Set("账户", pwdJson.Get(key.String()).Get("#connection").Get("user").String())
		info.Set("密码", pwdJson.Get(key.String()).Get("#connection").Get("password").String())
		tmp = append(tmp, info)
		return true
	})
	result.Set("DBeaver连接信息", tmp)
	return result, nil
}
