package service

// 本文件实现 Admin 登录失败计数与临时锁定。
// 设计要点：
//   - Redis 不可用时 fail-open（允许通过 + 写 warn 日志），不阻塞业务
//   - 账户维度与 IP 维度独立计数，二者任一达阈值即锁
//   - 登录成功调用 Reset 清零账户维度

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"gotribe/internal/core/cache"
	"gotribe/internal/core/config"
	applog "gotribe/internal/core/logger"
)

// LockoutTracker 跟踪登录失败次数并执行临时锁定。
type LockoutTracker struct {
	cfg    config.AdminLockoutConfig
	client redis.UniversalClient
	keys   *cache.KeyBuilder
}

// LockedError 表示被锁定时的具体状态。handler 可读取剩余秒数返回给客户端。
type LockedError struct {
	LockedUntil      time.Time
	RemainingSeconds int64
	Scope            string // "account" or "ip"
}

func (e *LockedError) Error() string {
	return fmt.Sprintf("locked by %s until %s", e.Scope, e.LockedUntil.Format(time.RFC3339))
}

// NewLockoutTracker 构造跟踪器。client / keys 可为空（用于禁用场景或单测）。
func NewLockoutTracker(cfg config.AdminLockoutConfig, client redis.UniversalClient, keys *cache.KeyBuilder) *LockoutTracker {
	return &LockoutTracker{cfg: cfg, client: client, keys: keys}
}

// CheckBeforeLogin 在校验密码前检查是否已被锁定。
// 返回 *LockedError 表示已锁；nil 表示可继续尝试。Redis 故障返回 nil 并日志告警（fail-open）。
func (t *LockoutTracker) CheckBeforeLogin(ctx context.Context, username, ip string) error {
	if !t.enabled() {
		return nil
	}
	if scope, until, ok := t.peekLock(ctx, username, ip); ok {
		return &LockedError{LockedUntil: until, RemainingSeconds: int64(time.Until(until).Seconds()), Scope: scope}
	}
	return nil
}

// RecordFailure 登录失败 +1。若达阈值则设置锁定。返回剩余尝试次数（账户维度）或锁定信息。
func (t *LockoutTracker) RecordFailure(ctx context.Context, username, ip string) (remainingAccountAttempts int, locked *LockedError) {
	if !t.enabled() || t.client == nil || t.keys == nil {
		return -1, nil
	}

	now := time.Now()
	accountKey := t.keys.AdminLoginFailAccountKey(username)
	ipKey := t.keys.AdminLoginFailIPKey(ip)

	accountCount := t.incrementWithTTL(ctx, accountKey, t.cfg.AccountLockDuration())
	ipCount := t.incrementWithTTL(ctx, ipKey, t.cfg.IPLockDuration())

	if accountCount < 0 || ipCount < 0 {
		// Redis 故障 fail-open
		return -1, nil
	}

	if accountCount >= int64(t.cfg.AccountMaxFails) {
		until := now.Add(t.cfg.AccountLockDuration())
		return 0, &LockedError{LockedUntil: until, RemainingSeconds: int64(time.Until(until).Seconds()), Scope: "account"}
	}
	if ipCount >= int64(t.cfg.IPMaxFails) {
		until := now.Add(t.cfg.IPLockDuration())
		return int(int64(t.cfg.AccountMaxFails) - accountCount), &LockedError{LockedUntil: until, RemainingSeconds: int64(time.Until(until).Seconds()), Scope: "ip"}
	}
	return int(int64(t.cfg.AccountMaxFails) - accountCount), nil
}

// Reset 登录成功后清零账户维度计数（IP 维度保留，因为可能是攻击中的合法登录）。
func (t *LockoutTracker) Reset(ctx context.Context, username string) {
	if !t.enabled() || t.client == nil || t.keys == nil {
		return
	}
	accountKey := t.keys.AdminLoginFailAccountKey(username)
	if err := t.client.Del(ctx, accountKey).Err(); err != nil && !errors.Is(err, redis.Nil) {
		applog.Sugared().Warnf("LockoutTracker: 重置账户计数失败 user=%s: %v", username, err)
	}
}

// ResetIP 一般用于运维手动解锁某 IP（暂未暴露端点，预留）。
func (t *LockoutTracker) ResetIP(ctx context.Context, ip string) {
	if !t.enabled() || t.client == nil || t.keys == nil {
		return
	}
	ipKey := t.keys.AdminLoginFailIPKey(ip)
	if err := t.client.Del(ctx, ipKey).Err(); err != nil && !errors.Is(err, redis.Nil) {
		applog.Sugared().Warnf("LockoutTracker: 重置 IP 计数失败 ip=%s: %v", ip, err)
	}
}

func (t *LockoutTracker) enabled() bool {
	return t.cfg.Enabled
}

// peekLock 不变更状态，仅检查是否已达阈值。返回 (scope, lockedUntil, locked)。
func (t *LockoutTracker) peekLock(ctx context.Context, username, ip string) (string, time.Time, bool) {
	if t.client == nil || t.keys == nil {
		return "", time.Time{}, false
	}

	accountKey := t.keys.AdminLoginFailAccountKey(username)
	ipKey := t.keys.AdminLoginFailIPKey(ip)

	accountCount, accountTTL := t.getCountAndTTL(ctx, accountKey)
	if accountCount >= int64(t.cfg.AccountMaxFails) && accountTTL > 0 {
		return "account", time.Now().Add(accountTTL), true
	}
	ipCount, ipTTL := t.getCountAndTTL(ctx, ipKey)
	if ipCount >= int64(t.cfg.IPMaxFails) && ipTTL > 0 {
		return "ip", time.Now().Add(ipTTL), true
	}
	return "", time.Time{}, false
}

// getCountAndTTL 读取计数和剩余 TTL；Redis 故障返回 0,0（即不视为锁定）。
func (t *LockoutTracker) getCountAndTTL(ctx context.Context, key string) (int64, time.Duration) {
	count, err := t.client.Get(ctx, key).Int64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			applog.Sugared().Warnf("LockoutTracker: 读取计数失败 key=%s: %v", key, err)
		}
		return 0, 0
	}
	ttl, err := t.client.TTL(ctx, key).Result()
	if err != nil {
		applog.Sugared().Warnf("LockoutTracker: 读取 TTL 失败 key=%s: %v", key, err)
		return count, 0
	}
	if ttl < 0 {
		// -1 表示 key 无 TTL，-2 表示不存在；视作不锁定
		return count, 0
	}
	return count, ttl
}

// incrementWithTTL 原子 INCR 并在首次创建时设 TTL；返回 -1 表示 Redis 故障。
func (t *LockoutTracker) incrementWithTTL(ctx context.Context, key string, ttl time.Duration) int64 {
	count, err := t.client.Incr(ctx, key).Result()
	if err != nil {
		applog.Sugared().Warnf("LockoutTracker: INCR 失败 key=%s: %v", key, err)
		return -1
	}
	if count == 1 {
		if err := t.client.Expire(ctx, key, ttl).Err(); err != nil {
			applog.Sugared().Warnf("LockoutTracker: EXPIRE 失败 key=%s: %v", key, err)
		}
	}
	return count
}
