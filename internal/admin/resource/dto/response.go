package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
	"gotribe/internal/admin/util/upload"
)

// ResourceResponse 返回给前端的资源
type ResourceResponse struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	URL           string `json:"url"`
	Path          string `json:"path"`
	FileType      int64  `json:"file_type"`
	FileExtension string `json:"file_extension"`
	Size          int64  `json:"size"`
	CreatedAt     string `json:"created_at"`
}

// ToResourceResponse 将单个 Resource 转换为 ResourceResponse
func ToResourceResponse(resource model.Resource) ResourceResponse {
	return ResourceResponse{
		ID:            resource.ID,
		Title:         resource.Title,
		Description:   resource.Description,
		URL:           resource.URL,
		Path:          resource.Path,
		FileType:      resource.FileType,
		FileExtension: resource.FileExtension,
		Size:          resource.Size,
		CreatedAt:     resource.CreatedAt.Format(constant.TIME_FORMAT),
	}
}

// ToResourceListResponse 将多个 Resource 转换为 ResourceResponse 列表
func ToResourceListResponse(resourceList []*model.Resource) []ResourceResponse {
	var resources []ResourceResponse
	for _, resource := range resourceList {
		resourceResponse := ResourceResponse{
			ID:            resource.ID,
			Title:         resource.Title,
			Description:   resource.Description,
			URL:           resource.URL,
			Path:          resource.Path,
			FileType:      resource.FileType,
			FileExtension: resource.FileExtension,
			Size:          resource.Size,
			CreatedAt:     resource.CreatedAt.Format(constant.TIME_FORMAT),
		}

		resources = append(resources, resourceResponse)
	}

	return resources
}

// UploadResourceResponse 上传资源返回结构
type UploadResourceResponse struct {
	FileExt  string `json:"file_ext"`
	Key      string `json:"key"`
	Domain   string `json:"domain"`
	FileType int    `json:"file_type"`
}

// ToUploadResourceResponse 将 UploadResource 转换为 UploadResourceResponse
func ToUploadResourceResponse(resource *upload.UploadResource) UploadResourceResponse {
	return UploadResourceResponse{
		FileExt:  resource.FileExt,
		Key:      resource.Key,
		Domain:   "",
		FileType: 0,
	}
}
