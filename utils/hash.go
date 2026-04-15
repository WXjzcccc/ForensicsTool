package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/OneOfOne/xxhash"
)

// HashType 定义哈希算法类型
type HashType int

const (
	MD5 HashType = iota
	SHA1
	SHA256
)

// Hash 计算输入数据的哈希值
// 输入可以是字符串或字节数组
// 返回字节数组
func Hash(data interface{}, hashType HashType) []byte {
	var input []byte

	switch v := data.(type) {
	case string:
		input = []byte(v)
	case []byte:
		input = v
	default:
		return nil
	}

	switch hashType {
	case MD5:
		h := md5.Sum(input)
		return h[:]
	case SHA1:
		h := sha1.Sum(input)
		return h[:]
	case SHA256:
		h := sha256.Sum256(input)
		return h[:]
	default:
		return nil
	}
}

// HashHex 计算输入数据的哈希值并返回16进制字符串
func HashHex(data interface{}, hashType HashType) string {
	hashBytes := Hash(data, hashType)
	return hex.EncodeToString(hashBytes)
}

// 以下是特定算法的快捷函数

// MD5Hash 计算MD5哈希值
func MD5Hash(data interface{}) []byte {
	return Hash(data, MD5)
}

// MD5HashHex 计算MD5哈希值并返回16进制字符串
func MD5HashHex(data interface{}) string {
	return HashHex(data, MD5)
}

// SHA1Hash 计算SHA1哈希值
func SHA1Hash(data interface{}) []byte {
	return Hash(data, SHA1)
}

// SHA1HashHex 计算SHA1哈希值并返回16进制字符串
func SHA1HashHex(data interface{}) string {
	return HashHex(data, SHA1)
}

// SHA256Hash 计算SHA256哈希值
func SHA256Hash(data interface{}) []byte {
	return Hash(data, SHA256)
}

// SHA256HashHex 计算SHA256哈希值并返回16进制字符串
func SHA256HashHex(data interface{}) string {
	return HashHex(data, SHA256)
}

func XXHash32(data []byte, seed uint32) uint32 {
	xxhash32 := xxhash.NewS32(seed)
	xxhash32.Write(data)
	return xxhash32.Sum32()
}
