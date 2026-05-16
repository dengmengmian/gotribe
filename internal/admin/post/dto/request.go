package dto

// CreatePostRequest 创建内容结构体
type CreatePostRequest struct {
	Title       string   `form:"title" json:"title" binding:"required,min=2,max=60"`
	Slug        string   `form:"slug" json:"slug"`
	Description string   `form:"description" json:"description" binding:"required,min=2,max=300"`
	CategoryID  int64    `form:"category_id" json:"category_id" binding:"required"`
	ProjectID   int64    `form:"project_id" json:"project_id" binding:"required"`
	UserID      int64    `form:"user_id" json:"user_id" binding:"required"`
	Author      string   `form:"author" json:"author" binding:"required"`
	Content     string   `form:"content" json:"content" binding:"required"`
	HtmlContent string   `form:"html_content" json:"html_content" binding:"required"`
	ColumnID    int64    `form:"column_id" json:"column_id"`
	Tag         string   `form:"tag" json:"tag"`
	Ext         string   `form:"ext" json:"ext"`
	Icon        string   `form:"icon" json:"icon"`
	Type        uint     `form:"type" json:"type" binding:"required"`
	IsTop       int64    `form:"is_top" json:"is_top"`
	IsPasswd    int64    `form:"is_passwd" json:"is_passwd"`
	Password    string   `form:"password" json:"password"`
	Status      uint     `form:"status" json:"status"`
	Location    string   `form:"location" json:"location"`
	People      string   `form:"people" json:"people"`
	Time        string   `form:"time" json:"time"`
	Images      []string `form:"images" json:"images"`
	UnitPrice   float64  `form:"unit_price" json:"unit_price"`
	Video       string   `form:"video" json:"video"`
	ShowTime    string   `form:"show_time" json:"show_time"`
}

// UpdatePostRequest 更新内容结构体
type UpdatePostRequest struct {
	Title       string   `form:"title" json:"title" binding:"required,min=2,max=60"`
	Slug        string   `form:"slug" json:"slug"`
	Description string   `form:"description" json:"description" binding:"required,min=2,max=300"`
	CategoryID  int64    `form:"category_id" json:"category_id" binding:"required"`
	ProjectID   int64    `form:"project_id" json:"project_id" binding:"required"`
	UserID      int64    `form:"user_id" json:"user_id" binding:"required"`
	Author      string   `form:"author" json:"author" binding:"required"`
	Content     string   `form:"content" json:"content" binding:"required"`
	HtmlContent string   `form:"html_content" json:"html_content" binding:"required"`
	ColumnID    int64    `form:"column_id" json:"column_id"`
	Tag         string   `form:"tag" json:"tag"`
	Ext         string   `form:"ext" json:"ext"`
	Icon        string   `form:"icon" json:"icon"`
	Type        uint     `form:"type" json:"type" binding:"required"`
	IsTop       int64    `form:"is_top" json:"is_top"`
	IsPasswd    int64    `form:"is_passwd" json:"is_passwd"`
	Password    string   `form:"password" json:"password"`
	Status      uint     `form:"status" json:"status"`
	Location    string   `form:"location" json:"location"`
	People      string   `form:"people" json:"people"`
	Time        string   `form:"time" json:"time"`
	Images      []string `form:"images" json:"images"`
	UnitPrice   float64  `form:"unit_price" json:"unit_price"`
	ShowTime    string   `form:"show_time" json:"show_time"`
	Video       string   `form:"video" json:"video"`
}

// PostListRequest 获取内容列表结构体
type PostListRequest struct {
	ID        int64  `form:"id" json:"id"`
	Title     string `form:"title" json:"title"`
	Status    uint   `form:"status" json:"status"`
	ProjectID int64  `form:"project_id" json:"project_id"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

// DeletePostsRequest 批量删除内容结构体
type DeletePostsRequest struct {
	PostIds []int64 `json:"post_ids" form:"post_ids"`
}
