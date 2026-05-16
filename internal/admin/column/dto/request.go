package dto

// CreateColumnRequest 创建专栏结构体
type CreateColumnRequest struct {
	Title       string `form:"title" json:"title" binding:"required,min=2,max=20"`
	Description string `form:"description" json:"description" binding:"required,min=2,max=300"`
	Info        string `form:"info" json:"info"`
	Icon        string `form:"icon" json:"icon" binding:"required,min=2,max=300"`
	ProjectID   int64  `form:"project_id" json:"project_id" binding:"required"`
}

// UpdateColumnRequest 更新专栏结构体
type UpdateColumnRequest struct {
	Title       string `form:"title" json:"title" binding:"required,min=2,max=20"`
	Description string `form:"description" json:"description" binding:"required,min=2,max=300"`
	Icon        string `form:"icon" json:"icon" binding:"required,min=2,max=300"`
	Info        string `form:"info" json:"info"`
	ProjectID   int64  `form:"project_id" json:"project_id" binding:"required"`
}

// ColumnListRequest 获取专栏列表结构体
type ColumnListRequest struct {
	ID        int64  `form:"id" json:"id"`
	ProjectID int64  `form:"project_id" json:"project_id"`
	Title     string `form:"title" json:"title"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
}

// DeleteColumnsRequest 批量删除专栏结构体
type DeleteColumnsRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
