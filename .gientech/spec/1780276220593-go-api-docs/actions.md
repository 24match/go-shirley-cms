# Go API 接口文档生成实施待办清单

## 实施计划

本待办清单基于已批准的需求文档和设计文档制定，将 API 文档生成功能的实现拆分为独立的、可执行的编码待办项。

---

## 待办清单

- [x] 1. 安装 Swagger 相关依赖
  - 安装 swag CLI 工具：`go install github.com/swaggo/swag/cmd/swag@latest`
  - 添加 swaggo/swag 依赖：`go get -u github.com/swaggo/swag`
  - 添加 gin-swagger 中间件：`go get -u github.com/swaggo/gin-swagger`
  - 添加 swagger 静态文件：`go get -u github.com/swaggo/files`
  - _需求：REQ-003_

- [x] 2. 创建 docs 目录并生成初始文档
  - 在项目根目录创建 docs 文件夹
  - 运行 `swag init -g main.go -o ./docs` 生成初始文档
  - 验证生成的 docs/docs.go、docs/swagger.json、docs/swagger.yaml 文件
  - _需求：REQ-001_

- [x] 3. 配置 main.go 添加 Swagger 初始化注释
  - 在 main.go 文件顶部添加 Swagger 主文档注释
  - 添加 @title、@version、@description 等元数据
  - 添加 @host、@BasePath 等配置信息
  - 添加 @securityDefinitions.apikey JWT 认证配置
  - _需求：REQ-001, REQ-005_

- [x] 4. 在 routes.go 中注册 Swagger UI 中间件
  - 导入 gin-swagger 和 swagger files 包
  - 在 RegisterRoutes 函数中注册 Swagger 路由
  - 配置路由路径为 `/swagger/*any`
  - 添加环境判断，仅开发环境启用 Swagger
  - _需求：REQ-002, REQ-008_

- [x] 5. 为 backend/common/response.go 添加 Swagger 注释
  - 为 Response 结构体添加 @Description 注释
  - 为字段添加 JSON 示例和说明
  - 定义通用错误响应结构
  - _需求：REQ-004, REQ-006_

- [x] 6. 为 backend/models/models.go 添加 Swagger 模型注释
  - 为 User 结构体添加模型注释和字段说明
  - 为 Content 结构体添加模型注释和字段说明
  - 为 Module 结构体添加模型注释和字段说明
  - 为 Image 结构体添加模型注释和字段说明
  - 为 LanguageText 结构体添加模型注释和字段说明
  - 为 SiteSetting 结构体添加模型注释和字段说明
  - _需求：REQ-004_

- [x] 7. 为 backend/controllers/user_controller.go 添加 API 注释
  - 为 Login 方法添加 @Summary、@Description、@Tags 注释
  - 添加 @Param 请求参数定义
  - 添加 @Success 和 @Failure 响应定义
  - 添加 @Router 路由定义
  - _需求：REQ-001, REQ-005, REQ-006_

- [x] 8. 为 backend/controllers/module_controller.go 添加 API 注释
  - 为 GetPageConfig 方法添加 API 注释
  - 为 UpdatePageConfig 方法添加 API 注释
  - 为 GetModuleConfigs 方法添加 API 注释
  - 为 GetModuleConfig 方法添加 API 注释
  - 为 SaveModuleConfig 方法添加 API 注释
  - 为 DeleteModuleConfig 方法添加 API 注释
  - 为 DeleteModuleImage 方法添加 API 注释
  - 为 GetPublicPageConfig 方法添加 API 注释
  - 为 GetPublicModuleConfigs 方法添加 API 注释
  - _需求：REQ-001, REQ-004, REQ-006_

- [x] 9. 为 backend/controllers/content_controller.go 添加 API 注释
  - 为 GetContentItems 方法添加 API 注释
  - 为 CreateContentItem 方法添加 API 注释
  - 为 UpdateContentItem 方法添加 API 注释
  - 为 DeleteContentItem 方法添加 API 注释
  - 为 GetPublicContentItems 方法添加 API 注释
  - _需求：REQ-001, REQ-004, REQ-006_

