package routes

import (
	"context"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	adminrepo "gotribe/internal/admin/admin_user/repository"
	"gotribe/internal/admin/common"
	"gotribe/internal/admin/middleware"
	"gotribe/internal/auth/core"
	coreconfig "gotribe/internal/core/config"
	"gotribe/internal/core/database"
	appMiddleware "gotribe/internal/core/middleware"
	adminweb "gotribe/web/admin"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Infra 汇总路由层所需的基础设施依赖。
type Infra struct {
	DB          *gorm.DB
	Tx          *database.TransactionManager
	Log         *zap.SugaredLogger
	Enforcer    *casbin.Enforcer
	AuthManager *core.Manager
}

// nonSPAPathPrefixes 列出不应回退到 SPA index.html 的请求路径前缀。
// 注册新的后端路由时需同步更新本列表。
var nonSPAPathPrefixes = []string{"/api/", "/swagger/", "/assets/", "/health"}

// trustedProxyCIDRs 与 API 端 (internal/bootstrap/router.go) 保持一致，限制内网代理 IP 段。
var trustedProxyCIDRs = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

// InitRoutes 初始化所有路由。
// 不使用 gin.Default()，避免与 setupMiddlewares 注册的 Logger / Recovery 重复。
func InitRoutes(infra *Infra, cfg coreconfig.Config, modules *AdminModules) *gin.Engine {
	gin.SetMode(cfg.Admin.Mode)
	engine := gin.New()
	if cfg.App.IsProduction() {
		_ = engine.SetTrustedProxies(trustedProxyCIDRs)
	}

	adminRepo := adminrepo.NewRepository(infra.Tx)
	setupMiddlewares(engine, infra, cfg)
	jwtMW := core.JWTMiddleware(infra.AuthManager, core.AudienceAdmin, nil)
	adminLoader := middleware.AdminUserLoader(adminRepo, infra.Log)
	// 旧 RegisterRoutes 接收一个 authMiddleware 形参，是历史 API。认证流程实际上由
	// protected.Use(jwtMW, adminLoader) 应用，这里传入 noopAuth 占位即可。
	noopAuth := gin.HandlerFunc(func(c *gin.Context) { c.Next() })
	setupStaticFiles(engine, infra)
	setupSwaggerRoutes(engine)
	registerAPIRoutes(engine, infra, cfg, jwtMW, adminLoader, noopAuth, modules)
	registerJobRoutes(engine, infra, cfg, jwtMW, adminLoader, noopAuth, modules)
	initHealthRoute(engine, infra)

	infra.Log.Info("初始化路由完成！")
	return engine
}

func setupMiddlewares(engine *gin.Engine, infra *Infra, cfg coreconfig.Config) {
	engine.Use(appMiddleware.RequestID())
	engine.Use(appMiddleware.SecurityHeaders())

	fillInterval := time.Duration(cfg.Admin.RateLimit.FillInterval)
	capacity := cfg.Admin.RateLimit.Capacity
	engine.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))
	engine.Use(appMiddleware.CORSWithMode(cfg.Admin.Mode, 600, false))
	engine.Use(middleware.LangMiddleware())
	engine.Use(appMiddleware.Logger(cfg.App))
	engine.Use(appMiddleware.Recovery())
	engine.Use(middleware.OperationLogMiddleware(infra.DB, infra.Enforcer, cfg.Admin.UrlPathPrefix, infra.Log))
}

func setupStaticFiles(engine *gin.Engine, infra *Infra) {
	webDist := os.Getenv("GOTRIBE_ADMIN_WEB_DIST")
	if webDist != "" {
		setupLocalStaticFiles(engine, infra, webDist)
		return
	}

	distFS, err := fs.Sub(adminweb.Dist, "dist")
	if err != nil {
		infra.Log.Warnf("读取内嵌 Admin 前端资源失败：%v", err)
		setupLocalStaticFiles(engine, infra, "web/admin/dist")
		return
	}

	indexData, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		infra.Log.Warnf("读取内嵌 Admin index.html 失败：%v", err)
	}

	engine.Use(embeddedStaticMiddleware(distFS))
	setupSPAFallback(engine, indexData)
}

func setupLocalStaticFiles(engine *gin.Engine, infra *Infra, webDist string) {
	indexPath := filepath.Join(webDist, "index.html")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		infra.Log.Warnf("读取 %s 失败：%v", indexPath, err)
	}

	engine.Use(staticCacheControlMiddleware())
	engine.Use(static.Serve("/", static.LocalFile(webDist, false)))
	setupSPAFallback(engine, indexData)
}

func setupSPAFallback(engine *gin.Engine, indexData []byte) {
	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, prefix := range nonSPAPathPrefixes {
			if strings.HasPrefix(path, prefix) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}
		if len(indexData) > 0 {
			setAdminStaticCacheHeaders(c, "index.html")
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, string(indexData))
		} else {
			c.String(http.StatusNotFound, "Index file not found")
		}
	})
}

func embeddedStaticMiddleware(distFS fs.FS) gin.HandlerFunc {
	fileServer := http.FileServer(http.FS(distFS))
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Next()
			return
		}

		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if _, err := fs.Stat(distFS, path); err != nil {
			c.Next()
			return
		}

		if contentType := mime.TypeByExtension(filepath.Ext(path)); contentType != "" {
			c.Header("Content-Type", contentType)
		}
		setAdminStaticCacheHeaders(c, path)
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func staticCacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			setAdminStaticCacheHeaders(c, strings.TrimPrefix(c.Request.URL.Path, "/"))
		}
		c.Next()
	}
}

