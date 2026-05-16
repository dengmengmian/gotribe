package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
	projectdto "gotribe/internal/admin/project/dto"
	userdto "gotribe/internal/admin/user/dto"
)

type FeedbackResponse struct {
	ID        int            `json:"id"`
	ProjectID int64           `json:"project_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	UserID int64           `json:"user_id"`
	User      userdto.UserResponse    `json:"user"`
	Phone     string         `json:"phone"`
	Project   projectdto.ProjectResponse `json:"project"`
	CreatedAt string         `json:"created_at"`
}

func toFeedbackResponse(feedback model.Feedback, domain string) FeedbackResponse {
	dtoItem := FeedbackResponse{
		ID:        int(feedback.ID),
		ProjectID: feedback.ProjectID,
		Content:   feedback.Content,
		Title:     feedback.Title,
		UserID:    feedback.UserID,
		Phone:     feedback.Phone,
		CreatedAt: feedback.CreatedAt.Format(constant.TIME_FORMAT),
	}
	if feedback.User != nil {
		dtoItem.User = userdto.ToUserResponse(feedback.User, domain)
	}
	if feedback.Project != nil {
		dtoItem.Project = projectdto.ToProjectResponse(feedback.Project)
	}
	return dtoItem
}

// ToFeedbackInfoResponse 将单个 Feedback 转换为 FeedbackResponse
func ToFeedbackInfoResponse(feedBack model.Feedback, domain string) FeedbackResponse {
	return toFeedbackResponse(feedBack, domain)
}

// ToFeedbackListResponse 将多个 Feedback 转换为 []FeedbackResponse
func ToFeedbackListResponse(feedBackList []*model.Feedback, domain string) []FeedbackResponse {
	var feedBacks = make([]FeedbackResponse, len(feedBackList))
	for i, feedBack := range feedBackList {
		feedBacks[i] = toFeedbackResponse(*feedBack, domain)
	}
	return feedBacks
}
