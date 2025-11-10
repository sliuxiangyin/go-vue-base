#!/usr/bin/env bash just --justfile
# justfile
set shell := ["C:\\Program Files\\Git\\bin\\bash.exe", "-c"]
default:
    @echo "Use 'just dev' or 'just build'"

hello:
  echo "hello world"

# 开发模式
dev:
	MODE=dev && go run main.go

# 构建前端
build-web:
	cd web && pnpm install && pnpm run build

# 构建后端（包含前端）
build: build-web
	go build -o databaseAi .
build_exe: build-web
	GOOS=windows && GOARCH=amd64 &&  go build -o databaseAi.exe  .

# 清理构建文件
clean:
	rm -f databaseAi
	rm -rf frontend/dist

# 安装前端依赖
install-web:
	cd wen && pnpm install

