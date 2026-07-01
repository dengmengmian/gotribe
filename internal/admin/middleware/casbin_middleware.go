package middleware

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	adminrepo "gotribe/internal/admin/admin_user/repository"
	"gotribe/internal/admin/common"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
)

// Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware(tx *database.TransactionManager, enforcer *casbin.SyncedEnforcer, urlPathPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ur := adminrepo.NewRepository(tx)
		actor, err := common.CurrentAdmin(c)
		if err != nil {
			response.Error(c, errs.Unauthorized("用户未登录"))
			c.Abort()
			return
		}
		admin, err := ur.Me(c.Request.Context(), actor)
		if err != nil {
			response.Error(c, errs.Unauthorized("用户未登录"))
			c.Abort()
			return
		}
		if admin.Status != 1 {
			response.Error(c, errs.Forbidden("当前用户已被禁用"))
			c.Abort()
			return
		}

		// 获得用户的全部角色
		roles := admin.Roles
		// 检查是否为超级管理员（拥有排序为1的角色）
		isSuperAdmin := false
		var subs []string
		for _, role := range roles {
			if role.Status == 1 { // 角色状态正常
				subs = append(subs, role.Keyword)
				// 超级管理员判断：排序为1
				if role.Sort == 1 {
					isSuperAdmin = true
				}
			}
		}

		// 超级管理员跳过权限检查
		if isSuperAdmin {
			c.Next()
			return
		}

		// 获得请求路径URL
		obj := strings.TrimPrefix(c.FullPath(), "/"+urlPathPrefix)
		// 获取请求方式
		act := c.Request.Method

		isPass := check(enforcer, subs, obj, act)
		if !isPass {
			response.Error(c, errs.Forbidden("权限不足"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func check(enforcer *casbin.SyncedEnforcer, subs []string, obj string, act string) bool {
	// SyncedEnforcer.Enforce 自身持内部读锁，天然支持并发鉴权，无需外部加锁。
	if enforcer == nil {
		return false
	}

	// 遍历用户的所有角色，只要有一个角色有权限就通过
	for _, sub := range subs {
		if pass, _ := enforcer.Enforce(sub, obj, act); pass {
			return true
		}
	}
	return false
}
