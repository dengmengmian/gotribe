package dto

// CreateMenuRequest 创建菜单结构体
type CreateMenuRequest struct {
	Name       string `json:"name" form:"name" binding:"required,min=1,max=50"`
	Title      string `json:"title" form:"title" binding:"required,min=1,max=50"`
	Icon       string `json:"icon" form:"icon" binding:"min=0,max=50"`
	Path       string `json:"path" form:"path" binding:"required,min=1,max=100"`
	Redirect   string `json:"redirect" form:"redirect" binding:"min=0,max=100"`
	Component  string `json:"component" form:"component" binding:"required,min=1,max=100"`
	Sort       uint   `json:"sort" form:"sort" binding:"gte=1,lte=999"`
	Status     uint   `json:"status" form:"status" binding:"oneof=1 2"`
	Hidden     uint   `json:"hidden" form:"hidden" binding:"oneof=1 2"`
	NoCache    uint   `json:"no_cache" form:"no_cache" binding:"oneof=1 2"`
	AlwaysShow uint   `json:"always_show" form:"always_show" binding:"oneof=1 2"`
	Breadcrumb int64  `json:"breadcrumb" form:"breadcrumb" binding:"oneof=1 2"`
	ActiveMenu string `json:"active_menu" form:"active_menu" binding:"min=0,max=100"`
	ParentID   int64  `json:"parent_id" form:"parent_id"`
}

// UpdateMenuRequest 更新菜单结构体
type UpdateMenuRequest struct {
	Name       string `json:"name" form:"name" binding:"required,min=1,max=50"`
	Title      string `json:"title" form:"title" binding:"required,min=1,max=50"`
	Icon       string `json:"icon" form:"icon" binding:"min=0,max=50"`
	Path       string `json:"path" form:"path" binding:"required,min=1,max=100"`
	Redirect   string `json:"redirect" form:"redirect" binding:"min=0,max=100"`
	Component  string `json:"component" form:"component" binding:"min=0,max=100"`
	Sort       uint   `json:"sort" form:"sort" binding:"gte=1,lte=999"`
	Status     uint   `json:"status" form:"status" binding:"oneof=1 2"`
	Hidden     uint   `json:"hidden" form:"hidden" binding:"oneof=1 2"`
	NoCache    uint   `json:"no_cache" form:"no_cache" binding:"oneof=1 2"`
	AlwaysShow uint   `json:"always_show" form:"always_show" binding:"oneof=1 2"`
	Breadcrumb int64  `json:"breadcrumb" form:"breadcrumb" binding:"oneof=1 2"`
	ActiveMenu string `json:"active_menu" form:"active_menu" binding:"min=0,max=100"`
	ParentID   int64  `json:"parent_id" form:"parent_id"`
}

// DeleteMenuRequest 批量删除菜单结构体
type DeleteMenuRequest struct {
	MenuIds []int64 `json:"menu_ids" form:"menu_ids"`
}
