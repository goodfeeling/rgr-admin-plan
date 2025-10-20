// src/application/services/sys/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	configDomain "github.com/gbrayhan/microservices-go/src/domain/sys/config"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type ISysConfigService interface {
	GetConfig() (*configDomain.ConfigResponse, error)
	Update(module string, dataMap map[string]interface{}) error
	GetConfigByModule(module string) (*map[string]string, error)
}

type SysConfigUseCase struct {
	Logger     *logger.Logger
	configPath string
	configData map[string]map[string]interface{}
	mutex      sync.RWMutex
}

func NewSysConfigUseCase(
	loggerInstance *logger.Logger) ISysConfigService {

	service := &SysConfigUseCase{
		Logger:     loggerInstance,
		configPath: "config.yaml", // 默认配置文件路径
		configData: make(map[string]map[string]interface{}),
	}

	// 初始化时加载配置
	service.loadConfigFile()

	return service
}

// Update 更新配置文件中指定模块的配置
func (s *SysConfigUseCase) Update(module string, dataMap map[string]interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Logger.Info("Updating config",
		zap.String("module", module),
		zap.Any("data", dataMap))

	// 如果该模块不存在，创建它
	if _, exists := s.configData[module]; !exists {
		s.configData[module] = make(map[string]interface{})
	}

	// 更新模块中的配置项
	for key, value := range dataMap {
		s.configData[module][key] = value
	}

	// 保存到文件
	return s.saveConfigFile()
}

// GetConfig 获取整个配置文件
func (s *SysConfigUseCase) GetConfig() (*configDomain.ConfigResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.Logger.Info("Get config to group")

	// 序列化再反序列化确保类型兼容
	data, err := json.Marshal(s.configData)
	if err != nil {
		return nil, err
	}

	var configData configDomain.Config
	if err := json.Unmarshal(data, &configData); err != nil {
		return nil, err
	}

	// 构造 domain 层的 Config 结构
	config := &configDomain.ConfigResponse{
		Data: configData,
	}

	return config, nil
}

// GetConfigByModule 根据模块名获取配置，返回 map[string]string
func (s *SysConfigUseCase) GetConfigByModule(module string) (*map[string]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.Logger.Info("get config by module", zap.String("module", module))

	moduleData, exists := s.configData[module]
	if !exists {
		// 模块不存在，返回空map而不是错误
		emptyMap := make(map[string]string)
		return &emptyMap, nil
	}

	// 转换为 string:string 的 map
	result := make(map[string]string)
	for key, value := range moduleData {
		// 将各种类型转换为字符串
		switch v := value.(type) {
		case string:
			result[key] = v
		case int:
			result[key] = string(rune(v))
		case float64:
			result[key] = string(rune(int(v)))
		case bool:
			if v {
				result[key] = "true"
			} else {
				result[key] = "false"
			}
		default:
			// 对于复杂类型，转换为JSON字符串
			if bytes, err := json.Marshal(v); err == nil {
				result[key] = string(bytes)
			} else {
				result[key] = ""
			}
		}
	}

	return &result, nil
}

// loadConfigFile 从文件加载配置
func (s *SysConfigUseCase) loadConfigFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查配置文件是否存在
	if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
		// 文件不存在，创建默认配置
		s.configData = make(map[string]map[string]interface{})
		return s.saveConfigFile()
	}

	// 读取配置文件
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		s.Logger.Error("Failed to read config file", zap.Error(err))
		return err
	}

	// 解析 YAML
	if err := yaml.Unmarshal(data, &s.configData); err != nil {
		s.Logger.Error("Failed to parse config file", zap.Error(err))
		return err
	}

	s.Logger.Info("Config file loaded successfully", zap.String("path", s.configPath))
	return nil
}

// saveConfigFile 保存配置到文件
func (s *SysConfigUseCase) saveConfigFile() error {
	// 序列化为 YAML
	data, err := yaml.Marshal(s.configData)
	if err != nil {
		s.Logger.Error("Failed to marshal config data", zap.Error(err))
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(s.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.Logger.Error("Failed to create config directory", zap.Error(err))
		return err
	}
	// 写入文件
	if err := os.WriteFile(s.configPath, data, 0644); err != nil {
		s.Logger.Error("Failed to write config file", zap.Error(err))
		return err
	}

	s.Logger.Info("Config file saved successfully", zap.String("path", s.configPath))
	return nil
}

// SetConfigPath 设置配置文件路径
func (s *SysConfigUseCase) SetConfigPath(path string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.configPath = path
}

// ReloadConfig 重新加载配置文件
func (s *SysConfigUseCase) ReloadConfig() error {
	return s.loadConfigFile()
}
