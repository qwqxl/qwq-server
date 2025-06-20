package main

import (
	"fmt"
	"qwqserver/internal/app"
	"qwqserver/internal/config"
)

type qwqLog struct {
}

func (q qwqLog) Log(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	fmt.Println(s)
}

func (q qwqLog) Info(msg string, args ...any) {
	q.Log(msg, args)
}

func (q qwqLog) Error(msg string, args ...any) {
	q.Log(msg, args)
}

func (q qwqLog) Warn(msg string, args ...any) {
	q.Log(msg)
}

func (q qwqLog) Debug(msg string, args ...any) {
	q.Log(msg)
}

func (q qwqLog) Fatal(msg string, args ...any) {
	q.Log(msg, args)
}

func (q qwqLog) Panic(msg string, args ...any) {
	q.Log(msg, args)
}

func main() {
	l := qwqLog{}
	cfg := config.New()
	a := app.New(&app.Config{
		L:      l,
		Config: cfg,
	})

	defer a.Close()

	// 启动服务
	a.Listen()

}
