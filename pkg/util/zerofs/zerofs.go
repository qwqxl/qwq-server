package zerofs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// 常量定义
const (
	// DefaultChunkSize 默认文件分片大小 (10MB)
	DefaultChunkSize = 10 << 20 // 10 MB

	// DefaultBufferSize 默认缓冲区大小 (64KB)
	DefaultBufferSize = 64 << 10 // 64 KB
)

// FilePart 文件分片信息结构体
type FilePart struct {
	Index  int    // 分片索引（从0开始）
	Offset int64  // 分片在文件中的起始位置（字节）
	Size   int64  // 分片大小（字节）
	Path   string // 分片文件路径（如果保存为文件）
}

// FileProcessor 文件处理接口
// 用于在读取文件分片时执行自定义操作（如加密、压缩、上传等）
type FileProcessor interface {
	// ProcessChunk 处理文件分片
	// chunk: 当前分片的数据
	// part: 当前分片的信息
	ProcessChunk(chunk []byte, part FilePart) error
}

// SplitOptions 文件分割选项
type SplitOptions struct {
	ChunkSize  int64         // 每个分片的大小（字节），0表示使用默认值
	OutputDir  string        // 分片文件输出目录（如果保存为文件）
	Processor  FileProcessor // 自定义分片处理器（如果不保存为文件）
	BufferSize int           // 读取缓冲区大小（字节），0表示使用默认值
}

// SplitFile 将文件分割为多个分片
// filePath: 要分割的文件路径
// opts: 分割选项
// 返回值: 分片信息列表和可能的错误
func SplitFile(filePath string, opts SplitOptions) ([]FilePart, error) {
	// 设置默认值
	if opts.ChunkSize <= 0 {
		opts.ChunkSize = DefaultChunkSize
	}
	if opts.BufferSize <= 0 {
		opts.BufferSize = DefaultBufferSize
	}

	// 打开源文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// 计算分片数量
	totalParts := fileSize / opts.ChunkSize
	if fileSize%opts.ChunkSize != 0 {
		totalParts++ // 如果文件大小不是分片大小的整数倍，增加一个分片
	}

	// 创建分片信息列表
	parts := make([]FilePart, 0, totalParts)

	// 创建读取缓冲区
	buffer := make([]byte, opts.BufferSize)

	// 遍历所有分片
	for i := int64(0); i < totalParts; i++ {
		// 计算当前分片的偏移量和大小
		offset := i * opts.ChunkSize
		remaining := fileSize - offset
		chunkSize := opts.ChunkSize
		if remaining < chunkSize {
			chunkSize = remaining
		}

		// 创建分片信息
		part := FilePart{
			Index:  int(i),
			Offset: offset,
			Size:   chunkSize,
		}

		// 处理分片
		if opts.Processor != nil {
			// 使用自定义处理器处理分片
			if err := processWithCustomHandler(file, buffer, part, opts.Processor); err != nil {
				return nil, err
			}
		} else {
			// 保存分片到文件
			if opts.OutputDir == "" {
				return nil, errors.New("输出目录未指定，请提供OutputDir或使用自定义处理器")
			}

			// 生成分片文件名
			part.Path = filepath.Join(opts.OutputDir, generateChunkFilename(fileInfo.Name(), i))
			if err := saveChunkToFile(file, buffer, part); err != nil {
				return nil, err
			}
		}

		// 添加到分片列表
		parts = append(parts, part)
	}

	return parts, nil
}

// MergeFiles 将多个分片文件合并为单个文件
// parts: 分片信息列表（必须按Index顺序）
// outputPath: 合并后的文件输出路径
// 返回值: 可能的错误
func MergeFiles(parts []FilePart, outputPath string) error {
	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 按顺序合并所有分片
	for _, part := range parts {
		if part.Path == "" {
			return errors.New("分片路径为空，无法合并")
		}

		// 打开分片文件
		chunkFile, err := os.Open(part.Path)
		if err != nil {
			return err
		}

		// 将分片内容复制到输出文件
		if _, err := io.Copy(outputFile, chunkFile); err != nil {
			chunkFile.Close()
			return err
		}

		// 关闭当前分片文件
		chunkFile.Close()
	}

	return nil
}

// ZeroCopyTransfer 使用零拷贝技术传输文件
// 在支持的系统上使用sendfile系统调用，避免数据在用户空间和内核空间之间复制
// src: 源文件
// dest: 目标文件
// 返回值: 传输的字节数和可能的错误
//func ZeroCopyTransfer(src *os.File, dest *os.File) (int64, error) {
//	// 获取源文件信息
//	srcInfo, err := src.Stat()
//	if err != nil {
//		return 0, err
//	}
//
//	// 在Linux系统上使用sendfile系统调用
//	return syscall.Sendfile(int(dest.Fd()), int(src.Fd()), nil, int(srcInfo.Size()))
//}

// generateChunkFilename 生成分片文件名
// baseName: 原始文件名
// index: 分片索引
func generateChunkFilename(baseName string, index int64) string {
	return baseName + ".part" + formatIndex(int(index))
}

// formatIndex 格式化分片索引为三位数字符串
// 例如：0 -> "000", 5 -> "005", 123 -> "123"
func formatIndex(i int) string {
	if i < 10 {
		return "00" + string(rune('0'+i))
	}
	if i < 100 {
		return "0" + string(rune('0'+i/10)) + string(rune('0'+i%10))
	}
	return string(rune('0'+i/100)) + string(rune('0'+(i/10)%10)) + string(rune('0'+i%10))
}

// processWithCustomHandler 使用自定义处理器处理分片
// file: 源文件
// buffer: 读取缓冲区
// part: 分片信息
// processor: 自定义处理器
func processWithCustomHandler(file *os.File, buffer []byte, part FilePart, processor FileProcessor) error {
	// 定位到分片起始位置
	if _, err := file.Seek(part.Offset, io.SeekStart); err != nil {
		return err
	}

	remaining := part.Size
	for remaining > 0 {
		// 计算本次读取的大小
		readSize := int64(len(buffer))
		if remaining < readSize {
			readSize = remaining
		}

		// 读取数据到缓冲区
		n, err := file.Read(buffer[:readSize])
		if err != nil && err != io.EOF {
			return err
		}

		// 处理读取的数据
		if err := processor.ProcessChunk(buffer[:n], part); err != nil {
			return err
		}

		// 更新剩余字节数
		remaining -= int64(n)
		if n == 0 {
			break // 已读取完所有数据
		}
	}

	return nil
}

// saveChunkToFile 将分片保存到文件
// file: 源文件
// buffer: 读取缓冲区
// part: 分片信息
func saveChunkToFile(file *os.File, buffer []byte, part FilePart) error {
	// 创建分片文件
	chunkFile, err := os.Create(part.Path)
	if err != nil {
		return err
	}
	defer chunkFile.Close()

	// 定位到分片起始位置
	if _, err := file.Seek(part.Offset, io.SeekStart); err != nil {
		return err
	}

	remaining := part.Size
	for remaining > 0 {
		// 计算本次读取的大小
		readSize := int64(len(buffer))
		if remaining < readSize {
			readSize = remaining
		}

		// 读取数据到缓冲区
		n, err := file.Read(buffer[:readSize])
		if err != nil && err != io.EOF {
			return err
		}

		// 将数据写入分片文件
		if _, err := chunkFile.Write(buffer[:n]); err != nil {
			return err
		}

		// 更新剩余字节数
		remaining -= int64(n)
		if n == 0 {
			break // 已读取完所有数据
		}
	}

	return nil
}
