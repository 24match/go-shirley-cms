# 前台后台配置联动检查 - 需求规格说明书

## 1. 引言

### 功能概述

本需求旨在检查系统中前台（frontend）与后台（admin）配置数据的联动机制，识别是否存在后台配置了数据但前台未能成功展示的问题。重点检查 Contact（联系我们）模块及其他核心模块（Banner、About、Products、Factory、Advantage、Events）的配置同步链路。

### 背景与目标

系统采用前后端分离架构：
- **后端**：Go + Gin 提供 RESTful API，数据存储在 SQLite 数据库
- **前台官网**：原生 HTML/CSS/JS，通过 `/api/public/*` 接口获取数据
- **后台管理**：独立管理界面，通过 `/api/admin/*` 接口管理数据

当前发现 Contact 模块存在配置字段（email、phone、whatsapp、address）在后台配置后，前台可能无法正确显示的问题。

### 关键利益相关者与用户

| 角色 | 职责 |
|------|------|
| 系统管理员 | 通过后台管理界面配置各模块数据 |
| 前台访客 | 访问官网查看配置内容，提交联系表单 |
| 开发维护人员 | 排查配置联动问题，修复数据同步链路 |

### 预期价值与收益

1. 确保后台配置的数据能够准确同步到前台展示
2. 建立完整的配置联动检查机制
3. 提升用户体验，避免配置失效问题

---

## 2. 需求

### 需求 1：后台配置数据管理

**用户故事**：作为一名系统管理员，我希望通过后台管理界面配置各模块的数据（如联系方式、Banner 内容等），以便控制前台官网展示的内容。

**验收标准（EARS 格式）**：

1. 当管理员在后台 Contact 模块配置页面填写 email、phone、whatsapp、address 字段并保存时，系统应将这些数据保存到 ModuleConfig 表的 extraData 字段中
2. 当管理员保存配置时，系统应返回成功响应，并在数据库中持久化配置数据
3. 当管理员再次打开配置页面时，系统应回显之前保存的配置数据

---

### 需求 2：后台配置 API 接口

**用户故事**：作为一名后台管理员，我希望通过 API 接口保存和获取模块配置，以便前端界面能够与后端进行数据交互。

**验收标准（EARS 格式）**：

1. 当收到 POST `/api/admin/modules` 请求时，系统应将配置数据保存到 ModuleConfig 表
2. 当收到 GET `/api/admin/modules` 请求时，系统应返回所有模块配置列表（需认证）
3. 当收到 GET `/api/public/modules` 请求时，系统应返回已启用模块的公开配置（无需认证）
4. 当模块配置包含 extraData 时，系统应将 extraData 中的字段合并到返回结果中，以便前端直接访问

---

### 需求 3：前台数据获取

**用户故事**：作为一名前台访客，我希望访问官网时能够看到管理员配置的最新内容，以便获取准确的联系信息和业务介绍。

**验收标准（EARS 格式）**：

1. 当前台页面加载时，系统应自动调用 `/api/public/modules` 接口获取模块配置
2. 当获取到 Contact 模块配置时，系统应解析 extraData 中的 email、phone、whatsapp、address 字段
3. 当配置数据成功加载后，系统应将数据显示在对应的 HTML 元素中（如 `#contactEmail`、`#contactPhone` 等）

---

### 需求 4：配置联动检查

**用户故事**：作为一名开发维护人员，我希望系统能够检查后台配置与前台展示之间的数据链路，以便快速定位配置未生效的问题。

**验收标准（EARS 格式）**：

1. 当后台保存配置后，系统应确保数据正确写入 ModuleConfig 表
2. 当前台请求公开接口时，系统应正确返回已保存的配置数据
3. 当 extraData 中包含联系信息字段时，系统应在返回结果中将这些字段展开到顶层，便于前端访问

---

### 需求 5：错误处理

**用户故事**：作为一名系统用户，我希望在配置保存或加载失败时收到明确的错误提示，以便了解问题原因并采取相应措施。

**验收标准（EARS 格式）**：

1. 当配置保存失败时，系统应返回具体的错误信息
2. 当前台数据加载失败时，系统应显示默认值或友好的错误提示
3. 当 extraData JSON 解析失败时，系统应忽略解析错误并继续处理其他字段

---

## 3. 配置联动链路分析

### 数据流链路

```
后台管理界面 → POST /api/admin/modules → ModuleController → ModuleService → Database
                                                              ↓
前台官网 ← GET /api/public/modules ← ModuleController ← ModuleService ← Database
```

### 当前实现分析

#### 后台保存流程（admin/admin.js）

```javascript
async function saveContactConfig() {
  const formData = new FormData();
  formData.append('moduleName', 'contact');
  formData.append('email', document.getElementById('contactEmail').value);
  formData.append('phone', document.getElementById('contactPhone').value);
  formData.append('whatsapp', document.getElementById('contactWhatsApp').value);
  formData.append('address', document.getElementById('contactAddress').value);
  const success = await saveModuleConfig('contact', formData);
}
```

