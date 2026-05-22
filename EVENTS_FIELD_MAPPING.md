# 展会配置字段命名规范 - 完整技术文档

## 🔍 问题背景

展会配置模块存在**前后端字段命名不一致**的问题，导致前台无法正确读取后台配置的数据。

---

## 📋 数据流完整追踪

### 1️⃣ **后台管理员提交数据**

**文件**: `admin/admin.js` - `saveEventsConfig()` 函数

**提交的字段名（驼峰 camelCase）**:
```javascript
formData.append('moduleName', 'events');
formData.append('zhName', '广州展会');        // 驼峰
formData.append('enName', 'Guangzhou');       // 驼峰
formData.append('booth', '24');
formData.append('startDate', '2026/06/01');   // 驼峰
formData.append('endDate', '2026/06/03');     // 驼峰
formData.append('zhLocation', '广州琶洲');     // 驼峰
formData.append('enLocation', 'Guangzhou');   // 驼峰
formData.append('zhDescription', '展会说明');  // 顶层字段
formData.append('enDescription', 'exhibition description'); // 顶层字段
formData.append('extraData', JSON.stringify({
    icon: '🏥',
    zhLeftTitle: 'WHX Guangzhou',
    enLeftTitle: 'WHX Guangzhou',
    zhLeftSubtitle: 'World Health Expo',
    enLeftSubtitle: 'World Health Expo'
}));
```

### 2️⃣ **后端接收并存储**

**文件**: `backend/controllers/module_controller.go` - `SaveModuleConfig()` 函数 (第86-197行)

#### 处理逻辑：

```go
// 第一部分：提取通用多语言字段（存入结构体字段）
updates["zhTitle"] = ctx.PostForm("zhTitle")       // → ModuleConfig.ZhTitle
updates["enTitle"] = ctx.PostForm("enTitle")       // → ModuleConfig.EnTitle
updates["zhDescription"] = ctx.PostForm("zhDescription") // → ModuleConfig.ZhDescription
updates["enDescription"] = ctx.PostForm("enDescription")   // → ModuleConfig.EnDescription

// 第二部分：提取展会特有字段（存入 extraData JSON）
extraDataMap := make(map[string]interface{})

if zhName := ctx.PostForm("zhName"); zhName != "" {
    extraDataMap["zh_name"] = zhName              // ⚠️ 转换为下划线 snake_case
}
if enName := ctx.PostForm("enName"); enName != "" {
    extraDataMap["en_name"] = enName              // ⚠️ 转换为下划线
}
if booth := ctx.PostForm("booth"); booth != "" {
    extraDataMap["booth"] = booth                 // 保持原样
}
if startDate := ctx.PostForm("startDate"); startDate != "" {
    extraDataMap["start_date"] = startDate         // ⚠️ 转换为下划线
}
if endDate := ctx.PostForm("endDate"); endDate != "" {
    extraDataMap["end_date"] = endDate             // ⚠️ 转换为下划线
}
if zhLocation := ctx.PostForm("zhLocation"); zhLocation != "" {
    extraDataMap["zh_location"] = zhLocation       // ⚠️ 转换为下划线
}
if enLocation := ctx.PostForm("enLocation"); enLocation != "" {
    extraDataMap["en_location"] = enLocation       // ⚠️ 转换为下划线
}

// 序列化并存入数据库
updates["extraData"] = json.Marshal(extraDataMap)
```

#### 最终数据库存储格式：

**表**: `module_configs`
**记录示例** (moduleName = 'events'):

| 字段名 | 值 | 类型 |
|-------|-----|------|
| id | 4 | int |
| module_name | "events" | string |
| enabled | true | boolean |
| zh_title | "" | string *(空，因为展会不用这个)* |
| en_title | "" | string |
| zh_description | "展会说明" | string *(来自顶层字段)* |
| en_description | "exhibition description" | string *(来自顶层字段)* |
| **extra_data** | **见下方JSON** | **text (JSON字符串)** |

**extra_data 字段存储的完整 JSON**:
```json
{
    "zh_name": "广州展会",
    "en_name": "Guangzhou",
    "booth": "24",
    "start_date": "2026/06/01",
    "end_date": "2026/06/03",
    "zh_location": "广州琶洲国际会展中心",
    "en_location": "Guangzhou Pazhou Convention Center",
    "icon": "🏥",
    "zhLeftTitle": "WHX Guangzhou 2026",
    "enLeftTitle": "WHX Guangzhou 2026",
    "zhLeftSubtitle": "World Health Expo",
    "enLeftSubtitle": "World Health Expo"
}
```

