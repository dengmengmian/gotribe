package cache

// 本文件封装 Redis 通用 JSON 缓存读写能力。

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store 封装 Redis 常用读写能力和统一键构建入口。
type Store struct {
	client *redis.Client
	keys   *KeyBuilder
}

// NewStore 创建缓存存储对象。
func NewStore(client *redis.Client, keys *KeyBuilder) *Store {
	return &Store{
		client: client,
		keys:   keys,
	}
}

// ProfileKey 返回当前用户资料缓存键。
func (s *Store) ProfileKey(projectID string, userID int64) string {
	return s.keys.ProfileKey(projectID, userID)
}

// PostDetailKey 返回文章详情缓存键。
func (s *Store) PostDetailKey(projectID, postID string) string {
	return s.keys.PostDetailKey(projectID, postID)
}

// PostListKey 返回文章列表缓存键。
func (s *Store) PostListKey(projectID, filter string) string {
	return s.keys.PostListKey(projectID, filter)
}

// PostListPattern 返回文章列表缓存键的匹配模式，用于批量删除。
func (s *Store) PostListPattern() string {
	return s.keys.PostListPattern()
}

// PostListPatternByProject 返回指定项目的文章列表缓存键匹配模式。
func (s *Store) PostListPatternByProject(projectID string) string {
	return s.keys.PostListPatternByProject(projectID)
}

// Keys 返回底层键构建器，供特殊场景复用。
func (s *Store) Keys() *KeyBuilder {
	return s.keys
}

// SetJSON 将结构化数据序列化后写入 Redis。
func (s *Store) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, body, ttl).Err()
}

// GetJSON 从 Redis 读取并反序列化结构化数据。
func (s *Store) GetJSON(ctx context.Context, key string, dst any) (bool, error) {
	body, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return false, err
	}
	return true, nil
}

// Delete 删除一个或多个缓存键。
func (s *Store) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}

// DeleteByPattern 使用 SCAN 按 pattern 批量删除缓存键，避免 KEYS 命令阻塞 Redis。
func (s *Store) DeleteByPattern(ctx context.Context, pattern string) error {
	if s.client == nil {
		return nil
	}
	var cursor uint64
	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := s.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
