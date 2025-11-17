package lesson

import (
	"databaseAi/internal/shared/repo"
)

type Service struct {
	lessonRepo *repo.LessonRepo
}

func NewService(lessonRepo *repo.LessonRepo) *Service {
	return &Service{
		lessonRepo: lessonRepo,
	}
}

// LessonDetailResponse 课程详情响应结构
type LessonDetailResponse struct {
	ID           uint        `json:"id"`
	Title        string      `json:"title"`
	AudioURL     string      `json:"audio_url"`
	SemanticJSON interface{} `json:"semantic_json"`
	Tags         []string    `json:"tags"`
	Level        int8        `json:"level"`
	UpdatedAt    string      `json:"updated_at"`
}

// GetLessonByID 根据 ID 获取课程详情
func (s *Service) GetLessonByID(id uint) (*LessonDetailResponse, error) {
	lesson, err := s.lessonRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &LessonDetailResponse{
		ID:           lesson.ID,
		Title:        lesson.Title,
		AudioURL:     lesson.AudioURL,
		SemanticJSON: lesson.SemanticJSON,
		Tags:         lesson.Tags,
		Level:        lesson.Level,
		UpdatedAt:    lesson.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