### 3️⃣ **API 接口返回数据**

#### ❌ **前台接口** (`GET /api/public/modules`)

**调用方法**: `GetPublicModuleConfigs()` (module_service.go:199-203)

**返回格式**: 原始 `[]models.ModuleConfig` 结构体（**未处理 extraData**）

```javascript
// 前台 cmsService.js 收到的数据结构
{
    moduleName: "events",
    enabled: true,
    zhTitle: "",                    // 空
    enTitle: "",                    // 空
    zhSubtitle: "",                 // 空
    enSubtitle: "",                 // 空
    zhContent: "",                  // 空
    enContent: "",                  // 空
    title: "",
    subtitle: "",
    content: "",
    imagePath: "",
    sortOrder: 0,
    zhDescription: "展会说明",      // ✅ 有值（来自顶层字段）
    enDescription: "exhibition description", // ✅ 有值（来自顶层字段）
    description: "",
    
    // ⚠️ extraData 是一个 JSON 字符串，未被解析！
    extraData: "{\"booth\":\"24\",\"zh_name\":\"广州展会\",\"en_name\":\"Guangzhou\",\"start_date\":\"2026/06/01\",\"end_date\":\"2026/06/03\",\"zh_location\":\"广州琶洲国际会展中心\",\"en_location\":\"Guangzhou Pazhou Convention Center\",\"icon\":\"🏥\",\"zhLeftTitle\":\"WHX Guangzhou 2026\",\"enLeftTitle\":\"WHX Guangzhou 2026\",\"zhLeftSubtitle\":\"World Health Expo\",\"enLeftSubtitle\":\"World Health Expo\"}"
}
```

#### ✅ **后台接口** (`GET /api/admin/modules`)

**调用方法**: `GetModuleConfigs()` (module_service.go:19-41)

**返回格式**: 经过 `mergeExtraData()` 处理后的 `map[string]interface{}`

```javascript
// 后台 admin.js 收到的数据结构
{
    id: 4,
    moduleName: "events",
    enabled: true,
    zhTitle: "",
    enTitle: "",
    zhDescription: "展会说明",
    enDescription: "exhibition description",
    
    // ✅ extraData 已被解析并合并到顶层！
    zh_name: "广州展会",            // 来自 extraData（下划线）
    en_name: "Guangzhou",           // 来自 extraData（下划线）
    booth: "24",
    start_date: "2026/06/01",      // 来自 extraData（下划线）
    end_date: "2026/06/03",        // 来自 extraData（下划线）
    zh_location: "广州琶洲国际会展中心",
    en_location: "Guangzhou Pazhou Convention Center",
    icon: "🏥",
    zhLeftTitle: "WHX Guangzhou 2026",
    enLeftTitle: "WHX Guangzhou 2026",
    zhLeftSubtitle: "World Health Expo",
    enLeftSubtitle: "World Health Expo",
    
    // 同时保留原始字符串
    extraData: "{json字符串...}"
}
```

---

## 🎯 核心问题总结

### 问题 #1：**前后端接口不一致**

| 接口路径 | 使用者 | 是否解析 extraData | 字段位置 |
|---------|-------|-------------------|---------|
| `/api/admin/modules` | 后台 admin.js | ✅ 是 (mergeExtraData) | 顶层 + extraData字符串 |
| `/api/public/modules` | 前台 cmsService.js | ❌ 否 | 只有 extraData字符串 |

### 问题 #2：**字段命名转换混乱**

| 阶段 | 命名规范 | 示例 |
|-----|---------|------|
| **HTML 表单** | camelCase (驼峰) | `zhName`, `startDate` |
| **extraData 存储** | snake_case (下划线) | `zh_name`, `start_date` |
| **Go 结构体** | PascalCase (首字母大写) | `ZhName`, `StartDate` *(不存在)* |
| **JSON tag** | camelCase (驼峰) | `zhTitle`, `moduleName` |

### 问题 #3：**前台需要手动解析 extraData**

由于 `/api/public/modules` 不自动解析 extraData，前台必须自己解析 JSON 字符串才能获取展会数据。

---

## ✅ 正确的字段映射关系

### 展会特有字段（全部在 extraData 中）

