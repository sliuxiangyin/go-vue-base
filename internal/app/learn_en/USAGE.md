# Learn English 模块 - 共享模型使用说明

## 概述

本模块展示了如何在 `learn_en` 中优化使用 `admin` 中的模型，通过创建 **共享模型层** (`internal/shared/models`) 实现模型复用和解耦。

## 架构设计

```
internal/
├── shared/
│   └── models/
│       └── english_lesson.go  # ✅ 共享模型（learn_en 和 admin 都可使用）
├── app/
│   ├── admin/
│   │   ├── migrations/
│   │   │   ├── migrate.go     # ✅ 统一迁移入口（包含 english_lessons）
│   │   │   └── rbac_migrate.go
│   │   └── models/
│   │       ├── admin_user.go
│   │       └── role.go
│   └── learn_en/
│       ├── repo/
│       │   ├── lesson_repo.go      # ✅ 使用共享模型的仓库
│       │   └── lesson_repo_test.go # ✅ 完整的测试用例
│       └── wire_provider.go
```

## 核心优势

### 1. **模型共享，避免重复**
- ✅ `EnglishLesson` 模型定义在 `internal/shared/models`
- ✅ `PhoneticDictionary` 发音词典独立存储，全局共享
- ✅ `admin` 和 `learn_en` 都可以直接使用
- ✅ 单一数据源，避免不一致

### 2. **迁移统一管理**
- ✅ 所有数据库迁移在 `admin/migrations` 中统一执行
- ✅ `learn_en` 无需独立迁移文件
- ✅ 启动时自动创建 `english_lessons` 表

### 3. **完整的 JSON 支持**
- ✅ 语义意群 (`SemanticChunks`)
- ✅ 逐词时间戳 (`WordTimestamps`)
- ✅ 标签数组 (`Tags`)
- ✅ 自动序列化/反序列化
- ✅ 发音信息独立存储在 `phonetic_dictionary` 表，全局共享

### 4. **类型安全**
```go
// 强类型的 JSON 字段
lesson.SemanticJSON = sharedModels.SemanticChunks{
    {Text: "Life was like a box", Start: 0.0, End: 2.5},
}

// 而不是使用 string 或 interface{}
```

## 使用示例

### 1. 创建课程

```go
import sharedModels "databaseAi/internal/shared/models"

lesson := &sharedModels.EnglishLesson{
    Title:     "Forrest Gump - Opening",
    AudioURL:  "/audio/forrest_gump.mp3",
    Duration:  23.6,
    ContentEN: "Life was like a box of chocolates...",
    ContentZH: "生活就像一盒巧克力...",
    SemanticJSON: sharedModels.SemanticChunks{
        {Text: "Life was like a box of chocolates", Start: 0.0, End: 2.5},
    },
    WordTimestampJSON: sharedModels.WordTimestamps{
        {Word: "Life", Start: 0.0, End: 0.3, Confidence: 0.98},
    },
    Tags:     sharedModels.Tags{"movie", "classic"},
    Level:    2,
    IsPublic: true,
}

err := lessonRepo.Create(lesson)
```

### 2. 查询课程

```go
// 按 ID 查询
lesson, err := lessonRepo.GetByID(1)

// 分页查询
lessons, total, err := lessonRepo.List(page, pageSize, nil, nil)

// 按难度过滤
level := int8(2)
lessons, total, err := lessonRepo.List(1, 10, &level, nil)

// 按标签搜索
lessons, err := lessonRepo.SearchByTags([]string{"movie", "classic"})
```

### 3. 更新和删除

```go
// 更新
lesson.Title = "New Title"
err := lessonRepo.Update(lesson)

// 软删除
err := lessonRepo.Delete(lessonID)
```

## 数据库表结构

```sql
CREATE TABLE `english_lessons` (
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime,
  `title` varchar(255) NOT NULL,
  `audio_url` varchar(512),
  `duration` real,
  `content_en` text,
  `content_zh` text,
  `semantic_json` json,
  `word_timestamp_json` json,
  `phonetic_json` json,
  `tags` json,
  `level` tinyint,
  `is_public` tinyint DEFAULT false
);
```

## 测试

所有测试已通过 ✅

```bash
# 运行所有测试
go test ./internal/app/learn_en/repo -v

# 运行特定测试
go test ./internal/app/learn_en/repo -v -run TestLessonRepo_Create
```

测试覆盖：
- ✅ 创建课程（包含复杂 JSON 字段）
- ✅ 查询单个课程
- ✅ 分页列表查询
- ✅ 条件过滤（level, is_public）
- ✅ 标签搜索
- ✅ 更新课程
- ✅ 软删除

## 依赖注入配置

在 `wire_provider.go` 中已添加：

```go
var ProviderLearnEnSet = wire.NewSet(
    // ... 其他依赖
    repo.NewLessonRepo, // ✅ 新增：课程仓库（使用共享模型）
    // ...
)
```

## 迁移执行

启动应用时，`admin` 的迁移会自动创建 `english_lessons` 表：

```go
// bin/main/main.go
if err := migrations.MigrateAdmin(adminApp.DB); err != nil {
    log.Fatalf("Failed to migrate admin database: %v", err)
}
```

## 最佳实践

1. **模型定义**：通用模型放在 `internal/shared/models`
2. **迁移管理**：统一在 `admin/migrations` 中执行
3. **类型安全**：为 JSON 字段创建专门的类型和方法
4. **测试完备**：为每个 Repo 方法编写测试
5. **解耦设计**：`learn_en` 通过共享模型使用，不直接依赖 `admin`

## 未来扩展

如果需要为 `learn_en` 添加特定字段：

1. **方案一**：扩展共享模型（如果字段通用）
2. **方案二**：创建 `learn_en` 专属模型，包含共享模型作为嵌入字段
3. **方案三**：使用视图模型（ViewModel）进行转换

---

**作者**: AI Assistant  
**更新时间**: 2025-11-11
