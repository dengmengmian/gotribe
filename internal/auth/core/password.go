package core

// 本文件封装密码强度校验、哈希和校验逻辑。
// 由原 internal/auth/password/ 迁移而来。

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	passwordMinLength = 8
	passwordMaxLength = 72
)

// HashPassword 对原始密码进行安全哈希。
func HashPassword(raw string) (string, error) {
	if err := ValidateNewPassword(raw); err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword 校验密码与哈希结果是否匹配。
func VerifyPassword(hashed, raw string) bool {
	if !isBCryptHash(hashed) {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}

// ValidateNewPassword 校验新密码是否满足复杂度要求。
func ValidateNewPassword(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if raw != strings.TrimSpace(raw) {
		return fmt.Errorf("password cannot start or end with whitespace")
	}
	if utf8.RuneCountInString(raw) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters", passwordMinLength)
	}
	if len(raw) > passwordMaxLength {
		return fmt.Errorf("password must be at most %d bytes", passwordMaxLength)
	}

	var hasLetter, hasDigit bool
	for _, r := range raw {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return fmt.Errorf("password must include both letters and digits")
	}
	return nil
}

// isBCryptHash 判断给定字符串是否为受支持的 bcrypt 哈希。
func isBCryptHash(hashed string) bool {
	return strings.HasPrefix(hashed, "$2a$") ||
		strings.HasPrefix(hashed, "$2b$") ||
		strings.HasPrefix(hashed, "$2y$")
}
