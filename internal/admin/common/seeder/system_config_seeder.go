package seeder

import (
	"gotribe/internal/model"

	"gorm.io/gorm"
)

// SystemConfigSeeder 系统配置种子
type SystemConfigSeeder struct {
	*BaseSeeder
}

// NewSystemConfigSeeder 创建系统配置种子
func NewSystemConfigSeeder() *SystemConfigSeeder {
	return &SystemConfigSeeder{
		BaseSeeder: NewBaseSeeder("system_config"),
	}
}

// Run 执行系统配置数据种子
func (s *SystemConfigSeeder) Run(db *gorm.DB, syncExisting bool) error {
	configs := []*model.SystemConfig{
		{
			Model:          model.Model{ID: 1},
			SystemConfigID: "245eko",
			Title:          "GoTribe管理后台",
			Logo:           "https://cdn.example.com/logo.png",
			Icon:           "https://cdn.example.com/icon.png",
		},
	}

	for _, config := range configs {
		if err := createIfNotExists(db, config, config.ID); err != nil {
			return err
		}
	}

	return nil
}
