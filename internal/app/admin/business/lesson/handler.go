package lesson

import (
	"databaseAi/internal/app/admin/middleware"
	"databaseAi/internal/shared/models"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service       *Service
	validate      *validator.Validate
	authValidator middleware.TokenValidator
}

func NewHandler(service *Service, authValidator middleware.TokenValidator) *Handler {
	return &Handler{
		service:       service,
		validate:      validator.New(),
		authValidator: authValidator,
	}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	authMiddleware := middleware.AuthMiddleware(h.authValidator)
	permMiddleware := middleware.PermissionMiddleware(h.service)

	// 课程管理路由
	lessons := r.Group("/lessons", authMiddleware, permMiddleware)
	lessons.Post("/", h.CreateLesson).Name("lesson.create")
	lessons.Put("/:id", h.UpdateLesson).Name("lesson.update")
	lessons.Delete("/:id", h.DeleteLesson).Name("lesson.delete")
	lessons.Get("/:id", h.GetLesson).Name("lesson.get")
	lessons.Get("/", h.ListLessons).Name("lesson.list")
	lessons.Get("/search/tags", h.SearchLessonsByTags).Name("lesson.search.tags")
}

// CreateLesson 创建课程
func (h *Handler) CreateLesson(c *fiber.Ctx) error {
	type CreateLessonRequest struct {
		Title             string                `json:"title" validate:"required,min=1,max=255"`
		AudioURL          string                `json:"audio_url" validate:"max=512"`
		Duration          float64               `json:"duration"`
		ContentEN         string                `json:"content_en" validate:"required"`
		ContentZH         string                `json:"content_zh"`
		SemanticJSON      models.SemanticChunks `json:"semantic_json"`
		WordTimestampJSON models.WordTimestamps `json:"word_timestamp_json"`
		PhoneticJSON      models.Phonetics      `json:"phonetic_json"`
		Tags              models.Tags           `json:"tags"`
		Level             int8                  `json:"level" validate:"min=1,max=5"`
		IsPublic          bool                  `json:"is_public"`
	}

	var req CreateLessonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	lesson := &models.EnglishLesson{
		Title:             req.Title,
		AudioURL:          req.AudioURL,
		Duration:          req.Duration,
		ContentEN:         req.ContentEN,
		ContentZH:         req.ContentZH,
		SemanticJSON:      req.SemanticJSON,
		WordTimestampJSON: req.WordTimestampJSON,
		PhoneticJSON:      req.PhoneticJSON,
		Tags:              req.Tags,
		Level:             req.Level,
		IsPublic:          req.IsPublic,
	}

	if err := h.service.CreateLesson(lesson); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "创建成功", "data": lesson})
}

// UpdateLesson 更新课程
func (h *Handler) UpdateLesson(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的课程ID"})
	}

	type UpdateLessonRequest struct {
		Title             string                `json:"title" validate:"required,min=1,max=255"`
		AudioURL          string                `json:"audio_url" validate:"max=512"`
		Duration          float64               `json:"duration"`
		ContentEN         string                `json:"content_en" validate:"required"`
		ContentZH         string                `json:"content_zh"`
		SemanticJSON      models.SemanticChunks `json:"semantic_json"`
		WordTimestampJSON models.WordTimestamps `json:"word_timestamp_json"`
		PhoneticJSON      models.Phonetics      `json:"phonetic_json"`
		Tags              models.Tags           `json:"tags"`
		Level             int8                  `json:"level" validate:"min=1,max=5"`
		IsPublic          bool                  `json:"is_public"`
	}

	var req UpdateLessonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	lesson := &models.EnglishLesson{
		ID:                uint(id),
		Title:             req.Title,
		AudioURL:          req.AudioURL,
		Duration:          req.Duration,
		ContentEN:         req.ContentEN,
		ContentZH:         req.ContentZH,
		SemanticJSON:      req.SemanticJSON,
		WordTimestampJSON: req.WordTimestampJSON,
		PhoneticJSON:      req.PhoneticJSON,
		Tags:              req.Tags,
		Level:             req.Level,
		IsPublic:          req.IsPublic,
	}

	if err := h.service.UpdateLesson(lesson); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "更新成功", "data": lesson})
}

// DeleteLesson 删除课程
func (h *Handler) DeleteLesson(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的课程ID"})
	}

	if err := h.service.DeleteLesson(uint(id)); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "删除成功"})
}

// GetLesson 获取课程详情
func (h *Handler) GetLesson(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的课程ID"})
	}

	lesson, err := h.service.GetLesson(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"code": 404, "message": "课程不存在"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": lesson})
}

// ListLessons 获取课程列表
func (h *Handler) ListLessons(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	var level *int8
	if levelStr := c.Query("level"); levelStr != "" {
		if l, err := strconv.ParseInt(levelStr, 10, 8); err == nil {
			levelVal := int8(l)
			level = &levelVal
		}
	}

	var isPublic *bool
	if isPublicStr := c.Query("is_public"); isPublicStr != "" {
		if p, err := strconv.ParseBool(isPublicStr); err == nil {
			isPublic = &p
		}
	}

	lessons, total, err := h.service.ListLessons(page, pageSize, level, isPublic)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取列表失败"})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"list":      lessons,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// SearchLessonsByTags 根据标签搜索课程
func (h *Handler) SearchLessonsByTags(c *fiber.Ctx) error {
	tagsStr := c.Query("tags", "")
	if tagsStr == "" {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "标签参数不能为空"})
	}

	// 按逗号分隔标签
	tags := strings.Split(tagsStr, ",")

	lessons, err := h.service.SearchLessonsByTags(tags)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "搜索失败"})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"list": lessons,
		},
	})
}
