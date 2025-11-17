package event

import (
	"databaseAi/internal/app/admin/business/auth"
	"databaseAi/internal/app/admin/middleware"
	"databaseAi/internal/app/admin/utils"
	sharedEvent "databaseAi/internal/shared/event"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService *auth.Service
}

func NewHandler(authService *auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(router fiber.Router) {
	authMiddleware := middleware.AuthMiddleware(h.authService)
	events := router.Group("/events", authMiddleware)
	events.Get("/stream", h.HandleEvents)
	events.Get("/stats", h.GetStats)
}

// HandleEvents SSE 事件流处理
func (h *Handler) HandleEvents(c *fiber.Ctx) error {
	// 获取当前用户ID
	userID := utils.GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}

	return sharedEvent.HandleSSE(c, userID)
}

// GetStats 获取事件系统统计信息
func (h *Handler) GetStats(c *fiber.Ctx) error {
	manager := sharedEvent.GetEventManager()
	return c.JSON(fiber.Map{
		"code": 0,
		"data": fiber.Map{
			"connected_clients": manager.GetClientCount(),
		},
	})
}
