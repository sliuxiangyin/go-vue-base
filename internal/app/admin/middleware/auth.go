package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// TokenValidator 定义 token 验证接口
type TokenValidator interface {
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
}

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(validator TokenValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"code":    401,
				"message": "未提供认证令牌",
			})
		}

		// 提取 token (Bearer xxx)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"code":    401,
				"message": "认证令牌格式错误",
			})
		}

		token := parts[1]
		claims, err := validator.ValidateToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"code":    401,
				"message": "认证令牌无效或已过期",
			})
		}

		// 将用户ID存入上下文
		userID := uint((*claims)["user_id"].(float64))
		c.Locals("user_id", userID)
		c.Locals("username", (*claims)["username"])

		return c.Next()
	}
}
