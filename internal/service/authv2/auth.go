package authv2

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"qwqserver/pkg/cache"
	"strings"
	"time"
)

// 路由路径
const (
	LoginPath    = "/login"
	RegisterPath = "/register"
	ProfilePath  = "/profile"
	LogoutPath   = "/logout"
	Identity     = "identity"
)

var (
	AccessSecret  = []byte("access_secret")
	RefreshSecret = []byte("refresh_secret")
)

// Token 有效期配置
const (
	AccessTokenExpire  = 15 * time.Minute   // AccessToken 有效期
	RefreshTokenExpire = 7 * 24 * time.Hour // RefreshToken 有效期
	JWTSecretKeyEnv    = "default_secret_key_please_change_in_production"
)

// JWT Claim 中的自定义字段
const (
	JwtCustomClaimsDeviceID  = "device_id"
	JwtCustomClaimsUserID    = "user_id"
	JwtCustomClaimsIpaddress = "ipaddress"
)

// Redis Key 前缀
const UserSessionCachePrefix = "user_session"

// Gin 上下文中设置的 Key
const (
	ContextKeyUserID   = "user_id"
	ContextKeyDeviceID = "device_id"
)

// Header 常量
const (
	HeaderAuthorization   = "Authorization"
	HeaderRefreshToken    = "X-Refresh-Token"
	HeaderDeviceID        = "X-Device-Id"
	HeaderNewAccessToken  = "New-Access-Token"
	HeaderNewRefreshToken = "New-Refresh-Token"
)

// Redis 字段常量
const (
	FieldRefreshToken = "refresh_token"
	FieldLastActive   = "last_active"
	FieldDeviceInfo   = "device_info"
)

var JWTSecretKey = []byte(JWTSecretKeyEnv)

type IdentityResult struct {
	UserID    string `json:"user_id"`
	DeviceID  string `json:"device_id"`
	IPAddress string `json:"ipaddress"`
	Status    bool   `json:"status"`
}

// 登录请求结构体
type LoginRequest struct {
	Name      string `json:"Name"`
	Password  string `json:"password"`
	DeviceID  string `json:"device_id"`
	IPAddress string `json:"ipaddress"`
}

// RegisterRequest 用户注册请求结构体
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// 自定义 JWT Claims
type CustomClaims struct {
	UserID    string `json:"user_id"`
	DeviceID  string `json:"device_id"`
	IPAddress string `json:"ipaddress"`
	jwt.RegisteredClaims
}

// 构建 Redis 的键名
func UserSessionCachePrefixToString(uid, deviceID, ipaddress string) string {
	return fmt.Sprintf("%s:%s:%s:%s", UserSessionCachePrefix, uid, deviceID, ipaddress)
}

// 清除同用户其他设备登录信息（单端登录）
func TerminateOtherSessions(userID, currentDeviceID string) error {
	ctx := context.Background()
	pattern := fmt.Sprintf("%s:%s:*", UserSessionCachePrefix, userID)
	keys, err := cache.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		parts := strings.Split(key, ":")
		if len(parts) < 4 {
			continue
		}
		if parts[2] != currentDeviceID {
			if err := cache.Delete(ctx, key); err != nil {
				return err
			}
		}
	}
	return nil
}

// 生成 AccessToken 和 RefreshToken
func GenerateTokenPair(userID, deviceID, ipaddr string) (string, string, error) {
	now := time.Now()

	accessClaims := CustomClaims{
		UserID:    userID,
		DeviceID:  deviceID,
		IPAddress: ipaddr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(JWTSecretKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := accessClaims
	refreshClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(RefreshTokenExpire))

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(JWTSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// 保存 Token 信息到 Redis
func saveTokenInfo(userID, deviceID, token, ipaddr, deviceInfo string) error {
	ctx := context.Background()
	key := UserSessionCachePrefixToString(userID, deviceID, ipaddr)
	lastActive := time.Now().Format(time.RFC3339)

	if err := cache.HMSet(ctx, key, gin.H{
		FieldDeviceInfo:          deviceInfo,
		FieldRefreshToken:        token,
		JwtCustomClaimsIpaddress: ipaddr,
		FieldLastActive:          lastActive,
	}); err != nil {
		return err
	}
	return cache.Expire(ctx, key, RefreshTokenExpire)
}

// 解析 Token
func parseToken(tokenString string) (*CustomClaims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, true
	}
	return nil, false
}

// 自动刷新 Token 的逻辑
func handleTokenRefresh(c *gin.Context, oldAccessToken, ipaddr string) bool {
	refreshToken := c.GetHeader(HeaderRefreshToken)
	deviceID := c.GetHeader(HeaderDeviceID)
	if refreshToken == "" || deviceID == "" {
		return false
	}

	refreshClaims, valid := parseToken(refreshToken)
	if !valid || !isDeviceValid(refreshClaims.UserID, deviceID, ipaddr) {
		return false
	}

	storedToken, err := GetStoredRefreshToken(refreshClaims.UserID, deviceID, ipaddr)
	if err != nil || storedToken != refreshToken {
		return false
	}

	newAccessToken, newRefreshToken, err := GenerateTokenPair(refreshClaims.UserID, deviceID, ipaddr)
	if err != nil {
		return false
	}

	if err := updateRefreshToken(refreshClaims.UserID, deviceID, newRefreshToken, ipaddr); err != nil {
		return false
	}

	c.Header(HeaderNewAccessToken, newAccessToken)
	c.Header(HeaderNewRefreshToken, newRefreshToken)
	c.Set(ContextKeyUserID, refreshClaims.UserID)
	c.Set(ContextKeyDeviceID, deviceID)
	c.Next()
	return true
}

// 从 Redis 获取 RefreshToken
func GetStoredRefreshToken(userID, deviceID, ipaddr string) (string, error) {
	ctx := context.Background()
	key := UserSessionCachePrefixToString(userID, deviceID, ipaddr)
	return cache.HGet(ctx, key, FieldRefreshToken)
}

// 更新 Redis 中的 RefreshToken
func updateRefreshToken(userID, deviceID, newToken, ipaddr string) error {
	ctx := context.Background()
	key := UserSessionCachePrefixToString(userID, deviceID, ipaddr)
	err := cache.HSet(ctx, key,
		FieldRefreshToken, newToken,
		FieldLastActive, time.Now().Format(time.RFC3339),
	)
	if err != nil {
		return err
	}
	return cache.Expire(ctx, key, RefreshTokenExpire)
}

// 校验 Redis 中设备 session 是否仍然存在（设备合法性）
func isDeviceValid(userID, deviceID, ipaddr string) bool {
	ctx := context.Background()
	key := UserSessionCachePrefixToString(userID, deviceID, ipaddr)
	exists, err := cache.Exists(ctx, key)
	return err == nil && exists
}
