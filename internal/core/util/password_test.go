package utils

import (
	"testing"
)

func TestPassword_GenPasswd(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		expectErr bool
	}{
		{"正常密码", "test1234", false},
		{"纯数字", "12345678", true},
		{"纯字母", "abcdefgh", true},
		{"空密码", "", true},
		{"过短", "ab1", true},
		{"特殊字符", "!@#$%^&*()", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PasswordUtil.GenPasswd(tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("GenPasswd() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && result == "" {
				t.Error("GenPasswd() returned empty string")
			}
		})
	}
}

func TestPassword_ComparePasswd(t *testing.T) {
	password := "test1234"
	hashedPassword, err := PasswordUtil.GenPasswd(password)
	if err != nil {
		t.Fatalf("GenPasswd() error = %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectErr      bool
	}{
		{"正确密码", hashedPassword, password, false},
		{"错误密码", hashedPassword, "wrong5678", true},
		{"空密码", hashedPassword, "", true},
		{"空哈希", "", password, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PasswordUtil.ComparePasswd(tt.hashedPassword, tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("ComparePasswd() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestPassword_Encrypt(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		expectErr bool
	}{
		{"正常字符串", "test5678", false},
		{"空字符串", "", true},
		{"纯数字", "12345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PasswordUtil.Encrypt(tt.source)
			if (err != nil) != tt.expectErr {
				t.Errorf("Encrypt() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && result == "" {
				t.Error("Encrypt() returned empty string")
			}
		})
	}
}

func TestPassword_Compare(t *testing.T) {
	source := "test5678"
	hashedPassword, err := PasswordUtil.Encrypt(source)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectErr      bool
	}{
		{"正确密码", hashedPassword, source, false},
		{"错误密码", hashedPassword, "wrong5678", true},
		{"空密码", hashedPassword, "", true},
		{"空哈希", "", source, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PasswordUtil.Compare(tt.hashedPassword, tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("Compare() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestPassword_Consistency(t *testing.T) {
	password := "test1234"

	hashes := make([]string, 5)
	for i := 0; i < 5; i++ {
		hash, err := PasswordUtil.GenPasswd(password)
		if err != nil {
			t.Fatalf("GenPasswd() error = %v", err)
		}
		hashes[i] = hash
	}

	for i := 0; i < len(hashes); i++ {
		for j := i + 1; j < len(hashes); j++ {
			if hashes[i] == hashes[j] {
				t.Errorf("Generated identical hashes: %s", hashes[i])
			}
		}
	}

	for _, hash := range hashes {
		err := PasswordUtil.ComparePasswd(hash, password)
		if err != nil {
			t.Errorf("ComparePasswd() failed for hash %s: %v", hash, err)
		}
	}
}

func BenchmarkPassword_GenPasswd(b *testing.B) {
	password := "test1234"
	for i := 0; i < b.N; i++ {
		_, err := PasswordUtil.GenPasswd(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPassword_ComparePasswd(b *testing.B) {
	password := "test1234"
	hashedPassword, err := PasswordUtil.GenPasswd(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := PasswordUtil.ComparePasswd(hashedPassword, password)
		if err != nil {
			b.Fatal(err)
		}
	}
}
