# 医疗设备官网 CMS 系统

一个功能完整的医疗设备官网内容管理系统，基于 Go + Gin + SQLite 构建。

## 功能特性

### 后端功能
- **用户认证**: JWT Token 认证，支持登录/登出
- **权限管理**: 基于用户角色的访问控制（超级管理员/租户管理员）
- **多租户支持**: 完整的 SaaS 多租户架构，支持租户隔离
- **租户配置管理**: 
  - 功能模块开关配置（图片管理、页面配置、多语言支持等）
  - 资源配额管理（图片数量、存储空间、内容项数量等）
  - 订阅计划管理（免费版/专业版/企业版）
- **配额检查**: 自动配额检查和限制，防止资源超限
- **图片管理**: 支持图片上传、分类、删除
- **页面配置**: 动态管理 Banner、关于我们、产品展示等内容
- **SQLite 内嵌数据库**: 无需额外配置数据库

### 前台功能
- **响应式官网**: 美观的医疗设备展示页面
- **动态内容**: 实时加载 CMS 配置和图片
- **产品展示**: 支持自定义产品图片和描述

### 后台功能
- **直观管理界面**: 简洁易用的管理后台
- **图片管理**: 上传新图片，按分类管理
- **页面配置**: 无需修改代码即可修改网站内容
- **租户管理**（超级管理员）:
  - 创建/编辑/删除租户
  - 配置租户功能模块
  - 设置资源配额
  - 管理订阅计划
  - 重置配额使用统计
- **租户 Dashboard**:
  - 功能模块状态展示
  - 资源配额使用情况可视化
  - 配额警告提示

## 项目结构

```
lihui/
├── main.go              # Go 后端主程序
├── go.mod              # Go 模块依赖
├── inde.html           # 官网首页
├── admin/
│   └── index.html      # CMS 后台管理界面
├── uploads/            # 上传的图片存储目录
├── medical.db          # SQLite 数据库文件（自动生成）
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 3. 访问页面

- **官网首页**: http://localhost:8080
- **CMS 后台**: http://localhost:8080/admin

### 4. 默认登录信息

- 用户名: `admin`
- 密码: `admin123`

## 租户配置功能

### 功能模块
系统支持以下功能模块的开关配置：

| 模块 | 说明 |
|------|------|
| image_management | 图片管理 |
| page_config | 页面配置 |
| multi_language | 多语言支持 |
| contact_form | 联系表单 |
| content_management | 内容管理 |
| analytics | 数据分析 |
| seo_tools | SEO 工具 |

### 资源配额
支持以下资源类型的配额管理：

| 资源 | 说明 | 默认值（免费版） |
|------|------|-----------------|
| max_images | 图片数量 | 50 |
| max_storage_mb | 存储空间 (MB) | 512 |
| max_content_items | 内容项数量 | 50 |
| max_users | 用户数量 | 5 |
| max_modules | 模块数量 | 10 |
| max_languages | 语言数量 | 2 |

### 订阅计划
| 计划 | 说明 |
|------|------|
| free | 免费版 - 基础功能，有限配额 |
| pro | 专业版 - 全部功能，更高配额 |
| enterprise | 企业版 - 无限制配额 |

## API 接口

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/public/config | 获取页面配置 |
| GET | /api/public/images | 获取公开图片 |
| POST | /api/login | 用户登录 |

### 管理接口（需要认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/admin/images | 获取所有图片 |
| POST | /api/admin/images | 上传图片 |
| DELETE | /api/admin/images/:id | 删除图片 |
| GET | /api/admin/config | 获取页面配置 |
| PUT | /api/admin/config | 更新页面配置 |

### 超级管理员接口（需要超级管理员权限）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/superadmin/tenants | 获取租户列表 |
| POST | /api/superadmin/tenants | 创建租户 |
| GET | /api/superadmin/tenants/:id | 获取租户详情 |
| PUT | /api/superadmin/tenants/:id | 更新租户 |
| DELETE | /api/superadmin/tenants/:id | 删除租户 |
| GET | /api/superadmin/tenants/:id/config | 获取租户配置 |
| PUT | /api/superadmin/tenants/:id/config | 更新租户配置 |
| POST | /api/superadmin/tenants/:id/quota/reset | 重置配额 |
| GET | /api/superadmin/tenants/:id/quota/usage | 获取配额使用情况 |

### 租户接口（需要租户管理员权限）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/tenant/config | 获取当前租户配置 |
| GET | /api/tenant/features | 获取功能模块状态 |
| GET | /api/tenant/quota | 获取配额使用情况 |

## 使用说明

### 图片管理
1. 登录后台后进入「图片管理」
2. 选择图片文件，填写描述和分类
3. 点击上传，图片将自动显示在列表中
4. 支持按分类筛选，可随时删除图片

### 页面配置
1. 进入「页面配置」页面
2. 修改各个模块的内容（标题、描述、开关等）
3. 点击「保存」按钮即可生效
4. 刷新官网首页即可看到更新后的内容

## 技术栈

- **后端**: Go + Gin
- **数据库**: SQLite (GORM)
- **认证**: JWT
- **前端**: 原生 HTML/CSS/JS

## 扩展开发

### 添加新的页面配置

在 `initDefaultPageConfig` 函数中添加新的默认配置项，然后在后台管理界面中添加对应的表单即可。

### 添加新的图片分类

在后台管理页面的下拉框中添加新的分类选项，后端会自动支持。