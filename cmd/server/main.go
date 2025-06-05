package main

import (
	"fmt"
	"qwqserver/internal/app"
	"qwqserver/internal/config"
)

type qwqLog struct {
}

func (q qwqLog) Log(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	fmt.Println(s)
}

func (q qwqLog) Info(msg string, args ...any) {
	q.Log(msg, args)
}

func (q qwqLog) Error(msg string, args ...any) {
	q.Log(msg, args)
}

func (q qwqLog) Warn(msg string, args ...any) {
	q.Log(msg)
}

func (q qwqLog) Debug(msg string, args ...any) {
	q.Log(msg)
}

func (q qwqLog) Fatal(msg string, args ...any) {
	q.Log(msg, args)
}

func (q qwqLog) Panic(msg string, args ...any) {
	q.Log(msg, args)
}

func main() {
	l := qwqLog{}
	cfg := config.New()
	apps := app.New(&app.Config{
		L:      l,
		Config: cfg,
	})

	defer apps.Close()

	// 启动服务
	apps.Listen()

}

//
//// === 数据结构定义 ===
//type ChatMessage struct {
//	Role    string `json:"role"`
//	Content string `json:"content"`
//}
//
//type ChatRequest struct {
//	Model            string            `json:"model"`
//	Messages         []ChatMessage     `json:"messages"`
//	Stream           bool              `json:"stream"`
//	MaxTokens        int               `json:"max_tokens"`
//	EnableThinking   bool              `json:"enable_thinking"`
//	ThinkingBudget   int               `json:"thinking_budget"`
//	MinP             float64           `json:"min_p"`
//	Stop             interface{}       `json:"stop"` // null
//	Temperature      float64           `json:"temperature"`
//	TopP             float64           `json:"top_p"`
//	TopK             int               `json:"top_k"`
//	FrequencyPenalty float64           `json:"frequency_penalty"`
//	N                int               `json:"n"`
//	ResponseFormat   map[string]string `json:"response_format"`
//	Tools            []Tool            `json:"tools"`
//}
//
//type Tool struct {
//	Type     string       `json:"type"`
//	Function ToolFunction `json:"function"`
//}
//
//type ToolFunction struct {
//	Description string                 `json:"description"`
//	Name        string                 `json:"name"`
//	Parameters  map[string]interface{} `json:"parameters"`
//	Strict      bool                   `json:"strict"`
//}
//
//// === 会话上下文缓存 ===
//var sessionContexts = struct {
//	sync.RWMutex
//	data map[string][]ChatMessage
//}{data: make(map[string][]ChatMessage)}
//
//func main() {
//	r := gin.Default()
//
//	// 主对话接口
//	r.GET("/invoke", func(c *gin.Context) {
//		sessionID := c.Query("session")
//		message := c.Query("message")
//
//		if sessionID == "" || message == "" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "session 和 message 参数不能为空"})
//			return
//		}
//
//		// 获取历史上下文
//		sessionContexts.RLock()
//		history := sessionContexts.data[sessionID]
//		sessionContexts.RUnlock()
//
//		// 添加用户消息
//		history = append(history, ChatMessage{
//			Role:    "user",
//			Content: message,
//		})
//
//		// 构造请求体
//		reqBody := ChatRequest{
//			Model:            "Qwen/QwQ-32B",
//			Messages:         history,
//			Stream:           false,
//			MaxTokens:        512,
//			EnableThinking:   false,
//			ThinkingBudget:   4096,
//			MinP:             0.05,
//			Stop:             nil,
//			Temperature:      0.7,
//			TopP:             0.7,
//			TopK:             50,
//			FrequencyPenalty: 0.5,
//			N:                1,
//			ResponseFormat:   map[string]string{"type": "text"},
//			Tools: []Tool{
//				{
//					Type: "function",
//					Function: ToolFunction{
//						Description: "<string>",
//						Name:        "<string>",
//						Parameters:  map[string]interface{}{},
//						Strict:      false,
//					},
//				},
//			},
//		}
//
//		jsonData, _ := json.Marshal(reqBody)
//
//		req, _ := http.NewRequest("POST", "https://api.siliconflow.cn/v1/chat/completions", bytes.NewBuffer(jsonData))
//
//		req.Header.Set("Authorization", "Bearer sk-lwethrftryyrbkcgbogihomelpbaiuaqvfquxssrnilqjunm")
//		req.Header.Set("Content-Type", "application/json")
//
//		client := &http.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "请求失败"})
//			return
//		}
//		defer resp.Body.Close()
//
//		body, _ := ioutil.ReadAll(resp.Body)
//
//		var respObj map[string]interface{}
//		err = json.Unmarshal(body, &respObj)
//		if err != nil {
//			c.JSON(http.StatusOK, gin.H{"raw": string(body)})
//			return
//		}
//
//		// 提取助手回复
//		reply := ""
//		if choices, ok := respObj["choices"].([]interface{}); ok && len(choices) > 0 {
//			if msg, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
//				reply = fmt.Sprintf("%v", msg["content"])
//				// 添加到上下文
//				history = append(history, ChatMessage{
//					Role:    "assistant",
//					Content: reply,
//				})
//			}
//		}
//
//		// 更新上下文
//		sessionContexts.Lock()
//		sessionContexts.data[sessionID] = history
//		sessionContexts.Unlock()
//
//		// 返回用户提问与模型回复
//		c.JSON(http.StatusOK, gin.H{
//			"session": sessionID,
//			"user":    message,
//			"reply":   reply,
//		})
//	})
//
//	// 清除对话接口
//	r.GET("/clear", func(c *gin.Context) {
//		sessionID := c.Query("session")
//		if sessionID == "" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "session 参数不能为空"})
//			return
//		}
//		sessionContexts.Lock()
//		delete(sessionContexts.data, sessionID)
//		sessionContexts.Unlock()
//
//		c.JSON(http.StatusOK, gin.H{"message": "上下文已清除", "session": sessionID})
//	})
//
//	r.Run(":8080")
//}
