package lesson

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	lessons := r.Group("/lessons")
	lessons.Get("/:id", h.GetLessonByID)
}

// GetLessonByID 根据 ID 获取课程详情
func (h *Handler) GetLessonByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid lesson id"})
	}

	lesson, err := h.service.GetLessonByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "lesson not found"})
	}

	return c.JSON(lesson)
}
