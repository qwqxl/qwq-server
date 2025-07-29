package captcha

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"qwqserver/pkg/cache"
	"time"
)

const (
	DefaultExpiration = 5 * time.Minute
	MaxAttempts       = 5
)

type DeviceType string

const (
	DevicePC     DeviceType = "pc"
	DeviceMobile DeviceType = "mobile"
	DeviceApp    DeviceType = "app"
)

type CaptchaService struct {
	cache *cache.Client
}

func NewCaptchaService(cacheClient *cache.Client) *CaptchaService {
	return &CaptchaService{cache: cacheClient}
}

// 生成6位数字验证码
func GenerateCode() string {
	rand.Int63n(time.Now().Unix())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// 生成Redis存储Key
func (s *CaptchaService) captchaKey(userID string, device DeviceType) string {
	return fmt.Sprintf("captcha:%s:%s", userID, device)
}

// 生成尝试次数Key
func (s *CaptchaService) attemptKey(userID string, device DeviceType) string {
	return fmt.Sprintf("captcha_attempts:%s:%s", userID, device)
}

// 生成并存储验证码（带分布式锁保护）
func (s *CaptchaService) GenerateCaptcha(ctx context.Context, userID string, device DeviceType) (string, error) {
	lockKey := fmt.Sprintf("lock:captcha:%s:%s", userID, device)
	attemptKey := s.attemptKey(userID, device)

	var code string
	//var err error

	// 使用分布式锁确保操作的原子性
	lockErr := s.cache.WithSafeLock(ctx, lockKey, 3*time.Second, func() error {
		// 检查尝试次数
		attempts, err := s.cache.Get(ctx, attemptKey)
		if err != nil && !errors.Is(err, cache.ErrKeyNotFound) {
			return err
		}

		attemptCount := 0
		if attempts != "" {
			if _, err := fmt.Sscan(attempts, &attemptCount); err != nil {
				return errors.Wrap(err, "parse attempt count")
			}
		}

		if attemptCount >= MaxAttempts {
			return fmt.Errorf("max attempts reached, please try later")
		}

		// 生成验证码
		code = GenerateCode()
		captchaKey := s.captchaKey(userID, device)

		// 使用事务保证原子操作
		return s.cache.Pipeline(ctx, func(pipe redis.Pipeliner) error {
			// 设置验证码
			pipe.Set(ctx, captchaKey, code, DefaultExpiration)

			// 增加尝试次数并设置过期时间
			pipe.Incr(ctx, attemptKey)
			pipe.Expire(ctx, attemptKey, DefaultExpiration)
			return nil
		})
	})

	if lockErr != nil {
		return "", lockErr
	}
	return code, nil
}

// 验证验证码
func (s *CaptchaService) VerifyCaptcha(ctx context.Context, userID string, device DeviceType, code string) (bool, error) {
	captchaKey := s.captchaKey(userID, device)
	attemptKey := s.attemptKey(userID, device)

	// 获取存储的验证码
	storedCode, err := s.cache.Get(ctx, captchaKey)
	switch {
	case errors.Is(err, cache.ErrKeyNotFound):
		return false, fmt.Errorf("captcha expired or not generated")
	case err != nil:
		return false, err
	case storedCode != code:
		return false, nil
	}

	// 验证成功，清除相关数据
	if err := s.cache.Del(ctx, captchaKey, attemptKey); err != nil {
		return true, errors.Wrap(err, "cleanup after verification")
	}

	return true, nil
}

// 获取剩余尝试次数
func (s *CaptchaService) GetRemainingAttempts(ctx context.Context, userID string, device DeviceType) (int, error) {
	attemptKey := s.attemptKey(userID, device)

	val, err := s.cache.Get(ctx, attemptKey)
	if errors.Is(err, cache.ErrKeyNotFound) {
		return MaxAttempts, nil
	}
	if err != nil {
		return 0, err
	}

	var attempts int
	if _, err := fmt.Sscan(val, &attempts); err != nil {
		return 0, errors.Wrap(err, "parse attempt count")
	}

	return MaxAttempts - attempts, nil
}

// 高级功能：重置验证码状态（管理员接口）
func (s *CaptchaService) ResetCaptchaState(ctx context.Context, userID string, device DeviceType) error {
	captchaKey := s.captchaKey(userID, device)
	attemptKey := s.attemptKey(userID, device)

	return s.cache.Del(ctx, captchaKey, attemptKey)
}
