package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	fib := fiber.New()
	fib.Use(cors.New())

	// 初始化 Learn English 应用
	app, err := InitializeLearnEnService(BuildEnv)
	if err != nil {
		log.Fatalf("Failed to initialize Learn English service: %v", err)
	}
	log.Printf("Learn English service initialized successfully")

	app.RegisterRoutes(fib)
	setupProductionWeb(fib)

	log.Println("Server starting on port 8080...")
	err = fib.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
