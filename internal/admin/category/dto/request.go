package dto

// CreateCategoryRequest 创建分类结构体
type CreateCategoryRequest struct {
	Title       string `json:"title" form:"title" binding:"required,min=1,max=30"`
	Slug        string `json:"slug" form:"slug" binding:"required,min=1,max=30"`
	Icon        string `json:"icon" form:"icon"`
	Path        string `json:"path" form:"path"`
	Sort        uint   `json:"sort" form:"sort" binding:"gte=1,lte=999"`
	Status      uint8  `json:"status" form:"status"`
	Hidden      uint8  `json:"hidden" form:"hidden" binding:"oneof=1 2"`
	Description string `json:"description" form:"description"`
	ParentID    int64  `json:"parent_id" form:"parent_id"`
}

// UpdateCategoryRequest 更新分类结构体
type UpdateCategoryRequest struct {
	Title       string `json:"title" form:"title" binding:"required,min=1,max=30"`
	Slug        string `json:"slug" form:"slug" binding:"required,min=1,max=30"`
	Icon        string `json:"icon" form:"icon"`
	Path        string `json:"path" form:"path"`
	Sort        uint   `json:"sort" form:"sort" binding:"gte=1,lte=999"`
	Status      uint8  `json:"status" form:"status"`
	Hidden      uint8  `json:"hidden" form:"hidden" binding:"oneof=1 2"`
	ParentID    int64  `json:"parent_id" form:"parent_id"`
	Description string `json:"description" form:"description"`
}

// DeleteCategoryRequest 删除分类结构体
type DeleteCategoryRequest struct {
	Ids []int64 `json:"ids" form:"ids"`
}
