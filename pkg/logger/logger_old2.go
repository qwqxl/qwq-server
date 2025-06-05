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
//	Name           string
//	Group          string
//	SaveToFile     bool
//	PrintToConsole bool
//	FilePath       string
//	MaxSize        int64  // 字节
//	FilePattern    string // 支持 {name}, {date}, {index}
//	RotateCount    int
//	currentIndex   int
//	currentSize    int64
//	file           *os.File
//	mu             sync.Mutex
//	// color
//}
//
//type LoggerManager struct {
//	instances map[string]*LogConfig
//	groups    map[string]*LogConfig
//	defaults  *LogConfig
//	mu        sync.RWMutex
//}
//
//var (
//	manager     *LoggerManager
//	managerOnce sync.Once
//)
//
//// 初始化管理器
//func Init(defaultSavePath string) *LoggerManager {
//	managerOnce.Do(func() {
//		defaultCfg := &LogConfig{
//			SaveToFile:     true,
//			PrintToConsole: true,
//			MaxSize:        100 * 1024 * 1024, // 100MB
//			FilePattern:    "{name}-{date}-{index}.log",
//			RotateCount:    10,
//			FilePath:       defaultSavePath,
//		}
//
//		manager = &LoggerManager{
//			instances: make(map[string]*LogConfig),
//			groups:    make(map[string]*LogConfig),
//			defaults:  defaultCfg,
//		}
//	})
//	return manager
//}
//
//// 配置选项类型
//type ConfigOption func(*LogConfig)
//
//// 公共API
//func Info(name string, format string, args ...interface{}) {
//	manager.mu.RLock()
//	defer manager.mu.RUnlock()
//
//	if cfg, ok := manager.instances[name]; ok {
//		cfg.log(format, args...)
//	} else {
//		// 使用默认配置创建临时实例
//		newCfg := manager.defaults.copy()
//		newCfg.Name = name
//		newCfg.log(format, args...)
//	}
//}
//
//// 创建或获取配置实例
//func Configure(name string, options ...ConfigOption) *LogConfig {
//	manager.mu.Lock()
//	defer manager.mu.Unlock()
//
//	if cfg, exists := manager.instances[name]; exists {
//		return cfg
//	}
//
//	// 创建新配置
//	newCfg := manager.defaults.copy()
//	newCfg.Name = name
//
//	// 应用配置选项
//	for _, opt := range options {
//		opt(newCfg)
//	}
//
//	// 处理群组配置
//	if newCfg.Group != "" {
//		if groupCfg, ok := manager.groups[newCfg.Group]; ok {
//			newCfg.mergeGroup(groupCfg)
//		}
//	}
//
//	manager.instances[name] = newCfg
//	return newCfg
//}
//
//// 日志写入核心方法
//func (lc *LogConfig) log(format string, args ...interface{}) {
//	lc.mu.Lock()
//	defer lc.mu.Unlock()
//
//	msg := fmt.Sprintf("[%s] %s - %s\n",
//		time.Now().Format("2006-01-02 15:04:05"),
//		strings.ToUpper(lc.Name),
//		fmt.Sprintf(format, args...),
//	)
//
//	// 控制台输出
//	if lc.PrintToConsole {
//		fmt.Print(msg)
//	}
//
//	// 文件输出
//	if lc.SaveToFile {
//		if err := lc.rotateIfNeeded(); err != nil {
//			fmt.Printf("日志轮换失败: %v\n", err)
//			return
//		}
//
//		if _, err := lc.file.WriteString(msg); err != nil {
//			fmt.Printf("写入日志失败: %v\n", err)
//		}
//		lc.currentSize += int64(len(msg))
//	}
//}
//
//// rotateIfNeeded 检查是否需要轮换日志文件，如果需要则进行轮换操作。
//func (lc *LogConfig) buildReplacements(t time.Time) map[string]string {
//	return map[string]string{
//		// 全称占位符
//		"{name}":  lc.Name,
//		"{group}": lc.Group,
//		"{index}": fmt.Sprintf("%d", lc.currentIndex),
//		"{date}":  t.Format("20060102"),
//		"{time}":  t.Format("150405"),
//		"{year}":  t.Format("2006"),
//		"{month}": t.Format("01"),
//		"{day}":   t.Format("02"),
//		"{hour}":  t.Format("15"),
//		"{min}":   t.Format("04"),
//		"{sec}":   t.Format("05"),
//		"{ts}":    fmt.Sprintf("%d", t.Unix()),
//		"{pid}":   fmt.Sprintf("%d", os.Getpid()),
//
//		// 简写占位符
//		"{_i}": fmt.Sprintf("%d", lc.currentIndex),
//		"{_d}": t.Format("20060102"),
//		"{_t}": t.Format("150405"),
//		"{_y}": t.Format("06"),
//		"{_m}": t.Format("01"),
//		"{_D}": t.Format("02"),
//		"{_H}": t.Format("15"),
//		"{_s}": t.Format("05"),
//	}
//}
//
//// 生成带格式的文件名
//func (lc *LogConfig) generateFileName() string {
//	now := time.Now()
//	replacements := lc.buildReplacements(now)
//
//	fileName := lc.FilePattern
//	fileName = replaceExactMatches(fileName, replacements)
//	fileName = replaceShortNotation(fileName, replacements)
//	fileName = filepath.Clean(fileName)
//
//	return fileName
//}
//
//// 精确替换全称占位符
//func replaceExactMatches(s string, reps map[string]string) string {
//	for k, v := range reps {
//		if strings.HasPrefix(k, "{") {
//			s = strings.ReplaceAll(s, k, v)
//		}
//	}
//	return s
//}
//
//// 替换简写占位符（带下划线前缀）
//func replaceShortNotation(s string, reps map[string]string) string {
//	for k, v := range reps {
//		if strings.HasPrefix(k, "{_") {
//			shortKey := "{" + strings.TrimPrefix(k, "{_")
//			s = strings.ReplaceAll(s, shortKey, v)
//		}
//	}
//	return s
//}
//
//// 修改文件轮换方法
//func (lc *LogConfig) rotateIfNeeded() error {
//	if lc.file == nil || lc.currentSize >= lc.MaxSize {
//		if lc.file != nil {
//			lc.file.Close()
//		}
//
//		// 生成完整文件路径
//		fullPath := lc.generateFileName()
//		dir := filepath.Dir(fullPath)
//
//		// 先创建目录
//		if err := os.MkdirAll(dir, 0755); err != nil {
//			return fmt.Errorf("创建目录失败 [%s]: %w", dir, err)
//		}
//
//		// 打开文件
//		file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
//		if err != nil {
//			return fmt.Errorf("打开文件失败 [%s]: %w", fullPath, err)
//		}
//
//		lc.file = file
//		lc.currentSize = 0
//		lc.currentIndex = (lc.currentIndex % lc.RotateCount) + 1
//	}
//	return nil
//}
//
//// 以下为配置选项函数
//func WithGroup(group string) ConfigOption {
//	return func(c *LogConfig) {
//		c.Group = group
//	}
//}
//
//func WithFileOutput(enabled bool) ConfigOption {
//	return func(c *LogConfig) {
//		c.SaveToFile = enabled
//	}
//}
//
//func WithConsoleOutput(enabled bool) ConfigOption {
//	return func(c *LogConfig) {
//		c.PrintToConsole = enabled
//	}
//}
//
//func WithFilePath(path string) ConfigOption {
//	return func(c *LogConfig) {
//		c.FilePath = path
//	}
//}
//
//func WithMaxSize(sizeMB int64) ConfigOption {
//	return func(c *LogConfig) {
//		c.MaxSize = sizeMB * 1024 * 1024
//	}
//}
//
//func WithFilePattern(pattern string) ConfigOption {
//	return func(c *LogConfig) {
//		c.FilePattern = pattern
//	}
//}
//
//func WithRotateCount(count int) ConfigOption {
//	return func(c *LogConfig) {
//		c.RotateCount = count
//	}
//}
//
//// 辅助方法
//func (lc *LogConfig) copy() *LogConfig {
//	return &LogConfig{
//		SaveToFile:     lc.SaveToFile,
//		PrintToConsole: lc.PrintToConsole,
//		FilePath:       lc.FilePath,
//		MaxSize:        lc.MaxSize,
//		FilePattern:    lc.FilePattern,
//		RotateCount:    lc.RotateCount,
//	}
//}
//
//func (lc *LogConfig) mergeGroup(groupCfg *LogConfig) {
//	if lc.FilePath == "" {
//		lc.FilePath = groupCfg.FilePath
//	}
//	if lc.MaxSize == 0 {
//		lc.MaxSize = groupCfg.MaxSize
//	}
//	if lc.FilePattern == "" {
//		lc.FilePattern = groupCfg.FilePattern
//	}
//	if lc.RotateCount == 0 {
//		lc.RotateCount = groupCfg.RotateCount
//	}
//}
