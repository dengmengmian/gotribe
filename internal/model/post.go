// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	Model
	Slug        string     `gorm:"type:varchar(255);uniqueIndex;comment:URL别名/Slug" json:"slug"`
	CategoryID  int64      `gorm:"index;comment:分类 ID" json:"category_id"`
	ProjectID   int64      `gorm:"index;comment:项目 ID" json:"project_id"`
	ColumnID    int64      `gorm:"index;comment:专栏ID" json:"column_id"`
	UserID      int64      `gorm:"index;comment:用户ID" json:"user_id"`
	Author      string     `gorm:"type:varchar(30);not null;index:idx_username;comment:作者" json:"author"`
	Title       string     `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Content     string     `gorm:"not null;type:text;comment:内容" json:"content"`
	HtmlContent string     `gorm:"not null;type:text;comment:html内容" json:"html_content"`
	Description string     `gorm:"not null;size:300;comment:描述" json:"description"`
	Ext         string     `gorm:"type:text;comment:'扩展字段'" json:"ext"`
	Icon        string     `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Tag         string     `gorm:"-" json:"tag"`
	View        int64      `gorm:"default:1;comment:'阅读量'" json:"view"`
	Type        uint       `gorm:"type:smallint;default:1;comment:类型，1.文章 2.page 3.短文" json:"type"`
	IsTop       int64      `gorm:"type:smallint;default:1;comment:是否置顶：1-禁用;2-启用" json:"is_top"`
	IsPasswd    int64      `gorm:"type:smallint;default:1;comment:是否加密：1-禁用;2-启用" json:"is_passwd"`
	PassWord    string     `gorm:"type:varchar(255);not null;comment:密码" json:"password"`
	Status      uint       `gorm:"type:smallint;not null;default:1;comment:状态，1-草稿；2-发布" json:"status"`
	UnitPrice   int64      `gorm:"type:integer;not null;comment:商品价格(分)" json:"unit_price"`
	Location    string     `gorm:"type:varchar(255);comment:地点" json:"location"`
	People      string     `gorm:"type:varchar(255);comment:人物" json:"people"`
	Time        *time.Time `gorm:"type:timestamp;comment:业务时间" json:"time"`
	Images      string     `gorm:"type:varchar(1000);comment:图片" json:"images"`
	ShowTime    *time.Time `gorm:"type:timestamp;comment:展示时间" json:"show_time"`
	Video       string     `gorm:"type:varchar(255);not null;comment:产品视频" json:"video"`
	// ToC-only fields
	PostID       string     `gorm:"column:post_id;type:varchar(255);comment:旧版字符ID/已废弃" json:"post_id"`
	DynamicType  string     `gorm:"column:dynamic_type;type:varchar(50);comment:动态类型" json:"dynamic_type"`
	Sort         uint       `gorm:"column:sort;default:0;comment:排序" json:"sort"`
	EventStartAt *time.Time `gorm:"column:event_start_at;type:timestamp;comment:活动开始时间" json:"event_start_at"`
	EventEndAt   *time.Time `gorm:"column:event_end_at;type:timestamp;comment:活动结束时间" json:"event_end_at"`
	RegisterURL  string     `gorm:"column:register_url;type:varchar(255);comment:报名链接" json:"register_url"`
	Category     *Category  `gorm:"-" json:"category"`
	Tags         []*Tag     `gorm:"-" json:"tags"`
	Project      *Project   `gorm:"-" json:"project"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	return nil
}
