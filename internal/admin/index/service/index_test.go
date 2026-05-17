package service

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gotribe/internal/core/database"
)

func TestNewIndexService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	svc := NewService(database.NewTransactionManager(db), nil)
	require.NotNil(t, svc)
}

func TestNewIndexServiceTreatsTypedNilRedisAsDisconnected(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	var redisClient *redis.Client
	svc := NewService(database.NewTransactionManager(db), redisClient)

	err = svc.CacheClear(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "redis 未连接")
}
