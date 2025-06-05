package singleton

import (
	"sync"
)

// Singleton 泛型单例容器
type Singleton[T any] struct {
	once     sync.Once // 确保初始化只执行一次
	instance *T        // 泛型实例指针
}

// Get 获取单例实例
// 参数: builder - 创建实例的函数
// 返回: 单例实例指针
func (s *Singleton[T]) Get(builder func() *T) *T {
	s.once.Do(func() {
		s.instance = builder()
	})
	return s.instance
}

// NewSingleton 创建新的泛型单例容器
func NewSingleton[T any]() *Singleton[T] {
	return &Singleton[T]{}
}
