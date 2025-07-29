package app

// Initializer 初始化接口
type Initializer interface {
	Initialize() error
}

// Closer 关闭接口
type Closer interface {
	Close()
}

// HealthChecker 健康检查接口
type HealthChecker interface {
	CheckHealth() bool
}

// MetricsProvider 指标提供接口
type MetricsProvider interface {
	GetMetrics() map[string]interface{}
}
