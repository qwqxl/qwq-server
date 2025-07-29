// // qwqlog.go
package qwqlog

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"gopkg.in/natefinch/lumberjack.v2"
//	"io"
//	"os"
//	"path/filepath"
//	"runtime"
//	"strconv"
//	"sync"
//	"sync/atomic"
//	"time"
//	"unsafe"
//)
//
//// 日志级别类型
//type LogLevel int
//
//const (
//	DEBUG LogLevel = iota
//	INFO
//	WARN
//	ERROR
//	FATAL
//)
//
//var (
//	levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
//
//	// 全局日志目录设置
//	globalLogDir     string
//	globalLogDirOnce sync.Once
//	globalLogDirMu   sync.RWMutex
//)
//
//// 设置全局日志目录
//func SetGlobalLogDir(dir string) error {
//	globalLogDirMu.Lock()
//	defer globalLogDirMu.Unlock()
//
//	// 确保目录存在
//	if err := os.MkdirAll(dir, 0755); err != nil {
//		return fmt.Errorf("failed to create log directory: %w", err)
//	}
//
//	globalLogDir = dir
//	return nil
//}
//
//// 获取日志文件的完整路径
//func getLogPath(filename string) string {
//	globalLogDirMu.RLock()
//	defer globalLogDirMu.RUnlock()
//
//	if globalLogDir == "" {
//		// 初始化默认日志目录为当前目录下的logs
//		globalLogDirOnce.Do(func() {
//			cwd, err := os.Getwd()
//			if err == nil {
//				globalLogDir = filepath.Join(cwd, "logs")
//				_ = os.MkdirAll(globalLogDir, 0755)
//			}
//		})
//	}
//
//	// 如果已经是绝对路径，直接返回
//	if filepath.IsAbs(filename) {
//		return filename
//	}
//
//	// 否则组合全局日志目录
//	return filepath.Join(globalLogDir, filename)
//}
//
//// 日志条目结构
//type logEntry struct {
//	level   LogLevel
//	time    time.Time
//	message string
//	caller  string
//}
//
//// 日志核心配置
//type Config struct {
//	Level          LogLevel
//	TimeFormat     string
//	CallerDepth    int
//	EnableCaller   bool // 是否启用调用者信息
//	AsyncQueueSize int  // 异步队列大小
//}
//
//// 日志输出器接口
//type LogWriter interface {
//	Write(entry *logEntry) error
//	Close() error
//	NeedsCaller() bool // 是否需要调用者信息
//}
//
//// 日志记录器主体
//type Logger struct {
//	config  atomic.Value // 使用atomic.Value实现无锁读取配置
//	writers []LogWriter
//	entryCh chan *logEntry
//	wg      sync.WaitGroup
//	closeCh chan struct{}
//	closed  int32 // 使用atomic操作保证线程安全
//	pool    sync.Pool
//}
//
//// 默认配置
//var defaultConfig = Config{
//	Level:          INFO,
//	TimeFormat:     "2006-01-02 15:04:05.000",
//	CallerDepth:    3,
//	EnableCaller:   true,
//	AsyncQueueSize: 8192, // 更大的默认缓冲区
//}
//
//// 创建新日志实例
//func NewLogger() *Logger {
//	logger := &Logger{
//		entryCh: make(chan *logEntry, defaultConfig.AsyncQueueSize),
//		closeCh: make(chan struct{}),
//	}
//
//	logger.config.Store(defaultConfig)
//	logger.pool = sync.Pool{
//		New: func() interface{} {
//			return &logEntry{}
//		},
//	}
//
//	logger.wg.Add(1)
//	go logger.processEntries()
//
//	return logger
//}
//
//// 核心处理协程
//func (l *Logger) processEntries() {
//	defer l.wg.Done()
//
//	for {
//		select {
//		case entry := <-l.entryCh:
//			l.writeToWriters(entry)
//			l.pool.Put(entry) // 归还对象到池
//		case <-l.closeCh:
//			// 关闭前处理剩余日志
//			for {
//				select {
//				case entry := <-l.entryCh:
//					l.writeToWriters(entry)
//					l.pool.Put(entry)
//				default:
//					return
//				}
//			}
//		}
//	}
//}
//
//// 写入所有输出器
//func (l *Logger) writeToWriters(entry *logEntry) {
//	cfg := l.config.Load().(Config)
//	for _, writer := range l.writers {
//		if entry.level < cfg.Level {
//			continue
//		}
//		if err := writer.Write(entry); err != nil {
//			fmt.Fprintf(os.Stderr, "Log write error: %v\n", err)
//		}
//	}
//}
//
//// 添加日志输出器
//func (l *Logger) AddWriter(writer LogWriter) {
//	l.writers = append(l.writers, writer)
//	l.updateCallerSetting()
//}
//
//// 更新调用者设置
//func (l *Logger) updateCallerSetting() {
//	cfg := l.config.Load().(Config)
//	enableCaller := false
//	for _, w := range l.writers {
//		if w.NeedsCaller() {
//			enableCaller = true
//			break
//		}
//	}
//
//	if enableCaller != cfg.EnableCaller {
//		newCfg := cfg
//		newCfg.EnableCaller = enableCaller
//		l.config.Store(newCfg)
//	}
//}
//
//// 设置日志级别
//func (l *Logger) SetLevel(level LogLevel) {
//	cfg := l.config.Load().(Config)
//	newCfg := cfg
//	newCfg.Level = level
//	l.config.Store(newCfg)
//}
//
//// 通用日志记录方法
//func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
//	cfg := l.config.Load().(Config)
//	if level < cfg.Level {
//		return
//	}
//
//	// 从池中获取entry对象
//	entry := l.pool.Get().(*logEntry)
//	entry.level = level
//	entry.time = time.Now()
//
//	// 优化：避免不必要的格式化
//	if len(args) == 0 {
//		entry.message = format
//	} else {
//		entry.message = fmt.Sprintf(format, args...)
//	}
//
//	// 只有需要时才获取调用者信息
//	if cfg.EnableCaller {
//		_, file, line, ok := runtime.Caller(cfg.CallerDepth)
//		if ok {
//			// 使用更高效的方式构建caller字符串
//			entry.caller = filepath.Base(file) + ":" + strconv.Itoa(line)
//		} else {
//			entry.caller = "???:0"
//		}
//	} else {
//		entry.caller = ""
//	}
//
//	select {
//	case l.entryCh <- entry:
//	default:
//		// 缓冲区满时优化处理
//		if atomic.LoadInt32(&l.closed) == 0 {
//			// 只输出消息核心部分避免递归
//			fmt.Fprintf(os.Stderr, "Log queue full! Message: %s\n", entry.message)
//		}
//		// 归还对象到池
//		l.pool.Put(entry)
//	}
//}
//
//// 分级日志方法
//func (l *Logger) Debug(format string, args ...interface{}) { l.log(DEBUG, format, args...) }
//func (l *Logger) Info(format string, args ...interface{})  { l.log(INFO, format, args...) }
//func (l *Logger) Warn(format string, args ...interface{})  { l.log(WARN, format, args...) }
//func (l *Logger) Error(format string, args ...interface{}) { l.log(ERROR, format, args...) }
//func (l *Logger) Fatal(format string, args ...interface{}) {
//	l.log(FATAL, format, args...)
//	l.Close()
//	os.Exit(1)
//}
//
//// 安全关闭日志系统
//func (l *Logger) Close() {
//	if !atomic.CompareAndSwapInt32(&l.closed, 0, 1) {
//		return
//	}
//
//	close(l.closeCh)
//	l.wg.Wait()
//
//	for _, writer := range l.writers {
//		if err := writer.Close(); err != nil {
//			fmt.Fprintf(os.Stderr, "Log writer close error: %v\n", err)
//		}
//	}
//}
//
//// ==================== 输出器实现 ====================
//
//// 控制台输出器
//type ConsoleWriter struct {
//	out    io.Writer
//	mu     sync.Mutex
//	isTerm bool // 是否是终端
//}
//
//func NewConsoleWriter() *ConsoleWriter {
//	w := &ConsoleWriter{out: os.Stdout}
//	// 检测是否是终端
//	if f, ok := w.out.(*os.File); ok {
//		fi, _ := f.Stat()
//		w.isTerm = (fi.Mode() & os.ModeCharDevice) != 0
//	}
//	return w
//}
//
//// 获取日志级别对应的颜色代码
//func getColorByLevel(level LogLevel) string {
//	switch level {
//	case DEBUG:
//		return "\033[36m" // 青色
//	case WARN:
//		return "\033[33m" // 黄色
//	case ERROR, FATAL:
//		return "\033[31m" // 红色
//	default:
//		return ""
//	}
//}
//
//func (w *ConsoleWriter) Write(entry *logEntry) error {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//
//	cfg := defaultConfig // 使用默认配置
//	var logLine string
//
//	if w.isTerm {
//		color := getColorByLevel(entry.level)
//		reset := "\033[0m"
//		logLine = fmt.Sprintf("%s%s [%s] %s %s%s\n",
//			color,
//			entry.time.Format(cfg.TimeFormat),
//			levelNames[entry.level],
//			entry.message,
//			entry.caller,
//			reset)
//	} else {
//		logLine = fmt.Sprintf("%s [%s] %s %s\n",
//			entry.time.Format(cfg.TimeFormat),
//			levelNames[entry.level],
//			entry.message,
//			entry.caller)
//	}
//
//	_, err := w.out.Write(unsafeBytes(logLine))
//	return err
//}
//
//func (w *ConsoleWriter) Close() error      { return nil }
//func (w *ConsoleWriter) NeedsCaller() bool { return true }
//
//// 文件输出器（带日志轮转）
//type FileWriter struct {
//	mu     sync.Mutex
//	lj     *lumberjack.Logger
//	buffer *bytes.Buffer
//}
//
//// 创建新的文件输出器，支持全局日志目录
//func NewFileWriter(filename string, maxSizeMB, maxBackups, maxAgeDays int) *FileWriter {
//	// 获取完整日志路径
//	fullPath := getLogPath(filename)
//
//	// 确保目录存在
//	dir := filepath.Dir(fullPath)
//	_ = os.MkdirAll(dir, 0755)
//
//	return &FileWriter{
//		lj: &lumberjack.Logger{
//			Filename:   fullPath,
//			MaxSize:    maxSizeMB,
//			MaxBackups: maxBackups,
//			MaxAge:     maxAgeDays,
//			Compress:   true,
//			LocalTime:  true,
//		},
//		buffer: bytes.NewBuffer(make([]byte, 0, 256)),
//	}
//}
//
//func (w *FileWriter) Write(entry *logEntry) error {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//
//	cfg := defaultConfig
//	w.buffer.Reset()
//
//	fmt.Fprintf(w.buffer, "%s [%s] %s %s\n",
//		entry.time.Format(cfg.TimeFormat),
//		levelNames[entry.level],
//		entry.message,
//		entry.caller)
//
//	_, err := w.lj.Write(w.buffer.Bytes())
//	return err
//}
//
//func (w *FileWriter) Close() error {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//	return w.lj.Close()
//}
//
//func (w *FileWriter) NeedsCaller() bool { return true }
//
//// JSON格式输出器
//type JSONWriter struct {
//	mu     sync.Mutex
//	writer io.Writer
//	buffer *bytes.Buffer
//}
//
//func NewJSONWriter(w io.Writer) *JSONWriter {
//	return &JSONWriter{
//		writer: w,
//		buffer: bytes.NewBuffer(make([]byte, 0, 256)),
//	}
//}
//
//// 创建JSON文件输出器，支持全局日志目录
//func NewJSONFileWriter(filename string) *JSONWriter {
//	// 获取完整日志路径
//	fullPath := getLogPath(filename)
//
//	// 确保目录存在
//	dir := filepath.Dir(fullPath)
//	_ = os.MkdirAll(dir, 0755)
//
//	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		// 失败时回退到标准错误
//		fmt.Fprintf(os.Stderr, "Failed to open JSON log file: %v\n", err)
//		file = os.Stderr
//	}
//
//	return NewJSONWriter(file)
//}
//
//func (w *JSONWriter) Write(entry *logEntry) error {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//
//	w.buffer.Reset()
//	w.buffer.WriteString(`{"time":"`)
//	w.buffer.WriteString(entry.time.Format(time.RFC3339Nano))
//	w.buffer.WriteString(`","level":"`)
//	w.buffer.WriteString(levelNames[entry.level])
//	w.buffer.WriteString(`","message":`)
//	jsonMessage, _ := json.Marshal(entry.message)
//	w.buffer.Write(jsonMessage)
//	w.buffer.WriteString(`,"caller":"`)
//	w.buffer.WriteString(entry.caller)
//	w.buffer.WriteString(`"}`)
//	w.buffer.WriteByte('\n')
//
//	_, err := w.writer.Write(w.buffer.Bytes())
//	return err
//}
//
//func (w *JSONWriter) Close() error {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//	if c, ok := w.writer.(io.Closer); ok {
//		return c.Close()
//	}
//	return nil
//}
//
//func (w *JSONWriter) NeedsCaller() bool { return true }
//
//// 高性能字符串转换
//func unsafeBytes(s string) []byte {
//	return *(*[]byte)(unsafe.Pointer(
//		&struct {
//			string
//			Cap int
//		}{s, len(s)},
//	))
//}
