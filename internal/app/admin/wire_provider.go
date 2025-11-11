package admin

import (
	"databaseAi/internal/app/admin/business/auth"
	"databaseAi/internal/app/admin/business/rbac"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"

	"github.com/google/wire"
)

var ProviderAdminSet = wire.NewSet(
	ProvideConfig,
	ProvideDB,
	// Auth module
	auth.NewRepo,
	ProvideAuthNewService,
	ProvideAuthNewHandler,

	// RBAC module
	rbac.NewRepo,
	rbac.NewService,
	ProvideRBACNewHandler,

	// App
	NewApp,
)

// ProvideConfig 提供配置实例
func ProvideConfig(buildEnv string) *config.Config {
	return config.NewConfig(buildEnv)
}

// ProvideDB 提供数据库实例
func ProvideDB(conf *config.Config) (*database.DB, error) {
	return database.NewDB(conf.DatabaseURL)
}
func ProvideAuthNewService(repo *auth.Repo, conf *config.Config) *auth.Service {
	return auth.NewService(repo, conf.JWTSecret)
}

func ProvideAuthNewHandler(service *auth.Service) *auth.Handler {
	return auth.NewHandler(service)
}

func ProvideRBACNewHandler(service *rbac.Service, authService *auth.Service) *rbac.Handler {
	return rbac.NewHandler(service, authService)
}
