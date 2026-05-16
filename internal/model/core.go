package model

import "time"

// AuthUser 定义认证场景下用户表映射，组合 Model 和 Core。
// Model.ID 为 int64，Core.ProjectID 为 string，分别适配 Admin int64 和 ToC string 项目 ID 体系。
type AuthUser struct {
	Model
	Core
}

// TableName 返回当前模型对应的数据表名。
func (AuthUser) TableName() string {
	return "user"
}

// Core 定义 user 表在多个模块间共享的字段映射。
type Core struct {
	Username  string     `gorm:"column:username"`
	ProjectID string     `gorm:"column:project_id"`
	Password  string     `gorm:"column:password" json:"-"`
	Nickname  string     `gorm:"column:nickname"`
	Email     string     `gorm:"column:email"`
	Phone     string     `gorm:"column:phone"`
	Sex       string     `gorm:"column:sex"`
	Status    int16      `gorm:"column:status"`
	Birthday  *time.Time `gorm:"column:birthday"`
	AvatarURL string     `gorm:"column:avatar_url"`
}
