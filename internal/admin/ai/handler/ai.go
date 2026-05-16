package handler

import (
	"gotribe/internal/admin/ai/dto"
	"gotribe/internal/admin/ai/service"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	aiService service.Service
}

// NewHandler 创建 AI 处理器。
func NewHandler(aiService service.Service) *Handler {
	return &Handler{aiService: aiService}
}

// Generate 执行通用 AI 生成任务。
func (h *Handler) Generate(c *gin.Context) {
	var req dto.GenerateRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.aiService.Generate(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}
