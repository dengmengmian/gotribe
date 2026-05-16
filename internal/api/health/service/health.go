// Package service implements liveness and readiness probe logic.
package service

// 本文件实现存活探针和就绪探针的检查逻辑。

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gotribe/internal/buildinfo"
)

// DependencyStatus 表示单个依赖组件的健康状态。
type DependencyStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Readiness 表示 readiness 探针返回的整体依赖状态。
type Readiness struct {
	Status       string                      `json:"status"`
	Service      string                      `json:"service"`
	Version      string                      `json:"version"`
	Timestamp    time.Time                   `json:"timestamp"`
	Dependencies map[string]DependencyStatus `json:"dependencies"`
}

// Liveness 表示 liveness 探针返回的基础存活状态。
type Liveness struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

// VersionInfo 表示对外暴露的版本元数据。
type VersionInfo struct {
	Service   string `json:"service"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
}

// Service 负责封装健康检查相关的业务逻辑。
type Service struct {
	db      *gorm.DB
	redis   *redis.Client
	appName string
}

// NewService 创建健康检查服务实例。
func NewService(db *gorm.DB, redis *redis.Client, appName string) *Service {
	return &Service{
		db:      db,
		redis:   redis,
		appName: appName,
	}
}

// Liveness 返回应用存活状态。
func (s *Service) Liveness() Liveness {
	return Liveness{
		Status:    "ok",
		Service:   s.appName,
		Version:   buildinfo.Version,
		Timestamp: time.Now().UTC(),
	}
}

// Version 返回当前服务的构建版本信息。
func (s *Service) Version() VersionInfo {
	info := buildinfo.Current()
	return VersionInfo{
		Service:   s.appName,
		Version:   info.Version,
		Commit:    info.Commit,
		BuildTime: info.BuildTime,
	}
}

// Readiness 检查数据库和 Redis 是否已准备就绪。
func (s *Service) Readiness(ctx context.Context) (Readiness, bool) {
	result := Readiness{
		Status:    "ready",
		Service:   s.appName,
		Version:   buildinfo.Version,
		Timestamp: time.Now().UTC(),
		Dependencies: map[string]DependencyStatus{
			"database": {Status: "up"},
			"redis":    {Status: "up"},
		},
	}

	sqlDB, err := s.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		result.Status = "not_ready"
		result.Dependencies["database"] = DependencyStatus{Status: "down", Message: "database ping failed"}
	}

	if err := s.redis.Ping(ctx).Err(); err != nil {
		result.Status = "not_ready"
		result.Dependencies["redis"] = DependencyStatus{Status: "down", Message: err.Error()}
	}

	return result, result.Status == "ready"
}
