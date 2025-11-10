package dbaccount

import (
	"databaseAi/internal/app/tools/models"
	"databaseAi/internal/infra/database"
	"errors"
	"gorm.io/gorm"
)

// DBAccountRepo 数据库账号仓库
type DBAccountRepo struct {
	db *database.DB
}

// NewDBAccountRepo 创建新的数据库账号仓库实例
func NewDBAccountRepo(db *database.DB) *DBAccountRepo {
	return &DBAccountRepo{
		db: db,
	}
}

// Create 创建数据库账号
func (r *DBAccountRepo) Create(account *models.DBAccount) error {
	return r.db.Create(account).Error
}

// GetByID 根据ID获取数据库账号
func (r *DBAccountRepo) GetByID(id uint) (*models.DBAccount, error) {
	var account models.DBAccount
	err := r.db.First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据库账号不存在")
		}
		return nil, err
	}
	return &account, nil
}

// GetAll 获取所有数据库账号
func (r *DBAccountRepo) GetAll() ([]models.DBAccount, error) {
	var accounts []models.DBAccount
	err := r.db.Find(&accounts).Error
	return accounts, err
}

// Update 更新数据库账号
func (r *DBAccountRepo) Update(account *models.DBAccount) error {
	return r.db.Save(account).Error
}

// Delete 删除数据库账号
func (r *DBAccountRepo) Delete(id uint) error {
	return r.db.Delete(&models.DBAccount{}, id).Error
}
