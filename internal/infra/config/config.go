package config

import (
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
	StoragePath string
	JWTSecret   string
}

func NewConfig(buildEnv string) *Config {
	var projectPath, _ = os.Getwd()

	// 尝试加载 .env 文件（如果存在）
	err := godotenv.Load(path.Join(projectPath, ".env"))
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", buildEnv),
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		GrpcDNS:     getEnv("GRPC_DNS", "localhost:50051"),
		OpenaiURl:   getEnv("OPENAI_URL", "https://dashscope.aliyuncs.com"),
		OpenaiKey:   getEnv("OPENAI_KEY", "sk-523cb71b9cbe475ab7e7f27cdab6d379"),
		StoragePath: getEnv("STORAGE_PATH", path.Join(projectPath, "storage")),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-please-change-in-production"),
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

// GetJWTSecret 获取 JWT 密钥
func (c *Config) GetJWTSecret() string {
	return c.JWTSecret
}
