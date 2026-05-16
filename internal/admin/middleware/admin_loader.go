package middleware

// 本文件实现 admin 用户加载中间件，必须挂在 core.JWTMiddleware 之后。
// JWT 校验已由 core 完成，本中间件读取 ContextKeyUserID，从 DB 加载完整 admin 实体写入 c.Set("user", admin)。

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"gotribe/internal/admin/admin_user/repository"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/errs"
	coremw "gotribe/internal/core/middleware"
	"gotribe/internal/core/response"
)

// AdminUserLoader 在 core.JWTMiddleware 验证完 JWT 之后，加载完整 admin 实体并写入 gin context。
// 必须在 core.JWTMiddleware 之后注册，否则读取不到 ContextKeyUserID。
func AdminUserLoader(adminRepo *repository.Repository, log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := coremw.GetUserID(c)
		if !ok {
			response.Error(c, errs.Unauthorized("missing user context"))
			c.Abort()
			return
		}

		admin, err := adminRepo.GetAdminByID(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.Error(c, errs.Unauthorized("用户不存在"))
			} else {
				if log != nil {
					log.Errorf("查询用户失败: userID=%d, err=%v", userID, err)
				}
				response.Error(c, errs.Internal("系统错误，请稍后重试", nil))
			}
			c.Abort()
			return
		}

		if admin.Status != constant.DEFAULT_ID {
			response.Error(c, errs.Forbidden("当前用户已被禁用"))
			c.Abort()
			return
		}

		c.Set("user", admin)
		c.Next()
	}
}
