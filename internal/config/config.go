package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"qwqserver/pkg/defaultvalue"
	"qwqserver/pkg/util/singleton"
	"reflect"
	"strings"
)

type Config struct {
	*Listen    `yaml:"listen"`
	*Database  `yaml:"database"`
	*Redis     `yaml:"redis"`
	*AdminUser `yaml:"admin_user"`
}

var (
	configSingleton = singleton.NewSingleton[Config]()
)

// DBDriverName 获取数据库驱动名称（已弃用）
func DBDriverName() string {
	return ""
}

// CacheDriverName 获取缓存驱动名称（已弃用）
func CacheDriverName() string {
	return "redis"
}

// WriteDefaultIfNotExists 检查文件路径是否存在，若不存在则创建目录并写入默认数据
// filePath: 目标文件路径
// defaultData: 要写入的默认数据
// dirPerm: 目录权限 (例如 0755)
// filePerm: 文件权限 (例如 0644)
func WriteDefaultIfNotExists(filePath string, defaultData []byte, dirPerm, filePerm os.FileMode) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 获取文件所在目录
		dir := filepath.Dir(filePath)

		// 创建目录（包括所有必要父目录）
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}

		// 写入默认数据到文件
		if err := os.WriteFile(filePath, defaultData, filePerm); err != nil {
			return fmt.Errorf("failed to write default data: %w", err)
		}
	}
	return nil
}

// LoadOrCreateYAML 加载或创建 YAML 配置文件
// path: 配置文件路径
// defaultConfig: 默认配置对象（将作为模板创建新配置）
// dirPerm: 目录权限 (默认 0755)
// filePerm: 文件权限 (默认 0644)
func LoadOrCreateYAML(path string, defaultConfig interface{}, dirPerm, filePerm os.FileMode) (interface{}, error) {
	// 确保传入的是指针类型
	// if reflect.ValueOf(defaultConfig).Kind() != reflect.Ptr {
	// 	return nil, fmt.Errorf("defaultConfig must be a pointer")
	// }

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 创建目录结构
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		// 序列化默认配置为 YAML
		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default config: %w", err)
		}

		// 写入配置文件
		if err := os.WriteFile(path, data, filePerm); err != nil {
			return nil, fmt.Errorf("failed to write default config: %w", err)
		}

		// 返回默认配置
		return defaultConfig, nil
	}

	// 读取现有配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 创建目标配置对象的副本（避免修改原始默认配置）
	// 注意：这里假设 defaultConfig 是结构体指针
	config := cloneConfig(defaultConfig)

	// 反序列化 YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// cloneConfig 创建配置对象的深拷贝
func cloneConfig(original interface{}) interface{} {
	// 在实际应用中，可以使用更高效的克隆方法
	// 这里使用序列化/反序列化作为通用解决方案
	data, err := yaml.Marshal(original)
	if err != nil {
		panic("failed to clone config: " + err.Error())
	}

	// 创建相同类型的新对象
	t := reflect.TypeOf(original)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	newObj := reflect.New(t).Interface()

	if err := yaml.Unmarshal(data, newObj); err != nil {
		panic("failed to clone config: " + err.Error())
	}

	return newObj
}

func New(s ...string) *Config {
	cfg := configSingleton.Get(func() *Config {
		var filename string
		if len(s) == 0 {
			filename = "configs/config.yaml"
		} else {
			filename = strings.Join(s, "")
		}

		defaultConfig := &Config{}
		err := defaultvalue.SetDefaults(defaultConfig)
		if err != nil {
			panic(fmt.Sprintf("配置初始化默认数据失败: %v", err))
		}

		// 加载或创建配置
		config, err := LoadOrCreateYAML(
			filename,
			defaultConfig,
			0755, // 目录权限
			0644, // 文件权限
		)

		if err != nil {
			panic(fmt.Sprintf("配置初始化失败: %v", err))
		}

		// 类型断言获取具体配置
		cfg2 := config.(*Config)
		return cfg2
	})
	return cfg
}
