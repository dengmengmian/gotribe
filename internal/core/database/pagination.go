package database

// 本文件提供分页参数规范化和分页元数据结构。

// Pagination 表示分页查询和响应所需的元数据。
type Pagination struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

// NormalizePagination 规范化分页参数，避免页码和页大小越界。
func NormalizePagination(page, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}
