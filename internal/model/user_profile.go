package model

// 本文件定义用户资料表映射结构。

// UserProfile 定义用户资料表的字段映射。
type UserProfile struct {
	Model
	Core
	Background string `gorm:"column:background"`
	Ext        string `gorm:"column:ext"`
}

// TableName 返回当前模型对应的数据表名。
func (UserProfile) TableName() string {
	return "user"
}
