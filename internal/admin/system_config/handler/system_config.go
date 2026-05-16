package handler

import (
	"gotribe/internal/admin/system_config/dto"
	"gotribe/internal/admin/system_config/service"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	systemConfigService service.Service
}

// NewHandler 创建系统配置处理器实例
func NewHandler(systemConfigService service.Service) *Handler {
	return &Handler{systemConfigService: systemConfigService}
}

// Detail 获取系统配置信息
// @Summary      获取系统配置信息
// @Description  获取当前系统配置信息
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /base/config [get]
func (h *Handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	systemConfig, err := h.systemConfigService.Detail(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	systemConfigInfoDto := dto.ToSystemConfigInfoResponse(&systemConfig)
	response.OK(c, gin.H{
		"systemConfig": systemConfigInfoDto,
	})
}

// Update 更新系统配置
// @Summary      更新系统配置
// @Description  更新系统配置信息
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSystemConfigRequest true "更新系统配置请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /systemConfig/update [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateSystemConfigRequest
	// 参数绑定
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 更新系统配置
	err := h.systemConfigService.Update(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
