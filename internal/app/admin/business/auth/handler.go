package auth

import (
	"databaseAi/internal/app/admin/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	authMiddleware := middleware.AuthMiddleware(h.service)
	r.Post("/auth/login", h.Login)
	r.Post("/auth/logout", h.Logout)
	r.Get("/auth/info", authMiddleware, h.GetUserInfo).Name("/auth/info")
	r.Post("/auth/change-password", authMiddleware, h.ChangePassword).Name("/auth/change-password")
}

// Login 登录接口
func (h *Handler) Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Password string `json:"password" validate:"required,min=6,max=50"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "请求参数错误",
		})
	}

	// 使用 validator 验证
	if err := h.validate.Struct(req); err != nil {
		// 解析验证错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch field {
				case "Username":
					switch e.Tag() {
					case "required":
						errorMessages["username"] = "用户名不能为空"
					case "min":
						errorMessages["username"] = "用户名至少需3个字符"
					case "max":
						errorMessages["username"] = "用户名最多50个字符"
					}
				case "Password":
					switch e.Tag() {
					case "required":
						errorMessages["password"] = "密码不能为空"
					case "min":
						errorMessages["password"] = "密码至少需6个字符"
					case "max":
						errorMessages["password"] = "密码最多50个字符"
					}
				}
			}
			return c.Status(400).JSON(fiber.Map{
				"code":    400,
				"message": "参数验证失败",
				"errors":  errorMessages,
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
		})
	}

	token, user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"code":    401,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "登录成功",
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"nickname": user.Nickname,
				"avatar":   user.Avatar,
			},
		},
	})
}

// Logout 登出接口
func (h *Handler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code":    0,
		"message": "登出成功",
	})
}

// GetUserInfo 获取用户信息
func (h *Handler) GetUserInfo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	user, err := h.service.GetUserInfo(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "用户不存在",
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"last_login": user.LastLogin,
		},
	})
}

// ChangePassword 修改密码
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordRequest struct {
		OldPassword string `json:"old_password" validate:"required,min=6"`
		NewPassword string `json:"new_password" validate:"required,min=6,max=50"`
	}

	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "请求参数错误",
		})
	}

	// 使用 validator 验证
	if err := h.validate.Struct(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch field {
				case "OldPassword":
					if e.Tag() == "required" {
						errorMessages["old_password"] = "旧密码不能为空"
					} else if e.Tag() == "min" {
						errorMessages["old_password"] = "旧密码至少需6个字符"
					}
				case "NewPassword":
					if e.Tag() == "required" {
						errorMessages["new_password"] = "新密码不能为空"
					} else if e.Tag() == "min" {
						errorMessages["new_password"] = "新密码至少需6个字符"
					} else if e.Tag() == "max" {
						errorMessages["new_password"] = "新密码最多50个字符"
					}
				}
			}
			return c.Status(400).JSON(fiber.Map{
				"code":    400,
				"message": "参数验证失败",
				"errors":  errorMessages,
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
		})
	}

	userID := c.Locals("user_id").(uint)
	if err := h.service.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "密码修改成功",
	})
}
