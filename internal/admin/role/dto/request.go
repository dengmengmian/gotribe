package dto

// CreateRoleRequest 新增角色结构体
type CreateRoleRequest struct {
	Name    string `json:"name" form:"name" binding:"required,min=1,max=20"`
	Keyword string `json:"keyword" form:"keyword" binding:"required,min=1,max=20"`
	Desc    string `json:"desc" form:"desc" binding:"min=0,max=100"`
	Status  uint   `json:"status" form:"status" binding:"oneof=1 2"`
	Sort    uint   `json:"sort" form:"sort" binding:"gte=1,lte=999"`
}

// RoleListRequest 获取用户角色结构体
type RoleListRequest struct {
	Name      string `json:"name" form:"name"`
	Keyword   string `json:"keyword" form:"keyword"`
	Status    uint   `json:"status" form:"status"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

// DeleteRoleRequest 批量删除角色结构体
type DeleteRoleRequest struct {
	RoleIds []int64 `json:"role_ids" form:"role_ids"`
}

// UpdateRoleMenusRequest 更新角色的权限菜单
type UpdateRoleMenusRequest struct {
	MenuIds []int64 `json:"menu_ids" form:"menu_ids"`
}

// UpdateRoleApisRequest 更新角色的权限接口
type UpdateRoleApisRequest struct {
	ApiIds []int64 `json:"api_ids" form:"api_ids"`
}
