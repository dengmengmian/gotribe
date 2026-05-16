package cache

// 本文件负责初始化 Redis 客户端。

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gotribe/internal/core/config"
)

// NewRedis 初始化 Redis 客户端并验证连接可用性。
func NewRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
