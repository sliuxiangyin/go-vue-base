package repo

import (
	"databaseAi/internal/infra/database"
	sharedModels "databaseAi/internal/shared/models"
)

// LessonRepo 课程仓库层，使用共享模型
type LessonRepo struct {
	db *database.DB
}

func NewLessonRepo(db *database.DB) *LessonRepo {
	return &LessonRepo{
		db: db,
	}
}

// Create 创建新课程
func (r *LessonRepo) Create(lesson *sharedModels.EnglishLesson) error {
	return r.db.Create(lesson).Error
}

// GetByID 根据 ID 获取课程
func (r *LessonRepo) GetByID(id uint) (*sharedModels.EnglishLesson, error) {
	var lesson sharedModels.EnglishLesson
	err := r.db.First(&lesson, id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// List 获取课程列表
func (r *LessonRepo) List(page, pageSize int, level *int8, isPublic *bool) ([]sharedModels.EnglishLesson, int64, error) {
	var lessons []sharedModels.EnglishLesson
	var total int64

	query := r.db.Model(&sharedModels.EnglishLesson{})

	// 条件过滤
	if level != nil {
		query = query.Where("level = ?", *level)
	}
	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&lessons).Error; err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}

// Update 更新课程
func (r *LessonRepo) Update(lesson *sharedModels.EnglishLesson) error {
	return r.db.Save(lesson).Error
}

// Delete 删除课程（软删除）
func (r *LessonRepo) Delete(id uint) error {
	return r.db.Delete(&sharedModels.EnglishLesson{}, id).Error
}

// SearchByTags 根据标签搜索课程
func (r *LessonRepo) SearchByTags(tags []string) ([]sharedModels.EnglishLesson, error) {
	var lessons []sharedModels.EnglishLesson
	// 注意：JSON 数组查询在不同数据库中语法不同
	// 这里提供 MySQL/SQLite 的示例
	query := r.db.Model(&sharedModels.EnglishLesson{}).Where("is_public = ?", true)

	for _, tag := range tags {
		// SQLite 和 MySQL 兼容的 JSON 查询方式
		// SQLite: 使用 json_extract 或 LIKE
		// MySQL: 使用 JSON_CONTAINS 或 JSON_SEARCH
		query = query.Where("json_extract(tags, '$') LIKE ? OR tags LIKE ?",
			"%\""+tag+"\"%",
			"%\""+tag+"\"%")
	}

	if err := query.Find(&lessons).Error; err != nil {
		return nil, err
	}

	return lessons, nil
}
