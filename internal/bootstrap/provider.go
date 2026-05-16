package bootstrap

// 本文件负责装配数据库、缓存、中间件和业务模块依赖。

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"gotribe/internal/auth/core"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/config"
	"gotribe/internal/core/database"
	applog "gotribe/internal/core/logger"
	"gotribe/internal/observability"
)

// Infra 汇总应用启动后共享的基础设施依赖。
type Infra struct {
	DB         *gorm.DB
	Redis      *redis.Client
	Keys       *cache.KeyBuilder
	Store      *cache.Store
	Tx         *database.TransactionManager
	JWT        *core.Manager
	AuthTokens *core.TokenStore
	AppName    string
	CacheTTL   int
	UserAuth   config.AuthAudienceConfig
}

// Providers 汇总应用基础设施和业务模块依赖。
type Providers struct {
	Infra   *Infra
	Modules *Modules
}

// NewProviders 创建应用运行所需的基础设施和业务依赖。
func NewProviders(cfg config.Config) (*Providers, error) {
	infra, err := newInfra(cfg)
	if err != nil {
		return nil, err
	}

	modules := buildModules(infra)

	return &Providers{
		Infra:   infra,
		Modules: modules,
	}, nil
}

func newInfra(cfg config.Config) (*Infra, error) {
	db, err := database.NewGORM(cfg.App, cfg.Database)
	if err != nil {
		return nil, err
	}
	if err := observability.InstrumentGORM(db); err != nil {
		applog.Warn(nil, "gorm observability disabled", "err", err)
	}

	redisClient, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		if sqlDB, sqlErr := db.DB(); sqlErr == nil {
			_ = sqlDB.Close()
		}
		return nil, err
	}
	if err := observability.InstrumentRedis(redisClient); err != nil {
		applog.Warn(nil, "redis observability disabled", "err", err)
	}

	keys := cache.NewKeyBuilder(cfg.App.Name)

	audiences := map[string]core.AudienceConfig{
		core.AudienceUser: {
			Audience:        cfg.Auth.User.Audience,
			AccessTokenTTL:  cfg.Auth.User.AccessTokenTTL(),
			RefreshTokenTTL: cfg.Auth.User.RefreshTokenTTL(),
		},
	}
	authManager, err := core.NewManager(cfg.Auth.Issuer, cfg.Auth.Secret, audiences)
	if err != nil {
		_ = redisClient.Close()
		if sqlDB, sqlErr := db.DB(); sqlErr == nil {
			_ = sqlDB.Close()
		}
		return nil, fmt.Errorf("auth manager: %w", err)
	}

	return &Infra{
		DB:         db,
		Redis:      redisClient,
		Keys:       keys,
		Store:      cache.NewStore(redisClient, keys),
		Tx:         database.NewTransactionManager(db),
		JWT:        authManager,
		AuthTokens: core.NewTokenStore(redisClient, keys),
		AppName:    cfg.App.Name,
		CacheTTL:   cfg.Redis.DefaultCacheTTLMins,
		UserAuth:   cfg.Auth.User,
	}, nil
}

// Close 关闭 Providers 管理的基础设施资源。
func (p *Providers) Close() error {
	if p == nil || p.Infra == nil {
		return nil
	}

	var closeErrs []error
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
