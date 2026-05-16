// Package main is the application entry point that loads configuration, bootstraps dependencies, and starts the HTTP server.
package main

// 本文件是应用入口，负责加载配置、启动 HTTP 服务并处理优雅退出。

import (
	"log"

	"gotribe/internal/bootstrap"
	"gotribe/internal/cli/runner"
	"gotribe/internal/core/config"
)

// main 负责加载配置并启动应用服务。
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		log.Fatalf("bootstrap app: %v", err)
	}

	if err := runner.Run("api", cfg.Server.ShutdownTimeout(), app); err != nil {
		log.Fatalf("api server: %v", err)
	}
}
