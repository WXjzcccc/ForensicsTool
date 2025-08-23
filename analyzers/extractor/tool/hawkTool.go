package tool

import (
	"encoding/base64"
	"fmt"
	"github.com/beevik/etree"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/iancoleman/orderedmap"
	"os"
	"strings"
)

func decryptHawk(encData, entity, pwd string) string {
	/*
		@param	encData:	文件中的##0V@后的值
		@param	entity:		密文对应的属性命名
		@param	pwd:		cipher_key，保存在crypto.KEY_256.xml或crypto.KEY_128.xml中
		@return:			解密结果
	*/
	pwdBytes, err := base64.StdEncoding.DecodeString(pwd)
	var data []byte
	if err != nil {
		return ""
	}
	encBytes, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return ""
	}
	iv := encBytes[2:14]
	tag := encBytes[len(encBytes)-16:]
	data = append(data, encBytes[:2][:]...)
	data = append(data, []byte(entity)[:]...)
	encDataBytes := encBytes[14:]
	return crypto.FromBytes(encDataBytes).WithKey(pwdBytes).WithIv(iv).GCMWithTagSize(len(tag), data).Decrypt().ToString()
}

func AnalyzeHawk2(filePath, pwd string) (*orderedmap.OrderedMap, error) {
	/*
		@param	filePath:	Hawk2.xml的文件路径
		@param	pwd:		cipher_key，保存在crypto.KEY_256.xml或crypto.KEY_128.xml中
		@return:			解析结果、错误
	*/
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a file", filePath)
	}
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	var tmp []*orderedmap.OrderedMap
	doc := etree.NewDocument()
	if err = doc.ReadFromFile(filePath); err != nil {
		return nil, err
	}
	root := doc.Root()
	for _, element := range root.ChildElements() {
		entity := element.SelectAttr("name").Value
		if strings.Contains(element.Text(), "##0V@") {
			info := orderedmap.New()
			info.SetEscapeHTML(false)
			encData := strings.Split(element.Text(), "##0V@")[1]
			info.Set("键", entity)
			info.Set("值", decryptHawk(encData, entity, pwd))
			tmp = append(tmp, info)
		}
	}
	result.Set("Hawk2解析", tmp)
	return result, nil
}
