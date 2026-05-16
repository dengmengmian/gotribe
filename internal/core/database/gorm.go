// Package database provides GORM initialization, connection pooling, transaction management, and pagination utilities.
package database

// 本文件负责初始化 GORM、连接池和 SQL 日志策略。

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gotribe/internal/core/config"
)

// NewGORM 初始化 GORM 数据库连接和日志策略。
func NewGORM(app config.AppConfig, cfg config.DatabaseConfig) (*gorm.DB, error) {
	if cfg.Type != "postgres" {
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	logMode := gormlogger.Error
	if app.IsDevelopment() && cfg.LogMode {
		logMode = gormlogger.Info
	}

	gormLog := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logMode,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormLog,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	configurePool(sqlDB, cfg)
	if err := sqlDB.PingContext(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}

// configurePool 根据配置调整数据库连接池参数。
func configurePool(db *sql.DB, cfg config.DatabaseConfig) {
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime())
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime())
}
