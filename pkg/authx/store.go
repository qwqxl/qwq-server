package authx

import (
	"context"
	"fmt"
	"time"

	"qwqserver/pkg/cache"
)

// SessionStore 负责 Redis 中的会话管理
type SessionStore struct {
	client *cache.Client
	prefix string
}

// NewSessionStore 创建一个新的 SessionStore
func NewSessionStore(client *cache.Client, prefix string) *SessionStore {
	return &SessionStore{
		client: client,
		prefix: prefix,
	}
}

// CreateSession 在 Redis 中创建会话
func (s *SessionStore) CreateSession(ctx context.Context, userID, platform, deviceSign string, ttl time.Duration) error {
	key := s.buildKey(userID, platform, deviceSign)
	return s.client.Set(ctx, key, "active", ttl)
}

// GetSession 从 Redis 中获取会-话
func (s *SessionStore) GetSession(ctx context.Context, userID, platform, deviceSign string) (string, error) {
	key := s.buildKey(userID, platform, deviceSign)
	return s.client.Get(ctx, key)
}

// DeleteSession 从 Redis 中删除会话
func (s *SessionStore) DeleteSession(ctx context.Context, userID, platform, deviceSign string) error {
	key := s.buildKey(userID, platform, deviceSign)
	return s.client.Del(ctx, key)
}

// UpdateSessionExpiration 更新会话的过期时间
func (s *SessionStore) UpdateSessionExpiration(ctx context.Context, userID, platform, deviceSign string, ttl time.Duration) error {
	key := s.buildKey(userID, platform, deviceSign)
	_, err := s.client.Expire(ctx, key, ttl)
	if err != nil {
		return err
	}

	return err
}

// buildKey 构建 Redis 键
func (s *SessionStore) buildKey(userID, platform, deviceSign string) string {
	return fmt.Sprintf("%s:%s:%s:%s", s.prefix, userID, platform, deviceSign)
}
