package qwqtest

import (
	"fmt"
	"qwqserver/pkg/util/multih256"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		// 生成盐值
		salt, err := multih256.GenerateSalt(multih256.DefaultSaltLength)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Generated salt: %s\n", salt)

		// 加密数据
		data := "sensitive_information@123"
		encrypted, err := multih256.Encrypt(data, salt, 5)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Encrypted data: %s\n", encrypted)

		// 验证数据
		valid, err := multih256.Verify(data, encrypted, salt, 5)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Verification result: %v\n", valid) // true

		// 一步完成盐值生成和加密
		encrypted, newSalt, err := multih256.EncryptWithSalt(data, 32, 4)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Encrypted with salt: %s\nSalt: %s\n", encrypted, newSalt)
	})
}
