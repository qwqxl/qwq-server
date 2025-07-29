package authx

import "github.com/golang-jwt/jwt/v5"

// CustomClaims 自定义 Claims 结构
type CustomClaims struct {
	UserID       string `json:"uid"`
	DeviceSign   string `json:"did"`   // 设备标识：pc/web、mobile/app
	PlatformSign string `json:"plat"`  // 平台标识：如 qwq、tv、draw
	jwt.RegisteredClaims
}