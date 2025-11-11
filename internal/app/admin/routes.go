package admin

import (
	"github.com/gofiber/fiber/v2"
)

func (a *App) RegisterRoutes(app *fiber.App) {
	// 注册 API 路由
	api := app.Group("/api/admin")
	for _, handler := range a.Handlers {
		handler.RegisterRoutes(api)
	}
}
