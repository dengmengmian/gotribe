package dto

import "gotribe/internal/model"

type ConfigResponse struct {
	ID          int64  `json:"id"`
	Alias       string `json:"alias"`
	Type        uint   `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Info        string `json:"info"`
	MDContent   string `json:"md_content"`
	ProjectID   int64  `json:"project_id"`
}

func ToConfigResponse(config model.Config) ConfigResponse {
	return ConfigResponse{
		ID:          config.ID,
		Alias:       config.Alias,
		Type:        config.Type,
		Title:       config.Title,
		Description: config.Description,
		Info:        config.Info,
		MDContent:   config.MDContent,
		ProjectID:   config.ProjectID,
	}
}
