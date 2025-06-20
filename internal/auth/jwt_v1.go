package auth

//
//import (
//	"errors"
//	"qwqserver/internal/common/constant"
//	"time"
//
//	"github.com/golang-jwt/jwt/v5"
//)
//
//type Claims struct {
//	UserID   uint   `json:"uid"`
//	DeviceID string `json:"did"`
//	jwt.RegisteredClaims
//}
//
//// GenerateToken 生成JWT Token
//func GenerateToken(userID uint, deviceID string) (string, string, error) {
//	now := time.Now()
//	expireTime := now.Add(constant.TokenExpire * time.Second)
//	refreshExpireTime := now.Add(constant.RefreshTokenExpire * time.Second)
//
//	// 生成Access Token
//	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
//		UserID:   userID,
//		DeviceID: deviceID,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(expireTime),
//			IssuedAt:  jwt.NewNumericDate(now),
//			Subject:   "access_token",
//		},
//	})
//	token, err := tokenClaims.SignedString([]byte(constant.JWTSecretKey))
//	if err != nil {
//		return "", "", err
//	}
//
//	// 生成Refresh Token
//	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
//		ExpiresAt: jwt.NewNumericDate(refreshExpireTime),
//		IssuedAt:  jwt.NewNumericDate(now),
//		Subject:   "refresh_token",
//		ID:        deviceID, // 将设备ID存储在jti中
//	})
//	refreshToken, err := refreshClaims.SignedString([]byte(constant.JWTSecretKey))
//	if err != nil {
//		return "", "", err
//	}
//
//	return token, refreshToken, nil
//}
//
//// ParseToken 解析Token
//func ParseToken(tokenString string) (*Claims, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//		return []byte(constant.JWTSecretKey), nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//		return claims, nil
//	}
//	return nil, errors.New("invalid token")
//}
//
//// ParseRefreshToken 解析Refresh Token
//func ParseRefreshToken(refreshToken string) (string, error) {
//	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
//		return []byte(constant.JWTSecretKey), nil
//	})
//	if err != nil {
//		return "", err
//	}
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		// 从jti获取设备ID
//		if deviceID, ok := claims["jti"].(string); ok {
//			return deviceID, nil
//		}
//	}
//	return "", errors.New("invalid refresh token")
//}
