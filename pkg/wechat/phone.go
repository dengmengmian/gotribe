// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package wechat

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"` // 修改为 string 类型
	Watermark       struct {
		Timestamp int    `json:"timestamp"`
		Appid     string `json:"appid"`
	} `json:"watermark"`
}

type Result struct {
	ErrCode   int       `json:"errcode"`
	ErrMsg    string    `json:"errmsg"`
	PhoneInfo PhoneInfo `json:"phone_info"`
}

type QueryNumber struct {
	AccessToken string `json:"access_token"`
}

func GetPhoneNumber(accessToken, Code string) (Result, error) {
	client := resty.New()

	query := QueryNumber{
		AccessToken: accessToken,
	}

	body := map[string]string{
		"code": Code,
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"access_token": query.AccessToken,
		}).
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		Post("https://api.weixin.qq.com/wxa/business/getuserphonenumber")

	if err != nil {
		log.Panicln(err)
	}

	var res Result
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return res, err
	}

	if resp.IsError() {
		return res, fmt.Errorf("failed to get phone number: %s", res.ErrMsg)
	}

	return res, nil
}
