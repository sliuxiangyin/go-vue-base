package test

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
	r.Post("/test", h.SetValue)
	r.Get("/test/audio", h.TestAudio)
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

func (h *Handler) TestAudio(c *fiber.Ctx) error {
	type req struct {
		Model string `json:"model,omitempty"`
		Voice string `json:"voice,omitempty"`
		Text  string `json:"text"`
	}
	var r req = req{
		Model: "cosyvoice-v3",
		Voice: "longxiaochun_v2",
		Text:  "My world is up to me",
	}

	// 默认值

	audioData, err := h.service.SynthesizeAudio(r.Model, r.Voice, r.Text)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 返回音频数据
	c.Set("Content-Type", "audio/mpeg")
	return c.Send(audioData)
}
