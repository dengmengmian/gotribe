package dto

import "gotribe/internal/model"

// CategoryResponse 分类响应。
type CategoryResponse struct {
	ID          int64  `json:"id"`
	ParentID    int64  `json:"parent_id"`
	Sort        uint   `json:"sort"`
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Path        string `json:"path"`
	Hidden      uint8  `json:"hidden"`
	Description string `json:"description,omitempty"`
	Status      uint8  `json:"status,omitempty"`
	Count       uint   `json:"count"`
}

// ToCategoryResponse 转换分类模型为响应。
func ToCategoryResponse(c *model.Category) CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		ParentID:    c.ParentID,
		Sort:        c.Sort,
		Icon:        c.Icon,
		Title:       c.Title,
		Slug:        c.Slug,
		Path:        c.Path,
		Hidden:      c.Hidden,
		Description: c.Description,
		Status:      c.Status,
		Count:       c.Count,
	}
}

// ToCategoryListResponse 批量转换。
func ToCategoryListResponse(categories []model.Category) []CategoryResponse {
	res := make([]CategoryResponse, 0, len(categories))
	for i := range categories {
		res = append(res, ToCategoryResponse(&categories[i]))
	}
	return res
}
