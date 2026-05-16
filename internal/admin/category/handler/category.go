package handler

import (
	"strconv"

	"gotribe/internal/admin/category/dto"
	categorieservice "gotribe/internal/admin/category/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	categoryService categorieservice.Service
}

func NewHandler(categoryService categorieservice.Service) *Handler {
	return &Handler{categoryService: categoryService}
}

// Detail 获取当前分类信息
// @Summary      获取分类详情
// @Description  根据分类ID获取分类详细信息
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        id path int true "分类ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category/{id} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("分类ID格式错误", err))
		return
	}
	ctx := c.Request.Context()
	category, err := h.categoryService.Detail(ctx, int64(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{
		"category": category,
	})
}

// List 获取分类列表
// @Summary      获取分类列表
// @Description  获取所有分类的列表
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	categories, err := h.categoryService.List(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"categories": categories})
}

// Tree 获取分类树
// @Summary      获取分类树
// @Description  获取树形结构的分类列表
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category/tree [get]
// @Security     BearerAuth
func (h *Handler) Tree(c *gin.Context) {
	ctx := c.Request.Context()
	categoryTree, err := h.categoryService.Tree(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"category_tree": categoryTree})
}

// Create 创建分类
// @Summary      创建分类
// @Description  创建一个新的分类
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateCategoryRequest true "创建分类请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.categoryService.Create(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新分类
// @Summary      更新分类
// @Description  根据分类ID更新分类信息
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        categoryID path string true "分类ID"
// @Param        request body dto.UpdateCategoryRequest true "更新分类请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category/{categoryID} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateCategoryRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("分类ID格式错误", err))
		return
	}
	ctx := c.Request.Context()
	err = h.categoryService.Update(ctx, int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

// Delete 批量删除分类
// @Summary      批量删除分类
// @Description  根据分类ID列表批量删除分类
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteCategoryRequest true "删除分类请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteCategoryRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	ctx := c.Request.Context()
	err := h.categoryService.Delete(ctx, req.Ids)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}
