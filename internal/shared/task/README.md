# Task 任务系统

基于 event_bus 解耦的异步任务执行系统，支持任务并行控制、自动重试、状态管理等功能。

## 核心特性

- ✅ **Event Bus 解耦**: 通过事件总线实现业务逻辑与任务执行的解耦
- ✅ **并行控制**: 支持设置任务的最大并行执行数量
- ✅ **自动重试**: 任务失败后自动重试，可配置最大重试次数
- ✅ **优先级调度**: 支持任务优先级，高优先级任务优先执行
- ✅ **状态管理**: 完整的任务状态流转（待执行 → 执行中 → 完成/失败）
- ✅ **事件通知**: 任务生命周期各个阶段发送事件，便于外部系统监听
- ✅ **持久化存储**: 任务状态持久化到数据库

## 架构设计

### 核心组件

```
┌─────────────────┐
│  业务代码       │
│  (Business)     │
└────────┬────────┘
         │ 1. SubmitTask
         ▼
┌─────────────────┐
│   EventBus      │◄─────── 解耦核心
└────────┬────────┘
         │ 2. EventTypeTaskSubmit
         ▼
┌─────────────────┐
│  TaskManager    │
│  (任务调度)     │
└────────┬────────┘
         │ 3. EventTypeTaskSave
         ▼
┌─────────────────┐
│   TaskRepo      │
│  (数据库存储)   │
└─────────────────┘
         │ 4. 检测并行数
         ▼
┌─────────────────┐
│  TaskExecutor   │
│  (任务执行)     │
└────────┬────────┘
         │ 5. EventTypeTaskSuccess/Failed
         ▼
┌─────────────────┐
│   EventBus      │
│  (通知外部)     │
└────────┬────────┘
         │ 6. EventTypeTaskUpdate
         ▼
┌─────────────────┐
│   TaskRepo      │
│  (更新状态)     │
└─────────────────┘
```

### 任务投递流程

1. **业务代码提交任务**: 通过 `EventBus.SubmitTask()` 投递任务，可设置并行数量
2. **Event 通知保存**: 监听 `TaskSubmit` 事件，将任务写入数据库为 `pending` 状态
3. **检测并行限制**: TaskManager 检测当前运行数量，满足条件则执行
4. **任务执行**: 调用 TaskExecutor.Run() 执行任务逻辑
5. **Event 触发外部逻辑**: 发送 `TaskStart`、`TaskProgress`、`TaskSuccess`、`TaskFailed` 事件
6. **完成后更新**: Event 通知更新数据库任务状态

## 使用方法

### 1. 定义任务执行器

实现 `TaskExecutor` 接口：

```go
type MyTaskExecutor struct{}

func (e *MyTaskExecutor) GetType() string {
    return "my_task_type"
}

func (e *MyTaskExecutor) Run(ctx context.Context, task *Task) (interface{}, error) {
    // 解析任务参数
    var payload MyPayload
    if err := UnmarshalPayload(task, &payload); err != nil {
        return nil, err
    }
    
    // 执行任务逻辑
    result := doSomething(payload)
    
    return result, nil
}
```

### 2. 注册任务执行器

```go
// 创建组件
repo := task.NewTaskRepo(db)
eventBus := task.NewEventBus()
registry := task.NewTaskExecutorRegistry()

// 注册执行器
registry.Register(task.NewExampleTaskExecutor())
registry.Register(&MyTaskExecutor{})

// 创建任务管理器
manager := task.NewTaskManager(repo, eventBus, registry)
manager.Start()
```

### 3. 提交任务

在业务代码中提交任务：

```go
err := eventBus.SubmitTask(task.SubmitTaskRequest{
    Type:   "my_task_type",
    UserID: 123,
    Payload: map[string]interface{}{
        "message": "Hello Task",
        "count":   10,
    },
    MaxParallel: 3,  // 最多同时执行3个该类型任务
    Priority:    5,  // 优先级
    MaxAttempts: 3,  // 最多重试3次
})
```

### 4. 监听任务事件

```go
// 监听任务开始事件
eventBus.Subscribe(task.EventTypeTaskStart, func(event task.TaskEvent) {
    fmt.Printf("任务开始: id=%d, type=%s\n", event.TaskID, event.Task.Type)
    // 执行外部逻辑，例如发送通知
})

// 监听任务成功事件
eventBus.Subscribe(task.EventTypeTaskSuccess, func(event task.TaskEvent) {
    fmt.Printf("任务成功: id=%d, result=%v\n", event.TaskID, event.Data["result"])
    // 执行外部逻辑，例如更新业务状态
})

// 监听任务失败事件
eventBus.Subscribe(task.EventTypeTaskFailed, func(event task.TaskEvent) {
    fmt.Printf("任务失败: id=%d, error=%v\n", event.TaskID, event.Data["error"])
    // 执行外部逻辑，例如告警通知
})
```

