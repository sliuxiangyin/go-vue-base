package task

import (
	"time"

	"gorm.io/gorm"
)

// JobStatus 任务状态
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"   // 等待执行
	JobStatusRunning   JobStatus = "running"   // 执行中
	JobStatusCompleted JobStatus = "completed" // 已完成
	JobStatusFailed    JobStatus = "failed"    // 失败
	JobStatusCancelled JobStatus = "cancelled" // 已取消
)

// Job 任务模型
type Job struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Type        string         `json:"type" gorm:"type:varchar(100);not null;index"`           // 任务类型
	UserID      uint           `json:"user_id" gorm:"index"`                                   // 用户ID
	Status      JobStatus      `json:"status" gorm:"type:varchar(20);default:'pending';index"` // 任务状态
	Priority    int            `json:"priority" gorm:"default:0;index"`                        // 优先级（数字越大优先级越高）
	MaxParallel int            `json:"max_parallel" gorm:"default:1"`                          // 最大并行数
	Payload     string         `json:"payload" gorm:"type:text"`                               // 任务参数（JSON）
	Result      string         `json:"result" gorm:"type:text"`                                // 任务结果（JSON）
	Error       string         `json:"error" gorm:"type:text"`                                 // 错误信息
	Attempts    int            `json:"attempts" gorm:"default:0"`                              // 尝试次数
	MaxAttempts int            `json:"max_attempts" gorm:"default:3"`                          // 最大尝试次数
	StartedAt   *time.Time     `json:"started_at"`                                             // 开始时间
	CompletedAt *time.Time     `json:"completed_at"`                                           // 完成时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (Job) TableName() string {
	return "jobs"
}
