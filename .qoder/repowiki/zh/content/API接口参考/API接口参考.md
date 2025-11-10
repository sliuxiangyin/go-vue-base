# API接口参考

<cite>
**本文档中引用的文件**  
- [response.go](file://internal/app/tools/infra/response/api/response.go)
- [dbaccount/handler.go](file://internal/app/tools/business/dbaccount/handler.go)
- [db_account_repo.go](file://internal/app/tools/business/dbaccount/db_account_repo.go)
- [database/handler.go](file://internal/app/tools/business/database/handler.go)
- [relationship/handler.go](file://internal/app/tools/business/relationship/handler.go)
- [relationship/service.go](file://internal/app/tools/business/relationship/service.go)
- [db_account.go](file://internal/app/tools/models/db_account.go)
- [relationship.go](file://internal/app/tools/models/relationship.go)
- [table_info.go](file://internal/app/tools/models/table_info.go)
- [model.go](file://internal/app/tools/business/dbaccount/model.go)
- [model.go](file://internal/app/tools/business/relationship/model.go)
</cite>

## 目录
1. [统一响应格式](#统一响应格式)
2. [/dbaccount 接口文档](#dbaccount-接口文档)
3. [/database 接口文档](#database-接口文档)
4. [/relationship 接口文档](#relationship-接口文档)
5. [异步分析流程与轮询机制](#异步分析流程与轮询机制)
6. [认证与速率限制](#认证与速率限制)
7. [调用示例](#调用示例)

## 统一响应格式

所有API接口遵循统一的响应结构，定义于 `response.go` 文件中。标准响应格式如下：

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {}
}
```

- **code**: 响应状态码，200表示成功，其他为错误码
- **msg**: 响应消息，描述请求结果
- **data**: 响应数据，仅在成功时返回，可选字段

成功响应通过 `api.Success()` 方法返回，错误响应通过 `api.Error()` 方法返回。

**Section sources**
- [response.go](file://internal/app/tools/infra/response/api/response.go#L7-L30)

## /dbaccount 接口文档

/dbaccount 接口用于管理数据库连接账号信息。

### GET /api/db-account

获取所有数据库账号列表。

- **请求方法**: GET
- **URL参数**: 无
- **请求体**: 无
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "获取数据库账号列表成功",
  "data": [
    {
      "id": 1,
      "name": "测试数据库",
      "host": "localhost",
      "account": "root",
      "password": "******",
      "port": 3306,
      "dbname": "test"
    }
  ]
}
```
- **错误码**:
  - 500: 获取数据库账号列表失败

**Section sources**
- [dbaccount/handler.go](file://internal/app/tools/business/dbaccount/handler.go#L84-L100)

### POST /api/db-account

创建新的数据库账号。

- **请求方法**: POST
- **URL参数**: 无
- **请求体schema**:
```json
{
  "name": "string",
  "host": "string",
  "account": "string",
  "password": "string",
  "port": "integer",
  "dbname": "string"
}
```
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "数据库账号创建成功",
  "data": { /* 创建的账号信息 */ }
}
```
- **错误码**:
  - 400: 请求参数解析失败或业务验证失败
  - 500: 内部服务器错误

**Section sources**
- [dbaccount/handler.go](file://internal/app/tools/business/dbaccount/handler.go#L34-L56)

### GET /api/db-account/:id

根据ID获取单个数据库账号信息。

- **请求方法**: GET
- **URL参数**:
  - id: 数据库账号ID（路径参数）
- **请求体**: 无
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "获取数据库账号成功",
  "data": { /* 账号信息 */ }
}
```
- **错误码**:
  - 400: 无效的ID参数
  - 404: 未找到指定账号

**Section sources**
- [dbaccount/handler.go](file://internal/app/tools/business/dbaccount/handler.go#L58-L82)

### PUT /api/db-account/:id

更新数据库账号信息。

- **请求方法**: PUT
- **URL参数**:
  - id: 数据库账号ID（路径参数）
- **请求体schema**: 同创建接口
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "数据库账号更新成功",
  "data": { /* 更新后的账号信息 */ }
}
```
- **错误码**:
  - 400: 无效参数或更新失败
  - 404: 未找到指定账号

**Section sources**
- [dbaccount/handler.go](file://internal/app/tools/business/dbaccount/handler.go#L102-L133)

### DELETE /api/db-account/:id

删除指定数据库账号。

- **请求方法**: DELETE
- **URL参数**:
  - id: 数据库账号ID（路径参数）
- **请求体**: 无
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "数据库账号删除成功",
  "data": null
}
```
- **错误码**:
  - 400: 无效的ID参数
  - 404: 未找到指定账号

**Section sources**
- [dbaccount/handler.go](file://internal/app/tools/business/dbaccount/handler.go#L135-L158)

## /database 接口文档

### GET /api/database/init/:id

初始化数据库连接并获取表列表。

- **请求方法**: GET
- **URL参数**:
  - id: 数据库账号ID（路径参数）
- **请求体**: 无
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "ok",
  "data": [
    {
      "name": "users",
      "comment": "用户表"
    },
    {
      "name": "orders",
      "comment": "订单表"
    }
  ]
}
```
- **错误码**:
  - 400: 请求参数解析失败或数据库连接失败

**Section sources**
- [database/handler.go](file://internal/app/tools/business/database/handler.go#L21-L32)

### GET /api/database/table/:id/:table

获取指定表的详细信息（注：此接口在当前代码中未实现，可能为规划中功能）。

- **请求方法**: GET
- **URL参数**:
  - id: 数据库账号ID
  - table: 表名
- **请求体**: 无
- **响应**: 表结构详情

## /relationship 接口文档

### POST /api/relationship

创建关系分析记录。

- **请求方法**: POST
- **URL参数**: 无
- **请求体schema**:
```json
{
  "account_id": 1,
  "title": "订单系统分析",
  "tables": "[\"users\",\"orders\"]",
  "ddls": "[\"CREATE TABLE...\", \"CREATE TABLE...\"]"
}
```
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "关系记录创建成功",
  "data": { /* 创建的记录 */ }
}
```
- **错误码**:
  - 400: 参数解析失败或业务验证失败

**Section sources**
- [relationship/handler.go](file://internal/app/tools/business/relationship/handler.go#L38-L59)

### GET /api/relationship/:id

根据ID获取关系记录。

- **请求方法**: GET
- **URL参数**:
  - id: 关系记录ID
- **请求体**: 无
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "获取关系记录成功",
  "data": { /* 关系记录 */ }
}
```
- **错误码**:
  - 400: 无效ID参数
  - 404: 记录不存在

**Section sources**
- [relationship/handler.go](file://internal/app/tools/business/relationship/handler.go#L61-L85)

### GET /api/relationship/account/:accountID

根据账号ID获取关系记录列表。

- **请求方法**: GET
- **URL参数**:
  - accountID: 数据库账号ID
- **请求体**: 无
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "获取关系记录列表成功",
  "data": [ /* 记录列表 */ ]
}
```
- **错误码**:
  - 400: 无效账号ID
  - 404: 未找到记录

**Section sources**
- [relationship/handler.go](file://internal/app/tools/business/relationship/handler.go#L87-L111)

### PUT /api/relationship/:id

更新关系记录。

- **请求方法**: PUT
- **URL参数**:
  - id: 关系记录ID
- **请求体schema**: 同创建接口
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "关系记录更新成功",
  "data": { /* 更新后的记录 */ }
}
```
- **错误码**:
  - 400: 参数错误或更新失败
  - 404: 记录不存在

**Section sources**
- [relationship/handler.go](file://internal/app/tools/business/relationship/handler.go#L131-L162)

### DELETE /api/relationship/:id

删除关系记录。

- **请求方法**: DELETE
- **URL参数**:
  - id: 关系记录ID
- **请求体**: 无
- **成功响应示例**:
```json
{
  "code": 200,
  "msg": "关系记录删除成功",
  "data": null
}
```
- **错误码**:
  - 400: 无效ID
  - 404: 记录不存在

**Section sources**
- [relationship/handler.go](file://internal/app/tools/business/relationship/handler.go#L164-L187)

## 异步分析流程与轮询机制

### GET /api/relationship/:id/analyze

触发AI分析指定关系记录的数据库表结构。

- **请求方法**: GET
- **URL参数**:
  - id: 关系记录ID
- **响应类型**: SSE (Server-Sent Events)
- **流程说明**:
  1. 客户端发起GET请求到 `/api/relationship/{id}/analyze`
  2. 服务端建立SSE连接，设置相应响应头
  3. 启动异步goroutine执行AI分析任务
  4. 分析过程中通过消息通道(channel)实时推送分析进度和结果
  5. 任务完成后推送完成事件并关闭连接

- **SSE消息格式**:
```
data: {"code":200,"msg":"分析中","data":"正在连接数据库..."}
data: {"code":200,"msg":"分析中","data":"正在分析表结构..."}
event: done
```

- **错误处理**:
  - 400: 无效的关系记录ID
  - 分析过程中的错误会通过SSE消息流返回

**Section sources**
- [relationship/handler.go](file://internal/app/tools/business/relationship/handler.go#L189-L217)
- [relationship/service.go](file://internal/app/tools/business/relationship/service.go#L1-L200)

## 认证与速率限制

### 认证机制

当前系统**未实现认证机制**，所有API接口均为公开访问。建议在生产环境中添加身份验证层以确保安全性。

### 速率限制

当前代码中未实现显式的速率限制策略。系统依赖于底层框架的默认行为。建议在生产环境中配置速率限制以防止滥用。

**Section sources**
- [routes.go](file://internal/app/tools/routes.go#L1-L11)
- [pnpm-lock.yaml](file://web/database_api/pnpm-lock.yaml#L1437-L1468)

## 调用示例

### curl命令示例

**创建数据库账号**:
```bash
curl -X POST http://localhost:3000/api/db-account \
  -H "Content-Type: application/json" \
  -d '{
    "name": "生产数据库",
    "host": "prod-db.example.com",
    "account": "admin",
    "password": "secure123",
    "port": 3306,
    "dbname": "production"
  }'
```

**触发AI分析**:
```bash
curl http://localhost:3000/api/relationship/1/analyze
```

### JavaScript调用片段

**使用fetch获取表列表**:
```javascript
async function getTables(accountId) {
  try {
    const response = await fetch(`/api/database/init/${accountId}`);
    const result = await response.json();
    if (result.code === 200) {
      console.log('表列表:', result.data);
    } else {
      console.error('错误:', result.msg);
    }
  } catch (error) {
    console.error('请求失败:', error);
  }
}
```

**使用SSE监听AI分析进度**:
```javascript
function startAnalysis(relationshipId) {
  const eventSource = new EventSource(`/api/relationship/${relationshipId}/analyze`);
  
  eventSource.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('分析进度:', data.msg, data.data);
  };
  
  eventSource.addEventListener('done', function(event) {
    console.log('分析完成');
    eventSource.close();
  });
  
  eventSource.onerror = function(error) {
    console.error('SSE连接错误:', error);
    eventSource.close();
  };
}
```

**Section sources**
- [relationshipApi.ts](file://web/src/lib/relationshipApi.ts)
- [api.ts](file://web/src/lib/api.ts)