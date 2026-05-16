package dto

import (
	"fmt"
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

// UserResponse 返回给前端的用户
type UserResponse struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Email     string  `json:"email"`
	AvatarURL string  `json:"avatar_url"`
	Sex       string  `json:"sex"`
	ProjectID int64   `json:"project_id"`
	Status    uint8   `json:"status"`
	Birthday  string  `json:"birthday"`
	Point     float64 `json:"point"`
	Phone     string  `json:"phone"`
	CreatedAt string  `json:"created_at"`
}

func toUserResponse(user *model.User, domain string) UserResponse {
	if user == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     stringValue(user.Email),
		Sex:       user.Sex,
		ProjectID: user.ProjectID,
		Birthday: func() string {
			if user.Birthday != nil {
				return user.Birthday.Format(constant.TIME_FORMAT)
			}
			return ""
		}(),
		AvatarURL: fmt.Sprintf("%s%s", domain, user.AvatarURL),
		CreatedAt: user.CreatedAt.Format(constant.TIME_FORMAT),
		Point:     user.Point,
		Status:    user.Status,
		Phone:     stringValue(user.Phone),
	}
}

// ToUserResponse 将单个 User 转换为 UserResponse
func ToUserResponse(user *model.User, domain string) UserResponse {
	return toUserResponse(user, domain)
}

// ToUserListResponse 将多个 User 转换为 UserResponse 列表
func ToUserListResponse(userList []*model.User, domain string) []UserResponse {
	var users []UserResponse
	for _, user := range userList {
		users = append(users, toUserResponse(user, domain))
	}
	return users
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
