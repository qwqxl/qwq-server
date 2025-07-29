package auth

import (
	"qwqserver/internal/common"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义Claims
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Platform string `json:"platform"`
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}

// 生成Token
func GenerateToken(userID, platform, deviceID string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Platform: platform,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(common.TokenExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(common.JWTSecretKey))
}

// 生成RefreshToken
func GenerateRefreshToken(userID, platform, deviceID string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Platform: platform,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(common.RefreshTokenExpire)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(common.JWTSecretKey))
}

// 解析Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(common.JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