| 显示内容 | HTML 表单名 | extraData Key | 前台读取方式 |
|---------|-----------|--------------|-------------|
| **展会名称(中文)** | `zhName` | `zh_name` | `extra.zh_name` |
| **展会名称(英文)** | `enName` | `en_name` | `extra.en_name` |
| **展位号** | `booth` | `booth` | `extra.booth` |
| **开始日期** | `startDate` | `start_date` | `extra.start_date` |
| **结束日期** | `endDate` | `end_date` | `extra.end_date` |
| **地点(中文)** | `zhLocation` | `zh_location` | `extra.zh_location` |
| **地点(英文)** | `enLocation` | `en_location` | `extra.en_location` |
| **图标** | `eventsIcon` → extraData.icon | `icon` | `extra.icon` |
| **左侧标题(中文)** | `eventsZhLeftTitle` → extraData.zhLeftTitle | `zhLeftTitle` | `extra.zhLeftTitle` |
| **左侧标题(英文)** | `eventsEnLeftTitle` → extraData.enLeftTitle | `enLeftTitle` | `extra.enLeftTitle` |
| **左侧副标题(中文)** | `eventsZhLeftSubtitle` → extraData.zhLeftSubtitle | `zhLeftSubtitle` | `extra.zhLeftSubtitle` |
| **左侧副标题(英文)** | `eventsEnLeftSubtitle` → extraData.enLeftSubtitle | `enLeftSubtitle` | `extra.enLeftSubtitle` |

### 通用多语言字段（直接在结构体中）

| 显示内容 | HTML 表单名 | 数据库字段 | 前台读取方式 |
|---------|-----------|-----------|-------------|
| **描述(中文)** | `eventsZhDescription` | `zh_description` | `module.zhDescription` |
| **描述(英文)** | `eventsEnDescription` | `en_description` | `module.enDescription` |

---

## 🔧 解决方案

### 方案 A：修改后端（推荐）✨

让 `/api/public/modules` 也调用 `mergeExtraData()`，保持与后台接口一致。

**优点**:
- 前端代码简单，不需要手动解析
- 前后台行为一致
- 易于维护

**缺点**:
- 需要修改后端代码
- 需要重新编译部署

### 方案 B：修复前端（当前采用）✅

在前端正确解析 extraData 并使用正确的字段名。

**优点**:
- 不需要改后端
- 立即生效
- 前端自主可控

**缺点**:
- 需要维护两套字段映射逻辑
- 如果后端变化，前端也要跟着变

---

## 📝 代码实现参考

### Events.js 正确的实现方式

```javascript
import { getCMSData, getCurrentLangFromStorage } from '../services/cmsService.js';

function parseExtraData(module) {
    if (!module || !module.extraData) return {};
    
    try {
        if (typeof module.extraData === 'string') {
            return JSON.parse(module.extraData);
        }
        return module.extraData;
    } catch (e) {
        console.warn('Failed to parse extraData:', e);
        return {};
    }
}

export function loadEventContent() {
    const { modules } = getCMSData();
    const eventsModule = modules['events'] || {};
    
    if (eventsModule.enabled === false) return;

    // ✅ 关键：必须先解析 extraData
    const extra = parseExtraData(eventsModule);
    
    console.log('🔍 Debug - Events Module:', eventsModule);
    console.log('🔍 Debug - Parsed ExtraData:', extra);

    // ✅ 使用正确的字段名（snake_case，来自 extraData）
    const eventName = getCurrentLangFromStorage() === 'zh' 
        ? (extra.zh_name || '') 
        : (extra.en_name || '');
        
    const eventDescription = getCurrentLangFromStorage() === 'zh' 
        ? (eventsModule.zhDescription || '')  // 注意：这个在顶层
        : (eventsModule.enDescription || ''); // 注意：这个在顶层
        
    const booth = extra.booth || '';
    const startDate = extra.start_date || '';
    const endDate = extra.end_date || '';
    const eventLocation = getCurrentLangFromStorage() === 'zh' 
        ? (extra.zh_location || '')
        : (extra.en_location || '');
        
    const icon = extra.icon || '';
    const leftTitle = getCurrentLangFromStorage() === 'zh' 
        ? (extra.zhLeftTitle || '')
        : (extra.enLeftTitle || '');
    const leftSubtitle = getCurrentLangFromStorage() === 'zh' 
        ? (extra.zhLeftSubtitle || '')
        : (extra.enLeftSubtitle || '');

    // ... 更新 DOM 元素
}
```

---

## 🧪 调试指南

### 在浏览器控制台执行：

