package qwqtest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path/filepath"
	"qwqserver/pkg/util"
	"testing"
)

const (
	mdDir       = "./markdown"  // Markdown 文件存储目录
	templateDir = "./templates" // HTML 模板目录
)

func TestMdToHTML(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		// 初始化 Gin 引擎
		r := gin.Default()

		// 加载 HTML 模板
		r.LoadHTMLGlob(filepath.Join(templateDir, "*.html"))

		// 创建 markdown 目录（如果不存在）
		if err := os.MkdirAll(mdDir, 0755); err != nil {
			panic(fmt.Sprintf("创建 markdown 目录失败: %v", err))
		}

		// 示例：写入一个测试 Markdown 文件
		exampleMD := []byte(`# Hello Markdown!
	
**这是一个示例文件**

- 列表项 1
- 列表项 2

[返回首页](/)`)
		if err := ioutil.WriteFile(filepath.Join(mdDir, "example.md"), exampleMD, 0644); err != nil {
			fmt.Printf("创建示例文件失败: %v\n", err)
		}

		// 路由定义
		r.GET("/", func(c *gin.Context) {
			util.ListMarkdownFiles(mdDir, c)
		})
		r.GET("/md/:filename", func(c *gin.Context) {
			util.RenderMarkdown(mdDir, c)
		})

		// 启动服务器
		fmt.Println("服务器运行在 http://localhost:8080")
		if err := r.Run(":8080"); err != nil {
			panic(fmt.Sprintf("服务器启动失败: %v", err))
		}

	})
}
