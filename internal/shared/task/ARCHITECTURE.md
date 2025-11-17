# Task 系统架构设计

## 系统概述

Task 系统是一个基于 Event Bus 解耦的异步任务执行框架，提供了完整的任务生命周期管理、并行控制、自动重试等企业级特性。

## 核心设计原则

### 1. 事件驱动架构（Event-Driven Architecture）

整个系统基于事件总线进行通信，实现了组件之间的松耦合：

```
业务代码 → EventBus → TaskManager → EventBus → 外部系统
```

**优势**：
- 业务代码无需依赖任务系统的具体实现
- 外部系统可以灵活监听任务事件
- 易于扩展和测试

### 2. 关注点分离（Separation of Concerns）

系统按职责划分为多个独立组件：

| 组件 | 职责 | 不关心 |
|------|------|--------|
| EventBus | 事件分发 | 任务如何存储、如何执行 |
| TaskRepo | 数据持久化 | 任务如何调度、如何执行 |
| TaskManager | 任务调度 | 任务的具体业务逻辑 |
| TaskExecutor | 业务逻辑 | 任务如何调度、如何存储 |

### 3. 依赖倒置原则（Dependency Inversion）

通过接口定义契约，高层模块不依赖低层模块：

```go
// TaskExecutor 是接口，具体实现由业务模块提供
type TaskExecutor interface {
    Run(ctx context.Context, task *Task) (interface{}, error)
    GetType() string
}
```

## 架构图

### 组件关系图

```
┌─────────────────────────────────────────────────────────────┐
│                        应用层 (Application)                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ FileHandler  │  │ EmailHandler │  │  AIHandler   │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                  │                  │              │
│         └──────────────────┼──────────────────┘              │
│                            │SubmitTask                       │
└────────────────────────────┼─────────────────────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      事件总线层 (Event Bus)                  │
│  ┌────────────────────────────────────────────────────┐    │
│  │              TaskEventBus                           │    │
│  │  • Subscribe(eventType, handler)                   │    │
│  │  • Publish(event)                                  │    │
│  │  • SubmitTask(request)                             │    │
│  └────────────────────────────────────────────────────┘    │
└────────────────────────────┬───────────────────────────────┘
                             │ Events
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                     任务管理层 (Task Manager)                │
│  ┌────────────────────────────────────────────────────┐    │
│  │              TaskManager                            │    │
│  │  • 监听任务事件                                     │    │
│  │  • 调度任务执行                                     │    │
│  │  • 控制并行数量                                     │    │
│  │  • 处理重试逻辑                                     │    │
│  └───────┬────────────────────────────────┬───────────┘    │
│          │                                 │                │
└──────────┼─────────────────────────────────┼────────────────┘
           ▼                                 ▼
┌──────────────────────┐         ┌──────────────────────┐
│   存储层 (Repo)      │         │   执行层 (Executor)  │
│  ┌────────────────┐  │         │  ┌────────────────┐  │
│  │   TaskRepo     │  │         │  │FileExecutor    │  │
│  │  • Create      │  │         │  ├────────────────┤  │
│  │  • Update      │  │         │  │EmailExecutor   │  │
│  │  • Query       │  │         │  ├────────────────┤  │
│  └────────────────┘  │         │  │AIExecutor      │  │
│          ▼            │         │  └────────────────┘  │
│    ┌──────────┐      │         │                      │
│    │ Database │      │         │                      │
│    └──────────┘      │         │                      │
└──────────────────────┘         └──────────────────────┘
```

### 数据流图

