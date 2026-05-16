package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

// CommentResponse 返回给前端的评论
type CommentResponse struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"project_id"`
	Status      uint   `json:"status"`
	UserID      int64  `json:"user_id"`
	ObjectID    string `json:"object_id"`
	ObjectType  int64  `json:"object_type"`
	Content     string `json:"comment"`
	HtmlContent string `json:"html_content"`
	Nickname    string `json:"nickname"`
	IP          string `json:"ip"`
	Country     string `json:"country"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// toCommentResponse converts a model.Comment to an CommentResponse.
func toCommentResponse(comment model.Comment) CommentResponse {
	var nickname string
	if comment.User != nil {
		nickname = comment.User.Nickname
	}
	return CommentResponse{
		ID:          comment.ID,
		ProjectID:   comment.ProjectID,
		UserID:      comment.UserID,
		ObjectID:    comment.ObjectID,
		ObjectType:  comment.ObjectType,
		Content:     comment.Content,
		HtmlContent: comment.HtmlContent,
		Status:      comment.Status,
		IP:          comment.IP,
		Country:     comment.Country,
		RegionName:  comment.RegionName,
		City:        comment.City,
		Nickname:    nickname,
		CreatedAt:   comment.CreatedAt.Format(constant.TIME_FORMAT),
		UpdatedAt:   comment.UpdatedAt.Format(constant.TIME_FORMAT),
	}
}

// ToCommentResponse converts a model.Comment to an CommentResponse.
func ToCommentResponse(comment model.Comment) CommentResponse {
	return toCommentResponse(comment)
}

// ToCommentListResponse converts a list of model.Comment to a list of CommentResponse.
func ToCommentListResponse(commentList []*model.Comment) []CommentResponse {
	var comments []CommentResponse
	for _, comment := range commentList {
		comments = append(comments, toCommentResponse(*comment))
	}
	return comments
}
