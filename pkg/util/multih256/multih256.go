package multih256

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"regexp"
)

const (
	// DefaultRounds 默认哈希轮次
	DefaultRounds = 3
	// DefaultSaltLength 默认盐值长度（字节）
	DefaultSaltLength = 16
)

var (
	// ErrInvalidRounds 轮次无效错误
	ErrInvalidRounds = errors.New("rounds must be at least 1")
	// ErrInvalidSaltLength 盐值长度无效错误
	ErrInvalidSaltLength = errors.New("salt length must be at least 1 byte")
	// ErrInvalidHexFormat 十六进制格式无效错误
	ErrInvalidHexFormat = errors.New("invalid hexadecimal format")
)

// GenerateSalt 生成密码学安全的随机盐值
//
// 参数:
//
//	length - 盐值长度（字节），建议至少16字节
//
// 返回:
//
//	十六进制格式的盐值字符串
//	错误信息（如果有）
func GenerateSalt(length int) (string, error) {
	if length < 1 {
		return "", ErrInvalidSaltLength
	}

	// 创建足够长度的字节切片
	saltBytes := make([]byte, length)

	// 使用 crypto/rand 填充随机字节
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}

	// 转换为十六进制字符串
	return hex.EncodeToString(saltBytes), nil
}

// validateHexString 验证十六进制字符串格式
func validateHexString(value string) error {
	// 检查是否为有效的十六进制字符串
	matched, err := regexp.MatchString(`^[0-9a-fA-F]+$`, value)
	if err != nil {
		return err
	}

	if !matched {
		return ErrInvalidHexFormat
	}

	// 检查长度是否为偶数（每个字节对应两个十六进制字符）
	if len(value)%2 != 0 {
		return ErrInvalidHexFormat
	}

	return nil
}

// Encrypt 执行多重SHA-256加密
//
// 参数:
//
//	data   - 要加密的原始字符串数据
//	salt   - 十六进制格式的盐值（可选）
//	rounds - 哈希轮次（至少1轮）
//
// 返回:
//
//	十六进制字符串格式的加密结果
//	错误信息（如果有）
func Encrypt(data string, salt string, rounds int) (string, error) {
	if rounds < 1 {
		return "", ErrInvalidRounds
	}

	// 处理盐值
	var saltBytes []byte
	if salt != "" {
		if err := validateHexString(salt); err != nil {
			return "", err
		}

		var err error
		saltBytes, err = hex.DecodeString(salt)
		if err != nil {
			return "", err
		}
	}

	// 准备初始数据：盐值 + 原始数据
	dataBytes := []byte(data)
	currentHash := append(saltBytes, dataBytes...)

	// 执行多轮哈希
	for i := 0; i < rounds; i++ {
		hash := sha256.Sum256(currentHash)
		currentHash = hash[:] // 将数组转换为切片
	}

	// 返回十六进制字符串结果
	return hex.EncodeToString(currentHash), nil
}

// Verify 验证数据是否与加密结果匹配
//
// 参数:
//
//	data      - 要验证的原始字符串数据
//	encrypted - 先前加密的结果（十六进制字符串）
//	salt      - 加密时使用的盐值
//	rounds    - 加密时使用的轮次
//
// 返回:
//
//	bool - 如果数据匹配加密结果则为true，否则为false
//	error - 错误信息（如果有）
func Verify(data string, encrypted string, salt string, rounds int) (bool, error) {
	// 验证加密结果格式
	if err := validateHexString(encrypted); err != nil {
		return false, err
	}

	// 对新数据加密
	newEncrypted, err := Encrypt(data, salt, rounds)
	if err != nil {
		return false, err
	}

	// 安全比较（防御时序攻击）
	return subtle.ConstantTimeCompare([]byte(newEncrypted), []byte(encrypted)) == 1, nil
}

// EncryptWithSalt 生成盐值并加密数据（一步完成）
//
// 参数:
//
//	data       - 要加密的原始字符串数据
//	saltLength - 要生成的盐值长度（字节）
//	rounds     - 哈希轮次（至少1轮）
//
// 返回:
//
//	加密结果（十六进制字符串）
//	盐值（十六进制字符串）
//	错误信息（如果有）
func EncryptWithSalt(data string, saltLength int, rounds int) (string, string, error) {
	salt, err := GenerateSalt(saltLength)
	if err != nil {
		return "", "", err
	}

	encrypted, err := Encrypt(data, salt, rounds)
	if err != nil {
		return "", "", err
	}

	return encrypted, salt, nil
}
