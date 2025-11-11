package test

import (
	"context"
	protos "databaseAi/internal/app/learn_en/proto"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetValue(key, value string) {
	s.repo.Save(key, value)
}

func (s *Service) GetValue(key string) (string, bool) {
	return s.repo.Get(key)
}

func (s *Service) SynthesizeAudio(model, voice, text, languageType string) (string, error) {
	return s.repo.SynthesizeAudio(model, voice, text, languageType)
}

// AnalyzeSemanticChunks 分析英文句子的语义意群
func (s *Service) AnalyzeSemanticChunks(ctx context.Context, sentence string) ([]string, error) {
	return s.repo.AnalyzeSemanticChunks(ctx, sentence)
}

// TranscribeAudio 音频转录
func (s *Service) TranscribeAudio(audioPath, language, modelSize string) (*protos.TranscribeReply, error) {
	return s.repo.TranscribeAudio(audioPath, language, modelSize)
}
