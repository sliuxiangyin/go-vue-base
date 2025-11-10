package relationship

import (
	"bufio"
	"databaseAi/internal/app/tools/infra/response/api"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"strconv"
)

// Handler 关系处理器
type Handler struct {
	service *Service
}

// NewHandler 创建新的关系处理器实例
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(app fiber.Router) {
	route := app.Group("/relationship")

	route.Post("/", h.Create)
	route.Get("/:id", h.GetByID)
	route.Get("/account/:accountID", h.GetByAccountID)
	route.Get("/", h.GetAll)
	route.Put("/:id", h.Update)
	route.Delete("/:id", h.Delete)
	route.Get("/:id/analyze", h.Analyze)
}

// Create 创建关系记录
// @Summary 创建关系记录
// @Description 创建一个新的关系记录
// @Tags 关系记录
// @Accept json
// @Produce json
// @Param relationship body models.Relationship true "关系记录信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/relationship [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var relationship RelationshipReq
	if err := c.BodyParser(&relationship); err != nil {
		return api.Error(c, 400, "请求参数解析失败")
	}
	model, err := h.service.Create(relationship.ToModel())
	if err != nil {
		return api.Error(c, 400, err.Error())
	}
	return api.Success(c, model, "关系记录创建成功")
}

// GetByID 根据ID获取关系记录
// @Summary 获取关系记录
// @Description 根据ID获取关系记录信息
// @Tags 关系记录
// @Accept json
// @Produce json
// @Param id path int true "关系记录ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/relationship/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return api.Error(c, 400, "无效的ID参数")
	}

	relationship, err := h.service.GetByID(uint(id))
	if err != nil {
		return api.Error(c, 404, err.Error())
	}

	return api.Success(c, relationship, "获取关系记录成功")
}

// GetByAccountID 根据账号ID获取关系记录列表
// @Summary 获取账号的关系记录列表
// @Description 根据账号ID获取关系记录列表
// @Tags 关系记录
// @Accept json
// @Produce json
// @Param accountID path int true "账号ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/relationship/account/{accountID} [get]
func (h *Handler) GetByAccountID(c *fiber.Ctx) error {
	accountID, err := strconv.ParseUint(c.Params("accountID"), 10, 32)
	if err != nil {
		return api.Error(c, 400, "无效的账号ID参数")
	}

	relationships, err := h.service.GetByAccountID(uint(accountID))
	if err != nil {
		return api.Error(c, 404, err.Error())
	}

	return api.Success(c, relationships, "获取关系记录列表成功")
}

// GetAll 获取所有关系记录
// @Summary 获取所有关系记录
// @Description 获取所有关系记录信息
// @Tags 关系记录
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/relationship [get]
func (h *Handler) GetAll(c *fiber.Ctx) error {
	relationships, err := h.service.GetAll()
	if err != nil {
		return api.Error(c, 500, "获取关系记录列表失败")
	}

	return api.Success(c, relationships, "获取关系记录列表成功")
}

// Update 更新关系记录
// @Summary 更新关系记录
// @Description 更新关系记录信息
// @Tags 关系记录
// @Accept json
// @Produce json
// @Param id path int true "关系记录ID"
// @Param relationship body models.Relationship true "关系记录信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/relationship/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return api.Error(c, 400, "无效的ID参数")
	}

	var relationship RelationshipReq
	if err := c.BodyParser(&relationship); err != nil {
		return api.Error(c, 400, "请求参数解析失败")
	}

	relationship.ID = uint(id)

	if err := h.service.Update(relationship.ToModel()); err != nil {
		return api.Error(c, 400, err.Error())
	}

	return api.Success(c, relationship, "关系记录更新成功")
}

// Delete 删除关系记录
// @Summary 删除关系记录
// @Description 删除指定ID的关系记录
// @Tags 关系记录
// @Accept json
// @Produce json
// @Param id path int true "关系记录ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/relationship/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return api.Error(c, 400, "无效的ID参数")
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return api.Error(c, 404, err.Error())
	}

	return api.Success(c, nil, "关系记录删除成功")
}

func (h *Handler) Analyze(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return api.Error(c, 400, "无效的ID参数")
	}
	// 设置 SSE 响应头
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	var msgChan = make(chan api.Response, 1)

	go func(msgChan <-chan api.Response) {
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			for text := range msgChan {
				marshal, _ := json.Marshal(text)
				msg := fmt.Sprintf("data: %s\n\n", marshal)
				(*w).Write([]byte(msg))
			}
			msg := fmt.Sprintf("event: %s\n\n", "done")
			(*w).Write([]byte(msg))
		})
	}(msgChan)

	_ = h.service.Analyze(uint(id), msgChan)
	close(msgChan)
	return nil

}
