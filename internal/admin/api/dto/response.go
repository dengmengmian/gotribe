package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

type ApiResponse struct {
	ID        int64  `json:"id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Category  string `json:"category"`
	Desc      string `json:"desc"`
	Creator   string `json:"creator"`
	CreatedAt string `json:"created_at"`
}

type ApiTreeResponse struct {
	ID       int          `json:"id"`
	Desc     string       `json:"desc"`
	Category string       `json:"category"`
	Children []*model.Api `json:"children"`
}

func toApiResponse(api *model.Api) ApiResponse {
	if api == nil {
		return ApiResponse{}
	}
	return ApiResponse{
		ID:        api.ID,
		Method:    api.Method,
		Path:      api.Path,
		Category:  api.Category,
		Desc:      api.Desc,
		Creator:   api.Creator,
		CreatedAt: api.CreatedAt.Format(constant.TIME_FORMAT),
	}
}

func ToApiResponse(api model.Api) ApiResponse {
	return toApiResponse(&api)
}

func ToApiListResponse(apiList []*model.Api) []ApiResponse {
	if apiList == nil {
		return []ApiResponse{}
	}

	apis := make([]ApiResponse, 0, len(apiList))
	for _, api := range apiList {
		apis = append(apis, toApiResponse(api))
	}

	return apis
}
