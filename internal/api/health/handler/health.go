// Package handler implements HTTP handlers for health check probes.
package handler

// 本文件实现健康检查探针的 HTTP 处理器。

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gotribe/internal/api/health/service"
	"gotribe/internal/core/response"
)

// Handler 负责处理健康检查相关的 HTTP 请求。
type Handler struct {
	service *service.Service
}

// NewHandler 创建健康检查 HTTP 处理器实例。
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Liveness 处理 liveness 探针请求。
func (h *Handler) Liveness(c *gin.Context) {
	response.OK(c, h.service.Liveness())
}

// Version 处理版本信息请求。
func (h *Handler) Version(c *gin.Context) {
	response.OK(c, h.service.Version())
}

// Readiness 处理 readiness 探针请求。
func (h *Handler) Readiness(c *gin.Context) {
	result, ok := h.service.Readiness(c.Request.Context())
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"data":       result,
			"request_id": c.GetString("request_id"),
		})
		return
	}
	response.OK(c, result)
}
