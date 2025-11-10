package compile

import (
	backend "databaseAi"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"io/fs"
	"log"
	"net/http"
)

func SetupProd(app *fiber.App, webEmbeds map[string]backend.EmbedsInfo) {

	for route, embed := range webEmbeds {
		var subFS fs.FS
		// 检查是否有嵌入的前端文件
		_, err := embed.Fs.ReadDir(embed.Path)
		if err != nil {
			// 如果没有嵌入的前端文件,提供一个简单的提示页面
			app.Get(route, func(c *fiber.Ctx) error {
				return c.SendString("web files not found. Please build the web first.")
			})
			return
		}

		// 提供嵌入的前端静态文件
		subFS, err = fs.Sub(embed.Fs, embed.Path)
		if err != nil {
			log.Printf("Failed to create sub filesystem: %v", err)
			app.Get(route, func(c *fiber.Ctx) error {
				return c.SendString("Failed to serve web files.")
			})
			return
		}
		app.Use(route, filesystem.New(filesystem.Config{
			Root:         http.FS(subFS),
			Browse:       false, // 建议关闭目录浏览
			Index:        "index.html",
			PathPrefix:   "",
			NotFoundFile: "index.html",
		}))

	}

}
