package handler

// 本文件实现 Admin 认证 HTTP 入口，依赖 internal/auth/core 提供的 audience-aware Manager。

import (
	authservice "gotribe/internal/auth/admin/service"
	"gotribe/internal/auth/admin/dto"
	"gotribe/internal/auth/core"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

// Handler 负责处理认证相关的 HTTP 请求
type Handler struct {
	audience    string
	authService authservice.Service
	manager     *core.Manager
}

// NewHandler 创建认证处理器实例。audience 通常传 core.AudienceAdmin。
func NewHandler(audience string, authService authservice.Service, manager *core.Manager) *Handler {
	return &Handler{audience: audience, authService: authService, manager: manager}
}

// Login 用户登录
// @Summary      用户登录
// @Description  管理员用户登录接口，返回JWT token
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterAndLoginRequest true "登录请求"
// @Success      200 {object} response.Response{data=object{token=string,expires=string}} "登录成功"
// @Failure      400 {object} response.Response "登录失败"
// @Router       /base/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.RegisterAndLoginRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, gin.H{
		"token":   result.Token,
		"expires": result.Expires.Format(constant.TIME_FORMAT),
	})
}

// Logout 用户登出
// @Summary      用户登出
// @Description  管理员用户登出接口
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response "登出成功"
// @Failure      400 {object} response.Response "登出失败"
// @Router       /base/logout [post]
// @Security     BearerAuth
func (h *Handler) Logout(c *gin.Context) {
	response.OK(c, nil)
}

// RefreshToken 刷新 token
// @Summary      刷新token
// @Description  刷新JWT token，延长登录状态
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response{data=object{token=string,expires=string}} "刷新成功"
// @Failure      400 {object} response.Response "刷新失败"
// @Router       /base/refreshToken [post]
// @Security     BearerAuth
func (h *Handler) RefreshToken(c *gin.Context) {
	bearer, err := core.ParseBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		response.Error(c, errs.Unauthorized("缺少或无效的认证令牌"))
		return
	}
	claims, err := h.manager.VerifyAccessTokenWithoutExpiry(h.audience, bearer)
	if err != nil {
		response.Error(c, errs.Unauthorized("缺少或无效的认证令牌"))
		return
	}

	result, err := h.authService.Refresh(c.Request.Context(), claims.UserID, claims.Username)
	if err != nil {
		response.Error(c, errs.Internal("刷新令牌失败，请稍后重试", err))
		return
	}

	response.OK(c, gin.H{
		"token":   result.Token,
		"expires": result.Expires.Format(constant.TIME_FORMAT),
	})
}
