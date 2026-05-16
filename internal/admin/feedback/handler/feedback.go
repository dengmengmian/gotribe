package handler

import (
	"gotribe/internal/admin/feedback/dto"
	"gotribe/internal/admin/feedback/service"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	feedbackService service.Service
	cdnDomain       string
}

// NewHandler 创建反馈处理器实例
func NewHandler(feedbackService service.Service, cdnDomain string) *Handler {
	return &Handler{feedbackService: feedbackService, cdnDomain: cdnDomain}
}

// List 获取反馈列表
// @Summary      获取反馈列表
// @Description  获取所有反馈的列表，支持分页和筛选
// @Tags         反馈管理
// @Accept       json
// @Produce      json
// @Param        request query dto.FeedbackListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /feedback [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.FeedbackListRequest
	// 参数绑定
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	// 获取
	feedbacks, total, err := h.feedbackService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.JSON(c, 200, gin.H{"feedbacks": dto.ToFeedbackListResponse(feedbacks, h.cdnDomain), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Delete 批量删除反馈。
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteFeedbackRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.feedbackService.Delete(c.Request.Context(), req.Ids); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
