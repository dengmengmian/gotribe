package bench

// Core 密码哈希与校验的微基准测试。

import (
	"testing"

	"gotribe/internal/auth/core"
)

func BenchmarkPassword_Hash(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.HashPassword("MyS3cur3P@ssw0rd")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPassword_Verify(b *testing.B) {
	hash, err := core.HashPassword("MyS3cur3P@ssw0rd")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ok := core.VerifyPassword(hash, "MyS3cur3P@ssw0rd")
		if !ok {
			b.Fatal("verify failed")
		}
	}
}
