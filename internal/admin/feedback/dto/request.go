package dto

// FeedbackListRequest 获取反馈列表结构体。
type FeedbackListRequest struct {
	ProjectID int64 `form:"project_id" json:"project_id"`
	PageNum   int64 `json:"page" form:"page"`
	PageSize  int64 `json:"per_page" form:"per_page"`
}

// DeleteFeedbackRequest 批量删除反馈请求。
type DeleteFeedbackRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
