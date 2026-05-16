package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

// ColumnResponse 返回给前端的专栏
type ColumnResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Info        string `json:"info"`
	ProjectID   int64  `json:"project_id"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func toColumnResponse(column *model.Column) ColumnResponse {
	if column == nil {
		return ColumnResponse{}
	}
	return ColumnResponse{
		ID:          column.ID,
		Title:       column.Title,
		Description: column.Description,
		Info:        column.Info,
		Icon:        column.Icon,
		ProjectID:   column.ProjectID,
		CreatedAt:   column.CreatedAt.Format(constant.TIME_FORMAT),
	}
}

// ToColumnResponse 将单个 Column 转换为 ColumnResponse
func ToColumnResponse(column model.Column) ColumnResponse {
	return toColumnResponse(&column)
}

// ToColumnListResponse 将多个 Column 转换为 ColumnResponse 列表
func ToColumnListResponse(columnList []*model.Column) []ColumnResponse {
	if columnList == nil {
		return []ColumnResponse{}
	}

	columns := make([]ColumnResponse, 0, len(columnList))
	for _, column := range columnList {
		columns = append(columns, toColumnResponse(column))
	}

	return columns
}
