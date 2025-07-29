package handler

//func Docs(relativePath string, r *gin.Engine, middlewares ...gin.HandlerFunc) {
//	group := r.Group(relativePath)
//	// 认证中间件
//	group.Use(middlewares...)
//
//	// 获取所有 Markdown 文件
//	mdDir := util.WorkDir(fmt.Sprintf("resources/docs"))
//	mdTemplateDir := util.WorkDir("resources/templates/markdown")
//
//	// 添加静态文件服务 - 关键部分
//	r.Static("/img", mdDir)
//
//	// 加载模板
//	r.LoadHTMLGlob(filepath.Join(mdTemplateDir, "*.html"))
//
//	// 获取所有 Markdown 文件
//	group.GET("", func(c *gin.Context) {
//		util.ListMarkdownFiles(mdDir, c)
//	})
//
//	// 渲染 Markdown 文件
//	group.GET("/:filename", func(c *gin.Context) {
//		util.RenderMarkdown(mdDir, c)
//	})
//}

//func Docs(relativePath string, r *server.HTTPServer, middlewares ...gin.HandlerFunc) {
//	group := r.Group(relativePath)
//	group.Use(middlewares...)
//
//	mdDir := util.WorkDir("resources/docs")
//	mdTemplateDir := util.WorkDir("resources/templates/markdown")
//
//	// 静态资源（如需鉴权，可以使用 group + 自定义 handler）
//	group.Static("/img", filepath.Join(mdDir, "img"))
//
//	// 加载所有 html 模板（含子目录）
//	r.LoadHTMLGlob(filepath.Join(mdTemplateDir, "*.html"))
//
//	group.GET("", func(c *gin.Context) {
//		util.ListMarkdownFiles(mdDir, c)
//	})
//
//	group.GET("/:filename", func(c *gin.Context) {
//		util.RenderMarkdown(mdDir, c)
//	})
//}
