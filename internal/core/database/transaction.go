package database

// 本文件封装数据库事务上下文和 DB 句柄获取逻辑。

import (
	"context"

	"gorm.io/gorm"
)

// txKey 用于在上下文中存放事务对象的键类型。
type txKey struct{}

// TransactionManager 表示数据库模块中的核心数据结构。
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器。
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// DB 返回当前上下文可用的数据库句柄。
func (m *TransactionManager) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return m.db.WithContext(ctx)
}

// WithinTransaction 在事务中执行指定业务逻辑。
func (m *TransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}
