package dto

// CreateAdSceneRequest 创建推广位结构体
type CreateAdSceneRequest struct {
	ProjectID   int64  `form:"project_id" json:"project_id" binding:"required"`
	Title       string `form:"title" json:"title" binding:"required,min=2,max=50"`
	Description string `form:"description" json:"description" binding:"min=0,max=150"`
}

// AdSceneListRequest 获取推广位列表结构体
type AdSceneListRequest struct {
	ProjectID int64 `form:"project_id" json:"project_id"`
	PageNum   int64 `json:"page" form:"page"`
	PageSize  int64 `json:"per_page" form:"per_page"`
}

// UpdateAdSceneRequest 更新推广位内容
type UpdateAdSceneRequest struct {
	Title       string `form:"title" json:"title" binding:"required,min=2,max=50"`
	Description string `form:"description" json:"description" binding:"required,min=2,max=150"`
	ProjectID   int64  `form:"project_id" json:"project_id" binding:"required"`
}

// DeleteAdScenesRequest 批量删除项目结构体
type DeleteAdScenesRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
