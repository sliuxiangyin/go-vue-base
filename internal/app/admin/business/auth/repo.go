package auth

import (
	"databaseAi/internal/app/admin/models"
	"databaseAi/internal/infra/database"
	"errors"
	"time"
)

type Repo struct {
	db *database.DB
}

func NewRepo(db *database.DB) *Repo {
	return &Repo{db: db}
}

// FindByUsername 根据用户名查找管理员
func (r *Repo) FindByUsername(username string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.Where("username = ? AND status = 1", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找管理员
func (r *Repo) FindByID(id uint) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.Where("id = ? AND status = 1", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *Repo) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.AdminUser{}).
		Where("id = ?", userID).
		Update("last_login", now).Error
}

// CreateUser 创建新用户
func (r *Repo) CreateUser(user *models.AdminUser) error {
	return r.db.Create(user).Error
}

// UpdateUser 更新用户信息
func (r *Repo) UpdateUser(user *models.AdminUser) error {
	if user.ID == 0 {
		return errors.New("user id is required")
	}
	return r.db.Save(user).Error
}
