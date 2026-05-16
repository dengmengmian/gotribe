package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

var emailPattern = regexp.MustCompile(`[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?`)

// IsEmail 校验字节串是否为合法 email 地址。
func IsEmail(b []byte) bool {
	return emailPattern.Match(b)
}

// RandString 生成指定长度的随机字符串。
func RandString(n int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(charset)
	result := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	return string(result)
}

// MD5 返回字节串的 MD5 十六进制摘要。
func MD5(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// JSONToMap 将 JSON 字符串解析为 map[string]interface{}。
func JSONToMap(s string) map[string]interface{} {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil
	}
	return m
}
