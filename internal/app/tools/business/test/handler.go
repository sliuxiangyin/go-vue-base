package test

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Post("/test", h.SetValue)
	r.Get("/test/:key", h.GetValue)
}

func (h *Handler) SetValue(c *fiber.Ctx) error {
	type req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	var r req
	if err := c.BodyParser(&r); err != nil {
		return err
	}
	h.service.SetValue(r.Key, r.Value)
	return c.JSON(fiber.Map{"msg": "ok"})
}

func (h *Handler) GetValue(c *fiber.Ctx) error {
	key := c.Params("key")
	value, ok := h.service.GetValue(key)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(fiber.Map{"key": key, "value": value})
}
