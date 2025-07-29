package httpconnpool

import (
	"net"
	"net/http"
	"time"
)

// 连接状态变化回调
func (p *ServerConnPool) ConnStateCallback(conn net.Conn, state http.ConnState) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch state {
	case http.StateNew:
		p.activeConns[conn] = struct{}{}

	case http.StateActive:
		// 从空闲池移除（如果存在）
		delete(p.idleConns, conn)

	case http.StateIdle:
		// 检查是否超过最大空闲连接数
		if len(p.idleConns) < p.maxIdle {
			p.idleConns[conn] = time.Now()
		} else {
			conn.Close()
			delete(p.activeConns, conn)
		}

	case http.StateHijacked, http.StateClosed:
		delete(p.idleConns, conn)
		delete(p.activeConns, conn)
	}
}

// 定期清理过期连接
func (p *ServerConnPool) StartCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		now := time.Now()

		for conn, idleSince := range p.idleConns {
			if now.Sub(idleSince) > p.idleTimeout {
				conn.Close()
				delete(p.idleConns, conn)
				delete(p.activeConns, conn)
			}
		}

		p.mu.Unlock()
	}
}

// 优雅关闭
func (p *ServerConnPool) GracefulShutdown() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for conn := range p.activeConns {
		conn.Close()
	}

	for conn := range p.idleConns {
		conn.Close()
	}

	p.idleConns = make(map[net.Conn]time.Time)
	p.activeConns = make(map[net.Conn]struct{})
}
