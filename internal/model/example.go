// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gotribe/internal/core/util"

	"gorm.io/gorm"
)

type Example struct {
	Model
	ExampleID   string `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"example_id"`
	ProjectID int64   `gorm:"not null;index;comment:项目ID;" json:"project_id"`
	Username    string `gorm:"type:varchar(30);not null;index:idx_username;comment:用户名" json:"username"`
	Title       string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Content     string `gorm:"not null;type:text;not null;comment:内容" json:"content"`
	Description string `gorm:"not null;size:300;not null;comment:描述" json:"description"`
	Status      uint8  `gorm:"type:smallint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
	// ToC-only fields
	UserID int64   `gorm:"column:user_id;index;comment:用户ID" json:"user_id"`
	OwnerUsername string `gorm:"column:owner_username;type:varchar(30);comment:拥有者用户名" json:"owner_username"`
	OwnerNickname string `gorm:"column:owner_nickname;type:varchar(30);comment:拥有者昵称" json:"owner_nickname"`
	Name          string `gorm:"column:name;type:varchar(255);comment:名称" json:"name"`
	PrimaryPostID string `gorm:"column:primary_post_id;type:varchar(255);comment:主要文章ID" json:"primary_post_id"`
}

func (e *Example) BeforeCreate(tx *gorm.DB) error {
	e.ExampleID = utils.GenShortID()

	return nil
}

// ExamplePost 表示示例业务单与文章的关联表。
type ExamplePost struct {
	Model
	ExampleRecordID int64   `gorm:"column:example_record_id;index;comment:示例业务单ID" json:"example_record_id"`
	ProjectID int64   `gorm:"column:project_id;not null;index;comment:项目ID" json:"project_id"`
	UserID int64   `gorm:"column:user_id;index;comment:用户ID" json:"user_id"`
	PostID          string `gorm:"column:post_id;type:varchar(255);comment:文章ID" json:"post_id"`
	PostTitle       string `gorm:"column:post_title;type:varchar(255);comment:文章标题" json:"post_title"`
	PostType        int16  `gorm:"column:post_type;comment:文章类型" json:"post_type"`
	PostStatus      int16  `gorm:"column:post_status;comment:文章状态" json:"post_status"`
	Sort            int    `gorm:"column:sort;default:0;comment:排序" json:"sort"`
}

func (ExamplePost) TableName() string {
	return "example_post"
}