```
1. 任务提交流程
   Business Code
        │ SubmitTask(request)
        ▼
   EventBus
        │ Publish(EventTypeTaskSubmit)
        ▼
   TaskManager (handleTaskSubmit)
        │ Publish(EventTypeTaskSave)
        ▼
   TaskManager (handleTaskSave)
        │ repo.Create(task)
        ▼
   Database
        │ task saved with status=pending
        ▼
   TaskManager (scheduleTask)

2. 任务执行流程
   TaskManager
        │ 检查并行数量
        ├─ runningCount < maxParallel? → Yes
        │
        │ Publish(EventTypeTaskUpdate) // status=running
        ▼
   Database
        │ task.status = running
        ▼
   TaskManager
        │ Publish(EventTypeTaskStart)
        │ executor.Run(ctx, task)
        ▼
   TaskExecutor
        │ 执行业务逻辑
        ├─ Success → Publish(EventTypeTaskSuccess)
        └─ Failed  → Publish(EventTypeTaskFailed)
        
3. 任务完成流程
   EventBus
        │ EventTypeTaskSuccess/Failed
        ▼
   TaskManager (handleTaskSuccess/Failed)
        │ Publish(EventTypeTaskUpdate)
        ▼
   Database
        │ task.status = completed/failed
        │ task.result = ...
        ▼
   TaskManager
        │ decrementRunningCount()
        │ scheduleTask() // 调度下一个任务
        ▼
   Next Task...
```

## 核心组件详解

### 1. EventBus（事件总线）

**职责**：
- 事件的订阅和发布
- 异步事件分发
- 任务提交入口

**关键方法**：
```go
Subscribe(eventType EventType, handler TaskEventHandler)
Publish(event TaskEvent)  // 异步
PublishSync(event TaskEvent)  // 同步
SubmitTask(req SubmitTaskRequest) error
```

**设计亮点**：
- 支持异步和同步两种发布模式
- 使用 goroutine 池避免 goroutine 泄漏
- panic recovery 保证系统稳定性

### 2. TaskManager（任务管理器）

**职责**：
- 监听任务相关事件
- 调度任务执行
- 控制并行数量
- 处理任务重试

**状态管理**：
```go
runningTasks map[string]int  // 记录每种任务类型的运行数量
```

**并行控制算法**：
```go
func scheduleTask(taskType string) {
    runningCount := getRunningCount(taskType)
    maxParallel := task.MaxParallel
    
    if runningCount >= maxParallel {
        return  // 已达上限，不调度
    }
    
    availableSlots := maxParallel - runningCount
    executeTasks(availableSlots)
}
```

**重试策略**：
```go
if task.Attempts < task.MaxAttempts {
    // 重试：status → pending
} else {
    // 最终失败：status → failed
}
```

### 3. TaskRepo（任务仓库）

**职责**：
- 任务的 CRUD 操作
- 按优先级和状态查询任务
- 统计运行中的任务数量

**关键查询**：
```go
GetPendingTasks(taskType, limit)
    → ORDER BY priority DESC, created_at ASC
    
GetRunningTasksCount(taskType)
    → COUNT(*) WHERE status=running AND type=?
```

### 4. TaskExecutor（任务执行器）

**职责**：
- 实现具体的业务逻辑
- 解析任务参数
- 返回执行结果

**接口设计**：
```go
type TaskExecutor interface {
    Run(ctx context.Context, task *Task) (interface{}, error)
    GetType() string
}
```

**最佳实践**：
- 检查 context 取消，支持超时控制
- 返回结构化的结果对象
- 明确区分可重试和不可重试错误

## 事件系统设计

### 事件类型及触发时机

| 事件类型 | 触发者 | 监听者 | 用途 |
|---------|--------|--------|------|
| TaskSubmit | 业务代码 | TaskManager | 开始任务处理流程 |
| TaskSave | TaskManager | TaskManager | 持久化任务 |
| TaskStart | TaskManager | 外部系统 | 通知任务开始 |
| TaskProgress | TaskExecutor | 外部系统 | 更新进度 |
| TaskSuccess | TaskExecutor | TaskManager, 外部系统 | 任务成功 |
| TaskFailed | TaskExecutor | TaskManager, 外部系统 | 任务失败 |
| TaskUpdate | TaskManager | TaskManager | 更新数据库 |

### 事件流转图

