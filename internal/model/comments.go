// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type Comment struct {
	Model
	ProjectID int64   `gorm:"not null;index;comment:项目ID;" json:"project_id"`
	Content     string `gorm:"not null;type:text;not null;comment:内容" json:"content"`
	HtmlContent string `gorm:"not null;type:text;not null;comment:HTML内容" json:"html_content"`
	Status uint   `gorm:"type:smallint;not null;index;default:1;comment:状态，1-待审核；2-审核通过" json:"status,omitempty"`
	ObjectID    string `gorm:"type:varchar(10);not null;index;comment:评论主题ID" json:"object_id"`
	ObjectType  int64   `gorm:"type:smallint;not null;default:1;index;comment:评论对象类型，1-文章；2-商品" json:"object_type"`
	Type uint   `gorm:"type:smallint;not null;default:1;comment:评论类型，1-评论；2-回复" json:"type"`
	UserID int64   `gorm:"not null;index;comment:用户ID" json:"user_id"`
	ToUserID    int64   `gorm:"not null;index;comment:被评论用户ID" json:"to_user_id"`
	ParentID    int    `gorm:"type:integer;not null;default:0;comment:父评论ID" json:"parent_id"`
	ReplyToID   int    `gorm:"type:integer;not null;default:0;comment:回复的评论ID" json:"reply_to_id"`
	Hot         int    `gorm:"type:integer;default:0;comment:热度" json:"hot"`
	Like        int    `gorm:"type:integer;default:0;comment:点赞数" json:"like"`
	Dislike     int    `gorm:"type:integer;default:0;comment:踩数" json:"dislike"`
	IP          string `gorm:"type:varchar(255);not null;comment:IP地址" json:"ip"`
	Country     string `gorm:"type:varchar(255);not null;comment:国家" json:"country"`
	RegionName  string `gorm:"type:varchar(255);not null;comment:地区" json:"region_name"`
	City        string `gorm:"type:varchar(255);not null;comment:城市" json:"city"`
	User        *User  `gorm:"-" json:"user"`
}
