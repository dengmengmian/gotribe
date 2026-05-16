// Package dto defines request/response structures for the tag module.
package dto

// ListQuery 标签列表查询参数。
type ListQuery struct {
	Keyword string `form:"keyword" binding:"omitempty,max=50"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=100"`
}
