package response

import "github.com/gofiber/fiber/v2"

// Success 成功返回
func (r *Response) Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Fail 失败返回
func (r *Response) Fail(c *fiber.Ctx, err error) error {
	return err // 会被全局 ErrorHandler 捕获
}

// Success Fiber Context 扩展方法
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}
func Fail(c *fiber.Ctx, err error) error {
	return err
}
