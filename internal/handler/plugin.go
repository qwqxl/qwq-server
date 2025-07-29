package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"qwqserver/pkg/plugin"
	pb "qwqserver/pkg/plugin/pb"
)

// RegisterPlugin 注册插件
func RegisterPlugin(c *gin.Context) {
	var req plugin.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "请求格式错误"})
		return
	}

	p := plugin.NewPluginsMap()
	p.Map(req.Name, req.Addr)
	fmt.Printf("✅ 插件注册成功: %s -> %s\n", req.Name, req.Addr)

	c.JSON(200, gin.H{"status": "ok"})
}

// LeavePlugin 注销插件
func LeavePlugin(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "请求格式错误"})
		return
	}
	p := plugin.NewPluginsMap()
	addr := p.Map(req.Name)
	fmt.Printf("✅ 插件注销成功: %s -> %s\n", req.Name, addr)
}

// ExecutePlugin 执行插件
func ExecutePlugin(c *gin.Context) {
	pluginName := c.Param("name")
	// 根据 name 获取插件地址
	//addr := "localhost:51001" + pluginName // 比如 "127.0.0.1:51001"
	p := plugin.NewPluginsMap()
	addr := p.Map(pluginName)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		c.JSON(500, gin.H{"error": "连接插件失败"})
		return
	}
	defer conn.Close()

	client := pb.NewPluginServiceClient(conn)

	// 构造请求
	var req struct {
		Action  string `json:"action"`
		Payload string `json:"payload"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "请求格式错误"})
		return
	}

	// 调用插件
	grpcResp, err := client.Execute(context.TODO(), &pb.PluginRequest{
		Action:  req.Action,
		Payload: []byte(req.Payload),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "插件执行失败"})
		return
	}

	c.JSON(200, gin.H{
		"success": grpcResp.Success,
		"message": grpcResp.Message,
		"data":    string(grpcResp.Data),
	})
}
