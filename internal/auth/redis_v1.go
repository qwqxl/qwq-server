package auth

//
//import (
//	"context"
//	"fmt"
//	"qwqserver/pkg/cache"
//	"strings"
//)
//
//func StoreTokens(userID, deviceID, accessToken, refreshToken string) error {
//	key := fmt.Sprintf("user:%s:device:%s", userID, deviceID)
//
//	// 存储Token并设置过期时间（与Refresh Token一致）
//	return cache.Set(context.Background(), key, fmt.Sprintf("%s:%s", accessToken, refreshToken), RefreshTokenExpire)
//	//return RedisClient.Set(context.Background(), key,
//	//	fmt.Sprintf("%s:%s", accessToken, refreshToken),
//	//	RefreshTokenExpire).Err()
//}
//
//func GetTokens(userID, deviceID string) (string, string, error) {
//	key := fmt.Sprintf("user:%s:device:%s", userID, deviceID)
//	//val, err := RedisClient.Get(context.Background(), key).Result()
//	val, err := cache.Get(context.Background(), key)
//	if err != nil {
//		return "", "", err
//	}
//
//	// 格式: accessToken:refreshToken
//	tokens := strings.Split(val, ":")
//	if len(tokens) != 2 {
//		return "", "", fmt.Errorf("invalid token format")
//	}
//	return tokens[0], tokens[1], nil
//}
//
//func DeleteTokens(userID, deviceID string) error {
//	key := fmt.Sprintf("user:%s:device:%s", userID, deviceID)
//	return cache.Delete(context.Background(), key)
//	//return RedisClient.Del(context.Background(), key).Err()
//}
