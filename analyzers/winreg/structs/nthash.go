package structs

import (
	"ForensicsTool/utils"
	"bytes"
	"encoding/hex"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

/*
解密逻辑和结构体详见https://github.com/skelsec/pypykatz
*/

func uint32To4BytesLittleEndian(i uint32) []byte {
	/*
		小端序转换
	*/
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(i))
	buf.WriteByte(byte(i >> 8))
	buf.WriteByte(byte(i >> 16))
	buf.WriteByte(byte(i >> 24))
	return buf.Bytes()
}

func expandDESKey(key []byte) []byte {
	// 确保密钥长度为7字节，不足补0
	if len(key) < 7 {
		padded := make([]byte, 7)
		copy(padded, key)
		key = padded
	} else {
		key = key[:7]
	}

	var result []byte

	// 第一字节
	b := ((key[0] >> 1) & 0x7f) << 1
	result = append(result, b)

	// 第二字节
	b = ((key[0]&0x01)<<6 | (key[1]>>2)&0x3f) << 1
	result = append(result, b)

	// 第三字节
	b = ((key[1]&0x03)<<5 | (key[2]>>3)&0x1f) << 1
	result = append(result, b)

	// 第四字节
	b = ((key[2]&0x07)<<4 | (key[3]>>4)&0x0f) << 1
	result = append(result, b)

	// 第五字节
	b = ((key[3]&0x0f)<<3 | (key[4]>>5)&0x07) << 1
	result = append(result, b)

	// 第六字节
	b = ((key[4]&0x1f)<<2 | (key[5]>>6)&0x03) << 1
	result = append(result, b)

	// 第七字节
	b = ((key[5]&0x3f)<<1 | (key[6]>>7)&0x01) << 1
	result = append(result, b)

	// 第八字节
	b = (key[6] & 0x7f) << 1
	result = append(result, b)

	return result
}

func rid2Key(rid uint32) ([]byte, []byte) {
	key := uint32To4BytesLittleEndian(rid)
	key1 := []byte{key[0], key[1], key[2], key[3], key[0], key[1], key[2]}
	key2 := []byte{key[3], key[0], key[1], key[2], key[3], key[0], key[1]}
	return expandDESKey(key1), expandDESKey(key2)
}

func GetHBootKey(key0 interface{}, bootKey []byte) []byte {
	QWERTY := []byte("!@#$%^&*()qwertyUIOPAzxcvbnmQQQQQQQQQQQQ)(*@&%\x00")
	DIGITS := []byte("0123456789012345678901234567890123456789\x00")
	switch key0.(type) {
	case *SamKeyData:
		key := key0.(*SamKeyData)
		buf := new(bytes.Buffer)
		buf.Write(key.Salt)
		buf.Write(QWERTY)
		buf.Write(bootKey)
		buf.Write(DIGITS)
		rc4Key := utils.MD5Hash(buf.Bytes())
		buf = new(bytes.Buffer)
		buf.Write(key.Key)
		buf.Write(key.CheckSum)
		hashedBootKey := crypto.FromBytes(buf.Bytes()).WithKey(rc4Key).RC4().Decrypt().ToBytes()
		buf = new(bytes.Buffer)
		buf.Write(hashedBootKey[:16])
		buf.Write(DIGITS)
		buf.Write(hashedBootKey[:16])
		buf.Write(QWERTY)
		checksum := utils.MD5Hash(buf.Bytes())
		if hex.EncodeToString(hashedBootKey[16:]) != hex.EncodeToString(checksum) {
			return nil
		}
		return hashedBootKey
	case *SamKeyDataAes:
		key := key0.(*SamKeyDataAes)
		hashedBootKey := new(bytes.Buffer)
		const n = 16
		var blocks [][]byte
		for i := 0; i < len(key.Data); i += n {
			blocks = append(blocks, key.Data[i:i+n])
		}
		for _, block := range blocks {
			hashedBootKey.Write(
				crypto.FromBytes(block).WithKey(bootKey).WithIv(key.Salt).Aes().CBC().NoPadding().Decrypt().ToBytes())
		}
		return hashedBootKey.Bytes()
	}
	return nil
}

const (
	NTHASH    = 0
	LMHASH    = 1
	NTDEFAULT = "31d6cfe0d16ae931b73c59d7e0c089c0"
	LMDEFAULT = "aad3b435b51404eeaad3b435b51404ee"
)

func DecryptHash(hashedBootKey []byte, rid uint32, hash interface{}, hashType int) string {
	NTPASSWORD := []byte("NTPASSWORD\x00")
	LMPASSWORD := []byte("LMPASSWORD\x00")
	var constant []byte
	defaultHash := ""
	var key []byte
	switch hashType {
	case NTHASH:
		constant = NTPASSWORD
		defaultHash = NTDEFAULT
	case LMHASH:
		constant = LMPASSWORD
		defaultHash = LMDEFAULT
	}
	key1, key2 := rid2Key(rid)
	switch hash.(type) {
	case *SAMHash:
		samHash := hash.(*SAMHash)
		if len(samHash.Hash) == 0 {
			return defaultHash
		}
		buf := new(bytes.Buffer)
		buf.Write(hashedBootKey[:16])
		buf.Write(uint32To4BytesLittleEndian(rid))
		buf.Write(constant)
		rc4Key := utils.MD5Hash(buf.Bytes())
		key = crypto.FromBytes(samHash.Hash).WithKey(rc4Key).RC4().Encrypt().ToBytes()
	case *SAMHashAes:
		samHashAes := hash.(*SAMHashAes)
		if len(samHashAes.Data) == 0 {
			return defaultHash
		}
		keyTmp := new(bytes.Buffer)
		const n = 16
		var blocks [][]byte
		for i := 0; i < len(samHashAes.Data); i += n {
			blocks = append(blocks, samHashAes.Data[i:i+n])
		}
		for _, block := range blocks {
			keyTmp.Write(
				crypto.FromBytes(block).WithKey(hashedBootKey[:16]).WithIv(samHashAes.Salt).Aes().NoPadding().CBC().Decrypt().ToBytes())
		}
		key = keyTmp.Bytes()[:16]
	}
	hash1 := crypto.FromBytes(key[:8]).WithKey(key1).Des().NoPadding().Decrypt().ToBytes()
	hash2 := crypto.FromBytes(key[8:]).WithKey(key2).Des().NoPadding().Decrypt().ToBytes()
	return hex.EncodeToString(hash1) + hex.EncodeToString(hash2)
}
