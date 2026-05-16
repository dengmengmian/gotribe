package bench

// 缓存读写操作的微基准测试。

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"gotribe/internal/core/cache"
)

type benchUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func setupBenchCache(b *testing.B) (*cache.Store, *miniredis.Miniredis, func()) {
	b.Helper()
	mr := miniredis.RunT(b)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	keys := cache.NewKeyBuilder("bench")
	store := cache.NewStore(rdb, keys)
	return store, mr, func() {
		rdb.Close()
		mr.Close()
	}
}

func BenchmarkCache_SetJSON(b *testing.B) {
	store, _, cleanup := setupBenchCache(b)
	defer cleanup()

	user := benchUser{ID: 1, Username: "bench", Email: "bench@test.com"}
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		keys := cache.NewKeyBuilder("bench")
		key := keys.ProfileKey("proj", 1)
		err := store.SetJSON(ctx, key, user, 60*time.Second)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCache_GetJSON_Hit(b *testing.B) {
	store, _, cleanup := setupBenchCache(b)
	defer cleanup()

	user := benchUser{ID: 1, Username: "bench", Email: "bench@test.com"}
	ctx := context.Background()
	keys := cache.NewKeyBuilder("bench")
	key := keys.ProfileKey("proj", 1)
	_ = store.SetJSON(ctx, key, user, 60*time.Second)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out benchUser
		_, err := store.GetJSON(ctx, key, &out)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCache_GetJSON_Miss(b *testing.B) {
	store, _, cleanup := setupBenchCache(b)
	defer cleanup()

	ctx := context.Background()
	keys := cache.NewKeyBuilder("bench")
	key := keys.ProfileKey("proj", 999999)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out benchUser
		_, _ = store.GetJSON(ctx, key, &out)
	}
}
