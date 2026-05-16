// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type ThirdPartyAccounts struct {
	Model
	UserID int64   `gorm:"index;comment:用户ID" json:"user_id"`
	Platform string `gorm:"type:varchar(50);not null;comment:平台" json:"platform"`
	BindFlag int64   `gorm:"type:smallint;default:1;comment:是否绑定,2绑定" json:"bind_flag"`
	OpenID   string `gorm:"type:varchar(255);uniqueIndex;not null;comment:openID" json:"open_id"`
}
