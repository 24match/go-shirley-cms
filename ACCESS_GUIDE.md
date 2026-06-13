# 系统访问指南

## 系统架构

本系统采用 SaaS 多租户架构，包含三种角色：

1. **超级管理员 (Super Admin)** - 系统最高权限，管理所有租户
2. **租户管理员 (Tenant Admin)** - 单个租户的管理员，管理自己租户的内容
3. **普通用户 (Public User)** - 访问前台网站的用户

## 访问入口

### 1. 前台网站（公开访问）

**URL**: `http://localhost:8080`

无需登录，直接访问企业展示网站。

**功能**:
- 浏览企业信息（Banner、About、Products、Factory 等）
- 查看展会活动
- 提交联系表单
- 切换语言（中文/英文）

### 2. 租户管理后台

**URL**: `http://localhost:8080/admin`

**登录方式**:
```bash
POST http://localhost:8080/api/login
Content-Type: application/json

{
  "email": "tenant@example.com",
  "password": "password123"
}
```

**功能**:
- 模块管理 (`/api/admin/modules`)
- 图片管理 (`/api/admin/images`)
- 内容管理 (`/api/admin/content`)
- 语言管理 (`/api/admin/lang`)
- 站点设置 (`/api/admin/site-settings`)
- 联系表单管理 (`/api/admin/contact-submissions`)

### 3. 超级管理员后台

**URL**: `http://localhost:8080/swagger/index.html` (API 文档)

**登录方式**:
```bash
POST http://localhost:8080/api/login
Content-Type: application/json

{
  "email": "superadmin@example.com",
  "password": "admin123"
}
```

**功能**:
- 租户管理 (`/api/superadmin/tenants`)
- 系统统计 (`/api/superadmin/stats`)
- 租户激活/禁用
- 模拟租户登录

## API 端点总览

### 公开接口（无需认证）

| 方法 | 端点 | 描述 |
|------|------|------|
| POST | `/api/login` | 用户登录 |
| GET | `/api/public/config` | 获取页面配置 |
| GET | `/api/public/modules` | 获取模块配置 |
| GET | `/api/public/images` | 获取图片列表 |
| GET | `/api/public/content` | 获取内容项 |
| GET | `/api/public/lang` | 获取翻译文本 |
| GET | `/api/public/site-settings` | 获取站点设置 |
| POST | `/api/contact/submit` | 提交联系表单 |

### 租户后台接口（需要认证）

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/api/admin/modules` | 获取模块配置列表 |
| GET | `/api/admin/modules/:name` | 获取单个模块配置 |
| PUT | `/api/admin/modules/:name` | 保存模块配置 |
| DELETE | `/api/admin/modules/:name` | 删除模块配置 |
| POST | `/api/admin/images` | 上传图片 |
| GET | `/api/admin/images` | 获取图片列表 |
| DELETE | `/api/admin/images/:id` | 删除图片 |
| GET | `/api/admin/content` | 获取内容项列表 |
| POST | `/api/admin/content` | 创建内容项 |
| PUT | `/api/admin/content/:id` | 更新内容项 |
| DELETE | `/api/admin/content/:id` | 删除内容项 |
| GET | `/api/admin/lang` | 获取语言文本 |
| POST | `/api/admin/lang` | 创建语言文本 |
| PUT | `/api/admin/lang/:id` | 更新语言文本 |
| GET | `/api/admin/site-settings` | 获取站点设置 |
| POST | `/api/admin/site-settings` | 保存站点设置 |
| GET | `/api/admin/contact-submissions` | 获取联系表单提交 |
| DELETE | `/api/admin/contact-submissions/:id` | 删除联系表单提交 |

### 超级管理员接口（需要超级管理员权限）

| 方法 | 端点 | 描述 |
|------|------|------|
| POST | `/api/superadmin/tenants` | 创建租户 |
| GET | `/api/superadmin/tenants` | 获取租户列表 |
| GET | `/api/superadmin/tenants/:id` | 获取单个租户详情 |
| PUT | `/api/superadmin/tenants/:id` | 更新租户信息 |
| DELETE | `/api/superadmin/tenants/:id` | 删除租户 |
| POST | `/api/superadmin/tenants/:id/activate` | 激活租户 |
| POST | `/api/superadmin/tenants/:id/disable` | 禁用租户 |
| POST | `/api/superadmin/tenants/:id/impersonate` | 模拟租户登录 |
| GET | `/api/superadmin/stats` | 获取系统统计 |

## 认证方式

所有需要认证的接口使用 JWT Token 认证：

1. 调用 `/api/login` 获取 token
2. 在后续请求的 `Authorization` 头中携带 token：
   ```
   Authorization: Bearer <your_token_here>
   ```

## 快速测试

### 使用 curl 测试登录

```bash
# 租户登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"tenant@example.com","password":"password123"}'

# 超级管理员登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"superadmin@example.com","password":"admin123"}'
```

### 使用 Swagger UI 测试

1. 访问 `http://localhost:8080/swagger/index.html`
2. 点击 "Authorize" 按钮
3. 输入 JWT token
4. 测试各个 API 端点

## 默认账号

系统初始化后，建议创建以下默认账号：

| 角色 | Email | 密码 | 说明 |
|------|-------|------|------|
| 超级管理员 | `admin@system.com` | `admin123` | 系统最高权限 |
| 租户管理员 | `tenant@example.com` | `password123` | 示例租户账号 |

## 租户隔离

系统采用租户隔离机制：
- 每个租户有独立的数据库记录（通过 `tenant_id` 字段区分）
- 租户只能访问和管理自己的数据
- 超级管理员可以跨租户管理

## 前端 React 应用

新的 React 前端应用构建后位于 `./frontend` 目录：

```bash
cd frontend-react
npm install
npm run build
```

构建产物自动输出到 `../frontend` 目录，通过 Go 后端静态文件服务访问。