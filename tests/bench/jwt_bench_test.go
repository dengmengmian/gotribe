package bench

// Core JWT 签发与校验的微基准测试。

import (
	"testing"
	"time"

	"gotribe/internal/auth/core"
)

func newBenchManager(b *testing.B, ttl time.Duration) *core.Manager {
	b.Helper()
	manager, err := core.NewManager("bench", "benchmark-secret-that-is-32-characters-long", map[string]core.AudienceConfig{
		core.AudienceUser: {
			Audience:        "bench.user",
			AccessTokenTTL:  ttl,
			RefreshTokenTTL: 24 * time.Hour,
		},
	})
	if err != nil {
		b.Fatal(err)
	}
	return manager
}

func BenchmarkJWT_SignAccessToken(b *testing.B) {
	m := newBenchManager(b, 15*time.Minute)
	subject := core.Subject{UserID: 1, Username: "bench_user", ProjectID: "bench-project"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := m.SignAccessToken(core.AudienceUser, subject)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJWT_VerifyAccessToken(b *testing.B) {
	// 使用较长的 TTL 避免 benchmark 期间 token 过期
	m := newBenchManager(b, time.Hour)
	subject := core.Subject{UserID: 1, Username: "bench_user", ProjectID: "bench-project"}
	token, _, err := m.SignAccessToken(core.AudienceUser, subject)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.VerifyAccessToken(core.AudienceUser, token)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJWT_GenerateRefreshToken(b *testing.B) {
	m := newBenchManager(b, 15*time.Minute)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.GenerateRefreshToken()
		if err != nil {
			b.Fatal(err)
		}
	}
}
