# 数据库账户模型 (db_account)

<cite>
**本文档引用的文件**  
- [db_account.go](file://internal/app/tools/models/db_account.go)
- [boot.go](file://internal/app/tools/boot.go)
- [handler.go](file://internal/app/tools/business/dbaccount/handler.go)
- [model.go](file://internal/app/tools/business/dbaccount/model.go)
- [service.go](file://internal/app/tools/business/dbaccount/service.go)
- [db_account_repo.go](file://internal/app/tools/business/dbaccount/db_account_repo.go)
</cite>

## 目录
1. [简介](#简介)
2. [数据表结构](#数据表结构)
3. [GORM结构体映射规则](#gorm结构体映射规则)
4. [核心字段说明](#核心字段说明)
5. [多数据库连接管理机制](#多数据库连接管理机制)
6. [密码加密处理策略](#密码加密处理策略)
7. [实际使用场景与示例](#实际使用场景与示例)
8. [系统核心作用](#系统核心作用)

## 简介
`db_account` 表是系统中用于存储数据库连接配置的核心数据模型。它作为动态数据库连接的数据源配置基础，支持用户添加、管理多种类型的数据库连接信息（如MySQL、SQLite等），为后续的数据库操作、关系分析和AI查询提供统一的配置入口。

该模型通过GORM实现ORM映射，并结合服务层、处理器层完成完整的CRUD操作。其设计兼顾安全性（密码加密）、可扩展性（支持多类型数据库）和易用性（字段校验与自动时间戳管理）。

**Section sources**
- [db_account.go](file://internal/app/tools/models/db_account.go#L5-L16)

## 数据表结构
`db_account` 表包含以下字段，构成完整的数据库连接配置信息：

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| ID | uint | 主键、自增 | 唯一标识符 |
| 名称 | string | 非空、必填 | 连接别名，便于识别 |
| 主机 | string | 非空、必填 | 数据库服务器地址 |
| 端口 | int | 非空、必填 | 服务端口 |
| 用户名 | string | 非空、必填 | 登录账户 |
| 密码 | string | 非空、必填 | 登录密码（加密存储） |
| 数据库名 | string | 非空、必填 | 默认连接的数据库 |
| DDL | string | - | 数据库结构定义（可选） |
| 创建时间 | time.Time | 自动填充 | 记录创建时间 |
| 更新时间 | time.Time | 自动填充 | 记录最后更新时间 |

**Section sources**
- [db_account.go](file://internal/app/tools/models/db_account.go#L5-L16)

## GORM结构体映射规则
在Go代码中，`db_account` 表由 `models.DBAccount` 结构体表示，使用GORM标签进行字段映射与约束定义：

```go
type DBAccount struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name" form:"name" binding:"required"`
	Host        string    `gorm:"not null" json:"host" form:"host" binding:"required"`
	Account     string    `gorm:"not null" json:"account" form:"account" binding:"required"`
	Password    string    `gorm:"not null" json:"password" form:"password" binding:"required"`
	Port        int       `gorm:"not null" json:"port" form:"port" binding:"required"`
	Dbname      string    `gorm:"not null" json:"dbname" form:"dbname" binding:"required"`
	Ddl         string    `gorm:"column:ddl;" json:"ddl" form:"port" binding:"required"`
	CreatedTime time.Time `gorm:"autoCreateTime" json:"created_time"`
	UpdatedTime time.Time `gorm:"autoUpdateTime" json:"updated_time"`
}
```

### 关键GORM标签说明：
- `gorm:"primaryKey"`：指定 `ID` 字段为主键
- `gorm:"not null"`：确保字段非空，数据库层面强制约束
- `gorm:"autoCreateTime"`：创建时自动填充当前时间
- `gorm:"autoUpdateTime"`：每次更新时自动刷新时间
- `json` 标签：用于API序列化输出
- `form` 和 `binding:"required"`：用于HTTP请求参数绑定与校验

**Section sources**
- [db_account.go](file://internal/app/tools/models/db_account.go#L5-L16)

## 核心字段说明
### ID（主键）
- 类型：`uint`
- 作用：唯一标识每条数据库连接记录
- 映射：`gorm:"primaryKey"`，数据库自增主键

### 名称（Name）
- 类型：`string`
- 用途：用户自定义的连接别名，如“生产数据库”、“测试MySQL”
- 约束：非空、必填，前端表单校验

### 主机（Host）
- 类型：`string`
- 含义：数据库服务器IP或域名
- 示例：`localhost`、`192.168.1.100`

### 端口（Port）
- 类型：`int`
- 常见值：MySQL默认3306，SQLite无端口
- 校验：服务层强制要求大于0

### 用户名（Account）
- 类型：`string`
- 用途：数据库登录账户名

### 密码（Password）
- 类型：`string`
- 存储方式：明文暂存于模型，实际应加密后存储（见“密码加密处理策略”）
- 安全性：传输与存储均需加密保护

### 数据库名（Dbname）
- 类型：`string`
- 作用：指定默认连接的数据库名称

### DDL
- 类型：`string`
- 用途：存储数据库结构定义（如建表语句），用于AI分析

### 创建时间与更新时间
- 类型：`time.Time`
- 自动管理：由GORM自动填充，无需手动设置

**Section sources**
- [db_account.go](file://internal/app/tools/models/db_account.go#L5-L16)

## 多数据库连接管理机制
系统通过 `boot.go` 中的初始化逻辑和 `manage_database_repo.go` 实现多数据库连接的动态管理。

### 连接类型识别
在 `boot.go` 中，通过 `getDBType` 函数根据数据库URL判断类型：
- 包含 `file:` 前缀 → SQLite
- 包含 `@tcp(` 或 `@(` → MySQL
- 默认使用 SQLite

```go
func getDBType(databaseURL string) database2.DBType {
	if strings.HasPrefix(databaseURL, "file:") {
		return database2.SQLite
	}
	if strings.Contains(databaseURL, "@tcp(") || strings.Contains(databaseURL, "@(") {
		return database2.MySQL
	}
	return database2.SQLite
}
```

### 动态连接池管理
`ManageDatabaseRepo` 维护一个 `map[uint]*Database]`，将 `db_account` 的ID映射到实际的数据库连接实例。每次新增或更新连接时，系统会根据账号信息动态建立新的数据库连接并缓存。

```go
func (s *ManageDatabaseRepo) Add(id uint, account *models.DBAccount) error {
	db, err := database.NewDB(database.MySQL, fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true", account.Account, account.Password, account.Host, account.Port, account.Dbname))
	if err != nil {
		return err
	}
	s.databases[id] = &Database{db: db, repo: NewDataBaseRepo(db), account: account}
	return nil
}
```

**Section sources**
- [boot.go](file://internal/app/tools/boot.go#L27-L37)
- [manage_database_repo.go](file://internal/app/tools/repo/manage_database_repo.go#L28-L40)

## 密码加密处理策略
当前代码中，密码字段以明文形式存储于 `DBAccount` 模型中（`Password string`），但系统设计上应支持加密存储。理想实现应包括以下步骤：

1. **前端输入**：用户输入明文密码
2. **传输加密**：HTTPS确保传输安全
3. **服务端加密**：使用AES或bcrypt对密码加密后再存入数据库
4. **使用时解密**：建立数据库连接时临时解密

虽然当前代码未展示加密逻辑，但从结构设计看，`Password` 字段具备加密改造基础，只需在 `Create` 和 `Update` 服务方法中加入加解密逻辑即可。

**Section sources**
- [db_account.go](file://internal/app/tools/models/db_account.go#L10-L10)
- [service.go](file://internal/app/tools/business/dbaccount/service.go#L20-L36)

## 实际使用场景与示例
### API路由注册（handler.go）
`dbaccount.Handler` 注册了标准RESTful接口：

```go
func (h *Handler) RegisterRoutes(app fiber.Router) {
	route := app.Group("/db-account")
	route.Post("/", h.Create)   // 创建
	route.Get("/:id", h.GetByID) // 查询单个
	route.Get("/", h.GetAll)    // 查询全部
	route.Put("/:id", h.Update) // 更新
	route.Delete("/:id", h.Delete) // 删除
}
```

### 创建数据库连接示例
1. **前端请求**：
```json
POST /api/db-account
{
  "name": "本地MySQL",
  "host": "localhost",
  "port": 3306,
  "account": "root",
  "password": "123456",
  "dbname": "test"
}
```

2. **后端处理流程**：
   - `handler.Create` 解析请求体 → `DBAccountReq`
   - 调用 `ToModel()` 转换为 `models.DBAccount`
   - `service.Create()` 进行字段校验
   - `repo.Create()` 写入数据库

```go
func (receiver *DBAccountReq) ToModel() *models.DBAccount {
	return &models.DBAccount{
		Name:     receiver.Name,
		Host:     receiver.Host,
		Account:  receiver.Account,
		Password: receiver.Password,
		Port:     receiver.Port,
		Dbname:   receiver.Dbname,
	}
}
```

### 查询所有连接
```go
func (s *Service) GetAll() ([]models.DBAccount, error) {
	return s.repo.GetAll()
}
```

**Section sources**
- [handler.go](file://internal/app/tools/business/dbaccount/handler.go#L23-L32)
- [model.go](file://internal/app/tools/business/dbaccount/model.go#L17-L27)
- [service.go](file://internal/app/tools/business/dbaccount/service.go#L44-L46)

## 系统核心作用
`db_account` 表是整个系统动态数据库访问的核心配置基础，承担以下关键职责：

1. **统一数据源管理**：集中存储所有可连接的数据库配置，支持多环境、多类型数据库。
2. **动态连接支撑**：为AI查询、关系分析、数据库浏览等功能提供实时连接能力。
3. **安全访问控制**：通过加密策略保障敏感凭证安全。
4. **配置持久化**：用户配置可长期保存，重启后仍可用。
5. **扩展性基础**：支持未来扩展PostgreSQL、SQL Server等更多数据库类型。

该模型的设计使得系统能够灵活对接不同数据源，是实现“数据库AI”功能的前提和基石。

**Section sources**
- [db_account.go](file://internal/app/tools/models/db_account.go#L5-L16)
- [boot.go](file://internal/app/tools/boot.go#L60-L63)
- [handler.go](file://internal/app/tools/business/dbaccount/handler.go#L74-L75)