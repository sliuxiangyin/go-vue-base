package learn_en

import (
	"databaseAi/internal/app/learn_en/business/test"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
	"databaseAi/internal/shared"
)

type App struct {
	Handlers []shared.HandlerInterfaces
	Config   *config.Config
	DB       *database.DB
}

func NewApp(
	conf *config.Config,
	db *database.DB,
	testHandler *test.Handler) *App {
	app := &App{
		Config: conf,
		DB:     db,
	}
	// 将 test handler 添加到 handlers 列表
	app.AddHandler(testHandler)
	return app
}

func (a *App) AddHandler(handler shared.HandlerInterfaces) {
	a.Handlers = append(a.Handlers, handler)
}
