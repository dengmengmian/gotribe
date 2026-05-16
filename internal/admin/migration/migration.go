package migration

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations 执行数据库迁移（up）
// dbType: "postgres" 或 "mysql"
// migrationsPath: 迁移文件目录的绝对路径，例如 "/app/migrations/postgres"
func RunMigrations(sqlDB *sql.DB, dbType string, migrationsPath string) error {
	var driver database.Driver
	var err error

	switch dbType {
	case "postgres":
		driver, err = postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			return fmt.Errorf("create postgres migration driver failed: %w", err)
		}
	default:
		// 默认 MySQL
		driver, err = mysql.WithInstance(sqlDB, &mysql.Config{})
		if err != nil {
			return fmt.Errorf("create mysql migration driver failed: %w", err)
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		dbType,
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance failed: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migration up failed: %w", err)
	}

	return nil
}

// CheckMigrationVersion 检查当前 migration 版本和是否存在待执行的迁移
// 返回 (version, dirty, error)
func CheckMigrationVersion(sqlDB *sql.DB, dbType string, migrationsPath string) (version int64, dirty bool, err error) {
	var driver database.Driver

	switch dbType {
	case "postgres":
		driver, err = postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			return 0, false, fmt.Errorf("create postgres migration driver failed: %w", err)
		}
	default:
		driver, err = mysql.WithInstance(sqlDB, &mysql.Config{})
		if err != nil {
			return 0, false, fmt.Errorf("create mysql migration driver failed: %w", err)
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		dbType,
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("create migrate instance failed: %w", err)
	}

	versionRaw, dirty, err := m.Version()
	version = int64(versionRaw)
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("get migration version failed: %w", err)
	}

	return version, dirty, nil
}
