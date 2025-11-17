package task

import (
	"context"
	"encoding/json"
)

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	// Run 执行任务
	// ctx: 上下文
	// job: 任务信息
	// 返回: 任务结果（会被序列化为JSON存储）, 错误信息
	Run(ctx context.Context, job *Job) (interface{}, error)

	// GetType 获取任务类型
	GetType() string
}

// TaskExecutorRegistry 任务执行器注册表
type TaskExecutorRegistry struct {
	executors map[string]TaskExecutor
}

// NewTaskExecutorRegistry 创建任务执行器注册表
func NewTaskExecutorRegistry() *TaskExecutorRegistry {
	return &TaskExecutorRegistry{
		executors: make(map[string]TaskExecutor),
	}
}

// Register 注册任务执行器
func (r *TaskExecutorRegistry) Register(executor TaskExecutor) {
	r.executors[executor.GetType()] = executor
}

// Get 获取任务执行器
func (r *TaskExecutorRegistry) Get(taskType string) (TaskExecutor, bool) {
	executor, ok := r.executors[taskType]
	return executor, ok
}

// GetAllTypes 获取所有已注册的任务类型
func (r *TaskExecutorRegistry) GetAllTypes() []string {
	types := make([]string, 0, len(r.executors))
	for t := range r.executors {
		types = append(types, t)
	}
	return types
}

// UnmarshalPayload 解析任务参数
func UnmarshalPayload(job *Job, v interface{}) error {
	return json.Unmarshal([]byte(job.Payload), v)
}

// MarshalResult 序列化任务结果
func MarshalResult(result interface{}) (string, error) {
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
