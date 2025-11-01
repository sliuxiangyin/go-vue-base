//go:build dev
// +build dev

package main

import (
	"databaseAi/compile"
	"github.com/gofiber/fiber/v2"
)

// setupProductionWeb 在开发模式下代理前端请求到 Vite 开发服务器
func setupProductionWeb(app *fiber.App) {
	compile.SetupDev(app)
}
