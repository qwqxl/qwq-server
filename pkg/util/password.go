package util

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"unicode"
)

func HashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(pwd, hash string) bool {
	
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}

// ValidatePasswordComplexity 检查密码是否符合复杂度要求
func ValidatePasswordComplexity(password string) error {
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("密码必须包含大写字母、小写字母、数字和特殊字符")
	}
	return nil
}

// ValidateEmailSuffix 检查邮箱是否属于允许的后缀
func ValidateEmailSuffix(email string, allowedSuffixes []string) error {
	email = strings.ToLower(email)
	for _, suffix := range allowedSuffixes {
		if strings.HasSuffix(email, strings.ToLower(suffix)) {
			return nil
		}
	}
	return fmt.Errorf("邮箱后缀不被允许，必须为以下之一: %v", allowedSuffixes)
}
