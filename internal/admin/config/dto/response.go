package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

type ConfigResponse struct {
	ID          int64            `json:"id"`
	Alias       string           `json:"alias"`
	Type        uint             `json:"type"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Info        string           `json:"info"`
	MDContent   string           `json:"md_content"`
	ProjectID   int64            `json:"project_id"`
	Project     *ProjectResponse `json:"project,omitempty"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

type ProjectResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
}

func ToConfigResponse(config model.Config) ConfigResponse {
	return ConfigResponse{
		ID:          config.ID,
		Alias:       config.Alias,
		Type:        config.Type,
		Info:        config.Info,
		Title:       config.Title,
		ProjectID:   config.ProjectID,
		Project:     toProjectResponse(config.Project),
		Description: config.Description,
		MDContent:   config.MDContent,
		CreatedAt:   config.CreatedAt.Format(constant.TIME_FORMAT),
		UpdatedAt:   config.UpdatedAt.Format(constant.TIME_FORMAT),
	}
}

func ToConfigListResponse(configList []*model.Config) []ConfigResponse {
	var configs []ConfigResponse
	for _, config := range configList {
		configDto := ConfigResponse{
			ID:          config.ID,
			Alias:       config.Alias,
			Type:        config.Type,
			Title:       config.Title,
			Description: config.Description,
			Info:        config.Info,
			ProjectID:   config.ProjectID,
			Project:     toProjectResponse(config.Project),
			MDContent:   config.MDContent,
			CreatedAt:   config.CreatedAt.Format(constant.TIME_FORMAT),
			UpdatedAt:   config.UpdatedAt.Format(constant.TIME_FORMAT),
		}

		configs = append(configs, configDto)
	}

	return configs
}

func toProjectResponse(project *model.Project) *ProjectResponse {
	if project == nil {
		return nil
	}
	return &ProjectResponse{
		ID:    project.ID,
		Title: project.Title,
		Name:  project.Name,
	}
}
