package jobs

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gotribe/internal/core/database"
)

// InitJobs 初始化所有任务
func InitJobs(db *gorm.DB, log *zap.SugaredLogger) error {
	tx := database.NewTransactionManager(db)
	log.Info("Initializing jobs...")

	// 初始化全局注册表
	registry := InitGlobalRegistry(log)

	// 注册站点地图任务
	if IsJobEnabled(DefaultJobsConfig(), "sitemap") {
		sitemapConfig, _ := GetJobConfig(DefaultJobsConfig(), "sitemap")
		sitemapJob := NewSitemapJob(sitemapConfig, tx, log)
		if err := RegisterJob(sitemapJob); err != nil {
			log.Errorf("Failed to register sitemap job: %v", err)
			return err
		}
	}

	// 注册示例任务
	if IsJobEnabled(DefaultJobsConfig(), "example") {
		exampleConfig, _ := GetJobConfig(DefaultJobsConfig(), "example")
		exampleJob := NewExampleJob(exampleConfig, log)
		if err := RegisterJob(exampleJob); err != nil {
			log.Errorf("Failed to register example job: %v", err)
			return err
		}
	}

	log.Info("Jobs initialized successfully")

	// 启动定时任务
	if err := registry.Start(); err != nil {
		log.Errorf("Failed to start jobs: %v", err)
		return err
	}

	log.Info("All jobs started successfully")
	return nil
}

// StartAllJobs 启动所有任务
func StartAllJobs(log *zap.SugaredLogger) error {
	registry := GetGlobalRegistry()
	if registry == nil {
		registry = InitGlobalRegistry(log)
	}
	log.Info("Starting all jobs...")

	if err := StartJobs(); err != nil {
		log.Errorf("Failed to start jobs: %v", err)
		return err
	}

	log.Info("All jobs started successfully")
	return nil
}

// StopAllJobs 停止所有任务
func StopAllJobs(log *zap.SugaredLogger) error {
	log.Info("Stopping all jobs...")

	if err := StopJobs(); err != nil {
		log.Errorf("Failed to stop jobs: %v", err)
		return err
	}

	log.Info("All jobs stopped successfully")
	return nil
}
