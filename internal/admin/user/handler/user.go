package handler

import (
	"strconv"

	"gotribe/internal/admin/user/dto"
	"gotribe/internal/admin/user/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService service.Service
	cdnDomain   string
}

// NewHandler 创建用户处理器实例
func NewHandler(userService service.Service, cdnDomain string) *Handler {
	return &Handler{userService: userService, cdnDomain: cdnDomain}
}

// Detail 获取用户信息
// @Summary      获取用户信息
// @Description  根据用户ID获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user/{id} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("用户ID格式错误", nil))
		return
	}
	user, err := h.userService.Detail(ctx, int64(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	userResponse := dto.ToUserResponse(&user, h.cdnDomain)
	response.OK(c, gin.H{
		"user": userResponse,
	})
}

// List 获取用户列表
// @Summary      获取用户列表
// @Description  获取所有用户的列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request query dto.UserListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.UserListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	user, total, err := h.userService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"users": dto.ToUserListResponse(user, h.cdnDomain), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Create 创建用户
// @Summary      创建用户
// @Description  创建一个新的用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserRequest true "创建用户请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.userService.Create(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新用户
// @Summary      更新用户
// @Description  根据用户ID更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Param        request body dto.UpdateUserRequest true "更新用户请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user/{id} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("用户ID格式错误", nil))
		return
	}

	err = h.userService.Update(ctx, int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除用户
// @Summary      批量删除用户
// @Description  根据用户ID列表批量删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteUsersRequest true "删除用户请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteUsersRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.userService.Delete(ctx, req.Ids)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Search 搜索用户
// @Summary      搜索用户
// @Description  根据昵称搜索用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        nickname query string true "用户昵称"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user/search [get]
// @Security     BearerAuth
func (h *Handler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := h.userService.Search(ctx, c.Query("nickname"))
	if err != nil {
		response.Error(c, errs.Internal("获取需要更新的用户信息失败", err))
		return
	}
	response.OK(c, gin.H{"users": dto.ToUserListResponse(user, h.cdnDomain)})
}
