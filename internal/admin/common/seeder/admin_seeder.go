package seeder

import (
	"gotribe/internal/core/util"
	"gotribe/internal/model"

	"gorm.io/gorm"
)

// AdminSeeder 管理员种子
type AdminSeeder struct {
	*BaseSeeder
}

// NewAdminSeeder 创建管理员种子
func NewAdminSeeder() *AdminSeeder {
	return &AdminSeeder{
		BaseSeeder: NewBaseSeeder("admin"),
	}
}

// Run 执行管理员数据种子
func (s *AdminSeeder) Run(db *gorm.DB, syncExisting bool) error {
	// 获取角色
	var adminRole model.Role
	if err := db.First(&adminRole, 1).Error; err != nil {
		return err
	}

	// 默认管理员密码，首次登录后请务必修改
	password, err := utils.PasswordUtil.GenPasswd("Gotribe!23456")
	if err != nil {
		return err
	}

	admin := &model.Admin{
		Model:        model.Model{ID: 1},
		Username:     "admin",
		Password:     password,
		Mobile:       "18888888888",
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Nickname:     new(string),
		Introduction: new(string),
		Status:       1,
		Creator:      "系统",
		Roles:        []*model.Role{&adminRole},
	}

	if err := createIfNotExists(db, admin, admin.ID); err != nil {
		return err
	}

	return ensureAdminRole(db, admin.ID, adminRole.ID)
}
