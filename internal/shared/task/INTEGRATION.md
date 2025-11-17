# Task 系统集成指南

本文档说明如何将 Task 任务系统集成到现有项目中。

## 快速集成步骤

### 1. Wire 依赖注入配置

在你的 Wire 配置文件中（例如 `bin/main/wire.go`）添加 Task 系统的提供者：

```go
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"your-project/internal/infra/database"
	"your-project/internal/shared/task"
	// ... 其他导入
)

func InitializeApp() (*YourApp, error) {
	wire.Build(
		// ... 现有的提供者
		
		// Task 系统提供者
		task.ProviderSet,
		
		// ... 其他提供者
	)
	return nil, nil
}
```

### 2. 在应用启动时初始化

在你的应用启动代码中（例如 `main.go`）：

```go
package main

import (
	"log"
	"your-project/internal/shared/task"
)

func main() {
	// 初始化应用（Wire 自动注入）
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	
	// Task Manager 会在 Wire 初始化时自动启动
	// 确保在程序退出时停止
	defer app.TaskManager.Stop()
	
	// 启动 Web 服务等
	app.Start()
}
```

### 3. 注册自定义任务执行器

修改 `internal/shared/task/wire_provider.go`：

```go
func ProvideTaskExecutorRegistry() *TaskExecutorRegistry {
	registry := NewTaskExecutorRegistry()
	
	// 注册内置的示例执行器
	registry.Register(NewExampleTaskExecutor())
	
	// 注册你的自定义执行器
	registry.Register(NewFileProcessExecutor())
	registry.Register(NewEmailSendExecutor())
	registry.Register(NewReportGenerateExecutor())
	
	return registry
}
```

### 4. 创建自定义执行器

在你的业务模块中创建任务执行器：

```go
// internal/app/tools/business/file/file_task_executor.go
package file

import (
	"context"
	"fmt"
	"your-project/internal/shared/task"
)

type FileProcessExecutor struct {
	fileService *FileService
}

func NewFileProcessExecutor(fileService *FileService) *FileProcessExecutor {
	return &FileProcessExecutor{
		fileService: fileService,
	}
}

func (e *FileProcessExecutor) GetType() string {
	return "file_process"
}

func (e *FileProcessExecutor) Run(ctx context.Context, t *task.Task) (interface{}, error) {
	// 解析参数
	var payload struct {
		FileID string `json:"file_id"`
		Action string `json:"action"`
	}
	
	if err := task.UnmarshalPayload(t, &payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}
	
	// 执行业务逻辑
	result, err := e.fileService.ProcessFile(ctx, payload.FileID, payload.Action)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}
```

### 5. 在业务代码中提交任务

在你的 Handler 或 Service 中使用任务系统：

```go
// internal/app/tools/business/file/handler.go
package file

import (
	"github.com/gofiber/fiber/v2"
	"your-project/internal/shared/task"
)

type FileHandler struct {
	eventBus *task.EventBus
}

func NewFileHandler(eventBus *task.EventBus) *FileHandler {
	return &FileHandler{
		eventBus: eventBus,
	}
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	// ... 文件上传逻辑
	
	// 提交异步处理任务
	err := h.eventBus.SubmitTask(task.SubmitTaskRequest{
		Type:   "file_process",
		UserID: userID,
		Payload: map[string]interface{}{
			"file_id": fileID,
			"action":  "convert_to_pdf",
		},
		MaxParallel: 5,  // 最多同时处理5个文件
		Priority:    0,
	})
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to submit task",
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "文件上传成功，正在处理中",
		"file_id": fileID,
	})
}
```

### 6. 监听任务事件

在需要的地方监听任务事件：

```go
// internal/app/tools/business/file/event_listener.go
package file

import (
	"log"
	"your-project/internal/shared/event"
	"your-project/internal/shared/task"
)

type FileTaskListener struct {
	eventManager *event.EventManager
	taskEventBus *task.EventBus
}

func NewFileTaskListener(
	eventManager *event.EventManager,
	taskEventBus *task.EventBus,
) *FileTaskListener {
	listener := &FileTaskListener{
		eventManager: eventManager,
		taskEventBus: taskEventBus,
	}
	
	// 注册监听器
	listener.registerListeners()
	
	return listener
}

func (l *FileTaskListener) registerListeners() {
	// 监听任务成功事件
	l.taskEventBus.Subscribe(task.EventTypeTaskSuccess, func(e task.TaskEvent) {
		if e.Task != nil && e.Task.Type == "file_process" {
			// 发送 SSE 通知给用户
			l.eventManager.SendToUser(e.Task.UserID, event.Event{
				Type:    event.EventTypeSuccess,
				Message: "文件处理完成",
				Data: map[string]interface{}{
					"task_id": e.TaskID,
					"result":  e.Data["result"],
				},
			})
		}
	})
	
	// 监听任务失败事件
	l.taskEventBus.Subscribe(task.EventTypeTaskFailed, func(e task.TaskEvent) {
		if e.Task != nil && e.Task.Type == "file_process" {
			// 发送错误通知给用户
			l.eventManager.SendToUser(e.Task.UserID, event.Event{
				Type:    event.EventTypeError,
				Message: "文件处理失败",
				Data: map[string]interface{}{
					"task_id": e.TaskID,
					"error":   e.Data["error"],
				},
			})
		}
	})
}
```

### 7. Wire 配置完整示例

如果使用 Wire，需要在提供者中添加监听器：

