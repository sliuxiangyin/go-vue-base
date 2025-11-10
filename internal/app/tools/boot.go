package tools

import (
	"databaseAi/internal/app/tools/business/database"
	"databaseAi/internal/app/tools/business/dbaccount"
	"databaseAi/internal/app/tools/business/relationship"
	"databaseAi/internal/app/tools/business/test"
	"databaseAi/internal/app/tools/repo"
	"databaseAi/internal/app/tools/repo/aliyun_bailian_relationship"
	"databaseAi/internal/infra/ai"
	"databaseAi/internal/infra/config"
	database2 "databaseAi/internal/infra/database"
	"databaseAi/internal/shared"
	"fmt"
	"strings"
)

type App struct {
	Handlers []shared.HandlerInterfaces
	DB       *database2.DB
}

func (a *App) AddHandler(handler shared.HandlerInterfaces) {
	a.Handlers = append(a.Handlers, handler)
}

// getDBType 根据数据库URL判断数据库类型
func getDBType(databaseURL string) database2.DBType {
	if strings.HasPrefix(databaseURL, "file:") {
		return database2.SQLite
	}
	// 检查是否为MySQL连接字符串
	if strings.Contains(databaseURL, "@tcp(") || strings.Contains(databaseURL, "@(") {
		return database2.MySQL
	}
	// 默认使用SQLite
	return database2.SQLite
}

func NewApp(buildEnv string) (*App, error) {
	app := &App{
		Handlers: make([]shared.HandlerInterfaces, 0),
	}
	// 加载配置
	cfg := config.NewConfig(buildEnv)

	// 连接数据库
	db, err := database2.NewDB(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	openai := ai.NewOpenai(cfg.OpenaiKey, cfg.OpenaiURl)
	app.DB = db

	// 初始化测试模块
	testRepo := test.NewRepo(db)

	// 初始化数据库账号模块
	dbAccountRepo := dbaccount.NewDBAccountRepo(db)
	// 初始化关系模块
	relationshipRepo := relationship.NewRelationshipRepo(db)
	// 数据库管理
	manageDatabaseRepo := repo.NewManageDatabaseRepo()
	//阿里千问图谱处理
	aliYunBaiLianRelationshipRepo := aliyun_bailian_relationship.NewAliYunBaiLianRelationshipRepo(openai)

	testService := test.NewService(testRepo)
	testHandler := test.NewHandler(testService)
	app.AddHandler(testHandler)

	dbAccountService := dbaccount.NewService(dbAccountRepo)
	dbAccountHandler := dbaccount.NewHandler(dbAccountService)
	app.AddHandler(dbAccountHandler)

	relationshipService := relationship.NewService(relationshipRepo, manageDatabaseRepo, dbAccountRepo, aliYunBaiLianRelationshipRepo)
	relationshipHandler := relationship.NewHandler(relationshipService)
	app.AddHandler(relationshipHandler)

	databaseService := database.NewService(dbAccountRepo, manageDatabaseRepo)
	databaseHandler := database.NewHandler(databaseService)
	app.AddHandler(databaseHandler)

	migrateRepo := repo.NewMigrateRepo(db)
	// 执行数据库迁移
	if err := migrateRepo.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return app, nil
}
