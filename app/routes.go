package app

import "github.com/gofiber/fiber/v2"

func (a *App) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	for _, handler := range a.Handlers {
		handler.RegisterRoutes(api)
	}
}
