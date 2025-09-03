package database

import (
	"ForensicsTool/utils"
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/WXjzcccc/go-sqlcipher"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

const (
	EnMicroMsgDecryptPragma = "_pragma_cipher_compatibility=1"
	FTSIndexDecryptPragma   = "_pragma_kdf_iter=64000&_pragma_cipher_kdf_algorithm=PBKDF2_HMAC_SHA1&_pragma_cipher_hmac_algorithm=HMAC_SHA1"
	SQLCipher4Pragma        = ""
	SQLCipher3Pragma        = "_pragma_cipher_page_size=1024&_pragma_kdf_iter=64000&_pragma_cipher_kdf_algorithm=PBKDF2_HMAC_SHA1&_pragma_cipher_hmac_algorithm=HMAC_SHA1"
	WcdbPragma              = "_pragma_cipher_page_size=4096&_pragma_kdf_iter=64000&_pragma_cipher_kdf_algorithm=PBKDF2_HMAC_SHA1&_pragma_cipher_hmac_algorithm=HMAC_SHA1"
	NtqqPragma              = "_pragma_cipher_page_size=4096&_pragma_kdf_iter=4000&_pragma_cipher_kdf_algorithm=PBKDF2_HMAC_SHA512"
	FTSIndexDB              = 0
	SQLCipher4DB            = 1
	SQLCipher3DB            = 2
	WCdbDB                  = 3
)

type DecryptDatabase struct {
	ctx context.Context
}

func NewDecryptDatabase() *DecryptDatabase {
	return &DecryptDatabase{}
}

type DecryptResult struct {
	SavePath string `json:"save_path"`
	Wxid     string `json:"wxid"`
	Err      string `json:"err"`
}

func (d *DecryptDatabase) InitCtx(ctx context.Context) {
	d.ctx = ctx
}

func moveElementToFirst(slice []string, element string) []string {
	idx := 0
	for i, v := range slice {
		if v == element {
			idx = i
			break
		}
	}
	// 创建一个新的切片，长度为原切片长度
	newSlice := make([]string, len(slice))

	// 将目标元素放到新切片的第一个位置
	newSlice[0] = slice[idx]

	copy(newSlice[1:idx+1], slice[0:idx])
	copy(newSlice[1+idx:], slice[idx+1:])

	return newSlice
}

func decryptSql(dbPath string) (string, string) {
	dbName := "a" + strings.Replace(filepath.Base(dbPath), ".", "", -1) //数字开头的变量不被接受
	savePath := dbPath + "_dec.db"
	return fmt.Sprintf("ATTACH DATABASE '%s' AS '%s_dec' KEY '';SELECT sqlcipher_export('%s_dec');DETACH DATABASE '%s_dec';", savePath, dbName, dbName, dbName), savePath
}

func (d *DecryptDatabase) DecryptEnMicroMsg(dbPath, password string) *DecryptResult {
	/*
		@param dbPath	: 数据库路径
		@param password	: 解密密码
		@return		 	: 解密后的数据库路径、wxid、error
	*/
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return &DecryptResult{"", "", err.Error()}
	}
	if fileInfo.IsDir() {
		return &DecryptResult{"", "", fmt.Sprintf("%s is not a file", dbPath)}
	}
	pwd := url.QueryEscape(password)
	dbName := fmt.Sprintf("%s?_key=%s&%s", dbPath, pwd, EnMicroMsgDecryptPragma)
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("打开%s失败：%v", dbPath, err)}
	}
	defer db.Close()
	// 获取wxid
	selectSql := "select value from userinfo where id = 2;"
	cur, err := db.Prepare(selectSql)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("准备SQL【%s】失败：%v", selectSql, err)}
	}
	var wxid string
	err = cur.QueryRow().Scan(&wxid)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("读取wxid失败: %v", err)}
	}
	decCmd, savePath := decryptSql(dbPath)
	_, err = db.Exec(decCmd)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("解密数据库失败：%v", err)}
	}
	return &DecryptResult{savePath, wxid, ""}
}

func (d *DecryptDatabase) decryptNormal(dbPath, password string, dbType int) *DecryptResult {
	/*
		@param dbPath	: 数据库路径
		@param password	: 解密密码
		@return		 	: 解密后的数据库路径、wxid（恒定为空字符串）、error
	*/
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return &DecryptResult{"", "", err.Error()}
	}
	if fileInfo.IsDir() {
		return &DecryptResult{"", "", fmt.Sprintf("%s is not a file", dbPath)}
	}
	var dbName string
	pwd := url.QueryEscape(password)
	switch dbType {
	case FTSIndexDB:
		dbName = fmt.Sprintf("%s?_key=%s&%s", dbPath, pwd, FTSIndexDecryptPragma)
	case SQLCipher4DB:
		dbName = fmt.Sprintf("%s?_key=%s", dbPath, pwd)
	case SQLCipher3DB:
		dbName = fmt.Sprintf("%s?_key=%s&%s", dbPath, pwd, SQLCipher3Pragma)
	case WCdbDB:
		dbName = fmt.Sprintf("%s?_key=%s&%s", dbPath, pwd, WcdbPragma)
	default:
		return &DecryptResult{"", "", fmt.Sprintf("不支持的数据库类型：%d", dbType)}
	}
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("打开%s失败：%v", dbPath, err)}
	}
	defer db.Close()
	decCmd, savePath := decryptSql(dbPath)
	_, err = db.Exec(decCmd)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("解密数据库失败：%v", err)}
	}
	return &DecryptResult{savePath, "", ""}
}

