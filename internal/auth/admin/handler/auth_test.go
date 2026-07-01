package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authservice "gotribe/internal/auth/admin/service"
	"gotribe/internal/auth/core"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"
)

// mockAuthService 模拟 auth.Service
type mockAuthService struct {
	loginFunc   func(ctx context.Context, username, password, clientIP string) (*authservice.LoginResult, error)
	refreshFunc func(ctx context.Context, userID int64, username string, issuedAt time.Time) (*authservice.LoginResult, error)
	logoutFunc  func(ctx context.Context, userID int64) error
}

func (m *mockAuthService) Login(ctx context.Context, username, password, clientIP string) (*authservice.LoginResult, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, username, password, clientIP)
	}
	return nil, errors.New("login not implemented")
}

func (m *mockAuthService) Refresh(ctx context.Context, userID int64, username string, issuedAt time.Time) (*authservice.LoginResult, error) {
	if m.refreshFunc != nil {
		return m.refreshFunc(ctx, userID, username, issuedAt)
	}
	return nil, errors.New("refresh not implemented")
}

func (m *mockAuthService) Logout(ctx context.Context, userID int64) error {
	if m.logoutFunc != nil {
		return m.logoutFunc(ctx, userID)
	}
	return nil
}

func newTestManager(t *testing.T, ttl time.Duration) *core.Manager {
	t.Helper()
	manager, err := core.NewManager("test", "test-secret-must-be-at-least-32-characters", map[string]core.AudienceConfig{
		core.AudienceAdmin: {
			Audience:        "test.admin",
			AccessTokenTTL:  ttl,
			RefreshTokenTTL: 24 * time.Hour,
		},
	})
	require.NoError(t, err)
	return manager
}

func setupTestHandler(t *testing.T) (*Handler, *core.Manager) {
	gin.SetMode(gin.TestMode)
	manager := newTestManager(t, time.Hour)
	svc := &mockAuthService{}
	h := NewHandler(core.AudienceAdmin, svc, manager, nil)
	return h, manager
}

func TestHandler_Login(t *testing.T) {
	h, _ := setupTestHandler(t)

	t.Run("登录成功", func(t *testing.T) {
		h.authService = &mockAuthService{
			loginFunc: func(ctx context.Context, username, password, clientIP string) (*authservice.LoginResult, error) {
				return &authservice.LoginResult{
					Token:   "fake-token",
					Expires: time.Now().Add(time.Hour),
					User:    &model.Admin{Username: username},
				}, nil
			},
		}

		payload := map[string]string{"username": "admin", "password": "Gotribe!23456"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/base/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.Login(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("缺少密码", func(t *testing.T) {
		payload := map[string]string{"username": "admin"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/base/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.Login(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("用户名或密码错误", func(t *testing.T) {
		h.authService = &mockAuthService{
			loginFunc: func(ctx context.Context, username, password, clientIP string) (*authservice.LoginResult, error) {
				return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "密码错误", nil)
			},
		}

		payload := map[string]string{"username": "admin", "password": "wrong"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/base/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.Login(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("服务内部错误", func(t *testing.T) {
		h.authService = &mockAuthService{
			loginFunc: func(ctx context.Context, username, password, clientIP string) (*authservice.LoginResult, error) {
				return nil, errors.New("db connection failed")
			},
		}

		payload := map[string]string{"username": "admin", "password": "Gotribe!23456"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/base/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.Login(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_Logout(t *testing.T) {
	h, _ := setupTestHandler(t)

	req := httptest.NewRequest("POST", "/base/logout", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Logout(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_RefreshToken(t *testing.T) {
	h, manager := setupTestHandler(t)

	t.Run("刷新成功", func(t *testing.T) {
		token, _, err := manager.SignAccessToken(core.AudienceAdmin, core.Subject{
			UserID:   1,
			Username: "testuser",
		})
		require.NoError(t, err)

		h.authService = &mockAuthService{
			refreshFunc: func(ctx context.Context, userID int64, username string, issuedAt time.Time) (*authservice.LoginResult, error) {
				return &authservice.LoginResult{
					Token:   "new-token",
					Expires: time.Now().Add(time.Hour),
				}, nil
			},
		}

		req := httptest.NewRequest("POST", "/base/refreshToken", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.RefreshToken(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("过期 token 也能刷新", func(t *testing.T) {
		// 使用极短 TTL 生成已过期令牌（VerifyAccessTokenWithoutExpiry 仍接受）
		expiredManager := newTestManager(t, time.Microsecond)
		token, _, err := expiredManager.SignAccessToken(core.AudienceAdmin, core.Subject{
			UserID:   1,
			Username: "testuser",
		})
		require.NoError(t, err)
		time.Sleep(50 * time.Millisecond)

		h.manager = expiredManager
		h.authService = &mockAuthService{
			refreshFunc: func(ctx context.Context, userID int64, username string, issuedAt time.Time) (*authservice.LoginResult, error) {
				return &authservice.LoginResult{
					Token:   "new-token",
					Expires: time.Now().Add(time.Hour),
				}, nil
			},
		}

		req := httptest.NewRequest("POST", "/base/refreshToken", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.RefreshToken(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("缺少 Authorization 头", func(t *testing.T) {
		// 用回正常 manager
		h.manager = newTestManager(t, time.Hour)
		req := httptest.NewRequest("POST", "/base/refreshToken", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.RefreshToken(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("刷新服务错误", func(t *testing.T) {
		h.manager = newTestManager(t, time.Hour)
		token, _, err := h.manager.SignAccessToken(core.AudienceAdmin, core.Subject{
			UserID:   1,
			Username: "testuser",
		})
		require.NoError(t, err)

		h.authService = &mockAuthService{
			refreshFunc: func(ctx context.Context, userID int64, username string, issuedAt time.Time) (*authservice.LoginResult, error) {
				return nil, errors.New("token generation failed")
			},
		}

		req := httptest.NewRequest("POST", "/base/refreshToken", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.RefreshToken(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
