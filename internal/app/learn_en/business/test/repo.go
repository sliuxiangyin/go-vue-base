package test

import (
	"context"
	protos "databaseAi/internal/app/learn_en/proto"
	"databaseAi/internal/app/learn_en/repo"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
)

type Repo struct {
	data         map[string]string
	config       *config.Config
	audioRepo    *repo.AudioRepo
	semanticRepo *repo.SemanticRepo
}

func NewRepo(db *database.DB, config *config.Config, audioRepo *repo.AudioRepo, semanticRepo *repo.SemanticRepo) *Repo {
	return &Repo{
		data:         make(map[string]string),
		config:       config,
		audioRepo:    audioRepo,
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

func (r *Repo) SynthesizeAudio(model, voice, text, languageType string) (string, error) {
	return r.audioRepo.Synthesizer(model, voice, text, languageType)
}

// AnalyzeSemanticChunks 分析英文句子的语义意群
func (r *Repo) AnalyzeSemanticChunks(ctx context.Context, sentence string) ([]string, error) {
	return r.semanticRepo.AnalyzeSemanticChunks(ctx, sentence)
}

// TranscribeAudio 音频转录，返回逐词时间戳
func (r *Repo) TranscribeAudio(audioPath, language, modelSize string) (*protos.TranscribeReply, error) {
	return r.audioRepo.Transcribe(audioPath, language, modelSize)
}
