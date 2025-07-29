package app

import (
	"sync"
)

// Manager 单例管理器
type Manager struct {
	instances map[string]interface{}
	mu        sync.RWMutex
}

var (
	managerInstance *Manager
	once            sync.Once
)

// GetManager 获取单例管理器实例
func GetManager() *Manager {
	once.Do(func() {
		managerInstance = &Manager{
			instances: make(map[string]interface{}),
		}
	})
	return managerInstance
}

// GetOrCreate 获取或创建单例实例
func (m *Manager) GetOrCreate(key string, constructor func() interface{}) interface{} {
	// 第一次检查 - 读锁
	m.mu.RLock()
	instance, exists := m.instances[key]
	m.mu.RUnlock()

	if exists {
		return instance
	}

	// 获取写锁准备创建
	m.mu.Lock()
	defer m.mu.Unlock()

	// 第二次检查 - 防止其他goroutine已经创建
	if instance, exists := m.instances[key]; exists {
		return instance
	}

	// 创建新实例
	instance = constructor()
	m.instances[key] = instance

	return instance
}

// CloseAll 关闭所有单例实例
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, instance := range m.instances {
		if closer, ok := instance.(Closer); ok {
			closer.Close()
		}
		delete(m.instances, key)
	}
}
