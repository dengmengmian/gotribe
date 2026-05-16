// Package handler implements HTTP request handlers for the authentication module.
package handler

// 本文件实现认证模块的 HTTP 处理逻辑。

import (
	"context"

	"github.com/gin-gonic/gin"
	"gotribe/internal/auth/user/dto"
	authservice "gotribe/internal/auth/user/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/middleware"
	"gotribe/internal/request"
	"gotribe/internal/core/response"
)

// authService 定义认证服务接口，便于测试时 mock。
type authService interface {
	Login(ctx context.Context, projectID string, req dto.LoginRequest) (*dto.AuthResponse, error)
	Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, currentUserID int64, req dto.LogoutRequest) error
}

// Handler 负责处理认证相关的 HTTP 请求。
type Handler struct {
	service authService
}

// NewHandler 创建认证 HTTP 处理器实例。
func NewHandler(svc *authservice.Service) *Handler {
	return &Handler{service: svc}
}

// Login 处理登录请求。
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	data, err := h.service.Login(c.Request.Context(), middleware.GetProjectID(c), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

// Refresh 处理刷新令牌请求。
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	data, err := h.service.Refresh(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

// Logout 处理退出登录请求。
func (h *Handler) Logout(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, errs.Unauthorized("missing user context"))
		return
	}
	var req dto.LogoutRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.Logout(c.Request.Context(), userID, req); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}
