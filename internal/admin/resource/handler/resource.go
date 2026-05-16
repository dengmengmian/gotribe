package handler

import (
	"strconv"

	"gotribe/internal/admin/resource/dto"
	"gotribe/internal/admin/resource/service"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	resourceService service.Service
	cdnDomain       string
}

// NewHandler 创建资源处理器实例
func NewHandler(resourceService service.Service, cdnDomain string) *Handler {
	return &Handler{resourceService: resourceService, cdnDomain: cdnDomain}
}

// Detail 获取资源信息
// @Summary      获取资源信息
// @Description  根据资源ID获取资源详细信息
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        id path int true "资源ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /resource/{id} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("资源ID格式错误", nil))
		return
	}
	resource, err := h.resourceService.Detail(ctx, int64(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	resourceInfoDto := dto.ToResourceResponse(resource)
	response.OK(c, gin.H{
		"resource": resourceInfoDto,
	})
}

// List 获取资源列表
// @Summary      获取资源列表
// @Description  获取资源列表，支持分页和筛选
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        request query dto.ResourceListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /resource [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.ResourceListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	resource, total, err := h.resourceService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"resources": dto.ToResourceListResponse(resource), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Update 更新资源
// @Summary      更新资源
// @Description  根据资源ID更新资源信息
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        id path int true "资源ID"
// @Param        request body dto.CreateResourceRequest true "更新资源请求参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /resource/{id} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.CreateResourceRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("资源ID格式错误", nil))
		return
	}

	err = h.resourceService.Update(ctx, int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Upload 上传资源
// @Summary      上传资源
// @Description  上传文件资源到服务器
// @Tags         资源管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "上传的文件"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /resource/upload [post]
// @Security     BearerAuth
func (h *Handler) Upload(c *gin.Context) {
	ctx := c.Request.Context()
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errs.BadRequest("上传文件失败", err))
		return
	}
	if fileHeader.Size > constant.DEFAULT_UPLOAD_SIZE {
		response.Error(c, errs.BadRequest("上传资源过大", nil))
		return
	}

	resource, fileRes, err := h.resourceService.Upload(ctx, fileHeader)
	if err != nil {
		response.Error(c, errs.Internal("上传资源失败", err))
		return
	}

	uploadRes := dto.ToUploadResourceResponse(fileRes)
	uploadRes.Domain = h.cdnDomain
	uploadRes.FileType = int(resource.FileType)

	response.OK(c, gin.H{"upload": uploadRes})
}

// Delete 删除资源
// @Summary      删除资源
// @Description  根据资源ID列表删除资源
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteResourcesRequest true "删除资源请求参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /resource [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteResourcesRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := h.resourceService.Delete(ctx, req.Ids); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}
