package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// PermissionChecker 定义权限检查接口
type PermissionChecker interface {
	CheckUserPermission(userID uint, permissionName string) (bool, error)
}

// PermissionMiddleware 权限验证中间件（基于路由名称）
func PermissionMiddleware(checker PermissionChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取当前用户ID
		userID, ok := c.Locals("user_id").(uint)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"code":    401,
				"message": "未认证的用户",
			})
		}

		// 获取路由名称（通过 .Name() 设置的名称）
		routeName := c.Route().Name
		if routeName == "" {
			// 如果路由没有设置名称，则跳过权限检查（或根据需求返回错误）
			return c.Next()
		}

		// 检查用户是否有该路由的权限
		hasPermission, err := checker.CheckUserPermission(userID, routeName)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"code":    500,
				"message": "权限检查失败",
			})
		}

		if !hasPermission {
			return c.Status(403).JSON(fiber.Map{
				"code":    403,
				"message": "没有访问权限",
			})
		}

		return c.Next()
	}
}
