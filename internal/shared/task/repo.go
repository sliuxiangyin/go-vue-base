package task

import (
	"fmt"

	"gorm.io/gorm"
)

// TaskRepo 任务数据访问层
type TaskRepo struct {
	db *gorm.DB
}

// NewTaskRepo 创建任务仓库
func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

// Create 创建任务
func (r *TaskRepo) Create(job *Job) error {
	return r.db.Create(job).Error
}

// Update 更新任务
func (r *TaskRepo) Update(job *Job) error {
	return r.db.Save(job).Error
}

// UpdateStatus 更新任务状态
func (r *TaskRepo) UpdateStatus(taskID uint, status JobStatus) error {
	return r.db.Model(&Job{}).Where("id = ?", taskID).Update("status", status).Error
}

// UpdateStatusWithFields 更新任务状态及其他字段
func (r *TaskRepo) UpdateStatusWithFields(taskID uint, updates map[string]interface{}) error {
	return r.db.Model(&Job{}).Where("id = ?", taskID).Updates(updates).Error
}

// GetByID 根据ID获取任务
func (r *TaskRepo) GetByID(taskID uint) (*Job, error) {
	var job Job
	err := r.db.First(&job, taskID).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetPendingTasks 获取待执行的任务（按优先级和创建时间排序）
func (r *TaskRepo) GetPendingTasks(taskType string, limit int) ([]*Job, error) {
	var jobs []*Job
	query := r.db.Where("status = ? AND type = ?", JobStatusPending, taskType).
		Order("priority DESC, created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&jobs).Error
	return jobs, err
}

// GetRunningTasksCount 获取正在运行的任务数量
func (r *TaskRepo) GetRunningTasksCount(taskType string) (int64, error) {
	var count int64
	err := r.db.Model(&Job{}).
		Where("status = ? AND type = ?", JobStatusRunning, taskType).
		Count(&count).Error
	return count, err
}

// GetTasksByUserID 获取用户的任务列表
func (r *TaskRepo) GetTasksByUserID(userID uint, status JobStatus, limit, offset int) ([]*Job, int64, error) {
	var jobs []*Job
	var total int64

	query := r.db.Model(&Job{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

// IncrementAttempts 增加任务尝试次数
func (r *TaskRepo) IncrementAttempts(taskID uint) error {
	return r.db.Model(&Job{}).
		Where("id = ?", taskID).
		UpdateColumn("attempts", gorm.Expr("attempts + ?", 1)).Error
}

// DeleteTask 删除任务（软删除）
func (r *TaskRepo) DeleteTask(taskID uint) error {
	return r.db.Delete(&Job{}, taskID).Error
}

// GetTasksByType 根据类型获取任务
func (r *TaskRepo) GetTasksByType(taskType string, limit, offset int) ([]*Job, int64, error) {
	var jobs []*Job
	var total int64

	query := r.db.Model(&Job{}).Where("type = ?", taskType)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

// AutoMigrate 自动迁移任务表
func (r *TaskRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&Job{})
}

// ResetStuckTasks 重置卡住的任务（例如：运行中但长时间未更新的任务）
func (r *TaskRepo) ResetStuckTasks(minutes int) error {
	// 将运行中且超过指定分钟未更新的任务重置为待执行
	return r.db.Model(&Job{}).
		Where("status = ? AND updated_at < datetime('now', ?)",
			JobStatusRunning,
			fmt.Sprintf("-%d minutes", minutes)).
		Updates(map[string]interface{}{
			"status": JobStatusPending,
			"error":  "Task reset due to timeout",
		}).Error
}
