# Database AI

这是一个使用 GoFiber 作为后端，Vue 3 + TypeScript + Shadcn-vue 作为前端的全栈项目。

## 项目结构

```
backend/
├── frontend/           # Vue 前端项目
│   ├── src/            # 前端源代码
│   ├── public/         # 静态资源
│   ├── package.json    # 前端依赖配置
│   └── ...             # 其他前端配置文件
├── main.go             # GoFiber 后端主程序
├── go.mod              # Go 模块配置
├── dev.bat             # 开发模式启动脚本
├── build.bat           # 构建脚本
└── README.md           # 项目说明文档
```

## 功能特性

1. 前后端一体化开发
2. 自动处理前后端接口地址
3. 开发模式下自动启动前端开发服务器
4. 构建时自动打包前端并嵌入到后端二进制文件中

## 开发环境要求

- Go 1.21+
- Node.js 16+
- npm 8+

## 快速开始

### 安装依赖

```bash
# 进入后端目录
cd backend

# 安装前端依赖
cd frontend
npm install
cd ..
```

### 开发模式

```bash
# Windows
dev.bat

# 或者手动启动
set MODE=dev
go run main.go
```

在开发模式下，后端服务启动在 `http://localhost:8080`，前端开发服务器启动在 `http://localhost:5173`，请求会自动代理到后端。

### 构建项目

```bash
# Windows
build.bat

# 或者手动构建
cd frontend
npm install
npm run build
cd ..

go build -o databaseAi.exe .
```

构建后的二进制文件会包含前端静态资源，可以直接运行。

## API 接口

- `/api/hello` - 示例接口，返回问候信息

## 配置说明

### 前端代理配置

前端通过 Vite 的代理功能将 API 请求转发到后端：

```javascript
// vite.config.ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```

### 前后端集成

1. 开发模式：前端和后端分别运行在不同的端口，通过代理通信
2. 生产模式：前端构建后嵌入到后端二进制文件中，统一部署

## 目录说明

- `frontend/dist/` - 前端构建输出目录（构建时自动生成）
- `frontend/src/` - 前端源代码目录
- `frontend/src/components/` - Vue 组件目录
- `frontend/src/views/` - 页面视图目录
- `frontend/src/assets/` - 静态资源目录