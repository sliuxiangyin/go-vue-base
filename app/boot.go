package app

import (
	"databaseAi/app/internal/test"
	"databaseAi/app/shared"
)

type App struct {
	Handlers []shared.HandlerInterfaces
}

func (a *App) AddHandler(handler shared.HandlerInterfaces) {
	a.Handlers = append(a.Handlers, handler)
}
func NewApp() (*App, error) {
	app := &App{
		Handlers: make([]shared.HandlerInterfaces, 0),
	}
	//db, err := infra.NewDB(cfg.DatabaseURL)
	//if err != nil {
	//	return nil, err
	//}

	testRepo := test.NewRepo()
	testService := test.NewService(testRepo)
	testHandler := test.NewHandler(testService)
	app.AddHandler(testHandler)
	return app, nil
}
