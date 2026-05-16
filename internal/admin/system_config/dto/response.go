package dto

import (
	"gotribe/internal/model"
)

type SystemConfigResponse struct {
	SystemConfigID string `json:"system_config_id"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Logo           string `json:"logo"`
	Icon           string `json:"icon"`
	Footer         string `json:"footer"`
}

func ToSystemConfigInfoResponse(systemConfig *model.SystemConfig) SystemConfigResponse {
	return SystemConfigResponse{
		SystemConfigID: systemConfig.SystemConfigID,
		Title:          systemConfig.Title,
		Content:        systemConfig.Content,
		Logo:           systemConfig.Logo,
		Icon:           systemConfig.Icon,
		Footer:         systemConfig.Footer,
	}
}
