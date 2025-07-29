// qwqlog.go
package qwqlog

import (
	"bytes"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"qwqserver/pkg/util/singleton"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 日志级别类型
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	DefaultLogLevel = INFO
	// default
)

var LogLevelStrName = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// StringToLevel 将字符串转换为对应的 LogLevel（忽略大小写）
// 如果找不到，返回 DefaultLogLevel 作为默认值
func StringToLevel(s string) LogLevel {
	if s != "" {
		for i := range LogLevelStrName {
			if strings.EqualFold(LogLevelStrName[i], s) {
				return LogLevel(i)
			}
		}
	}
	return DefaultLogLevel
}

var levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// 日志条目结构
type logEntry struct {
	level   LogLevel
	time    time.Time
	message string
	caller  string
	group   string // 群组字段
}

// 日志核心配置
type Config struct {
	Level       LogLevel
	TimeFormat  string
	CallerDepth int
	BaseDir     string // 全局根路径
}

// 日志输出器接口
type LogWriter interface {
	Write(entry *logEntry) error
	Close() error
}

// 日志记录器主体
type Logger struct {
	config       Config
	writers      []LogWriter
	entryCh      chan *logEntry
	wg           sync.WaitGroup
	closeCh      chan struct{}
	closed       bool
	mu           sync.Mutex
	defaultGroup string // 默认群组名称
}

// 全局日志实例
var (
	globalLogger     *Logger
	globalLoggerOnce = singleton.NewSingleton[Logger]()
)

// 获取默认全局日志实例
func Default() *Logger {
	l, err := globalLoggerOnce.Get(func() (*Logger, error) {
		globalLogger = NewLogger()
		return globalLogger, nil
	})
	if err != nil {
		panic("qwqlog.Global() 获取全局日志实例错误")
	}
	return l
}

// 配置全局日志实例
func ConfigureGlobal(opts ...Option) error {
	logger := Default()
	for _, opt := range opts {
		opt(logger)
	}
	return nil
}

// 默认配置
var defaultConfig = Config{
	Level:       INFO,
	TimeFormat:  "2006-01-02 15:04:05.000",
	CallerDepth: 3,
	BaseDir:     "", // 默认为空
}

// 创建新日志实例
func NewLogger(opts ...Option) *Logger {
	logger := &Logger{
		config:       defaultConfig,
		entryCh:      make(chan *logEntry, 4096),
		closeCh:      make(chan struct{}),
		defaultGroup: "default",
	}

	// 应用选项
	for _, opt := range opts {
		opt(logger)
	}

	logger.wg.Add(1)
	go logger.processEntries()

	return logger
}

// 配置选项模式
type Option func(*Logger)

// 设置日志级别
func WithLevel(levelStr string) Option {
	return func(l *Logger) {
		level := StringToLevel(levelStr)
		l.config.Level = level
	}
}

// 设置全局根路径
func WithBaseDir(dir string) Option {
	return func(l *Logger) {
		l.config.BaseDir = dir
	}
}

// 设置默认群组
func WithDefaultGroup(group string) Option {
	return func(l *Logger) {
		l.defaultGroup = group
	}
}

// 设置时间格式
func WithTimeFormat(format string) Option {
	return func(l *Logger) {
		l.config.TimeFormat = format
	}
}

// 设置调用者深度
func WithCallerDepth(depth int) Option {
	return func(l *Logger) {
		l.config.CallerDepth = depth
	}
}

// 核心处理协程
func (l *Logger) processEntries() {
	defer l.wg.Done()

	for {
		select {
		case entry := <-l.entryCh:
			l.writeToWriters(entry)
		case <-l.closeCh:
			// 关闭前处理剩余日志
			for len(l.entryCh) > 0 {
				entry := <-l.entryCh
				l.writeToWriters(entry)
			}
			return
		}
	}
}

// 写入所有输出器
func (l *Logger) writeToWriters(entry *logEntry) {
	for _, writer := range l.writers {
		if err := writer.Write(entry); err != nil {
			fmt.Fprintf(os.Stderr, "Log write error: %v\n", err)
		}
	}
}

// 添加日志输出器
func (l *Logger) AddWriter(writer LogWriter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = append(l.writers, writer)
}

// 添加控制台输出器
func (l *Logger) AddConsoleWriter() {
	l.AddWriter(NewConsoleWriter())
}

// 添加文件输出器（自动应用全局根路径）
func (l *Logger) AddFileWriter(filename string, maxSizeMB, maxBackups, maxAgeDays int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 应用全局根路径
	fullPath := filename
	if l.config.BaseDir != "" {
		fullPath = filepath.Join(l.config.BaseDir, filename)

		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
		}
	}

	writer := NewFileWriter(fullPath, maxSizeMB, maxBackups, maxAgeDays)
	l.writers = append(l.writers, writer)
}

// 添加JSON输出器（自动应用全局根路径）
func (l *Logger) AddJSONWriter(filename string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 应用全局根路径
	fullPath := filename
	if l.config.BaseDir != "" {
		fullPath = filepath.Join(l.config.BaseDir, filename)

		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
		}
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open JSON log file: %v\n", err)
		return
	}

	writer := NewJSONWriter(file)
	l.writers = append(l.writers, writer)
}

