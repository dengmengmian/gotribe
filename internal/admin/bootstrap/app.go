package bootstrap

// 本文件封装 Admin 端应用对象和 HTTP Server 生命周期；
// 与 internal/bootstrap/app.go 风格对齐。

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"

	"gotribe/internal/admin/routes"
	coreconfig "gotribe/internal/core/config"
	"gotribe/internal/request"
)

// AdminApp 封装 Admin 端应用运行时依赖和 HTTP 服务生命周期。
type AdminApp struct {
	cfg       coreconfig.Config
	srv       *http.Server
	providers *AdminProviders
}

// NewAdminApp 创建 Admin 应用实例并完成核心依赖初始化。
func NewAdminApp(cfg coreconfig.Config) (*AdminApp, error) {
	if err := request.Initialize(); err != nil {
		return nil, err
	}
	request.MaxBodyBytes = cfg.Server.MaxRequestBodyBytes

	providers, err := NewAdminProviders(cfg)
	if err != nil {
		return nil, err
	}

	router := routes.InitRoutes(&routes.Infra{
		DB:          providers.Infra.DB,
		Log:         providers.Infra.Log,
		Enforcer:    providers.Infra.Enforcer,
		Tx:          providers.Infra.Tx,
		AuthManager: providers.Infra.AuthManager,
	}, cfg, providers.Modules)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Admin.Host, cfg.Admin.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &AdminApp{
		cfg:       cfg,
		srv:       srv,
		providers: providers,
	}, nil
}

// Run 启动 HTTP 服务并开始监听请求。
func (a *AdminApp) Run() error {
	a.printBanner()
	return a.srv.ListenAndServe()
}

// Shutdown 优雅关闭应用：先停止接受请求，再走 providers.Close 链路（停 jobs、
// 关日志 channel、等待 worker、关 redis/db）。
func (a *AdminApp) Shutdown(ctx context.Context) error {
	var shutdownErr error
	if a != nil && a.srv != nil {
		shutdownErr = a.srv.Shutdown(ctx)
	}

	var closeErr error
	if a != nil && a.providers != nil {
		closeErr = a.providers.Close()
	}

	return errors.Join(shutdownErr, closeErr)
}

func (a *AdminApp) printBanner() {
	colorFg := color.New(color.FgCyan, color.Bold)
	colorFg.Println(`
			░██████╗░░█████╗░████████╗██████╗░██╗██████╗░███████╗
			██╔════╝░██╔══██╗╚══██╔══╝██╔══██╗██║██╔══██╗██╔════╝
			██║░░██╗░██║░░██║░░░██║░░░██████╔╝██║██████╦╝█████╗░░
			██║░░╚██╗██║░░██║░░░██║░░░██╔══██╗██║██╔══██╗██╔══╝░░
			╚██████╔╝╚█████╔╝░░░██║░░░██║░░██║██║██████╦╝███████╗
			░╚═════╝░░╚════╝░░░░╚═╝░░░╚═╝░░╚═╝╚═╝╚═════╝░╚══════╝`)
	fmt.Println("	App running at:")
	fmt.Printf("	- Local: http://localhost:%d\n", a.cfg.Admin.Port)
	fmt.Printf("	- Bind: %s:%d\n", a.cfg.Admin.Host, a.cfg.Admin.Port)
}
