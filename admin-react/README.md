# Medical CMS - 管理后台 (React 版)

将原有的管理后台改造为 React 实现。

## 技术栈

- **React 18** - UI 框架
- **TypeScript 5** - 类型系统
- **React Router 6** - 路由管理
- **Axios** - HTTP 客户端
- **Vite 5** - 构建工具

## 项目结构

```
admin-react/
├── src/
│   ├── components/
│   │   └── layout/
│   │       └── AdminLayout.tsx    # 管理后台布局
│   ├── contexts/
│   │   └── AuthContext.tsx        # 认证上下文
│   ├── pages/
│   │   ├── Login.tsx              # 登录页面
│   │   ├── Dashboard.tsx          # 仪表盘
│   │   ├── modules/
│   │   │   └── ModuleList.tsx     # 模块管理
│   │   └── tenants/
│   │       └── TenantList.tsx     # 租户管理（超级管理员）
│   ├── services/
│   │   └── api.ts                 # API 服务层
│   ├── types/
│   │   └── index.ts               # TypeScript 类型定义
│   ├── styles/
│   │   └── admin.css              # 管理后台样式
│   ├── App.tsx                    # 根组件
│   └── main.tsx                   # 入口文件
├── package.json
├── tsconfig.json
├── vite.config.ts
└── index.html
```

## 快速开始

### 1. 安装依赖

```bash
cd admin-react
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

开发服务器将在 http://localhost:3001 启动。

### 3. 构建生产版本

```bash
npm run build
```

构建产物将输出到 `../admin` 目录。

## 功能特性

### 租户管理员功能

- **仪表盘** - 查看模块概览和最近联系表单
- **模块管理** - 查看、编辑、启用/禁用模块
- **图片管理** - 上传、预览、删除图片
- **内容管理** - 管理内容项
- **语言管理** - 管理翻译文本
- **站点设置** - 配置站点信息
- **联系表单** - 查看和管理联系表单提交

### 超级管理员功能

- **仪表盘** - 查看系统统计
- **租户管理** - 创建、编辑、删除、激活/禁用租户
- **系统统计** - 查看全局统计数据

## 默认账号

| 角色 | Email | 密码 |
|------|-------|------|
| 租户管理员 | tenant@example.com | password123 |
| 超级管理员 | admin@system.com | admin123 |

## API 端点

管理后台使用以下 API 端点：

### 认证
- `POST /api/login` - 用户登录

### 租户管理员
- `GET /api/admin/modules` - 获取模块列表
- `PUT /api/admin/modules/:name` - 更新模块配置
- `DELETE /api/admin/modules/:name` - 删除模块
- `GET /api/admin/images` - 获取图片列表
- `POST /api/admin/images` - 上传图片
- `GET /api/admin/content` - 获取内容列表
- `GET /api/admin/lang` - 获取语言文本
- `GET /api/admin/site-settings` - 获取站点设置
- `GET /api/admin/contact-submissions` - 获取联系表单

### 超级管理员
- `GET /api/superadmin/tenants` - 获取租户列表
- `POST /api/superadmin/tenants` - 创建租户
- `PUT /api/superadmin/tenants/:id` - 更新租户
- `DELETE /api/superadmin/tenants/:id` - 删除租户
- `POST /api/superadmin/tenants/:id/activate` - 激活租户
- `POST /api/superadmin/tenants/:id/disable` - 禁用租户
- `GET /api/superadmin/stats` - 获取系统统计

## 路由说明

| 路径 | 说明 | 权限要求 |
|------|------|---------|
| `/login` | 登录页面 | 公开 |
| `/dashboard` | 仪表盘 | 已认证 |
| `/modules` | 模块管理 | 租户管理员 |
| `/tenants` | 租户管理 | 超级管理员 |

## 与后端集成

Vite 配置已将构建输出目录设置为 `../admin`，运行 `npm run build` 后，构建产物将自动复制到后端静态文件目录。

启动 Go 后端服务后，即可通过 `http://localhost:8080/admin` 访问管理后台。