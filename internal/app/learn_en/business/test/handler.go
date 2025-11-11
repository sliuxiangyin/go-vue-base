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
	r.Post("/test/semantic-chunks", h.AnalyzeSemanticChunks)
	r.Post("/test/transcribe", h.TranscribeAudio)
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
		Model        string `json:"model,omitempty"`
		Voice        string `json:"voice,omitempty"`
		Text         string `json:"text"`
		LanguageType string `json:"language_type,omitempty"`
	}
	var r req = req{
		Model:        "qwen3-tts-flash",
		Voice:        "Jennifer",
		Text:         "Life was like a box of chocolates, you never know what you're gonna get.",
		LanguageType: "English",
	}

	// 默认值

	audioURL, err := h.service.SynthesizeAudio(r.Model, r.Voice, r.Text, r.LanguageType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 返回音频文件相对路径
	return c.JSON(fiber.Map{
		"path":          audioURL,
		"model":         r.Model,
		"voice":         r.Voice,
		"language_type": r.LanguageType,
	})
}

// AnalyzeSemanticChunks 英文语义分析器：将英文句子按意群分割
func (h *Handler) AnalyzeSemanticChunks(c *fiber.Ctx) error {
	type req struct {
		Sentence string `json:"sentence"`
	}
	var r req
	if err := c.BodyParser(&r); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if r.Sentence == "" {
		return c.Status(400).JSON(fiber.Map{"error": "sentence is required"})
	}

	chunks, err := h.service.AnalyzeSemanticChunks(c.Context(), r.Sentence)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"sentence": r.Sentence,
		"chunks":   chunks,
	})
}

// TranscribeAudio 音频转录：将音频文件转写为文字并返回逐词时间戳
func (h *Handler) TranscribeAudio(c *fiber.Ctx) error {
	type req struct {
		AudioPath string `json:"audio_path"`
		Language  string `json:"language,omitempty"`
		ModelSize string `json:"model_size,omitempty"`
	}
	var r req
	if err := c.BodyParser(&r); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if r.AudioPath == "" {
		return c.Status(400).JSON(fiber.Map{"error": "audio_path is required"})
	}

	// 设置默认值
	if r.Language == "" {
		r.Language = "en"
	}
	if r.ModelSize == "" {
		r.ModelSize = "small"
	}

	// 调用转录服务
	result, err := h.service.TranscribeAudio(r.AudioPath, r.Language, r.ModelSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 构造响应
	words := make([]map[string]interface{}, 0, len(result.Words))
	for _, w := range result.Words {
		words = append(words, map[string]interface{}{
			"word":       w.Word,
			"start":      w.Start,
			"end":        w.End,
			"confidence": w.Confidence,
		})
	}

	return c.JSON(fiber.Map{
		"text":     result.Text,
		"words":    words,
		"language": result.Language,
	})
}
