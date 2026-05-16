// Package response provides standardized HTTP response formatting.
package response

// 本文件封装统一的 HTTP 响应输出格式。

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
)

// envelope 表示统一响应结构的基础外层包装。
type envelope struct {
	Data      any            `json:"data,omitempty"`
	Meta      any            `json:"meta,omitempty"`
	Code      string         `json:"code,omitempty"`
	Message   string         `json:"message,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
}

// OK 返回标准成功响应。
func OK(c *gin.Context, data any) {
	JSON(c, http.StatusOK, data, nil)
}

// Created 返回创建成功响应。
func Created(c *gin.Context, data any) {
	JSON(c, http.StatusCreated, data, nil)
}

// JSON 按统一结构返回指定状态码的响应。
func JSON(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, envelope{
		Data:      data,
		Meta:      meta,
		RequestID: c.GetString("request_id"),
	})
}

// NoContent 返回无响应体的成功结果。
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 按统一结构输出错误响应，并将错误挂入请求上下文。
func Error(c *gin.Context, err error) {
	appErr := errs.As(err)
	lang := c.GetString("lang")
	if lang == "" {
		lang = "zh"
	}
	appErr.Localize(lang)
	_ = c.Error(appErr)
	c.JSON(appErr.Status, envelope{
		Code:      string(appErr.Code),
		Message:   appErr.Message,
		Details:   appErr.Details,
		RequestID: c.GetString("request_id"),
	})
}
