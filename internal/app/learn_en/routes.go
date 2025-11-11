package learn_en

import (
	"databaseAi/internal/utils"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func (a *App) RegisterRoutes(app *fiber.App) {
	// 注册 API 路由
	api := app.Group("/api")
	for _, handler := range a.Handlers {
		handler.RegisterRoutes(api)
	}

	// 注册静态文件服务，提供音频文件访问
	storagePath := filepath.Join(utils.ProjectRoot(), "storage")
	app.Static("/storage", storagePath, fiber.Static{
		Browse: false,
		Index:  "",
	})
}
