package httpconnpool

import (
	"net"
	"sync"
	"time"
)

// 连接状态
type ConnState int

const (
	StateNew ConnState = iota
	StateActive
	StateIdle
	StateClosed
)

// 连接池结构
type ServerConnPool struct {
	mu          sync.Mutex
	idleConns   map[net.Conn]time.Time // 空闲连接及进入空闲时间
	activeConns map[net.Conn]struct{}  // 活跃连接

	maxIdle     int           // 最大空闲连接数
	maxOpen     int           // 最大打开连接数
	idleTimeout time.Duration // 空闲超时时间
}

func NewServerConnPool(maxIdle, maxOpen int, idleTimeout time.Duration) *ServerConnPool {
	return &ServerConnPool{
		idleConns:   make(map[net.Conn]time.Time),
		activeConns: make(map[net.Conn]struct{}),
		maxIdle:     maxIdle,
		maxOpen:     maxOpen,
		idleTimeout: idleTimeout,
	}
}
