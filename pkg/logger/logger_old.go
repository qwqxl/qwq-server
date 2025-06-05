package logger

//
//import (
//	"fmt"
//	"os"
//	"path/filepath"
//	"strings"
//	"sync"
//	"time"
//)
//
//// 日志配置结构体
//type LogConfig struct {
//	Name         string
//	SaveFile     bool
//	FilePath     string
//	MaxSize      int64 // 字节
//	FilePattern  string
//	Group        string
//	RotateCount  int
//	currentIndex int
//	file         *os.File
//	size         int64
//	mu           sync.Mutex
//}
//
//// 日志管理器
//type LoggerManager struct {
//	configs    map[string]*LogConfig
//	groups     map[string]*LogConfig
//	mu         sync.RWMutex
//	defaultDir string
//}
//
//var (
//	manager     *LoggerManager
//	managerOnce sync.Once
//)
//
//// 初始化日志管理器
//func NewManager(defaultDir string) *LoggerManager {
//	managerOnce.Do(func() {
//		manager = &LoggerManager{
//			configs:    make(map[string]*LogConfig),
//			groups:     make(map[string]*LogConfig),
//			defaultDir: defaultDir,
//		}
//	})
//	return manager
//}
//
//// 获取或创建配置
//func (lm *LoggerManager) GetConfig(name string, options ...ConfigOption) *LogConfig {
//	lm.mu.Lock()
//	defer lm.mu.Unlock()
//
//	if cfg, exists := lm.configs[name]; exists {
//		return cfg
//	}
//
//	// 创建新配置
//	cfg := &LogConfig{
//		Name:        name,
//		SaveFile:    true,
//		MaxSize:     100 * 1024 * 1024, // 默认100MB
//		FilePattern: "{name}-{date}-{index}.log",
//		RotateCount: 10,
//		FilePath:    lm.defaultDir,
//	}
//
//	// 应用选项
//	for _, opt := range options {
//		opt(cfg)
//	}
//
//	// 处理群组配置
//	if cfg.Group != "" {
//		if groupCfg, ok := lm.groups[cfg.Group]; ok {
//			// 合并群组配置
//			cfg.FilePath = groupCfg.FilePath
//			cfg.MaxSize = groupCfg.MaxSize
//			cfg.FilePattern = groupCfg.FilePattern
//			cfg.RotateCount = groupCfg.RotateCount
//		}
//	}
//
//	lm.configs[name] = cfg
//	return cfg
//}
//
//// 配置选项类型
//type ConfigOption func(*LogConfig)
//
//// 配置选项函数
//func WithSaveFile(save bool) ConfigOption {
//	return func(c *LogConfig) {
//		c.SaveFile = save
//	}
//}
//
//func WithFilePath(path string) ConfigOption {
//	return func(c *LogConfig) {
//		c.FilePath = path
//	}
//}
//
//func WithMaxSize(size int64) ConfigOption {
//	return func(c *LogConfig) {
//		c.MaxSize = size
//	}
//}
//
//func WithFilePattern(pattern string) ConfigOption {
//	return func(c *LogConfig) {
//		c.FilePattern = pattern
//	}
//}
//
//func WithGroup(group string) ConfigOption {
//	return func(c *LogConfig) {
//		c.Group = group
//	}
//}
//
//func WithRotateCount(count int) ConfigOption {
//	return func(c *LogConfig) {
//		c.RotateCount = count
//	}
//}
//
//// 日志写入方法
//func (lc *LogConfig) Info(format string, args ...interface{}) {
//	lc.mu.Lock()
//	defer lc.mu.Unlock()
//
//	msg := fmt.Sprintf("[%s] INFO - %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
//
//	// 控制台输出
//	fmt.Print(msg)
//
//	// 文件记录
//	if lc.SaveFile {
//		if err := lc.rotateIfNeeded(); err != nil {
//			fmt.Printf("日志文件轮换失败: %v\n", err)
//			return
//		}
//
//		if _, err := lc.file.WriteString(msg); err != nil {
//			fmt.Printf("写入日志失败: %v\n", err)
//		}
//		lc.size += int64(len(msg))
//	}
//}
//
//// 文件轮换逻辑
//func (lc *LogConfig) rotateIfNeeded() error {
//	if lc.file == nil || lc.size >= lc.MaxSize {
//		if lc.file != nil {
//			lc.file.Close()
//		}
//
//		// 生成新文件名
//		newPath := lc.generateFilename()
//		if err := os.MkdirAll(lc.FilePath, 0755); err != nil {
//			return fmt.Errorf("创建目录失败: %w", err)
//		}
//
//		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
//		if err != nil {
//			return fmt.Errorf("打开文件失败: %w", err)
//		}
//
//		lc.file = file
//		lc.size = 0
//		lc.currentIndex = (lc.currentIndex % lc.RotateCount) + 1
//	}
//	return nil
//}
//
//// 生成带格式的文件名
//func (lc *LogConfig) generateFilename() string {
//	now := time.Now()
//	replacements := map[string]string{
//		"{name}":  lc.Name,
//		"{date}":  now.Format("20060102"),
//		"{index}": fmt.Sprintf("%d", lc.currentIndex),
//	}
//
//	fileName := lc.FilePattern
//	for k, v := range replacements {
//		fileName = strings.ReplaceAll(fileName, k, v)
//	}
//
//	return filepath.Join(lc.FilePath, fileName)
//}
//
//// 公共API
//func Info(name string, format string, args ...interface{}) {
//	manager.mu.RLock()
//	defer manager.mu.RUnlock()
//
//	if cfg, ok := manager.configs[name]; ok {
//		cfg.Info(format, args...)
//	} else {
//		fmt.Printf("未找到日志配置: %s\n", name)
//	}
//}
//
//// 示例用法
//func ExampleUsage() {
//	// 初始化
//	NewManager("./logs")
//
//	// 注册mysql日志配置
//	manager.GetConfig("mysql",
//		WithGroup("database"),
//		WithMaxSize(10*1024*1024), // 10MB
//		WithFilePattern("{name}-{index}.log"),
//		WithRotateCount(5),
//	)
//
//	// 注册群组配置
//	manager.GetConfig("database-group",
//		WithGroup("database"),
//		WithFilePath("/var/logs/db"),
//		WithRotateCount(10),
//	)
//
//	// 使用日志
//	Info("mysql", "Connection established. Host: %s, Port: %d", "127.0.0.1", 3306)
//}
