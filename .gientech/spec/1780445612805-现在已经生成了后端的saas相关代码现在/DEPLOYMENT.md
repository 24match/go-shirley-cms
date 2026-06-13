# 部署指南

## 前置条件

确保已安装以下工具：
- Node.js 18+ 
- Go 1.21+

## 步骤 1：安装前端依赖

```bash
cd frontend-react
npm install
```

## 步骤 2：构建前端

```bash
npm run build
```

构建产物将自动输出到 `../frontend` 目录。

## 步骤 3：启动后端服务

```bash
# 返回项目根目录
cd ..

# 启动 Go 后端
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

## 步骤 4：验证应用

1. 打开浏览器访问 `http://localhost:8080`
2. 验证页面正常加载
3. 测试语言切换功能（EN ↔ 中文）
4. 验证 CMS 数据加载
5. 测试联系表单提交

## 常见问题

### 构建失败

确保 node_modules 已正确安装：
```bash
rm -rf node_modules package-lock.json
npm install
npm run build
```

### API 请求失败

确保后端服务正在运行，并且 API 端点可访问。
检查浏览器控制台的网络请求。

### 样式不生效

清除浏览器缓存并刷新页面。

## 生产环境部署

### 1. 构建优化版本

```bash
cd frontend-react
npm run build
```

### 2. 部署到服务器

将构建产物（`../frontend` 目录）部署到 Web 服务器。

### 3. 配置反向代理

如果使用 Nginx，配置如下：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        root /path/to/frontend;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 环境变量

生产环境可配置以下环境变量：

- `VITE_API_BASE_URL` - API 基础 URL（默认：`/api`）