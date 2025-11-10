package repo

import (
	"databaseAi/internal/app/tools/models"
	"databaseAi/internal/infra/database"
)

// MigrateRepo 数据库账号仓库
type MigrateRepo struct {
	db *database.DB
}

// NewMigrateRepo 创建新的数据库账号仓库实例
func NewMigrateRepo(db *database.DB) *MigrateRepo {
	return &MigrateRepo{
		db: db,
	}
}

// Migrate 迁移数据库账号表
func (r *MigrateRepo) Migrate() error {

	_ = r.db.AutoMigrate(&models.DBAccount{})
	_ = r.db.AutoMigrate(&models.Relationship{})

	return nil
}
