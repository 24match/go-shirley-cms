# Swagger API 文档

## 概述

本项目使用 **swaggo/swag** 生成 OpenAPI 2.0 (Swagger) 规范的 API 文档，提供类似 Java Swagger + Knife4j 的交互式文档界面。

## 快速开始

### 1. 安装依赖

```bash
# 安装 swag CLI 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 添加 Go 依赖
go get -u github.com/swaggo/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 2. 生成文档

```bash
# 生成 Swagger 文档
~/go/bin/swag init -g main.go -o ./docs --parseDependency

# 格式化注释（可选）
~/go/bin/swag fmt
```

### 3. 访问文档

启动应用后，访问：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Swagger JSON**: http://localhost:8080/swagger/doc.json
- **Swagger YAML**: http://localhost:8080/swagger/doc.yaml

## Swagger 注释规范

### 控制器方法注释

```go
// @Summary 接口摘要
// @Description 接口详细描述
// @Tags 标签/分组
// @Accept json
// @Produce json
// @Param param_name param_type required "description"
// @Success 200 {object} ResponseType
// @Failure 400 {object} ErrorResponseType
// @Security APIKey
// @Router /api/endpoint [method]
func (c *Controller) Method(ctx *gin.Context) {
    // ...
}
```

### 参数类型

| 类型 | 说明 | 示例 |
|------|------|------|
| query | URL 查询参数 | `?page=1` |
| path | URL 路径参数 | `/api/users/{id}` |
| formData | 表单提交参数 | `multipart/form-data` |
| body | JSON 请求体 | `{"name": "test"}` |
| header | HTTP 请求头 | `Authorization: Bearer xxx` |

### 模型注释

```go
// ResponseType 响应结构
// @Description 结构体描述
type ResponseType struct {
    // 字段说明
    Field string `json:"field" example:"示例值"`
}
```

### 常用注释标签

| 标签 | 说明 |
|------|------|
| @Summary | 接口摘要 |
| @Description | 接口详细描述 |
| @Tags | 接口分组标签 |
| @Accept | 请求格式（json, multipart/form-data） |
| @Produce | 响应格式（json, xml） |
| @Param | 请求参数定义 |
| @Success | 成功响应定义 |
| @Failure | 失败响应定义 |
| @Security | 安全认证要求 |
| @Router | 路由定义 |
| @Description | 模型/字段描述 |

## API 接口分类

### 公开接口（无需认证）

| 接口 | 说明 |
|------|------|
| GET /api/public/config | 获取页面配置 |
| GET /api/public/modules | 获取模块配置 |
| GET /api/public/images | 获取图片列表 |
| GET /api/public/content | 获取内容项 |
| GET /api/public/lang | 获取多语言文本 |
| GET /api/public/site-settings | 获取站点设置 |

### 管理接口（需要 JWT 认证）

#### 用户管理

| 接口 | 说明 |
|------|------|
| POST /api/login | 用户登录 |

#### 模块管理

| 接口 | 说明 |
|------|------|
| GET /api/admin/config | 获取页面配置 |
| PUT /api/admin/config | 更新页面配置 |
| GET /api/admin/modules | 获取模块配置列表 |
| GET /api/admin/modules/{name} | 获取单个模块配置 |
| POST /api/admin/modules | 保存模块配置 |
| PUT /api/admin/modules/{name} | 更新模块配置 |
| DELETE /api/admin/modules/{name} | 删除模块配置 |
| DELETE /api/admin/modules/{name}/image | 删除模块图片 |

#### 内容管理

| 接口 | 说明 |
|------|------|
| GET /api/admin/content | 获取内容项列表 |
| POST /api/admin/content | 创建内容项 |
| PUT /api/admin/content/{id} | 更新内容项 |
| DELETE /api/admin/content/{id} | 删除内容项 |

#### 图片管理

| 接口 | 说明 |
|------|------|
| GET /api/admin/images | 获取图片列表 |
| POST /api/admin/images | 上传单张图片 |
| POST /api/admin/images/batch | 上传多张图片 |
| PUT /api/admin/images/{id} | 更新图片信息 |
| DELETE /api/admin/images/{id} | 删除图片 |
| DELETE /api/admin/images/by-filename/{filename} | 根据文件名删除图片 |

#### 多语言管理

| 接口 | 说明 |
|------|------|
| GET /api/admin/lang | 获取多语言文本列表 |
| GET /api/admin/lang/{id} | 获取单个多语言文本 |
| POST /api/admin/lang | 创建多语言文本 |
| PUT /api/admin/lang/{id} | 更新多语言文本 |
| DELETE /api/admin/lang/{id} | 删除多语言文本 |
| GET /api/admin/lang/{id}/versions | 获取版本历史 |
| POST /api/admin/lang/{id}/restore/{version} | 恢复历史版本 |

#### 站点设置

| 接口 | 说明 |
|------|------|
| GET /api/admin/site-settings | 获取站点设置 |
| POST /api/admin/site-settings | 保存站点设置 |

## 认证说明

### JWT Token 认证

所有管理接口需要在请求头中携带 JWT Token：

```
Authorization: Bearer {token}
```

### 获取 Token

调用 `POST /api/login` 接口获取 Token：

```json
{
  "username": "admin",
  "password": "password"
}
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin"
    }
  }
}
```

## 响应格式

所有接口统一返回格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": 1717200000
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 405 | 请求方法不允许 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |
| 501 | 数据库错误 |
| 502 | 数据验证错误 |
| 503 | 服务不可用 |

## 常见问题

### Q: 如何更新文档？

修改代码注释后，重新运行：
```bash
~/go/bin/swag init -g main.go -o ./docs --parseDependency
```

### Q: 如何在 Swagger UI 中测试需要认证的接口？

1. 先调用登录接口获取 Token
2. 点击 Swagger UI 右上角的 "Authorize" 按钮
3. 在 Value 中输入 `Bearer {token}`
4. 点击 "Authorize" 确认

### Q: 生产环境如何禁用 Swagger？

在 `backend/routes/routes.go` 中，Swagger 中间件注册时添加环境判断：

```go
if config.Env == "development" {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

## 参考链接

- [swaggo/swag GitHub](https://github.com/swaggo/swag)
- [gin-swagger GitHub](https://github.com/swaggo/gin-swagger)
- [Swagger 2.0 规范](https://swagger.io/specification/v2/)
- [OpenAPI Initiative](https://openapis.org/)