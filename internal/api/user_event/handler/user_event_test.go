package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	usereventrepo "gotribe/internal/api/user_event/repository"
	usereventservice "gotribe/internal/api/user_event/service"
	"gotribe/internal/core/database"
	"gotribe/internal/core/middleware"
)

func TestHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	tx := database.NewTransactionManager(db)
	require.NoError(t, db.Exec(`CREATE TABLE user_events (
		id integer primary key autoincrement,
		user_id integer,
		project_id text,
		event_type integer,
		event_detail text,
		duration integer,
		ip text,
		user_agent text,
		referer text,
		platform text,
		created_at datetime,
		updated_at datetime,
		deleted_at datetime
	)`).Error)

	handler := NewHandler(usereventservice.NewService(usereventrepo.NewRepository(tx)))
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyAuth, &middleware.AuthContext{UserID: 1, ProjectID: "1", Username: "tester"})
		c.Set(middleware.ContextKeyProjectID, "1")
		c.Next()
	})
	engine.POST("/events", handler.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{"event_type":1,"event_detail":"click"}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}
