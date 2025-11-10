package database

import (
	"databaseAi/internal/app/tools/business/dbaccount"
	"databaseAi/internal/app/tools/models"
	"databaseAi/internal/app/tools/repo"
)

type Service struct {
	repo               *dbaccount.DBAccountRepo
	manageDatabaseRepo *repo.ManageDatabaseRepo
}

func NewService(repo *dbaccount.DBAccountRepo, manageDatabaseRepo *repo.ManageDatabaseRepo) *Service {
	return &Service{
		repo:               repo,
		manageDatabaseRepo: manageDatabaseRepo,
	}
}

func (s *Service) GetTables(id uint) ([]models.TableInfo, error) {
	return s.manageDatabaseRepo.GetDatabaseRepo(id).GetAllTables()
}
func (s *Service) Add(id uint) error {
	account, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.manageDatabaseRepo.Add(id, account)

}
func (s *Service) Remove(id uint) error {
	return s.manageDatabaseRepo.Remove(id)
}
