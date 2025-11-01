package config

import (
	"fmt"
	"log"
	"os"

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

func LoadConfig() *Config {
	// 尝试加载 .env 文件（如果存在）
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "production"),
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "file:app.db?_foreign_keys=on"),
		GrpcDNS:     getEnv("GRPC_DNS", "file:app.db?_foreign_keys=on"),
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
