package dto

// CreateSystemConfigRequest 创建系统配置结构体
type CreateSystemConfigRequest struct {
	Title   string `form:"title" json:"title" binding:"required,min=2,max=20"`
	Content string `form:"content" json:"content"`
	Logo    string `form:"logo" json:"logo" binding:"required"`
	Icon    string `form:"icon" json:"icon"`
	Footer  string `form:"footer" json:"footer"`
}
