package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

type AdSceneResponse struct {
	ID int64   `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ProjectID int64   `json:"project_id"`
	ProjectTitle string `json:"project_title"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// toAdSceneResponse converts a model.AdScene to an AdSceneResponse.
func toAdSceneResponse(adScene model.AdScene) AdSceneResponse {
	var projectTitle string
	if adScene.Project != nil {
		projectTitle = adScene.Project.Title
	}
	return AdSceneResponse{
		ID:           adScene.ID,
		Title:        adScene.Title,
		Description:  adScene.Description,
		ProjectID:    adScene.ProjectID,
		CreatedAt:    adScene.CreatedAt.Format(constant.TIME_FORMAT),
		UpdatedAt:    adScene.UpdatedAt.Format(constant.TIME_FORMAT),
		ProjectTitle: projectTitle,
	}
}

// ToAdSceneInfoResponse converts a model.AdScene to an AdSceneResponse.
func ToAdSceneInfoResponse(adScene model.AdScene) AdSceneResponse {
	return toAdSceneResponse(adScene)
}

// ToAdSceneListResponse converts a list of model.AdScene to a list of AdSceneResponse.
func ToAdSceneListResponse(adSceneList []*model.AdScene) []AdSceneResponse {
	var adScenes []AdSceneResponse
	for _, adScene := range adSceneList {
		adScenes = append(adScenes, toAdSceneResponse(*adScene))
	}
	return adScenes
}
