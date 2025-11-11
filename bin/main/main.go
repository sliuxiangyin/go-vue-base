package main

import (
	"databaseAi/internal/app/admin/migrations"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	fib := fiber.New()
	fib.Use(cors.New())

	// 初始化 Learn English 应用
	learnEnApp, err := InitializeLearnEnService(BuildEnv)
	if err != nil {
		log.Fatalf("Failed to initialize Learn English service: %v", err)
	}
	log.Printf("Learn English service initialized successfully")

	// 初始化 Admin 应用
	adminApp, err := InitializeAdminService(BuildEnv)
	if err != nil {
		log.Fatalf("Failed to initialize Admin service: %v", err)
	}
	log.Printf("Admin service initialized successfully")

	// 运行 Admin 数据库迁移
	if err := migrations.MigrateAdmin(adminApp.DB); err != nil {
		log.Fatalf("Failed to migrate admin database: %v", err)
	}

	// 注册路由
	learnEnApp.RegisterRoutes(fib)
	adminApp.RegisterRoutes(fib)
	setupProductionWeb(fib)

	log.Println("Server starting on port 8080...")
	err = fib.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
