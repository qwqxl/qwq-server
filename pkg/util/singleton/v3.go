package singleton

import (
	"sync"
	"sync/atomic"
)

// Singleton 泛型单例容器，提供线程安全的单例管理
// T: 单例实例类型
type Singleton[T any] struct {
	mu       sync.Mutex        // 保护初始化过程
	once     sync.Once         // 确保初始化只执行一次
	instance atomic.Pointer[T] // 原子指针保证可见性
}

// Get 获取单例实例（双重检查锁定优化版）
// 参数:
//
//	builder - 创建实例的函数，返回实例指针及可能的错误
//
// 返回:
//
//	*T - 单例实例指针
//	error - 构造函数返回的错误
//
// 注意: 当实例已存在时，builder不会被调用
func (s *Singleton[T]) Get(builder func() (*T, error)) (*T, error) {
	// 第一重检查：快速路径
	if instance := s.instance.Load(); instance != nil {
		return instance, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 第二重检查：避免锁竞争后重复创建
	if instance := s.instance.Load(); instance != nil {
		return instance, nil
	}

	// 初始化实例
	instance, err := builder()
	if err != nil {
		return nil, err
	}
	s.instance.Store(instance)
	return instance, nil
}

// MustGet 安全获取单例实例（构造函数出错时panic）
// 参数:
//
//	builder - 创建实例的函数
//
// 返回:
//
//	*T - 单例实例指针
func (s *Singleton[T]) MustGet(builder func() (*T, error)) *T {
	instance, err := s.Get(builder)
	if err != nil {
		panic("singleton initialization failed: " + err.Error())
	}
	return instance
}

// IsInitialized 检查单例是否已初始化
// 返回:
//
//	bool - true表示实例已初始化
func (s *Singleton[T]) IsInitialized() bool {
	return s.instance.Load() != nil
}

// Reset 重置单例状态（并发安全）
// 注意: 仅应在测试或特殊场景使用，重置后下次Get将重新创建实例
func (s *Singleton[T]) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.instance.Store(nil)
	s.once = sync.Once{} // 重置sync.Once
}

// NewSingleton 创建新的泛型单例容器
// 返回:
//
//	*Singleton[T] - 单例容器实例
func NewSingleton[T any]() *Singleton[T] {
	return &Singleton[T]{}
}
