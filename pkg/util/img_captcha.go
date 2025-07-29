// Package util 提供图形验证码生成工具
package util

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/mojocn/base64Captcha"
)

// ImgCaptchaType 表示验证码类型
//   - CaptchaTypeString：默认字符验证码
//   - CaptchaTypeChinese：中文验证码
//   - CaptchaTypeMath：数学算式验证码
type ImgCaptchaType int

const (
	ImgCaptchaTypeString  ImgCaptchaType = iota // 字符验证码
	ImgCaptchaTypeChinese                       // 中文验证码
	ImgCaptchaTypeMath                          // 数学验证码
)

const (
	ImgCaptchaSource     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"  // 验证码字符源
	ImgCaptchaFontFile   = "resources/fonst/WenQuanYiMicroHei.ttf" // 默认字体文件
	ImgCaptchaImgHeight  = 80                                      // 验证码图片高度
	ImgCaptchaImgWidth   = 200                                     // 验证码图片宽度
	ImgCaptchaNoiseCount = 0                                       // 干扰点数量
	ImgCaptchaLength     = 4                                       // 验证码字符长度
)

var imgCaptchaStore = base64Captcha.DefaultMemStore

// ImgCaptchaConfig 表示图形验证码配置项
type ImgCaptchaConfig struct {
	Type       ImgCaptchaType // 验证码类型：字符 / 中文 / 数学
	Width      int            // 图片宽度
	Height     int            // 图片高度
	Length     int            // 验证码长度
	NoiseCount int            // 干扰点数量
	Source     string         // 字符验证码字符源
	FontFiles  []string       // 字体文件
	ShowLine   int            // 干扰线类型
}

// ImgCaptchaDefaultConfig 默认验证码配置
var ImgCaptchaDefaultConfig = ImgCaptchaConfig{
	Type:       ImgCaptchaTypeString,
	Width:      ImgCaptchaImgWidth,
	Height:     ImgCaptchaImgHeight,
	Length:     ImgCaptchaLength,
	NoiseCount: ImgCaptchaNoiseCount,
	Source:     ImgCaptchaSource,
	FontFiles:  []string{ImgCaptchaFontFile},
	ShowLine:   base64Captcha.OptionShowSineLine,
}

// GenImgCaptchaBase64 生成图形验证码并返回 Base64 编码字符串
// 参数：
//   - cfg: 可选配置项，为 nil 则使用默认配置
//
// 返回值：
//   - string: 图片的 Base64 字符串
//   - string: 验证码正确答案
//   - error: 错误信息
func GenImgCaptchaBase64(cfg *ImgCaptchaConfig) (string, string, error) {
	imgBytes, answer, err := GenImgCaptchaBytes(cfg)
	if err != nil {
		return "", "", err
	}
	b64 := base64.StdEncoding.EncodeToString(imgBytes)
	return b64, answer, nil
}

// GenImgCaptchaBytes 生成图形验证码并返回图片字节内容
// 参数：
//   - cfg: 可选配置项，为 nil 则使用默认配置
//
// 返回值：
//   - []byte: 图片的字节内容
//   - string: 验证码正确答案
//   - error: 错误信息
func GenImgCaptchaBytes(cfg *ImgCaptchaConfig) ([]byte, string, error) {
	if cfg == nil {
		cfg = &ImgCaptchaDefaultConfig
	}

	driver := createDriver(cfg)
	captcha := base64Captcha.NewCaptcha(driver, imgCaptchaStore)

	_, content, answer := captcha.Driver.GenerateIdQuestionAnswer()
	item, err := captcha.Driver.DrawCaptcha(content)
	if err != nil {
		return nil, "", fmt.Errorf("生成图形验证码失败: %v", err)
	}

	var buf bytes.Buffer
	_, err = item.WriteTo(&buf)
	if err != nil {
		return nil, "", fmt.Errorf("图像写入失败: %v", err)
	}
	return buf.Bytes(), answer, nil
}

// createDriver 根据配置创建验证码驱动
// 参数：
//   - cfg: 配置项
//
// 返回值：
//   - base64Captcha.Driver 接口实例（内部自动选择）
func createDriver(cfg *ImgCaptchaConfig) base64Captcha.Driver {
	switch cfg.Type {
	case ImgCaptchaTypeChinese:
		return &base64Captcha.DriverChinese{
			Height:          cfg.Height,
			Width:           cfg.Width,
			Length:          cfg.Length,
			NoiseCount:      cfg.NoiseCount,
			ShowLineOptions: cfg.ShowLine,
			Fonts:           cfg.FontFiles,
		}
	case ImgCaptchaTypeMath:
		return &base64Captcha.DriverMath{
			Height:          cfg.Height,
			Width:           cfg.Width,
			NoiseCount:      cfg.NoiseCount,
			ShowLineOptions: cfg.ShowLine,
			Fonts:           cfg.FontFiles,
		}
	default:
		return &base64Captcha.DriverString{
			Height:          cfg.Height,
			Width:           cfg.Width,
			Length:          cfg.Length,
			NoiseCount:      cfg.NoiseCount,
			ShowLineOptions: cfg.ShowLine,
			Source:          cfg.Source,
			Fonts:           cfg.FontFiles,
		}
	}
}
