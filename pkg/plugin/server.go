package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	pb "qwqserver/pkg/plugin/pb"
	"time"
)

// ServiceServer 实现接口
type ServiceServer struct {
	pb.UnimplementedPluginServiceServer // 推荐嵌入，防止接口变动报错
}

type ServiceRequest struct {
	Action  string `json:"action"`
	Payload []byte `json:"payload"`
}

type ServiceResponse struct {
	Result []byte `json:"result"`
}

// Execute 实现方法
func (s *ServiceServer) Execute(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	log.Printf("🔧 收到执行请求: action=%s, payload=%d字节", req.GetAction(), len(req.GetPayload()))

	// 这里你可以加上实际的插件调用逻辑
	// 暂时只返回成功
	return &pb.PluginResponse{
		Success: true,
		Message: fmt.Sprintf("已成功执行 action: %s", req.GetAction()),
		Data:    []byte("plugin-result"), // 假装返回了一些数据
	}, nil
}

// HealthCheck 实现方法
func (s *ServiceServer) HealthCheck(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	log.Println("💓 插件健康检查请求收到")
	return &pb.HealthResponse{Alive: true}, nil
}

// RegisterToMainService 插件启动时主动注册
func RegisterToMainService(pluginName, pluginAddr string) {
	body := map[string]string{
		"name": pluginName,
		"addr": pluginAddr,
	}
	jsonData, _ := json.Marshal(body)

	resp, err := http.Post("http://127.0.0.1:8080/plugin/register", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Fatalf("❌ 插件注册失败: %v", err)
	}
	defer resp.Body.Close()
	log.Println("🔔 插件注册成功")
}

// CheckPluginHealth 插件心跳检测
func CheckPluginHealth(name, addr string) bool {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return false
	}
	defer conn.Close()

	client := pb.NewPluginServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.HealthCheck(ctx, &pb.HealthRequest{})
	return err == nil && resp.Alive
}
