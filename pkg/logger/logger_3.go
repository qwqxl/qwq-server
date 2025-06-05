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
//// 日志级别类型
//type LogLevel int
//
//const (
//	LevelDebug LogLevel = iota
//	LevelInfo
//	LevelWarn
//	LevelError
//)
//
//// 对外暴露的日志函数
//var (
//	Debug = createLogFunc(LevelDebug)
//	Info  = createLogFunc(LevelInfo)
//	Warn  = createLogFunc(LevelWarn)
//	Error = createLogFunc(LevelError)
//)
//
//func createLogFunc(level LogLevel) func(string, string, ...interface{}) {
//	return func(name string, format string, args ...interface{}) {
//		manager.mu.RLock()
//		defer manager.mu.RUnlock()
//
//		if cfg, ok := manager.instances[name]; ok {
//			cfg.log(level, format, args...)
//		} else {
//			newCfg := manager.defaults.copy()
//			newCfg.Name = name
//			newCfg.log(level, format, args...)
//		}
//	}
//}
//
//func (l LogLevel) String() string {
//	switch l {
//	case LevelDebug:
//		return "DEBUG"
//	case LevelInfo:
//		return "INFO"
//	case LevelWarn:
//		return "WARN"
//	case LevelError:
//		return "ERROR"
//	default:
//		return "UNKNOWN"
//	}
//}
//
//// 日志配置结构体
//type LogConfig struct {
//	Name           string
//	Group          string
//	SaveToFile     bool
//	PrintToConsole bool
//	FilePath       string
//	MaxSize        int64 // 字节
//	FilePattern    string
//	RotateCount    int
//	currentIndex   int
//	currentSize    int64
//	file           *os.File
//	mu             sync.Mutex
//	EnableColors   bool
//	Colors         map[LogLevel]string
//	// 队列
//	asyncQueue     chan string    // 异步消息队列
//	asyncQueueSize int            // 队列缓冲区大小
//	asyncWG        sync.WaitGroup // 等待队列消费完成
//	rotateIndex    any
//}
//
//// 构建用于替换路径和文件名中占位符的映射
//func (lc *LogConfig) buildReplacements(t time.Time) map[string]string {
//	return map[string]string{
//		"{name}":  lc.Name,
//		"{group}": lc.Group,
//		"{year}":  fmt.Sprintf("%04d", t.Year()),
//		"{month}": fmt.Sprintf("%02d", int(t.Month())),
//		"{day}":   fmt.Sprintf("%02d", t.Day()),
//		"{_d}":    t.Format("2006-01-02"),
//		"{_H}":    fmt.Sprintf("%02d", t.Hour()),
//		"{_m}":    fmt.Sprintf("%02d", t.Minute()),
//		"{_s}":    fmt.Sprintf("%02d", t.Second()),
//		"{_i}":    fmt.Sprintf("%d", lc.rotateIndex), // 轮转编号
//	}
//}
//
//// 在 generateFileName 中修复目录路径处理
//func (lc *LogConfig) generateFileName() string {
//	now := time.Now()
//	replacements := lc.buildReplacements(now)
//
//	// 替换目录路径中的占位符
//	dirPath := lc.FilePath
//	for placeholder, value := range replacements {
//		dirPath = strings.ReplaceAll(dirPath, placeholder, value)
//	}
//
//	// 替换文件名中的占位符
//	fileName := lc.FilePattern
//	for placeholder, value := range replacements {
//		fileName = strings.ReplaceAll(fileName, placeholder, value)
//	}
//
//	// 确保目录路径是绝对路径
//	if !filepath.IsAbs(dirPath) {
//		// 使用基础路径（Init设置的路径）
//		basePath := manager.defaults.FilePath
//		dirPath = filepath.Join(basePath, dirPath)
//	}
//
//	// 创建目录（如果不存在）
//	if err := os.MkdirAll(dirPath, 0755); err != nil {
//		fmt.Printf("创建目录失败: %v\n", err)
//	}
//
//	return filepath.Join(dirPath, fileName)
//}
//
//// 改进 mergeConfigs 函数
//func mergeConfigs(base, override *LogConfig) *LogConfig {
//	merged := base.copy()
//
//	// 应用覆盖设置（如果覆盖配置中有设置）
//	if override.SaveToFile {
//		merged.SaveToFile = true
//	}
//	if override.PrintToConsole {
//		merged.PrintToConsole = true
//	}
//	if override.FilePath != "" {
//		merged.FilePath = override.FilePath
//	}
//	if override.FilePattern != "" {
//		merged.FilePattern = override.FilePattern
//	}
//	if override.MaxSize != 0 {
//		merged.MaxSize = override.MaxSize
//	}
//	if override.RotateCount != 0 {
//		merged.RotateCount = override.RotateCount
//	}
//	if override.EnableColors {
//		merged.EnableColors = true
//	}
//
//	// 合并颜色配置
//	for level, color := range override.Colors {
//		merged.Colors[level] = color
//	}
//
//	// 保留实例特定的属性
//	merged.Name = override.Name
//	merged.Group = override.Group
//
//	return merged
//}
//
//// 在 rotateIfNeeded 中移除目录创建逻辑
//func (lc *LogConfig) rotateIfNeeded() error {
//	if lc.file == nil || lc.currentSize >= lc.MaxSize {
//		if lc.file != nil {
//			lc.file.Close()
//		}
//
//		fullPath := lc.generateFileName() // 这里已经包含目录创建
//
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
//// 新增配置选项
//func WithAsync(enabled bool, queueSize int) ConfigOption {
//	return func(c *LogConfig) {
//		if enabled {
//			c.asyncQueue = make(chan string, queueSize)
//			c.asyncQueueSize = queueSize
//			go c.asyncWriter()
//		}
//	}
//}
//
//// 异步写入器
//func (lc *LogConfig) asyncWriter() {
//	lc.asyncWG.Add(1)
//	defer lc.asyncWG.Done()
//
//	for msg := range lc.asyncQueue {
//		// 执行实际的日志写入
//		lc.writeToFile(msg)
//	}
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
//			EnableColors:   true,
//			Colors: map[LogLevel]string{
//				LevelDebug: "\x1b[37m",
//				LevelInfo:  "\x1b[32m",
//				LevelWarn:  "\x1b[33m",
//				LevelError: "\x1b[31m",
//			},
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
//func logWithLevel(name string, level LogLevel, format string, args ...interface{}) {
//	manager.mu.RLock()
//	defer manager.mu.RUnlock()
//
//	if cfg, ok := manager.instances[name]; ok {
//		cfg.log(level, format, args...)
//	} else {
//		newCfg := manager.defaults.copy()
//		newCfg.Name = name
//		newCfg.log(level, format, args...)
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
//	newCfg := manager.defaults.copy()
//	newCfg.Name = name
//
//	for _, opt := range options {
//		opt(newCfg)
//	}
//
//	if newCfg.Group != "" {
//		if groupCfg, ok := manager.groups[newCfg.Group]; ok {
//			// 深度合并组配置
//			newCfg = mergeConfigs(groupCfg, newCfg)
//		}
//		//if groupCfg, ok := manager.groups[newCfg.Group]; ok {
//		//	newCfg.mergeGroup(groupCfg)
//		//}
//	}
//
//	manager.instances[name] = newCfg
//	return newCfg
//}
//
//func mergeConfigs2(base, override *LogConfig) *LogConfig {
//	merged := base.copy()
//	// 应用覆盖设置
//	if override.SaveToFile {
//		merged.SaveToFile = true
//	}
//	if override.FilePath != "" {
//		merged.FilePath = override.FilePath
//	}
//	// ...其他字段类似处理
//	return merged
//}
//
//// 日志写入核心方法
//
//// 修改日志记录方法
//func (lc *LogConfig) log(level LogLevel, format string, args ...interface{}) {
//	// ...构造消息部分保持不变...
//	lc.mu.Lock()
//	defer lc.mu.Unlock()
//
//	formattedMsg := fmt.Sprintf(format, args...)
//	timeStr := time.Now().Format("2006-01-02 15:04:05")
//	levelStr := strings.ToUpper(level.String())
//	rawMsg := fmt.Sprintf("[%s] [%s] %s - %s\n", timeStr, levelStr, lc.Name, formattedMsg)
//
//	// 控制台输出
//	if lc.PrintToConsole {
//		var outputMsg string
//		if lc.EnableColors {
//			colorCode := lc.Colors[level]
//			resetCode := "\x1b[0m"
//			outputMsg = colorCode + rawMsg + resetCode
//		} else {
//			outputMsg = rawMsg
//		}
//		fmt.Print(outputMsg)
//	}
//
//	// 异步处理逻辑
//	if lc.asyncQueue != nil {
//		select {
//		case lc.asyncQueue <- rawMsg:
//			// 成功写入队列
//		default:
//			// 队列满时直接写入（保证不丢失日志）
//			lc.writeToFile(rawMsg)
//		}
//		return
//	}
//
//	// 同步写入
//	lc.writeToFile(rawMsg)
//}
//
//// 独立出实际文件写入方法
//func (lc *LogConfig) writeToFile(rawMsg string) {
//	if err := lc.rotateIfNeeded(); err != nil {
//		fmt.Printf("日志轮换失败: %v\n", err)
//		return
//	}
//
//	if _, err := lc.file.WriteString(rawMsg); err != nil {
//		fmt.Printf("写入日志失败: %v\n", err)
//	}
//	lc.currentSize += int64(len(rawMsg))
//}
//
//// 新增关闭方法
//func (lc *LogConfig) Close() {
//	if lc.asyncQueue != nil {
//		close(lc.asyncQueue)
//		lc.asyncWG.Wait() // 等待队列消费完成
//	}
//	if lc.file != nil {
//		lc.file.Close()
//	}
//}
//
//// 以下辅助方法与之前相同，需添加颜色相关处理
//func (lc *LogConfig) rotateIfNeeded2() error {
//	if lc.file == nil || lc.currentSize >= lc.MaxSize {
//		if lc.file != nil {
//			lc.file.Close()
//		}
//
//		fullPath := lc.generateFileName()
//		dir := filepath.Dir(fullPath)
//
//		if err := os.MkdirAll(dir, 0755); err != nil {
//			return fmt.Errorf("创建目录失败 [%s]: %w", dir, err)
//		}
//
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
//// ConfigureGroup 配置日志组
//func ConfigureGroup(group string, options ...ConfigOption) {
//	manager.mu.Lock()
//	defer manager.mu.Unlock()
//
//	// 创建或更新组配置
//	if groupCfg, exists := manager.groups[group]; exists {
//		// 更新现有组配置
//		for _, opt := range options {
//			opt(groupCfg)
//		}
//	} else {
//		// 创建新组配置
//		newCfg := manager.defaults.copy()
//		newCfg.Group = group
//		for _, opt := range options {
//			opt(newCfg)
//		}
//		manager.groups[group] = newCfg
//	}
//}
//
//func (lc *LogConfig) generateFileName2() string {
//	now := time.Now()
//	replacements := map[string]string{
//		// 基础占位符
//		"{name}":  lc.Name,
//		"{group}": lc.Group,
//		"{index}": fmt.Sprintf("%d", lc.currentIndex),
//		"{pid}":   fmt.Sprintf("%d", os.Getpid()),
//
//		// 完整时间格式
//		"{date}":  now.Format("20060102"),
//		"{time}":  now.Format("150405"),
//		"{year}":  now.Format("2006"),
//		"{month}": now.Format("01"), // 数字月份 (01-12)
//		"{day}":   now.Format("02"),
//		"{hour}":  now.Format("15"), // 24小时制
//		"{min}":   now.Format("04"), // 分钟
//		"{sec}":   now.Format("05"),
//		"{ts}":    fmt.Sprintf("%d", now.Unix()),
//
//		// 短格式别名 (修正 {_m} -> 分钟)
//		"{_i}": fmt.Sprintf("%d", lc.currentIndex), // 索引
//		"{_d}": now.Format("20060102"),             // 完整日期
//		"{_t}": now.Format("150405"),               // 完整时间
//		"{_y}": now.Format("06"),                   // 两位年份
//		"{_M}": now.Format("01"),                   // 月份 (大写 M)
//		"{_m}": now.Format("04"),                   // 分钟 (小写 m) <- 关键修复
//		"{_D}": now.Format("02"),                   // 日
//		"{_H}": now.Format("15"),                   // 小时
//		"{_s}": now.Format("05"),                   // 秒
//	}
//
//	fileName := lc.FilePattern
//	for placeholder, value := range replacements {
//		fileName = strings.ReplaceAll(fileName, placeholder, value)
//	}
//
//	return filepath.Join(lc.FilePath, fileName)
//}
//
//// 配置选项函数
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
//func WithColor(level LogLevel, colorCode string) ConfigOption {
//	return func(c *LogConfig) {
//		c.Colors[level] = colorCode
//	}
//}
//
//func WithEnableColors(enable bool) ConfigOption {
//	return func(c *LogConfig) {
//		c.EnableColors = enable
//	}
//}
//
//// 辅助方法
//func (lc *LogConfig) copy() *LogConfig {
//	newColors := make(map[LogLevel]string)
//	for k, v := range lc.Colors {
//		newColors[k] = v
//	}
//	return &LogConfig{
//		Name:           lc.Name,
//		Group:          lc.Group,
//		SaveToFile:     lc.SaveToFile,
//		PrintToConsole: lc.PrintToConsole,
//		FilePath:       lc.FilePath,
//		MaxSize:        lc.MaxSize,
//		FilePattern:    lc.FilePattern,
//		RotateCount:    lc.RotateCount,
//		currentIndex:   lc.currentIndex,
//		currentSize:    lc.currentSize,
//		file:           nil,
//		mu:             sync.Mutex{},
//		EnableColors:   lc.EnableColors,
//		Colors:         newColors,
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
//	if !lc.EnableColors {
//		lc.EnableColors = groupCfg.EnableColors
//	}
//	for level, color := range groupCfg.Colors {
//		if _, exists := lc.Colors[level]; !exists {
//			lc.Colors[level] = color
//		}
//	}
//}
