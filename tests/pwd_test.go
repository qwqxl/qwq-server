package qwqtest

import (
	"fmt"
	"qwqserver/pkg/security"
	"testing"
)

func TestPwdMain(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		pwdStr := "jsJiiu223892_"

		if err := security.PwdCheckPasswordStrength(pwdStr); err != nil {
			t.Error(err)
			return
		}

		/*
			Time = 1       // 1次迭代（实际计算量由内存控制）
			Memory = 64MB  // 推荐64MB-128MB（平衡安全性与性能）
			Threads = 4    // 4线程并行
			KeyLen = 32    // 256位密钥输出
			SaltLen = 16   // 128位随机盐
		*/

		// 创建企业级参数
		pwdParams := &security.PwdArgon2Params{
			Time:    3,
			Memory:  128 * 1024,
			Threads: 6,
			KeyLen:  32,
			SaltLen: 16,
		}

		//pwdParams := security.NewPwdDefaultParams()

		hashed, err := security.HashPassword(pwdStr, pwdParams)
		if err != nil {
			fmt.Println("密码哈希生成失败：", err)
			return
		}

		fmt.Println("密码哈希：", hashed)

		// 密码验证
		err = security.VerifyPassword("Secur3P@ss!", hashed)
		if err == nil {
			fmt.Println("Password verified!")
		}
	})
}
