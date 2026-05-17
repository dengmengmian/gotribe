// Package cache provides Redis client initialization, key building, and generic JSON cache storage.
package cache

// 本文件定义 Redis key 构造器，统一缓存和限流 key 命名。

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const defaultKeyPrefix = "application"

// KeyBuilder 负责统一生成 Redis 键名，避免不同模块各自拼接前缀。
type KeyBuilder struct {
	prefix string
}

// NewKeyBuilder 创建统一的 Redis 键构建器。
func NewKeyBuilder(prefix string) *KeyBuilder {
	normalized := strings.TrimSpace(prefix)
	if normalized == "" {
		normalized = defaultKeyPrefix
	}
	return &KeyBuilder{prefix: normalized}
}

// RefreshTokenKey 返回刷新令牌内容对应的缓存键。
// Deprecated: 使用 AudienceRefreshTokenKey，等迁移完成后删除。
func (b *KeyBuilder) RefreshTokenKey(token string) string {
	return fmt.Sprintf("%s:auth:refresh:%s", b.prefix, hash(token))
}

// RefreshTokenIndexKey 返回指定项目和用户的刷新令牌索引键。
// Deprecated: 使用 AudienceRefreshTokenIndexKey。
func (b *KeyBuilder) RefreshTokenIndexKey(projectID string, userID int64) string {
	return fmt.Sprintf("%s:auth:refresh_index:%s:%d", b.prefix, strings.TrimSpace(projectID), userID)
}

// AccessTokenInvalidBeforeKey 返回指定用户访问令牌失效时间的缓存键。
// Deprecated: 使用 AudienceAccessTokenInvalidBeforeKey。
func (b *KeyBuilder) AccessTokenInvalidBeforeKey(projectID string, userID int64) string {
	return fmt.Sprintf("%s:auth:access_invalid_before:%s:%d", b.prefix, strings.TrimSpace(projectID), userID)
}

// AudienceRefreshTokenKey 返回带 audience 命名空间的刷新令牌键。
// 用于 internal/auth/core 多 audience 隔离存储。
func (b *KeyBuilder) AudienceRefreshTokenKey(audience, token string) string {
	return fmt.Sprintf("%s:auth:%s:refresh:%s", b.prefix, strings.TrimSpace(audience), hash(token))
}

// AudienceRefreshTokenIndexKey 返回带 audience 命名空间的用户索引键。
func (b *KeyBuilder) AudienceRefreshTokenIndexKey(audience, projectID string, userID int64) string {
	return fmt.Sprintf("%s:auth:%s:refresh_index:%s:%d", b.prefix, strings.TrimSpace(audience), strings.TrimSpace(projectID), userID)
}

// AudienceAccessTokenInvalidBeforeKey 返回带 audience 命名空间的 access token 失效时间键。
func (b *KeyBuilder) AudienceAccessTokenInvalidBeforeKey(audience, projectID string, userID int64) string {
	return fmt.Sprintf("%s:auth:%s:access_invalid_before:%s:%d", b.prefix, strings.TrimSpace(audience), strings.TrimSpace(projectID), userID)
}

// RateLimitKey 返回限流场景使用的缓存键。
func (b *KeyBuilder) RateLimitKey(scope, identity string) string {
	return fmt.Sprintf("%s:rate_limit:%s:%s", b.prefix, scope, hash(identity))
}

// ProfileKey 返回当前用户资料缓存键。
func (b *KeyBuilder) ProfileKey(projectID string, userID int64) string {
	return fmt.Sprintf("%s:profile:%s:%d", b.prefix, strings.TrimSpace(projectID), userID)
}

// PostDetailKey 返回文章详情缓存键。
func (b *KeyBuilder) PostDetailKey(projectID, postID string) string {
	return fmt.Sprintf("%s:post:detail:%s:%s", b.prefix, strings.TrimSpace(projectID), strings.TrimSpace(postID))
}

// AdminLoginFailAccountKey 返回 admin 账户维度登录失败计数键。
// 使用 hash 隐藏明文用户名，避免 Redis 中泄露用户名枚举。
func (b *KeyBuilder) AdminLoginFailAccountKey(username string) string {
	return fmt.Sprintf("%s:admin:login_fail:account:%s", b.prefix, hash(strings.ToLower(strings.TrimSpace(username))))
}

// AdminLoginFailIPKey 返回 admin IP 维度登录失败计数键。
func (b *KeyBuilder) AdminLoginFailIPKey(ip string) string {
	return fmt.Sprintf("%s:admin:login_fail:ip:%s", b.prefix, hash(strings.TrimSpace(ip)))
}

// AdminTOTPStepJTIKey 返回 TOTP step_token 的 jti 黑名单键，用于防重放。
func (b *KeyBuilder) AdminTOTPStepJTIKey(jti string) string {
	return fmt.Sprintf("%s:admin:totp:step_jti:%s", b.prefix, strings.TrimSpace(jti))
}

// hash 计算字符串的稳定哈希值，用于生成较短的缓存键后缀。
func hash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
