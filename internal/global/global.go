package global

import (
	"qwqserver/pkg/qwqlog"
	"qwqserver/pkg/util/singleton"
)

type HTTPResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Global struct {
	Logger *qwqlog.Logger
}

var (
	App = singleton.NewSingleton[Global]()
)

func Apps() (*Global, error) {
	return App.Get(func() (*Global, error) {
		apps := &Global{}

		return apps, nil
	})
}
