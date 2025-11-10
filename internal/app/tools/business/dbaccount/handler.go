package dbaccount

import (
	api2 "databaseAi/internal/app/tools/infra/response/api"
	"databaseAi/internal/app/tools/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler 数据库账号处理器
type Handler struct {
	service *Service
}

// NewHandler 创建新的数据库账号处理器实例
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(app fiber.Router) {
	route := app.Group("/db-account")

	route.Post("/", h.Create)
	route.Get("/:id", h.GetByID)
	route.Get("/", h.GetAll)
	route.Put("/:id", h.Update)
	route.Delete("/:id", h.Delete)
}

// Create 创建数据库账号
// @Summary 创建数据库账号
// @Description 创建一个新的数据库账号
// @Tags 数据库账号
// @Accept json
// @Produce json
// @Param account body models.DBAccount true "数据库账号信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/db-account [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var account DBAccountReq
	if err := c.BodyParser(&account); err != nil {
		return api2.Error(c, 400, "请求参数解析失败")
	}

	if err := h.service.Create(account.ToModel()); err != nil {
		return api2.Error(c, 400, err.Error())
	}

	return api2.Success(c, account, "数据库账号创建成功")
}

// GetByID 根据ID获取数据库账号
// @Summary 获取数据库账号
// @Description 根据ID获取数据库账号信息
// @Tags 数据库账号
// @Accept json
// @Produce json
// @Param id path int true "数据库账号ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/db-account/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return api2.Error(c, 400, "无效的ID参数")
	}

	account, err := h.service.GetByID(uint(id))
	if err != nil {
		return api2.Error(c, 404, err.Error())
	}

	return api2.Success(c, account, "获取数据库账号成功")
}

// GetAll 获取所有数据库账号
// @Summary 获取所有数据库账号
// @Description 获取所有数据库账号信息
// @Tags 数据库账号
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/db-account [get]
func (h *Handler) GetAll(c *fiber.Ctx) error {
	accounts, err := h.service.GetAll()
	if err != nil {
		return api2.Error(c, 500, "获取数据库账号列表失败")
	}

	return api2.Success(c, accounts, "获取数据库账号列表成功")
}

// Update 更新数据库账号
// @Summary 更新数据库账号
// @Description 更新数据库账号信息
// @Tags 数据库账号
// @Accept json
// @Produce json
// @Param id path int true "数据库账号ID"
// @Param account body models.DBAccount true "数据库账号信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/db-account/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return api2.Error(c, 400, "无效的ID参数")
	}

	var account models.DBAccount
	if err := c.BodyParser(&account); err != nil {
		return api2.Error(c, 400, "请求参数解析失败")
	}

	account.ID = uint(id)

	if err := h.service.Update(&account); err != nil {
		return api2.Error(c, 400, err.Error())
	}

	return api2.Success(c, account, "数据库账号更新成功")
}

// Delete 删除数据库账号
// @Summary 删除数据库账号
// @Description 删除指定ID的数据库账号
// @Tags 数据库账号
// @Accept json
// @Produce json
// @Param id path int true "数据库账号ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/db-account/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return api2.Error(c, 400, "无效的ID参数")
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return api2.Error(c, 404, err.Error())
	}

	return api2.Success(c, nil, "数据库账号删除成功")
}
