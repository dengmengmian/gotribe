package dto

// 本文件定义个人资料模块的响应结构与跨模块共享的用户快照。

// MeResponse 表示当前用户资料接口的响应数据。
type MeResponse struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	ProjectID  string `json:"project_id"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Sex        string `json:"sex"`
	Status     int16  `json:"status"`
	Birthday   string `json:"birthday"`
	Background string `json:"background"`
	Ext        string `json:"ext"`
	AvatarURL  string `json:"avatar_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// UserSnapshot 表示供其他模块复用的用户快照信息。
type UserSnapshot struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}
