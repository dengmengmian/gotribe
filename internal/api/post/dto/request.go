// Package dto defines request, response, and query structures for the post module.
package dto

// 本文件定义帖子模块的请求与查询参数结构。

// CreatePostRequest 预留文章创建接口的请求参数结构。
type CreatePostRequest struct{}

// ListQuery 表示文章列表接口支持的查询条件。
type ListQuery struct {
	Page        int    `form:"page" binding:"omitempty,min=1"`
	PerPage     int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Keyword     string `form:"keyword" binding:"omitempty,max=100"`
	Tag         string `form:"tag" binding:"omitempty,max=50"`
	Status      *int16 `form:"status" binding:"omitempty,oneof=0 1 2"`
	Type        *int16 `form:"type" binding:"omitempty"`
	DynamicType string `form:"dynamic_type" binding:"omitempty,max=50"`
	CategoryID  *int64 `form:"category_id" binding:"omitempty,min=1"`
}

// DetailQuery 表示文章详情接口支持的查询条件。
type DetailQuery struct {
	Password string `form:"password" binding:"omitempty,max=128"`
}
