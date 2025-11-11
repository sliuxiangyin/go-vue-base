package learn_en

import (
	"databaseAi/internal/app/learn_en/business/test"
	"databaseAi/internal/app/learn_en/repo"
	"databaseAi/internal/infra/ai"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
	"databaseAi/internal/infra/grpc"
	"databaseAi/internal/infra/storage"
	sharedRepo "databaseAi/internal/shared/repo"

	"github.com/google/wire"
)

var ProviderLearnEnSet = wire.NewSet(
	ProvideConfig,
	ProvideDB,
	ProvideGrpc,
	ProvideOpenai,
	ProvideFileStorage,
	// Shared Repo layer
	sharedRepo.NewLessonRepo,
	// Repo layer
	repo.NewAudioRepo,
	repo.NewSemanticRepo,
	// Test module
	test.NewRepo,
	test.NewService,
	test.NewHandler,

	// App
	NewApp,
)

// ProvideConfig 提供配置实例
func ProvideConfig(buildEnv string) *config.Config {
	return config.NewConfig(buildEnv)
}

func ProvideGrpc(conf *config.Config) *grpc.GrpcFactory {
	return grpc.NewGrpc(conf.GrpcDNS)
}

func ProvideDB(conf *config.Config) (*database.DB, error) {
	return database.NewDB(conf.DatabaseURL)
}

func ProvideOpenai(conf *config.Config) *ai.Openai {
	return ai.ProvideOpenai(conf)
}

func ProvideFileStorage(conf *config.Config) *storage.FileStorage {
	return storage.NewFileStorage(conf.AppEnv)
}