```
[业务代码]
    │
    │ SubmitTask
    ▼
[EventBus] ─────► EventTypeTaskSubmit
    │
    ▼
[TaskManager.handleTaskSubmit]
    │
    │ Publish
    ▼
[EventBus] ─────► EventTypeTaskSave
    │
    ▼
[TaskManager.handleTaskSave]
    │
    │ repo.Create
    ▼
[Database]
    │
    │ scheduleTask
    ▼
[TaskManager.executeTask]
    │
    │ Publish
    ▼
[EventBus] ─────► EventTypeTaskStart ────► [外部监听器]
    │
    │ executor.Run
    ▼
[TaskExecutor]
    │
    ├─ Success ──► EventTypeTaskSuccess ─► [外部监听器]
    │                      │
    │                      ▼
    │              [TaskManager.handleTaskSuccess]
    │                      │
    │                      │ Publish
    │                      ▼
    │              EventTypeTaskUpdate
    │                      │
    │                      ▼
    │              [Database] status=completed
    │
    └─ Failed ───► EventTypeTaskFailed ──► [外部监听器]
                           │
                           ▼
                   [TaskManager.handleTaskFailed]
                           │
                           ├─ 可重试 → EventTypeTaskUpdate (status=pending)
                           └─ 最终失败 → EventTypeTaskUpdate (status=failed)
```

## 并发控制机制

### 并行数量控制

```go
// 每种任务类型独立控制
runningTasks = {
    "file_process": 3,  // 当前运行3个
    "email_send":   5,  // 当前运行5个
    "ai_chat":      2,  // 当前运行2个
}

// 调度时检查
if runningTasks["file_process"] < task.MaxParallel {
    executeTask()
    incrementRunningCount("file_process")
}
```

### 优先级调度

```sql
SELECT * FROM tasks 
WHERE status = 'pending' AND type = ?
ORDER BY priority DESC, created_at ASC
LIMIT ?
```

**调度策略**：
1. 优先级高的任务先执行
2. 相同优先级按创建时间先后顺序
3. 每次调度取出多个任务（批量调度）

### 任务状态机

```
                   ┌──────────┐
                   │ pending  │
                   └────┬─────┘
                        │
          ┌─────────────┼─────────────┐
          │             │             │
          │             ▼             │
          │      ┌──────────┐         │
          │      │ running  │         │
          │      └────┬─────┘         │
          │           │               │
          │    ┌──────┼──────┐        │
          │    │             │        │
          │    ▼             ▼        │
          │ ┌────────┐  ┌────────┐   │ retry
          │ │success │  │ failed │───┘ (attempts < maxAttempts)
          │ └────────┘  └───┬────┘
          │                 │
          │                 ▼
          │            ┌────────┐
          │            │ failed │ (final)
          │            └────────┘
          │
          ▼
     ┌──────────┐
     │cancelled │
     └──────────┘
```

## 扩展性设计

### 1. 插件式执行器

通过注册表模式，支持动态注册执行器：

```go
registry.Register(NewFileExecutor())
registry.Register(NewEmailExecutor())
registry.Register(NewAIExecutor())
```

### 2. 事件监听扩展

外部系统可以监听任何事件：

```go
eventBus.Subscribe(EventTypeTaskSuccess, func(event TaskEvent) {
    // 发送通知
    // 更新统计
    // 触发其他任务
})
```

### 3. 自定义调度策略

可以通过扩展 TaskManager 实现自定义调度：

```go
type CustomTaskManager struct {
    *TaskManager
}

func (m *CustomTaskManager) scheduleTask(taskType string) {
    // 自定义调度逻辑
}
```

## 性能优化

### 1. 数据库优化

**索引**：
- `idx_tasks_type_status` - 加速任务查询
- `idx_tasks_priority` - 优化优先级排序
- `idx_tasks_user_id` - 用户任务查询

**连接池**：
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 2. 并发优化

