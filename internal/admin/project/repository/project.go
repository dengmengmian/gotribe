package repository

import (
	"context"
	"fmt"
	"strings"

	"gotribe/internal/admin/project/dto"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建项目仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// GetProjectByID 获取单个项目
func (r *Repository) GetProjectByID(ctx context.Context, id int64) (model.Project, error) {
	var project model.Project
	err := r.tx.DB(ctx).Where("id = ?", id).First(&project).Error
	return project, err
}

// List 获取项目列表
func (r *Repository) List(ctx context.Context, req *dto.ProjectListRequest) ([]*model.Project, int64, error) {
	var list []*model.Project
	db := r.tx.DB(ctx).Model(&model.Project{})

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	if req.ID != 0 {
		db = db.Where("id = ?", req.ID)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order("created_at DESC")

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	return list, total, err
}

// Create 创建项目
func (r *Repository) Create(ctx context.Context, project *model.Project) error {
	err := r.tx.DB(ctx).Create(project).Error
	return err
}

// UpdateProject 更新项目
func (r *Repository) UpdateProject(ctx context.Context, project *model.Project) error {
	err := r.tx.DB(ctx).Model(project).Updates(project).Error
	if err != nil {
		return err
	}

	return err
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var projects []model.Project
	for _, id := range ids {
		// 根据ID获取项目
		project, err := r.GetProjectByID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的项目", id)
		}
		projects = append(projects, project)
	}

	err := r.tx.DB(ctx).Delete(&projects).Error

	return err
}

// GetProjectsBySitemap 获取sitmap所需要的 projects 信息
func (r *Repository) GetProjectsBySitemap(ctx context.Context) ([]*model.Project, error) {
	var list []*model.Project
	err := r.tx.DB(ctx).Model(&model.Project{}).Order("created_at DESC").Find(&list).Error
	return list, err
}
