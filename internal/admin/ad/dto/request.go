package dto

// CreateAdRequest 创建推广位结构体
type CreateAdRequest struct {
	Title       string `form:"title" json:"title" binding:"required,min=2,max=50"`
	Description string `form:"description" json:"description" binding:"min=0,max=150"`
	URL         string `form:"url" json:"url" binding:"required,min=2,max=255"`
	URLType     int64  `form:"url_type" json:"url_type" binding:"required"`
	Image       string `form:"image" json:"image"`
	Video       string `form:"video" json:"video"`
	Sort        uint   `form:"sort" json:"sort" binding:"required"`
	Status      uint   `form:"status" json:"status" binding:"oneof=1 2"`
	SceneID     int64  `form:"scene_id" json:"scene_id" binding:"required"`
	Ext         string `form:"ext" json:"ext"`
}

// AdListRequest 获取推广位列表结构体
type AdListRequest struct {
	SceneID   int64  `form:"scene_id" json:"scene_id"`
	Title     string `form:"title" json:"title"`
	Status    uint   `form:"status" json:"status"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

// UpdateAdRequest 更新推广位内容
type UpdateAdRequest struct {
	Title       string `form:"title" json:"title" binding:"required,min=2,max=50"`
	Description string `form:"description" json:"description" binding:"min=0,max=150"`
	URL         string `form:"url" json:"url" binding:"required,min=2,max=255"`
	URLType     int64  `form:"url_type" json:"url_type" binding:"required"`
	Image       string `form:"image" json:"image"`
	Video       string `form:"video" json:"video"`
	Sort        uint   `form:"sort" json:"sort" binding:"required"`
	Status      uint   `form:"status" json:"status" binding:"oneof=1 2"`
	SceneID     int64  `form:"scene_id" json:"scene_id" binding:"required"`
	Ext         string `form:"ext" json:"ext"`
}

// DeleteAdsRequest 批量删除项目结构体
type DeleteAdsRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
