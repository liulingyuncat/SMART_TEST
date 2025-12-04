package utils

import "github.com/gin-gonic/gin"

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseSuccess 成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ResponseSuccessWithCode 自定义HTTP状态码的成功响应
func ResponseSuccessWithCode(c *gin.Context, httpCode int, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ResponseError 错误响应
func ResponseError(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
	})
}

// SuccessResponse 成功响应别名
func SuccessResponse(c *gin.Context, data interface{}) {
	ResponseSuccess(c, data)
}

// ErrorResponse 错误响应别名
func ErrorResponse(c *gin.Context, httpCode int, message string) {
	ResponseError(c, httpCode, message)
}

// MessageResponse 仅消息响应
func MessageResponse(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
	})
}
