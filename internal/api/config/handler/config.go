package handler

import (
	configdto "gotribe/internal/api/config/dto"
	configservice "gotribe/internal/api/config/service"
	"gotribe/internal/core/middleware"
	"gotribe/internal/core/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *configservice.Service
}

func NewHandler(service *configservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) DetailByAlias(c *gin.Context) {
	config, err := h.service.DetailByAlias(c.Request.Context(), middleware.GetProjectID(c), c.Param("alias"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"config": configdto.ToConfigResponse(config)})
}
