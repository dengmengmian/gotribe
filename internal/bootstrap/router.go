package bootstrap

// 本文件负责创建总路由骨架并挂接各模块路由。

import (
	"time"

	"github.com/gin-gonic/gin"

	"gotribe/internal/auth/core"
	"gotribe/internal/core/config"
	"gotribe/internal/core/middleware"
	"gotribe/internal/observability"
)

// NewRouter 创建总路由骨架并挂接各模块路由。
func NewRouter(cfg config.Config, providers *Providers) (*gin.Engine, error) {
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	if cfg.App.IsProduction() {
		_ = router.SetTrustedProxies([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	}
	router.Use(
		middleware.RequestID(),
		observability.TracingMiddleware(),
		observability.MetricsMiddleware(),
		middleware.CORS(cfg.Server.CORS.AllowedOrigins, cfg.Server.CORS.AllowedHeaders, cfg.Server.CORS.AllowedMethods),
		middleware.SecurityHeaders(),
		middleware.Logger(cfg.App),
		middleware.Recovery(),
	)

	if cfg.Observability.MetricsEnabled {
		router.GET("/metrics", gin.WrapH(observability.MetricsHandler()))
	}

	router.GET("/version", providers.Modules.Health.Handler.Version)
	router.GET("/livez", providers.Modules.Health.Handler.Liveness)
	router.GET("/readyz", providers.Modules.Health.Handler.Readiness)

	authLimiter := middleware.RateLimit(
		providers.Infra.Redis,
		providers.Infra.Keys,
		cfg.RateLimit.Enabled,
		"auth",
		cfg.RateLimit.AuthPerMinute,
		time.Minute,
		true,
		middleware.ClientIPResolver(),
	)
	apiLimiter := middleware.RateLimit(
		providers.Infra.Redis,
		providers.Infra.Keys,
		cfg.RateLimit.Enabled,
		"api",
		cfg.RateLimit.APIForMinute,
		time.Minute,
		false,
		middleware.UserOrIPResolver(),
	)
	eventLimiter := middleware.RateLimit(
		providers.Infra.Redis,
		providers.Infra.Keys,
		cfg.RateLimit.Enabled,
		"user_event",
		cfg.RateLimit.EventPerMin,
		time.Minute,
		false,
		middleware.UserOrIPResolver(),
	)

	jwtMiddleware := core.JWTMiddleware(providers.Infra.JWT, core.AudienceUser, providers.Infra.AuthTokens)
	currentUserMiddleware := middleware.CurrentUser(providers.Modules.Profile.Service)

	authRoutes := router.Group("/api/v1/auth")
	authRoutes.Use(middleware.ProjectID(cfg.App.DefaultProjectID))
	providers.Modules.Auth.Handler.RegisterRoutes(authRoutes, authLimiter, jwtMiddleware)

	api := router.Group("/api/v1")
	api.Use(middleware.ProjectID(cfg.App.DefaultProjectID))

	public := api.Group("/")
	secured := api.Group("/")
	secured.Use(jwtMiddleware)
	currentUser := secured.Group("/")
	currentUser.Use(currentUserMiddleware)

	providers.Modules.Post.Handler.RegisterRoutes(public, apiLimiter)
	providers.Modules.Config.Handler.RegisterRoutes(public, apiLimiter)
	providers.Modules.Profile.Handler.RegisterRoutes(secured, currentUser, apiLimiter, authLimiter)
	providers.Modules.Tag.Handler.RegisterRoutes(public)
	providers.Modules.Category.Handler.RegisterRoutes(public)
	providers.Modules.Example.Handler.RegisterRoutes(currentUser, apiLimiter)
	providers.Modules.UserEvent.Handler.RegisterRoutes(public, eventLimiter)

	return router, nil
}
