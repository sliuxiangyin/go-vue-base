package migrations

import (
	"databaseAi/internal/app/admin/models"
	"databaseAi/internal/infra/database"
	sharedModels "databaseAi/internal/shared/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// createDefaultRoles 创建默认角色和权限
func createDefaultRoles(db *database.DB, adminUserID uint) error {
	log.Println("[Migration] Creating default roles and permissions...")

	// 创建超级管理员角色
	superAdminRole := &models.Role{
		Name:        "super_admin",
		DisplayName: "超级管理员",
		Description: "拥有所有权限的系统管理员",
		Status:      1,
		IsSystem:    true,
	}
	if err := db.Create(superAdminRole).Error; err != nil {
		return err
	}

	// 创建默认权限
	defaultPermissions := []models.Permission{
		// 认证相关 - 后端API
		{Name: "/auth/info", DisplayName: "获取用户信息", Category: "认证", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "/auth/change-password", DisplayName: "修改密码", Category: "认证", Type: models.PermissionTypeBackend, Status: 1},
		// 角色管理 - 后端API
		{Name: "role.create", DisplayName: "创建角色", Category: "角色管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "role.update", DisplayName: "更新角色", Category: "角色管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "role.delete", DisplayName: "删除角色", Category: "角色管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "role.get", DisplayName: "查看角色", Category: "角色管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "role.list", DisplayName: "角色列表", Category: "角色管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "role.assign.permissions", DisplayName: "分配权限", Category: "角色管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "role.permissions", DisplayName: "查看角色权限", Category: "角色管理", Type: models.PermissionTypeBackend, Status: 1},
		// 权限管理 - 后端API
		{Name: "permission.create", DisplayName: "创建权限", Category: "权限管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "permission.update", DisplayName: "更新权限", Category: "权限管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "permission.delete", DisplayName: "删除权限", Category: "权限管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "permission.get", DisplayName: "查看权限", Category: "权限管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "permission.list", DisplayName: "权限列表", Category: "权限管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "permission.list.all", DisplayName: "所有权限", Category: "权限管理", Type: models.PermissionTypeBackend, Status: 1},
		// 用户角色管理 - 后端API
		{Name: "user.assign.roles", DisplayName: "分配角色", Category: "用户管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "user.roles", DisplayName: "查看用户角色", Category: "用户管理", Type: models.PermissionTypeBackend, Status: 1},
		{Name: "user.permissions", DisplayName: "查看用户权限", Category: "用户管理", Type: models.PermissionTypeBackend, Status: 1},
		// 前端菜单权限
		{Name: "menu.dashboard", DisplayName: "仪表板", Category: "菜单", Type: models.PermissionTypeFrontend, Path: "/dashboard", Icon: "LayoutDashboard", Sort: 1, Status: 1},
		{Name: "menu.rbac", DisplayName: "权限管理", Category: "菜单", Type: models.PermissionTypeFrontend, Path: "/rbac", Icon: "Shield", Sort: 2, Status: 1},
		{Name: "menu.rbac.roles", DisplayName: "角色管理", Category: "菜单", Type: models.PermissionTypeFrontend, Path: "/rbac/roles", Icon: "ShieldCheck", Sort: 1, Status: 1},
		{Name: "menu.rbac.permissions", DisplayName: "权限管理", Category: "菜单", Type: models.PermissionTypeFrontend, Path: "/rbac/permissions", Icon: "Shield", Sort: 2, Status: 1},
	}

	for _, perm := range defaultPermissions {
		if err := db.Create(&perm).Error; err != nil {
			return err
		}
	}

	// 为超级管理员角色分配所有权限
	var permissionIDs []uint
	for _, perm := range defaultPermissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	for _, permID := range permissionIDs {
		if err := db.Create(&models.RolePermission{
			RoleID:       superAdminRole.ID,
			PermissionID: permID,
		}).Error; err != nil {
			return err
		}
	}

	// 为默认管理员分配超级管理员角色
	if err := db.Create(&models.UserRole{
		UserID: adminUserID,
		RoleID: superAdminRole.ID,
	}).Error; err != nil {
		return err
	}

	log.Println("[Migration] Default roles and permissions created successfully")
	return nil
}

// MigrateAdmin 执行 admin 模块的数据库迁移
func MigrateAdmin(db *database.DB) error {
	log.Println("[Migration] Starting admin module migration...")

	// 自动迁移表结构
	if err := db.AutoMigrate(&models.AdminUser{}); err != nil {
		return fmt.Errorf("failed to migrate admin tables: %w", err)
	}

	// 迁移 RBAC 相关表
	if err := MigrateRBAC(db); err != nil {
		return fmt.Errorf("failed to migrate RBAC tables: %w", err)
	}

	// 迁移 English Lessons 表（learn_en 模块共享模型）
	if err := db.AutoMigrate(&sharedModels.EnglishLesson{}); err != nil {
		return fmt.Errorf("failed to migrate english_lessons table: %w", err)
	}

	// 检查是否已存在默认管理员
	var count int64
	if err := db.Model(&models.AdminUser{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count admin users: %w", err)
	}

	// 如果没有管理员，创建默认管理员和权限数据
	if count == 0 {
		log.Println("[Migration] Creating default admin user...")

		// 生成密码哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		defaultAdmin := &models.AdminUser{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Nickname: "系统管理员",
			Status:   1,
		}

		if err := db.Create(defaultAdmin).Error; err != nil {
			return fmt.Errorf("failed to create default admin: %w", err)
		}

		log.Println("[Migration] Default admin created successfully (username: admin, password: admin123)")

		// 创建默认超级管理员角色
		if err := createDefaultRoles(db, defaultAdmin.ID); err != nil {
			return fmt.Errorf("failed to create default roles: %w", err)
		}
	}

	log.Println("[Migration] Admin module migration completed")
	return nil
}
