package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	_shortIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

var (
	_serverHash string
	_chars      [62]rune
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Failed to get hostname:", err)
		return
	}
	_serverHash = _sha256hash(hostname)[0:2]
	for i, char := range _shortIDChars {
		_chars[i] = char
	}
}

// ShortIDOptions 短 ID 生成选项。
type ShortIDOptions struct {
	Number        int
	StartWithYear bool
	EndWithHost   bool
}

// ShortIDOption 选项函数。
type ShortIDOption func(*ShortIDOptions)

// WithNumber 设置随机字符数量。
func WithNumber(number int) ShortIDOption {
	return func(o *ShortIDOptions) {
		o.Number = number
	}
}

// WithStartWithYear 以年份开头。
func WithStartWithYear(startWithYear bool) ShortIDOption {
	return func(o *ShortIDOptions) {
		o.StartWithYear = startWithYear
	}
}

// WithEndWithHost 以主机 hash 结尾。
func WithEndWithHost(endWithHost bool) ShortIDOption {
	return func(o *ShortIDOptions) {
		o.EndWithHost = endWithHost
	}
}

// GenShortID 生成短唯一 ID，默认 6 位字符。
func GenShortID(options ...ShortIDOption) string {
	opt := ShortIDOptions{Number: 6}
	for _, option := range options {
		option(&opt)
	}
	id, err := _generateShortID(opt)
	if err != nil {
		return ""
	}
	return strings.ToLower(id)
}

func _generateShortID(opt ShortIDOptions) (string, error) {
	var buffer bytes.Buffer
	randomLength := opt.Number

	if opt.StartWithYear {
		year := time.Now().UTC().Format("06")
		buffer.WriteString(year)
		randomLength -= len(year)
	}

	if opt.EndWithHost {
		buffer.WriteString(_serverHash)
		randomLength -= len(_serverHash)
	}

	if randomLength <= 0 {
		return "", fmt.Errorf("generated ID length is too short")
	}

	data, err := _generateRandomBytes(randomLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	for _, b := range data {
		pick := int(b) % 62
		buffer.WriteRune(_chars[pick])
	}

	return buffer.String(), nil
}

func _sha256hash(text string) string {
	rawBytes := []byte(text)
	h := sha256.Sum256(rawBytes)
	return base64.URLEncoding.EncodeToString(h[:])[:2]
}

func _generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
