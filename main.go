package main

import (
	"databaseAi/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	fib := fiber.New()
	fib.Use(cors.New())
	setupProductionWeb(fib)
	api, err := app.NewApp()
	if err != nil {
		return
	}
	api.RegisterRoutes(fib)
	log.Println("Server starting on port 8080...")
	err = fib.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
