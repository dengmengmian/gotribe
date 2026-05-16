package handler

import (
	"strconv"

	"gotribe/internal/admin/common"
	"gotribe/internal/admin/menu/dto"
	"gotribe/internal/admin/menu/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	menuService service.Service
}

// NewHandler 创建菜单处理器实例
func NewHandler(menuService service.Service) *Handler {
	return &Handler{menuService: menuService}
}

// List 获取菜单列表
// @Summary      获取菜单列表
// @Description  获取所有菜单的列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/list [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := h.menuService.List(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"menus": dto.ToMenuListResponse(menus)})
}

// Tree 获取菜单树
// @Summary      获取菜单树
// @Description  获取树形结构的菜单列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/tree [get]
// @Security     BearerAuth
func (h *Handler) Tree(c *gin.Context) {
	ctx := c.Request.Context()
	menu_tree, err := h.menuService.Tree(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"menu_tree": dto.ToMenuTreeResponse(menu_tree)})
}

// Create 创建菜单
// @Summary      创建菜单
// @Description  创建一个新的菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateMenuRequest true "创建菜单请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/create [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.menuService.Create(ctx, actor, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新菜单
// @Summary      更新菜单
// @Description  根据菜单ID更新菜单信息
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        menuID path string true "菜单ID"
// @Param        request body dto.UpdateMenuRequest true "更新菜单请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/update/{menuID} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateMenuRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 获取路径中的menuID
	menuID, _ := strconv.Atoi(c.Param("menuID"))
	if menuID <= 0 {
		response.Error(c, errs.BadRequest("菜单ID不正确", nil))
		return
	}

	ctx := c.Request.Context()
	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.menuService.Update(ctx, actor, int64(menuID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除菜单
// @Summary      批量删除菜单
// @Description  根据菜单ID列表批量删除菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteMenuRequest true "删除菜单请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/delete/batch [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteMenuRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.menuService.Delete(ctx, req.MenuIds)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// GetUserMenusByUserID 获取用户的可访问菜单列表
// @Summary      获取用户的可访问菜单列表
// @Description  根据用户ID获取用户的可访问菜单列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        userID path string true "用户ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/access/list/{userID} [get]
// @Security     BearerAuth
func (h *Handler) GetUserMenusByUserID(c *gin.Context) {
	// 获取路径中的userID
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID <= 0 {
		response.Error(c, errs.BadRequest("用户ID不正确", nil))
		return
	}

	ctx := c.Request.Context()
	menus, err := h.menuService.GetUserMenusByUserID(ctx, int64(userID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"menus": dto.ToMenuListResponse(menus)})
}

// GetUserMenuTreeByUserID 获取用户的可访问菜单树
// @Summary      获取用户的可访问菜单树
// @Description  根据用户ID获取用户的可访问菜单树
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        userID path string true "用户ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/access/tree/{userID} [get]
// @Security     BearerAuth
func (h *Handler) GetUserMenuTreeByUserID(c *gin.Context) {
	// 获取路径中的userID
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID <= 0 {
		response.Error(c, errs.BadRequest("用户ID不正确", nil))
		return
	}

	ctx := c.Request.Context()
	menu_tree, err := h.menuService.GetUserMenuTreeByUserID(ctx, int64(userID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"menu_tree": dto.ToMenuTreeResponse(menu_tree)})
}
