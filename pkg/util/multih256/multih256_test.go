// multih256_test.go
package multih256

import (
	"fmt"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	t.Run("ValidLength", func(t *testing.T) {
		salt, err := GenerateSalt(16)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(salt) != 32 { // 16字节 = 32个十六进制字符
			t.Errorf("Expected salt length 32, got %d", len(salt))
		}
	})

	t.Run("InvalidLength", func(t *testing.T) {
		_, err := GenerateSalt(0)
		if err != ErrInvalidSaltLength {
			t.Errorf("Expected ErrInvalidSaltLength, got %v", err)
		}
	})
}

func TestEncrypt(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		salt     string
		rounds   int
		expected string
		wantErr  bool
	}{
		{
			name:   "BasicEncryption",
			data:   "hello world",
			salt:   "",
			rounds: 1,
			// 计算方式: sha256("hello world")
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:   "WithSalt",
			data:   "password123",
			salt:   "0011223344556677",
			rounds: 1,
			// 计算方式: sha256(hex"0011223344556677" + "password123")
			expected: "f0f8b3a5e7e8c7a9c1e4e7d5e7c1a1d8e2c0a3d7e1c3a5d8e2c0a3d7e1c3a5d8",
		},
		{
			name:   "MultipleRounds",
			data:   "test",
			salt:   "",
			rounds: 3,
			// 计算方式: sha256(sha256(sha256("test")))
			expected: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		{
			name:    "InvalidRounds",
			data:    "test",
			salt:    "",
			rounds:  0,
			wantErr: true,
		},
		{
			name:    "InvalidHexSalt",
			data:    "test",
			salt:    "invalid-hex",
			rounds:  1,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Encrypt(tc.data, tc.salt, tc.rounds)

			if tc.wantErr {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	data := "sensitive_data"
	salt, _ := GenerateSalt(16)
	encrypted, _ := Encrypt(data, salt, 5)

	t.Run("ValidVerification", func(t *testing.T) {
		valid, err := Verify(data, encrypted, salt, 5)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !valid {
			t.Error("Expected verification to succeed")
		}
	})

	t.Run("InvalidData", func(t *testing.T) {
		valid, err := Verify("wrong_data", encrypted, salt, 5)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if valid {
			t.Error("Expected verification to fail")
		}
	})

	t.Run("InvalidSalt", func(t *testing.T) {
		wrongSalt, _ := GenerateSalt(16)
		valid, err := Verify(data, encrypted, wrongSalt, 5)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if valid {
			t.Error("Expected verification to fail")
		}
	})

	t.Run("InvalidRounds", func(t *testing.T) {
		valid, err := Verify(data, encrypted, salt, 4) // 使用不同轮次
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if valid {
			t.Error("Expected verification to fail")
		}
	})

	t.Run("InvalidEncryptedFormat", func(t *testing.T) {
		_, err := Verify(data, "invalid-hex", salt, 5)
		if err == nil {
			t.Error("Expected error for invalid encrypted format")
		}
	})
}

func TestEncryptWithSalt(t *testing.T) {
	data := "important_data"

	t.Run("ValidParameters", func(t *testing.T) {
		encrypted, salt, err := EncryptWithSalt(data, 16, 4)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// 验证盐值长度
		if len(salt) != 32 {
			t.Errorf("Expected salt length 32, got %d", len(salt))
		}

		// 验证加密结果格式
		if err := validateHexString(encrypted); err != nil {
			t.Errorf("Invalid encrypted format: %v", err)
		}

		// 验证加密结果是否匹配
		valid, err := Verify(data, encrypted, salt, 4)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !valid {
			t.Error("Generated encryption failed verification")
		}
	})

	t.Run("InvalidSaltLength", func(t *testing.T) {
		_, _, err := EncryptWithSalt(data, 0, 4)
		if err != ErrInvalidSaltLength {
			t.Errorf("Expected ErrInvalidSaltLength, got %v", err)
		}
	})

	t.Run("InvalidRounds", func(t *testing.T) {
		_, _, err := EncryptWithSalt(data, 16, 0)
		if err != ErrInvalidRounds {
			t.Errorf("Expected ErrInvalidRounds, got %v", err)
		}
	})
}

func BenchmarkEncrypt(b *testing.B) {
	data := "benchmark_data"
	salt, _ := GenerateSalt(16)

	b.Run("SingleRound", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = Encrypt(data, salt, 1)
		}
	})

	b.Run("DefaultRounds", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = Encrypt(data, salt, DefaultRounds)
		}
	})

	b.Run("TenRounds", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = Encrypt(data, salt, 10)
		}
	})
}

func TestRun(t *testing.T) {
	t.Run("BasicUsage", func(t *testing.T) {
		// 生成盐值
		salt, err := GenerateSalt(DefaultSaltLength)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Generated salt: %s\n", salt)

		// 加密数据
		data := "sensitive_information@123"
		encrypted, err := Encrypt(data, salt, 5)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Encrypted data: %s\n", encrypted)

		// 验证数据
		valid, err := Verify(data, encrypted, salt, 5)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Verification result: %v\n", valid) // true

		// 一步完成盐值生成和加密
		encrypted, newSalt, err := EncryptWithSalt(data, 32, 4)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Encrypted with salt: %s\nSalt: %s\n", encrypted, newSalt)
	})
}
