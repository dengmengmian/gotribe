// Package utils provides AES-GCM encryption helpers for sensitive value storage.
package utils

// 本文件实现 AES-256-GCM 对称加密，用于落库前对 TOTP secret 等敏感字段做透明加密。
// 密钥固定 32 字节（AES-256 要求），nonce 每次随机 12 字节并与密文拼接存储。

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// AESGCMCipher 是一个线程安全的 AES-256-GCM 加解密器。
type AESGCMCipher struct {
	gcm cipher.AEAD
}

// NewAESGCMCipher 创建 AES-256-GCM 加解密器。key 必须为 32 字节。
func NewAESGCMCipher(key []byte) (*AESGCMCipher, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("aes-gcm key must be 32 bytes, got %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}
	return &AESGCMCipher{gcm: gcm}, nil
}

// Encrypt 加密 plaintext，返回 nonce||ciphertext||tag 的字节序列。
func (c *AESGCMCipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("read nonce: %w", err)
	}
	return c.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt 还原 Encrypt 输出。data 长度必须 >= nonce_size。
func (c *AESGCMCipher) Decrypt(data []byte) ([]byte, error) {
	ns := c.gcm.NonceSize()
	if len(data) < ns {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := data[:ns], data[ns:]
	plaintext, err := c.gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("aes-gcm open: %w", err)
	}
	return plaintext, nil
}

// EncryptToString 加密并返回 base64 字符串，方便存数据库 varchar / text。
func (c *AESGCMCipher) EncryptToString(plaintext string) (string, error) {
	out, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

// DecryptFromString 解码 base64 后解密为字符串。
func (c *AESGCMCipher) DecryptFromString(encoded string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	plaintext, err := c.Decrypt(raw)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
