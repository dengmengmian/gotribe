package utils

import (
	"net"
	"net/http"
	"strings"
)

// ClientIP 解析请求的客户端 IP，优先读取 X-Forwarded-For 和 X-Real-IP 头。
func ClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if ip := strings.TrimSpace(strings.Split(xff, ",")[0]); ip != "" {
		return ip
	}
	if ip := strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
