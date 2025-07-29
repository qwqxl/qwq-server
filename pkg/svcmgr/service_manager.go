package svcmgr

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Svc 定义服务接口
type Svc interface {
	// Run 持续运行服务，直到ctx被取消或返回错误
	Run(ctx context.Context) error
	// Name 返回服务名称
	Name() string
}

// Mgr 服务管理器
type Mgr struct {
	svcs    []Svc
	wg      sync.WaitGroup
	cancel  context.CancelFunc
	errChan chan error
	mu      sync.Mutex
	running bool

	// 配置项
	stopTimeout    time.Duration
	restartDelay   time.Duration
	maxRestarts    int
	restartOnError bool
}

// Config 管理器配置
type Config struct {
	StopTimeout    time.Duration // 停止超时(默认10s)
	RestartDelay   time.Duration // 重启延迟(默认1s)
	MaxRestarts    int           // 最大重启次数(默认3)
	RestartOnError bool          // 是否自动重启(默认true)
	ErrBufSize     int           // 错误缓冲大小(默认20)
}

// New 创建服务管理器
func New(cfg ...Config) *Mgr {
	c := Config{
		StopTimeout:    10 * time.Second,
		RestartDelay:   time.Second,
		MaxRestarts:    3,
		RestartOnError: true,
		ErrBufSize:     20,
	}
	if len(cfg) > 0 {
		c = cfg[0]
	}

	return &Mgr{
		errChan:        make(chan error, c.ErrBufSize),
		stopTimeout:    c.StopTimeout,
		restartDelay:   c.RestartDelay,
		maxRestarts:    c.MaxRestarts,
		restartOnError: c.RestartOnError,
	}
}

// Add 添加服务
func (m *Mgr) Add(s Svc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.svcs = append(m.svcs, s)
}

// Run 启动并管理所有服务
func (m *Mgr) Run() error {
	if err := m.Start(); err != nil {
		return err
	}
	m.Wait()
	return nil
}

// Start 启动所有服务
func (m *Mgr) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return errors.New("services already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.running = true

	for _, svc := range m.svcs {
		m.wg.Add(1)
		go m.runService(ctx, svc)
	}

	return nil
}

// runService 运行单个服务(带自动重启逻辑)
func (m *Mgr) runService(ctx context.Context, svc Svc) {
	defer m.wg.Done()

	var restarts int
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := svc.Run(ctx)
			if err != nil {
				m.errChan <- fmt.Errorf("%s failed: %w", svc.Name(), err)
			}

			// 检查是否需要重启
			if ctx.Err() != nil || !m.shouldRestart(svc.Name(), restarts) {
				return
			}

			restarts++
			time.Sleep(m.restartDelay)
		}
	}
}

// shouldRestart 判断服务是否应该重启
func (m *Mgr) shouldRestart(name string, restarts int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.restartOnError {
		return false
	}

	if m.maxRestarts > 0 && restarts >= m.maxRestarts {
		m.errChan <- fmt.Errorf("%s reached max restarts (%d)", name, m.maxRestarts)
		return false
	}

	return true
}

// Stop 停止所有服务
func (m *Mgr) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running || m.cancel == nil {
		return
	}

	m.cancel()

	// 等待服务停止或超时
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(m.stopTimeout):
		m.errChan <- fmt.Errorf("stop timeout after %s", m.stopTimeout)
	}

	m.running = false
}

// Wait 等待所有服务结束
func (m *Mgr) Wait() {
	m.wg.Wait()
	close(m.errChan)
}

// ErrChan 返回错误通道
func (m *Mgr) ErrChan() <-chan error {
	return m.errChan
}

// IsRunning 检查是否在运行
func (m *Mgr) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

// Count 返回服务数量
func (m *Mgr) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.svcs)
}

// UpdateConfig 更新管理器配置
func (m *Mgr) UpdateConfig(cfg Config) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cfg.StopTimeout > 0 {
		m.stopTimeout = cfg.StopTimeout
	}
	if cfg.RestartDelay > 0 {
		m.restartDelay = cfg.RestartDelay
	}
	if cfg.MaxRestarts >= 0 {
		m.maxRestarts = cfg.MaxRestarts
	}
	// RestartOnError 可以设置为 false
	m.restartOnError = cfg.RestartOnError
}
