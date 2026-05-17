package bootstrap

// 本文件负责装配 Admin 端基础设施和业务模块依赖；与 internal/bootstrap/provider.go 风格对齐。

import (
	"errors"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"

	"gotribe/internal/admin/common"
	"gotribe/internal/admin/jobs"
	"gotribe/internal/admin/middleware"
	ologrepo "gotribe/internal/admin/operation_log/repository"
	resourceservice "gotribe/internal/admin/resource/service"
	"gotribe/internal/admin/routes"
	"gotribe/internal/auth/core"
	"gotribe/internal/core/cache"
	coreconfig "gotribe/internal/core/config"
	"gotribe/internal/core/database"
	"gotribe/internal/core/logger"
)

// operationLogWorkerCount 操作日志消费者协程数。
const operationLogWorkerCount = 3

// AdminInfra 汇总 Admin 端启动后共享的基础设施依赖。
type AdminInfra struct {
	Cfg         coreconfig.Config
	DB          *gorm.DB
	Redis       *redis.Client
	Tx          *database.TransactionManager
	Log         *zap.SugaredLogger
	Enforcer    *casbin.Enforcer
	AuthManager *core.Manager
	UploadCfg   resourceservice.UploadConfig
	CDNDomain   string
}

// AdminProviders 汇总 Admin 端基础设施和业务模块依赖。
type AdminProviders struct {
	Infra      *AdminInfra
	Modules    *routes.AdminModules
	logWorkers sync.WaitGroup
}

// NewAdminProviders 创建 Admin 端运行所需的基础设施和业务依赖。
// 启动期任何关键步骤失败都会返回错误，已分配的资源会回滚。
func NewAdminProviders(cfg coreconfig.Config) (*AdminProviders, error) {
	logger.InitWithRotation(
		cfg.Admin.Logs.Path,
		zapcore.Level(int8(cfg.Admin.Logs.Level)),
		cfg.Admin.Logs.MaxSize,
		cfg.Admin.Logs.MaxBackups,
		cfg.Admin.Logs.MaxAge,
		cfg.Admin.Logs.Compress,
	)
	applog := logger.Sugared()

	db, err := database.NewGORM(cfg.App, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database init: %w", err)
	}
	closeDB := func() {
		if sqlDB, derr := db.DB(); derr == nil {
			_ = sqlDB.Close()
		}
	}

	common.RunMigrations(db, &cfg, applog)

	enforcer := common.InitCasbinEnforcer(cfg.Admin.Casbin.ModelPath, db, applog)

	if err := common.InitData(db, &cfg, applog); err != nil {
		closeDB()
		return nil, fmt.Errorf("init seed data: %w", err)
	}

	if err := jobs.InitJobs(db, applog); err != nil {
		closeDB()
		return nil, fmt.Errorf("init jobs: %w", err)
	}
	if err := jobs.StartAllJobs(applog); err != nil {
		closeDB()
		return nil, fmt.Errorf("start jobs: %w", err)
	}

	// Redis 不是 Admin 端核心依赖，仅用于仪表盘缓存，初始化失败保持告警继续启动。
	redisClient, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		applog.Warnf("Redis 初始化失败，仪表盘缓存管理将不可用: %v", err)
	}

	tx := database.NewTransactionManager(db)

	audiences := map[string]core.AudienceConfig{
		core.AudienceAdmin: {
			Audience:        cfg.Auth.Admin.Audience,
			AccessTokenTTL:  cfg.Auth.Admin.AccessTokenTTL(),
			RefreshTokenTTL: cfg.Auth.Admin.RefreshTokenTTL(),
		},
	}
	authManager, err := core.NewManager(cfg.Auth.Issuer, cfg.Auth.Secret, audiences)
	if err != nil {
		closeDB()
		return nil, fmt.Errorf("auth manager: %w", err)
	}

	uploadCfg, cdnDomain := resolveUploadConfig(cfg)

	infra := &AdminInfra{
		Cfg:         cfg,
		DB:          db,
		Redis:       redisClient,
		Tx:          tx,
		Log:         applog,
		Enforcer:    enforcer,
		AuthManager: authManager,
		UploadCfg:   uploadCfg,
		CDNDomain:   cdnDomain,
	}

	modules := routes.BuildAdminModules(tx, enforcer, applog, authManager, cdnDomain, uploadCfg, redisClient, cfg.AI)

	p := &AdminProviders{
		Infra:   infra,
		Modules: modules,
	}

	logRepository := ologrepo.NewRepository(tx, applog)
	for i := 0; i < operationLogWorkerCount; i++ {
		p.logWorkers.Add(1)
		go func() {
			defer p.logWorkers.Done()
			logRepository.SaveOperationLogChannel(middleware.OperationLogChan)
		}()
	}

	return p, nil
}

// Close 关闭 Admin 端基础设施资源：停定时任务、关闭操作日志 channel 并等待 worker，
// 最后关闭 redis / db 连接。
func (p *AdminProviders) Close() error {
	if p == nil || p.Infra == nil {
		return nil
	}

	var closeErrs []error

	if err := jobs.StopAllJobs(p.Infra.Log); err != nil {
		p.Infra.Log.Errorf("Failed to stop jobs: %v", err)
		closeErrs = append(closeErrs, err)
	}

	close(middleware.OperationLogChan)
	p.logWorkers.Wait()

	if p.Infra.Redis != nil {
		if err := p.Infra.Redis.Close(); err != nil {
			closeErrs = append(closeErrs, err)
		}
	}
	if p.Infra.DB != nil {
		sqlDB, err := p.Infra.DB.DB()
		if err != nil {
			closeErrs = append(closeErrs, err)
		} else if err := sqlDB.Close(); err != nil {
			closeErrs = append(closeErrs, err)
		}
	}
	return errors.Join(closeErrs...)
}

// resolveUploadConfig 解析公用 / Admin 专属上传配置的 fallback 关系。
func resolveUploadConfig(cfg coreconfig.Config) (resourceservice.UploadConfig, string) {
	uploadConfig := cfg.Upload
	if uploadConfig.Provider == "" && cfg.Admin.Upload.Provider != "" {
		uploadConfig = cfg.Admin.Upload
	}
	cdnDomain := uploadConfig.CDNDomain
	if cdnDomain == "" {
		cdnDomain = cfg.Admin.CDNDomain
	}

	uploadCfg := resourceservice.UploadConfig{
		Provider:  uploadConfig.Provider,
		Endpoint:  uploadConfig.Endpoint,
		AccessKey: uploadConfig.AccessKey,
		SecretKey: uploadConfig.SecretKey,
		Bucket:    uploadConfig.Bucket,
		Region:    uploadConfig.Region,
		CDNDomain: cdnDomain,
	}
	if uploadCfg.Provider == "" {
		if cfg.Admin.EnableOss {
			uploadCfg.Provider = "oss"
		} else {
			uploadCfg.Provider = "qiniu"
		}
	}
	return uploadCfg, cdnDomain
}
