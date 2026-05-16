// @title           Gotribe Admin API
// @version         1.0
// @description     gotribe Admin 端接口文档。
// @contact.name   gotribe
// @contact.url    https://github.com/dengmengmian/gotribe
// @license.name  Apache 2.0
// @license.url   https://www.apache.org/licenses/LICENSE-2.0
// @BasePath  /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 Bearer {token}，注意 Bearer 和 token 之间有一个空格

package main

import (
	"log"

	_ "gotribe/docs/admin/swagger"
	"gotribe/internal/admin/bootstrap"
	"gotribe/internal/cli/runner"
	"gotribe/internal/core/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := bootstrap.NewAdminApp(cfg)
	if err != nil {
		log.Fatalf("bootstrap admin app: %v", err)
	}

	if err := runner.Run("admin", cfg.Server.ShutdownTimeout(), app); err != nil {
		log.Fatalf("admin server: %v", err)
	}
}
