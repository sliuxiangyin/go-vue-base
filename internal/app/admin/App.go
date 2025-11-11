package admin

import (
	"databaseAi/internal/app/admin/business/auth"
	"databaseAi/internal/app/admin/business/rbac"
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
	authHandler *auth.Handler,
	rbacHandler *rbac.Handler) *App {
	app := &App{
		Config: conf,
		DB:     db,
	}
	// 添加 handlers 到 handlers 列表
	app.AddHandler(authHandler)
	app.AddHandler(rbacHandler)
	return app
}

func (a *App) AddHandler(handler shared.HandlerInterfaces) {
	a.Handlers = append(a.Handlers, handler)
}
