package lesson

import (
	"databaseAi/internal/app/admin/business/rbac"
	"databaseAi/internal/shared/models"
	"databaseAi/internal/shared/repo"
)

type Service struct {
	lessonRepo  *repo.LessonRepo
	rbacService *rbac.Service // 用于权限检查
}

func NewService(lessonRepo *repo.LessonRepo, rbacService *rbac.Service) *Service {
	return &Service{
		lessonRepo:  lessonRepo,
		rbacService: rbacService,
	}
}

// CreateLesson 创建课程
func (s *Service) CreateLesson(lesson *models.EnglishLesson) error {
	return s.lessonRepo.Create(lesson)
}

// GetLesson 获取课程详情
func (s *Service) GetLesson(id uint) (*models.EnglishLesson, error) {
	return s.lessonRepo.GetByID(id)
}

// ListLessons 获取课程列表
func (s *Service) ListLessons(page, pageSize int, level *int8, isPublic *bool) ([]models.EnglishLesson, int64, error) {
	return s.lessonRepo.List(page, pageSize, level, isPublic)
}

// UpdateLesson 更新课程
func (s *Service) UpdateLesson(lesson *models.EnglishLesson) error {
	// 先查询确保课程存在
	existing, err := s.lessonRepo.GetByID(lesson.ID)
	if err != nil {
		return err
	}

	// 更新字段
	existing.Title = lesson.Title
	existing.AudioURL = lesson.AudioURL
	existing.Duration = lesson.Duration
	existing.ContentEN = lesson.ContentEN
	existing.ContentZH = lesson.ContentZH
	existing.SemanticJSON = lesson.SemanticJSON
	existing.WordTimestampJSON = lesson.WordTimestampJSON
	existing.PhoneticJSON = lesson.PhoneticJSON
	existing.Tags = lesson.Tags
	existing.Level = lesson.Level
	existing.IsPublic = lesson.IsPublic

	return s.lessonRepo.Update(existing)
}

// DeleteLesson 删除课程（软删除）
func (s *Service) DeleteLesson(id uint) error {
	return s.lessonRepo.Delete(id)
}

// SearchLessonsByTags 根据标签搜索课程
func (s *Service) SearchLessonsByTags(tags []string) ([]models.EnglishLesson, error) {
	return s.lessonRepo.SearchByTags(tags)
}

// CheckUserPermission 实现 PermissionChecker 接口
func (s *Service) CheckUserPermission(userID uint, permissionName string) (bool, error) {
	return s.rbacService.CheckUserPermission(userID, permissionName)
}
