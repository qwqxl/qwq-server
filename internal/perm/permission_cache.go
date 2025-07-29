package perm

//
//import (
//	"context"
//	"sync"
//	"time"
//
//	"qwqserver/internal/model"
//)
//
//// PermissionCache 带缓存的权限检查器
//type PermissionCache struct {
//	manager     *PermissionManager
//	cache       map[uint]map[model.PermissionKey]bool
//	cacheExpiry map[uint]time.Time
//	cacheTTL    time.Duration
//	mu          sync.RWMutex
//}
//
//func NewPermissionCache(manager *PermissionManager, ttl time.Duration) *PermissionCache {
//	return &PermissionCache{
//		manager:     manager,
//		cache:       make(map[uint]map[model.PermissionKey]bool),
//		cacheExpiry: make(map[uint]time.Time),
//		cacheTTL:    ttl,
//	}
//}
//
//// HasPermission 检查用户是否拥有指定权限（带缓存）
//func (c *PermissionCache) HasPermission(
//	ctx context.Context,
//	userID uint,
//	permissionKey model.PermissionKey,
//) (bool, error) {
//	c.mu.RLock()
//
//	// 检查缓存是否有效
//	if perms, ok := c.cache[userID]; ok {
//		if expiry, ok := c.cacheExpiry[userID]; ok && time.Now().Before(expiry) {
//			if hasPerm, ok := perms[permissionKey]; ok {
//				c.mu.RUnlock()
//				return hasPerm, nil
//			}
//		}
//	}
//
//	c.mu.RUnlock()
//
//	// 缓存未命中，查询数据库
//	hasPerm, err := c.manager.CheckPermission(ctx, userID, permissionKey)
//	if err != nil {
//		return false, err
//	}
//
//	// 更新缓存
//	c.mu.Lock()
//	defer c.mu.Unlock()
//
//	if _, ok := c.cache[userID]; !ok {
//		c.cache[userID] = make(map[model.PermissionKey]bool)
//	}
//
//	c.cache[userID][permissionKey] = hasPerm
//	c.cacheExpiry[userID] = time.Now().Add(c.cacheTTL)
//
//	return hasPerm, nil
//}
//
//// InvalidateUserCache 使指定用户的缓存失效
//func (c *PermissionCache) InvalidateUserCache(userID uint) {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//
//	delete(c.cache, userID)
//	delete(c.cacheExpiry, userID)
//}
//
//// RefreshUserPermissions 刷新用户权限缓存
//func (c *PermissionCache) RefreshUserPermissions(ctx context.Context, userID uint) error {
//	// 获取用户所有权限
//	permissions, err := c.manager.GetUserPermissions(ctx, userID)
//	if err != nil {
//		return err
//	}
//
//	// 构建权限映射
//	permMap := make(map[model.PermissionKey]bool)
//	for _, perm := range permissions {
//		permMap[perm.Key] = true
//	}
//
//	// 更新缓存
//	c.mu.Lock()
//	defer c.mu.Unlock()
//
//	c.cache[userID] = permMap
//	c.cacheExpiry[userID] = time.Now().Add(c.cacheTTL)
//
//	return nil
//}
