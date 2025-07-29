package httpcore

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"os/signal"
	"qwqserver/pkg/httpcore/httpconnpool"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"qwqserver/pkg/logger"
	"qwqserver/pkg/util/singleton"
)

// HandlerFunc 自定义 handler 类型
// 完全屏蔽 gin.HandlerFunc

type HandlerFunc func(Context)

// Context 自定义请求上下文接口（可扩展）
type Context interface {
	JSON(code int, obj any)
	Bind(obj any) error
	Param(key string) string
	Query(key string) string
	Set(key string, value any)
	Get(key string) (any, bool)
	GetHeader(key string) string
	Header(key, value string)
	Status(code int)
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, obj any)
	Logger() logger.Logger
	Request() *http.Request
	Response() http.ResponseWriter
	Abort()
	Next()
	ClientIP() string
	ShouldBindJSON(v any) error
	NotFound()       // 手动触发 404 响应
	IsHandled() bool // 检查请求是否已被处理
}

// ResponseWriter 预留接口（若需自定义响应输出）
type ResponseWriter interface {
	Write([]byte) (int, error)
}

// Engine 是对 gin.Engine 的封装
// 封装所有对外使用接口，隐藏 gin 实现
type Engine struct {
	engine          *gin.Engine
	httpServer      *http.Server
	logger          logger.Logger
	connPool        *httpconnpool.ServerConnPool
	staticDirs      map[string]string // 静态文件映射：URL路径前缀 -> 本地目录
	notFoundHandler HandlerFunc       // 自定义 404 处理函数
}

// contextWrapper 封装 gin.Context，实现自定义 Context 接口
type contextWrapper struct {
	ctx       *gin.Context
	logIns    logger.Logger
	requestID string
}

func (c *contextWrapper) JSON(code int, obj any) {
	c.ctx.JSON(code, obj)
}
func (c *contextWrapper) Bind(obj any) error {
	return c.ctx.Bind(obj)
}
func (c *contextWrapper) Param(key string) string {
	return c.ctx.Param(key)
}
func (c *contextWrapper) Query(key string) string {
	return c.ctx.Query(key)
}
func (c *contextWrapper) Set(key string, value any) {
	c.ctx.Set(key, value)
}
func (c *contextWrapper) Get(key string) (any, bool) {
	return c.ctx.Get(key)
}
func (c *contextWrapper) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}
func (c *contextWrapper) Header(key, value string) {
	c.ctx.Header(key, value)
}
func (c *contextWrapper) Status(code int) {
	c.ctx.Status(code)
}
func (c *contextWrapper) AbortWithStatus(code int) {
	c.ctx.AbortWithStatus(code)
}
func (c *contextWrapper) AbortWithStatusJSON(code int, obj any) {
	c.ctx.AbortWithStatusJSON(code, obj)
}
func (c *contextWrapper) Logger() logger.Logger {
	return c.logIns
}
func (c *contextWrapper) Request() *http.Request {
	return c.ctx.Request
}
func (c *contextWrapper) Response() http.ResponseWriter {
	return c.ctx.Writer
}
func (c *contextWrapper) Abort() {
	c.ctx.Abort()
}
func (c *contextWrapper) Next() {
	c.ctx.Next()
}
func (c *contextWrapper) ClientIP() string {
	return c.ctx.ClientIP()
}
func (c *contextWrapper) ShouldBindJSON(v any) error {
	return c.ctx.ShouldBindJSON(v)
}

func (c *contextWrapper) RequestID() string {
	if c.requestID == "" {
		c.requestID = uuid.New().String()
		c.ctx.Header("X-Request-ID", c.requestID)
	}
	return c.requestID
}

func (c *contextWrapper) NotFound() {
	c.ctx.AbortWithStatus(http.StatusNotFound)
}

func (c *contextWrapper) IsHandled() bool {
	return c.ctx.IsAborted() || c.ctx.Writer.Status() != http.StatusOK
}

// 在中间件中设置
func requestIDMiddleware() HandlerFunc {
	return func(c Context) {
		if wc, ok := c.(*contextWrapper); ok {
			wc.RequestID() // 初始化requestID
		}
		c.Next()
	}
}

