package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/admin/admin_user/repository"
	"gotribe/internal/core/database"
	coremw "gotribe/internal/core/middleware"
	"gotribe/internal/model"
	utils "gotribe/internal/core/util"
)

// setupTestRepo 创建内存 SQLite 数据库并写入一个种子 admin。
func setupTestRepo(t *testing.T, seedStatus int) (*repository.Repository, *zap.SugaredLogger, int64) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Admin{}))

	hashedPassword, err := utils.PasswordUtil.GenPasswd("Gotribe!23456")
	require.NoError(t, err)
	admin := &model.Admin{
		Username: "testuser",
		Password: hashedPassword,
	}
	require.NoError(t, db.Create(admin).Error)
	// 显式设置 status，避免 GORM 跳过零值
	require.NoError(t, db.Model(admin).Update("status", seedStatus).Error)

	tx := database.NewTransactionManager(db)
	repo := repository.NewRepository(tx)
	log, _ := zap.NewDevelopment()
	return repo, log.Sugar(), admin.ID
}

// setUserContext 模拟 core.JWTMiddleware 已经写入 user 上下文。
func setUserContext(c *gin.Context, userID int64, username string) {
	coremw.SetAuthContext(c, &coremw.AuthContext{UserID: userID, Username: username})
	c.Set(coremw.ContextKeyUserID, userID)
	c.Set(coremw.ContextKeyUsername, username)
}

func TestAdminUserLoader(t *testing.T) {
	t.Run("缺少 user 上下文", func(t *testing.T) {
		repo, log, _ := setupTestRepo(t, 1)
		mw := AdminUserLoader(repo, log)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		mw(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("用户不存在", func(t *testing.T) {
		repo, log, _ := setupTestRepo(t, 1)
		mw := AdminUserLoader(repo, log)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		setUserContext(c, 999, "nonexistent")
		mw(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("用户被禁用", func(t *testing.T) {
		repo, log, adminID := setupTestRepo(t, 0) // status = 0 (disabled)
		mw := AdminUserLoader(repo, log)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		setUserContext(c, adminID, "testuser")
		mw(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("加载成功", func(t *testing.T) {
		repo, log, adminID := setupTestRepo(t, 1)
		mw := AdminUserLoader(repo, log)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		setUserContext(c, adminID, "testuser")
		mw(c)
		assert.Equal(t, http.StatusOK, w.Code)

		user, exists := c.Get("user")
		assert.True(t, exists)
		admin, ok := user.(model.Admin)
		assert.True(t, ok)
		assert.Equal(t, adminID, admin.ID)
	})
}
