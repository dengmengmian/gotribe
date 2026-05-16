package cache

// 本文件验证 Redis key 构造规则是否符合预期。

import (
	"strings"
	"testing"
)

// TestNewKeyBuilder 验证缓存相关逻辑是否符合预期。
func TestNewKeyBuilder(t *testing.T) {
	t.Parallel()

	builder := NewKeyBuilder("gotribe-api")
	if got := builder.RefreshTokenKey("token"); !strings.HasPrefix(got, "gotribe-api:") {
		t.Fatalf("RefreshTokenKey prefix = %q, want prefix %q", got, "gotribe-api")
	}

	fallback := NewKeyBuilder("   ")
	if got := fallback.RateLimitKey("auth", "user:1"); !strings.HasPrefix(got, "application:") {
		t.Fatalf("RateLimitKey prefix = %q, want prefix %q", got, "application")
	}
}

// TestRefreshTokenIndexKeyContainsProjectID 验证缓存相关逻辑是否符合预期。
func TestRefreshTokenIndexKeyContainsProjectID(t *testing.T) {
	t.Parallel()

	builder := NewKeyBuilder("gotribe-api")
	got := builder.RefreshTokenIndexKey("project-a", 99)
	if !strings.Contains(got, ":project-a:99") {
		t.Fatalf("RefreshTokenIndexKey = %q, want project-aware suffix", got)
	}
}

// TestAccessTokenInvalidBeforeKeyContainsProjectID 验证访问令牌失效键包含项目和用户信息。
func TestAccessTokenInvalidBeforeKeyContainsProjectID(t *testing.T) {
	t.Parallel()

	builder := NewKeyBuilder("gotribe-api")
	got := builder.AccessTokenInvalidBeforeKey("project-a", 99)
	if !strings.Contains(got, ":project-a:99") {
		t.Fatalf("AccessTokenInvalidBeforeKey = %q, want project-aware suffix", got)
	}
}
