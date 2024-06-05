// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

// PostM 是数据库中 post 记录 struct 格式的映射.
type PostM struct {
	gorm.Model
	PostID      string `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"postID"`
	CategoryID  string `gorm:"type:varchar(10);Index;comment:分类 ID" json:"categoryID"`
	ProjectID   string `gorm:"type:varchar(10);Index;comment:项目 ID" json:"projectID"`
	ColumnID    string `gorm:"type:varchar(10);Index;comment:专栏ID" json:"columnID"`
	UserID      string `gorm:"type:varchar(10);Index;comment:用户ID" json:"userID"`
	Author      string `gorm:"type:varchar(30);not null;index:idx_username;comment:作者" json:"author"`
	Title       string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Content     string `gorm:"not null;type:longtext;comment:内容" json:"content"`
	HtmlContent string `gorm:"not null;type:longtext;comment:html内容" json:"htmlContent"`
	Description string `gorm:"not null;size:300;comment:描述" json:"description"`
	Ext         string `gorm:"type:text;comment:'扩展字段'" json:"ext"`
	Icon        string `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Tag         string `gorm:"type:varchar(30);comment:tag" json:"tag"`
	View        uint   `gorm:"default:1;comment:'阅读量'" json:"view"`
	Type        uint   `gorm:"type:tinyint;default:1;comment:类型，1-文章" json:"type"`
	IsTop       uint   `gorm:"type:tinyint;default:2;comment:是否置顶：1-启用;2-禁用" json:"isTop"`
	IsPasswd    uint   `gorm:"type:tinyint;default:2;comment:是否加密：1-启用;2-禁用" json:"isPasswd"`
	PassWord    string `gorm:"type:varchar(255);not null;comment:密码" json:"password"`
	Status      uint   `gorm:"type:tinyint(1);not null;default:1;comment:状态，1-启用；2-发布" json:"status"`
}

// TableName 用来指定映射的 MySQL 表名.
func (p *PostM) TableName() string {
	return "post"
}

// BeforeCreate 在创建数据库记录之前生成 postID.
func (p *PostM) BeforeCreate(tx *gorm.DB) error {
	p.PostID = gid.GenShortID()

	return nil
}
