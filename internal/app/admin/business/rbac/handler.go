package rbac

import (
	"databaseAi/internal/app/admin/middleware"
	"databaseAi/internal/app/admin/models"
	"strconv"

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

	// 角色管理路由
	roles := r.Group("/roles", authMiddleware, permMiddleware)
	roles.Post("/", h.CreateRole).Name("role.create")
	roles.Put("/:id", h.UpdateRole).Name("role.update")
	roles.Delete("/:id", h.DeleteRole).Name("role.delete")
	roles.Get("/:id", h.GetRole).Name("role.get")
	roles.Get("/", h.ListRoles).Name("role.list")

	// 权限管理路由
	permissions := r.Group("/permissions", authMiddleware, permMiddleware)
	permissions.Post("/", h.CreatePermission).Name("permission.create")
	permissions.Put("/:id", h.UpdatePermission).Name("permission.update")
	permissions.Delete("/:id", h.DeletePermission).Name("permission.delete")
	permissions.Get("/:id", h.GetPermission).Name("permission.get")
	permissions.Get("/", h.ListPermissions).Name("permission.list")
	permissions.Get("/all", h.ListAllPermissions).Name("permission.list.all")

	// 角色权限关联路由
	roles.Post("/:id/permissions", h.AssignPermissionsToRole).Name("role.assign.permissions")
	roles.Get("/:id/permissions", h.GetRolePermissions).Name("role.permissions")

	// 用户角色关联路由
	users := r.Group("/users", authMiddleware, permMiddleware)
	users.Post("/", h.CreateAdminUser).Name("admin.user.create")
	users.Put("/:id", h.UpdateAdminUser).Name("admin.user.update")
	users.Put("/:id/password", h.UpdateAdminUserPassword).Name("admin.user.password.update")
	users.Delete("/:id", h.DeleteAdminUser).Name("admin.user.delete")
	users.Get("/:id", h.GetAdminUser).Name("admin.user.get")
	users.Get("/", h.ListAdminUsers).Name("admin.user.list")
	users.Post("/:id/roles", h.AssignRolesToUser).Name("user.assign.roles")
	users.Get("/:id/roles", h.GetUserRoles).Name("user.roles")
	users.Get("/:id/permissions", h.GetUserPermissions).Name("user.permissions")
	users.Get("/:id/permissions/frontend", h.GetUserFrontendPermissions).Name("user.permissions.frontend")
}

// ========== 角色管理接口 ==========

// CreateRole 创建角色
func (h *Handler) CreateRole(c *fiber.Ctx) error {
	type CreateRoleRequest struct {
		Name        string `json:"name" validate:"required,min=2,max=50"`
		DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
		Description string `json:"description" validate:"max=255"`
		Status      int8   `json:"status" validate:"oneof=0 1"`
	}

	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	role, err := h.service.CreateRole(req.Name, req.DisplayName, req.Description, req.Status)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "创建成功", "data": role})
}

// UpdateRole 更新角色
func (h *Handler) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的角色ID"})
	}

	type UpdateRoleRequest struct {
		DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
		Description string `json:"description" validate:"max=255"`
		Status      int8   `json:"status" validate:"oneof=0 1"`
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	role, err := h.service.UpdateRole(uint(id), req.DisplayName, req.Description, req.Status)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "更新成功", "data": role})
}

// DeleteRole 删除角色
func (h *Handler) DeleteRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的角色ID"})
	}

	if err := h.service.DeleteRole(uint(id)); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "删除成功"})
}

// GetRole 获取角色详情
func (h *Handler) GetRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的角色ID"})
	}

	role, err := h.service.GetRole(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"code": 404, "message": "角色不存在"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": role})
}

// ListRoles 获取角色列表
func (h *Handler) ListRoles(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	name := c.Query("name", "")

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			statusVal := int8(s)
			status = &statusVal
		}
	}

	roles, total, err := h.service.ListRoles(page, pageSize, name, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取列表失败"})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"list":      roles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ========== 权限管理接口 ==========

// CreatePermission 创建权限
func (h *Handler) CreatePermission(c *fiber.Ctx) error {
	type CreatePermissionRequest struct {
		Name        string                `json:"name" validate:"required,min=2,max=100"`
		DisplayName string                `json:"display_name" validate:"required,min=2,max=100"`
		Description string                `json:"description" validate:"max=255"`
		Category    string                `json:"category" validate:"max=50"`
		Type        models.PermissionType `json:"type" validate:"required,oneof=backend frontend"`
		Path        string                `json:"path" validate:"max=255"`
		Icon        string                `json:"icon" validate:"max=50"`
		ParentID    *uint                 `json:"parent_id"`
		Sort        int                   `json:"sort"`
		Status      int8                  `json:"status" validate:"oneof=0 1"`
	}

	var req CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	permission, err := h.service.CreatePermission(
		req.Name, req.DisplayName, req.Description, req.Category,
		req.Type, req.Path, req.Icon, req.ParentID, req.Sort, req.Status,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "创建成功", "data": permission})
}

// UpdatePermission 更新权限
func (h *Handler) UpdatePermission(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的权限ID"})
	}

	type UpdatePermissionRequest struct {
		DisplayName string                `json:"display_name" validate:"required,min=2,max=100"`
		Description string                `json:"description" validate:"max=255"`
		Category    string                `json:"category" validate:"max=50"`
		Type        models.PermissionType `json:"type" validate:"required,oneof=backend frontend"`
		Path        string                `json:"path" validate:"max=255"`
		Icon        string                `json:"icon" validate:"max=50"`
		ParentID    *uint                 `json:"parent_id"`
		Sort        int                   `json:"sort"`
		Status      int8                  `json:"status" validate:"oneof=0 1"`
	}

	var req UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	permission, err := h.service.UpdatePermission(
		uint(id), req.DisplayName, req.Description, req.Category,
		req.Type, req.Path, req.Icon, req.ParentID, req.Sort, req.Status,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "更新成功", "data": permission})
}

