package core

// 本文件实现 audience 参数化的 JWT 鉴权中间件。
// 通过显式 audience 参数复用同一份代码服务 ToC API 与 Admin。

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"gotribe/internal/core/errs"
	"gotribe/internal/core/logger"
	coremw "gotribe/internal/core/middleware"
	"gotribe/internal/core/response"
)

// AccessTokenChecker 抽象 access token 失效检查（如 logout 后失效）。
// nil 表示跳过检查。
type AccessTokenChecker interface {
	IsAccessTokenValid(ctx context.Context, audience, projectID string, userID int64, issuedAt time.Time) (bool, error)
}

// JWTMiddleware 创建针对指定 audience 的 JWT 鉴权中间件。
// 解析失败 / audience 不匹配 / token 已被 logout 失效 → 401。
func JWTMiddleware(manager *Manager, audience string, checker AccessTokenChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := ParseBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			response.Error(c, errs.Unauthorized("missing bearer token"))
			c.Abort()
			return
		}
		claims, err := manager.VerifyAccessToken(audience, token)
		if err != nil {
			response.Error(c, errs.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}
		if claims.IssuedAt == nil {
			response.Error(c, errs.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}
		if checker != nil {
			valid, err := checker.IsAccessTokenValid(c.Request.Context(), audience, claims.ProjectID, claims.UserID, claims.IssuedAt.Time)
			if err != nil {
				response.Error(c, errs.ServiceUnavailable("access token validation unavailable", err))
				c.Abort()
				return
			}
			if !valid {
				response.Error(c, errs.Unauthorized("invalid or expired token"))
				c.Abort()
				return
			}
		}

		coremw.SetAuthContext(c, &coremw.AuthContext{
			UserID:    claims.UserID,
			Username:  claims.Username,
			ProjectID: claims.ProjectID,
		})
		c.Set(coremw.ContextKeyUserID, claims.UserID)
		c.Set(coremw.ContextKeyUsername, claims.Username)
		ctx := logger.WithUserID(c.Request.Context(), claims.UserID)
		ctx = logger.WithUsername(ctx, claims.Username)
		if claims.ProjectID != "" {
			c.Set(coremw.ContextKeyProjectID, claims.ProjectID)
			ctx = logger.WithProjectID(ctx, claims.ProjectID)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
