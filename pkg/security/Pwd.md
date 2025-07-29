# Password Hashing

```
Time = 1       // 1次迭代（实际计算量由内存控制）
Memory = 64MB  // 推荐64MB-128MB（平衡安全性与性能）
Threads = 4    // 4线程并行
KeyLen = 32    // 256位密钥输出
SaltLen = 16   // 128位随机盐
```


```go

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

const (
	DefaultIterations  = 1          // 迭代次数（建议 >=1）
	DefaultMemory      = 64 * 1024  // 内存消耗：64MB，单位KB
	DefaultParallelism = 4          // 并行度（线程数）
	DefaultKeyLength   = 32         // 输出的密钥长度
	DefaultSaltLength  = 16         // 盐的长度
)

// ------------------------- 错误信息定义 -------------------------

var (
	ErrHash格式无效        = errors.New("哈希格式无效")
	ErrArgon2版本不兼容    = errors.New("argon2 版本不兼容")
	Err密码不匹配          = errors.New("密码不正确")
)

// ------------------------- 参数结构体定义 -------------------------

// Argon2参数结构体，支持定制
type Argon2参数 struct {
	迭代次数  uint32
	内存消耗  uint32
	线程数    uint8
	密钥长度  uint32
	盐长度    uint32
}

// 返回默认安全配置参数
func 获取默认参数() *Argon2参数 {
	return &Argon2参数{
		迭代次数:  DefaultIterations,
		内存消耗:  DefaultMemory,
		线程数:    DefaultParallelism,
		密钥长度:  DefaultKeyLength,
		盐长度:    DefaultSaltLength,
	}
}

// ------------------------- 密码加密与校验逻辑 -------------------------

// Hash密码：生成带盐的 Argon2id 哈希字符串，包含所有参数（PHC格式）
func 生成密码哈希(明文密码 string, 参数 *Argon2参数) (string, error) {
	// 生成随机盐
	盐 := make([]byte, 参数.盐长度)
	if _, err := rand.Read(盐); err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}

	// 使用 Argon2id 生成哈希
	哈希 := argon2.IDKey(
		[]byte(明文密码),
		盐,
		参数.迭代次数,
		参数.内存消耗,
		参数.线程数,
		参数.密钥长度,
	)

	// Base64 编码盐与哈希
	盐编码 := base64.RawStdEncoding.EncodeToString(盐)
	哈希编码 := base64.RawStdEncoding.EncodeToString(哈希)

	// 按 PHC 标准格式化最终存储的哈希字符串
	返回字符串 := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		参数.内存消耗,
		参数.迭代次数,
		参数.线程数,
		盐编码,
		哈希编码,
	)

	return 返回字符串, nil
}

// 校验密码：判断输入的明文密码是否匹配存储的哈希
func 校验密码(明文密码, 哈希字符串 string) error {
	参数, 盐, 原始哈希, err := 解析哈希字符串(哈希字符串)
	if err != nil {
		return err
	}

	// 使用相同参数重新生成哈希
	新哈希 := argon2.IDKey(
		[]byte(明文密码),
		盐,
		参数.迭代次数,
		参数.内存消耗,
		参数.线程数,
		参数.密钥长度,
	)

	// 使用常量时间比较（防时序攻击）
	if subtle.ConstantTimeCompare(原始哈希, 新哈希) != 1 {
		return Err密码不匹配
	}

	return nil
}

// ------------------------- 哈希字符串解析 -------------------------

// 从哈希字符串中解析出参数、盐、哈希值
func 解析哈希字符串(哈希字符串 string) (*Argon2参数, []byte, []byte, error) {
	parts := strings.Split(哈希字符串, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrHash格式无效
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrHash格式无效
	}

	// 解析版本
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, fmt.Errorf("版本解析失败: %w", err)
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrArgon2版本不兼容
	}

	// 解析算法参数
	参数 := &Argon2参数{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &参数.内存消耗, &参数.迭代次数, &参数.线程数); err != nil {
		return nil, nil, nil, fmt.Errorf("参数解析失败: %w", err)
	}

	// 解码盐值
	盐, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("盐解码失败: %w", err)
	}
	参数.盐长度 = uint32(len(盐))

	// 解码哈希值
	哈希, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("哈希解码失败: %w", err)
	}
	参数.密钥长度 = uint32(len(哈希))

	return 参数, 盐, 哈希, nil
}

// ------------------------- 密码强度检查 -------------------------

// 检查密码强度：长度 + 至少包含三类字符（大写、小写、数字、特殊符号）
func 检查密码强度(密码 string) error {
	if len(密码) < 12 {
		return errors.New("密码长度至少为12个字符")
	}

	匹配数量 := 0
	规则 := []*regexp.Regexp{
		regexp.MustCompile(`[A-Z]`),        // 至少一个大写字母
		regexp.MustCompile(`[a-z]`),        // 至少一个小写字母
		regexp.MustCompile(`[0-9]`),        // 至少一个数字
		regexp.MustCompile(`[^a-zA-Z0-9]`), // 至少一个特殊字符
	}

	for _, 正则 := range 规则 {
		if 正则.MatchString(密码) {
			匹配数量++
		}
	}

	if 匹配数量 < 3 {
		return errors.New("密码必须包含以下任意三项：大写字母、小写字母、数字、特殊字符")
	}

	return nil
}

```