// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gorm.io/gorm"
	"time"
)

// Model 基础模型
// swagger:model Model
type Model struct {
	// ID 主键
	// example: 1
	ID int64 `gorm:"primarykey" json:"id"`
	// 创建时间
	// example: 2023-01-01T00:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// 更新时间
	// example: 2023-01-01T00:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
	// 删除时间
	// swagger:type string
	// format: date-time
	// example: 2023-01-01T00:00:00Z
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" format:"date-time"`
}
