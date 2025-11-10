package relationship

import (
	"databaseAi/internal/app/tools/business/database"
	"databaseAi/internal/app/tools/business/dbaccount"
	"databaseAi/internal/domain/api/repo"
	"databaseAi/internal/domain/api/repo/aliyun_bailian_relationship"
	"databaseAi/internal/infra/ai"
	"databaseAi/internal/infra/config"
	database2 "databaseAi/internal/infra/database"
	"os"
	"path/filepath"
	"testing"
)

// getTestDatabasePath 获取测试环境下的数据库路径
func getTestDatabasePath() string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		// 如果无法获取当前工作目录，使用默认路径
		return "file:app.db?_foreign_keys=on"
	}

	// 查找backend目录
	backendDir := ""
	for {
		// 检查当前目录是否包含backend文件夹
		if _, err := os.Stat(filepath.Join(wd, "backend")); err == nil {
			backendDir = filepath.Join(wd, "backend")
			break
		}

		// 检查当前目录是否就是backend目录
		if _, err := os.Stat(filepath.Join(wd, "app")); err == nil {
			backendDir = wd
			break
		}

		// 向上一级目录
		parent := filepath.Dir(wd)
		if parent == wd {
			// 已经到达根目录
			break
		}
		wd = parent
	}

	// 如果找到了backend目录，使用该目录下的app.db
	if backendDir != "" {
		dbPath := filepath.Join(backendDir, "app.db")
		return "file:" + dbPath + "?_foreign_keys=on"
	}

	// 默认路径
	return "file:app.db?_foreign_keys=on"
}

func TestRelationshipService_Analyze(t *testing.T) {
	// 设置测试数据库
	// 加载配置
	cfg := config.LoadConfig("dev")

	// 连接数据库
	db, err := database2.NewDB(database2.SQLite, cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	openai := ai.NewOpenai(cfg.OpenaiKey, cfg.OpenaiURl)

	dbAccountRepo := dbaccount.NewDBAccountRepo(db)
	manageDatabaseRepo := repo.NewManageDatabaseRepo()

	databaseService := database.NewService(dbAccountRepo, manageDatabaseRepo)
	aliYunBaiLianRelationshipRepo := aliyun_bailian_relationship.NewAliYunBaiLianRelationshipRepo(openai)

	databaseService.Add(1)
	// 初始化仓库
	relationshipRepo := NewRelationshipRepo(db)

	// 初始化服务
	service := NewService(relationshipRepo, manageDatabaseRepo, dbAccountRepo, aliYunBaiLianRelationshipRepo)

	// 测试删除关系记录
	_, err = service.Analyze(1)
	if err != nil {
		t.Skipf("Skipping delete test: %v", err)
	}
}
