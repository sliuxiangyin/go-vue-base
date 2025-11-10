package database

import (
	"databaseAi/internal/app/tools/infra/response/api"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) RegisterRoutes(r fiber.Router) {
	route := r.Group("/database")
	route.Get("/init/:id", h.Init)
}

func (h *Handler) Init(c *fiber.Ctx) error {
	id := cast.ToInt(c.Params("id"))
	err := h.service.Add(uint(id))
	if err != nil {
		return api.Error(c, 400, "请求参数解析失败")
	}
	tables, err := h.service.GetTables(uint(id))
	if err != nil {
		return err
	}
	return api.Success(c, tables, "ok")
}

func (h *Handler) Leave(c *fiber.Ctx) error {
	id := cast.ToInt(c.Params("id"))
	err := h.service.Remove(uint(id))
	if err != nil {
		return api.Error(c, 400, "请求参数解析失败")
	}
	return api.Success(c, nil, "ok")
}
