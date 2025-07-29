// passsec 提供密码安全处理功能，包括哈希、验证和强度评估
package passsec

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// 默认工作因子 - 平衡安全性和性能
const DefaultCost = 12

// PasswordStrength 密码强度级别
type PasswordStrength int

const (
	VeryWeak PasswordStrength = iota
	Weak
	Moderate
	Strong
	VeryStrong
)

// StrengthDescription 返回密码强度的文本描述
func (s PasswordStrength) String() string {
	switch s {
	case VeryWeak:
		return "非常弱"
	case Weak:
		return "弱"
	case Moderate:
		return "中等"
	case Strong:
		return "强"
	case VeryStrong:
		return "非常强"
	default:
		return "未知"
	}
}

// Hash 生成密码的bcrypt哈希
// 参数:
//
//	password: 明文密码
//	cost: bcrypt工作因子(4-31)，默认使用DefaultCost
//
// 返回:
//
//	哈希字符串, 错误
//
// 实现细节:
//  1. 验证工作因子是否在有效范围(4-31)
//  2. 使用bcrypt.GenerateFromPassword生成哈希
//  3. bcrypt自动处理盐值生成和存储
//  4. 哈希结果包含算法标识、工作因子和盐值
func Hash(password string, cost ...int) (string, error) {
	// 确定工作因子
	bcryptCost := DefaultCost
	if len(cost) > 0 {
		bcryptCost = cost[0]
	}

	// 验证工作因子范围
	if bcryptCost < bcrypt.MinCost || bcryptCost > bcrypt.MaxCost {
		return "", fmt.Errorf("无效的工作因子: %d (必须在%d-%d之间)",
			bcryptCost, bcrypt.MinCost, bcrypt.MaxCost)
	}

	// 将密码转换为字节切片
	bytes := []byte(password)

	// 生成带盐的bcrypt哈希
	hashedBytes, err := bcrypt.GenerateFromPassword(bytes, bcryptCost)
	if err != nil {
		return "", fmt.Errorf("生成密码哈希失败: %w", err)
	}

	return string(hashedBytes), nil
}

// Check 验证密码是否匹配哈希
// 参数:
//
//	password: 待验证的明文密码
//	hashedPassword: 之前存储的哈希密码
//
// 返回:
//
//	bool: 是否匹配, 错误
//
// 实现细节:
//  1. 使用bcrypt.CompareHashAndPassword函数
//  2. 函数安全地比较哈希值，防止时序攻击
//  3. 自动从哈希中提取盐值和工作因子
//  4. 错误处理包括密码不匹配和无效哈希格式
func Check(password, hashedPassword string) (bool, error) {
	// 将输入转换为字节切片
	hashedBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)

	// 比较密码和哈希
	err := bcrypt.CompareHashAndPassword(hashedBytes, passwordBytes)

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// 密码不匹配
			return false, nil
		}
		// 其他错误（如哈希格式无效）
		return false, fmt.Errorf("密码验证失败: %w", err)
	}

	// 密码匹配
	return true, nil
}

// GenerateRandomSalt 生成随机盐值
// 参数:
//
//	length: 盐值长度（字节）
//
// 返回:
//
//	base64编码的盐字符串, 错误
//
// 实现细节:
//  1. 使用crypto/rand生成加密安全的随机字节
//  2. 使用base64.URLEncoding进行编码
//  3. 确保生成的盐值适合存储
func GenerateRandomSalt(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("盐值长度必须大于0")
	}

	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("生成随机盐失败: %w", err)
	}

	return base64.URLEncoding.EncodeToString(salt), nil
}

// CheckStrength 评估密码强度
// 参数:
//
//	password: 要评估的密码
//
// 返回:
//
//	PasswordStrength: 密码强度级别
//
// 实现细节:
//  1. 检查长度要求
//  2. 检查字符种类（小写、大写、数字、特殊字符）
//  3. 使用正则表达式识别特殊字符
//  4. 根据满足的条件数量和长度确定强度
func CheckStrength(password string) PasswordStrength {
	var (
		hasLower   = false
		hasUpper   = false
		hasDigit   = false
		hasSpecial = false
		length     = len(password)
	)

	// 特殊字符正则表达式
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*()\-_=+{}[\]:;'"<>,.?/|\\]`)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case specialCharRegex.MatchString(string(char)):
			hasSpecial = true
		}
	}

	// 计算满足的条件数量
	conditionsMet := 0
	if hasLower {
		conditionsMet++
	}
	if hasUpper {
		conditionsMet++
	}
	if hasDigit {
		conditionsMet++
	}
	if hasSpecial {
		conditionsMet++
	}

	// 根据长度和条件数量确定强度
	switch {
	case length < 6:
		return VeryWeak
	case length < 8:
		return Weak
	case length < 10 && conditionsMet < 3:
		return Moderate
	case length < 12 && conditionsMet < 4:
		return Strong
	case length >= 12 && conditionsMet == 4:
		return VeryStrong
	case length >= 10:
		return Strong
	default:
		return Moderate
	}
}

// GetCost 从哈希中提取工作因子
// 参数:
//
//	hashedPassword: bcrypt哈希字符串
//
// 返回:
//
//	工作因子值, 错误
//
// 实现细节:
//  1. 解析bcrypt哈希格式($2a$cost$...)
//  2. 提取成本参数并转换为整数
func GetCost(hashedPassword string) (int, error) {
	// bcrypt哈希格式: $2a$cost$salt+hash
	if len(hashedPassword) < 7 || hashedPassword[0] != '$' {
		return 0, errors.New("无效的bcrypt哈希格式")
	}

	// 查找第二个$符号的位置
	dollarPos := 4
	for i := 4; i < len(hashedPassword); i++ {
		if hashedPassword[i] == '$' {
			dollarPos = i
			break
		}
	}

	if dollarPos == 4 || dollarPos > len(hashedPassword)-1 {
		return 0, errors.New("无效的bcrypt哈希格式")
	}

	// 提取成本部分(例如: $2a$10$...)
	costStr := hashedPassword[4:dollarPos]
	cost := 0
	for _, c := range costStr {
		if c < '0' || c > '9' {
			return 0, errors.New("无效的成本格式")
		}
		cost = cost*10 + int(c-'0')
	}

	return cost, nil
}
