package dto

// 本文件定义认证模块的响应结构。

// UserSummary 表示登录响应中返回的用户摘要信息。
type UserSummary struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
	ProjectID string `json:"project_id"`
}

// AuthResponse 表示登录或刷新令牌接口的响应数据。
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	User         UserSummary `json:"user"`
}
