package dto

// CreateUserRequest 创建用户结构体
type CreateUserRequest struct {
	Username  string `form:"username" json:"username" binding:"required,min=2,max=20,alphanum"`
	Nickname  string `form:"nickname" json:"nickname" binding:"required,min=2,max=20"`
	AvatarURL string `form:"avatar_url" json:"avatar_url"`
	Email     string `form:"email" json:"email" binding:"omitempty,email,max=254"`
	Phone     string `form:"phone" json:"phone" binding:"omitempty,max=32"`
	ProjectID int64  `form:"project_id" json:"project_id" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

// UserListRequest 获取用户列表结构体
type UserListRequest struct {
	UserID    int64  `form:"user_id" json:"user_id"`
	ProjectID int64  `form:"project_id" json:"project_id"`
	Nickname  string `form:"nickname" json:"nickname"`
	Username  string `form:"username" json:"username"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
}

// DeleteUsersRequest 批量删除用户结构体
type DeleteUsersRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}

// UpdateUserRequest 更新用户结构体
type UpdateUserRequest struct {
	Nickname  string `form:"nickname" json:"nickname"`
	AvatarURL string `form:"avatar_url" json:"avatar_url"`
	Email     string `form:"email" json:"email" binding:"omitempty,email,max=254"`
	Phone     string `form:"phone" json:"phone" binding:"omitempty,max=32"`
	Password  string `form:"password" json:"password"`
}