- [x] 10. 为 backend/controllers/image_controller.go 添加 API 注释
  - 为 GetImages 方法添加 API 注释
  - 为 UploadImage 方法添加 API 注释（包含文件上传参数）
  - 为 UploadMultipleImages 方法添加 API 注释
  - 为 UpdateImage 方法添加 API 注释
  - 为 DeleteImage 方法添加 API 注释
  - 为 DeleteImageByFilename 方法添加 API 注释
  - _需求：REQ-001, REQ-004, REQ-006_

- [x] 11. 为 backend/controllers/language_controller.go 添加 API 注释
  - 为 GetLanguageTexts 方法添加 API 注释
  - 为 GetLanguageText 方法添加 API 注释
  - 为 CreateLanguageText 方法添加 API 注释
  - 为 UpdateLanguageText 方法添加 API 注释
  - 为 DeleteLanguageText 方法添加 API 注释
  - 为 GetLanguageTextVersions 方法添加 API 注释
  - 为 RestoreLanguageTextVersion 方法添加 API 注释
  - 为 GetPublicLanguageTexts 方法添加 API 注释
  - _需求：REQ-001, REQ-004, REQ-006_

- [x] 12. 为 backend/controllers/site_setting_controller.go 添加 API 注释
  - 为 GetSettings 方法添加 API 注释
  - 为 SaveSettings 方法添加 API 注释
  - 为 GetPublicSettings 方法添加 API 注释
  - _需求：REQ-001, REQ-004, REQ-006_

- [x] 13. 为 backend/controllers/auth.go 中间件添加安全注释
  - 在需要认证的接口添加 @Security APIKey 注释
  - 公开接口明确标注无需认证
  - _需求：REQ-005_

- [x] 14. 重新生成并验证 Swagger 文档
  - 运行 `swag init -g main.go -o ./docs --parseDependency`
  - 运行 `swag fmt` 格式化注释
  - 启动应用访问 `/swagger/index.html` 验证界面
  - 访问 `/swagger/doc.json` 验证 JSON 文档有效性
  - _需求：REQ-001, REQ-002, REQ-007_

- [x] 15. 配置环境隔离和构建优化
  - 在 config 中添加 Swagger 启用配置
  - 生产环境自动禁用 Swagger 中间件
  - 在 Makefile 或构建脚本中添加文档生成命令
  - _需求：REQ-008_

- [x] 16. 编写 Swagger 注释编写规范文档
  - 创建 docs/README.md 说明注释规范
  - 提供常用注释模板示例
  - 记录常见问题和解决方案
  - _需求：REQ-002_

---

## 执行顺序说明

### 第一阶段：基础设置（待办 1-2）
- 安装必要的依赖包
- 创建文档生成基础设施

### 第二阶段：核心配置（待办 3-4）
- 配置 main.go 入口文件
- 注册 Swagger 中间件

### 第三阶段：模型注释（待办 5-6）
- 为统一响应结构添加注释
- 为所有数据模型添加注释

### 第四阶段：API 注释（待办 7-13）
- 按控制器逐个添加 API 注释
- 区分公开接口和认证接口

### 第五阶段：验证优化（待办 14-16）
- 生成并验证最终文档
- 配置环境隔离
- 编写使用文档

---

## 质量标准

1. **完整性**：所有 API 端点都有对应的 Swagger 注释
2. **准确性**：注释描述的参数、响应与实际代码一致
3. **可测试性**：通过 Swagger UI 可直接测试所有接口
4. **可维护性**：注释格式统一，易于后续更新
5. **安全性**：生产环境默认禁用 Swagger UI

---

## 验收标准

- [ ] 访问 `http://localhost:8080/swagger/index.html` 能正常展示 Swagger UI
- [ ] 所有 API 接口在文档中都有对应条目
- [ ] 每个接口都有清晰的摘要、描述和标签
- [ ] 请求参数和响应结构完整展示
- [ ] 认证接口正确标注 JWT 要求
- [ ] 通过 Swagger UI 能成功调用 API
- [ ] 生产环境构建时 Swagger 自动禁用