// DeletePermission 删除权限
func (h *Handler) DeletePermission(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的权限ID"})
	}

	if err := h.service.DeletePermission(uint(id)); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "删除成功"})
}

// GetPermission 获取权限详情
func (h *Handler) GetPermission(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的权限ID"})
	}

	permission, err := h.service.GetPermission(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"code": 404, "message": "权限不存在"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": permission})
}

// ListPermissions 获取权限列表
func (h *Handler) ListPermissions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	category := c.Query("category", "")

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			statusVal := int8(s)
			status = &statusVal
		}
	}

	var permType *models.PermissionType
	if typeStr := c.Query("type"); typeStr != "" {
		pt := models.PermissionType(typeStr)
		permType = &pt
	}

	permissions, total, err := h.service.ListPermissions(page, pageSize, category, permType, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取列表失败"})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"list":      permissions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ListAllPermissions 获取所有权限（不分页）
func (h *Handler) ListAllPermissions(c *fiber.Ctx) error {
	permissions, err := h.service.ListAllPermissions()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取权限列表失败"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": permissions})
}

// ========== 角色权限关联接口 ==========

// AssignPermissionsToRole 为角色分配权限
func (h *Handler) AssignPermissionsToRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的角色ID"})
	}

	type AssignPermissionsRequest struct {
		PermissionIDs []uint `json:"permission_ids" validate:"required"`
	}

	var req AssignPermissionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	if err := h.service.AssignPermissionsToRole(uint(id), req.PermissionIDs); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "分配成功"})
}

// GetRolePermissions 获取角色的权限列表
func (h *Handler) GetRolePermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的角色ID"})
	}

	permissions, err := h.service.GetRolePermissions(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取权限列表失败"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": permissions})
}

// ========== 用户角色关联接口 ==========

// AssignRolesToUser 为用户分配角色
func (h *Handler) AssignRolesToUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	type AssignRolesRequest struct {
		RoleIDs []uint `json:"role_ids" validate:"required"`
	}

	var req AssignRolesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	if err := h.service.AssignRolesToUser(uint(id), req.RoleIDs); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "分配成功"})
}

// GetUserRoles 获取用户的角色列表
func (h *Handler) GetUserRoles(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	roles, err := h.service.GetUserRoles(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取角色列表失败"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": roles})
}

// GetUserPermissions 获取用户的权限列表
func (h *Handler) GetUserPermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	permissions, err := h.service.GetUserPermissions(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取权限列表失败"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": permissions})
}

// GetUserFrontendPermissions 获取用户的前端权限列表(用于菜单渲染)
func (h *Handler) GetUserFrontendPermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	permissions, err := h.service.GetUserFrontendPermissions(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取菜单权限列表失败"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": permissions})
}

// ========== 管理员管理接口 ==========

// CreateAdminUser 创建管理员
func (h *Handler) CreateAdminUser(c *fiber.Ctx) error {
	type CreateAdminUserRequest struct {
		Username string `json:"username" validate:"required,min=2,max=50"`
		Password string `json:"password" validate:"required,min=6,max=50"`
		Email    string `json:"email" validate:"omitempty,email,max=100"`
		Nickname string `json:"nickname" validate:"max=50"`
		Avatar   string `json:"avatar" validate:"max=255"`
		Status   int8   `json:"status" validate:"oneof=0 1"`
	}

	var req CreateAdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	user, err := h.service.CreateAdminUser(req.Username, req.Password, req.Email, req.Nickname, req.Avatar, req.Status)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "创建成功", "data": user})
}

// UpdateAdminUser 更新管理员
func (h *Handler) UpdateAdminUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	type UpdateAdminUserRequest struct {
		Email    string `json:"email" validate:"omitempty,email,max=100"`
		Nickname string `json:"nickname" validate:"max=50"`
		Avatar   string `json:"avatar" validate:"max=255"`
		Status   int8   `json:"status" validate:"oneof=0 1"`
	}

	var req UpdateAdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	user, err := h.service.UpdateAdminUser(uint(id), req.Email, req.Nickname, req.Avatar, req.Status)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "更新成功", "data": user})
}

// UpdateAdminUserPassword 更新管理员密码
func (h *Handler) UpdateAdminUserPassword(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	type UpdatePasswordRequest struct {
		Password string `json:"password" validate:"required,min=6,max=50"`
	}

	var req UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "请求参数错误"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "参数验证失败", "errors": err.Error()})
	}

	if err := h.service.UpdateAdminUserPassword(uint(id), req.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "密码更新成功"})
}

// DeleteAdminUser 删除管理员
func (h *Handler) DeleteAdminUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	if err := h.service.DeleteAdminUser(uint(id)); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "删除成功"})
}

// GetAdminUser 获取管理员详情
func (h *Handler) GetAdminUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 400, "message": "无效的用户ID"})
	}

	user, err := h.service.GetAdminUser(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"code": 404, "message": "管理员不存在"})
	}

	return c.JSON(fiber.Map{"code": 0, "message": "success", "data": user})
}

// ListAdminUsers 获取管理员列表
func (h *Handler) ListAdminUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	keyword := c.Query("keyword", "")

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			statusVal := int8(s)
			status = &statusVal
		}
	}

	users, total, err := h.service.ListAdminUsers(page, pageSize, keyword, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": 500, "message": "获取列表失败"})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"list":      users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
