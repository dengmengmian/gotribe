package seeder

import (
	"gotribe/internal/model"

	"gorm.io/gorm"
)

// UserSeeder 用户种子
type UserSeeder struct {
	*BaseSeeder
}

// NewUserSeeder 创建用户种子
func NewUserSeeder() *UserSeeder {
	return &UserSeeder{
		BaseSeeder: NewBaseSeeder("user"),
	}
}

// Run 执行用户数据种子
func (s *UserSeeder) Run(db *gorm.DB, syncExisting bool) error {
	users := []*model.User{
		{
			Model:     model.Model{ID: 1},
			Username:  "gotribe",
			Nickname:  "gotribe",
			ProjectID: 1,
		},
	}

	for _, user := range users {
		if err := createIfNotExists(db, user, user.ID); err != nil {
			return err
		}
	}

	return nil
}