### 5. 查询任务状态

```go
// 获取任务状态
task, err := manager.GetTaskStatus(taskID)
if err != nil {
    return err
}
fmt.Printf("任务状态: %s\n", task.Status)

// 取消任务
err = manager.CancelTask(taskID)
```

## 事件类型说明

| 事件类型 | 说明 | 何时触发 |
|---------|------|---------|
| `EventTypeTaskSubmit` | 任务提交 | 业务代码调用 SubmitTask 时 |
| `EventTypeTaskSave` | 保存任务 | 收到提交事件后，保存到数据库前 |
| `EventTypeTaskStart` | 任务开始 | 任务开始执行时 |
| `EventTypeTaskProgress` | 任务进度 | 任务执行过程中（需执行器主动发送） |
| `EventTypeTaskSuccess` | 任务成功 | 任务执行成功时 |
| `EventTypeTaskFailed` | 任务失败 | 任务执行失败时 |
| `EventTypeTaskUpdate` | 更新任务 | 需要更新任务状态到数据库时 |

## 任务状态流转

```
pending → running → completed
                 ↘ failed (可重试) → pending
                               ↘ failed (最终失败)
                 ↘ cancelled
```

## 并行控制机制

- 每个任务可以设置 `MaxParallel` 参数
- TaskManager 维护每种任务类型的运行计数
- 只有当 `当前运行数 < MaxParallel` 时才会调度新任务
- 任务完成后自动减少计数，触发下一个任务调度

## 数据库表结构

```sql
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(100) NOT NULL,
    user_id INTEGER,
    status VARCHAR(20) DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    max_parallel INTEGER DEFAULT 1,
    payload TEXT,
    result TEXT,
    error TEXT,
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE INDEX idx_tasks_type ON tasks(type);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_priority ON tasks(priority);
```

## 最佳实践

1. **任务粒度**: 任务应该是独立的工作单元，避免过大或过小
2. **幂等性**: 任务执行器应该是幂等的，重试不会产生副作用
3. **超时控制**: 在 Run 方法中检查 context 取消，避免长时间阻塞
4. **错误处理**: 明确区分可重试错误和不可重试错误
5. **监听事件**: 利用事件系统实现业务逻辑解耦，避免直接依赖

## 示例：完整使用流程

参考 `example_executor.go` 查看完整示例代码。

```go
// 1. 初始化
db := database.NewDB("your_database_url")
repo := task.NewTaskRepo(db.GetDB())
repo.AutoMigrate()

eventBus := task.NewEventBus()
registry := task.NewTaskExecutorRegistry()

// 2. 注册执行器
registry.Register(task.NewExampleTaskExecutor())

// 3. 创建并启动管理器
manager := task.NewTaskManager(repo, eventBus, registry)
manager.Start()
defer manager.Stop()

// 4. 监听事件（可选）
eventBus.Subscribe(task.EventTypeTaskStart, func(event task.TaskEvent) {
    log.Printf("任务开始: %d", event.TaskID)
})

// 5. 提交任务
eventBus.SubmitTask(task.SubmitTaskRequest{
    Type:   "example_task",
    UserID: 1,
    Payload: map[string]interface{}{
        "message": "Hello",
        "count":   5,
    },
    MaxParallel: 2,
    Priority:    1,
})
```

## 扩展开发

### 自定义任务执行器

继承 `TaskExecutor` 接口，实现自己的业务逻辑：

```go
type EmailTaskExecutor struct {
    emailService *EmailService
}

func (e *EmailTaskExecutor) GetType() string {
    return "send_email"
}

func (e *EmailTaskExecutor) Run(ctx context.Context, task *Task) (interface{}, error) {
    var payload struct {
        To      string `json:"to"`
        Subject string `json:"subject"`
        Body    string `json:"body"`
    }
    
    if err := UnmarshalPayload(task, &payload); err != nil {
        return nil, err
    }
    
    // 发送邮件
    err := e.emailService.Send(payload.To, payload.Subject, payload.Body)
    if err != nil {
        return nil, err
    }
    
    return map[string]interface{}{
        "sent_at": time.Now(),
        "to":      payload.To,
    }, nil
}
```

## 技术栈

- Go 语言
- GORM (数据库 ORM)
- Event Bus (事件驱动)
- Context (超时控制)
