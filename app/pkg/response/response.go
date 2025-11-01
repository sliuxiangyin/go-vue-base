package response

import (
	"github.com/gofiber/fiber/v2"
)

// 全局响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 自定义业务错误
type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string {
	return e.Msg
}

// 全局 ErrorHandler，注册到 Fiber Config
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := 500
	msg := "Internal Server Error"

	switch e := err.(type) {
	case *fiber.Error:
		code = e.Code
		msg = e.Message
	case *BizError:
		code = e.Code
		msg = e.Msg
	default:
		msg = e.Error()
	}

	return c.Status(code).JSON(Response{
		Code: code,
		Msg:  msg,
	})
}
