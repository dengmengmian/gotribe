package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

// CategoryResponse 返回给前端的分类
type CategoryResponse struct {
	ID          int64             `json:"id"`
	ParentID    int64             `json:"parent_id"`
	Sort        uint              `json:"sort"`
	Icon        string            `json:"icon"`
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	Path        string            `json:"path"`
	Hidden      uint8             `json:"hidden"`
	Description string            `json:"description,omitempty"`
	Ext         string            `json:"ext"`
	Status      uint8             `json:"status,omitempty"`
	Count       uint              `json:"count"`
	Children    []CategoryResponse `json:"children"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

func toCategoryResponse(category *model.Category) CategoryResponse {
	if category == nil {
		return CategoryResponse{}
	}
	children := make([]CategoryResponse, 0, len(category.Children))
	for _, child := range category.Children {
		children = append(children, toCategoryResponse(child))
	}
	return CategoryResponse{
		ID:          category.ID,
		ParentID:    category.ParentID,
		Sort:        category.Sort,
		Icon:        category.Icon,
		Title:       category.Title,
		Slug:        category.Slug,
		Path:        category.Path,
		Hidden:      category.Hidden,
		Description: category.Description,
		Ext:         category.Ext,
		Status:      category.Status,
		Count:       category.Count,
		Children:    children,
		CreatedAt:   category.CreatedAt.Format(constant.TIME_FORMAT),
		UpdatedAt:   category.UpdatedAt.Format(constant.TIME_FORMAT),
	}
}

// ToCategoryResponse 将单个 Category 转换为 CategoryResponse
func ToCategoryResponse(category model.Category) CategoryResponse {
	return toCategoryResponse(&category)
}

// ToCategoryListResponse 将多个 Category 转换为 CategoryResponse 列表
func ToCategoryListResponse(categoryList []*model.Category) []CategoryResponse {
	if categoryList == nil {
		return []CategoryResponse{}
	}

	categories := make([]CategoryResponse, 0, len(categoryList))
	for _, category := range categoryList {
		categories = append(categories, toCategoryResponse(category))
	}

	return categories
}
