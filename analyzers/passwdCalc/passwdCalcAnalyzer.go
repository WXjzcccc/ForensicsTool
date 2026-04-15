package passwdCalc

import (
	"ForensicsTool/utils"
	"context"
	"crypto/aes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

type PasswdCalc struct {
	ctx context.Context
}

func NewPasswdCalc() *PasswdCalc {
	return &PasswdCalc{}
}

func (p *PasswdCalc) InitCtx(ctx context.Context) {
	p.ctx = ctx
}

func (p *PasswdCalc) startup(ctx context.Context) {
	p.ctx = ctx
}

func unpad(ciphertext []byte) ([]byte, error) {
	length := len(ciphertext)
	if length%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext block size is invalid")
	}
	padding := ciphertext[length-1]
	if int(padding) > length {
		return nil, errors.New("padding is invalid")
	}
	return ciphertext[:length-int(padding)], nil
}

func decryptAES(token string, key, iv []byte) string {
	return crypto.FromBase64String(token).
		Aes().CBC().
		WithKey(key).WithIv(iv).
		Decrypt().ToString()
}

func (p *PasswdCalc) CalWechat(uin, imei string) string {
	/*
		@param uin: 用户id
		@param imei: 手机imei，传入空则默认为1234567890ABCDEF
		@return: 安卓微信EncMicroMsg.db的解密密钥
	*/
	if imei == "" {
		imei = "1234567890ABCDEF"
	}
	pwd := utils.MD5HashHex(imei + uin)[:7]
	return pwd
}

func (p *PasswdCalc) CalWechatIndex(uin, wxid, imei string) string {
	/*
			@param uin: 用户id
			@param imei: 手机imei，传入空则默认为1234567890ABCDEF
		    @param wxid: 微信用户的内部wxid
			@return: 安卓微信EncMicroMsg.db的解密密钥
	*/
	if imei == "" {
		imei = "1234567890ABCDEF"
	}
	atoi, err := strconv.Atoi(uin)
	if err != nil {
		return ""
	}
	if atoi < 0 {
		atoi += 4294967296
	}
	uin = strconv.Itoa(atoi)
	return utils.MD5HashHex(uin + imei + wxid)[:7]
}

func (p *PasswdCalc) CalWildFire(token string) []string {
	/*
		@param token: 野火IM系应用的用户token，在应用目录下shared_prefs/config.xml中的token的值
		@return: data的解密密钥
	*/
	// 新版本密钥
	key, _ := hex.DecodeString("001122334455667778797A7B7C7D7E7F")
	iv := key
	// 旧版本密钥
	key2, _ := hex.DecodeString("7F7E7D7C7B7A79787766554433221100")
	iv2 := key2
	pwd := decryptAES(token, key, iv)
	flag := "使用SQLCipher4进行解密"
	if pwd == "" {
		pwd = decryptAES(token, key2, iv2)
		if pwd == "" {
			return []string{"解密失败", ""}
		}
		flag = "使用SQLCipher3进行解密"
	}
	parts := strings.Split(pwd, "|")
	if len(parts) > 0 {
		pwd = parts[len(parts)-1]
	}
	return []string{pwd, flag}
}

func (p *PasswdCalc) CalMostone(uid string) []string {
	/*
		@param uid: 默往的用户uid，在shared_prefs/im.xml中的userId的值
		@return: 默往数据库msg.db的解密密钥
	*/
	hash := utils.MD5HashHex(uid)
	ret := strings.ToUpper(hash[:6])
	return []string{ret, "使用SQLCipher3进行解密"}
}

func (p *PasswdCalc) CalTiktok(uid string) []string {
	/*
		@param uid: 抖音的用户uid，一般在数据库文件名中就有
		@return: 抖音数据库的解密密钥
	*/
	en := fmt.Sprintf("byte%simwcdb%sdance", uid, uid)
	return []string{en, "使用wcdb进行解密"}
}

func longToBytes(uid int64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(uid))
	return bytes
}

func (p *PasswdCalc) CalBatChat(uid string) []string {
	/*
		 * @param uid: 用户id，数据库名中的数字
			@return: 蝙蝠聊天数据库的解密密钥
	*/
	long_uid, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return []string{"", "uid不合法，应该是一串数字！"}
	}
	sd := -1756908916
	seed := uint32(sd)
	xxhash32_result := utils.XXHash32(longToBytes(long_uid), seed)
	fmt.Println(xxhash32_result)
	pwd := utils.MD5HashHex(string(rune(xxhash32_result)) + uid)
	return []string{strings.ToUpper(pwd), "使用wcdb进行解密"}
}
