package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// SEO SEO工具类，用于处理SEO相关操作
type SEO struct{}

// PushBaidu 推送URL到百度
func (s *SEO) PushBaidu(site, token string, urls string) (bool, error) {
	api := "http://data.zz.baidu.com/urls?site=" + site + "&token=" + token
	resp, err := http.Post(api, "text/plain", strings.NewReader(urls))
	if err != nil {
		return false, fmt.Errorf("推送失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("解析响应失败: %v", err)
	}

	log.Printf("推送记录: %s", string(body))
	return true, nil
}

// SEOUtil 全局 SEO 工具实例
var SEOUtil = &SEO{}
