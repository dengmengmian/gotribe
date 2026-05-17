package handler

// 本文件实现 TOTP 二次校验相关的 HTTP 入口：
//   - VerifyTOTP：登录后两步校验（公开，凭 step_token）
//   - Status / Bind / Confirm / Delete：登录态下的自助管理
//   - AdminResetTOTP：超管强制重置他人 TOTP

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gotribe/internal/admin/common"
	authservice "gotribe/internal/auth/admin/service"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"
)

// TOTPVerifyRequest 登录后两步校验请求体。
type TOTPVerifyRequest struct {
	StepToken string `json:"step_token" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// TOTPConfirmRequest 完成绑定请求体。
type TOTPConfirmRequest struct {
	Code string `json:"code" binding:"required"`
}

// TOTPDeleteRequest 自助解绑请求体。
type TOTPDeleteRequest struct {
	Code string `json:"code" binding:"required"`
}

// TOTPEnrollRequest 登录中途首次绑定的请求体（仅 step_token）。
type TOTPEnrollRequest struct {
	StepToken string `json:"step_token" binding:"required"`
}

// TOTPEnrollConfirmRequest 登录中途首次绑定的确认请求体。
type TOTPEnrollConfirmRequest struct {
	StepToken string `json:"step_token" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// VerifyTOTP 登录后两步校验，通过后签发 access_token。
// @Summary  TOTP 二次校验
// @Tags     认证管理
// @Accept   json
// @Produce  json
// @Param    request body TOTPVerifyRequest true "step_token + code"
// @Success  200 {object} response.Response{data=object{token=string,expires=string}}
// @Router   /base/totp/verify [post]
func (h *Handler) VerifyTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	var req TOTPVerifyRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.totpService.VerifyAndIssue(c.Request.Context(), req.StepToken, req.Code)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{
		"stage":   string(authservice.LoginStageOK),
		"token":   result.Token,
		"expires": result.Expires.Format(constant.TIME_FORMAT),
	})
}

// StatusTOTP 返回当前账户的 TOTP 状态。
// @Summary  查询 TOTP 绑定状态
// @Tags     认证管理
// @Produce  json
// @Success  200 {object} response.Response{data=authservice.TOTPStatus}
// @Router   /base/totp/status [get]
// @Security BearerAuth
func (h *Handler) StatusTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	admin, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	status, err := h.totpService.Status(c.Request.Context(), admin.ID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, status)
}

// BindTOTP 发起绑定流程，返回 secret + QR URL + 一次性备份码（仅本次返回）。
// @Summary  发起 TOTP 绑定
// @Tags     认证管理
// @Produce  json
// @Success  200 {object} response.Response{data=authservice.TOTPBindResult}
// @Router   /base/totp/bind [post]
// @Security BearerAuth
func (h *Handler) BindTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	admin, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.totpService.Bind(c.Request.Context(), &admin)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

// ConfirmTOTP 确认绑定（用一次 6 位码激活）。
// @Summary  确认 TOTP 绑定
// @Tags     认证管理
// @Accept   json
// @Produce  json
// @Param    request body TOTPConfirmRequest true "code"
// @Success  200 {object} response.Response
// @Router   /base/totp/confirm [post]
// @Security BearerAuth
func (h *Handler) ConfirmTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	admin, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req TOTPConfirmRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.totpService.Confirm(c.Request.Context(), admin.ID, req.Code); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// DeleteTOTP 自助解绑，需提供当前 6 位码。
// @Summary  自助解绑 TOTP
// @Tags     认证管理
// @Accept   json
// @Produce  json
// @Param    request body TOTPDeleteRequest true "code"
// @Success  200 {object} response.Response
// @Router   /base/totp [delete]
// @Security BearerAuth
func (h *Handler) DeleteTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	admin, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req TOTPDeleteRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.totpService.Delete(c.Request.Context(), &admin, req.Code); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// EnrollTOTP 登录后台强制绑定首步：用 step_token(purpose=totp_bind) 换取 secret/QR/备份码。
// @Summary  登录中途首次绑定 TOTP（公开，凭 step_token）
// @Tags     认证管理
// @Accept   json
// @Produce  json
// @Param    request body TOTPEnrollRequest true "step_token"
// @Success  200 {object} response.Response{data=authservice.TOTPBindResult}
// @Router   /base/totp/enroll [post]
func (h *Handler) EnrollTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	var req TOTPEnrollRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.totpService.EnrollPending(c.Request.Context(), req.StepToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

// ConfirmEnrollTOTP 登录中途首次绑定的最终确认：用同一 step_token + 6 位码激活并签发 access_token。
// @Summary  确认登录中途绑定并签发 access_token
// @Tags     认证管理
// @Accept   json
// @Produce  json
// @Param    request body TOTPEnrollConfirmRequest true "step_token + code"
// @Success  200 {object} response.Response{data=object{token=string,expires=string}}
// @Router   /base/totp/enroll/confirm [post]
func (h *Handler) ConfirmEnrollTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	var req TOTPEnrollConfirmRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.totpService.ConfirmEnrollPending(c.Request.Context(), req.StepToken, req.Code)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{
		"stage":   string(authservice.LoginStageOK),
		"token":   result.Token,
		"expires": result.Expires.Format(constant.TIME_FORMAT),
	})
}

// AdminResetTOTP 超管强制重置他人 TOTP。Casbin 已对 super admin 全部放行；
// 非 super admin 需通过 Casbin 单独授权。
// @Summary  超管重置他人 TOTP
// @Tags     认证管理
// @Param    id path int true "目标管理员 ID"
// @Success  200 {object} response.Response
// @Router   /admin/{id}/totp/reset [post]
// @Security BearerAuth
func (h *Handler) AdminResetTOTP(c *gin.Context) {
	if h.totpService == nil {
		response.Error(c, errs.Internal("TOTP 服务未启用", nil))
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, errs.BadRequest("无效的管理员 ID", err))
		return
	}
	if err := h.totpService.AdminReset(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
