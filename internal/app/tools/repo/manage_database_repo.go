package repo

import (
	"databaseAi/internal/app/tools/models"
	"databaseAi/internal/infra/database"
	"fmt"
)

type ManageDatabaseRepo struct {
	databases map[uint]*Database
}

func NewManageDatabaseRepo() *ManageDatabaseRepo {
	return &ManageDatabaseRepo{databases: make(map[uint]*Database)}
}

type Database struct {
	db      *database.DB
	repo    *DataBaseRepo
	account *models.DBAccount
}

func (s *ManageDatabaseRepo) GetDatabaseRepo(id uint) *DataBaseRepo {
	if database, ok := s.databases[id]; ok {
		return database.repo
	}
	return nil
}
func (s *ManageDatabaseRepo) Add(id uint, account *models.DBAccount) error {
	if _, ok := s.databases[id]; ok {
		return nil
	}

	db, err := database.NewDB(database.MySQL, fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true", account.Account, account.Password, account.Host, account.Port, account.Dbname))
	if err != nil {
		return err
	}
	s.databases[id] = &Database{db: db, repo: NewDataBaseRepo(db), account: account}
	return nil

}
func (s *ManageDatabaseRepo) Remove(id uint) error {
	s.databases[id].db.Close()
	delete(s.databases, id)
	return nil
}
