package dbaccount

import (
	"databaseAi/internal/app/tools/models"
	"errors"
)

// Service 数据库账号服务
type Service struct {
	repo *DBAccountRepo
}

// NewService 创建新的数据库账号服务实例
func NewService(repo *DBAccountRepo) *Service {
	return &Service{
		repo: repo,
	}
}

// Create 创建数据库账号
func (s *Service) Create(account *models.DBAccount) error {
	// 验证必要字段
	if account.Name == "" {
		return errors.New("账号名称不能为空")
	}
	if account.Host == "" {
		return errors.New("主机地址不能为空")
	}
	if account.Password == "" {
		return errors.New("密码不能为空")
	}
	if account.Port <= 0 {
		return errors.New("端口号必须大于0")
	}

	return s.repo.Create(account)
}

// GetByID 根据ID获取数据库账号
func (s *Service) GetByID(id uint) (*models.DBAccount, error) {
	return s.repo.GetByID(id)
}

// GetAll 获取所有数据库账号
func (s *Service) GetAll() ([]models.DBAccount, error) {
	return s.repo.GetAll()
}

// Update 更新数据库账号
func (s *Service) Update(account *models.DBAccount) error {
	// 验证必要字段
	if account.Name == "" {
		return errors.New("账号名称不能为空")
	}
	if account.Host == "" {
		return errors.New("主机地址不能为空")
	}
	if account.Password == "" {
		return errors.New("密码不能为空")
	}
	if account.Port <= 0 {
		return errors.New("端口号必须大于0")
	}

	// 检查账号是否存在
	existing, err := s.repo.GetByID(account.ID)
	if err != nil {
		return errors.New("数据库账号不存在")
	}

	// 保留创建时间
	account.CreatedTime = existing.CreatedTime

	return s.repo.Update(account)
}

// Delete 删除数据库账号
func (s *Service) Delete(id uint) error {
	// 检查账号是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("数据库账号不存在")
	}

	return s.repo.Delete(id)
}
