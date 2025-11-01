# 开发模式
dev:
	MODE=dev go run main.go

# 构建前端
build-frontend:
	cd frontend && npm install && npm run build

# 构建后端（包含前端）
build: build-frontend
	go build -o databaseAi .

# 清理构建文件
clean:
	rm -f databaseAi
	rm -rf frontend/dist

# 安装前端依赖
install-frontend:
	cd frontend && npm install

.PHONY: dev build build-frontend clean install-frontend