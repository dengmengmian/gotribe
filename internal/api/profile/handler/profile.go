// Package handler implements HTTP request handlers for the current user profile module.
package handler

// 本文件实现当前用户模块的 HTTP 处理逻辑。

import (
	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/middleware"
	profiledto "gotribe/internal/api/profile/dto"
	profileservice "gotribe/internal/api/profile/service"
	profileview "gotribe/internal/api/profile/view"
	"gotribe/internal/request"
	"gotribe/internal/core/response"
)

// Handler 负责处理当前用户资料相关的 HTTP 请求。
type Handler struct {
	service *profileservice.Service
}

// NewHandler 创建当前用户资料 HTTP 处理器实例。
func NewHandler(service *profileservice.Service) *Handler {
	return &Handler{service: service}
}

// GetMe 处理获取当前用户资料请求。
func (h *Handler) GetMe(c *gin.Context) {
	if currentUser, ok := middleware.GetCurrentUser(c); ok {
		response.OK(c, toMeResponse(*currentUser))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, errs.Unauthorized("missing user context"))
		return
	}
	data, err := h.service.GetMe(c.Request.Context(), middleware.GetProjectID(c), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, toMeResponse(*data))
}

// UpdateMe 处理更新当前用户资料请求。
func (h *Handler) UpdateMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, errs.Unauthorized("missing user context"))
		return
	}

	var req profiledto.UpdateProfileRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	data, err := h.service.UpdateMe(c.Request.Context(), middleware.GetProjectID(c), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, toMeResponse(*data))
}

// ChangePassword 处理修改当前用户密码请求。
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, errs.Unauthorized("missing user context"))
		return
	}

	var req profiledto.ChangePasswordRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), middleware.GetProjectID(c), userID, req); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// toMeResponse 将 profile 模块内部视图转换为 HTTP 响应结构。
func toMeResponse(view profileview.MeView) profiledto.MeResponse {
	return profiledto.MeResponse{
		ID:         view.ID,
		Username:   view.Username,
		ProjectID:  view.ProjectID,
		Nickname:   view.Nickname,
		Email:      view.Email,
		Phone:      view.Phone,
		Sex:        view.Sex,
		Status:     view.Status,
		Birthday:   view.Birthday,
		Background: view.Background,
		Ext:        view.Ext,
		AvatarURL:  view.AvatarURL,
		CreatedAt:  view.CreatedAt,
		UpdatedAt:  view.UpdatedAt,
	}
}