// RouterGroup routerGroupWrapper 实现 RouterGroup 接口
type RouterGroup interface {
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	Group(relativePath string) RouterGroup
	Use(middleware ...HandlerFunc)
}

type routerGroupWrapper struct {
	group *gin.RouterGroup
	log   logger.Logger
}

func (g *routerGroupWrapper) GET(path string, handler HandlerFunc) {
	g.group.GET(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: g.log})
	})
}

func (g *routerGroupWrapper) POST(path string, handler HandlerFunc) {
	g.group.POST(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: g.log})
	})
}

func (g *routerGroupWrapper) PUT(path string, handler HandlerFunc) {
	g.group.PUT(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: g.log})
	})
}

func (g *routerGroupWrapper) DELETE(path string, handler HandlerFunc) {
	g.group.DELETE(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: g.log})
	})
}

func (g *routerGroupWrapper) Group(relativePath string) RouterGroup {
	return &routerGroupWrapper{
		group: g.group.Group(relativePath),
		log:   g.log,
	}
}
func (g *routerGroupWrapper) UseV1(middleware ...HandlerFunc) {
	handlers := make([]gin.HandlerFunc, 0, len(middleware))
	for _, h := range middleware {
		handlers = append(handlers, func(c *gin.Context) {
			h(&contextWrapper{ctx: c, logIns: g.log})
		})
	}
	g.group.Use(handlers...)
}

func (g *routerGroupWrapper) Use(middleware ...HandlerFunc) {
	handlers := make([]gin.HandlerFunc, len(middleware))
	for i, h := range middleware {
		// 避免闭包捕获问题
		handler := h
		handlers[i] = func(c *gin.Context) {
			handler(&contextWrapper{ctx: c, logIns: g.log})
		}
	}
	g.group.Use(handlers...)
}

// 单例引擎
var globalEngine = singleton.NewSingleton[Engine]()

func DefaultWithLogger(log logger.Logger) *Engine {
	v, _ := globalEngine.Get(func() (*Engine, error) {
		ginEngine := gin.New()
		eng := &Engine{
			engine:          ginEngine,
			httpServer:      &http.Server{Handler: ginEngine},
			logger:          log,
			staticDirs:      make(map[string]string),
			notFoundHandler: defaultNotFoundHandler,
		}

		// 添加默认中间件
		eng.Use(eng.loggingMiddleware())
		eng.engine.Use(gin.Recovery()) // Gin 的恢复中间件

		// 添加静态文件中间件（在日志中间件之后）
		eng.Use(eng.staticMiddleware())

		// 设置 Gin 的 NoRoute 处理
		//ginEngine.NoRoute(func(c *gin.Context) {
		//	defaultNotFoundHandler(&contextWrapper{ctx: c, logIns: log})
		//})

		// 添加自定义 404 中间件（作为最后一个中间件）
		eng.Use(eng.notFoundMiddleware())

		return eng, nil
	})
	return v
}

// loggingMiddleware 日志中间件封装
func (e *Engine) loggingMiddlewareV1() HandlerFunc {
	return func(c Context) {
		start := time.Now()
		sc := c.(*contextWrapper)
		sc.ctx.Next()
		duration := time.Since(start)
		ip := GetClientIP(c.Request())
		e.logger.Info("%s %s %s %d [%s]", ip, sc.ctx.Request.Method, sc.ctx.Request.URL.Path, sc.ctx.Writer.Status(), duration)
	}
}

// 增强日志中间件
func (e *Engine) loggingMiddleware() HandlerFunc {
	return func(c Context) {
		start := time.Now()
		req := c.Request()
		clientIP := GetClientIP(req)

		// 记录请求信息
		e.logger.Debug("Request started: %s %s %s", clientIP, req.Method, req.URL.Path)

		defer func() {
			latency := time.Since(start)
			status := c.Response().(gin.ResponseWriter).Status()
			clientIP = GetClientIP(req)

			e.logger.Info("[%s] %s %s %d %v",
				c.(*contextWrapper).RequestID(),
				clientIP,
				req.URL.Path,
				status,
				latency,
			)
		}()

		c.Next()
	}
}

// 注册通用路由

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.engine.GET(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: e.logger})
	})
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.engine.POST(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: e.logger})
	})
}

