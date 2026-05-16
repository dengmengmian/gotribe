package dto

import "gotribe/internal/model"

// TagResponse 标签响应。
type TagResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
	Count uint   `json:"count"`
}

// ToTagResponse 转换标签。
func ToTagResponse(tag *model.Tag) TagResponse {
	return TagResponse{
		ID:    tag.ID,
		Title: tag.Title,
		Slug:  tag.Slug,
		Color: tag.Color,
		Count: tag.Count,
	}
}

// ToTagListResponse 批量转换。
func ToTagListResponse(tags []model.Tag) []TagResponse {
	res := make([]TagResponse, 0, len(tags))
	for i := range tags {
		res = append(res, ToTagResponse(&tags[i]))
	}
	return res
}
