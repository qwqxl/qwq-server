package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	UploadDir = "./uploads" // 上传文件存储目录
	TempDir   = "./temp"    // 临时分片存储目录
	MaxMemory = 32 << 20    // 32MB (最大内存使用)
)

func main3() {
	// 确保上传目录存在
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		log.Fatalf("无法创建上传目录: %v", err)
	}
	if err := os.MkdirAll(TempDir, 0755); err != nil {
		log.Fatalf("无法创建临时目录: %v", err)
	}

	router := gin.Default()

	// 设置模板
	router.LoadHTMLGlob("web/html/*")

	// 首页 - 上传表单
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "大文件上传示例",
		})
	})

	// 检查文件分片是否已存在（用于断点续传）
	router.HEAD("/upload/:fileId/:chunkIndex", checkChunkExists)

	// 上传文件分片
	router.POST("/upload/:fileId/:chunkIndex", uploadChunk)

	// 合并文件分片
	router.POST("/merge/:fileId", mergeChunks)

	// 获取上传进度
	router.GET("/progress/:fileId", getUploadProgress)

	// 启动服务器
	fmt.Println("文件上传服务启动，访问 http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 检查分片是否存在
func checkChunkExists(c *gin.Context) {
	fileId := c.Param("fileId")
	chunkIndex := c.Param("chunkIndex")
	chunkPath := filepath.Join(TempDir, fileId, chunkIndex)

	if _, err := os.Stat(chunkPath); err == nil {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// 上传文件分片
func uploadChunk(c *gin.Context) {
	fileId := c.Param("fileId")
	chunkIndex := c.Param("chunkIndex")
	totalChunks := c.PostForm("totalChunks")
	chunkSize := c.PostForm("chunkSize")
	fileName := c.PostForm("fileName")

	// 创建文件ID目录
	chunkDir := filepath.Join(TempDir, fileId)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建分片目录"})
		return
	}

	// 处理文件上传
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取上传文件"})
		return
	}
	defer file.Close()

	// 创建分片文件
	chunkPath := filepath.Join(chunkDir, chunkIndex)
	dst, err := os.Create(chunkPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建分片文件"})
		return
	}
	defer dst.Close()

	// 保存分片数据
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存分片数据"})
		return
	}

	// 更新进度信息
	updateProgress(fileId, fileName, chunkIndex, totalChunks, chunkSize)

	c.JSON(http.StatusOK, gin.H{
		"message":   "分片上传成功",
		"fileId":    fileId,
		"chunk":     chunkIndex,
		"fileName":  fileName,
		"chunkSize": header.Size,
	})
}

// 合并文件分片
func mergeChunks(c *gin.Context) {
	fileId := c.Param("fileId")
	fileName := c.PostForm("fileName")
	totalChunks, _ := strconv.Atoi(c.PostForm("totalChunks"))

	// 创建最终文件
	finalPath := filepath.Join(UploadDir, fileName)
	finalFile, err := os.Create(finalPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建最终文件"})
		return
	}
	defer finalFile.Close()

	// 合并所有分片
	chunkDir := filepath.Join(TempDir, fileId)
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, strconv.Itoa(i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开分片文件"})
			return
		}

		// 复制分片内容到最终文件
		if _, err := io.Copy(finalFile, chunkFile); err != nil {
			chunkFile.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "合并文件失败"})
			return
		}
		chunkFile.Close()
	}

	// 清理临时文件
	if err := os.RemoveAll(chunkDir); err != nil {
		log.Printf("警告: 无法清理临时文件: %v", err)
	}

	// 删除进度信息
	deleteProgress(fileId)

	c.JSON(http.StatusOK, gin.H{
		"message":  "文件合并成功",
		"fileId":   fileId,
		"fileName": fileName,
		"path":     finalPath,
	})
}

// 获取上传进度
func getUploadProgress(c *gin.Context) {
	fileId := c.Param("fileId")
	progress := getProgress(fileId)

	if progress == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到上传进度"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// 进度跟踪数据结构
type UploadProgress struct {
	FileID      string    `json:"fileId"`
	FileName    string    `json:"fileName"`
	TotalChunks int       `json:"totalChunks"`
	Uploaded    []int     `json:"uploaded"` // 已上传的分片索引
	ChunkSize   int64     `json:"chunkSize"`
	StartTime   time.Time `json:"startTime"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

var progressMap = make(map[string]*UploadProgress)

// 更新进度信息
func updateProgress(fileId, fileName, chunkIndex, totalChunks, chunkSize string) {
	progress, exists := progressMap[fileId]
	if !exists {
		total, _ := strconv.Atoi(totalChunks)
		chunkSizeVal, _ := strconv.ParseInt(chunkSize, 10, 64)
		progress = &UploadProgress{
			FileID:      fileId,
			FileName:    fileName,
			TotalChunks: total,
			Uploaded:    []int{},
			ChunkSize:   chunkSizeVal,
			StartTime:   time.Now(),
			LastUpdate:  time.Now(),
		}
		progressMap[fileId] = progress
	}

	// 添加已上传的分片索引
	index, _ := strconv.Atoi(chunkIndex)
	progress.Uploaded = append(progress.Uploaded, index)
	progress.LastUpdate = time.Now()
}

// 获取进度信息
func getProgress(fileId string) *UploadProgress {
	return progressMap[fileId]
}

// 删除进度信息
func deleteProgress(fileId string) {
	delete(progressMap, fileId)
}
