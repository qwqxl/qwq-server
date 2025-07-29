package apiresponse

type Response struct {
	Code    BaseCode    `json:"code"`           // 模块化业务码，如 U1001
	Message string      `json:"message"`        // 默认提示信息，可自定义
	Data    interface{} `json:"data,omitempty"` // 可选数据内容
}
