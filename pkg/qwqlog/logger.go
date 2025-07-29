// // qwqlog.go
package qwqlog

//
//import (
//	"bytes"
//	"fmt"
//	"gopkg.in/natefinch/lumberjack.v2"
//	"io"
//	"os"
//	"path/filepath"
//	"runtime"
//	"strings"
//	"sync"
//	"time"
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
//var levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
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
//	Level       LogLevel
//	TimeFormat  string
//	CallerDepth int
//}
//
//// 日志输出器接口
//type LogWriter interface {
//	Write(entry *logEntry) error
//	Close() error
//}
//
//// 日志记录器主体
//type Logger struct {
//	config  Config
//	writers []LogWriter
//	entryCh chan *logEntry
//	wg      sync.WaitGroup
//	closeCh chan struct{}
//	closed  bool
//	mu      sync.Mutex
//}
//
//// 默认配置
//var defaultConfig = Config{
//	Level:       INFO,
//	TimeFormat:  "2006-01-02 15:04:05.000",
//	CallerDepth: 3,
//}
//
//// 创建新日志实例
//func NewLogger() *Logger {
//	logger := &Logger{
//		config:  defaultConfig,
//		entryCh: make(chan *logEntry, 4096), // 缓冲队列提升并发性能
//		closeCh: make(chan struct{}),
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
//		case <-l.closeCh:
//			// 关闭前处理剩余日志
//			for len(l.entryCh) > 0 {
//				entry := <-l.entryCh
//				l.writeToWriters(entry)
//			}
//			return
//		}
//	}
//}
//
//// 写入所有输出器
//func (l *Logger) writeToWriters(entry *logEntry) {
//	for _, writer := range l.writers {
//		if err := writer.Write(entry); err != nil {
//			fmt.Fprintf(os.Stderr, "Log write error: %v\n", err)
//		}
//	}
//}
//
//// 添加日志输出器
//func (l *Logger) AddWriter(writer LogWriter) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.writers = append(l.writers, writer)
//}
//
//// 设置日志级别
//func (l *Logger) SetLevel(level LogLevel) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.config.Level = level
//}
//
//// 通用日志记录方法
//func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
//	if level < l.config.Level {
//		return
//	}
//
//	message := fmt.Sprintf(format, args...)
//	_, file, line, ok := runtime.Caller(l.config.CallerDepth)
//	caller := ""
//	if ok {
//		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
//	}
//
//	entry := &logEntry{
//		level:   level,
//		time:    time.Now(),
//		message: message,
//		caller:  caller,
//	}
//
//	select {
//	case l.entryCh <- entry:
//	default:
//		// 缓冲区满时直接输出到stderr防止阻塞
//		fmt.Fprintf(os.Stderr, "Log queue full! %s\n", message)
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
//	os.Exit(1)
//}
//
//// 安全关闭日志系统
//func (l *Logger) Close() {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//
//	if l.closed {
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
//	l.closed = true
//}
//
//// ==================== 输出器实现 ====================
//
//// 控制台输出器
//type ConsoleWriter struct {
//	out io.Writer
//}
//
//func NewConsoleWriter() *ConsoleWriter {
//	return &ConsoleWriter{out: os.Stdout}
//}
//
//func (w *ConsoleWriter) Write(entry *logEntry) error {
//	color := getColorByLevel(entry.level)
//	reset := "\033[0m"
//
//	logLine := fmt.Sprintf("%s [%s] %s %s\n",
//		entry.time.Format(defaultConfig.TimeFormat),
//		levelNames[entry.level],
//		entry.message,
//		entry.caller)
//
//	if color != "" {
//		logLine = color + logLine + reset
//	}
//
//	_, err := fmt.Fprint(w.out, logLine)
//	return err
//}
//
//func (w *ConsoleWriter) Close() error { return nil }
//
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
//// 文件输出器（带日志轮转）
//type FileWriter struct {
//	writer io.Writer
//	lj     *lumberjack.Logger
//}
//
//func NewFileWriter(filename string, maxSizeMB, maxBackups int, maxAgeDays int) *FileWriter {
//	lj := &lumberjack.Logger{
//		Filename:   filename,
//		MaxSize:    maxSizeMB,  // 日志文件最大大小(MB)
//		MaxBackups: maxBackups, // 保留旧文件最大个数
//		MaxAge:     maxAgeDays, // 保留旧文件最大天数
//		Compress:   true,       // 压缩旧文件
//		LocalTime:  true,       // 使用本地时间
//	}
//
//	return &FileWriter{
//		writer: lj,
//		lj:     lj,
//	}
//}
//
//func (w *FileWriter) Write(entry *logEntry) error {
//	buf := bytes.NewBuffer(nil)
//	_, err := fmt.Fprintf(buf, "%s [%s] %s %s\n",
//		entry.time.Format(defaultConfig.TimeFormat),
//		levelNames[entry.level],
//		entry.message,
//		entry.caller)
//	if err != nil {
//		return err
//	}
//
//	_, err = buf.WriteTo(w.writer)
//	return err
//}
//
//func (w *FileWriter) Close() error {
//	return w.lj.Close()
//}
//
//// JSON格式输出器（适合ELK采集）
//type JSONWriter struct {
//	writer io.Writer
//}
//
//func NewJSONWriter(w io.Writer) *JSONWriter {
//	return &JSONWriter{writer: w}
//}
//
//func (w *JSONWriter) Write(entry *logEntry) error {
//	logLine := fmt.Sprintf(
//		`{"time":"%s","level":"%s","message":"%s","caller":"%s"}`+"\n",
//		entry.time.Format(time.RFC3339Nano),
//		levelNames[entry.level],
//		escapeJSON(entry.message),
//		entry.caller)
//
//	_, err := w.writer.Write([]byte(logLine))
//	return err
//}
//
//func (w *JSONWriter) Close() error {
//	if c, ok := w.writer.(io.Closer); ok {
//		return c.Close()
//	}
//	return nil
//}
//
//func escapeJSON(s string) string {
//	return strings.ReplaceAll(s, `"`, `\"`)
//}
