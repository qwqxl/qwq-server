

    gin-ddos-shield/
    ├── core/
    │   ├── limiter.go          # 限流核心接口
    │   ├── token_bucket.go     # 令牌桶实现
    │   ├── concurrency.go      # 并发控制
    │   ├── adaptive.go         # 自适应限流
    │   └── distributed.go      # 分布式限流
    ├── protection/
    │   ├── reputation.go       # IP信誉系统
    │   ├── challenge.go        # 挑战机制
    │   ├── behavior.go         # 行为分析
    │   └── intelligence.go     # 威胁情报
    ├── middleware/
    │   └── gin.go              # Gin中间件
    ├── config/
    │   └── config.go           # 配置管理
    └── utils/
    ├── ip.go               # IP工具
    └── metrics.go          # 监控指标