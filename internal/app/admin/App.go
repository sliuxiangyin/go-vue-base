package admin

import (
	"databaseAi/internal/app/admin/business/auth"
	"databaseAi/internal/app/admin/business/event"
	"databaseAi/internal/app/admin/business/file"
	"databaseAi/internal/app/admin/business/lesson"
	"databaseAi/internal/app/admin/business/rbac"
	"databaseAi/internal/app/admin/queue"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
	"databaseAi/internal/shared"
	"databaseAi/internal/shared/task"
)

type App struct {
	Handlers []shared.HandlerInterfaces
	Config   *config.Config
	DB       *database.DB
}

func NewApp(
	conf *config.Config,
	db *database.DB,

	authHandler *auth.Handler,
	rbacHandler *rbac.Handler,
	lessonHandler *lesson.Handler,
	eventHandler *event.Handler,
	fileHandler *file.Handler,

	lessonExecutor *queue.LessonExecutor,
) *App {
	app := &App{
		Config: conf,
		DB:     db,
	}
	executors := make([]task.TaskExecutor, 0)
	executors = append(executors, lessonExecutor)
	//任务系统
	task.Init(db.DB, executors)
	// 添加 handlers 到 handlers 列表
	app.AddHandler(authHandler)
	app.AddHandler(rbacHandler)
	app.AddHandler(lessonHandler)
	app.AddHandler(eventHandler)
	app.AddHandler(fileHandler)
	return app
}

func (a *App) AddHandler(handler shared.HandlerInterfaces) {
	a.Handlers = append(a.Handlers, handler)
}
