// Package util 雪花算法ID生成
package util

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	snowflakeNode *snowflake.Node
	snowflakeOnce sync.Once
)

// GenerateID 生成雪花算法 ID
// 返回值：
//   - int64: 生成的雪花算法 ID
//   - error: 操作过程中的错误
func GenerateID() (int64, error) {
	snowflakeOnce.Do(func() {
		var err error
		snowflakeNode, err = snowflake.NewNode(1)
		if err != nil {
			fmt.Printf("初始化雪花算法节点失败: %v", err)
		}
	})

	switch {
	case snowflakeNode != nil:
		return snowflakeNode.Generate().Int64(), nil

	default:
		// 雪花格式: 41 位时间戳 + 10 位节点 ID + 12 位序列号
		// 标准雪花纪元，节点 ID 1，序列号使用当前纳秒的低 12 位
		ts := time.Now().UnixMilli() - 1288834974657
		nodeID := int64(1)
		seq := time.Now().UnixNano() & 0xFFF

		return (ts << 22) | (nodeID << 12) | seq, nil
	}
}
