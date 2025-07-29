package auth

import "fmt"

// BuildTokenKey 构建token key
func BuildTokenKey(prefix, uid, platform, deviceID string) string {
	return fmt.Sprintf("%s%s:%s:%s", prefix, uid, platform, deviceID)
}
