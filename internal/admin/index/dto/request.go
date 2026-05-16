package dto

// IndexInfoRequest 首页信息查询参数
type IndexInfoRequest struct {
	ProjectID string `form:"project_id" json:"project_id"`
}

// TimeRangeRequest 时间范围查询参数
type TimeRangeRequest struct {
	ProjectID string `form:"project_id" json:"project_id"`
	TimeRange string `form:"time_range" json:"time_range"`
}
