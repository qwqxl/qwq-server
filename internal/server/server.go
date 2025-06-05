package server

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Server struct {
	// 单例实现
	router *gin.Engine
	//  线程安全实现
	Once *sync.Once
}

var serverInstance = &Server{
	router: nil,
	Once:   &sync.Once{},
}

func New() *gin.Engine {
	serverInstance.Once.Do(func() {
		serverInstance.router = gin.New()
	})
	return serverInstance.router
}

func Run(addr string) error {
	return New().Run(addr)
}
