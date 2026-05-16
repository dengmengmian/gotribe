package response

// 本文件验证统一错误响应会将业务错误挂入 Gin 上下文。

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
)

// TestErrorWritesResponseAndContext 验证错误响应会写出统一结构并保留错误上下文。
func TestErrorWritesResponseAndContext(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("request_id", "req_test")

	Error(ctx, errs.BadRequest("invalid request", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if len(ctx.Errors) != 1 {
		t.Fatalf("errors count = %d, want 1", len(ctx.Errors))
	}

	appErr := errs.As(ctx.Errors.Last().Err)
	if appErr == nil {
		t.Fatalf("last error = nil, want app error")
	}
	if appErr.Code != errs.CodeBadRequest {
		t.Fatalf("error code = %q, want %q", appErr.Code, errs.CodeBadRequest)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got := body["code"]; got != string(errs.CodeBadRequest) {
		t.Fatalf("response code = %#v, want %q", got, errs.CodeBadRequest)
	}
	if got := body["request_id"]; got != "req_test" {
		t.Fatalf("request_id = %#v, want %q", got, "req_test")
	}
}
