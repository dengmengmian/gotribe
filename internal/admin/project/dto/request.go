package dto

// CreateProjectRequest 创建项目结构体
type CreateProjectRequest struct {
	Name           string `form:"name" json:"name" binding:"required,min=2,max=20"`
	Title          string `form:"title" json:"title"`
	Description    string `form:"description" json:"description"`
	Keywords       string `form:"keywords" json:"keywords"`
	Domain         string `form:"domain" json:"domain"`
	PostURL        string `form:"post_url" json:"post_url"`
	ICP            string `form:"icp" json:"icp"`
	BaiduAnalytics string `form:"baidu_analytics" json:"baidu_analytics"`
	Favicon        string `form:"favicon" json:"favicon"`
	PublicSecurity string `form:"public_security" json:"public_security"`
	Author         string `form:"author" json:"author"`
	NavImage       string `form:"nav_image" json:"nav_image"`
	Info           string `form:"info" json:"info"`
	PushToken      string `form:"push_token" json:"push_token"`
}

// ProjectListRequest 获取项目列表结构体
type ProjectListRequest struct {
	ID       int64  `form:"id" json:"id"`
	Title    string `form:"title" json:"title"`
	PageNum  int64  `json:"page" form:"page"`
	PageSize int64  `json:"per_page" form:"per_page"`
}

// DeleteProjectsRequest 批量删除项目结构体
type DeleteProjectsRequest struct {
	ProjectIDs []int64 `json:"project_ids" form:"project_ids"`
}
