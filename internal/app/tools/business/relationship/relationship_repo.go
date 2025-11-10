package relationship

import (
	"databaseAi/internal/app/tools/models"
	"databaseAi/internal/infra/database"
	"errors"
	"gorm.io/gorm"
)

// RelationshipRepo 关系仓库
type RelationshipRepo struct {
	db *database.DB
}

// NewRelationshipRepo 创建新的关系仓库实例
func NewRelationshipRepo(db *database.DB) *RelationshipRepo {
	return &RelationshipRepo{
		db: db,
	}
}

// Create 创建关系记录
func (r *RelationshipRepo) Create(relationship *models.Relationship) (*models.Relationship, error) {
	err := r.db.Create(&relationship).Error
	return relationship, err
}

// GetByID 根据ID获取关系记录
func (r *RelationshipRepo) GetByID(id uint) (*models.Relationship, error) {
	var relationship models.Relationship
	err := r.db.First(&relationship, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("关系记录不存在")
		}
		return nil, err
	}
	return &relationship, nil
}

// GetByAccountID 根据账号ID获取关系记录列表
func (r *RelationshipRepo) GetByAccountID(accountID uint) ([]models.Relationship, error) {
	var relationships []models.Relationship
	err := r.db.Where("account_id = ?", accountID).Find(&relationships).Error
	return relationships, err
}

// GetAll 获取所有关系记录
func (r *RelationshipRepo) GetAll() ([]models.Relationship, error) {
	var relationships []models.Relationship
	err := r.db.Find(&relationships).Error
	return relationships, err
}

// Update 更新关系记录
func (r *RelationshipRepo) Update(relationship *models.Relationship) error {
	return r.db.Save(relationship).Error
}

// Delete 删除关系记录
func (r *RelationshipRepo) Delete(id uint) error {
	return r.db.Delete(&models.Relationship{}, id).Error
}
