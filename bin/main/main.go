package main

import (
	"databaseAi/internal/app/tools"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	fib := fiber.New()
	fib.Use(cors.New())
	api, err := tools.NewApp(BuildEnv)
	if err != nil {
		return
	}
	api.RegisterRoutes(fib)

	setupProductionWeb(fib)

	log.Println("Server starting on port 8080...")
	err = fib.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
