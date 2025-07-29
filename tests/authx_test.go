package qwqtest

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"qwqserver/pkg/authx"
)

func setupTest(t *testing.T) (*authx.AuthX, *redis.Client) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "139.159.145.78:6379",
		Password: "redis_KF26xN",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("could not connect to redis: %v", err)
	}

	config := &authx.Config{
		JWTSecret:    "test-secret",
		AccessTTL:    time.Minute * 15,
		RefreshTTL:   time.Hour * 24 * 7,
		EnableSSO:    true,
		RedisPrefix:  "authx_test",
		CacheClient:  redisClient,
		Hooks:        authx.LifecycleHooks{},
		UseBlacklist: false,
	}

	ax, err := authx.New(config)
	if err != nil {
		t.Fatalf("failed to create AuthX instance: %v", err)
	}

	return ax, redisClient
}

func TestAuthX_LoginAndValidate(t *testing.T) {
	ax, redisClient := setupTest(t)
	defer redisClient.Close()

	loginInput := &authx.LoginInput{
		UserID:     "user-123",
		Platform:   "web",
		DeviceSign: "chrome-108",
	}

	tokenPair, err := ax.Login(context.Background(), loginInput)
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)

	claims, err := ax.ValidateToken(context.Background(), tokenPair.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, loginInput.UserID, claims.UserID)
	assert.Equal(t, loginInput.Platform, claims.PlatformSign)
	assert.Equal(t, loginInput.DeviceSign, claims.DeviceSign)
}
