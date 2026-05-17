package service

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"gotribe/internal/core/cache"
	"gotribe/internal/core/config"
)

func newRedisHarness(t *testing.T) (redis.UniversalClient, *miniredis.Miniredis, *cache.KeyBuilder) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return client, mr, cache.NewKeyBuilder("test")
}

func defaultLockoutCfg() config.AdminLockoutConfig {
	return config.AdminLockoutConfig{
		Enabled:            true,
		AccountMaxFails:    3,
		AccountLockMinutes: 10,
		IPMaxFails:         5,
		IPLockMinutes:      30,
	}
}

func TestLockout_DisabledReturnsNil(t *testing.T) {
	tracker := NewLockoutTracker(config.AdminLockoutConfig{Enabled: false}, nil, nil)
	require.NoError(t, tracker.CheckBeforeLogin(context.Background(), "alice", "1.1.1.1"))
	rem, locked := tracker.RecordFailure(context.Background(), "alice", "1.1.1.1")
	require.Equal(t, -1, rem)
	require.Nil(t, locked)
}

func TestLockout_AccountThreshold(t *testing.T) {
	client, _, keys := newRedisHarness(t)
	tracker := NewLockoutTracker(defaultLockoutCfg(), client, keys)
	ctx := context.Background()

	rem, locked := tracker.RecordFailure(ctx, "alice", "1.1.1.1")
	require.Equal(t, 2, rem)
	require.Nil(t, locked)

	rem, locked = tracker.RecordFailure(ctx, "alice", "1.1.1.1")
	require.Equal(t, 1, rem)
	require.Nil(t, locked)

	rem, locked = tracker.RecordFailure(ctx, "alice", "1.1.1.1")
	require.Equal(t, 0, rem)
	require.NotNil(t, locked)
	require.Equal(t, "account", locked.Scope)

	// Subsequent check must report lock
	err := tracker.CheckBeforeLogin(ctx, "alice", "2.2.2.2")
	require.Error(t, err)
	var l *LockedError
	require.True(t, errors.As(err, &l))
	require.Equal(t, "account", l.Scope)
}

func TestLockout_IPThresholdScansMultipleAccounts(t *testing.T) {
	client, _, keys := newRedisHarness(t)
	cfg := defaultLockoutCfg()
	cfg.AccountMaxFails = 100 // 让账户阈值不会先触发
	tracker := NewLockoutTracker(cfg, client, keys)
	ctx := context.Background()

	for i := 0; i < cfg.IPMaxFails; i++ {
		tracker.RecordFailure(ctx, "user"+string(rune('a'+i)), "1.1.1.1")
	}
	err := tracker.CheckBeforeLogin(ctx, "anybody", "1.1.1.1")
	require.Error(t, err)
	var l *LockedError
	require.True(t, errors.As(err, &l))
	require.Equal(t, "ip", l.Scope)
}

func TestLockout_ResetClearsAccount(t *testing.T) {
	client, _, keys := newRedisHarness(t)
	tracker := NewLockoutTracker(defaultLockoutCfg(), client, keys)
	ctx := context.Background()

	tracker.RecordFailure(ctx, "alice", "1.1.1.1")
	tracker.RecordFailure(ctx, "alice", "1.1.1.1")
	tracker.Reset(ctx, "alice")

	rem, locked := tracker.RecordFailure(ctx, "alice", "1.1.1.1")
	require.Equal(t, 2, rem, "reset should clear account counter back to fresh state")
	require.Nil(t, locked)
}
