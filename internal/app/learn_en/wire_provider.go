package learn_en

import (
	"databaseAi/internal/app/learn_en/business/test"
	"databaseAi/internal/app/learn_en/repo"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
	"databaseAi/internal/infra/grpc"

	"github.com/google/wire"
)

var ProviderLearnEnSet = wire.NewSet(
	ProvideConfig,
	ProvideDB,
	ProvideGrpc,
	// Repo layer
	repo.NewAudioRepo,
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
