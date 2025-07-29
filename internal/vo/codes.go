// vo/codes.go
package vo

const (
	StatusSuccess       = 0           // 请求成功
	StatusInvalidParams = 1000 + iota // 参数错误
	StatusUnauthorized                // 未授权
	StatusInternalError               // 内部错误
	StatusNotFound                    // 资源未找到
	// 更多业务码...
)

const (
	R404 = iota
	R500
)
