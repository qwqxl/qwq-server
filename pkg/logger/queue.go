package logger

//
//import "sync"
//
//type AsyncQueue struct {
//	asyncEnabled   bool           // 是否启用异步
//	asyncQueue     chan string    // 异步消息队列
//	asyncQueueSize int            // 队列大小（默认1000）
//	asyncWG        sync.WaitGroup // 等待异步处理完成（关闭前用）
//}
//
//// 启用异步写入
//func WithAsyncQueue(size int) ConfigOption {
//	return func(c *LogConfig) {
//		c.asyncEnabled = true
//		if size <= 0 {
//			size = 1000 // 默认队列大小
//		}
//		c.asyncQueueSize = size
//	}
//}
//
//// 关闭异步写入
//func WithDisableAsync() ConfigOption {
//	return func(c *LogConfig) {
//		c.asyncEnabled = false
//	}
//}
