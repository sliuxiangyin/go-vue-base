package admin

import (
	"databaseAi/internal/app/admin/business/auth"
	"databaseAi/internal/app/admin/business/event"
	"databaseAi/internal/app/admin/business/file"
	"databaseAi/internal/app/admin/business/lesson"
	"databaseAi/internal/app/admin/business/rbac"
	"databaseAi/internal/app/admin/queue"
	"databaseAi/internal/infra/ai"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
	"databaseAi/internal/infra/grpc"
	"databaseAi/internal/infra/storage"
	"databaseAi/internal/shared"
	"databaseAi/internal/shared/repo"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderAdminSet = wire.NewSet(
	ProvideConfig,
	ProvideDB,
	ProvideGormDB,
	ProvideGrpc,
	ProvideOpenai,
	ProvideFileStorage,
	// Shared Repo
	repo.NewLessonRepo,
	repo.NewPhoneticDictionaryRepo,
	repo.NewAudioRepo,
	repo.NewSemanticRepo,

	// Auth module
	auth.NewRepo,
	ProvideAuthNewService,
	ProvideAuthNewHandler,

	// RBAC module
	rbac.NewRepo,
	rbac.NewService,
	ProvideRBACNewHandler,

	// Lesson module
	lesson.NewService,
	ProvideLessonNewHandler,
	queue.NewLessonExecutor,
	// File module
	ProvideFileRepository,
	ProvideLocalFileStorage,
	ProvideFileNewService,
	ProvideFileNewHandler,

	// Event module
	ProvideEventNewHandler,

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

// ProvideGormDB 从 database.DB 提供 *gorm.DB
func ProvideGormDB(db *database.DB) *gorm.DB {
	return db.DB
}
func ProvideGrpc(conf *config.Config) *grpc.GrpcFactory {
	return grpc.NewGrpc(conf.GrpcDNS)
}
func ProvideOpenai(conf *config.Config) *ai.Openai {
	return ai.ProvideOpenai(conf)
}

// ProvideFileRepository 提供 FileRepository 接口实现
func ProvideFileRepository(db *gorm.DB) shared.FileRepository {
	return repo.NewFileRepo(db)
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

func ProvideLessonNewHandler(service *lesson.Service, authService *auth.Service) *lesson.Handler {
	return lesson.NewHandler(service, authService)
}

func ProvideFileNewService(repo shared.FileRepository, storage shared.FileStorage) *file.Service {
	return file.NewService(repo, storage)
}

func ProvideLocalFileStorage(conf *config.Config) shared.FileStorage {
	return storage.NewLocalFileStorage(conf.AppEnv)
}

func ProvideFileNewHandler(service *file.Service, authService *auth.Service) *file.Handler {
	return file.NewHandler(service, authService)
}
func ProvideFileStorage(conf *config.Config) *storage.FileStorage {
	return storage.NewFileStorage(conf.AppEnv)
}

func ProvideEventNewHandler(authService *auth.Service) *event.Handler {
	return event.NewHandler(authService)
}
