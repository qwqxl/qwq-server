package util

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 在渲染后修复HTML中的图片路径
func fixImagePaths(htmlContent []byte, baseDir string, relativeBase string) []byte {
	// 使用正则表达式查找所有图片标签
	re := regexp.MustCompile(`<img[^>]+src="([^"]+)"[^>]*>`)

	return re.ReplaceAllFunc(htmlContent, func(match []byte) []byte {
		submatches := re.FindSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		src := string(submatches[1])

		// 检查是否是相对路径
		if !strings.HasPrefix(src, "http://") &&
			!strings.HasPrefix(src, "https://") &&
			!strings.HasPrefix(src, "/") {

			// 创建绝对文件路径
			absPath := filepath.Join(baseDir, src)

			// 检查文件是否存在
			if _, err := os.Stat(absPath); err == nil {
				// 转换为URL路径
				urlPath := "/static/docs/" + filepath.ToSlash(filepath.Join(relativeBase, src))
				return bytes.Replace(match, submatches[1], []byte(urlPath), 1)
			}
		}

		return match
	})
}
