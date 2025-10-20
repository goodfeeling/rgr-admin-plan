package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// loadYAMLConfigToEnv loads configuration from config.yaml and sets as environment variables
func LoadYAMLConfigToEnv() error {
	configPath := "config.yaml"

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("config.yaml not found, skipping...")
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析 YAML
	var configData map[string]map[string]interface{}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 将配置设置为环境变量
	for module, moduleData := range configData {
		for key, value := range moduleData {
			// 构造环境变量名称，例如: MODULE_KEY
			envKey := strings.ToUpper(module + "_" + key)

			// 将值转换为字符串
			envValue := convertToString(value)

			// 设置环境变量
			os.Setenv(envKey, envValue)
		}
	}

	fmt.Println("Configuration loaded from config.yaml to environment variables")
	return nil
}

// convertToString converts various types to string
func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		// 对于复杂类型，转换为JSON字符串
		if bytes, err := json.Marshal(v); err == nil {
			return string(bytes)
		}
		return ""
	}
}
