package middleware

// middleware 包的单元测试。

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gotribe/internal/core/errs"
	profileview "gotribe/internal/api/profile/view"
)

func setupMiddlewareTest() (*gin.Engine, *httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return nil, w, c
}

// === Recovery ===

func TestRecovery_RecoversPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(Recovery())
	engine.GET("/panic", func(c *gin.Context) {
		panic("something went wrong")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRecovery_NoPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(Recovery())
	engine.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// === ProjectID ===

func TestProjectID_FromHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ProjectID("default"))
	engine.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, GetProjectID(c))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-ID", "my-project")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "my-project", w.Body.String())
}

func TestProjectID_FromDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ProjectID("fallback-project"))
	engine.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, GetProjectID(c))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "fallback-project", w.Body.String())
}

func TestProjectID_Required(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ProjectID(""))
	engine.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// === RequestID ===

func TestRequestID_GeneratesID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID())
	engine.GET("/test", func(c *gin.Context) {
		id, _ := c.Get(ContextKeyRequestID)
		c.String(http.StatusOK, id.(string))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body.String())
	assert.True(t, len(w.Body.String()) > 3)
}

func TestRequestID_PreservesHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID())
	engine.GET("/test", func(c *gin.Context) {
		id, _ := c.Get(ContextKeyRequestID)
		c.String(http.StatusOK, id.(string))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "custom-req-id")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "custom-req-id", w.Body.String())
}

// === CurrentUser ===

// mockCurrentUserReader 是 CurrentUserReader 的 mock 实现。
type mockCurrentUserReader struct {
	mock.Mock
}

func (m *mockCurrentUserReader) GetMe(ctx context.Context, projectID string, userID int64) (*profileview.MeView, error) {
	args := m.Called(ctx, projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*profileview.MeView), args.Error(1)
}

func TestCurrentUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reader := new(mockCurrentUserReader)
	user := &profileview.MeView{ID: 1, Username: "testuser", Nickname: "Test"}
	reader.On("GetMe", mock.Anything, "proj-1", int64(42)).Return(user, nil)

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set(ContextKeyAuth, &AuthContext{UserID: 42, ProjectID: "proj-1", Username: "testuser"})
		c.Next()
	})
	engine.Use(CurrentUser(reader))
	engine.GET("/me", func(c *gin.Context) {
		u, _ := GetCurrentUser(c)
		c.JSON(http.StatusOK, u)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
	reader.AssertExpectations(t)
}

func TestCurrentUser_MissingAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reader := new(mockCurrentUserReader)

	engine := gin.New()
	engine.Use(CurrentUser(reader))
	engine.GET("/me", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCurrentUser_ReaderError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reader := new(mockCurrentUserReader)
	reader.On("GetMe", mock.Anything, "proj-1", int64(42)).Return(nil, errs.NotFound("user not found", errors.New("db error")))

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set(ContextKeyAuth, &AuthContext{UserID: 42, ProjectID: "proj-1", Username: "testuser"})
		c.Next()
	})
	engine.Use(CurrentUser(reader))
	engine.GET("/me", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	reader.AssertExpectations(t)
}

// === AuthContext ===

func TestGetSetAuthContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, _, c := setupMiddlewareTest()

	auth := &AuthContext{UserID: 1, Username: "test", ProjectID: "proj"}
	SetAuthContext(c, auth)

	got, ok := GetAuthContext(c)
	assert.True(t, ok)
	assert.Equal(t, auth, got)
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, _, c := setupMiddlewareTest()

	c.Set(ContextKeyUserID, int64(99))
	id, ok := GetUserID(c)
	assert.True(t, ok)
	assert.Equal(t, int64(99), id)
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(SecurityHeaders())
	engine.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
}
