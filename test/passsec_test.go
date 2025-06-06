package qwqtest

import (
	"fmt"
	"log"
	"qwqserver/pkg/util/passsec"
	"testing"
)

func TestPasssecRun1(t *testing.T) {
	t.Run("test passsec run", func(t *testing.T) {
		
	})
}

func TestPasssecRun2(t *testing.T) {
	t.Run("test passsec run", func(t *testing.T) {
		// 示例1: 哈希密码
		password := "MySecureP@ssw0rd!"

		// 哈希密码（使用默认成本）
		hashed, err := passsec.Hash(password)
		if err != nil {
			log.Fatalf("哈希失败: %v", err)
		}
		fmt.Printf("哈希密码 (默认成本): %s\n", hashed)

		// 哈希密码（指定成本14）
		hashedHighCost, err := passsec.Hash(password, 14)
		if err != nil {
			log.Fatalf("哈希失败: %v", err)
		}
		fmt.Printf("哈希密码 (高成本): %s\n", hashedHighCost)

		// 示例2: 验证密码
		match, err := passsec.Check(password, hashed)
		if err != nil {
			log.Fatalf("验证失败: %v", err)
		}
		fmt.Printf("密码匹配: %v\n", match)

		// 示例3: 检查错误密码
		wrongPassword := "wrongpassword"
		match, err = passsec.Check(wrongPassword, hashed)
		if err != nil {
			log.Fatalf("验证失败: %v", err)
		}
		fmt.Printf("错误密码匹配: %v\n", match)

		// 示例4: 检查密码强度
		strength := passsec.CheckStrength(password)
		fmt.Printf("密码强度: %s\n", strength)

		// 示例5: 生成随机盐
		salt, err := passsec.GenerateRandomSalt(16)
		if err != nil {
			log.Fatalf("生成盐失败: %v", err)
		}
		fmt.Printf("随机盐: %s\n", salt)

		// 示例6: 从哈希中提取成本
		cost, err := passsec.GetCost(hashed)
		if err != nil {
			log.Fatalf("提取成本失败: %v", err)
		}
		fmt.Printf("哈希成本因子: %d\n", cost)

		// 示例7: 测试不同密码强度
		testPasswords := []string{
			"12345",                             // 非常弱
			"password",                          // 弱
			"Password1",                         // 中等
			"P@ssw0rd",                          // 强
			"VeryL0ngP@ssw0rd!WithSpecialChars", // 非常强
		}

		fmt.Println("\n密码强度测试:")
		for _, pwd := range testPasswords {
			strength := passsec.CheckStrength(pwd)
			fmt.Printf("- '%s': %s\n", pwd, strength)
		}
	})
}
