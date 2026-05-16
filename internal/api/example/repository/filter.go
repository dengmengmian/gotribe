package repository

// 本文件定义 example 仓储层使用的内部查询条件。

// ListFilter 表示示例业务单列表查询的过滤条件。
type ListFilter struct {
	Page    int
	PerPage int
	Keyword string
	Status  *int16
}
