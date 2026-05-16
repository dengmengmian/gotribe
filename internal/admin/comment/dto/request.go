package dto

// CommentListRequest 获取评论列表请求
type CommentListRequest struct {
	ProjectID  int64  `form:"project_id" json:"project_id"`
	ObjectID   string `form:"object_id" json:"object_id"`
	ObjectType int64  `form:"object_type" json:"object_type"`
	Status     uint   `form:"status" json:"status"`
	Nickname   string `form:"nickname" json:"nickname"`
	PageNum    int64  `json:"page" form:"page"`
	PageSize   int64  `json:"per_page" form:"per_page"`
	SortBy     string `json:"sort_by" form:"sort_by"`
	SortOrder  string `json:"sort_order" form:"sort_order"`
}

// UpdateCommentRequest 更新评论请求。
type UpdateCommentRequest struct {
	Status uint `form:"status" json:"status" binding:"oneof=1 2"`
}

// DeleteCommentRequest 批量删除评论请求。
type DeleteCommentRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