```javascript
// 1. 查看 raw 数据
const { modules } = getCMSData();
console.log('Raw events module:', modules['events']);

// 2. 查看 extraData 内容
const rawExtra = modules['events'].extraData;
console.log('ExtraData string:', rawExtra);
console.log('Parsed extraData:', JSON.parse(rawExtra));

// 3. 检查具体字段
const parsed = JSON.parse(rawExtra);
console.log('zh_name:', parsed.zh_name);          // 应该有值
console.log('en_name:', parsed.en_name);          // 应该有值
console.log('booth:', parsed.booth);              // 应该有值
console.log('start_date:', parsed.start_date);    // 应该有值
console.log('end_date:', parsed.end_date);        // 应该有值
```

**预期输出**:
```
zh_name: "广州展会"              ✅
en_name: "Guangzhou"             ✅
booth: "24"                      ✅
start_date: "2026/06/01"         ✅
end_date: "2026/06/03"           ✅
```

如果输出都是 `undefined`，说明：
- extraData 格式不对
- 或者字段名不是预期的

---

## 🐛 常见错误排查

### 错误1：字段名用错了

```javascript
// ❌ 错误：使用了驼峰（这是表单名，不是存储名）
eventsModule.zhName        // undefined!
eventsModule.startDate     // undefined!

// ✅ 正确：使用下划线（这是 extraData 中的 key）
extra.zh_name              // "广州展会" ✅
extra.start_date           // "2026/06/01" ✅
```

### 错误2：忘记解析 extraData

```javascript
// ❌ 错误：直接访问模块对象的属性
modules['events'].zh_name  // undefined! (因为没解析)

// ✅ 正确：先解析再访问
const extra = JSON.parse(modules['events'].extraData);
extra.zh_name              // "广州展会" ✅
```

### 错误3：混淆顶层字段和 extraData 字段

```javascript
// zhDescription 和 enDescription 在顶层（不在 extraData 中）
const eventsModule = modules['events'];

// ✅ 正确
eventsModule.zhDescription    // "展会说明" ✅ (顶层字段)
eventsModule.enDescription    // "exhibition description" ✅ (顶层字段)

// ❌ 错误
extra.zhDescription           // undefined! (不在 extraData 中)
```

---

## 📊 性能影响评估

| 操作 | 影响 | 说明 |
|-----|------|------|
| JSON.parse(extraData) | < 0.1ms | 极快，可忽略 |
| 额外内存占用 | < 1KB | 解析后的对象很小 |
| 代码复杂度 | 低 | 只需几行解析代码 |

---

## 🎯 最佳实践建议

### 1️⃣ **统一命名规范**（长期方案）

建议整个项目统一使用一种命名规范：
- **推荐**: snake_case (下划线) - 与数据库、JSON API 一致
- 或: camelCase (驼峰) - 与 JavaScript 惯例一致

**避免混用**导致的混乱。

### 2️⃣ **API 文档化**

为每个接口明确标注：
- 请求参数格式
- 响应数据结构
- 字段命名规范
- 示例数据

### 3️⃣ **类型定义**

使用 TypeScript 或 JSDoc 定义清晰的数据类型：
```typescript
interface EventsModule {
    moduleName: string;
    enabled: boolean;
    zhDescription: string;      // 顶层字段
    enDescription: string;      // 顶层字段
    extraData: string;           // JSON 字符串
    
    // 以下字段在 extraData 解析后可用
    // zh_name: string;
    // en_name: string;
    // booth: string;
    // start_date: string;
    // end_date: string;
    // ...
}

interface ParsedEventsExtraData {
    zh_name: string;
    en_name: string;
    booth: string;
    start_date: string;
    end_date: string;
    zh_location: string;
    en_location: string;
    icon: string;
    zhLeftTitle: string;
    enLeftTitle: string;
    zhLeftSubtitle: string;
    enLeftSubtitle: string;
}
```

---

## 📝 更新历史

| 版本 | 日期 | 作者 | 变更内容 |
|-----|------|------|---------|
| v1.0 | 2026-05-21 | AI Assistant | 初始版本，完整字段映射文档 |

---

**最后更新**: 2026-05-21  
**适用范围**: 医疗设备 CMS 系统 - 展会配置模块  
**相关文件**:
- `frontend/js/components/Events.js`
- `admin/admin.js`
- `backend/controllers/module_controller.go`
- `backend/services/module_service.go`
- `backend/models/models.go`
