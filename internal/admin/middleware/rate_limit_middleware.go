package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
)

func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	if fillInterval <= 0 || capacity <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Error(c, errs.TooManyRequests("访问限流"))
			c.Abort()
			return
		}
		c.Next()
	}
}
