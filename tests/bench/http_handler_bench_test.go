package bench

// HTTP handler 级别的基准测试，使用 httptest 无需外部依赖。

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gotribe/internal/api/health/handler"
	healthsvc "gotribe/internal/api/health/service"
	"gotribe/internal/core/response"
)

func BenchmarkHTTP_Health_Liveness(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	h := handler.NewHandler(healthsvc.NewService(nil, nil, "bench"))
	engine.GET("/livez", h.Liveness)

	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: %d", w.Code)
		}
	}
}

func BenchmarkResponse_JSON(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	type item struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	data := item{ID: 1, Title: "benchmark item"}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response.OK(c, data)
	}
}
