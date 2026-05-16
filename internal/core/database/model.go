package database

// 本文件定义各数据表共用的基础模型字段。

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 封装数据库公共字段，供各表结构复用。
type BaseModel struct {
	ID        int64          `gorm:"column:id;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
