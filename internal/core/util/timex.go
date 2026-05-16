package utils

// 本文件提供通用的时间格式化辅助函数。

import "time"

// ToDateString 将时间指针格式化为日期字符串。
func ToDateString(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
