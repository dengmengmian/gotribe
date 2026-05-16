package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

type ProjectResponse struct {
	ID int64   `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CreatedAt      string `json:"created_at"`
	Name           string `json:"name"`
	Keywords       string `json:"keywords"`
	Domain         string `json:"domain"`
	PostURL        string `json:"post_url"`
	ICP            string `json:"icp"`
	Author         string `json:"author"`
	BaiduAnalytics string `json:"baidu_analytics"`
	Favicon        string `json:"favicon"`
	PublicSecurity string `json:"public_security"`
	NavImage       string `json:"nav_image"`
	Info           string `json:"info"`
	PushToken      string `json:"push_token"`
}

func ToProjectResponse(project *model.Project) ProjectResponse {
	return ProjectResponse{
		ID:             project.ID,
		Title:          project.Title,
		Description:    project.Description,
		CreatedAt:      project.CreatedAt.Format(constant.TIME_FORMAT),
		Name:           project.Name,
		Keywords:       project.Keywords,
		Domain:         project.Domain,
		PostURL:        project.PostURL,
		ICP:            project.ICP,
		Author:         project.Author,
		BaiduAnalytics: project.BaiduAnalytics,
		Favicon:        project.Favicon,
		PublicSecurity: project.PublicSecurity,
		NavImage:       project.NavImage,
		Info:           project.Info,
		PushToken:      project.PushToken,
	}
}

func ToProjectListResponse(projectList []*model.Project) []ProjectResponse {
	var projects []ProjectResponse
	for _, project := range projectList {
		projectDto := ProjectResponse{
			ID:             project.ID,
			Title:          project.Title,
			Description:    project.Description,
			CreatedAt:      project.CreatedAt.Format(constant.TIME_FORMAT),
			Name:           project.Name,
			Keywords:       project.Keywords,
			Domain:         project.Domain,
			PostURL:        project.PostURL,
			ICP:            project.ICP,
			Author:         project.Author,
			BaiduAnalytics: project.BaiduAnalytics,
			Favicon:        project.Favicon,
			PublicSecurity: project.PublicSecurity,
			NavImage:       project.NavImage,
			Info:           project.Info,
			PushToken:      project.PushToken,
		}

		projects = append(projects, projectDto)
	}

	return projects
}
