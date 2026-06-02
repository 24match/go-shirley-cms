# 数据结构说明

## 存储角色概览

本项目使用 **SQLite** 作为嵌入式关系型数据库，通过 **GORM** ORM 框架进行数据访问和 schema 管理。

| 存储类型 | 文件位置 | 职责说明 |
|----------|----------|----------|
| SQLite 数据库 | `medical.db`（项目根目录） | 存储用户、图片、页面配置、模块配置、内容项、多语言文本等核心业务数据 |
| 文件存储 | `uploads/` 目录 | 存储上传的图片文件，数据库仅记录文件路径和元数据 |

数据库在应用启动时通过 GORM 的 `AutoMigrate` 自动初始化，无需手动执行 SQL 脚本。

## 关键存储实体

### 用户表（users）

对应模型：[User](../../../../backend/models/models.go)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ID | uint | 主键 | 自增 ID |
| CreatedAt | time.Time | - | 创建时间 |
| UpdatedAt | time.Time | - | 更新时间 |
| DeletedAt | gorm.DeletedAt | 索引 | 软删除时间 |
| Username | string | 唯一索引，非空 | 用户名 |
| Password | string | 非空 | 加密后的密码（bcrypt） |
| Role | string | 默认 'admin' | 用户角色 |

### 图片表（images）

对应模型：[Image](../../../../backend/models/models.go)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ID | uint | 主键 | 自增 ID |
| CreatedAt | time.Time | - | 创建时间 |
| UpdatedAt | time.Time | - | 更新时间 |
| DeletedAt | gorm.DeletedAt | 索引 | 软删除时间 |
| Filename | string | 非空 | 原始文件名 |
| FilePath | string | 非空 | 存储路径 |
| FileSize | int64 | - | 文件大小（字节） |
| Description | string | - | 简短描述 |
| LongDescription | string | - | 详细描述 |
| Category | string | - | 分类 |
| SortOrder | int | 默认 0 | 排序顺序 |

### 页面配置表（page_configs）

对应模型：[PageConfig](../../../../backend/models/models.go)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ID | uint | 主键 | 自增 ID |
| CreatedAt | time.Time | - | 创建时间 |
| UpdatedAt | time.Time | - | 更新时间 |
| DeletedAt | gorm.DeletedAt | 索引 | 软删除时间 |
| PageName | string | 唯一索引，非空 | 页面标识（如 banner、about、products） |
| ConfigData | string | text 类型 | JSON 格式的配置数据 |

### 模块配置表（module_configs）

对应模型：[ModuleConfig](../../../../backend/models/models.go)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ID | uint | 主键 | 自增 ID |
| CreatedAt | time.Time | - | 创建时间 |
| UpdatedAt | time.Time | - | 更新时间 |
| DeletedAt | gorm.DeletedAt | 索引 | 软删除时间 |
| ModuleName | string | 唯一索引，非空 | 模块名称 |
| Enabled | bool | 默认 true | 是否启用 |
| ZhTitle / EnTitle | string | - | 中/英文标题 |
| ZhSubtitle / EnSubtitle | string | - | 中/英文副标题 |
| ZhContent / EnContent | string | text 类型 | 中/英文内容 |
| Title / Subtitle / Content | string | - | 向后兼容字段 |
| ImagePath | string | - | 关联图片路径 |
| SortOrder | int | 默认 0 | 排序顺序 |
| ExtraData | string | text 类型 | 扩展数据（JSON） |
| ZhDescription / EnDescription | string | - | 中/英文描述 |
| Description | string | - | 向后兼容字段 |

### 内容项表（content_items）

对应模型：[ContentItem](../../../../backend/models/models.go)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ID | uint | 主键 | 自增 ID |
| CreatedAt | time.Time | - | 创建时间 |
| UpdatedAt | time.Time | - | 更新时间 |
| DeletedAt | gorm.DeletedAt | 索引 | 软删除时间 |
| Section | string | 非空 | 所属版块 |
| ZhTitle / EnTitle | string | - | 中/英文标题 |
| ZhDescription / EnDescription | string | - | 中/英文描述 |
| Title / Description | string | - | 向后兼容字段 |
| Icon | string | - | 图标 |
| ImagePath | string | - | 关联图片路径 |
| SortOrder | int | 默认 0 | 排序顺序 |

