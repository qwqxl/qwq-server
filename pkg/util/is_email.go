package util

import "regexp"

// IsEmail 判断字符串是否是合法的邮箱格式
func IsEmail(email string) bool {
	// 这是一个通用的 Email 正则表达式
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
