package middleware

// 本文件提供统一的跨域处理中间件。

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 创建跨域处理中间件。
func CORS(allowedOrigins, allowedHeaders, allowedMethods []string) gin.HandlerFunc {
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	// 预计算允许的 origins 集合用于快速查找
	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[strings.TrimSpace(o)] = struct{}{}
	}

	// 预计算 headers 和 methods 字符串
	headers := strings.Join(allowedHeaders, ", ")
	methods := strings.Join(allowedMethods, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 确定是否允许该 origin
		if allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			if _, ok := originSet[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		c.Writer.Header().Set("Access-Control-Allow-Methods", methods)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// CORSWithMode 创建支持 debug 模式的跨域中间件。
// mode: gin.DebugMode 时允许 localhost 来源，其他模式仅允许显式配置的来源。
func CORSWithMode(mode string, maxAge int, allowCredentials bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			if isAllowedOrigin(origin, mode) {
				c.Header("Vary", "Origin")
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, Session, Content-Type, Accept-Language")
				c.Header("Access-Control-Expose-Headers", "Content-Length")
				c.Header("Access-Control-Max-Age", formatMaxAge(maxAge))
				if allowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
			} else if method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func isAllowedOrigin(origin, mode string) bool {
	if strings.EqualFold(mode, gin.ReleaseMode) {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

func formatMaxAge(maxAge int) string {
	if maxAge <= 0 {
		return "600"
	}
	return strconv.Itoa(maxAge)
}
