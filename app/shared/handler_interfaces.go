package shared

import "github.com/gofiber/fiber/v2"

type HandlerInterfaces interface {
	RegisterRoutes(r fiber.Router)
}
