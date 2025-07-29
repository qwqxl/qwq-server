# Cache Package Usage

`cache` 库中主要方法的使用示例，按功能模块组织，每段示例都是“如何调用 + 解释用途”。你可以直接 copy 到你的项目中作为参考 🚀

***

## 1. 获取 Redis 客户端实例

```go

cfg := &config.Cache{ /* 加载你的配置，比如从 env 或 文件 */ }
pool := cache.NewPool(cfg)

clientWrapper, err := pool.GetClient()
if err != nil {
    log.Fatalf("获取 Redis 客户端失败：%v", err)
}
client := clientWrapper // 便于调用
```

***

## 2. 基本字符串操作（Set / Get / Del / Exists）

```go
ctx := context.Background()
err = client.Set(ctx, "foo", "bar", 10*time.Minute)
if err != nil {
    log.Printf("Set 出错：%v", err)
}

value, err := client.Get(ctx, "foo")
if err != nil {
    if errors.Is(err, cache.ErrKeyNotFound) {
        log.Println("key 不存在")
    } else {
        log.Printf("Get 出错：%v", err)
    }
} else {
    fmt.Println("获取 foo =", value)
}

exists, _ := client.Exists(ctx, "foo")
fmt.Println("是否存在 foo:", exists)

_ = client.Del(ctx, "foo")
```

***

## 3. 哈希（Hash）命令示例

```go
ctx := context.Background()
_ = client.HSet(ctx, "user:1", "name", "Alice")
name, err := client.HGet(ctx, "user:1", "name")
if err != nil {
    log.Println("未找到字段或出错")
} else {
    fmt.Println("user:1 name =", name)
}

all, _ := client.HGetAll(ctx, "user:1")
fmt.Println("完整哈希内容：", all)

_ = client.HDel(ctx, "user:1", "name")
```

***

## 4. 条件设置 & TTL 操作（SetNX / Expire / TTL）

```go
ctx := context.Background()
ok, _ := client.SetNX(ctx, "once_key", "v", 1*time.Hour)
fmt.Println("是否首次设置成功:", ok)

ok2, _ := client.Expire(ctx, "once_key", 30*time.Minute)
fmt.Println("是否设置新 TTL:", ok2)

ttl, _ := client.TTL(ctx, "once_key")
fmt.Println("剩余 TTL:", ttl)
```

***

## 5. 增量操作（自增 / 自减）

```go
ctx := context.Background()
count, _ := client.Incr(ctx, "counter")
fmt.Println("增加后的值：", count)

count, _ = client.Decr(ctx, "counter")
fmt.Println("减少后的值：", count)
```

***

## 6. 缓存穿透保护（GetOrSet）

```go
ctx := context.Background()
val, err := client.GetOrSet(ctx, "cache:user:2", func() (interface{}, time.Duration, error) {
    // 模拟加载过程（例如数据库查询）
    return "Bob", 15 * time.Minute, nil
})
if err != nil {
    log.Println("GetOrSet 出错：", err)
} else {
    fmt.Println("最终结果:", val) // 来自缓存或函数
}
```

***

## 7. 管道操作（Pipeline）

```go
ctx := context.Background()
err := client.Pipeline(ctx, func(pipe redis.Pipeliner) error {
    pipe.Set(ctx, "a", 1, 0)
    pipe.Set(ctx, "b", 2, 0)
    return nil
})
if err != nil {
    log.Println("Pipeline 出错：", err)
}
```

***

## 8. 简易分布式锁（旧版：WithLock）

```go
err := client.WithLock(ctx, "old_lock", 10*time.Second, func() error {
    fmt.Println("旧版锁逻辑执行中")
    return nil
})
if err != nil {
    log.Println("WithLock 错误：", err)
}
```

***

## 9. 安全分布式锁（UUID + WithSafeLock）

```go
err := client.WithSafeLock(ctx, "safe_task_lock", 15*time.Second, func(args ...any) error {
    // 回调没有使用 args，但示范可扩展
    fmt.Println("安全锁机制下执行任务")
    return nil
})
if err != nil {
    log.Println("执行失败或锁未获得：", err)
}
```

***

## 10. 手动 Lock / Unlock 使用样例

```go
locked, token, err := client.Lock(ctx, "manual_lock", 10*time.Second)
if err != nil {
    log.Fatalf("Lock 出错: %v", err)
}
if !locked {
    log.Println("锁已被其他实例持有，退出")
    return
}
defer func(){
    if extErr := client.Unlock(ctx, "manual_lock", token); extErr != nil {
        log.Println("Unlock 出错：", extErr)
    }
}()

// 临界区代码
fmt.Println("我手动控制锁")
```

***

## 11. 关闭 Redis 连接池

```go
if err := pool.Close(); err != nil {
    log.Printf("关闭 Redis 失败: %v", err)
}
```

***

### ✅ 总结

* **获取客户端**：`NewPool + GetClient`

* **基础操作**：Set/Get/Del/Exists

* **扩展功能**：HSet/HGet, SetNX, Expire、Incr/Decr

* **高级功能**：GetOrSet, Pipeline

* **锁机制**：

  * `WithLock`：简单实现

  * `WithSafeLock`, `Lock` + `Unlock`：安全版本

* **关闭连接非必需但推荐释放资源**：`pool.Close()`

如需扩展错误处理重试、日志中间件、Cluster/Sentinel 支持等功能，随时告诉我，我可以继续帮你完善 🔨🤖🔧
