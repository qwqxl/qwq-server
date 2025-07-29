package qqbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SendMsg 发送消息函数
func SendMsg(apiBaseURL string, groupID int64, text string) {
	url := fmt.Sprintf("%s/send_group_msg", apiBaseURL)

	payload := map[string]interface{}{
		"group_id": groupID,
		"message": []map[string]interface{}{
			{
				"type": "text",
				"data": map[string]string{
					"text": text,
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Println("请求创建失败:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("请求发送失败:", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("响应:", string(respBody))
}
