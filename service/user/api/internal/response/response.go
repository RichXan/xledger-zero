package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`        // 请求是否成功
	Code    int64       `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 业务数据（可选）
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Response{
		Success: true,
		Code:    200,
		Message: "Success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(w http.ResponseWriter, message string, data interface{}) {
	httpx.OkJson(w, Response{
		Success: true,
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(w http.ResponseWriter, code int64, message string) {
	httpx.OkJson(w, Response{
		Success: false,
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BusinessError 业务错误响应（从 RPC 响应转换）
func BusinessError(w http.ResponseWriter, code int64, message string) {
	httpx.OkJson(w, Response{
		Success: code == 200,
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// FromRPC 从 RPC 响应转换为统一格式
func FromRPC(w http.ResponseWriter, code int64, message string, data interface{}) {
	httpx.OkJson(w, Response{
		Success: code == 200,
		Code:    code,
		Message: message,
		Data:    data,
	})
}
