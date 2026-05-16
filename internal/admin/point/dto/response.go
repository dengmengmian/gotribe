package dto

import (
	"gotribe/internal/model"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/util"
)

type PointResponse struct {
	ID        int64   `json:"id"`
	Point     float64 `json:"point"`
	UserID int64    `json:"user_id"`
	Reason    string  `json:"reason"`
	Nickname  string  `json:"nickname"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// toPointResponse converts a model.PointLog to a PointResponse.
func toPointResponse(point model.PointLog) PointResponse {
	var nickname string
	if point.User != nil {
		nickname = point.User.Nickname
	}
	return PointResponse{
		ID:        int64(point.ID),
		Point:     utils.MoneyUtil.CentsToYuan(point.Points),
		UserID:    point.UserID,
		Nickname:  nickname,
		Reason:    point.Reason,
		CreatedAt: point.CreatedAt.Format(constant.TIME_FORMAT),
		UpdatedAt: point.UpdatedAt.Format(constant.TIME_FORMAT),
	}
}

// ToPointInfoResponse converts a model.PointLog to a PointResponse.
func ToPointInfoResponse(point model.PointLog) PointResponse {
	return toPointResponse(point)
}

// ToPointListResponse converts a list of model.PointLog to a list of PointResponse.
func ToPointListResponse(pointList []*model.PointLog) []PointResponse {
	var points []PointResponse
	for _, point := range pointList {
		points = append(points, toPointResponse(*point))
	}
	return points
}
