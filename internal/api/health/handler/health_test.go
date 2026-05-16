package handler

// 本文件展示 handler 层的 HTTP 单元测试写法，使用 httptest 和 gin 测试模式。

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/buildinfo"
	"gotribe/internal/api/health/service"
)

func setupHealthHandler(t *testing.T) (*Handler, *gin.Engine, func()) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	srv := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: srv.Addr()})

	svc := service.NewService(db, redisClient, "gotribe-test")
	h := NewHandler(svc)

	engine := gin.New()
	engine.GET("/version", h.Version)
	engine.GET("/livez", h.Liveness)
	engine.GET("/readyz", h.Readiness)

	cleanup := func() {
		_ = redisClient.Close()
		srv.Close()
	}

	return h, engine, cleanup
}

func TestHandler_Liveness(t *testing.T) {
	_, engine, cleanup := setupHealthHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/livez", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data struct {
			Status  string `json:"status"`
			Service string `json:"service"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body.Data.Status)
	assert.Equal(t, "gotribe-test", body.Data.Service)
}

func TestHandler_Readiness_AllUp(t *testing.T) {
	_, engine, cleanup := setupHealthHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/readyz", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data struct {
			Status       string `json:"status"`
			Dependencies map[string]struct {
				Status string `json:"status"`
			} `json:"dependencies"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ready", body.Data.Status)
	assert.Equal(t, "up", body.Data.Dependencies["database"].Status)
	assert.Equal(t, "up", body.Data.Dependencies["redis"].Status)
}

func TestHandler_Version(t *testing.T) {
	originalVersion := buildinfo.Version
	originalCommit := buildinfo.Commit
	originalBuildTime := buildinfo.BuildTime
	buildinfo.Version = "v1.2.3"
	buildinfo.Commit = "abc1234"
	buildinfo.BuildTime = "2026-04-09T12:00:00Z"
	defer func() {
		buildinfo.Version = originalVersion
		buildinfo.Commit = originalCommit
		buildinfo.BuildTime = originalBuildTime
	}()

	_, engine, cleanup := setupHealthHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/version", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data service.VersionInfo `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "gotribe-test", body.Data.Service)
	assert.Equal(t, "v1.2.3", body.Data.Version)
	assert.Equal(t, "abc1234", body.Data.Commit)
	assert.Equal(t, "2026-04-09T12:00:00Z", body.Data.BuildTime)
}
