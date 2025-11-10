package test

import (
	"databaseAi/internal/infra/database"
)

type Repo struct {
	data map[string]string
	db   *database.DB
}

func NewRepo(db *database.DB) *Repo {
	return &Repo{
		data: make(map[string]string),
	}
}

func (r *Repo) Save(key, value string) {
	r.data[key] = value
}

func (r *Repo) Get(key string) (string, bool) {
	val, ok := r.data[key]
	return val, ok
}
