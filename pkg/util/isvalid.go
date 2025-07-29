package util

import (
	"regexp"
	"strings"
)

// IsValidEmailOld 简单的邮箱格式验证
func IsValidEmailOld(email string) bool {
	// 这里可以使用更复杂的正则表达式验证
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// IsValidEmailFormat 使用正则表达式验证邮箱地址格式是否有效
// 例如有效邮箱: user@example.com、user.name+tag@sub.domain.co.uk
func IsValidEmailFormat(email string) bool {
	// emailFormatRegex 是用于验证邮箱格式的正则表达式：
	// - 允许字母、数字、点、下划线、加号、短横线作为用户名
	// - 域名部分允许子域名，TLD 最少两位字母
	emailFormatRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailFormatRegex.MatchString(email)
}

// IsAllowedEmailSuffix 验证邮箱是否属于指定域名后缀
func IsAllowedEmailSuffix(email string, allowedSuffixes []string) bool {
	for _, suffix := range allowedSuffixes {
		if strings.HasSuffix(email, suffix) {
			return true
		}
	}
	return false
}
