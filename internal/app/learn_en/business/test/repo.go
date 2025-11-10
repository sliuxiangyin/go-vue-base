package test

import (
	"databaseAi/internal/app/learn_en/repo"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
)

type Repo struct {
	data      map[string]string
	db        *database.DB
	config    *config.Config
	audioRepo *repo.AudioRepo
}

func NewRepo(db *database.DB, config *config.Config, audioRepo *repo.AudioRepo) *Repo {
	return &Repo{
		data:      make(map[string]string),
		config:    config,
		audioRepo: audioRepo,
	}
}

func (r *Repo) Save(key, value string) {
	r.data[key] = value
}

func (r *Repo) Get(key string) (string, bool) {
	val, ok := r.data[key]
	return val, ok
}

func (r *Repo) SynthesizeAudio(model, voice, text string) ([]byte, error) {
	return r.audioRepo.Synthesizer(model, voice, text)
}
