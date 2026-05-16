package common

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// CheckDBHealth 检查显式传入的数据库连接健康状态。
func CheckDBHealth(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("database is not initialized")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
