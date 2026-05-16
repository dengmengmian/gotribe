package repository

import (
	"context"
	"strings"

	"gotribe/internal/core/database"
	"gotribe/internal/model"
)

type Repository struct {
	tx *database.TransactionManager
}

func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

func (r *Repository) DetailByAlias(ctx context.Context, projectID int64, alias string) (model.Config, error) {
	var config model.Config
	err := r.tx.DB(ctx).
		Where("project_id = ? AND alias = ? AND status = ?", projectID, strings.TrimSpace(alias), 1).
		First(&config).Error
	return config, err
}
