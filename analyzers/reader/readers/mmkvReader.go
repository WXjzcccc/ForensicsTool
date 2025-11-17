package readers

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"

	go_mmkv "github.com/WXjzcccc/go-mmkv"
	"github.com/donnie4w/go-logger/logger"
	"github.com/iancoleman/orderedmap"
)

func MMKVReader(path, password string) (*orderedmap.OrderedMap, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Errorf("【mmkvReader】文件不存在")
		} else {
			logger.Errorf("【mmkvReader】文件读取失败")
		}
		return nil, err
	}
	fileMode := fileInfo.Mode()
	dirPath := path
	files := []string{}
	switch {
	case fileMode.IsDir():
		logger.Infof("【mmkvReader】传入的%s是目录", path)
		files = getDirFiles(path)
	case fileMode.IsRegular():
		logger.Infof("【mmkvReader】传入的%s是文件", path)
		dirPath = filepath.Dir(path)
		files = append(files, filepath.Base(path))
	default:
		logger.Errorf("【mmkvReader】文件%s类型错误", path)
		return nil, errors.New("文件类型错误")
	}
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	mr, err := go_mmkv.NewManager(dirPath)
	if err != nil {
		logger.Errorf("【mmkvReader】打开目录%s失败", dirPath)
		return nil, err
	}
	var openErr error
	for _, mmkvFile := range files {
		var tmp []*orderedmap.OrderedMap
		var vault go_mmkv.Vault
		if password == "" {
			vault, openErr = mr.OpenVault(mmkvFile)
		} else {
			vault, openErr = mr.OpenVaultCrypto(mmkvFile, password)
		}
		if openErr != nil {
			logger.Errorf("【mmkvReader】打开文件%s失败：%v", mmkvFile, openErr)
			continue
		}
		for _, key := range vault.Keys() {
			info := orderedmap.New()
			info.SetEscapeHTML(false)
			info.Set("键", key)
			varData, err := vault.GetVarint(key)
			info.Set("VARINT", varData)
			bytesData, err := vault.GetBytes(key)
			if err != nil {
				logger.Errorf("【mmkvReader】获取%s的bytes数据失败：%v", key, err)
				info.Set("BYTES", "")
			} else if len(bytesData) > 0 {
				info.Set("BYTES", hex.EncodeToString(bytesData))
			} else {
				info.Set("BYTES", "")
			}
			stringData, err := vault.GetString(key)
			if err != nil {
				logger.Errorf("【mmkvReader】获取%s的string数据失败：%v", key, err)
				info.Set("STRING", "")
			} else {
				info.Set("STRING", stringData)
			}
			tmp = append(tmp, info)
		}

		result.Set(mmkvFile, tmp)
	}
	return result, openErr
}

func getDirFiles(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		logger.Errorf("【mmkvReader】读取目录%s失败", path)
		return []string{}
	}
	var fileNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".crc") {
			fileNames = append(fileNames, strings.ReplaceAll(file.Name(), ".crc", ""))
		}
	}
	return fileNames
}
