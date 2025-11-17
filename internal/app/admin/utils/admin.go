package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// getUserID 从上下文获取用户ID

func GetUserID(c *fiber.Ctx) uint {
	userID := c.Locals("user_id") // 注意：与中间件保持一致，使用 user_id
	if userID == nil {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	if id, ok := userID.(float64); ok {
		return uint(id)
	}
	// 尝试从字符串解析
	if idStr, ok := userID.(string); ok {
		id, _ := strconv.ParseUint(idStr, 10, 32)
		return uint(id)
	}
	return 0
}
