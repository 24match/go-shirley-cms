# Contact 模块配置联动测试指南

## 1. 概述

本文档描述 Contact（联系我们）模块的配置联动机制，包括数据流、调试步骤和常见问题排查指南。

---

## 2. 配置联动数据流

### 2.1 完整数据流图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Contact 模块配置联动数据流                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  后台管理界面  │────▶│  后端 API    │────▶│   数据库     │────▶│  后端 API    │
│  admin/     │     │  /api/admin/ │     │  ModuleConfig│     │  /api/public/│
│  admin.js   │     │  modules     │     │  .ExtraData  │     │  modules     │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
       │                    │                    │                    │
       │ 1. 填写联系信息     │ 2. 解析 FormData   │ 3. 存储 JSON      │ 4. 查询 enabled=true
       │    email, phone    │    构建 extraData  │    到 ExtraData   │    的模块
       │    whatsapp,       │    Map             │    字段           │
       │    address         │                    │                    │
       │                    │                    │                    │
       ▼                    ▼                    ▼                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              mergeExtraData() 展开字段                        │
│                         email, phone, whatsapp, address                      │
│                         从 ExtraData 展开到返回结果顶层                        │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  前台官网     │◀────│  前端 CMS    │◀────│  前端组件     │
│  显示联系信息  │     │  数据服务     │     │  Contact.js  │
│  #contactEmail│     │  cmsService  │     │  loadContact │
│  #contactPhone│     │  .js         │     │  Content()   │
└──────────────┘     └──────────────┘     └──────────────┘
       ▲                    ▲                    ▲
       │ 7. 更新 DOM        │ 6. 获取模块数据     │ 5. 调用 API
       │    显示配置值       │    以 moduleName    │    /api/public/
       │                    │    为键组织         │    modules
```

### 2.2 关键数据格式

#### 后台保存请求（FormData）
```
POST /api/admin/modules
Content-Type: multipart/form-data

moduleName: contact
enabled: true
email: info@company.com
phone: +86 123 4567 8900
whatsapp: +1 234 567 8900
address: No. 123, Industrial Park, City, Country
```

#### 数据库存储（ModuleConfig 表）
```sql
INSERT INTO module_configs (module_name, enabled, extra_data, ...)
VALUES ('contact', 1, '{"email":"info@company.com","phone":"+86 123 4567 8900","whatsapp":"+1 234 567 8900","address":"No. 123, Industrial Park, City, Country"}', ...);
```

#### 公开 API 响应
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "moduleName": "contact",
      "enabled": true,
      "email": "info@company.com",
      "phone": "+86 123 4567 8900",
      "whatsapp": "+1 234 567 8900",
      "address": "No. 123, Industrial Park, City, Country",
      "extraData": "{\"email\":\"info@company.com\",\"phone\":\"+86 123 4567 8900\",\"whatsapp\":\"+1 234 567 8900\",\"address\":\"No. 123, Industrial Park, City, Country\"}",
      ...
    }
  ]
}
```

---

## 3. 调试步骤和检查点

### 3.1 后台配置保存检查

**步骤 1：打开浏览器开发者工具**
- 按 F12 打开开发者工具
- 切换到 Network（网络）标签

**步骤 2：保存 Contact 配置**
- 登录后台管理界面
- 进入"联系我们配置"页面
- 填写联系信息
- 点击保存

**步骤 3：检查网络请求**
- 查找 `POST /api/admin/modules` 请求
- 检查请求头：`Content-Type: multipart/form-data`
- 检查请求体（Payload）：
  ```
  moduleName: contact
  enabled: true
  email: [填写的值]
  phone: [填写的值]
  whatsapp: [填写的值]
  address: [填写的值]
  ```

**预期结果**：
- 请求成功，响应状态码 200
- 响应数据包含保存的配置

### 3.2 数据库验证检查

**步骤 1：使用 SQLite 工具查看数据库**
```bash
sqlite3 medical.db
```

**步骤 2：查询 Contact 模块配置**
```sql
SELECT id, module_name, enabled, extra_data, created_at, updated_at 
FROM module_configs 
WHERE module_name = 'contact';
```

**预期结果**：
- 存在 `module_name = 'contact'` 的记录
- `enabled = 1`（true）
- `extra_data` 包含 JSON 格式的联系信息

### 3.3 公开 API 验证检查

**步骤 1：使用 curl 或浏览器测试 API**
```bash
curl http://localhost:8080/api/public/modules
```

**步骤 2：检查响应数据**
- 查找 `moduleName: "contact"` 的模块
- 验证以下字段在响应顶层（不在 extraData 内）：
  - `email`
  - `phone`
  - `whatsapp`
  - `address`
  - `enabled`

**预期结果**：
- 只返回 `enabled = true` 的模块
- contact 模块的联系信息字段在顶层可直接访问

### 3.4 前台显示验证检查

**步骤 1：打开浏览器开发者工具**
- 按 F12 打开开发者工具
- 切换到 Console（控制台）标签

**步骤 2：访问前台官网**
- 打开 `http://localhost:8080/`
- 滚动到 Contact（联系我们）区域

**步骤 3：检查控制台日志**
```
[Contact] Module config: {moduleName: "contact", enabled: true, email: "...", ...}
[Contact] Parsed data from top-level fields: {email: "...", phone: "...", ...}
[Contact] Final contact data: {email: "...", phone: "...", ...}
```

**步骤 4：检查页面元素**
- 检查 `#contactEmail` 元素是否显示配置的 email
- 检查 `#contactPhone` 元素是否显示配置的 phone
- 检查 `#contactWhatsApp` 元素是否显示配置的 whatsapp
- 检查 `#contactAddress` 元素是否显示配置的 address

