package dto

// CreatePointLogRequest 创建积分结构体
type CreatePointLogRequest struct {
	ProjectID int64   `form:"project_id" json:"project_id" binding:"required"`
	UserID    int64   `form:"user_id" json:"user_id" binding:"required"`
	Point     float64 `form:"point" json:"point" binding:"required"`
}

// UpdatePointLogRequest 更新积分请求。
type UpdatePointLogRequest struct {
	ProjectID int64   `form:"project_id" json:"project_id"`
	UserID    int64   `form:"user_id" json:"user_id"`
	Point     float64 `form:"point" json:"point"`
}

// PointLogListRequest 获取积分列表结构体。
type PointLogListRequest struct {
	UserID    int64  `form:"user_id" json:"user_id"`
	Nickname  string `form:"nickname" json:"nickname"`
	ProjectID int64  `form:"project_id" json:"project_id"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
}

// DeletePointLogRequest 批量删除积分请求。
type DeletePointLogRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
