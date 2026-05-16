// Package dto defines request, response, and query structures for the example module.
package dto

// 本文件定义 example 模块的请求结构。

// CreateRequest 表示创建示例业务单的请求参数。
type CreateRequest struct {
	Name          string   `json:"name" binding:"required,max=100"`
	Description   string   `json:"description" binding:"omitempty,max=500"`
	Status        *int16   `json:"status" binding:"omitempty,oneof=0 1"`
	PrimaryPostID string   `json:"primary_post_id" binding:"required,max=64"`
	PostIDs       []string `json:"post_ids" binding:"required,min=1,max=20,dive,required,max=64"`
}

// UpdateRequest 表示更新示例业务单的请求参数。
type UpdateRequest struct {
	Name          *string   `json:"name"`
	Description   *string   `json:"description"`
	Status        *int16    `json:"status" binding:"omitempty,oneof=0 1"`
	PrimaryPostID *string   `json:"primary_post_id"`
	PostIDs       *[]string `json:"post_ids"`
}

// ListQuery 表示示例业务单列表接口支持的查询条件。
type ListQuery struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Keyword string `form:"keyword" binding:"omitempty,max=100"`
	Status  *int16 `form:"status" binding:"omitempty,oneof=0 1"`
}
