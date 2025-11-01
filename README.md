# Database AI

这是一个使用 GoFiber 作为后端，Vue 3 + TypeScript 作为前端的全栈项目。

## 项目结构

```
backend/
├── app/                # 应用核心代码
│   ├── boot.go         # 应用启动配置
│   ├── routes.go       # 路由注册
│   ├── config/         # 配置管理
│   ├── infra/          # 基础设施层（数据库等）
│   ├── internal/       # 业务逻辑模块
│   │   └── test/       # 测试模块示例
│   ├── pkg/            # 公共包
│   │   └── response/   # 响应处理
│   └── shared/         # 共享接口定义
├── compile/            # 编译相关代码
│   ├── dev.go          # 开发模式配置
│   ├── prod.go         # 生产模式配置
│   └── start_web.go    # 前端开发服务器启动
├── web/                # Vue 前端项目
│   ├── src/            # 前端源代码
│   ├── public/         # 静态资源
│   ├── package.json    # 前端依赖配置
│   └── ...             # 其他前端配置文件
├── main.go             # GoFiber 后端主程序
├── build_dev.go        # 开发模式构建配置
├── build_prod.go       # 生产模式构建配置
├── go.mod              # Go 模块配置
├── justfile            # 构建命令配置
└── README.md           # 项目说明文档
```

## 功能特性

1. 前后端一体化开发
2. 自动处理前后端接口地址
3. 开发模式下自动启动前端开发服务器
4. 构建时自动打包前端并嵌入到后端二进制文件中
5. 支持热重载开发体验

## 开发环境要求

- Go 1.21+
- Node.js 16+
- npm 8+
- pnpm (推荐) 或 npm

## 快速开始

### 安装依赖

```bash
# 进入后端目录
cd backend

# 安装前端依赖
cd web
pnpm install
cd ..
```

### 开发模式

```bash
# 使用 just 命令（推荐）
just dev

# 或者手动启动
go run -tags dev main.go
```

在开发模式下，后端服务启动在 `http://localhost:8080`，前端开发服务器启动在 `http://localhost:5173`。
所有非 API 请求都会被代理到前端开发服务器，因此你可以通过访问 `http://localhost:8080` 来查看前端页面。

### 构建项目

```bash
# 使用 just 命令（推荐）
just build

# 或者手动构建
cd web
pnpm install
pnpm run build
cd ..

go build -o databaseAi .
```

构建后的二进制文件会包含前端静态资源，可以直接运行。

## API 接口

- `POST /api/test` - 设置键值对
- `GET /api/test/:key` - 获取键值对

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

1. 开发模式：前端和后端分别运行在不同的端口，后端自动代理前端请求
2. 生产模式：前端构建后嵌入到后端二进制文件中，统一部署

## 目录说明

- `web/dist/` - 前端构建输出目录（构建时自动生成）
- `web/src/` - 前端源代码目录
- `web/src/components/` - Vue 组件目录
- `app/internal/` - 业务逻辑模块目录
- `app/pkg/` - 公共功能包目录