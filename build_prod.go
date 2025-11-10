//go:build !dev
// +build !dev

package main

import (
	"databaseAi/compile"
	"embed"
	"github.com/gofiber/fiber/v2"
)

// 在生产模式下嵌入前端静态文件
//
//go:embed web/dist/*
var webFiles embed.FS
var BuildEnv = "prod"

// setupProductionWeb 在生产模式下设置前端静态文件服务
func setupProductionWeb(app *fiber.App) {
	compile.SetupProd(app, webFiles)
}
