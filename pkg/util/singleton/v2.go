package singleton

//
//import (
//	"errors"
//	"sync"
//	"sync/atomic"
//)
//
//var (
//	ErrNotInitialized = errors.New("singleton not initialized")
//)
//
//type Singleton[T any] struct {
//	once     sync.Once
//	instance atomic.Value // *T
//}
//
//// NewSingleton 创建单例容器
//// builder: 实例构造函数 (必须非nil)
//func NewSingleton[T any](builder func() *T) *Singleton[T] {
//	if builder == nil {
//		panic("singleton: builder cannot be nil")
//	}
//
//	s := &Singleton[T]{}
//	s.once.Do(func() {
//		s.instance.Store(builder())
//	})
//	return s
//}
//
//// MustGet 获取单例实例 (未初始化时panic)
//func (s *Singleton[T]) MustGet() *T {
//	instance := s.instance.Load()
//	if instance == nil {
//		panic(ErrNotInitialized)
//	}
//	return instance.(*T)
//}
//
//// Get 安全获取单例实例
//// 返回: (实例指针, 错误)
//func (s *Singleton[T]) Get() (*T, error) {
//	if instance := s.instance.Load(); instance != nil {
//		return instance.(*T), nil
//	}
//	return nil, ErrNotInitialized
//}
//
//// IsInitialized 检查是否已初始化
//func (s *Singleton[T]) IsInitialized() bool {
//	return s.instance.Load() != nil
//}
