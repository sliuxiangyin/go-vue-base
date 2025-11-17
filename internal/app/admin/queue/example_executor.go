package queue

import (
	"context"
	"databaseAi/internal/shared/repo"
	"databaseAi/internal/shared/task"
	"fmt"
	"time"
)

// ExampleTaskPayload 示例任务参数
type ExampleTaskPayload struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

// ExampleTaskResult 示例任务结果
type ExampleTaskResult struct {
	ProcessedMessage string    `json:"processed_message"`
	ProcessedCount   int       `json:"processed_count"`
	ProcessedAt      time.Time `json:"processed_at"`
}

// ExampleTaskExecutor 示例任务执行器
type ExampleTaskExecutor struct {
	lessonRepo *repo.LessonRepo
}

// NewExampleTaskExecutor 创建示例任务执行器
func NewExampleTaskExecutor(lessonRepo *repo.LessonRepo) *ExampleTaskExecutor {
	return &ExampleTaskExecutor{
		lessonRepo: lessonRepo,
	}
}

// GetType 获取任务类型
func (e *ExampleTaskExecutor) GetType() string {
	return "example_task"
}

// Run 执行任务
func (e *ExampleTaskExecutor) Run(ctx context.Context, job *task.Job) (interface{}, error) {
	// 解析任务参数
	var payload ExampleTaskPayload

	if err := task.UnmarshalPayload(job, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// 模拟任务处理
	fmt.Printf("[ExampleTask] 处理任务: id=%d, message=%s, count=%d\n", job.ID, payload.Message, payload.Count)

	// 模拟耗时操作
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Second):
		// 继续执行
	}

	// 返回结果
	result := ExampleTaskResult{
		ProcessedMessage: fmt.Sprintf("Processed: %s", payload.Message),
		ProcessedCount:   payload.Count * 2,
		ProcessedAt:      time.Now(),
	}

	return result, nil
}
