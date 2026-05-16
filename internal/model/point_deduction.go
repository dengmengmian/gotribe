// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

// PointDeduction 扣减积分表
type PointDeduction struct {
	Model
	ProjectID int64  `gorm:"not null;comment:项目ID;" json:"project_id"`
	UserID int64  `gorm:"index:idx_point_deduction_available_user,priority:2;comment:用户ID" json:"user_id"`
	Points            int64 `gorm:"type:bigint;NOT NULL;comment:积分数值(分)"`
	PointsDetailID    int64  `gorm:"comment:'积分明细ID'"`
	AvailablePointsID int64  `gorm:"not null;index:idx_point_deduction_available_user,priority:1;comment:'可用积分表ID'"`
}
