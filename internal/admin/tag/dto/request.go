package dto

// CreateTagRequest 创建标签结构体
type CreateTagRequest struct {
	Title       string `form:"title" json:"title" binding:"required,min=2,max=20"`
	Slug        string `form:"slug" json:"slug" binding:"required,min=2,max=30"`
	Description string `form:"description" json:"description"`
	Color       string `form:"color" json:"color"`
	Sort        uint   `form:"sort" json:"sort"`
	Status      uint8  `form:"status" json:"status"`
}

// TagListRequest 获取标签列表结构体
type TagListRequest struct {
	ID        int64  `form:"id" json:"id"`
	Title     string `form:"title" json:"title"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

// DeleteTagsRequest 批量删除标签结构体
type DeleteTagsRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
