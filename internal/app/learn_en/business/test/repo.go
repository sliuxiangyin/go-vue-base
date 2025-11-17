package test

import (
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
	"databaseAi/internal/shared/repo"
)

type Repo struct {
	data         map[string]string
	config       *config.Config
	semanticRepo *repo.SemanticRepo
}

func NewRepo(db *database.DB, config *config.Config, semanticRepo *repo.SemanticRepo) *Repo {
	return &Repo{
		data:         make(map[string]string),
		config:       config,
		semanticRepo: semanticRepo,
	}
}

func (r *Repo) Save(key, value string) {
	r.data[key] = value
}

func (r *Repo) Get(key string) (string, bool) {
	val, ok := r.data[key]
	return val, ok
}
