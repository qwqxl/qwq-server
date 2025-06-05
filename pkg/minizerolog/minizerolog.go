package minizerolog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ======== 配置结构体 ========

type LogConfig struct {
	Group           string
	Async           bool
	AsyncQueueSize  int
	RotateCount     int    // 最大轮转文件数量
	FilePathPattern string // 路径模板，如 ./logs/{name}/{year}/{month}/{day}
	FileNamePattern string // 文件名模板，如 {name}_{_H}{_m}{_s}-{_i}.log
	MaxFileSizeMB   int    // 单文件最大大小MB，超过轮转
	ConsoleOutput   bool
	Level           Level
	Colors          map[Level]string // ANSI颜色码
}

func DefaultConfig() LogConfig {
	return LogConfig{
		Async:           true,
		AsyncQueueSize:  1000,
		RotateCount:     5,
		FilePathPattern: "./logs/{name}/{year}/{month}/{day}",
		FileNamePattern: "{name}_{_H}{_m}{_s}-{_i}.log",
		MaxFileSizeMB:   50,
		ConsoleOutput:   true,
		Level:           DEBUG,
		Colors: map[Level]string{
			DEBUG: "\033[36m", // 青色
			INFO:  "\033[32m", // 绿色
			WARN:  "\033[33m", // 黄色
			ERROR: "\033[31m", // 红色
		},
	}
}

// ======== 异步队列 ========

type logEntry struct {
	level   Level
	time    time.Time
	message string
}

type asyncQueue struct {
	ch      chan logEntry
	wg      sync.WaitGroup
	quit    chan struct{}
	handler func(entry logEntry)
}

func newAsyncQueue(size int, handler func(entry logEntry)) *asyncQueue {
	q := &asyncQueue{
		ch:      make(chan logEntry, size),
		quit:    make(chan struct{}),
		handler: handler,
	}
	q.wg.Add(1)
	go q.loop()
	return q
}

func (q *asyncQueue) loop() {
	defer q.wg.Done()
	for {
		select {
		case entry := <-q.ch:
			q.handler(entry)
		case <-q.quit:
			for entry := range q.ch {
				q.handler(entry)
			}
			return
		}
	}
}

func (q *asyncQueue) push(entry logEntry) {
	select {
	case q.ch <- entry:
	default:
		// 队列满丢弃
	}
}

func (q *asyncQueue) close() {
	close(q.quit)
	q.wg.Wait()
}

// ======== Logger 实例 ========

type Logger struct {
	name        string
	config      LogConfig
	asyncQueue  *asyncQueue
	fileHandle  *os.File
	currentSize int64
	fileIndex   int
	mu          sync.Mutex
}

func newLogger(name string, cfg LogConfig) *Logger {
	logger := &Logger{
		name:   name,
		config: cfg,
	}

	err := logger.openNewFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Logger open file err: %v\n", err)
	}

	if cfg.Async {
		logger.asyncQueue = newAsyncQueue(cfg.AsyncQueueSize, logger.output)
	}

	return logger
}

// 打开一个新的日志文件（轮转）
func (l *Logger) openNewFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fileHandle != nil {
		l.fileHandle.Close()
	}

	dir := l.resolvePathPattern()
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// 找文件名，轮转编号递增
	for i := 0; i < l.config.RotateCount; i++ {
		fileName := l.resolveFileNamePattern(i)
		fullPath := filepath.Join(dir, fileName)
		info, err := os.Stat(fullPath)
		if err != nil || info.Size() < int64(l.config.MaxFileSizeMB)*1024*1024 {
			l.fileIndex = i
			fh, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			l.fileHandle = fh
			if err == nil && info != nil {
				l.currentSize = info.Size()
			} else {
				l.currentSize = 0
			}
			return nil
		}
	}
	// 如果都满了，覆盖第一个
	fileName := l.resolveFileNamePattern(0)
	fullPath := filepath.Join(dir, fileName)
	fh, err := os.OpenFile(fullPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.fileHandle = fh
	l.fileIndex = 0
	l.currentSize = 0
	return nil
}

