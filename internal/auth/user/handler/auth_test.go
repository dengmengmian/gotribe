package handler

// 认证 handler 的 HTTP 层单元测试。

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gotribe/internal/auth/user/dto"
	"gotribe/internal/core/errs"
	"gotribe/internal/request"
)

// mockAuthService 是 authService 接口的 mock 实现。
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, projectID string, req dto.LoginRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *mockAuthService) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *mockAuthService) Logout(ctx context.Context, currentUserID int64, req dto.LogoutRequest) error {
	args := m.Called(ctx, currentUserID, req)
	return args.Error(0)
}

func setupAuthHandlerTest(t *testing.T) (*gin.Engine, *mockAuthService, *Handler) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	_ = request.Initialize()

	mockSvc := new(mockAuthService)
	h := &Handler{service: mockSvc}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("project_id", "test-project")
		c.Next()
	})
	router.POST("/login", h.Login)
	router.POST("/refresh", h.Refresh)

	logout := router.Group("/logout")
	logout.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("project_id", "test-project")
		c.Next()
	})
	logout.POST("", h.Logout)

	return router, mockSvc, h
}

func TestAuthHandler_LoginSuccess(t *testing.T) {
	router, mockSvc, _ := setupAuthHandlerTest(t)

	resp := &dto.AuthResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User:         dto.UserSummary{ID: 1, Username: "testuser"},
	}
	mockSvc.On("Login", mock.Anything, "test-project", dto.LoginRequest{Identity: "testuser", Password: "TestPass123"}).Return(resp, nil)

	body := []byte(`{"identity":"testuser","password":"TestPass123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Project-ID", "test-project")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test-access-token")
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_LoginInvalidCredentials(t *testing.T) {
	router, mockSvc, _ := setupAuthHandlerTest(t)

	mockSvc.On("Login", mock.Anything, "test-project", dto.LoginRequest{Identity: "testuser", Password: "wrongpass"}).Return(nil, errs.Unauthorized("invalid identity or password"))

	body := []byte(`{"identity":"testuser","password":"wrongpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Project-ID", "test-project")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_LoginInvalidBody(t *testing.T) {
	router, _, _ := setupAuthHandlerTest(t)

	body := []byte(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Project-ID", "test-project")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_RefreshSuccess(t *testing.T) {
	router, mockSvc, _ := setupAuthHandlerTest(t)

	resp := &dto.AuthResponse{AccessToken: "new-access", RefreshToken: "new-refresh", TokenType: "Bearer", ExpiresIn: 900}
	mockSvc.On("Refresh", mock.Anything, dto.RefreshRequest{RefreshToken: "old-refresh"}).Return(resp, nil)

	body := []byte(`{"refresh_token":"old-refresh"}`)
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "new-access")
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_LogoutSuccess(t *testing.T) {
	router, mockSvc, _ := setupAuthHandlerTest(t)

	mockSvc.On("Logout", mock.Anything, int64(1), dto.LogoutRequest{RefreshToken: "refresh-to-revoke"}).Return(nil)

	body := []byte(`{"refresh_token":"refresh-to-revoke"}`)
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