#### 后端处理流程（module_controller.go）

```go
// 处理 contact 模块的特殊字段
if moduleName == "contact" {
    if email := ctx.PostForm("email"); email != "" {
        extraDataMap["email"] = email
    }
    if phone := ctx.PostForm("phone"); phone != "" {
        extraDataMap["phone"] = phone
    }
    if whatsapp := ctx.PostForm("whatsapp"); whatsapp != "" {
        extraDataMap["whatsapp"] = whatsapp
    }
    if address := ctx.PostForm("address"); address != "" {
        extraDataMap["address"] = address
    }
}
```

#### 后端返回流程（module_service.go）

```go
func mergeExtraData(mc *models.ModuleConfig) map[string]interface{} {
    // 解析并合并 ExtraData 中的字段
    if mc.ExtraData != "" {
        var extra map[string]interface{}
        if err := json.Unmarshal([]byte(mc.ExtraData), &extra); err == nil {
            for k, v := range extra {
                result[k] = v  // 将 extraData 中的字段展开到顶层
            }
        }
    }
    return result
}
```

#### 前台加载流程（frontend/js/components/Contact.js）

```javascript
export function loadContactContent() {
    const { modules } = getCMSData();
    const contactModule = modules['contact'];
    
    // 方式 1：后端 mergeExtraData 已经解析出来的字段
    if (contactModule.email) contactData.email = contactModule.email;
    if (contactModule.phone) contactData.phone = contactModule.phone;
    if (contactModule.whatsapp) contactData.whatsapp = contactModule.whatsapp;
    if (contactModule.address) contactData.address = contactModule.address;
    
    // 方式 2：从 extraData JSON 字符串中解析
    if (contactModule.extraData) {
        try {
            const extraData = JSON.parse(contactModule.extraData);
            if (extraData.email && !contactData.email) contactData.email = extraData.email;
            // ...
        } catch (e) { /* 忽略解析错误 */ }
    }
}
```

---

## 4. 潜在问题识别

### 问题 1：enabled 字段处理

**现象**：后台保存 contact 配置时，`enabled` 字段可能未正确传递

**分析**：
- 后台 `saveContactConfig()` 函数中未设置 `enabled` 字段
- 后端 `SaveModuleConfig` 中 `enabled := getStringValue(updates, "enabled", "") == "true"`，如果未传递则为 false

**影响**：如果 enabled 为 false，前台 `/api/public/modules` 接口不会返回该模块（因为只返回 enabled=true 的模块）

---

### 问题 2：extraData 字段重复处理

**现象**：extraData 可能被双重序列化

**分析**：
- 后台保存时，contact 的特殊字段（email、phone 等）被写入 extraDataMap，然后序列化为 JSON 存入 extraData 字段
- 但同时，`formData.append('extraData', JSON.stringify({...}))` 也可能被调用

**影响**：可能导致数据覆盖或解析失败

---

### 问题 3：前台 enabled 检查

**现象**：前台 `loadContactContent()` 检查 `contactModule.enabled !== false`

**分析**：
- 如果后台未传递 enabled 字段，默认为 false
- 前台将跳过内容加载

---

## 5. 质量检查

| 检查项 | 状态 |
|--------|------|
| 用户故事清晰完整 | ✅ |
| 验收标准可测量 | ✅ |
| 覆盖正常路径和异常路径 | ✅ |
| 与用户概念保持一致 | ✅ |
| 功能完整 | ✅ |

---

## 6. 附录

### 相关接口列表

| 接口 | 方法 | 认证 | 说明 |
|------|------|------|------|
| `/api/admin/modules` | POST | 是 | 保存模块配置 |
| `/api/admin/modules` | GET | 是 | 获取所有模块配置 |
| `/api/public/modules` | GET | 否 | 获取公开模块配置 |
| `/api/admin/modules/:name` | GET | 是 | 获取单个模块配置 |
| `/api/admin/modules/:name` | DELETE | 是 | 删除模块配置 |
| `/api/admin/modules/:name/image` | DELETE | 是 | 删除模块图片 |

### 数据模型

| 模型 | 表名 | 说明 |
|------|------|------|
| ModuleConfig | module_configs | 存储模块配置，包含 extraData 字段 |

### 关键文件

| 文件 | 职责 |
|------|------|
| [`admin/admin.js`](admin/admin.js) | 后台配置界面逻辑 |
| [`frontend/js/components/Contact.js`](frontend/js/components/Contact.js) | 前台联系模块内容加载 |
| [`backend/controllers/module_controller.go`](backend/controllers/module_controller.go) | 模块配置 API 处理 |
| [`backend/services/module_service.go`](backend/services/module_service.go) | 模块配置业务逻辑 |
| [`backend/models/models.go`](backend/models/models.go) | 数据模型定义 |