// resolveFileNamePattern 根据配置文件解析文件名
func (l *Logger) resolvePathPattern() string {
	now := time.Now()
	s := l.config.FilePathPattern

	// 定义支持的占位符
	replacements := map[string]string{
		"{name}":    l.name,                            // 日志名称
		"{group}":   l.config.Group,                    // 日志组名称
		"{year}":    fmt.Sprintf("%04d", now.Year()),   // 年份
		"{month}":   fmt.Sprintf("%02d", now.Month()),  // 月份
		"{day}":     fmt.Sprintf("%02d", now.Day()),    // 日
		"{hour}":    fmt.Sprintf("%02d", now.Hour()),   // 小时
		"{min}":     fmt.Sprintf("%02d", now.Minute()), // 分钟
		"{sec}":     fmt.Sprintf("%02d", now.Second()), // 秒
		"{weekday}": now.Weekday().String(),            // 星期几

		// 简写
		"{yyyy}": fmt.Sprintf("%04d", now.Year()),   // 年
		"{mm}":   fmt.Sprintf("%02d", now.Month()),  // 月
		"{dd}":   fmt.Sprintf("%02d", now.Day()),    // 日
		"{hh}":   fmt.Sprintf("%02d", now.Hour()),   // 小时
		"{ii}":   fmt.Sprintf("%02d", now.Minute()), // 分钟
		"{ss}":   fmt.Sprintf("%02d", now.Second()), // 秒
	}

	// 替换所有占位符
	for key, val := range replacements {
		s = strings.ReplaceAll(s, key, val)
	}

	// 如果不是绝对路径，加上 baseLogDir 前缀
	if !filepath.IsAbs(s) {
		s = filepath.Join(manager.baseLogDir, s)
	}

	return s
}

func (l *Logger) resolvePattern(pattern string, index ...int) string {
	now := time.Now()

	replacements := map[string]string{
		"{name}":    l.name,                            // 日志名称
		"{group}":   l.config.Group,                    // 日志组名称
		"{year}":    fmt.Sprintf("%04d", now.Year()),   // 年份
		"{month}":   fmt.Sprintf("%02d", now.Month()),  // 月份
		"{day}":     fmt.Sprintf("%02d", now.Day()),    // 日
		"{hour}":    fmt.Sprintf("%02d", now.Hour()),   // 小时
		"{min}":     fmt.Sprintf("%02d", now.Minute()), // 分钟
		"{sec}":     fmt.Sprintf("%02d", now.Second()), // 秒
		"{index}":   fmt.Sprintf("%d", index),          // 文件序号
		"{i}":       fmt.Sprintf("%d", index),          // 文件序号
		"{weekday}": now.Weekday().String(),            // 星期几

		// 简写
		"{yyyy}": fmt.Sprintf("%04d", now.Year()),   // 年
		"{mm}":   fmt.Sprintf("%02d", now.Month()),  // 月
		"{dd}":   fmt.Sprintf("%02d", now.Day()),    // 日
		"{hh}":   fmt.Sprintf("%02d", now.Hour()),   // 小时
		"{ii}":   fmt.Sprintf("%02d", now.Minute()), // 分钟
		"{ss}":   fmt.Sprintf("%02d", now.Second()), // 秒

		//"{name}":      l.name,
		//"{group}":     l.config.Group,
		//"{year}":      fmt.Sprintf("%04d", now.Year()),
		//"{month}":     fmt.Sprintf("%02d", now.Month()),
		//"{day}":       fmt.Sprintf("%02d", now.Day()),
		//"{hour}":      fmt.Sprintf("%02d", now.Hour()),
		//"{min}":       fmt.Sprintf("%02d", now.Minute()),
		//"{sec}":       fmt.Sprintf("%02d", now.Second()),
		//"{index}":     fmt.Sprintf("%d", index),
		//"{i}":         fmt.Sprintf("%d", index),
		//"{weekday}":   now.Weekday().String(),
		//"{yyyy}":      fmt.Sprintf("%04d", now.Year()),
		//"{mm}":        fmt.Sprintf("%02d", now.Month()),
		//"{dd}":        fmt.Sprintf("%02d", now.Day()),
		//"{hh}":        fmt.Sprintf("%02d", now.Hour()),
		//"{ii}":        fmt.Sprintf("%02d", now.Minute()),
		//"{ss}":        fmt.Sprintf("%02d", now.Second()),
		//"{timestamp}": fmt.Sprintf("%d", now.Unix()),
	}

	for key, val := range replacements {
		pattern = strings.ReplaceAll(pattern, key, val)
	}

	return pattern
}

