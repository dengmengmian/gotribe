package handler

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(public *gin.RouterGroup, apiLimiter gin.HandlerFunc) {
	public.GET("/configs/:alias", apiLimiter, h.DetailByAlias)
}
