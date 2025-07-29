package logger

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	DebugGroup(group, format string, args ...interface{})
	InfoGroup(group, format string, args ...interface{})
	WarnGroup(group, format string, args ...interface{})
	ErrorGroup(group, format string, args ...interface{})
	FatalGroup(group, format string, args ...interface{})
	Close()
}

type Loggers struct {
	Log Logger
}

func (l Loggers) Debug(format string, args ...interface{}) {
	l.Debug(format, args...)
}

func (l Loggers) Info(format string, args ...interface{}) {
	l.Info(format, args...)
}

func (l Loggers) Warn(format string, args ...interface{}) {
	l.Warn(format, args...)
}

func (l Loggers) Error(format string, args ...interface{}) {
	l.Error(format, args...)
}

func (l Loggers) Fatal(format string, args ...interface{}) {
	l.Fatal(format, args...)
}

func (l Loggers) DebugGroup(group, format string, args ...interface{}) {
	l.DebugGroup(group, format, args...)
}

func (l Loggers) InfoGroup(group, format string, args ...interface{}) {
	l.InfoGroup(group, format, args...)
}

func (l Loggers) WarnGroup(group, format string, args ...interface{}) {
	l.WarnGroup(group, format, args...)
}

func (l Loggers) ErrorGroup(group, format string, args ...interface{}) {
	l.ErrorGroup(group, format, args...)
}

func (l Loggers) FatalGroup(group, format string, args ...interface{}) {
	l.FatalGroup(group, format, args...)
}

func (l Loggers) Close() {
	l.Close()
}
