package tool

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/iancoleman/orderedmap"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"strings"
)

func randomKey(head []byte) []byte {
	/*	使用Java的随机算法生成key
		@param	head:	加密数据的前8字节
		@return:		生成的key
	*/
	ks := int64(3680984568597093857) / int64(NewJavaRandom(int64(head[5])).NextInt(127))
	random := NewJavaRandom(ks)
	t := int(head[0])
	for i := 0; i < t; i++ {
		random.NextInt64()
	}
	n := random.NextInt64()
	r2 := NewJavaRandom(n)
	ld := []int64{int64(head[4]), r2.NextInt64(), int64(head[7]), int64(head[3]), r2.NextInt64(), int64(head[1]), random.NextInt64(), int64(head[2])}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, ld)
	if err != nil {
		return nil
	}
	keyData := calMd5(buf.Bytes())
	return keyData[:8]
}

func calMd5(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

func decryptFinalShell(encData string) string {
	/*
		@param	encData:	加密的密文
		@return:			解密结果
	*/
	encBytes, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return ""
	}
	return crypto.FromBytes(encBytes[8:]).WithKey(randomKey(encBytes[:8])).Des().ECB().PKCS5Padding().Decrypt().ToString()
}

func AnalyzeFinalShell(folder string) (*orderedmap.OrderedMap, error) {
	/*
		@param	folder:	保存连接配置的文件夹路径
		@return:		解析结果、错误
	*/
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	var tmp []*orderedmap.OrderedMap
	var conJson gjson.Result
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filename := info.Name()
		if strings.HasSuffix(filename, "_connect_config.json") {
			conData, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			conJson = gjson.Parse(string(conData))
			json := orderedmap.New()
			json.SetEscapeHTML(false)
			json.Set("连接名", conJson.Get("name").String())
			json.Set("地址", conJson.Get("host").String())
			json.Set("端口", conJson.Get("port").String())
			json.Set("用户名", conJson.Get("user_name").String())
			json.Set("密码", decryptFinalShell(conJson.Get("password").String()))
			tmp = append(tmp, json)
		} else if strings.HasSuffix(filename, ".json") {
			conData, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			conJson = gjson.Parse(string(conData))
			json := orderedmap.New()
			json.SetEscapeHTML(false)
			json.Set("连接名", strings.Replace(filename, ".json", "", -1))
			json.Set("地址", conJson.Get("host").String())
			json.Set("端口", conJson.Get("port").String())
			json.Set("用户名", conJson.Get("user_name").String())
			json.Set("密码", decryptFinalShell(conJson.Get("password").String()))
			tmp = append(tmp, json)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	result.Set("FinalShell连接信息", tmp)
	return result, nil
}
