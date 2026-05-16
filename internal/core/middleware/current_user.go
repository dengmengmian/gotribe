package middleware

// 本文件提供按需加载当前完整用户信息的中间件。

import (
	"context"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
	profileview "gotribe/internal/api/profile/view"
	"gotribe/internal/core/response"
)

// CurrentUserReader 定义按身份读取当前用户资料的能力契约。
type CurrentUserReader interface {
	GetMe(ctx context.Context, projectID string, userID int64) (*profileview.MeView, error)
}

// CurrentUser 创建按需加载当前用户资料的中间件。
func CurrentUser(reader CurrentUserReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, ok := GetAuthContext(c)
		if !ok {
			response.Error(c, errs.Unauthorized("missing auth context"))
			c.Abort()
			return
		}

		currentUser, err := reader.GetMe(c.Request.Context(), auth.ProjectID, auth.UserID)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		c.Set(ContextKeyCurrentUser, currentUser)
		c.Next()
	}
}

// GetCurrentUser 从上下文中读取已加载的当前用户资料。
func GetCurrentUser(c *gin.Context) (*profileview.MeView, bool) {
	value, ok := c.Get(ContextKeyCurrentUser)
	if !ok {
		return nil, false
	}
	currentUser, ok := value.(*profileview.MeView)
	return currentUser, ok && currentUser != nil
}