```go
// internal/app/tools/wire_provider.go
package tools

import (
	"github.com/google/wire"
	"your-project/internal/app/tools/business/file"
	"your-project/internal/shared/event"
	"your-project/internal/shared/task"
)

var ProviderSet = wire.NewSet(
	// ... 现有的提供者
	
	// Task 执行器
	file.NewFileProcessExecutor,
	
	// Task 事件监听器
	file.NewFileTaskListener,
)
```

修改 `internal/shared/task/wire_provider.go`，支持依赖注入自定义执行器：

```go
func ProvideTaskExecutorRegistry(
	fileExecutor *file.FileProcessExecutor,
	// ... 其他执行器
) *TaskExecutorRegistry {
	registry := NewTaskExecutorRegistry()
	
	// 注册执行器
	registry.Register(NewExampleTaskExecutor())
	registry.Register(fileExecutor)
	// registry.Register(emailExecutor)
	// registry.Register(reportExecutor)
	
	return registry
}
```

## API 接口示例

创建任务查询接口：

```go
// internal/app/tools/business/task/handler.go
package taskhandler

import (
	"github.com/gofiber/fiber/v2"
	"your-project/internal/shared/task"
	"strconv"
)

type TaskHandler struct {
	repo    *task.TaskRepo
	manager *task.TaskManager
}

func NewTaskHandler(repo *task.TaskRepo, manager *task.TaskManager) *TaskHandler {
	return &TaskHandler{
		repo:    repo,
		manager: manager,
	}
}

// GetMyTasks 获取我的任务列表
func (h *TaskHandler) GetMyTasks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	status := task.TaskStatus(c.Query("status", ""))
	
	tasks, total, err := h.repo.GetTasksByUserID(userID, status, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"tasks": tasks,
		"total": total,
	})
}

// GetTaskStatus 获取任务状态
func (h *TaskHandler) GetTaskStatus(c *fiber.Ctx) error {
	taskID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}
	
	task, err := h.manager.GetTaskStatus(uint(taskID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
		})
	}
	
	return c.JSON(task)
}

// CancelTask 取消任务
func (h *TaskHandler) CancelTask(c *fiber.Ctx) error {
	taskID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}
	
	err = h.manager.CancelTask(uint(taskID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Task cancelled successfully",
	})
}
```

注册路由：

```go
// internal/app/tools/routes.go
func RegisterRoutes(app *fiber.App, handler *taskhandler.TaskHandler) {
	api := app.Group("/api")
	
	// Task 相关路由
	tasks := api.Group("/tasks")
	tasks.Get("/", handler.GetMyTasks)
	tasks.Get("/:id", handler.GetTaskStatus)
	tasks.Delete("/:id", handler.CancelTask)
}
```

## 前端集成示例

前端通过 SSE 接收任务进度通知：

```typescript
// 提交任务
async function uploadFile(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('/api/files/upload', {
    method: 'POST',
    body: formData,
  });
  
  const data = await response.json();
  console.log('任务已提交:', data.file_id);
}

// 通过 SSE 接收任务完成通知
const eventSource = new EventSource('/api/events/stream');

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'success' && data.data?.task_id) {
    console.log('任务完成:', data.data.task_id);
    console.log('结果:', data.data.result);
    
    // 显示通知
    showNotification('文件处理完成', 'success');
  }
  
  if (data.type === 'error' && data.data?.task_id) {
    console.error('任务失败:', data.data.error);
    showNotification('文件处理失败', 'error');
  }
};
```

## 常见问题

### Q: 如何设置任务优先级？

A: 在提交任务时设置 `Priority` 字段，数值越大优先级越高：

```go
eventBus.SubmitTask(task.SubmitTaskRequest{
    Priority: 10,  // 高优先级
    // ...
})
```

### Q: 如何控制任务的并行数量？

A: 设置 `MaxParallel` 参数：

```go
eventBus.SubmitTask(task.SubmitTaskRequest{
    MaxParallel: 5,  // 最多同时执行5个该类型任务
    // ...
})
```

### Q: 任务执行失败后会自动重试吗？

A: 是的，可以通过 `MaxAttempts` 设置最大重试次数：

```go
eventBus.SubmitTask(task.SubmitTaskRequest{
    MaxAttempts: 3,  // 失败后最多重试3次
    // ...
})
```

### Q: 如何监控任务执行情况？

A: 有两种方式：
1. 订阅任务事件（实时）
2. 查询数据库中的任务记录（历史）

### Q: 任务执行器可以访问其他服务吗？

A: 可以，通过依赖注入将所需服务注入到执行器中：

```go
type MyExecutor struct {
    dbService   *DatabaseService
    fileService *FileService
}

func NewMyExecutor(db *DatabaseService, fs *FileService) *MyExecutor {
    return &MyExecutor{
        dbService:   db,
        fileService: fs,
    }
}
```

## 性能优化建议

1. **合理设置并行数**: 根据任务类型和系统资源设置合适的并行数
2. **使用索引**: 任务表已创建必要的索引，定期清理历史任务
3. **异步处理**: 利用事件总线的异步特性，避免阻塞主流程
4. **监控日志**: 关注任务执行时长，及时发现性能瓶颈

## 总结

Task 系统通过 Event Bus 实现了完全解耦的异步任务执行，具有以下优势：

- ✅ 业务代码无需关心任务如何存储、调度和执行
- ✅ 通过事件监听实现灵活的扩展
- ✅ 支持并行控制、优先级调度、自动重试
- ✅ 完整的状态管理和持久化
- ✅ 易于集成和测试
