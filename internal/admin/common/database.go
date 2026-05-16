package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gotribe/internal/admin/migration"
	coreconfig "gotribe/internal/core/config"
	"gotribe/internal/model/migrate"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// buildPostgresDSN 构建 PostgreSQL 连接字符串（用于日志展示）。
func buildPostgresDSN(dbCfg *coreconfig.DatabaseConfig, maskPassword bool) string {
	sslMode := dbCfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	queryPart := ""
	if dbCfg.Query != "" {
		queryPart = " " + dbCfg.Query
	}

	password := dbCfg.Password
	if maskPassword {
		password = "******"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s%s",
		dbCfg.Host, dbCfg.Username, password, dbCfg.Database, dbCfg.Port, sslMode, queryPart)
}

// RunMigrations 运行数据库迁移和健康检查。
// 数据库连接已由调用方通过 core/database.NewGORM 创建。
func RunMigrations(db *gorm.DB, cfg *coreconfig.Config, log *zap.SugaredLogger) {
	dbCfg := &cfg.Database
	adminCfg := &cfg.Admin

	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("获取数据库连接失败: %v", err)
	}

	migrationsPath := filepath.Join("migrations", dbCfg.Type)
	hasMigrations := false
	if _, err := os.Stat(migrationsPath); err == nil {
		hasMigrations = true
	}

	if adminCfg.Mode == "release" {
		if adminCfg.EnableMigrate {
			log.Warn("生产环境已禁用自动迁移（enable-migrate=true 在 release 模式下被忽略），请使用 migrations/ 目录下的 SQL migration")
		}
		if hasMigrations {
			version, dirty, err := migration.CheckMigrationVersion(sqlDB, dbCfg.Type, migrationsPath)
			if err != nil {
				log.Warnf("migration 状态检查失败: %v，请确保已手动运行 migrations/ 下的 SQL 文件", err)
			} else if dirty {
				log.Warnf("migration 处于 dirty 状态 (version=%d)，请手动修复", version)
			} else {
				log.Infof("当前 migration 版本: %d", version)
			}
		}
	} else {
		if hasMigrations {
			log.Info("检测到 migrations 目录，使用 golang-migrate 执行迁移")
			if version, dirty, err := migration.CheckMigrationVersion(sqlDB, dbCfg.Type, migrationsPath); err != nil {
				log.Warnf("migration 状态检查失败，将继续尝试执行迁移: %v", err)
			} else if dirty {
				log.Panicf("migration 处于 dirty 状态 (version=%d)，请先重置开发库或手动修复 schema_migrations；开发环境可执行: make dev-db-reset", version)
			}
			if err := migration.RunMigrations(sqlDB, dbCfg.Type, migrationsPath); err != nil {
				log.Panicf("golang-migrate 执行失败: %v", err)
			}
			log.Info("golang-migrate 迁移完成")
			if dbCfg.Type == "postgres" {
				if err := resetPostgresSequencesCore(db, log); err != nil {
					log.Errorf("重置 PostgreSQL 序列失败: %v", err)
				}
			}
		} else {
			log.Warn("migrations 目录不存在，跳过数据库迁移（请确保 migrations/ 目录包含 SQL 迁移文件）")
		}
	}

	migrate.RunHealthChecks(db)
	log.Info("数据库健康检查通过")
}

// sequenceInfoCore 自增序列元数据。
type sequenceInfoCore struct {
	TableName    string
	ColumnName   string
	SequenceName string
}

// resetPostgresSequencesCore 重置所有自增序列。
func resetPostgresSequencesCore(db *gorm.DB, log *zap.SugaredLogger) error {
	var sequences []sequenceInfoCore
	sql := `
		SELECT
			c.relname AS table_name,
			a.attname AS column_name,
			pg_get_serial_sequence(c.relname, a.attname) AS sequence_name
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		JOIN pg_attribute a ON a.attrelid = c.oid
		WHERE c.relkind = 'r'
		  AND n.nspname = 'public'
		  AND a.attnum > 0
		  AND NOT a.attisdropped
		  AND pg_get_serial_sequence(c.relname, a.attname) IS NOT NULL
	`
	if err := db.Raw(sql).Scan(&sequences).Error; err != nil {
		return fmt.Errorf("查询自增序列失败: %w", err)
	}

	var errs []error
	for _, seq := range sequences {
		resetSQL := fmt.Sprintf(
			`SELECT setval('%s', COALESCE(MAX("%s"), 0) + 1, false) FROM "%s"`,
			seq.SequenceName, seq.ColumnName, seq.TableName,
		)
		if err := db.Exec(resetSQL).Error; err != nil {
			log.Warnf("重置序列失败 %s.%s (%s): %v", seq.TableName, seq.ColumnName, seq.SequenceName, err)
			errs = append(errs, fmt.Errorf("%s.%s: %w", seq.TableName, seq.ColumnName, err))
		}
	}
	return errors.Join(errs...)
}
