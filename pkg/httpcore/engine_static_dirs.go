package httpcore

import (
	"crypto/md5"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type StaticConfig struct {
	CacheMaxAge     time.Duration     // 缓存最大时长
	EnableETag      bool              // 是否启用ETag
	EnableBrotli    bool              // 是否支持Brotli压缩
	EnableGzip      bool              // 是否支持Gzip压缩
	IndexFiles      []string          // 默认索引文件
	FileExtensions  []string          // 允许的文件扩展名
	ExcludePatterns []string          // 排除路径模式
	SecurityHeaders map[string]string // 安全头设置
}

// 安全增强
func (e *Engine) addSecurityHeaders(c Context) {
	wc := c.(*contextWrapper)

	// 设置安全头
	wc.ctx.Header("X-Content-Type-Options", "nosniff")
	wc.ctx.Header("X-Frame-Options", "DENY")
	wc.ctx.Header("X-XSS-Protection", "1; mode=block")

	// 设置CSP（内容安全策略）
	csp := "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';"
	wc.ctx.Header("Content-Security-Policy", csp)

	// 对于特定文件类型添加额外安全头
	if strings.HasSuffix(wc.ctx.Request.URL.Path, ".wasm") {
		wc.ctx.Header("Content-Type", "application/wasm")
	}
}

// 文件缓存
type fileCache struct {
	content     []byte
	contentType string
	etag        string
	lastMod     time.Time
}

var (
	fileCacheMap sync.Map
	cacheMutex   sync.RWMutex
)

func (e *Engine) serveCachedFile(c Context, filePath string) {
	// 检查缓存
	if entry, ok := fileCacheMap.Load(filePath); ok {
		cached := entry.(*fileCache)

		// 检查 ETag
		if match := c.GetHeader("If-None-Match"); match == cached.etag {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}

		// 设置缓存头
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.Header("ETag", cached.etag)
		c.Header("Content-Type", cached.contentType)
		c.Status(http.StatusOK)
		c.Response().Write(cached.content)
		return
	}

	// 缓存未命中，读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// 获取文件信息
	stat, err := os.Stat(filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 计算 ETag
	etag := fmt.Sprintf("%x", md5.Sum(content))

	// 确定内容类型
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 缓存文件
	cached := &fileCache{
		content:     content,
		contentType: contentType,
		etag:        etag,
		lastMod:     stat.ModTime(),
	}

	fileCacheMap.Store(filePath, cached)

	// 响应文件
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", etag)
	c.Header("Content-Type", contentType)
	c.Status(http.StatusOK)
	c.Response().Write(content)
}

// 压缩
func (e *Engine) compressionMiddleware() HandlerFunc {
	return func(c Context) {
		wc := c.(*contextWrapper)

		// 检查是否支持压缩
		acceptEncoding := wc.ctx.GetHeader("Accept-Encoding")

		// 对于静态文件请求，优先使用预压缩文件
		if strings.Contains(acceptEncoding, "br") {
			if e.servePrecompressedFile(wc, ".br", "br") {
				return
			}
		}

		if strings.Contains(acceptEncoding, "gzip") {
			if e.servePrecompressedFile(wc, ".gz", "gzip") {
				return
			}
		}

		wc.ctx.Next()
	}
}

func (e *Engine) servePrecompressedFile(wc *contextWrapper, ext, encoding string) bool {
	// 检查预压缩文件是否存在
	compressedPath := wc.ctx.Request.URL.Path + ext
	for prefix, dir := range e.staticDirs {
		if strings.HasPrefix(compressedPath, prefix) {
			filePath := filepath.Join(dir, strings.TrimPrefix(compressedPath, prefix))
			if _, err := os.Stat(filePath); err == nil {
				wc.ctx.Header("Content-Encoding", encoding)
				wc.ctx.Header("Vary", "Accept-Encoding")
				wc.ctx.File(filePath)
				return true
			}
		}
	}
	return false
}

// AddStatic 添加单个静态文件映射
func (e *Engine) AddStatic(urlPrefix, localDir string) {
	if e.staticDirs == nil {
		e.staticDirs = make(map[string]string)
	}

	// 确保路径格式正确
	if !strings.HasSuffix(urlPrefix, "/") {
		urlPrefix += "/"
	}

	// 检查目录是否存在
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		e.logger.Warn("Static directory not found: %s", localDir)
		return
	}

	e.staticDirs[urlPrefix] = localDir
	e.logger.Info("Mapped static path: %s -> %s", urlPrefix, localDir)
}

// AddStaticFS 添加文件系统映射（支持嵌入文件系统）
func (e *Engine) AddStaticFS(urlPrefix string, fs http.FileSystem) {
	e.engine.StaticFS(urlPrefix, fs)
	e.logger.Info("Mapped static FS: %s", urlPrefix)
}

// AddStaticFile 添加单个静态文件
func (e *Engine) AddStaticFile(urlPath, localFile string) {
	if _, err := os.Stat(localFile); os.IsNotExist(err) {
		e.logger.Warn("Static file not found: %s", localFile)
		return
	}

	e.engine.StaticFile(urlPath, localFile)
	e.logger.Info("Mapped static file: %s -> %s", urlPath, localFile)
}

// middleware

// staticMiddleware 静态文件服务中间件
func (e *Engine) staticMiddleware() HandlerFunc {
	// 创建文件服务器处理器
	fileHandler := func(c Context) {
		wc := c.(*contextWrapper)
		reqUrlPath := wc.ctx.Request.URL.Path

		// 检查是否匹配静态路径
		for prefix, dir := range e.staticDirs {
			if strings.HasPrefix(reqUrlPath, prefix) {
				// 构建本地文件路径
				filePath := filepath.Join(dir, strings.TrimPrefix(reqUrlPath, prefix))

				// 检查文件是否存在
				if stat, err := os.Stat(filePath); err == nil {
					// 设置缓存控制头
					wc.ctx.Header("Cache-Control", "public, max-age=3600")

					// 设置 ETag
					etag := fmt.Sprintf("%x", md5.Sum([]byte(filePath+stat.ModTime().String())))
					wc.ctx.Header("ETag", etag)

					// 检查 If-None-Match
					if match := wc.ctx.GetHeader("If-None-Match"); match == etag {
						wc.ctx.AbortWithStatus(http.StatusNotModified)
						return
					}

					// 发送文件
					wc.ctx.File(filePath)
					return
				}
			}
		}

		// 不是静态文件请求，继续处理
		wc.ctx.Next()
	}

	return fileHandler
}

// formatFileSize 将字节大小格式化为人类可读的字符串
func formatFileSize(bytes int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
