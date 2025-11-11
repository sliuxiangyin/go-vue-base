//go:build wireinject
// +build wireinject

// wire.go
package main

import (
	"databaseAi/internal/app/admin"
	"databaseAi/internal/app/learn_en"

	"github.com/google/wire"
)

// 注入函数
func InitializeLearnEnService(buildEnv string) (*learn_en.App, error) {
	wire.Build(learn_en.ProviderLearnEnSet)
	return &learn_en.App{}, nil
}
func InitializeAdminService(buildEnv string) (*admin.App, error) {
	wire.Build(admin.ProviderAdminSet)
	return &admin.App{}, nil
}
