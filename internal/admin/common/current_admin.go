package common

import (
	"gotribe/internal/core/errs"
	"gotribe/internal/model"

	"github.com/gin-gonic/gin"
)

// CurrentAdmin extracts the authenticated admin that the JWT middleware placed
// on the Gin context. Keep this at the HTTP boundary so services stay Gin-free.
func CurrentAdmin(c *gin.Context) (model.Admin, error) {
	var admin model.Admin
	if c == nil {
		return admin, errs.Unauthorized("用户未登录")
	}

	value, ok := c.Get("user")
	if !ok {
		return admin, errs.Unauthorized("用户未登录")
	}

	admin, ok = value.(model.Admin)
	if !ok || admin.ID == 0 || admin.Username == "" {
		return model.Admin{}, errs.Unauthorized("用户未登录")
	}
	return admin, nil
}
