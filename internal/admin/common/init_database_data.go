package common

import (
	"gotribe/internal/admin/common/seeder"
	coreconfig "gotribe/internal/core/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 初始化数据库数据
func InitData(db *gorm.DB, cfg *coreconfig.Config, log *zap.SugaredLogger) error {
	// 是否初始化数据
	if !cfg.Admin.InitData {
		return nil
	}

	// 注册所有种子
	registerAllSeeders()

	// 运行所有种子
	if err := seeder.RunSeeders(db, cfg.Admin.SyncSeedData); err != nil {
		log.Errorf("数据库种子执行失败: %v", err)
		return err
	}

	// 种子执行后重置 PostgreSQL 序列，避免自增主键冲突
	if cfg.Database.Type == "postgres" {
		if err := resetPostgresSequencesCore(db, log); err != nil {
			log.Errorf("重置 PostgreSQL 序列失败: %v", err)
			return err
		}
		log.Info("PostgreSQL 序列已重置")
	}
	return nil
}

// registerAllSeeders 注册所有种子
func registerAllSeeders() {
	// 基础数据种子
	seeder.RegisterSeeder(seeder.NewRoleSeeder())
	seeder.RegisterSeeder(seeder.NewAdminSeeder())
	seeder.RegisterSeeder(seeder.NewMenuSeeder())
	seeder.RegisterSeeder(seeder.NewApiSeeder())
	seeder.RegisterSeeder(seeder.NewSystemConfigSeeder())

	// 内容管理种子
	seeder.RegisterSeeder(seeder.NewCategorySeeder())
	seeder.RegisterSeeder(seeder.NewTagSeeder())
	seeder.RegisterSeeder(seeder.NewPostSeeder())

	// 项目管理种子
	seeder.RegisterSeeder(seeder.NewProjectSeeder())
	seeder.RegisterSeeder(seeder.NewUserSeeder())

	// 可以继续添加其他种子...
}
