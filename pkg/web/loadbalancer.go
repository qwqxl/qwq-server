package web

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type backendStatus struct {
	url   *url.URL
	alive bool
	mu    sync.RWMutex
}

type loadBalancer struct {
	policy    string
	statuses  []*backendStatus
	index     uint32
	transport *http.Transport
}

func (lb *loadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 尝试找到可用后端
	for i := 0; i < len(lb.statuses)*2; i++ {
		backend := lb.nextBackend()

		backend.mu.RLock()
		alive := backend.alive
		backend.mu.RUnlock()

		if alive {
			proxy := httputil.NewSingleHostReverseProxy(backend.url)
			proxy.Transport = lb.transport
			proxy.ErrorHandler = lb.errorHandler(backend)
			proxy.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, "服务不可用", http.StatusServiceUnavailable)
}

func (lb *loadBalancer) nextBackend() *backendStatus {
	switch lb.policy {
	case "random":
		return lb.statuses[rand.Intn(len(lb.statuses))]
	default: // round_robin
		idx := atomic.AddUint32(&lb.index, 1) % uint32(len(lb.statuses))
		return lb.statuses[idx]
	}
}

func (lb *loadBalancer) errorHandler(backend *backendStatus) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		backend.mu.Lock()
		backend.alive = false
		backend.mu.Unlock()

		// 重试其他后端
		lb.ServeHTTP(w, r)
	}
}

// 健康检查实现
func (s *Server) healthCheckOld(status *backendStatus) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := http.Get(status.url.String() + "/health")
		alive := err == nil && resp.StatusCode == http.StatusOK

		status.mu.Lock()
		status.alive = alive
		status.mu.Unlock()
	}
}

func (s *Server) healthCheck(status *backendStatus, path string, intervalSec int) {
	if path == "" {
		path = "/health"
	}
	if intervalSec <= 0 {
		intervalSec = 15
	}
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := http.Get(status.url.String() + path)
		alive := err == nil && resp.StatusCode == http.StatusOK

		status.mu.Lock()
		status.alive = alive
		status.mu.Unlock()
	}
}
