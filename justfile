#!/usr/bin/env bash just --justfile
# justfile
set shell := ["C:\\Program Files\\Git\\bin\\bash.exe", "-c"]
default:
    @echo "Use 'just dev' or 'just build'"

hello:
  echo "hello world"

buf:
    cd  bin/gen_proto && buf generate

wire:
    cd  bin/main && wire

# 开发模式
dev:
	go run -tags dev ./bin/main/

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

py-run:
	cd ./python/learn_en && .venv/Scripts/activate.bat && .venv/Scripts/python.exe main.py
	@echo "Python script executed"