**Goroutine 池**：
- 使用 sync.WaitGroup 管理 goroutine 生命周期
- 避免无限制创建 goroutine

**锁优化**：
```go
// 读多写少场景使用 RWMutex
sync.RWMutex
```

### 3. 批量调度

每次调度取出多个待执行任务，减少数据库查询：

```go
tasks := repo.GetPendingTasks(taskType, 100)
for _, task := range tasks[:availableSlots] {
    executeTask(task)
}
```

## 容错机制

### 1. Panic Recovery

```go
defer func() {
    if r := recover(); r != nil {
        log.Printf("Handler panic: %v", r)
    }
}()
```

### 2. 任务重试

```go
if task.Attempts < task.MaxAttempts {
    // 自动重试
    repo.UpdateStatus(task.ID, TaskStatusPending)
}
```

### 3. 超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

result, err := executor.Run(ctx, task)
```

### 4. 卡住任务恢复

```go
// 定期重置长时间未更新的任务
repo.ResetStuckTasks(30) // 30分钟未更新的任务
```

## 测试策略

### 1. 单元测试

- Model 测试
- Repo 测试
- EventBus 测试
- Executor 测试

### 2. 集成测试

- 完整任务流程测试
- 并发控制测试
- 重试机制测试

### 3. Mock 测试

```go
type MockTaskExecutor struct {
    executed bool
    result   interface{}
    err      error
}
```

## 监控与运维

### 1. 日志记录

```go
log.Printf("[TaskManager] 任务开始: id=%d, type=%s", task.ID, task.Type)
log.Printf("[TaskManager] 任务完成: id=%d, duration=%v", task.ID, duration)
log.Printf("[TaskManager] 任务失败: id=%d, error=%v", task.ID, err)
```

### 2. 指标统计

建议监控的指标：
- 待执行任务数量
- 运行中任务数量
- 任务执行成功率
- 任务平均执行时长
- 任务重试次数分布

### 3. 告警机制

- 任务堆积告警（待执行任务过多）
- 任务失败率告警
- 任务执行超时告警

## 最佳实践

### 1. 任务粒度

✅ **推荐**：
- 单个文件处理
- 单封邮件发送
- 单次 AI 对话

❌ **不推荐**：
- 批量处理1000个文件（应拆分为1000个任务）
- 复杂的多步骤流程（应拆分为多个任务）

### 2. 幂等性

任务执行器应该是幂等的：

```go
func (e *FileExecutor) Run(ctx context.Context, task *Task) (interface{}, error) {
    // 检查是否已处理
    if alreadyProcessed(task.Payload.FileID) {
        return getExistingResult(task.Payload.FileID), nil
    }
    
    // 处理文件
    result := processFile(task.Payload.FileID)
    return result, nil
}
```

### 3. 错误处理

明确区分可重试和不可重试错误：

```go
// 不可重试错误（如参数错误）
if invalidPayload {
    return nil, fmt.Errorf("invalid payload: %w", err)
}

// 可重试错误（如网络错误）
if networkError {
    return nil, fmt.Errorf("network error: %w", err)
}
```

### 4. 资源清理

确保资源正确释放：

```go
func (e *FileExecutor) Run(ctx context.Context, task *Task) (interface{}, error) {
    file, err := openFile(task.Payload.FileID)
    if err != nil {
        return nil, err
    }
    defer file.Close()  // 确保关闭
    
    // 处理文件
    return result, nil
}
```

## 总结

Task 系统通过事件驱动架构实现了高度解耦的异步任务执行框架，具有以下特点：

✅ **松耦合**：通过 EventBus 解耦各组件
✅ **可扩展**：插件式执行器注册
✅ **高性能**：并发控制、批量调度
✅ **高可靠**：重试机制、容错设计
✅ **易集成**：Wire 依赖注入
✅ **易测试**：接口驱动、Mock 友好

适用场景：
- 异步任务处理
- 后台任务调度
- 批量数据处理
- 定时任务执行
- 工作流引擎
