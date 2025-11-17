package task

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// TaskManager 任务管理器
type TaskManager struct {
	repo         *TaskRepo
	eventBus     *EventBus
	registry     *TaskExecutorRegistry
	runningTasks map[string]int // 记录每种任务类型当前运行的数量
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewTaskManager 创建任务管理器
func NewTaskManager(repo *TaskRepo, eventBus *EventBus, registry *TaskExecutorRegistry) *TaskManager {
	ctx, cancel := context.WithCancel(context.Background())

	tm := &TaskManager{
		repo:         repo,
		eventBus:     eventBus,
		registry:     registry,
		runningTasks: make(map[string]int),
		ctx:          ctx,
		cancel:       cancel,
	}

	// 订阅事件
	tm.subscribeEvents()

	return tm
}

// subscribeEvents 订阅事件
func (tm *TaskManager) subscribeEvents() {
	// 监听任务提交事件
	tm.eventBus.Subscribe(EventTypeTaskSubmit, tm.handleTaskSubmit)

	// 监听任务保存事件
	tm.eventBus.Subscribe(EventTypeTaskSave, tm.handleTaskSave)

	// 监听任务成功事件
	tm.eventBus.Subscribe(EventTypeTaskSuccess, tm.handleTaskSuccess)

	// 监听任务失败事件
	tm.eventBus.Subscribe(EventTypeTaskFailed, tm.handleTaskFailed)

	// 监听任务更新事件
	tm.eventBus.Subscribe(EventTypeTaskUpdate, tm.handleTaskUpdate)
}

// handleTaskSubmit 处理任务提交事件
func (tm *TaskManager) handleTaskSubmit(event TaskEvent) {
	log.Printf("[TaskManager] 收到任务提交事件: type=%s", event.Job.Type)

	// 发布保存任务事件
	tm.eventBus.Publish(TaskEvent{
		Type: EventTypeTaskSave,
		Job:  event.Job,
		Data: event.Data,
	})
}

// handleTaskSave 处理任务保存事件
func (tm *TaskManager) handleTaskSave(event TaskEvent) {
	log.Printf("[TaskManager] 保存任务到数据库: type=%s", event.Job.Type)

	// 保存任务到数据库
	err := tm.repo.Create(event.Job)
	if err != nil {
		log.Printf("[TaskManager] 保存任务失败: %v", err)
		return
	}

	log.Printf("[TaskManager] 任务保存成功: id=%d, type=%s", event.Job.ID, event.Job.Type)

	// 尝试调度任务
	tm.scheduleTask(event.Job.Type)
}

// handleTaskSuccess 处理任务成功事件
func (tm *TaskManager) handleTaskSuccess(event TaskEvent) {
	log.Printf("[TaskManager] 任务执行成功: id=%d", event.TaskID)

	// 发布更新事件
	tm.eventBus.Publish(TaskEvent{
		Type:   EventTypeTaskUpdate,
		TaskID: event.TaskID,
		Data: map[string]interface{}{
			"status":       JobStatusCompleted,
			"result":       event.Data["result"],
			"completed_at": time.Now(),
		},
	})

	// 减少运行计数
	task, _ := tm.repo.GetByID(event.TaskID)
	if task != nil {
		tm.decrementRunningCount(task.Type)
		// 调度下一个任务
		tm.scheduleTask(task.Type)
	}
}

// handleTaskFailed 处理任务失败事件
func (tm *TaskManager) handleTaskFailed(event TaskEvent) {
	log.Printf("[TaskManager] 任务执行失败: id=%d, error=%v", event.TaskID, event.Data["error"])

	task, err := tm.repo.GetByID(event.TaskID)
	if err != nil {
		log.Printf("[TaskManager] 获取任务失败: %v", err)
		return
	}

	// 增加尝试次数
	tm.repo.IncrementAttempts(event.TaskID)
	task.Attempts++

	// 判断是否需要重试
	if task.Attempts < task.MaxAttempts {
		log.Printf("[TaskManager] 任务将重试: id=%d, attempts=%d/%d", event.TaskID, task.Attempts, task.MaxAttempts)
		// 重置为待执行状态
		tm.eventBus.Publish(TaskEvent{
			Type:   EventTypeTaskUpdate,
			TaskID: event.TaskID,
			Data: map[string]interface{}{
				"status": JobStatusPending,
				"error":  event.Data["error"],
			},
		})
	} else {
		log.Printf("[TaskManager] 任务最终失败: id=%d, attempts=%d/%d", event.TaskID, task.Attempts, task.MaxAttempts)
		// 标记为失败
		tm.eventBus.Publish(TaskEvent{
			Type:   EventTypeTaskUpdate,
			TaskID: event.TaskID,
			Data: map[string]interface{}{
				"status":       JobStatusFailed,
				"error":        event.Data["error"],
				"completed_at": time.Now(),
			},
		})
	}

	// 减少运行计数
	tm.decrementRunningCount(task.Type)

	// 调度下一个任务
	tm.scheduleTask(task.Type)
}

// handleTaskUpdate 处理任务更新事件
func (tm *TaskManager) handleTaskUpdate(event TaskEvent) {
	log.Printf("[TaskManager] 更新任务状态: id=%d", event.TaskID)

	err := tm.repo.UpdateStatusWithFields(event.TaskID, event.Data)
	if err != nil {
		log.Printf("[TaskManager] 更新任务失败: %v", err)
	}
}

// scheduleTask 调度任务
func (tm *TaskManager) scheduleTask(taskType string) {
	// 获取执行器
	executor, ok := tm.registry.Get(taskType)
	if !ok {
		log.Printf("[TaskManager] 未找到任务执行器: type=%s", taskType)
		return
	}

	// 获取待执行的任务
	tasks, err := tm.repo.GetPendingTasks(taskType, 100)
	if err != nil {
		log.Printf("[TaskManager] 获取待执行任务失败: %v", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	// 获取第一个任务的并行限制
	maxParallel := tasks[0].MaxParallel

	// 检查当前运行数量
	runningCount := tm.getRunningCount(taskType)
	if runningCount >= maxParallel {
		log.Printf("[TaskManager] 任务并行数已达上限: type=%s, running=%d, max=%d", taskType, runningCount, maxParallel)
		return
	}

	// 计算可以执行的任务数量
	availableSlots := maxParallel - runningCount
	if availableSlots > len(tasks) {
		availableSlots = len(tasks)
	}

	// 执行任务
	for i := 0; i < availableSlots; i++ {
		task := tasks[i]
		tm.incrementRunningCount(taskType)
		tm.wg.Add(1)
		go tm.executeTask(executor, task)
	}
}

// executeTask 执行任务
func (tm *TaskManager) executeTask(executor TaskExecutor, job *Job) {
	defer tm.wg.Done()

	log.Printf("[TaskManager] 开始执行任务: id=%d, type=%s", job.ID, job.Type)

	// 更新任务状态为运行中
	now := time.Now()
	tm.eventBus.PublishSync(TaskEvent{
		Type:   EventTypeTaskUpdate,
		TaskID: job.ID,
		Data: map[string]interface{}{
			"status":     JobStatusRunning,
			"started_at": now,
		},
	})

	// 发布任务开始事件（触发外部逻辑）
	tm.eventBus.Publish(TaskEvent{
		Type:   EventTypeTaskStart,
		TaskID: job.ID,
		Job:    job,
	})

	// 执行任务
	result, err := executor.Run(tm.ctx, job)

	if err != nil {
		// 任务失败
		tm.eventBus.Publish(TaskEvent{
			Type:   EventTypeTaskFailed,
			TaskID: job.ID,
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		})
	} else {
		// 任务成功
		resultJSON, _ := MarshalResult(result)
		tm.eventBus.Publish(TaskEvent{
			Type:   EventTypeTaskSuccess,
			TaskID: job.ID,
			Data: map[string]interface{}{
				"result": resultJSON,
			},
		})
	}
}

// getRunningCount 获取运行中的任务数量
func (tm *TaskManager) getRunningCount(taskType string) int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.runningTasks[taskType]
}

// incrementRunningCount 增加运行计数
func (tm *TaskManager) incrementRunningCount(taskType string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.runningTasks[taskType]++
}

// decrementRunningCount 减少运行计数
func (tm *TaskManager) decrementRunningCount(taskType string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tm.runningTasks[taskType] > 0 {
		tm.runningTasks[taskType]--
	}
}

// Start 启动任务管理器
func (tm *TaskManager) Start() {
	log.Println("[TaskManager] 任务管理器启动")

	// 启动定时调度器
	tm.wg.Add(1)
	go tm.scheduler()
}

// scheduler 定时调度器
func (tm *TaskManager) scheduler() {
	defer tm.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-tm.ctx.Done():
			log.Println("[TaskManager] 调度器停止")
			return
		case <-ticker.C:
			// 遍历所有注册的任务类型
			for _, taskType := range tm.registry.GetAllTypes() {
				tm.scheduleTask(taskType)
			}
		}
	}
}

// Stop 停止任务管理器
func (tm *TaskManager) Stop() {
	log.Println("[TaskManager] 正在停止任务管理器...")
	tm.cancel()
	tm.wg.Wait()
	log.Println("[TaskManager] 任务管理器已停止")
}

// GetTaskStatus 获取任务状态
func (tm *TaskManager) GetTaskStatus(taskID uint) (*Job, error) {
	return tm.repo.GetByID(taskID)
}

// CancelTask 取消任务
func (tm *TaskManager) CancelTask(taskID uint) error {
	job, err := tm.repo.GetByID(taskID)
	if err != nil {
		return err
	}

	if job.Status == JobStatusRunning {
		return fmt.Errorf("cannot cancel running task")
	}

	if job.Status == JobStatusCompleted || job.Status == JobStatusFailed {
		return fmt.Errorf("cannot cancel finished task")
	}

	return tm.repo.UpdateStatusWithFields(taskID, map[string]interface{}{
		"status":       JobStatusCancelled,
		"completed_at": time.Now(),
	})
}