func (e *Engine) PUT(path string, handler HandlerFunc) {
	e.engine.PUT(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: e.logger})
	})
}

func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.engine.DELETE(path, func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: e.logger})
	})
}

func (e *Engine) Use(handlers ...HandlerFunc) {
	for _, h := range handlers {
		e.engine.Use(func(c *gin.Context) {
			h(&contextWrapper{ctx: c, logIns: e.logger})
		})
	}
}

// ModuleRegistrar 模块化路由注册
type ModuleRegistrar interface {
	RegisterRoutes(e *Engine)
}

func (e *Engine) RegisterModules(modules ...ModuleRegistrar) {
	for _, m := range modules {
		m.RegisterRoutes(e)
	}
}

// 设置连接池
func (e *Engine) SetConnPool(pool *httpconnpool.ServerConnPool) {
	e.connPool = pool
	e.httpServer.ConnState = pool.ConnStateCallback // 设置状态回调

	// 启动连接池清理协程
	go pool.StartCleanup()
}

func (e *Engine) GetConnPool() (*httpconnpool.ServerConnPool, error) {
	if e.connPool == nil {
		return nil, errors.New("connPool is nil")
	}
	return e.connPool, nil
}

// RunWithGracefulShutdown 优雅关停
func (e *Engine) RunWithGracefulShutdown(addr string, shutdownTimeout time.Duration) error {
	e.httpServer.Addr = addr
	e.logger.Info("Starting server at %s", addr)

	// 启动连接池清理（如果存在）
	if e.connPool != nil {
		go e.connPool.StartCleanup()
	}

	// 启动服务器
	serverErr := make(chan error, 1)
	go func() {
		if err := e.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	// 等待服务器启动或出错
	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("server failed to start: %w", err)
		}
	case <-time.After(100 * time.Millisecond):
		// 短暂等待确保服务器开始监听
	}

	return e.waitForShutdown(shutdownTimeout)
}

// RunTLSWithGracefulShutdown 优雅关停的HTTPS启动方法
func (e *Engine) RunTLSWithGracefulShutdown(
	addr, certFile, keyFile string,
	shutdownTimeout time.Duration,
) error {
	e.httpServer.Addr = addr
	e.logger.Info("Starting HTTPS server at %s", addr)

	// 启动连接池清理（如果存在）
	if e.connPool != nil {
		go e.connPool.StartCleanup()
	}

	// 启动服务器
	serverErr := make(chan error, 1)
	go func() {
		if err := e.httpServer.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	// 等待服务器启动或出错
	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("HTTPS server failed to start: %w", err)
		}
	case <-time.After(100 * time.Millisecond):
		// 短暂等待确保服务器开始监听
	}

	return e.waitForShutdown(shutdownTimeout)
}

// waitForShutdown 统一等待关闭信号
func (e *Engine) waitForShutdown(shutdownTimeout time.Duration) error {
	// 创建带缓冲的信号通道
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待关闭信号
	sig := <-quit
	e.logger.Info("Received signal: %v, shutting down server...", sig)

	// 创建带超时的关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// 先关闭连接池（如果存在）
	if e.connPool != nil {
		e.logger.Debug("Gracefully shutting down connection pool...")
		e.connPool.GracefulShutdown()
	}

	// 然后关闭HTTP服务器
	e.logger.Debug("Shutting down HTTP server...")
	if err := e.httpServer.Shutdown(ctx); err != nil {
		e.logger.Error("HTTP server shutdown error: %v", err)
		return err
	}

	e.logger.Info("Server gracefully stopped")
	return nil
}

// Group 路由组封装
func (e *Engine) Group(prefix string) RouterGroup {
	return &routerGroupWrapper{
		group: e.engine.Group(prefix),
		log:   e.logger,
	}
}

// RunTLS TLS 启动支持
func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	e.httpServer.Addr = addr
	e.logger.Info("Starting HTTPS server at %s", addr)
	return e.httpServer.ListenAndServeTLS(certFile, keyFile)
}

// Get Field

// NewHTTPServer get http.Server
func (e *Engine) NewHTTPServer() *http.Server {
	if e.httpServer == nil {
		e.httpServer = &http.Server{}
	}
	return e.httpServer
}
