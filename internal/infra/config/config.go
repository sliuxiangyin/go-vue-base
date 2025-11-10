package config

import (
	"databaseAi/internal/utils"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     string
	DatabaseURL string
	GrpcDNS     string
	OpenaiURl   string
	OpenaiKey   string
}

func NewConfig(buildEnv string) *Config {
	var projectPath, _ = os.Getwd()

	databaseUrl := ""

	if buildEnv == "dev" {
		projectPath = utils.ProjectRoot()
		databaseUrl = fmt.Sprintf("file:%s?_foreign_keys=on", path.Join(projectPath, "app.db"))
	} else {
		dns := fmt.Sprintf("file:%s?_foreign_keys=on", path.Join(projectPath, "app.db"))
		databaseUrl = getEnv("DATABASE_URL", dns)
	}
	// 尝试加载 .env 文件（如果存在）
	err := godotenv.Load(path.Join(projectPath, ".env"))
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", buildEnv),
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: databaseUrl,
		GrpcDNS:     getEnv("GRPC_DNS", "pkg"),
		OpenaiURl:   getEnv("OPENAI_URL", "https://dashscope.aliyuncs.com"),
		OpenaiKey:   getEnv("OPENAI_KEY", "sk-523cb71b9cbe475ab7e7f27cdab6d379"),
	}

	log.Printf("[config] Loaded configuration: env=%s, port=%s", cfg.AppEnv, cfg.AppPort)
	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
