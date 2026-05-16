package dto

// CreateAdminRequest 创建用户结构体
type CreateAdminRequest struct {
	Username     string  `form:"username" json:"username" binding:"required,min=2,max=20"`
	Password     string  `form:"password" json:"password"`
	Mobile       string  `form:"mobile" json:"mobile" binding:"required,checkMobile"`
	Avatar       string  `form:"avatar" json:"avatar"`
	Nickname     string  `form:"nickname" json:"nickname" binding:"min=0,max=20"`
	Introduction string  `form:"introduction" json:"introduction" binding:"min=0,max=255"`
	Status       uint    `form:"status" json:"status" binding:"oneof=1 2"`
	RoleIds      []int64 `form:"role_ids" json:"role_ids" binding:"required"`
}

// AdminListRequest 获取用户列表结构体
type AdminListRequest struct {
	Username  string `json:"username" form:"username"`
	Mobile    string `json:"mobile" form:"mobile"`
	Nickname  string `json:"nickname" form:"nickname"`
	Status    uint   `json:"status" form:"status"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

// DeleteAdminRequest 批量删除用户结构体
type DeleteAdminRequest struct {
	UserIds []int64 `json:"user_ids" form:"user_ids"`
}

// ChangePwdRequest 更新密码结构体
type ChangePwdRequest struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
}