**预期结果**：
- 控制台显示模块配置和解析数据
- 页面显示后台配置的联系信息

---

## 4. 常见问题排查指南

### 问题 1：前台不显示联系信息

**可能原因**：
1. Contact 模块 `enabled = false`
2. 数据库中无 Contact 模块配置
3. 前端 JS 加载失败

**排查步骤**：
1. 检查控制台是否有 `[Contact] Module config: undefined` 日志
2. 检查 API 响应中是否有 contact 模块
3. 检查数据库中 `module_name = 'contact'` 的记录

**解决方案**：
- 在后台启用 Contact 模块
- 保存一次 Contact 配置以创建记录

### 问题 2：联系信息显示默认值而非配置值

**可能原因**：
1. `extraData` 字段未正确解析
2. 后端 `mergeExtraData()` 未展开字段

**排查步骤**：
1. 检查后端日志是否有 `Contact module extraData parsed` 日志
2. 检查 API 响应中 `email`、`phone` 等字段是否在顶层

**解决方案**：
- 重新保存 Contact 配置
- 检查后端 `module_service.go` 中 `mergeExtraData()` 函数

### 问题 3：保存配置后前台仍显示旧数据

**可能原因**：
1. 浏览器缓存
2. 前端 JS 模块缓存

**排查步骤**：
1. 硬刷新页面（Ctrl+F5 或 Cmd+Shift+R）
2. 清除浏览器缓存

**解决方案**：
- 硬刷新页面
- 在开发模式下禁用缓存

### 问题 4：API 返回 401 未授权错误

**可能原因**：
- 错误地调用了 `/api/admin/modules` 而非 `/api/public/modules`

**排查步骤**：
1. 检查 Network 标签中的请求 URL
2. 确认前台调用的是 `/api/public/modules`

**解决方案**：
- 前台应使用公开接口 `/api/public/modules`

---

## 5. API 接口测试示例

### 5.1 获取公开模块配置

```bash
# 使用 curl 测试
curl -X GET http://localhost:8080/api/public/modules \
  -H "Accept: application/json"
```

**预期响应**：
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "moduleName": "contact",
      "enabled": true,
      "email": "info@company.com",
      "phone": "+86 123 4567 8900",
      "whatsapp": "+1 234 567 8900",
      "address": "No. 123, Industrial Park"
    }
  ]
}
```

### 5.2 保存 Contact 模块配置

```bash
# 使用 curl 测试（需要先登录获取 token）
curl -X POST http://localhost:8080/api/admin/modules \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "moduleName=contact" \
  -F "enabled=true" \
  -F "email=info@company.com" \
  -F "phone=+86 123 4567 8900" \
  -F "whatsapp=+1 234 567 8900" \
  -F "address=No. 123, Industrial Park"
```

**预期响应**：
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "moduleName": "contact",
    "enabled": true,
    ...
  }
}
```

### 5.3 查询数据库验证

```bash
# 使用 sqlite3 命令行
sqlite3 medical.db "SELECT module_name, enabled, extra_data FROM module_configs WHERE module_name='contact';"
```

**预期输出**：
```
contact|1|{"email":"info@company.com","phone":"+86 123 4567 8900","whatsapp":"+1 234 567 8900","address":"No. 123, Industrial Park"}
```

---

## 6. 代码位置参考

| 组件 | 文件路径 | 关键函数 |
|------|----------|----------|
| 后台配置界面 | `admin/admin.js` | `loadContactConfig()`, `saveContactConfig()` |
| 前台组件 | `frontend/js/components/Contact.js` | `loadContactContent()` |
| CMS 数据服务 | `frontend/js/services/cmsService.js` | `loadCMSData()`, `getCMSData()` |
| 后端控制器 | `backend/controllers/module_controller.go` | `SaveModuleConfig()`, `GetPublicModuleConfigs()` |
| 后端服务 | `backend/services/module_service.go` | `mergeExtraData()`, `GetPublicModuleConfigs()` |
| 数据模型 | `backend/models/models.go` | `ModuleConfig` 结构体 |

---

## 7. 日志说明

### 前端日志（浏览器控制台）

| 日志前缀 | 说明 |
|----------|------|
| `[Contact] Module config:` | 显示从 API 获取的 Contact 模块完整配置 |
| `[Contact] Parsed data from top-level fields:` | 显示从顶层字段解析的联系信息 |
| `[Contact] Parsed data from extraData:` | 显示从 extraData JSON 解析的联系信息 |
| `[Contact] Final contact data:` | 显示最终用于显示的联系信息 |

### 后端日志（服务器日志）

| 日志前缀 | 说明 |
|----------|------|
| `[ModuleService] Contact module extraData parsed:` | 显示 Contact 模块 extraData 解析结果 |
| `[ModuleService] Failed to parse extraData for module:` | 显示 extraData JSON 解析错误 |

---

## 8. 总结

Contact 模块配置联动的完整链路：

1. **后台保存** → `admin/admin.js` → `POST /api/admin/modules`
2. **后端处理** → `module_controller.go` → 解析 FormData → 构建 extraDataMap → 序列化 JSON
3. **数据存储** → `ModuleConfig.ExtraData` 字段
4. **公开查询** → `GET /api/public/modules` → 过滤 `enabled=true` → `mergeExtraData()` 展开字段
5. **前台加载** → `cmsService.js` → 组织数据 → `Contact.js` → 更新 DOM

任何环节出现问题都可能导致前台无法显示配置的联系信息。按照本文档的调试步骤逐一排查可快速定位问题。