func (l *Logger) resolveFileNamePattern(index int) string {
	now := time.Now()
	s := l.config.FileNamePattern

	// 新增更多占位符
	replacements := map[string]string{
		"{name}":    l.name,                            // 日志名称
		"{group}":   l.config.Group,                    // 日志组名称
		"{year}":    fmt.Sprintf("%04d", now.Year()),   // 年份
		"{month}":   fmt.Sprintf("%02d", now.Month()),  // 月份
		"{day}":     fmt.Sprintf("%02d", now.Day()),    // 日
		"{hour}":    fmt.Sprintf("%02d", now.Hour()),   // 小时
		"{min}":     fmt.Sprintf("%02d", now.Minute()), // 分钟
		"{sec}":     fmt.Sprintf("%02d", now.Second()), // 秒
		"{index}":   fmt.Sprintf("%d", index),          // 文件序号
		"{weekday}": now.Weekday().String(),            // 星期几

		// 简写
		"{yyyy}": fmt.Sprintf("%04d", now.Year()),   // 年
		"{mm}":   fmt.Sprintf("%02d", now.Month()),  // 月
		"{dd}":   fmt.Sprintf("%02d", now.Day()),    // 日
		"{hh}":   fmt.Sprintf("%02d", now.Hour()),   // 小时
		"{ii}":   fmt.Sprintf("%02d", now.Minute()), // 分钟
		"{ss}":   fmt.Sprintf("%02d", now.Second()), // 秒
	}

	for key, val := range replacements {
		s = strings.ReplaceAll(s, key, val)
	}

	return s
}

