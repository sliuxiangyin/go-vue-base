package task

import "gorm.io/gorm"

func Init(db *gorm.DB, executors []TaskExecutor) *TaskManager {
	repo := NewTaskRepo(db)

	// 自动迁移数据库表
	if err := repo.AutoMigrate(); err != nil {
		panic(err)
	}
	eventBus := NewEventBus()
	registry := NewTaskExecutorRegistry()
	// 注册所有通过依赖注入传入的执行器
	for _, executor := range executors {
		registry.Register(executor)
	}
	manager := NewTaskManager(repo, eventBus, registry)
	return manager
}
