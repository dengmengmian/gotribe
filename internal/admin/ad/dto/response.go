package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

type AdResponse struct {
	ID int64   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SceneID int64   `json:"scene_id"`
	SceneTitle  string `json:"scene_title"`
	Status uint   `json:"status"`
	Image       string `json:"image"`
	Video       string `json:"video"`
	Sort uint   `json:"sort"`
	URL         string `json:"url"`
	URLType     int64   `json:"url_type"`
	Ext         string `json:"ext"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// toAdResponse converts a model.Ad to an AdResponse.
func toAdResponse(ad model.Ad) AdResponse {
	var sceneTitle string
	if ad.Scene != nil {
		sceneTitle = ad.Scene.Title
	}
	return AdResponse{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		SceneID:     ad.SceneID,
		Status:      ad.Status,
		Image:       ad.Image,
		Video:       ad.Video,
		Sort:        ad.Sort,
		URL:         ad.URL,
		URLType:     ad.URLType,
		Ext:         ad.Ext,
		SceneTitle:  sceneTitle,
		CreatedAt:   ad.CreatedAt.Format(constant.TIME_FORMAT),
		UpdatedAt:   ad.UpdatedAt.Format(constant.TIME_FORMAT),
	}
}

// ToAdInfoResponse converts a model.Ad to an AdResponse.
func ToAdInfoResponse(ad model.Ad) AdResponse {
	return toAdResponse(ad)
}

// ToAdListResponse converts a list of model.Ad to a list of AdResponse.
func ToAdListResponse(adList []*model.Ad) []AdResponse {
	var ads []AdResponse
	for _, ad := range adList {
		ads = append(ads, toAdResponse(*ad))
	}
	return ads
}
