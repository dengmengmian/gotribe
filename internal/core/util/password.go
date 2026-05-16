// Package utils provides common utility helpers.
package utils

import (
	"gotribe/internal/auth/core"
)

// Password 密码工具类，委托给 internal/auth/core 统一实现。
type Password struct{}

// GenPasswd 密码加密，使用 bcrypt 自适应哈希。明文需满足密码强度要求。
func (p *Password) GenPasswd(passwd string) (string, error) {
	return core.HashPassword(passwd)
}

// ComparePasswd 校验密文与明文是否匹配。
func (p *Password) ComparePasswd(hashPasswd string, passwd string) error {
	if !core.VerifyPassword(hashPasswd, passwd) {
		return errMismatch
	}
	return nil
}

var errMismatch = &passwordMismatch{}

type passwordMismatch struct{}

func (e *passwordMismatch) Error() string { return "password mismatch" }

// Encrypt 使用 bcrypt 加密纯文本。委托给统一密码实现。
func (p *Password) Encrypt(source string) (string, error) {
	return core.HashPassword(source)
}

// Compare 比较密文和明文是否相同。
func (p *Password) Compare(hashedPassword, password string) error {
	return p.ComparePasswd(hashedPassword, password)
}

// PasswordUtil 全局密码工具实例，委托给 internal/auth/core 统一实现。
var PasswordUtil = &Password{}