func (d *DecryptDatabase) DecryptFTSIndexDB(dbPath, password string) *DecryptResult {
	return d.decryptNormal(dbPath, password, FTSIndexDB)
}

func (d *DecryptDatabase) DecryptSQLCipher4DB(dbPath, password string) *DecryptResult {
	return d.decryptNormal(dbPath, password, SQLCipher4DB)
}

func (d *DecryptDatabase) DecryptSQLCipher3DB(dbPath, password string) *DecryptResult {
	return d.decryptNormal(dbPath, password, SQLCipher3DB)
}

func (d *DecryptDatabase) DecryptWCDB(dbPath, password string) *DecryptResult {
	return d.decryptNormal(dbPath, password, WCdbDB)
}

func (d *DecryptDatabase) DecryptNtqqDB(dbPath, password string) *DecryptResult {
	/*
		@param dbPath	: 数据库路径
		@param password	: 解密密码
		@return		 	: 解密后的数据库路径、wxid（恒定为空字符串）、error
	*/
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return &DecryptResult{"", "", err.Error()}
	}
	if fileInfo.IsDir() {
		return &DecryptResult{"", "", fmt.Sprintf("%s is not a file", dbPath)}
	}
	pathBak := dbPath + ".bak"
	fr, err := os.OpenFile(dbPath, os.O_RDWR, 0666)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("打开文件%s失败：%v", dbPath, err)}
	}
	defer fr.Close()
	_, err = os.Stat(pathBak)
	// 读取源文件内容
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("读取文件%s失败：%v", dbPath, err)}
	}
	if os.IsNotExist(err) {
		// 备份不存在就写备份，否则不再备份，以免原始数据被修改
		fw, err := os.Create(pathBak)
		if err != nil {
			return &DecryptResult{"", "", fmt.Sprintf("创建文件%s失败：%v", pathBak, err)}
		}
		defer fw.Close()
		// 写入备份文件
		_, err = fw.Write(data)
		if err != nil {
			return &DecryptResult{"", "", fmt.Sprintf("写入文件%s失败：%v", pathBak, err)}
		}
	}
	// 提取头部和加密内容
	head := data[:1024]
	encryptedContent := data[1024:]

	// 查找 salt 的起始和结束位置
	saltStart := bytes.Index(head, []byte{0x12, 0x08}) + 2
	if saltStart < 2 {
		return &DecryptResult{"", "", "Salt start not found"}
	}
	saltEnd := bytes.Index(head, []byte{0x1a, 0x07})
	if saltEnd == -1 {
		return &DecryptResult{"", "", "Salt end not found"}
	}

	salt := string(head[saltStart:saltEnd])

	_, err = fr.Seek(0, io.SeekStart)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("Error seeking in source file:%v", err)}
	}
	_, err = fr.Write(encryptedContent)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("Error writing in source file:%v", err)}
	}
	key := url.QueryEscape(utils.MD5HashHex(utils.MD5HashHex(password) + salt))
	methods := []string{"HMAC_SHA1", "HMAC_SHA256", "HMAC_SHA512"}
	for _, method := range methods {
		pragma := NtqqPragma + fmt.Sprintf("&_pragma_cipher_hmac_algorithm=%s", method)
		dbName := fmt.Sprintf("%s?_key=%s&%s", dbPath, key, pragma)
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			continue
		}
		decCmd, savePath := decryptSql(dbPath)
		_, err = db.Exec(decCmd)
		if err != nil {
			err2 := db.Close()
			if err2 != nil {
				continue
			}
			continue
		}
		err = db.Close()
		if err != nil {
			continue
		}
		return &DecryptResult{savePath, "", ""}
	}
	return &DecryptResult{"", "", "解密数据库失败"}
}

