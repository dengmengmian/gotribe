// Package bootstrap provides application lifecycle management, dependency wiring, and HTTP routing initialization.
package bootstrap

// 本文件封装应用对象和 HTTP Server 生命周期。

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/config"
	"gotribe/internal/core/logger"
	"gotribe/internal/observability"
	"gotribe/internal/request"
)

// App 封装应用运行时依赖和 HTTP 服务生命周期。
type App struct {
	engine    *gin.Engine
	server    *http.Server
	providers *Providers
	obsClose  observability.ShutdownFn
}

// NewApp 创建应用实例并完成核心依赖初始化。
func NewApp(cfg config.Config) (*App, error) {
	logger.Init(cfg.App)

	if err := request.Initialize(); err != nil {
		return nil, err
	}
	request.MaxBodyBytes = cfg.Server.MaxRequestBodyBytes

	obsClose := observability.Init(cfg)

	providers, err := NewProviders(cfg)
	if err != nil {
		return nil, err
	}

	engine, err := NewRouter(cfg, providers)
	if err != nil {
		_ = providers.Close()
		return nil, err
	}

	server := &http.Server{
		Addr:           cfg.Server.Address(),
		Handler:        engine,
		ReadTimeout:    cfg.Server.ReadTimeout(),
		WriteTimeout:   cfg.Server.WriteTimeout(),
		IdleTimeout:    cfg.Server.IdleTimeout(),
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	return &App{
		engine:    engine,
		server:    server,
		providers: providers,
		obsClose:  obsClose,
	}, nil
}

// Run 启动 HTTP 服务并开始监听请求。
func (a *App) Run() error {
	return a.server.ListenAndServe()
}

// Shutdown 执行应用优雅停机。
func (a *App) Shutdown(ctx context.Context) error {
	var shutdownErr error
	if a != nil && a.server != nil {
		shutdownErr = a.server.Shutdown(ctx)
	}

	var closeErr error
	if a != nil {
		closeErr = a.providers.Close()
	}

	var obsErr error
	if a != nil && a.obsClose != nil {
		obsErr = a.obsClose(ctx)
	}

	return errors.Join(shutdownErr, closeErr, obsErr)
}