func setAdminStaticCacheHeaders(c *gin.Context, path string) {
	switch {
	case path == "" || path == "index.html" || path == "sw.js":
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
	case strings.HasPrefix(path, "assets/"):
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	}
}

func setupSwaggerRoutes(engine *gin.Engine) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func registerAPIRoutes(engine *gin.Engine, infra *Infra, cfg coreconfig.Config, jwtMW, adminLoader, authMiddleware gin.HandlerFunc, m *AdminModules) {
	prefix := cfg.Admin.UrlPathPrefix
	if prefix == "" {
		prefix = "api"
	}
	apiGroup := engine.Group("/" + prefix)

	apiGroup.POST("/base/login", m.Auth.Login)
	apiGroup.POST("/base/logout", m.Auth.Logout)
	apiGroup.POST("/base/refreshToken", m.Auth.RefreshToken)
	apiGroup.POST("/base/totp/verify", m.Auth.VerifyTOTP)
	apiGroup.POST("/base/totp/enroll", m.Auth.EnrollTOTP)
	apiGroup.POST("/base/totp/enroll/confirm", m.Auth.ConfirmEnrollTOTP)
	apiGroup.GET("/base/config", m.SystemConfig.Detail)

	protected := apiGroup.Group("")
	protected.Use(jwtMW)
	protected.Use(adminLoader)
	protected.Use(middleware.CasbinMiddleware(infra.Tx, infra.Enforcer, prefix))

	protected.GET("/base/totp/status", m.Auth.StatusTOTP)
	protected.POST("/base/totp/bind", m.Auth.BindTOTP)
	protected.POST("/base/totp/confirm", m.Auth.ConfirmTOTP)
	protected.DELETE("/base/totp", m.Auth.DeleteTOTP)
	protected.POST("/admin/:id/totp/reset", m.Auth.AdminResetTOTP)

	m.AI.RegisterRoutes(protected, authMiddleware)
	m.Admin.RegisterRoutes(protected, authMiddleware)
	m.Role.RegisterRoutes(protected, authMiddleware)
	m.Menu.RegisterRoutes(protected, authMiddleware)
	m.API.RegisterRoutes(protected, authMiddleware)
	m.OperationLog.RegisterRoutes(protected, authMiddleware)
	m.Project.RegisterRoutes(protected, authMiddleware)
	m.Config.RegisterRoutes(protected, authMiddleware)
	m.Tag.RegisterRoutes(protected, authMiddleware)
	m.Category.RegisterRoutes(protected, authMiddleware)
	m.Post.RegisterRoutes(protected, authMiddleware)
	m.User.RegisterRoutes(protected, authMiddleware)
	m.Resource.RegisterRoutes(protected, authMiddleware)
	m.Column.RegisterRoutes(protected, authMiddleware)
	m.AdScene.RegisterRoutes(protected, authMiddleware)
	m.Ad.RegisterRoutes(protected, authMiddleware)
	m.Comment.RegisterRoutes(protected, authMiddleware)
	m.Point.RegisterRoutes(protected, authMiddleware)
	m.SystemConfig.RegisterRoutes(protected, authMiddleware)
	m.Feedback.RegisterRoutes(protected, authMiddleware)
	m.Index.RegisterRoutes(protected, authMiddleware)
}

func registerJobRoutes(engine *gin.Engine, infra *Infra, cfg coreconfig.Config, jwtMW, adminLoader, authMiddleware gin.HandlerFunc, m *AdminModules) {
	_ = authMiddleware // 与 registerAPIRoutes 签名对齐，job 路由当前未使用 inline 校验
	jobGroup := engine.Group("/api/jobs")
	jobGroup.Use(jwtMW)
	jobGroup.Use(adminLoader)
	jobGroup.Use(middleware.CasbinMiddleware(infra.Tx, infra.Enforcer, cfg.Admin.UrlPathPrefix))
	{
		jobGroup.GET("/", m.Job.ListJobs)
		jobGroup.GET("/:name/status", m.Job.GetJobStatus)
		jobGroup.GET("/:name/history", m.Job.GetJobHistory)
		jobGroup.POST("/:name/enable", m.Job.EnableJob)
		jobGroup.POST("/:name/disable", m.Job.DisableJob)
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp int64             `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func initHealthRoute(engine *gin.Engine, infra *Infra) {
	engine.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		checks := make(map[string]string)
		status := "healthy"
		httpStatus := http.StatusOK

		if err := common.CheckDBHealth(ctx, infra.DB); err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			status = "unhealthy"
			httpStatus = http.StatusServiceUnavailable
		} else {
			checks["database"] = "healthy"
		}

		if infra.Enforcer == nil {
			checks["casbin"] = "unhealthy: not initialized"
			status = "unhealthy"
			httpStatus = http.StatusServiceUnavailable
		} else {
			checks["casbin"] = "healthy"
		}

		c.JSON(httpStatus, HealthResponse{
			Status:    status,
			Timestamp: time.Now().Unix(),
			Checks:    checks,
		})
	})
}
