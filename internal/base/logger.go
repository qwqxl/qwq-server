package base

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)
}
