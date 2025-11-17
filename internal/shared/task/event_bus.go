package task

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// EventType 任务事件类型
type EventType string

const (
	EventTypeTaskSubmit   EventType = "task.submit"   // 任务提交
	EventTypeTaskSave     EventType = "task.save"     // 保存任务到数据库
	EventTypeTaskStart    EventType = "task.start"    // 开始执行任务
	EventTypeTaskProgress EventType = "task.progress" // 任务进度更新
	EventTypeTaskSuccess  EventType = "task.success"  // 任务成功
	EventTypeTaskFailed   EventType = "task.failed"   // 任务失败
	EventTypeTaskUpdate   EventType = "task.update"   // 更新任务状态到数据库
)

// TaskEvent 任务事件
type TaskEvent struct {
	Type      EventType              `json:"type"`
	TaskID    uint                   `json:"task_id,omitempty"`
	Job       *Job                   `json:"job,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// TaskEventHandler 任务事件处理器
type TaskEventHandler func(event TaskEvent)

// EventBus 事件总线
type EventBus struct {
	handlers map[EventType][]TaskEventHandler
	mu       sync.RWMutex
}

var (
	globalEventBus *EventBus
	eventBusOnce   sync.Once
)

// GetEventBus 获取全局事件总线实例（单例）
func GetEventBus() *EventBus {
	eventBusOnce.Do(func() {
		globalEventBus = &EventBus{
			handlers: make(map[EventType][]TaskEventHandler),
		}
	})
	return globalEventBus
}

// NewEventBus 创建新的事件总线（已废弃，使用 GetEventBus 代替）
// Deprecated: 使用 GetEventBus() 获取全局单例
func NewEventBus() *EventBus {
	return GetEventBus()
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType EventType, handler TaskEventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Publish 发布事件（异步）
func (eb *EventBus) Publish(event TaskEvent) {
	event.Timestamp = time.Now()

	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	for _, handler := range handlers {
		go func(h TaskEventHandler) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[EventBus] Handler panic: %v\n", r)
				}
			}()
			h(event)
		}(handler)
	}
}

// PublishSync 发布事件（同步）
func (eb *EventBus) PublishSync(event TaskEvent) {
	event.Timestamp = time.Now()

	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	for _, handler := range handlers {
		func(h TaskEventHandler) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[EventBus] Handler panic: %v\n", r)
				}
			}()
			h(event)
		}(handler)
	}
}

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
	Type        string                 `json:"type"`         // 任务类型
	UserID      uint                   `json:"user_id"`      // 用户ID
	Payload     map[string]interface{} `json:"payload"`      // 任务参数
	MaxParallel int                    `json:"max_parallel"` // 最大并行数，默认为1
	Priority    int                    `json:"priority"`     // 优先级，默认为0
	MaxAttempts int                    `json:"max_attempts"` // 最大重试次数，默认为3
}

// SubmitTask 提交任务（业务代码调用）
func (eb *EventBus) SubmitTask(req SubmitTaskRequest) error {
	// 序列化 Payload
	payloadJSON, err := json.Marshal(req.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 设置默认值
	if req.MaxParallel <= 0 {
		req.MaxParallel = 1
	}
	if req.MaxAttempts <= 0 {
		req.MaxAttempts = 3
	}

	task := &Job{
		Type:        req.Type,
		UserID:      req.UserID,
		Status:      JobStatusPending,
		Priority:    req.Priority,
		MaxParallel: req.MaxParallel,
		Payload:     string(payloadJSON),
		MaxAttempts: req.MaxAttempts,
	}

	// 发布任务提交事件
	eb.Publish(TaskEvent{
		Type: EventTypeTaskSubmit,
		Job:  task,
		Data: req.Payload,
	})

	return nil
}
