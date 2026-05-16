package handler

import (
	"gotribe/internal/admin/operation_log/dto"
	"gotribe/internal/admin/operation_log/service"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	operationLogService service.Service
}

// NewHandler 创建操作日志处理器实例
func NewHandler(operationLogService service.Service) *Handler {
	return &Handler{operationLogService: operationLogService}
}

// List 获取操作日志列表
// @Summary      获取操作日志列表
// @Description  获取所有操作日志的列表
// @Tags         操作日志管理
// @Accept       json
// @Produce      json
// @Param        request query dto.OperationLogListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /operation-log [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.OperationLogListRequest
	// 绑定参数
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	ctx := c.Request.Context()
	// 获取
	logs, total, err := h.operationLogService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"logs": dto.ToOperationLogListResponse(logs), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Delete 批量删除操作日志
// @Summary      批量删除操作日志
// @Description  根据操作日志ID列表批量删除操作日志
// @Tags         操作日志管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteOperationLogRequest true "操作日志ID列表"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /operation-log [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteOperationLogRequest
	// 参数绑定
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	// 删除接口
	err := h.operationLogService.Delete(ctx, req.OperationLogIds)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}
