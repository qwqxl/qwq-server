package cache

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"qwqserver/internal/config"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// Client 封装Redis客户端并提供安全的方法调用
type Client struct {
	Redis *redis.Client
}

// Pool 管理Redis连接池
type Pool struct {
	config *config.Cache
	client *redis.Client
	mu     sync.Mutex
}

var (
	instance *Pool
	once     sync.Once
)

// NewPool 创建全局唯一的Redis连接池实例 (单例模式)
func NewPool(cfgs ...*config.Cache) *Pool {
	once.Do(func() {
		if len(cfgs) > 0 {
			cfg := cfgs[0]

			client := redis.NewClient(&redis.Options{
				Addr:     cfg.Redis.Addr,
				Password: cfg.Redis.Password,
				DB:       cfg.Redis.DB,
			})
			instance = &Pool{
				config: cfg,
				client: client,
			}
		}
	})
	//instance.startHealthCheck(10 * time.Second)
	return instance
}

// GetClient 从连接池获取Redis客户端
func (rp *Pool) GetClient() (*Client, error) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	// 如果连接池已初始化，直接返回
	if rp.client != nil {
		return &Client{rp.client}, nil
	}

	// 初始化连接池
	opt := &redis.Options{
		Addr:            rp.config.Redis.Addr,
		Password:        rp.config.Redis.Password,
		DB:              rp.config.Redis.DB,
		PoolSize:        rp.config.Redis.PoolSize,
		MinIdleConns:    rp.config.Redis.MinIdleConns,
		MaxRetries:      rp.config.Redis.MaxRetries,
		DialTimeout:     rp.config.Redis.DialTimeout,
		ReadTimeout:     rp.config.Redis.ReadTimeout,
		WriteTimeout:    rp.config.Redis.WriteTimeout,
		ConnMaxIdleTime: rp.config.Redis.ConnMaxIdleTime,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			//log.Println("[Redis] ✅ New connection established")
			// OnConnect 回调函数，实现 Redis 自动连接建立时的日志记录 ✅
			return nil
		},
	}

	client := redis.NewClient(opt)

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "redis connection failed")
	}

	rp.client = client
	return &Client{client}, nil
}

// ReloadConfig 增加热更新配置能力
func (rp *Pool) ReloadConfig(newCfg *config.Cache) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	rp.config = newCfg
	_ = rp.Close()  // 关闭旧连接
	rp.client = nil // 强制下次GetClient时重建
}

// Close 关闭连接池
func (rp *Pool) Close() error {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	if rp.client != nil {
		err := rp.client.Close()
		rp.client = nil
		return err
	}
	return nil
}

// 基本操作方法
func (rc *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := rc.Redis.Set(ctx, key, value, expiration).Err(); err != nil {
		return errors.Wrapf(err, "set key %s failed", key)
	}
	return nil
}

func (rc *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.Redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrKeyNotFound
	} else if err != nil {
		return "", errors.Wrapf(err, "get key %s failed", key)
	}
	return val, nil
}

func (rc *Client) Del(ctx context.Context, keys ...string) error {
	if err := rc.Redis.Del(ctx, keys...).Err(); err != nil {
		return errors.Wrapf(err, "delete keys %v failed", keys)
	}
	return nil
}

func (rc *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := rc.Redis.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.Wrapf(err, "check exists for key %s failed", key)
	}
	return result > 0, nil
}

// 哈希表操作方法
func (rc *Client) HSet(ctx context.Context, key string, field string, value interface{}) error {
	if err := rc.Redis.HSet(ctx, key, field, value).Err(); err != nil {
		return errors.Wrapf(err, "hset key %s field %s failed", key, field)
	}
	return nil
}

func (rc *Client) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := rc.Redis.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", ErrFieldNotFound
	} else if err != nil {
		return "", errors.Wrapf(err, "hget key %s field %s failed", key, field)
	}
	return val, nil
}

func (rc *Client) HDel(ctx context.Context, key string, fields ...string) error {
	if err := rc.Redis.HDel(ctx, key, fields...).Err(); err != nil {
		return errors.Wrapf(err, "hdel key %s fields %v failed", key, fields)
	}
	return nil
}

func (rc *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := rc.Redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "hgetall key %s failed", key)
	}
	return result, nil
}

// 高级操作方法
func (rc *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := rc.Redis.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return false, errors.Wrapf(err, "setnx key %s failed", key)
	}
	return result, nil
}

func (rc *Client) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return errors.Wrapf(err, "json marshal for key %s failed", key)
	}
	return rc.Set(ctx, key, bs, expiration)
}

func (rc *Client) GetJSON(ctx context.Context, key string, target interface{}) error {
	val, err := rc.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), target)
}

func (rc *Client) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	result, err := rc.Redis.Expire(ctx, key, expiration).Result()
	if err != nil {
		return false, errors.Wrapf(err, "expire key %s failed", key)
	}
	return result, nil
}

