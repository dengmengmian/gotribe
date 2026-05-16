package handler

// 本文件实现用户行为上报的 HTTP 处理逻辑。

import (
	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/middleware"
	"gotribe/internal/request"
	"gotribe/internal/core/response"
	usereventdto "gotribe/internal/api/user_event/dto"
	usereventservice "gotribe/internal/api/user_event/service"
)

// Handler 负责处理用户行为事件相关的 HTTP 请求。
type Handler struct {
	service *usereventservice.Service
}

// NewHandler 创建用户行为事件 HTTP 处理器实例。
func NewHandler(service *usereventservice.Service) *Handler {
	return &Handler{service: service}
}

// Create 处理用户行为事件上报请求。
func (h *Handler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, errs.Unauthorized("missing user context"))
		return
	}

	var req usereventdto.CreateRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.service.Create(
		c.Request.Context(),
		middleware.GetProjectID(c),
			int64(userID),
		req,
		c.ClientIP(),
		c.Request.UserAgent(),
	); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}
