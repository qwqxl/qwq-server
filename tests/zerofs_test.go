package qwqtest

import (
	"fmt"
	"log"
	"os"
	"qwqserver/pkg/util/zerofs"
	"testing"
)

func TestZeroFsRun1(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		// 定义分割选项
		opts := zerofs.SplitOptions{
			ChunkSize:  5 * 1024 * 1024, // 5MB分片
			OutputDir:  "chunks",        // 分片保存目录
			BufferSize: 32 * 1024,       // 32KB缓冲区
		}

		// 分割文件
		parts, err := zerofs.SplitFile("largefile.zip", opts)
		if err != nil {
			log.Fatal("分割文件失败:", err)
		}

		fmt.Printf("文件已分割为 %d 个分片:\n", len(parts))
		for _, part := range parts {
			fmt.Printf("分片 %d: 偏移 %d, 大小 %d 字节, 路径 %s\n",
				part.Index, part.Offset, part.Size, part.Path)
		}
	})
}

func TestZeroFsRun(t *testing.T) {
	t.Run("test", func(t *testing.T) {

		// 示例1: 基本文件分割
		parts, err := zerofs.SplitFile("largefile.zip", zerofs.SplitOptions{
			ChunkSize: 50 << 20, // 50MB分片
			OutputDir: "chunks",
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("分割为 %d 个分片\n", len(parts))

		// 示例2: 使用自定义处理器
		customProcessor := &CustomProcessor{}
		_, err = zerofs.SplitFile("data.bin", zerofs.SplitOptions{
			ChunkSize: 5 << 20,
			Processor: customProcessor,
		})

		// 示例3: 合并文件
		err = zerofs.MergeFiles(parts, "restored.zip")
		if err != nil {
			panic(err)
		}

		// 示例4: 零拷贝传输
		src, _ := os.Open("source.iso")
		defer src.Close()
		dest, _ := os.Create("destination.iso")
		defer dest.Close()
		//_, err = zerofs.ZeroCopyTransfer(src, dest)
	})
}

// 自定义处理器示例
type CustomProcessor struct{}

func (p *CustomProcessor) ProcessChunk(chunk []byte, part zerofs.FilePart) error {
	// 这里实现自定义处理逻辑
	fmt.Printf("处理分片 %d (大小: %d 字节)\n", part.Index, len(chunk))
	// 例如: 加密、压缩、上传等
	return nil
}
