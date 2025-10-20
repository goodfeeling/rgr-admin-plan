package utils

import (
	"os"
	"path/filepath"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 初始化Casbin执行器
func InitCasbinEnforcer(db *gorm.DB, logger *logger.Logger) (*casbin.Enforcer, error) {
	// 创建Casbin适配器
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, &gormadapter.CasbinRule{}, "casbin_rule")
	if err != nil {
		return nil, err
	}

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		logger.Error("Error getting working directory", zap.Error(err))
		return nil, err
	}

	// 构建配置文件路径
	modelPath := filepath.Join(wd, "config", "model.conf")

	// 检查文件是否存在
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		// 如果config/model.conf不存在，尝试其他可能的路径
		alternativePaths := []string{
			filepath.Join(wd, "src", "config", "model.conf"),
			filepath.Join(wd, "..", "config", "model.conf"),
			"config/model.conf",
		}

		found := false
		for _, path := range alternativePaths {
			if _, err := os.Stat(path); err == nil {
				modelPath = path
				found = true
				break
			}
		}

		if !found {
			logger.Error("Model config file not found", zap.String("attempted_path", modelPath))
			return nil, err
		}
	}

	// 从配置文件创建模型
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		logger.Error("Error creating model from file", zap.Error(err), zap.String("path", modelPath))
		return nil, err
	}

	// 创建执行器
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		logger.Error("Error creating enforcer", zap.Error(err))
		return nil, err
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		logger.Error("Error loading policy", zap.Error(err))
		return nil, err
	}

	return enforcer, nil
}
