package handler

import (
	"strconv"

	"gotribe/internal/admin/comment/dto"
	"gotribe/internal/admin/comment/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	commentService service.Service
}

// NewHandler 创建评论处理器实例
func NewHandler(commentService service.Service) *Handler {
	return &Handler{commentService: commentService}
}

// List 获取评论列表
// @Summary      获取评论列表
// @Description  获取所有评论的列表，支持分页和筛选
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        request query dto.CommentListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /comment [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.CommentListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	// 获取
	comment, total, err := h.commentService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"comments": dto.ToCommentListResponse(comment), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Update 更新评论
// @Summary      更新评论状态
// @Description  根据评论ID更新评论的审核状态
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        id path int true "评论ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /comment/{id} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateCommentRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	// 根据path中的ID获取评论信息
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	err = h.commentService.Update(ctx, int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除评论。
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteCommentRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	for _, id := range req.Ids {
		if err := h.commentService.Delete(c.Request.Context(), id); err != nil {
			response.Error(c, err)
			return
		}
	}
	response.OK(c, nil)
}
