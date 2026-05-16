package core

// 本文件实现基于 Redis 的多 audience refresh token 存储。
// 各 audience 的 key 通过 KeyBuilder 提供的 Audience* 方法显式隔离。

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"gotribe/internal/core/cache"
)

// RefreshSession 表示 refresh token 在 Redis 中的会话信息。
// Audience 字段同时用于校验和 key 隔离。
type RefreshSession struct {
	Audience  string `json:"audience"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	ProjectID string `json:"project_id"`
}

// TokenStore 管理多 audience 的 refresh token 与 access token 失效时间。
type TokenStore struct {
	client *redis.Client
	keys   *cache.KeyBuilder
}

// NewTokenStore 创建新的 token 存储。
func NewTokenStore(client *redis.Client, keys *cache.KeyBuilder) *TokenStore {
	return &TokenStore{client: client, keys: keys}
}

var rotateRefreshTokenScript = redis.NewScript(`
if redis.call("EXISTS", KEYS[1]) == 0 then
  return 0
end
redis.call("DEL", KEYS[1])
redis.call("SREM", KEYS[2], KEYS[1])
redis.call("SET", KEYS[3], ARGV[1], "EX", ARGV[2])
redis.call("SADD", KEYS[4], KEYS[3])
redis.call("EXPIRE", KEYS[4], ARGV[2])
return 1
`)

// Save 保存 refresh token 与 session，并维护 user 索引。
func (s *TokenStore) Save(ctx context.Context, audience, token string, session RefreshSession, ttl time.Duration) error {
	body, err := json.Marshal(session)
	if err != nil {
		return err
	}
	tokenKey := s.keys.AudienceRefreshTokenKey(audience, token)
	indexKey := s.keys.AudienceRefreshTokenIndexKey(audience, session.ProjectID, session.UserID)
	pipe := s.client.TxPipeline()
	pipe.Set(ctx, tokenKey, body, ttl)
	pipe.SAdd(ctx, indexKey, tokenKey)
	pipe.Expire(ctx, indexKey, ttl)
	_, err = pipe.Exec(ctx)
	return err
}

// Get 读取 refresh token 对应的 session。
func (s *TokenStore) Get(ctx context.Context, audience, token string) (RefreshSession, bool, error) {
	body, err := s.client.Get(ctx, s.keys.AudienceRefreshTokenKey(audience, token)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return RefreshSession{}, false, nil
		}
		return RefreshSession{}, false, err
	}
	var session RefreshSession
	if err := json.Unmarshal(body, &session); err != nil {
		return RefreshSession{}, false, err
	}
	return session, true, nil
}

// Rotate 在旧 token 仍有效时原子地替换为新 token。
func (s *TokenStore) Rotate(
	ctx context.Context,
	audience, oldToken string, oldSession RefreshSession,
	newToken string, newSession RefreshSession,
	ttl time.Duration,
) (bool, error) {
	body, err := json.Marshal(newSession)
	if err != nil {
		return false, err
	}
	result, err := rotateRefreshTokenScript.Run(
		ctx, s.client,
		[]string{
			s.keys.AudienceRefreshTokenKey(audience, oldToken),
			s.keys.AudienceRefreshTokenIndexKey(audience, oldSession.ProjectID, oldSession.UserID),
			s.keys.AudienceRefreshTokenKey(audience, newToken),
			s.keys.AudienceRefreshTokenIndexKey(audience, newSession.ProjectID, newSession.UserID),
		},
		body,
		int(ttl.Seconds()),
	).Int64()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// Delete 删除指定 audience 下的 refresh token。
func (s *TokenStore) Delete(ctx context.Context, audience, token string) error {
	tokenKey := s.keys.AudienceRefreshTokenKey(audience, token)
	session, ok, err := s.Get(ctx, audience, token)
	if err != nil {
		return err
	}
	pipe := s.client.TxPipeline()
	pipe.Del(ctx, tokenKey)
	if ok {
		pipe.SRem(ctx, s.keys.AudienceRefreshTokenIndexKey(audience, session.ProjectID, session.UserID), tokenKey)
	}
	_, err = pipe.Exec(ctx)
	return err
}

// InvalidateUserSessions 撤销指定 audience 下某用户的所有 refresh token，
// 并设置 access token 失效时间。
func (s *TokenStore) InvalidateUserSessions(
	ctx context.Context,
	audience, projectID string,
	userID int64,
	invalidBefore time.Time,
	accessTTL time.Duration,
) error {
	indexKey := s.keys.AudienceRefreshTokenIndexKey(audience, projectID, userID)
	tokenKeys, err := s.client.SMembers(ctx, indexKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	pipe := s.client.TxPipeline()
	if len(tokenKeys) > 0 {
		keys := make([]string, 0, len(tokenKeys))
		keys = append(keys, tokenKeys...)
		pipe.Del(ctx, keys...)
	}
	pipe.Del(ctx, indexKey)
	pipe.Set(ctx, s.keys.AudienceAccessTokenInvalidBeforeKey(audience, projectID, userID), invalidBefore.UTC().Unix(), accessTTL)
	_, err = pipe.Exec(ctx)
	return err
}

// IsAccessTokenValid 判断 access token 是否仍处于可用状态。
func (s *TokenStore) IsAccessTokenValid(
	ctx context.Context,
	audience, projectID string,
	userID int64,
	issuedAt time.Time,
) (bool, error) {
	cutoff, err := s.client.Get(ctx, s.keys.AudienceAccessTokenInvalidBeforeKey(audience, projectID, userID)).Int64()
	if err != nil {
		if err == redis.Nil {
			return true, nil
		}
		return false, err
	}
	return issuedAt.UTC().Unix() > cutoff, nil
}
