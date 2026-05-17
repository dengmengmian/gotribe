// Package repository implements data access and filtering logic for post queries.
package repository

// 本文件定义帖子仓储层使用的内部查询条件结构。

// ListFilter 表示帖子列表查询在数据访问层使用的过滤条件。
type ListFilter struct {
	Page        int
	PerPage     int
	Keyword     string
	Status      *int16
	Type        *int16
	DynamicType string
	TagIDs      []int64
	CategoryID  *int64
}
