package utils

// 本文件提供通用的字符串辅助函数。

import "strings"

// TrimOrEmpty 返回去除首尾空白后的字符串。
func TrimOrEmpty(value string) string {
	return strings.TrimSpace(value)
}
