package migrations

import (
	"databaseAi/internal/app/admin/models"
	"databaseAi/internal/infra/database"
	"fmt"
)

// MigrateRBAC 执行 RBAC 相关表的迁移
func MigrateRBAC(db *database.DB) error {
	// 删除旧的 user_roles 表（如果存在）以便重新创建正确的外键
	if db.Migrator().HasTable("user_roles") {
		fmt.Println("[Migration] Dropping old user_roles table...")
		if err := db.Migrator().DropTable("user_roles"); err != nil {
			fmt.Printf("[Migration] Warning: failed to drop user_roles table: %v\n", err)
		}
	}

	return db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
	)
}
