package handler

import (
	"strconv"
	"strings"

	"gotribe/internal/admin/common"
	"gotribe/internal/admin/role/dto"
	"gotribe/internal/admin/role/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

// Handler 角色处理器。
type Handler struct {
	roleService service.Service
}

// NewHandler 创建角色处理器实例。
func NewHandler(roleService service.Service) *Handler {
	return &Handler{roleService: roleService}
}

// List 获取角色列表。
func (h *Handler) List(c *gin.Context) {
	var req dto.RoleListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	roles, total, err := h.roleService.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"roles": roles, "total": total})
}

// Create 创建角色。
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.roleService.Create(c.Request.Context(), actor, &req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "不能创建比自己等级高或相同等级的角色" {
			response.Error(c, errs.Forbidden(errMsg))
			return
		}
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新角色。
func (h *Handler) Update(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Error(c, errs.BadRequest("角色ID不正确", nil))
		return
	}

	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.roleService.Update(c.Request.Context(), actor, int64(roleID), &req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "未获取到角色信息" {
			response.Error(c, errs.NotFound(errMsg, nil))
			return
		}
		if errMsg == "不能更新比自己角色等级高的角色" || errMsg == "不能把角色等级更新得比当前用户的等级高或相同" {
			response.Error(c, errs.Forbidden(errMsg))
			return
		}
		if errMsg == "更新角色成功，但角色关键字关联的权限接口更新失败" || errMsg == "更新角色成功，但角色关键字关联角色的权限接口策略加载失败" {
			response.Error(c, errs.Internal(errMsg, nil))
			return
		}
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// GetRoleMenusByID 获取角色的权限菜单。
func (h *Handler) GetRoleMenusByID(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Error(c, errs.BadRequest("角色ID不正确", nil))
		return
	}
	menus, err := h.roleService.GetRoleMenusByID(c.Request.Context(), int64(roleID))
	if err != nil {
		response.Error(c, errs.Internal("获取角色的权限菜单失败", err))
		return
	}
	response.OK(c, gin.H{"menus": menus})
}

// UpdateRoleMenusByID 更新角色的权限菜单。
func (h *Handler) UpdateRoleMenusByID(c *gin.Context) {
	var req dto.UpdateRoleMenusRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Error(c, errs.BadRequest("角色ID不正确", nil))
		return
	}

	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.roleService.UpdateRoleMenusByID(c.Request.Context(), actor, int64(roleID), &req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "未获取到角色信息" {
			response.Error(c, errs.NotFound(errMsg, nil))
			return
		}
		if errMsg == "不能更新比自己角色等级高或相等角色的权限菜单" || strings.HasPrefix(errMsg, "无权设置") {
			response.Error(c, errs.Forbidden(errMsg))
			return
		}
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// GetRoleApisByID 获取角色的权限接口。
func (h *Handler) GetRoleApisByID(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Error(c, errs.BadRequest("角色ID不正确", nil))
		return
	}
	apis, err := h.roleService.GetRoleApisByID(c.Request.Context(), int64(roleID))
	if err != nil {
		errMsg := err.Error()
		if errMsg == "未获取到角色信息" {
			response.Error(c, errs.NotFound(errMsg, nil))
			return
		}
		response.Error(c, errs.Internal(errMsg, nil))
		return
	}
	response.OK(c, gin.H{"apis": apis})
}

// UpdateRoleApisByID 更新角色的权限接口。
func (h *Handler) UpdateRoleApisByID(c *gin.Context) {
	var req dto.UpdateRoleApisRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Error(c, errs.BadRequest("角色ID不正确", nil))
		return
	}
	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.roleService.UpdateRoleApisByID(c.Request.Context(), actor, int64(roleID), &req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "未获取到角色信息" {
			response.Error(c, errs.NotFound(errMsg, nil))
			return
		}
		if errMsg == "不能更新比自己角色等级高或相等角色的权限接口" {
			response.Error(c, errs.Forbidden(errMsg))
			return
		}
		if errMsg == "根据接口ID获取接口信息失败" {
			response.Error(c, errs.Internal(errMsg, nil))
			return
		}
		if strings.HasPrefix(errMsg, "无权设置") {
			response.Error(c, errs.Forbidden(errMsg))
			return
		}
		response.Error(c, errs.Internal(errMsg, nil))
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除角色。
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteRoleRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.roleService.Delete(c.Request.Context(), actor, req.RoleIds)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "未获取到角色信息" {
			response.Error(c, errs.NotFound(errMsg, nil))
			return
		}
		if errMsg == "不能删除比自己角色等级高或相等的角色" {
			response.Error(c, errs.Forbidden(errMsg))
			return
		}
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