### 多语言文本表（language_texts）

对应模型：[LanguageText](../../../../backend/models/models.go)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ID | uint | 主键 | 自增 ID |
| CreatedAt | time.Time | - | 创建时间 |
| UpdatedAt | time.Time | - | 更新时间 |
| DeletedAt | gorm.DeletedAt | 索引 | 软删除时间 |
| Key | string | 唯一索引，非空 | 文本键 |
| Module | string | 非空 | 所属模块 |
| EnText | string | text 类型 | 英文文本 |
| ZhText | string | text 类型 | 中文文本 |
| Description | string | - | 描述说明 |
| Version | int | 默认 1 | 版本号 |

### 多语言文本版本历史表（language_text_versions）

对应模型：[LanguageTextVersion](../../../../backend/models/models.go)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| ID | uint | 主键 | 自增 ID |
| CreatedAt | time.Time | - | 创建时间 |
| UpdatedAt | time.Time | - | 更新时间 |
| DeletedAt | gorm.DeletedAt | 索引 | 软删除时间 |
| LanguageTextID | uint | 非空 | 关联 language_texts.ID |
| Key | string | 非空 | 文本键 |
| Module | string | 非空 | 所属模块 |
| EnText | string | text 类型 | 英文文本 |
| ZhText | string | text 类型 | 中文文本 |
| Description | string | - | 描述说明 |
| Version | int | 非空 | 版本号 |
| UpdatedAt | time.Time | - | 版本更新时间 |

## 迁移与演进方式

### Schema 初始化

数据库 schema 通过 GORM 的 `AutoMigrate` 在应用启动时自动创建或更新：

```go
// 来源：backend/config/database.go
err = DB.AutoMigrate(
    &models.User{},
    &models.Image{},
    &models.PageConfig{},
    &models.ModuleConfig{},
    &models.ContentItem{},
    &models.LanguageText{},
    &models.LanguageTextVersion{},
)
```

### 默认数据初始化

系统启动时会自动初始化以下默认数据（如果不存在）：

1. **默认管理员用户**
   - 用户名：`admin`
   - 密码：`admin123`（bcrypt 加密存储）
   - 角色：`admin`

2. **默认页面配置**
   - `banner`：横幅配置（启用状态、标题、副标题）
   - `about`：关于我们配置（启用状态、标题、内容）
   - `products`：产品展示配置（启用状态、标题、显示数量）

### 升级策略

- **新增字段**：在模型结构体中添加字段，重启应用后 GORM 会自动添加新列
- **修改约束**：部分约束（如唯一索引）的修改可能需要手动干预
- **数据迁移**：复杂的数据迁移需编写自定义迁移逻辑

## 代码入口

| 文件 | 说明 |
|------|------|
| [backend/config/database.go](../../../../backend/config/database.go) | 数据库初始化、AutoMigrate、默认数据初始化 |
| [backend/models/models.go](../../../../backend/models/models.go) | 所有数据模型定义 |

## 注意事项

1. **软删除机制**：所有模型均嵌入 `gorm.DeletedAt` 字段，删除操作默认为软删除，查询时需注意 `Unscoped()` 的使用。

2. **向后兼容字段**：`ModuleConfig` 和 `ContentItem` 包含多组中英文字段（如 `Title`/`ZhTitle`/`EnTitle`），新开发建议优先使用带语言前缀的字段。

3. **JSON 配置数据**：`PageConfig.ConfigData` 和 `ModuleConfig.ExtraData` 以 JSON 字符串形式存储，读写时需手动序列化/反序列化。

4. **多语言版本控制**：`LanguageText` 与 `LanguageTextVersion` 配合使用，支持文本变更的历史追溯，具体版本管理逻辑待结合服务层实现进一步确认。

5. **文件与数据库分离**：图片文件存储在 `uploads/` 目录，数据库仅记录元数据；删除图片时需同步清理文件和数据库记录。

## 本文档引用的文件

- [README.md](../../../../README.md) - 项目概述
- [backend/config/database.go](../../../../backend/config/database.go) - 数据库配置与迁移逻辑
- [backend/models/models.go](../../../../backend/models/models.go) - 数据模型定义