package api

import (
	"github.com/gofiber/fiber/v2"
)

// 全局响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}, msg string) error {
	return c.JSON(Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(Response{
		Code: code,
		Msg:  msg,
	})
}
