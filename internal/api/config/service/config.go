package service

import (
	"context"
	"strconv"

	"gotribe/internal/api/config/repository"
	"gotribe/internal/model"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) DetailByAlias(ctx context.Context, projectIDRaw string, alias string) (model.Config, error) {
	projectID, err := strconv.ParseInt(projectIDRaw, 10, 64)
	if err != nil {
		return model.Config{}, err
	}
	return s.repo.DetailByAlias(ctx, projectID, alias)
}