// 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Level = level
}

// 设置默认群组
func (l *Logger) SetDefaultGroup(group string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.defaultGroup = group
}

// 设置全局根路径
func (l *Logger) SetBaseDir(dir string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.BaseDir = dir
}

// 通用日志记录方法（带群组）
func (l *Logger) logWithGroup(group string, level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	cfg := l.config
	l.mu.Unlock()

	if level < cfg.Level {
		return
	}

	message := fmt.Sprintf(format, args...)
	_, file, line, ok := runtime.Caller(cfg.CallerDepth)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// 使用默认群组如果未指定
	if group == "" {
		group = l.defaultGroup
	}

	entry := &logEntry{
		level:   level,
		time:    time.Now(),
		message: message,
		caller:  caller,
		group:   group,
	}

	select {
	case l.entryCh <- entry:
	default:
		// 缓冲区满时直接输出到stderr防止阻塞
		fmt.Fprintf(os.Stderr, "Log queue full! [%s] %s\n", group, message)
	}
}

// 分级日志方法（使用默认群组）
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logWithGroup("", DEBUG, format, args...)
}
func (l *Logger) Info(format string, args ...interface{}) { l.logWithGroup("", INFO, format, args...) }
func (l *Logger) Warn(format string, args ...interface{}) { l.logWithGroup("", WARN, format, args...) }
func (l *Logger) Error(format string, args ...interface{}) {
	l.logWithGroup("", ERROR, format, args...)
}
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logWithGroup("", FATAL, format, args...)
	os.Exit(1)
}

// 带群组的分级日志方法
func (l *Logger) DebugGroup(group, format string, args ...interface{}) {
	l.logWithGroup(group, DEBUG, format, args...)
}
func (l *Logger) InfoGroup(group, format string, args ...interface{}) {
	l.logWithGroup(group, INFO, format, args...)
}
func (l *Logger) WarnGroup(group, format string, args ...interface{}) {
	l.logWithGroup(group, WARN, format, args...)
}
func (l *Logger) ErrorGroup(group, format string, args ...interface{}) {
	l.logWithGroup(group, ERROR, format, args...)
}
func (l *Logger) FatalGroup(group, format string, args ...interface{}) {
	l.logWithGroup(group, FATAL, format, args...)
	os.Exit(1)
}

// 安全关闭日志系统
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return
	}

	close(l.closeCh)
	l.wg.Wait()

	for _, writer := range l.writers {
		if err := writer.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Log writer close error: %v\n", err)
		}
	}
	l.closed = true
}

// ==================== 输出器实现 ====================

// 控制台输出器
type ConsoleWriter struct {
	out io.Writer
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{out: os.Stdout}
}

func (w *ConsoleWriter) Write(entry *logEntry) error {
	color := getColorByLevel(entry.level)
	reset := "\033[0m"

	// 添加群组信息
	logLine := fmt.Sprintf("%s [%s] [%s] %s %s\n",
		entry.time.Format(defaultConfig.TimeFormat),
		levelNames[entry.level],
		entry.group,
		entry.message,
		entry.caller)

	if color != "" {
		logLine = color + logLine + reset
	}

	_, err := fmt.Fprint(w.out, logLine)
	return err
}

func (w *ConsoleWriter) Close() error { return nil }

func getColorByLevel(level LogLevel) string {
	switch level {
	case DEBUG:
		return "\033[36m" // 青色
	case WARN:
		return "\033[33m" // 黄色
	case ERROR, FATAL:
		return "\033[31m" // 红色
	default:
		return ""
	}
}

// 文件输出器（带日志轮转）
type FileWriter struct {
	writer io.Writer
	lj     *lumberjack.Logger
}

func NewFileWriter(filename string, maxSizeMB, maxBackups, maxAgeDays int) *FileWriter {
	lj := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   true,
		LocalTime:  true,
	}

	return &FileWriter{
		writer: lj,
		lj:     lj,
	}
}

func (w *FileWriter) Write(entry *logEntry) error {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s [%s] [%s] %s %s\n",
		entry.time.Format(defaultConfig.TimeFormat),
		levelNames[entry.level],
		entry.group,
		entry.message,
		entry.caller)

	_, err := buf.WriteTo(w.writer)
	return err
}

func (w *FileWriter) Close() error {
	return w.lj.Close()
}

// JSON格式输出器
type JSONWriter struct {
	writer io.Writer
}

func NewJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{writer: w}
}

func (w *JSONWriter) Write(entry *logEntry) error {
	logLine := fmt.Sprintf(
		`{"time":"%s","level":"%s","group":"%s","message":"%s","caller":"%s"}`+"\n",
		entry.time.Format(time.RFC3339Nano),
		levelNames[entry.level],
		entry.group,
		escapeJSON(entry.message),
		entry.caller)

	_, err := w.writer.Write([]byte(logLine))
	return err
}

func (w *JSONWriter) Close() error {
	if c, ok := w.writer.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
