package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS

func main() {
	app := fiber.New()

	// 添加 CORS 中间件
	app.Use(cors.New())

	// API 路由示例
	api := app.Group("/api")
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from GoFiber!",
		})
	})

	// 在开发模式下启动前端开发服务器
	if os.Getenv("MODE") == "dev" {
		go startFrontendDevServer()
	}

	// 在生产模式下提供嵌入的前端静态文件
	subFS, err := fs.Sub(frontendFiles, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}
	
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(subFS),
		Browse: true,
	}))

	// 启动服务器
	log.Println("Server starting on port 8080...")
	app.Listen(":8080")
}

// 启动前端开发服务器
func startFrontendDevServer() {
	cmd := exec.Command("npm", "run", "dev")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to start frontend dev server: %v", err)
	}
}