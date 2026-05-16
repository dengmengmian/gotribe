package handler

import (
	"errors"
	"strconv"

	"gotribe/internal/admin/admin_user/dto"
	"gotribe/internal/admin/admin_user/service"
	"gotribe/internal/admin/common"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	adminService service.Service
}

func NewHandler(adminService service.Service) *Handler {
	return &Handler{adminService: adminService}
}

func (h *Handler) Me(c *gin.Context) {
	ctx := c.Request.Context()
	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	user, err := h.adminService.Me(ctx, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"admin": dto.ToAdminInfoResponse(user)})
}

func (h *Handler) List(c *gin.Context) {
	var req dto.AdminListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	ctx := c.Request.Context()
	users, total, err := h.adminService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200,
		gin.H{"admins": dto.ToAdminListResponse(users), "total": total},
		gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total},
	)
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	var req dto.ChangePwdRequest
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
	err = h.adminService.UpdatePassword(ctx, actor, req.OldPassword, req.NewPassword)
	if err != nil {
		var ae *errs.AppError
		if errors.As(err, &ae) && ae.Code == errs.CodeUnauthorized {
			response.Error(c, errs.Unauthorized("原密码有误"))
			return
		}
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateAdminRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if req.Password != "" && len(req.Password) < 6 {
		response.Error(c, errs.BadRequest("密码长度至少为6位", nil))
		return
	}
	ctx := c.Request.Context()
	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.adminService.Create(ctx, actor, &req)
	if err != nil {
		if err.Error() == "未获取到角色信息" || err.Error() == "用户不能创建比自己等级高的或者相同等级的用户" {
			response.Error(c, errs.Forbidden(err.Error()))
			return
		}
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

func (h *Handler) Update(c *gin.Context) {
	var req dto.CreateAdminRequest
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
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的用户ID", err))
		return
	}
	err = h.adminService.Update(ctx, actor, int64(userID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteAdminRequest
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
	err = h.adminService.Delete(ctx, actor, req.UserIds)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Detail 获取管理员详情。
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	admin, err := h.adminService.Detail(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"admin": admin})
}
