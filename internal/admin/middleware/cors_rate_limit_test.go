package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	coreMiddleware "gotribe/internal/core/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCORSMiddlewareAllowsLocalhostOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(coreMiddleware.CORSWithMode(gin.DebugMode, 600, true))
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	require.Equal(t, "Origin", w.Header().Get("Vary"))
}

func TestCORSMiddlewareRejectsUnknownPreflightOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(coreMiddleware.CORSWithMode(gin.ReleaseMode, 600, true))
	router.OPTIONS("/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestRateLimitMiddlewareReturnsTooManyRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RateLimitMiddleware(time.Hour, 1))
	router.GET("/limited", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/limited", nil))
	require.Equal(t, http.StatusOK, first.Code)

	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/limited", nil))
	require.Equal(t, http.StatusTooManyRequests, second.Code)
	require.Contains(t, second.Body.String(), "访问限流")
}

func TestShouldSkipOperationLogPaths(t *testing.T) {
	require.True(t, shouldSkipLog(""))
	require.True(t, shouldSkipLog("/assets/app.js"))
	require.True(t, shouldSkipLog("/favicon.ico"))
	require.False(t, shouldSkipLog("/post"))
	require.Equal(t, http.StatusTooManyRequests, http.StatusTooManyRequests)
}
