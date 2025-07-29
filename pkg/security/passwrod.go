// Package security 提供密码加密与校验工具，使用 Argon2id 算法，符合 OWASP 推荐标准
package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ------------------------- Argon2 默认参数配置（安全级别高） -------------------------

// 默认加密参数（符合OWASP推荐标准）
const (
	PwdDefaultTime    = 1         // 迭代次数（建议 >=1）
	PwdDefaultMemory  = 64 * 1024 // 内存消耗：64MB，单位KB
	PwdDefaultThreads = 4         // 并行度（线程数）
	PwdDefaultKeyLen  = 32        // 输出的密钥长度
	PwdDefaultSaltLen = 16        // 盐的长度
)

// ------------------------- 错误信息定义 -------------------------

var (
	PwdErrInvalidHash         = errors.New("invalid hash format")
	PwdErrIncompatibleVersion = errors.New("incompatible argon2 version")
	PwdErrPasswordMismatch    = errors.New("password does not match hash")
)

// ------------------------- 参数结构体定义 -------------------------

// PwdArgon2Params 存储Argon2id算法参数
type PwdArgon2Params struct {
	Time    uint32 // 迭代次数
	Memory  uint32 // 内存消耗(KB)
	Threads uint8  // 并行线程数
	KeyLen  uint32 // 输出密钥长度
	SaltLen uint32 // 盐值长度
}

// NewPwdDefaultParams 返回安全的默认参数配置
func NewPwdDefaultParams() *PwdArgon2Params {
	return &PwdArgon2Params{
		Time:    PwdDefaultTime,
		Memory:  PwdDefaultMemory,
		Threads: PwdDefaultThreads,
		KeyLen:  PwdDefaultKeyLen,
		SaltLen: PwdDefaultSaltLen,
	}
}

// HashPassword 生成带盐的Argon2id密码哈希 哈希字符串，包含所有参数（PHC格式）
func HashPassword(password string, params *PwdArgon2Params) (string, error) {
	// 生成随机盐
	// 盐 := make([]byte, 参数.盐长度)
	salt := make([]byte, params.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}

	// 生成Argon2id哈希
	hash := argon2.IDKey(
		[]byte(password), // 明文密码
		salt,             // 盐
		params.Time,      // 参数.迭代次数
		params.Memory,    // 参数.内存消耗
		params.Threads,   // 参数.线程数
		params.KeyLen,    // 参数.密钥长度
	)

	// Base64编码
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式化存储字符串 (遵循PHC字符串格式)
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Time,
		params.Threads,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

// VerifyPassword 验证密码与哈希是否匹配
func VerifyPassword(password, encodedHash string) error {
	// 解析哈希字符串
	params, salt, hash, err := pwdDecodeHash(encodedHash)
	if err != nil {
		return err
	}

	// 使用相同参数生成新哈希
	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLen,
	)

	// 安全比较哈希值（防止时序攻击）
	if subtle.ConstantTimeCompare(hash, newHash) != 1 {
		return PwdErrPasswordMismatch
	}

	return nil
}

// 解析存储的哈希字符串
func pwdDecodeHash(encodedHash string) (*PwdArgon2Params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, PwdErrInvalidHash
	}

	// 验证算法标识
	if parts[1] != "argon2id" {
		return nil, nil, nil, PwdErrInvalidHash
	}

	// 解析版本
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, PwdErrIncompatibleVersion
	}

	// 解析参数
	params := &PwdArgon2Params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Time, &params.Threads)
	if err != nil {
		return nil, nil, nil, err
	}

	// 解码盐值
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLen = uint32(len(salt))

	// 解码哈希值
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLen = uint32(len(hash))

	return params, salt, hash, nil
}

// PwdCheckPasswordStrength 企业级密码强度校验
func PwdCheckPasswordStrength(password string) error {
	if len(password) < 12 {
		return errors.New("密码长度至少为12个字符")
		//return errors.New("password must be at least 12 characters")
	}

	var checks int
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`[A-Z]`),        // 大写字母
		regexp.MustCompile(`[a-z]`),        // 小写字母
		regexp.MustCompile(`[0-9]`),        // 数字
		regexp.MustCompile(`[^a-zA-Z0-9]`), // 特殊字符
	}

	for _, pattern := range patterns {
		if pattern.MatchString(password) {
			checks++
		}
	}

	if checks < 3 {
		return errors.New("password must include 3 of: uppercase, lowercase, number, special character")
	}

	return nil
}
