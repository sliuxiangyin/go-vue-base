set shell := ["bash", "-cu"]

default:
    @echo "Use 'just dev' or 'just build'"

hello:
  echo "hello world"

# 开发模式
dev:
	MODE=dev && go run main.go

# 构建前端
build-web:
	cd web && npm install && npm run build

# 构建后端（包含前端）
build: build-web
	go build -o databaseAi .

# 清理构建文件
clean:
	rm -f databaseAi
	rm -rf frontend/dist

# 安装前端依赖
install-web:
	cd wen && pnpm install

.PHONY: dev build build-frontend clean install-frontend