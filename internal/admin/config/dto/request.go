package dto

// CreateConfigRequest 创建配置结构体
type CreateConfigRequest struct {
	ProjectID   int64  `form:"project_id" json:"project_id" binding:"required"`
	Alias       string `form:"alias" json:"alias" binding:"required,min=2,max=20"`
	Type        uint   `form:"type" json:"type" binding:"required"`
	Title       string `form:"title" json:"title" binding:"required,min=2,max=20"`
	Description string `form:"description" json:"description" binding:"required,min=2,max=150"`
	Info        string `form:"info" json:"info" binding:"required,min=2,max=3000"`
	MDContent   string `form:"md_content" json:"md_content"`
}

// ConfigListRequest 获取配置列表结构体
type ConfigListRequest struct {
	ID        int64  `form:"id" json:"id"`
	Alias     string `form:"alias" json:"alias"`
	ProjectID int64  `form:"project_id" json:"project_id"`
	Title     string `form:"title" json:"title"`
	Type      uint   `form:"type" json:"type"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
}

// UpdateConfigRequest 更新配置内容
type UpdateConfigRequest struct {
	ProjectID   int64  `form:"project_id" json:"project_id" binding:"required"`
	Title       string `form:"title" json:"title" binding:"required,min=2,max=20"`
	Description string `form:"description" json:"description" binding:"required,min=2,max=150"`
	MDContent   string `form:"md_content" json:"md_content"`
	Info        string `form:"info" json:"info" binding:"required,min=2,max=3000"`
}

// DeleteConfigsRequest 批量删除项目结构体
type DeleteConfigsRequest struct {
	Ids []int64 `json:"ids" form:"ids" binding:"required"`
}
