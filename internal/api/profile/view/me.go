// Package view defines reusable internal view structures for the profile module.
package view

// 本文件定义 profile 模块内部复用的用户资料视图。

// MeView 表示 profile 模块在服务层和中间件之间传递的用户资料视图。
type MeView struct {
	ID         int64
	Username   string
	ProjectID  string
	Nickname   string
	Email      string
	Phone      string
	Sex        string
	Status     int16
	Birthday   string
	Background string
	Ext        string
	AvatarURL  string
	CreatedAt  string
	UpdatedAt  string
}