func (d *DecryptDatabase) DecryptAMapDB(dbPath string) *DecryptResult {
	/*
		@param dbPath	: 数据库路径
		@return		 	: 解密后的数据库路径、wxid（恒定为空字符串）、error
	*/
	// 获取文件大小并格式化
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("获取文件%s状态失败：%v", dbPath, err)}
	}
	if fileInfo.IsDir() {
		return &DecryptResult{"", "", fmt.Sprintf("%s is not a file", dbPath)}
	}
	size := fileInfo.Size() / 1024
	uintNum := uint32(size)
	// 创建一个4字节的字节切片
	byteSlice := make([]byte, 4)

	// 将uint32转换为字节切片
	byteSlice[0] = byte(uintNum >> 24)
	byteSlice[1] = byte(uintNum >> 16)
	byteSlice[2] = byte(uintNum >> 8)
	byteSlice[3] = byte(uintNum)

	// 将字节切片转换为hex字符串
	sizeHex := hex.EncodeToString(byteSlice)
	// 固定解密密钥
	key := "a4a11bb9ef4b2f4c"

	// SQLite3固定文件头，16字节
	head := "53514C69746520666F726D6174203300"

	// 输出解密后的文件路径
	outPath := strings.Replace(dbPath, "girf_sync.db", "girf_sync_dec.db", -1)
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("读取文件%s失败：%v", dbPath, err)}
	}

	// 获取页大小，读写、缓存等信息，原数据明文，非加密，在第3个8字节数据中给出
	prop := hex.EncodeToString(data[16:24])

	// AES解密
	decData := crypto.
		FromBytes(data).
		SetKey(key).
		Decrypt().
		ToBytes()

	// 在第93-96字节存在第25-28字节的备份数据，拿来即可
	magic := hex.EncodeToString(decData[92:96]) + sizeHex

	// 拼出新的32字节文件头
	trueHead, err := hex.DecodeString(head + prop + magic)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("拼接文件头失败：%v", err)}
	}

	// 写入新文件
	newFileData := append(trueHead, decData[32:]...)
	err = os.WriteFile(outPath, newFileData, 0644)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("写入文件%s失败：%v", outPath, err)}
	}

	return &DecryptResult{outPath, "", ""}
}

func (d *DecryptDatabase) DecryptDingTalkDB(dbPath, password string) *DecryptResult {
	/*
		@param dbPath	: 数据库路径
		@param password	: 解密密码
		@return		 	: 解密后的数据库路径、wxid（恒定为空字符串）、error
	*/
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return &DecryptResult{"", "", err.Error()}
	}
	if fileInfo.IsDir() {
		return &DecryptResult{"", "", fmt.Sprintf("%s is not a file", dbPath)}
	}
	cpus := []string{"armeabi", "armeabi-v7a", "arm64-v8a", "x86", "x86_64", "mips", "mips64"}
	device := strings.Split(password, "/")
	if len(device) != 5 {
		return &DecryptResult{"", "", "密码组成不正确"}
	}
	inCpu := device[1]
	// 将cpus中与device[1]相同的元素位置移到第一个
	newCpus := moveElementToFirst(cpus, inCpu)
	outPath := strings.Replace(dbPath, ".db", "_dec.db", -1)

	// 读取文件内容
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("读取文件%s失败：%v", dbPath, err)}
	}
	sqliteHead := "53514C69746520666F726D6174203300"
	for _, cpu := range newCpus {
		device[1] = cpu
		newDevice := strings.Join(device, "/")
		key := utils.MD5HashHex(newDevice)[:16]
		decData := crypto.
			FromBytes(data).
			SetKey(key).
			Decrypt().
			ToBytes()
		if hex.EncodeToString(decData[:16]) == strings.ToLower(sqliteHead) {
			err = os.WriteFile(outPath, decData, 0644)
			if err != nil {
				return &DecryptResult{"", "", fmt.Sprintf("写入文件%s失败: %v", outPath, err)}
			}
			return &DecryptResult{outPath, "", ""}
		}
	}
	return &DecryptResult{"", "", "解密数据库失败"}
}

func (d *DecryptDatabase) DecryptSystemDataSQLite(dbPath, password string) *DecryptResult {
	/*
		@param dbPath	: 数据库路径
		@param password	: 解密密码
		@return		 	: 解密后的数据库路径、wxid（恒定为空字符串）、error
	*/
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return &DecryptResult{"", "", err.Error()}
	}
	if fileInfo.IsDir() {
		return &DecryptResult{"", "", fmt.Sprintf("%s is not a file", dbPath)}
	}
	file, err := os.Open(dbPath)
	savePath := dbPath + "_dec"
	decFile, err := os.Create(savePath)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("打开文件%s失败：%v", dbPath, err)}
	}
	defer file.Close()
	defer decFile.Close()
	// 每1024字节进行解密
	buffer := make([]byte, 1024)
	chunkCount := 0
	for {
		// 读取数据到缓冲区
		_, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return &DecryptResult{"", "", fmt.Sprintf("读取%s失败：%v", dbPath, err)}
			}
			break
		}
		chunkCount++
		key := utils.SHA1Hash(password)[:16]                                                    // 计算密钥，sha1的前16个字节
		_, err = decFile.Write(crypto.FromBytes(buffer).WithKey(key).RC4().Decrypt().ToBytes()) // 解密，写入数据
		if err != nil {
			return &DecryptResult{"", "", fmt.Sprintf("解密失败：%v", err)}
		}
	}
	data, err := os.ReadFile(savePath)
	if err != nil {
		return &DecryptResult{"", "", fmt.Sprintf("打开%s失败：%v", savePath, err)}
	}
	if hex.EncodeToString(data[:16]) != hex.EncodeToString([]byte("SQLite format 3\x00")) { // 校验前16个字节
		return &DecryptResult{"", "", "解密失败，密钥不正确"}
	}
	return &DecryptResult{savePath, "", ""}
}
