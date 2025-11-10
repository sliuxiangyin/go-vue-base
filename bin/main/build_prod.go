//go:build !dev
// +build !dev

package main

import (
	"databaseAi"
	"databaseAi/bin/main/compile"

	"github.com/gofiber/fiber/v2"
)

var BuildEnv = "prod"

// setupProductionWeb 在生产模式下设置前端静态文件服务
func setupProductionWeb(app *fiber.App) {
	compile.SetupProd(app, databaseAi.WebEmbeds)
}
