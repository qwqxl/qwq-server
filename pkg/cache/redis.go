package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

type Config struct {
	Addr         string        `qwq-default:"localhost"`
	Password     string        `qwq-default:"123456"`
	DB           int           `default-value:"0"`
	PoolSize     int           `default-value:"10"`
	MinIdleConns int           `default-value:"10"`
	MaxRetries   int           `default-value:"3"`
	DialTimeout  time.Duration `default-value:"5s"`
	ReadTimeout  time.Duration `default-value:"3s"`
	WriteTimeout time.Duration `default-value:"3s"`
	IdleTimeout  time.Duration `default-value:"5m"`
}

// InitRedis 初始化Redis连接池
func InitRedis(cfg *Config) (*redis.Client, error) {
	var err error

	redisOnce.Do(func() {
		options := &redis.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			MaxRetries:   cfg.MaxRetries,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		}

		client := redis.NewClient(options)

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Ping(ctx).Result()
		if err != nil {
			return
		}

		redisClient = client
		// redis 连接成功
	})

	return redisClient, err
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}

func Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return redisClient.Keys(ctx, pattern)
	//return redisClient.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

// GetBytes 获取字节数组缓存
func GetBytes(ctx context.Context, key string) ([]byte, error) {
	return redisClient.Get(ctx, key).Bytes()
}

// GetJSON 获取JSON缓存并解析到结构体
func GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := GetBytes(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// SetJSON 将结构体序列化为JSON并设置缓存
func SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Set(ctx, key, data, expiration)
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := redisClient.Exists(ctx, key).Result()
	return result == 1, err
}

// Delete 删除缓存
func Delete(ctx context.Context, keys ...string) error {
	return redisClient.Del(ctx, keys...).Err()
}

// HSet 设置哈希字段值
func HSet(ctx context.Context, key string, field string, value ...interface{}) error {
	return redisClient.HSet(ctx, key, value...).Err()
}

// HGet 获取哈希字段值
func HGet(ctx context.Context, key string, field string) (string, error) {
	return redisClient.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段值
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return redisClient.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return redisClient.HDel(ctx, key, fields...).Err()
}

// HMSET 设置多个哈希字段值
func HMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	return redisClient.HMSet(ctx, key, fields).Err()
}

// Incr 自增
func Incr(ctx context.Context, key string) (int64, error) {
	return redisClient.Incr(ctx, key).Result()
}

// Decr 自减
func Decr(ctx context.Context, key string) (int64, error) {
	return redisClient.Decr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return redisClient.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return redisClient.TTL(ctx, key).Result()
}

// Lock 获取分布式锁
func Lock(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return redisClient.SetNX(ctx, key, value, expiration).Result()
}

// Unlock 释放分布式锁
func Unlock(ctx context.Context, key string) error {
	return Delete(ctx, key)
}

// Pipeline 执行管道操作
func Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) error {
	_, err := redisClient.Pipelined(ctx, fn)
	return err
}

// WithLock 使用分布式锁执行操作
func WithLock(ctx context.Context, key string, expiration time.Duration, fn func() error) error {
	lockKey := fmt.Sprintf("lock:%s", key)

	// 尝试获取锁
	locked, err := Lock(ctx, lockKey, "1", expiration)
	if err != nil {
		return fmt.Errorf("获取锁失败: %w", err)
	}
	if !locked {
		return errors.New("资源被锁定")
	}

	// 确保释放锁
	defer func() {
		if err := Unlock(ctx, lockKey); err != nil {
			//logger.Error("释放锁失败", "key", lockKey, "error", err)
		}
	}()

	// 执行操作
	return fn()
}
