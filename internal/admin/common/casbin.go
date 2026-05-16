package common

import (
	"fmt"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 初始化casbin策略管理器
func InitCasbinEnforcer(modelPath string, db *gorm.DB, log *zap.SugaredLogger) *casbin.Enforcer {
	e, err := databaseCasbin(modelPath, db, log)
	if err != nil {
		log.Panicf("初始化Casbin失败：%v", err)
		panic(fmt.Sprintf("初始化Casbin失败：%v", err))
	}

	log.Info("初始化Casbin完成!")
	return e
}

func databaseCasbin(modelPath string, db *gorm.DB, log *zap.SugaredLogger) (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 外部路径优先：GOTRIBE_ADMIN_RBAC_MODEL > configs/rbac_model.conf > 配置文件 > 内置
	envModelPath := os.Getenv("GOTRIBE_ADMIN_RBAC_MODEL")
	if envModelPath == "" {
		envModelPath = "configs/rbac_model.conf"
	}
	if _, statErr := os.Stat(envModelPath); statErr != nil {
		envModelPath = modelPath
	} else {
		modelPath = envModelPath
	}

	var e *casbin.Enforcer
	// 如果外部文件存在则优先使用
	if modelPath != "" {
		if _, statErr := os.Stat(modelPath); statErr == nil {
			log.Infof("加载外部 RBAC 模型: %s", modelPath)
			e, err = casbin.NewEnforcer(modelPath, a)
			if err != nil {
				log.Warnf("加载外部 RBAC 模型失败(%s): %v，使用内置默认", modelPath, err)
				e = nil
			}
		} else {
			log.Warnf("外部 RBAC 模型不可用(%s): %v，使用内置默认", modelPath, statErr)
		}
	}

	// 兜底：使用内置默认模型
	if e == nil {
		if embeddedRBACModel == "" {
			return nil, fmt.Errorf("内置 RBAC 模型为空，无法初始化")
		}
		m := model.NewModel()
		if err := m.LoadModelFromText(embeddedRBACModel); err != nil {
			return nil, fmt.Errorf("加载内置 RBAC 模型失败: %w", err)
		}
		log.Info("加载内置默认 RBAC 模型")
		e, err = casbin.NewEnforcer(m, a)
		if err != nil {
			return nil, err
		}
	}

	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}
	return e, nil
}
