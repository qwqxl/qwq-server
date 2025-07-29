package util

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// ListMarkdownFiles 列出所有 Markdown 文件（支持子目录）
func ListMarkdownFiles(mdDir string, c *gin.Context) {
	files, err := getMarkdownFiles(mdDir)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "无法读取文件列表: " + err.Error(),
		})
		return
	}

	// 按目录分组
	dirMap := make(map[string][]FileInfo)
	for _, file := range files {
		dir := filepath.Dir(file.RelativePath)
		dirMap[dir] = append(dirMap[dir], file)
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"Title":      "文档列表",
		"BaseDir":    mdDir,
		"DirMap":     dirMap,
		"Breadcrumb": getBreadcrumb(c.Request.URL.Path),
	})
}

// RenderMarkdown 渲染 Markdown 文件（支持子目录）
func RenderMarkdown(mdDir string, c *gin.Context) {
	// 获取路径参数，支持子目录
	pathParam := c.Param("filename")

	// 确保路径安全
	if strings.Contains(pathParam, "..") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "无效的文件路径",
		})
		return
	}

	// 添加 .md 扩展名（如果不存在）
	filename := pathParam
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	filePath := filepath.Join(mdDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "文件未找到: " + filename,
		})
		return
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "无法读取文件: " + err.Error(),
		})
		return
	}

	// 转换 Markdown 为 HTML
	// 然后在 RenderMarkdown 中使用：
	htmlContent := mdToHTML(content)
	//htmlContent = fixImagePaths(htmlContent, "sss", "sss")
	//htmlContent := mdToHTML(content)

	// 获取不带扩展名的文件名
	baseName := strings.TrimSuffix(filepath.Base(filename), ".md")

	c.HTML(http.StatusOK, "markdown.html", gin.H{
		"Title":      baseName,
		"Content":    template.HTML(htmlContent),
		"FilePath":   filename,
		"Breadcrumb": getBreadcrumb(c.Request.URL.Path),
	})
}

// FileInfo 存储文件信息
type FileInfo struct {
	Name         string
	RelativePath string // 相对于根目录的路径
	IsDir        bool
	ModTime      time.Time
}

// 递归获取所有 Markdown 文件
func getMarkdownFiles(rootDir string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// 跳过根目录本身
		if relPath == "." {
			return nil
		}

		// 处理目录
		if info.IsDir() {
			// 跳过隐藏目录（以 . 开头）
			if strings.HasPrefix(filepath.Base(relPath), ".") {
				return filepath.SkipDir
			}
			files = append(files, FileInfo{
				Name:         filepath.Base(relPath),
				RelativePath: relPath,
				IsDir:        true,
				ModTime:      info.ModTime(),
			})
			return nil
		}

		// 处理文件 - 只添加 Markdown 文件
		if strings.HasSuffix(strings.ToLower(path), ".md") {
			files = append(files, FileInfo{
				Name:         strings.TrimSuffix(filepath.Base(relPath), ".md"),
				RelativePath: relPath,
				IsDir:        false,
				ModTime:      info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// 将 Markdown 转换为 HTML
func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.ToHTML(md, p, renderer)
}

// 生成面包屑导航
func getBreadcrumb(path string) []BreadcrumbItem {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var items []BreadcrumbItem
	accum := ""

	for i, part := range parts {
		accum += "/" + part
		isLast := i == len(parts)-1

		items = append(items, BreadcrumbItem{
			Name: part,
			Path: accum,
			Last: isLast,
		})
	}

	// 添加首页
	if len(items) > 0 {
		items = append([]BreadcrumbItem{{Name: "首页", Path: "/"}}, items...)
	}

	return items
}

// BreadcrumbItem 面包屑导航项
type BreadcrumbItem struct {
	Name string
	Path string
	Last bool
}

//package util
//
//import (
//	"html/template"
//	"io/ioutil"
//	"net/http"
//	"path/filepath"
//	"strings"
//
//	"github.com/gin-gonic/gin"
//	"github.com/gomarkdown/markdown"
//	"github.com/gomarkdown/markdown/html"
//	"github.com/gomarkdown/markdown/parser"
//)
//
//// ListMarkdownFiles 列出所有 Markdown 文件
//func ListMarkdownFiles(mdDir string, c *gin.Context) {
//	files, err := getMarkdownFiles(mdDir)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
//			"error": "无法读取文件列表",
//		})
//		return
//	}
//
//	c.HTML(http.StatusOK, "list.html", gin.H{
//		"Title": "Markdown 文件列表",
//		"Files": files,
//	})
//}
//
//// RenderMarkdown 渲染 Markdown 文件
//func RenderMarkdown(mdDir string, c *gin.Context) {
//	filename := c.Param("filename")
//	if !strings.HasSuffix(filename, ".md") {
//		filename += ".md"
//	}
//
//	// 防止路径遍历
//	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
//			"error": "无效的文件名",
//		})
//		return
//	}
//
//	filePath := filepath.Join(mdDir, filename)
//	content, err := ioutil.ReadFile(filePath)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
//			"error": "文件未找到",
//		})
//		return
//	}
//
//	// 转换 Markdown 为 HTML
//	htmlContent := mdToHTML(content)
//
//	// 获取不带扩展名的文件名
//	baseName := strings.TrimSuffix(filename, ".md")
//
//	c.HTML(http.StatusOK, "markdown.html", gin.H{
//		"Title":   baseName,
//		"Content": template.HTML(htmlContent),
//	})
//}
//
//// 获取所有 Markdown 文件
//func getMarkdownFiles(mdDir string) ([]string, error) {
//	files, err := ioutil.ReadDir(mdDir)
//	if err != nil {
//		return nil, err
//	}
//
//	var mdFiles []string
//	for _, file := range files {
//		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
//			mdFiles = append(mdFiles, strings.TrimSuffix(file.Name(), ".md"))
//		}
//	}
//	return mdFiles, nil
//}
//
//// 将 Markdown 转换为 HTML
//func mdToHTML(md []byte) []byte {
//	// 创建 Markdown 解析器
//	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
//	p := parser.NewWithExtensions(extensions)
//
//	// 创建 HTML 渲染器
//	htmlFlags := html.CommonFlags | html.HrefTargetBlank
//	opts := html.RendererOptions{Flags: htmlFlags}
//	renderer := html.NewRenderer(opts)
//
//	// 转换
//	return markdown.ToHTML(md, p, renderer)
//}
