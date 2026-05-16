package utils

import "github.com/rs/xid"

// UUID 生成全局唯一字符串 ID。
func UUID() string {
	return xid.New().String()
}
