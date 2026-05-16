package dto

// CreateResourceRequest 创建资源结构体
type CreateResourceRequest struct {
	Title       string `form:"title" json:"title" binding:"required,min=2,max=20"`
	Description string `form:"description" json:"description" binding:"required,min=2,max=150"`
}

// ResourceListRequest 获取资源列表结构体
type ResourceListRequest struct {
	ID       int64 `form:"id" json:"id"`
	Type     uint  `form:"type" json:"type"`
	PageNum  int64 `json:"page" form:"page"`
	PageSize int64 `json:"per_page" form:"per_page"`
}

// DeleteResourcesRequest 批量删除资源结构体
type DeleteResourcesRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
