package task

// 此文件展示如何在你的业务模块中正确配置 Wire 依赖注入

/*
================================
方式一：在业务模块中提供执行器切片
================================

// internal/app/tools/wire_provider.go
package tools

import (
	"github.com/google/wire"
	"your-project/internal/app/tools/business/file"
	"your-project/internal/shared/task"
	"your-project/internal/shared/repo"
)

// ProvideTaskExecutors 提供所有任务执行器
func ProvideTaskExecutors(
	fileRepo *repo.FileRepo,
	lessonRepo *repo.LessonRepo,
	// ... 其他依赖
) []task.TaskExecutor {
	return []task.TaskExecutor{
		// 注册需要依赖注入的执行器
		file.NewFileProcessExecutor(fileRepo),
		file.NewAudioProcessExecutor(fileRepo, lessonRepo),
		// ... 更多执行器
	}
}

// ProviderSet Wire 提供者集合
var ProviderSet = wire.NewSet(
	ProvideTaskExecutors,
	// ... 其他提供者
)

================================
全局事件总线使用示例
================================

// 方式1：通过 Wire 注入（推荐）
type FileHandler struct {
	eventBus *task.EventBus
}

func NewFileHandler(eventBus *task.EventBus) *FileHandler {
	return &FileHandler{
		eventBus: eventBus,
	}
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	// 使用注入的 eventBus
	h.eventBus.SubmitTask(task.SubmitTaskRequest{
		Type: "file_process",
		// ...
	})
	return nil
}

// 方式2：直接获取全局单例
func SomeBusinessFunction() {
	// 直接使用全局单例
	eventBus := task.GetEventBus()
	eventBus.SubmitTask(task.SubmitTaskRequest{
		Type: "ai_chat",
		// ...
	})
}

// 监听事件示例
func init() {
	// 在应用启动时注册事件监听器
	eventBus := task.GetEventBus()

	eventBus.Subscribe(task.EventTypeTaskSuccess, func(event task.TaskEvent) {
		// 处理任务成功事件
		log.Printf("任务成功: %d", event.TaskID)
	})
}

================================
方式二：使用多个 Provide 函数
================================

// internal/app/tools/wire_provider.go
package tools

// 为每个执行器提供独立的 Provide 函数
func ProvideFileProcessExecutor(fileRepo *repo.FileRepo) task.TaskExecutor {
	return file.NewFileProcessExecutor(fileRepo)
}

func ProvideAudioProcessExecutor(
	fileRepo *repo.FileRepo,
	lessonRepo *repo.LessonRepo,
) task.TaskExecutor {
	return file.NewAudioProcessExecutor(fileRepo, lessonRepo)
}

// 聚合所有执行器
func ProvideTaskExecutors(
	fileExecutor task.TaskExecutor,
	audioExecutor task.TaskExecutor,
) []task.TaskExecutor {
	return []task.TaskExecutor{
		fileExecutor,
		audioExecutor,
	}
}

var ProviderSet = wire.NewSet(
	ProvideFileProcessExecutor,
	ProvideAudioProcessExecutor,
	ProvideTaskExecutors,
)

================================
方式三：直接在主 Wire 配置中注册
================================

// bin/main/wire.go
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"your-project/internal/app/tools/business/file"
	"your-project/internal/shared/task"
)

func provideTaskExecutors(
	fileRepo *repo.FileRepo,
	lessonRepo *repo.LessonRepo,
) []task.TaskExecutor {
	return []task.TaskExecutor{
		file.NewFileProcessExecutor(fileRepo),
		task.NewExampleTaskExecutor(),
	}
}

func InitializeApp() (*App, error) {
	wire.Build(
		// 基础设施
		database.NewDB,

		// Repo 层
		repo.NewFileRepo,
		repo.NewLessonRepo,

		// Task 系统
		task.ProviderSet,
		provideTaskExecutors,  // 提供执行器列表

		// 应用层
		NewApp,
	)
	return nil, nil
}

================================
完整示例：业务执行器实现
================================

// internal/app/tools/business/file/file_executor.go
package file

import (
	"context"
	"fmt"
	"your-project/internal/shared/task"
	"your-project/internal/shared/repo"
)

// FileProcessExecutor 文件处理执行器（可以注入依赖）
type FileProcessExecutor struct {
	fileRepo *repo.FileRepo
}

// NewFileProcessExecutor 创建文件处理执行器
func NewFileProcessExecutor(fileRepo *repo.FileRepo) task.TaskExecutor {
	return &FileProcessExecutor{
		fileRepo: fileRepo,
	}
}

func (e *FileProcessExecutor) GetType() string {
	return "file_process"
}

func (e *FileProcessExecutor) Run(ctx context.Context, t *task.Task) (interface{}, error) {
	var payload struct {
		FileID string `json:"file_id"`
		Action string `json:"action"`
	}

	if err := task.UnmarshalPayload(t, &payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	// 使用注入的 fileRepo
	file, err := e.fileRepo.GetByID(payload.FileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 处理文件...
	result := map[string]interface{}{
		"file_id": file.ID,
		"status":  "processed",
	}

	return result, nil
}

================================
AudioProcessExecutor - 多依赖示例
================================

// AudioProcessExecutor 音频处理执行器（注入多个依赖）
type AudioProcessExecutor struct {
	fileRepo   *repo.FileRepo
	lessonRepo *repo.LessonRepo
}

func NewAudioProcessExecutor(
	fileRepo *repo.FileRepo,
	lessonRepo *repo.LessonRepo,
) task.TaskExecutor {
	return &AudioProcessExecutor{
		fileRepo:   fileRepo,
		lessonRepo: lessonRepo,
	}
}

func (e *AudioProcessExecutor) GetType() string {
	return "audio_process"
}

func (e *AudioProcessExecutor) Run(ctx context.Context, t *task.Task) (interface{}, error) {
	// 可以使用 fileRepo 和 lessonRepo
	// ...
	return nil, nil
}

*/
