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