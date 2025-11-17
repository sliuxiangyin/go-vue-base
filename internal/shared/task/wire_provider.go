package task

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet Wire 提供者集合
var ProviderSet = wire.NewSet(
	ProvideTaskRepo,
	ProvideEventBus,
	ProvideTaskExecutorRegistry,
	ProvideTaskManager,
)

// ProvideTaskRepo 提供任务仓库
func ProvideTaskRepo(db *gorm.DB) *TaskRepo {
	repo := NewTaskRepo(db)

	// 自动迁移数据库表
	if err := repo.AutoMigrate(); err != nil {
		panic(err)
	}

	return repo
}

// ProvideEventBus 提供全局事件总线单例
func ProvideEventBus() *EventBus {
	return GetEventBus()
}

// ProvideTaskExecutorRegistry 提供任务执行器注册表
// 注意：执行器应该通过 Wire 注入，然后在这里注册
// 使用方式：
//  1. 在你的业务模块中定义 Provide 函数返回执行器切片
//  2. 将执行器切片作为参数传入此函数
//
// 示例：
//
//	func ProvideExecutors(fileRepo *FileRepo) []TaskExecutor {
//	    return []TaskExecutor{
//	        NewFileExecutor(fileRepo),
//	        NewEmailExecutor(),
//	    }
//	}
func ProvideTaskExecutorRegistry(executors []TaskExecutor) *TaskExecutorRegistry {
	registry := NewTaskExecutorRegistry()

	// 注册所有通过依赖注入传入的执行器
	for _, executor := range executors {
		registry.Register(executor)
	}

	return registry
}

// ProvideTaskManager 提供任务管理器
func ProvideTaskManager(
	repo *TaskRepo,
	eventBus *EventBus,
	registry *TaskExecutorRegistry,
) *TaskManager {
	manager := NewTaskManager(repo, eventBus, registry)

	// 启动任务管理器
	manager.Start()

	return manager
}
