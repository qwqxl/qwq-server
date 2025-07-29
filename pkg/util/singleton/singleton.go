package singleton

//
//import (
//	"errors"
//	"sync"
//)
//
//// SafeSingleton 读写安全的单例管理器
//type SafeSingleton[T any] struct {
//	instance    *T                 // 单例实例
//	initFn      func() (*T, error) // 初始化函数
//	mu          sync.RWMutex       // 读写锁
//	once        sync.Once          // 确保只初始化一次
//	initialized bool               // 标记是否已初始化
//}
//
//// NewSafeSingleton 创建安全的单例管理器
//func NewSafeSingleton[T any](initFn func() (*T, error)) *SafeSingleton[T] {
//	return &SafeSingleton[T]{
//		initFn: initFn,
//	}
//}
//
//// Get 安全读取单例实例
//func (s *SafeSingleton[T]) Get() (*T, error) {
//	s.mu.RLock()         // 获取读锁
//	defer s.mu.RUnlock() // 释放读锁
//
//	if !s.initialized {
//		// 双重检查锁定模式
//		s.mu.RUnlock()
//		s.mu.Lock()
//		defer s.mu.Unlock()
//
//		if !s.initialized {
//			// 初始化实例
//			instance, err := s.initFn()
//			if err != nil {
//				return nil, err
//			}
//			s.instance = instance
//			s.initialized = true
//		}
//
//		// 重新获取读锁
//		s.mu.RLock()
//	}
//
//	return s.instance, nil
//}
//
//// Update 安全更新单例实例（写操作期间阻塞所有读取）
//func (s *SafeSingleton[T]) Update(updateFn func(current *T) (*T, error)) error {
//	s.mu.Lock() // 获取写锁（独占）
//	defer s.mu.Unlock()
//
//	if !s.initialized {
//		return ErrNotInitialized
//	}
//
//	newInstance, err := updateFn(s.instance)
//	if err != nil {
//		return err
//	}
//
//	s.instance = newInstance
//	return nil
//}
//
//// MustGet 无错误处理的安全读取
//func (s *SafeSingleton[T]) MustGet() *T {
//	instance, err := s.Get()
//	if err != nil {
//		panic(err)
//	}
//	return instance
//}
//
//// Reload 重新加载单例（完全重新初始化）
//func (s *SafeSingleton[T]) Reload() error {
//	s.mu.Lock() // 获取写锁
//	defer s.mu.Unlock()
//
//	instance, err := s.initFn()
//	if err != nil {
//		return err
//	}
//
//	s.instance = instance
//	s.initialized = true
//	return nil
//}
//
//// ErrNotInitialized Errors
//var (
//	ErrNotInitialized = errors.New("singleton not initialized")
//)
