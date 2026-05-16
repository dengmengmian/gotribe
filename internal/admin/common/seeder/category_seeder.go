package seeder

import (
	"gotribe/internal/model"

	"gorm.io/gorm"
)

// CategorySeeder 分类种子
type CategorySeeder struct {
	*BaseSeeder
}

// NewCategorySeeder 创建分类种子
func NewCategorySeeder() *CategorySeeder {
	return &CategorySeeder{
		BaseSeeder: NewBaseSeeder("category"),
	}
}

// Run 执行分类数据种子
func (s *CategorySeeder) Run(db *gorm.DB, syncExisting bool) error {
	categories := []*model.Category{
		{
			Model:       model.Model{ID: 1},
			Title:       "默认分类",
			Slug:        "default",
			Description: "默认分类",
			Status:      1,
		},
	}

	for _, category := range categories {
		if err := createIfNotExists(db, category, category.ID); err != nil {
			return err
		}
	}

	return nil
}
