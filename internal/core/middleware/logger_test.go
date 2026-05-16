package middleware

// 本文件验证请求日志中间件提取统一业务错误的能力。

import (
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
)

// TestLatestAppError 验证日志中间件可从 Gin 错误上下文中提取统一业务错误。
func TestLatestAppError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	_ = ctx.Error(errs.BadRequest("invalid request", nil))

	appErr := latestAppError(ctx)
	if appErr == nil {
		t.Fatalf("latestAppError() = nil, want app error")
	}
	if appErr.Code != errs.CodeBadRequest {
		t.Fatalf("error code = %q, want %q", appErr.Code, errs.CodeBadRequest)
	}
}

// TestSanitizeQueryRedactsSensitiveValues 验证日志查询串会对敏感值进行脱敏。
func TestSanitizeQueryRedactsSensitiveValues(t *testing.T) {
	t.Parallel()

	got := sanitizeQuery("password=secret&foo=bar&refresh_token=token")
	if got == "" {
		t.Fatalf("sanitizeQuery() = empty, want redacted query")
	}
	values, err := url.ParseQuery(got)
	if err != nil {
		t.Fatalf("ParseQuery() error = %v", err)
	}
	if values.Get("password") != "REDACTED" {
		t.Fatalf("password = %q, want REDACTED", values.Get("password"))
	}
	if values.Get("refresh_token") != "REDACTED" {
		t.Fatalf("refresh_token = %q, want REDACTED", values.Get("refresh_token"))
	}
	if values.Get("foo") != "bar" {
		t.Fatalf("foo = %q, want %q", values.Get("foo"), "bar")
	}
}