func (l *Logger) output(entry logEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	line := fmt.Sprintf("[%s] %s %s\n", entry.time.Format(time.RFC3339), entry.level.String(), entry.message)

	// 写文件
	if l.fileHandle != nil {
		n, err := l.fileHandle.WriteString(line)
		if err == nil {
			l.currentSize += int64(n)
		}
		if l.currentSize > int64(l.config.MaxFileSizeMB)*1024*1024 {
			// 轮转文件
			l.fileIndex = (l.fileIndex + 1) % l.config.RotateCount
			l.openNewFile()
		}
	}

	// 控制台输出
	if l.config.ConsoleOutput {
		color, ok := l.config.Colors[entry.level]
		if !ok {
			color = "\033[0m"
		}
		reset := "\033[0m"
		fmt.Printf("%s%s%s", color, line, reset)
	}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.config.Level {
		return
	}
	entry := logEntry{
		level:   level,
		time:    time.Now(),
		message: fmt.Sprintf(format, args...),
	}

	if l.config.Async && l.asyncQueue != nil {
		l.asyncQueue.push(entry)
	} else {
		l.output(entry)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) { l.log(DEBUG, format, args...) }
func (l *Logger) Info(format string, args ...interface{})  { l.log(INFO, format, args...) }
func (l *Logger) Warn(format string, args ...interface{})  { l.log(WARN, format, args...) }
func (l *Logger) Error(format string, args ...interface{}) { l.log(ERROR, format, args...) }

func (l *Logger) Close() {
	if l.asyncQueue != nil {
		l.asyncQueue.close()
	}
	l.mu.Lock()
	if l.fileHandle != nil {
		l.fileHandle.Close()
		l.fileHandle = nil
	}
	l.mu.Unlock()
}

func Close() {
	if manager == nil {
		return
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	for _, logger := range manager.loggers {
		logger.Close()
	}
	manager.loggers = make(map[string]*Logger)
}

// ======== LoggerManager 单例 ========

type LoggerManager struct {
	mu          sync.RWMutex
	loggers     map[string]*Logger
	groupConfig map[string]LogConfig
	defaultCfg  LogConfig
	baseLogDir  string //  日志根目录
}

var manager *LoggerManager
var once sync.Once

func Init(defaultLogDir string) {
	once.Do(func() {
		manager = &LoggerManager{
			baseLogDir:  defaultLogDir,
			loggers:     make(map[string]*Logger),
			groupConfig: make(map[string]LogConfig),
			defaultCfg:  DefaultConfig(),
		}
		manager.defaultCfg.FilePathPattern = defaultLogDir + "{name}/{year}/{month}/{day}"
	})
}

func Get(name string) *Logger {
	manager.mu.RLock()
	logger, ok := manager.loggers[name]
	manager.mu.RUnlock()
	if ok {
		return logger
	}
	// 没有则自动配置
	return Configure(name)
}

func (m *LoggerManager) ConfigureGroup(group string, options ...Option) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cfg := m.defaultCfg
	for _, opt := range options {
		opt(&cfg)
	}
	m.groupConfig[group] = cfg
}

func (m *LoggerManager) Configure(name string, options ...Option) *Logger {
	m.mu.Lock()
	defer m.mu.Unlock()

	cfg := m.defaultCfg

	// 查找 group
	var group string
	for _, opt := range options {
		if g, ok := extractGroup(opt); ok {
			group = g
		}
	}
	if group != "" {
		if gcfg, ok := m.groupConfig[group]; ok {
			cfg = gcfg
		}
	}

	// 应用所有配置覆盖
	for _, opt := range options {
		opt(&cfg)
	}

	logger := newLogger(name, cfg)
	m.loggers[name] = logger
	return logger
}

// Option 定义

type Option func(*LogConfig)

func WithGroup(g string) Option {
	return func(cfg *LogConfig) { cfg.Group = g }
}

func WithAsync(enabled bool, queueSize int) Option {
	return func(cfg *LogConfig) {
		cfg.Async = enabled
		cfg.AsyncQueueSize = queueSize
	}
}

func WithRotateCount(count int) Option {
	return func(cfg *LogConfig) { cfg.RotateCount = count }
}

func WithFilePath(pattern string) Option {
	return func(cfg *LogConfig) { cfg.FilePathPattern = pattern }
}

func WithFilePattern(pattern string) Option {
	return func(cfg *LogConfig) { cfg.FileNamePattern = pattern }
}

func WithMaxSize(mb int) Option {
	return func(cfg *LogConfig) { cfg.MaxFileSizeMB = mb }
}

func WithConsoleOutput(enabled bool) Option {
	return func(cfg *LogConfig) { cfg.ConsoleOutput = enabled }
}

func WithLogLevel(level Level) Option {
	return func(cfg *LogConfig) { cfg.Level = level }
}

func WithColors(colors map[Level]string) Option {
	return func(cfg *LogConfig) { cfg.Colors = colors }
}

func extractGroup(opt Option) (string, bool) {
	dummy := LogConfig{}
	opt(&dummy)
	if dummy.Group != "" {
		return dummy.Group, true
	}
	return "", false
}

func extractGroup2(opt Option) (string, bool) {
	marker := struct {
		Group string
	}{}
	opt(&LogConfig{
		Group: "___TEMP_MARKER_GROUP___",
	})
	if marker.Group != "" && marker.Group != "___TEMP_MARKER_GROUP___" {
		return marker.Group, true
	}
	return "", false
}

// 顶层方便调用

func ConfigureGroup(group string, options ...Option) {
	manager.ConfigureGroup(group, options...)
}

func Configure(name string, options ...Option) *Logger {
	return manager.Configure(name, options...)
}

func Debug(name, format string, args ...interface{}) {
	Get(name).Debug(format, args...)
}
func Info(name, format string, args ...interface{}) {
	Get(name).Info(format, args...)
}
func Warn(name, format string, args ...interface{}) {
	Get(name).Warn(format, args...)
}
func Error(name, format string, args ...interface{}) {
	Get(name).Error(format, args...)
}
