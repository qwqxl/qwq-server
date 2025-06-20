package cache

//
//import (
//	"context"
//	"crypto/rand"
//	"encoding/base64"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/redis/go-redis/v9"
//	"sync"
//	"time"
//)
//
//var (
//	redisClient *redis.Client
//	redisOnce   sync.Once
//)
//
//// 定义锁获取错误
//var (
//	ErrLockAcquireFailed = errors.New("lock acquire failed")
//)
//
//type Config struct {
//	Addr         string        `qwq-default:"localhost"`
//	Password     string        `qwq-default:"123456"`
//	DB           int           `default-value:"0"`
//	PoolSize     int           `default-value:"10"`
//	MinIdleConns int           `default-value:"10"`
//	MaxRetries   int           `default-value:"3"`
//	DialTimeout  time.Duration `default-value:"5s"`
//	ReadTimeout  time.Duration `default-value:"3s"`
//	WriteTimeout time.Duration `default-value:"3s"`
//	IdleTimeout  time.Duration `default-value:"5m"`
//}
//
//// InitRedis 初始化Redis连接池
//func InitRedis(cfg *Config) (*redis.Client, error) {
//	var err error
//
//	redisOnce.Do(func() {
//		options := &redis.Options{
//			Addr:         cfg.Addr,
//			Password:     cfg.Password,
//			DB:           cfg.DB,
//			PoolSize:     cfg.PoolSize,
//			MinIdleConns: cfg.MinIdleConns,
//			MaxRetries:   cfg.MaxRetries,
//			DialTimeout:  cfg.DialTimeout,
//			ReadTimeout:  cfg.ReadTimeout,
//			WriteTimeout: cfg.WriteTimeout,
//			// 使用正确的字段名 ConnMaxIdleTime
//			ConnMaxIdleTime: cfg.IdleTimeout,
//		}
//
//		client := redis.NewClient(options)
//
//		// 测试连接
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		_, err = client.Ping(ctx).Result()
//		if err != nil {
//			return
//		}
//
//		redisClient = client
//	})
//
//	return redisClient, err
//}
//
//// CloseRedis 关闭Redis连接
//func CloseRedis() error {
//	if redisClient != nil {
//		return redisClient.Close()
//	}
//	return nil
//}
//
//// Set 设置缓存
//func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
//	return redisClient.Set(ctx, key, value, expiration).Err()
//}
//
//// Keys 查询匹配的key
//func Keys(ctx context.Context, pattern string) ([]string, error) {
//	return redisClient.Keys(ctx, pattern).Result()
//}
//
//// Get 获取缓存
//func Get(ctx context.Context, key string) (string, error) {
//	return redisClient.Get(ctx, key).Result()
//}
//
//// GetBytes 获取字节数组缓存
//func GetBytes(ctx context.Context, key string) ([]byte, error) {
//	return redisClient.Get(ctx, key).Bytes()
//}
//
//// GetJSON 获取JSON缓存并解析到结构体
//func GetJSON(ctx context.Context, key string, dest interface{}) error {
//	data, err := GetBytes(ctx, key)
//	if err != nil {
//		return err
//	}
//	return json.Unmarshal(data, dest)
//}
//
//// SetJSON 将结构体序列化为JSON并设置缓存
//func SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
//	data, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	return Set(ctx, key, data, expiration)
//}
//
//// Exists 检查键是否存在
//func Exists(ctx context.Context, key string) (bool, error) {
//	result, err := redisClient.Exists(ctx, key).Result()
//	return result == 1, err
//}
//
//// Delete 删除缓存
//func Delete(ctx context.Context, keys ...string) error {
//	return redisClient.Del(ctx, keys...).Err()
//}
//
//// HSet 设置哈希字段值
//func HSet(ctx context.Context, key string, field string, value ...interface{}) error {
//	return redisClient.HSet(ctx, key, value...).Err()
//}
//
//// HGet 获取哈希字段值
//func HGet(ctx context.Context, key string, field string) (string, error) {
//	return redisClient.HGet(ctx, key, field).Result()
//}
//
//// HGetAll 获取所有哈希字段值
//func HGetAll(ctx context.Context, key string) (map[string]string, error) {
//	return redisClient.HGetAll(ctx, key).Result()
//}
//
//// HDel 删除哈希字段
//func HDel(ctx context.Context, key string, fields ...string) error {
//	return redisClient.HDel(ctx, key, fields...).Err()
//}
//
//// HMSet 设置多个哈希字段值
//func HMSet(ctx context.Context, key string, fields map[string]interface{}) error {
//	return redisClient.HMSet(ctx, key, fields).Err()
//}
//
//// Incr 自增
//func Incr(ctx context.Context, key string) (int64, error) {
//	return redisClient.Incr(ctx, key).Result()
//}
//
//// Decr 自减
//func Decr(ctx context.Context, key string) (int64, error) {
//	return redisClient.Decr(ctx, key).Result()
//}
//
//// Expire 设置过期时间
//func Expire(ctx context.Context, key string, expiration time.Duration) error {
//	return redisClient.Expire(ctx, key, expiration).Err()
//}
//
//// TTL 获取剩余过期时间
//func TTL(ctx context.Context, key string) (time.Duration, error) {
//	return redisClient.TTL(ctx, key).Result()
//}
//
//// generateToken 生成随机令牌
//func generateToken() (string, error) {
//	b := make([]byte, 16)
//	_, err := rand.Read(b)
//	if err != nil {
//		return "", err
//	}
//	return base64.URLEncoding.EncodeToString(b), nil
//}
//
//// Lock 获取分布式锁（安全版本）
//func Lock(ctx context.Context, key string, expiration time.Duration) (token string, err error) {
//	token, err = generateToken()
//	if err != nil {
//		return "", fmt.Errorf("生成token失败: %w", err)
//	}
//
//	ok, err := redisClient.SetNX(ctx, key, token, expiration).Result()
//	if err != nil {
//		return "", err
//	}
//	if !ok {
//		return "", ErrLockAcquireFailed
//	}
//	return token, nil
//}
//
//// Unlock 释放分布式锁（安全版本）
//func Unlock(ctx context.Context, key string, token string) error {
//	script := `
//	if redis.call("GET", KEYS[1]) == ARGV[1] then
//		return redis.call("DEL", KEYS[1])
//	else
//		return 0
//	end
//	`
//	result, err := redisClient.Eval(ctx, script, []string{key}, token).Int64()
//	if err != nil {
//		return err
//	}
//	if result == 0 {
//		return errors.New("解锁失败：锁不存在或token不匹配")
//	}
//	return nil
//}
//
//// Pipeline 执行管道操作
//func Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) error {
//	_, err := redisClient.Pipelined(ctx, fn)
//	return err
//}
//
//// WithLock 使用分布式锁执行操作（安全版本）
//func WithLock(ctx context.Context, key string, expiration time.Duration, fn func() error) error {
//	lockKey := fmt.Sprintf("lock:%s", key)
//	token, err := Lock(ctx, lockKey, expiration)
//	if err != nil {
//		if errors.Is(err, ErrLockAcquireFailed) {
//			return errors.New("资源被锁定")
//		}
//		return fmt.Errorf("获取锁失败: %w", err)
//	}
//
//	defer func() {
//		if unlockErr := Unlock(ctx, lockKey, token); unlockErr != nil {
//			// 记录日志：logger.Error("释放锁失败", "key", lockKey, "error", unlockErr)
//		}
//	}()
//
//	return fn()
//}
