package seeder

import (
	"gotribe/internal/model"

	"gorm.io/gorm"
)

// ProjectSeeder 项目种子
type ProjectSeeder struct {
	*BaseSeeder
}

// NewProjectSeeder 创建项目种子
func NewProjectSeeder() *ProjectSeeder {
	return &ProjectSeeder{
		BaseSeeder: NewBaseSeeder("project"),
	}
}

// Run 执行项目数据种子
func (s *ProjectSeeder) Run(db *gorm.DB, syncExisting bool) error {
	projects := []*model.Project{
		{
			Model:       model.Model{ID: 1},
			Name:        "default",
			Title:       "默认项目",
			Description: "默认项目",
		},
	}

	for _, project := range projects {
		if err := createIfNotExists(db, project, project.ID); err != nil {
			return err
		}
	}

	return nil
}
