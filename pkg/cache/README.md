```go
// internal/service/user_service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
	"qwqserver/pkg/cache"
	"qwqserver/pkg/logger"
)

const (
	userCachePrefix    = "user:"
	userCacheDuration  = 30 * time.Minute
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	cacheKey := userCachePrefix + string(id)
	
	// 尝试从缓存获取
	var user model.User
	if exists, _ := cache.GetJSON(ctx, cacheKey, &user); exists == nil {
		logger.Debug("从缓存获取用户", "id", id)
		return &user, nil
	}
	
	// 缓存未命中，从数据库获取
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// 设置缓存
	if err := cache.SetJSON(ctx, cacheKey, user, userCacheDuration); err != nil {
		logger.Warn("设置用户缓存失败", "id", id, "error", err)
	}
	
	return &user, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID uint, profile *model.UserProfile) error {
	// 使用分布式锁确保并发安全
	return cache.WithLock(ctx, "user_profile_update:"+string(userID), 10*time.Second, func() error {
		// 获取当前用户
		user, err := s.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		
		// 更新资料
		user.Profile = *profile
		if err := s.userRepo.Update(ctx, &user); err != nil {
			return err
		}
		
		// 清除缓存
		cacheKey := userCachePrefix + string(userID)
		if err := cache.Delete(ctx, cacheKey); err != nil {
			logger.Warn("清除用户缓存失败", "id", userID, "error", err)
		}
		
		return nil
	})
}
```



## New Cache Redis Use

```go
// 1. 创建配置
cfg := &cache.Config{
	Addr:         "localhost:6379",
	Password:     "mypassword",
	DB:           0,
	PoolSize:     20,
	MinIdleConns: 5,
	MaxRetries:   3,
	DialTimeout:  5 * time.Second,
	IdleTimeout:  5 * time.Minute,
}

// 2. 初始化全局连接池
pool := cache.NewRedisPool(cfg)

// 3. 在需要的地方创建客户端
func processRequest() {
	// 每次使用创建新客户端（轻量级）
	client, err := cache.NewRedisClient(pool)
	if err != nil {
		// 处理错误
	}
	defer func() {
		// 注意：这里不需要关闭客户端，因为共享连接池
	}()

	// 执行操作
	ctx := context.Background()
	err = client.Set(ctx, "user:1001", "john", 24*time.Hour)
	val, err := client.Get(ctx, "user:1001")
	
	// 使用高级功能
	locked, err := client.SetNX(ctx, "lock:resource", "1", 10*time.Second)
	if locked {
		defer client.Del(ctx, "lock:resource")
	}
}

// 4. 程序退出时关闭连接池
func main() {
	defer pool.Close()
	// ...
}
```