func (rc *Client) Increment(ctx context.Context, key string, value int64) (int64, error) {
	result, err := rc.Redis.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "increment key %s failed", key)
	}
	return result, nil
}

// 连接管理方法
func (rc *Client) Ping(ctx context.Context) error {
	return rc.Redis.Ping(ctx).Err()
}

// 自定义错误类型
var (
	ErrKeyNotFound   = errors.New("key does not exist")
	ErrFieldNotFound = errors.New("field does not exist")
)

// 新增实用方法

// GetOrSet 获取缓存，不存在时设置缓存
func (rc *Client) GetOrSet(ctx context.Context, key string, fn func() (interface{}, time.Duration, error)) (string, error) {
	// 尝试获取缓存
	val, err := rc.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return "", err
	}

	// 缓存未命中，执行加载函数
	newVal, expiration, err := fn()
	if err != nil {
		return "", errors.Wrap(err, "getOrSet loader function failed")
	}

	// 设置新值
	if err := rc.Set(ctx, key, newVal, expiration); err != nil {
		return "", errors.Wrap(err, "set value in getOrSet failed")
	}

	// 返回新值（转换为字符串）
	switch v := newVal.(type) {
	case string:
		return v, nil
	default:
		return "", errors.New("unsupported value type in getOrSet")
	}
}

// WithLock 分布式锁执行函数
func (rc *Client) WithLock(ctx context.Context, lockKey string, ttl time.Duration, fn func() error) error {
	// 尝试获取锁
	locked, err := rc.SetNX(ctx, lockKey, "1", ttl)
	if err != nil {
		return errors.Wrap(err, "acquire lock failed")
	}
	if !locked {
		return errors.New("lock already held")
	}

	// 确保释放锁
	defer func() {
		_ = rc.Del(ctx, lockKey)
	}()

	// 执行受保护的函数
	return fn()
}

// ------- cache 新增 --------- //

// Incr 自增
func (rc *Client) Incr(ctx context.Context, key string) (int64, error) {
	result, err := rc.Redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "incr key %s failed", key)
	}
	return result, nil
}

// Decr 自减
func (rc *Client) Decr(ctx context.Context, key string) (int64, error) {
	result, err := rc.Redis.Decr(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "decr key %s failed", key)
	}
	return result, nil
}

// TTL 获取剩余过期时间
func (rc *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	result, err := rc.Redis.TTL(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "ttl key %s failed", key)
	}
	return result, nil
}

/* ---------------------- 分布式锁强化 ---------------------- */

// LockOld 获取分布式锁（安全版本）
func (rc *Client) LockOld(ctx context.Context, lockKey string, ttl time.Duration) (bool, error) {
	// 尝试获取锁
	locked, err := rc.SetNX(ctx, lockKey, "1", ttl)
	if err != nil {
		return false, errors.Wrap(err, "acquire lock failed")
	}
	if !locked {
		return false, errors.New("lock already held")
	}

	// 确保释放锁
	defer func() {
		_ = rc.Del(ctx, lockKey)
	}()

	return true, nil
}

// Lock 获取分布式锁（安全版本）
// 返回值：是否加锁成功、token值、error
func (rc *Client) Lock(ctx context.Context, lockKey string, ttl time.Duration) (bool, string, error) {
	token := uuid.NewString()

	locked, err := rc.SetNX(ctx, lockKey, token, ttl)
	if err != nil {
		return false, "", errors.Wrap(err, "acquire lock failed")
	}
	if !locked {
		return false, "", nil
	}

	return true, token, nil
}

// Unlock 释放分布式锁（安全版本）
func (rc *Client) Unlock(ctx context.Context, lockKey, token string) error {
	script := `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`
	result, err := rc.Redis.Eval(ctx, script, []string{lockKey}, token).Int64()
	if err != nil {
		return errors.Wrap(err, "unlock eval failed")
	}
	if result == 0 {
		return errors.New("unlock failed: token mismatch or lock not held")
	}
	return nil
}

// Pipeline 执行管道操作
func (rc *Client) Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) error {
	pipe := rc.Redis.Pipeline()
	if err := fn(pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	return err
}

// WithSafeLock 获取分布式锁，执行回调函数，并自动释放锁
func (rc *Client) WithSafeLock(ctx context.Context, lockKey string, ttl time.Duration, fn func(args ...any) error) error {
	// 尝试加锁
	token := uuid.NewString()
	locked, err := rc.SetNX(ctx, lockKey, token, ttl)
	if err != nil {
		return errors.Wrap(err, "acquire safe lock failed")
	}
	if !locked {
		return errors.New("safe lock already held")
	}

	// 确保释放锁
	defer func() {
		if err := rc.Unlock(ctx, lockKey, token); err != nil {
			log.Printf("⚠️ unlock error for key [%s]: %v", lockKey, err)
		}
	}()

	// 执行业务逻辑
	return fn()
}
