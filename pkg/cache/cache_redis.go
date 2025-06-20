package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

type Config struct {
	Addr     string
	DB       int
	Password string
}

func InitRedis(c *Config) *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	return RedisClient
}

// 设置键值对（带过期时间）
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// 获取值
func Get(key string) (string, error) {
	ctx := context.Background()
	return RedisClient.Get(ctx, key).Result()
}

// 删除键
func Del(key string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}

// 检查键是否存在
func Exists(key string) (bool, error) {
	ctx := context.Background()
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// 设置哈希表字段
func HSet(key, field string, value interface{}) error {
	ctx := context.Background()
	return RedisClient.HSet(ctx, key, field, value).Err()
}

// 获取哈希表字段值
func HGet(key, field string) (string, error) {
	ctx := context.Background()
	return RedisClient.HGet(ctx, key, field).Result()
}

// 删除哈希表字段
func HDel(key string, fields ...string) error {
	ctx := context.Background()
	return RedisClient.HDel(ctx, key, fields...).Err()
}

// 获取哈希表所有字段
func HGetAll(key string) (map[string]string, error) {
	ctx := context.Background()
	return RedisClient.HGetAll(ctx, key).Result()
}
