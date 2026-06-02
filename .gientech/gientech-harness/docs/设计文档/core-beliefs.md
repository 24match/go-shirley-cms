# 核心设计信念

## 当前可验证的设计信念

- **信念 1：轻量级内嵌数据库优先**  
  项目采用 SQLite 作为唯一数据源，通过 GORM 自动迁移管理表结构，无需额外部署数据库服务。数据库初始化在 [`main.go`](../../../../main.go) 的 `config.InitDB()` 中完成，支持 7 个核心模型的自动建表（User、Image、PageConfig、ModuleConfig、ContentItem、LanguageText、LanguageTextVersion），证据见 [`backend/config/database.go`](../../../../backend/config/database.go) 的 `AutoMigrate` 调用。

- **信念 2：Controller-Service 分层明确**  
  业务逻辑严格遵循三层架构：Controller 负责 HTTP 请求解析与响应封装，Service 封装核心业务规则，Model 定义数据结构。例如 [`ImageController`](../../../../backend/controllers/image_controller.go) 仅处理参数提取和文件上传，实际数据操作委托给 [`ImageService`](../../../../backend/services/image_service.go)。

- **信念 3：公开接口与管理接口物理隔离**  
  路由在 [`backend/routes/routes.go`](../../../../backend/routes/routes.go) 中按访问权限分为两组：`/api/public/*` 无需认证，`/api/admin/*` 必须通过 `AuthMiddleware()` 中间件。这种分离确保前台展示与后台管理的安全边界。

- **信念 4：配置即数据，支持运行时动态修改**  
  页面配置（Banner、About、Products 等）以 JSON 形式存储在 `PageConfig` 表中，而非硬编码在配置文件。管理员可通过 API 实时更新配置，无需重启服务，证据见 [`database.go`](../../../../backend/config/database.go) 的 `initDefaultPageConfig` 函数。

- **信念 5：多语言支持内建于数据模型**  
  `ModuleConfig`、`ContentItem`、`LanguageText` 等模型均包含 `Zh*`/`En*` 成对字段，支持中英文双语内容存储。`LanguageTextVersion` 表提供版本追溯能力，确保多语言文本的变更可审计。

- **信念 6：文件上传采用本地存储策略**  
  上传的图片直接保存至 `uploads/` 目录，文件路径记录在数据库 `Image.FilePath` 字段。删除操作同步清理物理文件，证据见 [`ImageService.DeleteImage`](../../../../backend/services/image_service.go)。

- **信念 7：统一响应格式与错误处理**  
  所有 API 响应通过 `common` 包的 `JSONSuccess`、`JSONError` 等函数封装，确保返回结构一致。全局错误处理中间件在 [`backend/middleware/error_handler.go`](../../../../backend/middleware/error_handler.go) 中定义。

## 这些信念如何指导后续工作

- **新增功能**：应遵循现有分层模式，新增 Controller 时必须配套 Service 实现，不得在 Controller 中直接操作数据库。
- **跨模块协作**：模块间通过 Service 层交互，避免 Controller 直接调用其他 Controller。
- **数据模型扩展**：新增表结构需在 `InitDB()` 的 `AutoMigrate` 列表中注册，并考虑是否需要多语言字段支持。
- **配置管理**：新增可配置项应优先采用 `PageConfig` 的 JSON 存储模式，而非新增配置文件。
- **文档沉淀**：新增 API 接口需在路由注册处保持 `public/admin` 分类一致性，并在 README 的 API 表格中补充说明。