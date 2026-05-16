package repository

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gotribe/internal/core/database"

	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

type Repository struct {
	tx  *database.TransactionManager
	log *zap.SugaredLogger
}

func NewRepository(tx *database.TransactionManager, log *zap.SugaredLogger) *Repository {
	return &Repository{tx: tx, log: log}
}

// 获取单个分类详情
func (r *Repository) Detail(ctx context.Context, id int64) (model.Category, error) {
	var category model.Category
	err := r.tx.DB(ctx).Where("id = ?", id).First(&category).Error
	return category, err
}

// 获取分类列表
func (r *Repository) List(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.tx.DB(ctx).Order("sort").Find(&categories).Error
	return categories, err
}

// 获取分类树
func (r *Repository) Tree(ctx context.Context) ([]*model.Category, error) {
	var categorys []*model.Category
	err := r.tx.DB(ctx).Order("sort").Find(&categorys).Error
	return GenCategoryTree(0, categorys), err
}

func GenCategoryTree(parentID int64, categorys []*model.Category) []*model.Category {
	tree := make([]*model.Category, 0)

	for _, m := range categorys {
		if m.ParentID == parentID {
			children := GenCategoryTree(m.ID, categorys)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 创建分类
func (r *Repository) Create(ctx context.Context, category *model.Category) error {
	err := r.tx.DB(ctx).Create(category).Error
	return err
}

// 更新分类
func (r *Repository) Update(ctx context.Context, id int64, category *model.Category) error {
	err := r.tx.DB(ctx).Model(category).Where("id = ?", id).Updates(category).Error
	return err
}

// 批量删除分类
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var categorys []*model.Category

	err := r.tx.DB(ctx).Where("id IN (?)", ids).Find(&categorys).Error
	if err != nil {
		return err
	}
	// 校验分类是否可以删除
	for _, category := range categorys {
		if category.ID == constant.DEFAULT_ID {
			return errors.New("默认分类不允许删除")
		}
		if r.isPID(ctx, int64(category.ID)) {
			return errors.New("该分类下包含子分类，请先删除子分类")
		}
	}

	err = r.tx.DB(ctx).Unscoped().Delete(&categorys).Error
	return err
}

// isPID 判断是否为别人的父类 ID
// 存在 true 不存在 false
func (r *Repository) isPID(ctx context.Context, ID int64) bool {
	var category model.Category
	if err := r.tx.DB(ctx).Where("parent_id = ?", ID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		} else {
			r.log.Error(err.Error())
			return false
		}
	}
	return true
}
