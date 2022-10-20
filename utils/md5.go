package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5Encode 小写
func Md5Encode(data string) string {
	h := md5.New()        // 新建一个对象
	h.Write([]byte(data)) // 转换成一个byte的字节流
	tempstr := h.Sum(nil)
	return hex.EncodeToString(tempstr)
}

// MD5Encode 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// MakePassword 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// VaildPassword  解密
func VaildPassword(plainpwd, salt string, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
