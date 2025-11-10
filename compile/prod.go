package compile

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"io/fs"
	"log"
	"net/http"
)

func SetupProd(app *fiber.App, webFiles embed.FS) {
	// 检查是否有嵌入的前端文件
	_, err := webFiles.ReadDir("web/dist")
	if err != nil {
		// 如果没有嵌入的前端文件，提供一个简单的提示页面
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("web files not found. Please build the web first.")
		})
		return
	}

	// 提供嵌入的前端静态文件
	subFS, err := fs.Sub(webFiles, "web/dist")
	if err != nil {
		log.Printf("Failed to create sub filesystem: %v", err)
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Failed to serve web files.")
		})
		return
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(subFS),
		Browse: true,
	}))
}
