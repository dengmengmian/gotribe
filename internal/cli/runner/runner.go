// Package runner provides a graceful-shutdown helper shared by API and Admin entry points.
package runner

// 本文件抽取 cmd/api 与 cmd/admin 共用的信号监听与优雅停机逻辑。

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	applog "gotribe/internal/core/logger"
)

// defaultShutdownTimeout 在调用方未提供超时时回退到此值。
const defaultShutdownTimeout = 5 * time.Second

// Server 描述支持优雅停机的 HTTP 应用。
type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}

// Run 启动 server 并处理 SIGINT / SIGTERM 优雅退出。
// name 用于日志标识；shutdownTimeout <= 0 时回退到 defaultShutdownTimeout。
func Run(name string, shutdownTimeout time.Duration, server Server) error {
	if shutdownTimeout <= 0 {
		shutdownTimeout = defaultShutdownTimeout
	}

	runErr := make(chan error, 1)
	go func() {
		applog.Info(context.Background(), name+" starting")
		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			runErr <- err
			return
		}
		runErr <- nil
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-runErr:
		if err != nil {
			return err
		}
		applog.Info(context.Background(), name+" stopped")
		return nil
	case <-quit:
	}

	applog.Info(context.Background(), name+" shutting down", "timeout", shutdownTimeout.String())
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	applog.Info(context.Background(), name+" stopped")
	return nil
}
