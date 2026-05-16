package middleware

// 本文件实现基于 Redis 的接口限流中间件。

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
)

// KeyResolver 定义限流中用于提取身份标识的函数签名。
type KeyResolver func(*gin.Context) (string, bool)

var rateLimitScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("TTL", KEYS[1])
return {current, ttl}
`)

// RateLimit 创建基于 Redis 的限流中间件。
func RateLimit(client *redis.Client, keys *cache.KeyBuilder, enabled bool, scope string, limit int64, window time.Duration, failClosed bool, resolver KeyResolver) gin.HandlerFunc {
	if !enabled || client == nil || limit <= 0 {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		identity, ok := resolver(c)
		if !ok || identity == "" {
			c.Next()
			return
		}

		key := keys.RateLimitKey(scope, identity)
		result, err := rateLimitScript.Run(c.Request.Context(), client, []string{key}, int(window.Seconds())).Result()
		if err != nil {
			if failClosed {
				response.Error(c, errs.ServiceUnavailable("rate limiter unavailable", err))
				c.Abort()
				return
			}
			c.Next()
			return
		}
		values, ok := result.([]any)
		if !ok || len(values) != 2 {
			if failClosed {
				response.Error(c, errs.ServiceUnavailable("rate limiter unavailable", nil))
				c.Abort()
				return
			}
			c.Next()
			return
		}
		count, ok := values[0].(int64)
		if !ok {
			if failClosed {
				response.Error(c, errs.ServiceUnavailable("rate limiter unavailable", nil))
				c.Abort()
				return
			}
			c.Next()
			return
		}
		ttlSeconds, ok := values[1].(int64)
		if !ok || ttlSeconds < 0 {
			ttlSeconds = int64(window.Seconds())
		}
		remaining := limit - count
		if remaining < 0 {
			remaining = 0
		}
		reset := time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix()

		c.Header("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))

		if count > limit {
			c.Header("Retry-After", strconv.FormatInt(ttlSeconds, 10))
			response.Error(c, errs.TooManyRequests("too many requests"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// ClientIPResolver 返回按客户端 IP 进行限流的身份提取器。
func ClientIPResolver() KeyResolver {
	return func(c *gin.Context) (string, bool) {
		return c.ClientIP(), true
	}
}

// UserOrIPResolver 返回优先按用户 ID、否则按 IP 限流的身份提取器。
func UserOrIPResolver() KeyResolver {
	return func(c *gin.Context) (string, bool) {
		if userID, ok := GetUserID(c); ok {
			return fmt.Sprintf("user:%d", userID), true
		}
		return fmt.Sprintf("ip:%s", c.ClientIP()), true
	}
}
