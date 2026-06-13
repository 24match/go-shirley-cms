go# Medical CMS Frontend - React 版本

将原有基于原生 JavaScript/HTML/CSS 的前端架构迁移到 React 18 + TypeScript + Vite。

## 技术栈

- **React 18** - UI 框架
- **TypeScript 5** - 类型系统
- **Vite 5** - 构建工具

## 项目结构

```
frontend-react/
├── src/
│   ├── components/
│   │   ├── common/        # 通用组件（Button, Loading, ErrorBoundary 等）
│   │   ├── layout/        # 布局组件（Header, Footer）
│   │   └── sections/      # 页面区块组件（Banner, About, Products 等）
│   ├── contexts/          # React Context（LanguageContext, CMSContext）
│   ├── services/          # API 服务和 i18n 服务
│   ├── hooks/             # 自定义 Hooks
│   ├── types/             # TypeScript 类型定义
│   ├── utils/             # 工具函数
│   ├── styles/            # CSS 样式
│   ├── App.tsx            # 根组件
│   └── main.tsx           # 入口文件
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 快速开始

### 1. 安装依赖

```bash
cd frontend-react
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

开发服务器将在 http://localhost:3000 启动，并自动代理 `/api` 请求到后端 `http://localhost:8080`。

### 3. 构建生产版本

```bash
npm run build
```

构建产物将输出到 `../frontend` 目录，与现有 Go 后端集成。

### 4. 预览生产构建

```bash
npm run preview
```

## 功能特性

- ✅ 响应式布局（支持移动端）
- ✅ 中英文语言切换
- ✅ CMS 数据动态加载
- ✅ 图片懒加载
- ✅ 表单验证
- ✅ 错误边界处理
- ✅ 加载状态管理

## 组件列表

### 布局组件
- `Header` - 导航栏、Logo、语言切换
- `Footer` - 页脚信息

### 页面区块
- `Banner` - 首页 Banner
- `Stats` - 统计数据
- `About` - 关于我们
- `Products` - 产品展示
- `Factory` - 工厂实力
- `Advantage` - 核心优势
- `Events` - 展会活动
- `Contact` - 联系方式和表单

### 通用组件
- `Button` - 按钮
- `Loading` - 加载指示器
- `ErrorBoundary` - 错误边界
- `LazyImage` - 懒加载图片
- `WhatsAppFloat` - WhatsApp 浮动按钮

## API 端点

前端需要以下后端 API 端点：

- `/api/public/config` - 页面配置
- `/api/public/modules` - 模块配置
- `/api/public/images` - 图片列表
- `/api/public/content` - 内容项
- `/api/public/site-settings` - 站点设置
- `/api/public/lang` - 翻译数据
- `/api/public/contact` - 联系表单提交

## 与现有后端集成

Vite 配置已将构建输出目录设置为 `../frontend`，运行 `npm run build` 后，构建产物将自动复制到后端静态文件目录。

启动 Go 后端服务后，即可通过 `http://localhost:8080` 访问 React 前端应用。