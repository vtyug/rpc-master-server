package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构
type Response struct {
	Code  int         `json:"code"`            // 业务码
	Msg   string      `json:"msg"`             // 提示信息
	Data  interface{} `json:"data,omitempty"`  // 数据
	Error string      `json:"error,omitempty"` // 错误信息，用于调试
}

// Result 统一返回结果
type Result struct {
	ctx *gin.Context
}

// NewResult 创建返回结果实例
func NewResult(ctx *gin.Context) *Result {
	return &Result{ctx: ctx}
}

// Success 成功返回
func (r *Result) Success(data interface{}) {
	r.ctx.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  GetMessage(SUCCESS),
		Data: data,
	})
}

// Fail 失败返回
func (r *Result) Fail(code int) {
	r.ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  GetMessage(code),
	})
}

// FailWithMsg 失败返回带自定义消息
func (r *Result) FailWithMsg(code int, msg string) {
	r.ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// FailWithError 失败返回带错误信息
func (r *Result) FailWithError(code int, err string) {
	r.ctx.JSON(http.StatusOK, Response{
		Code:  code,
		Msg:   GetMessage(code),
		Error: err,
	})
}
