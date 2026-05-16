package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

// TagResponse 返回给前端的标签
type TagResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Sort        uint   `json:"sort"`
	Count       uint   `json:"count"`
	Status      uint8  `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func toTagResponse(tag *model.Tag) TagResponse {
	if tag == nil {
		return TagResponse{}
	}
	return TagResponse{
		ID:          tag.ID,
		Title:       tag.Title,
		Slug:        tag.Slug,
		Color:       tag.Color,
		Description: tag.Description,
		Sort:        tag.Sort,
		Count:       tag.Count,
		Status:      tag.Status,
		CreatedAt:   tag.CreatedAt.Format(constant.TIME_FORMAT),
	}
}

// ToTagResponse 将单个 Tag 转换为 TagResponse
func ToTagResponse(tag model.Tag) TagResponse {
	return toTagResponse(&tag)
}

// ToTagListResponse 将多个 Tag 转换为 TagResponse 列表
func ToTagListResponse(tagList []*model.Tag) []TagResponse {
	if tagList == nil {
		return []TagResponse{}
	}

	tags := make([]TagResponse, 0, len(tagList))
	for _, tag := range tagList {
		tags = append(tags, toTagResponse(tag))
	}

	return tags
}
