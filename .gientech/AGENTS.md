# AGENTS.md - AI Agent 导航地图

> 这是 lihui 的 Harness 导航文件。代码仓库本身是记录系统，本文档负责把稳定信息源组织成可导航入口。

## 快速定位

| 我想做什么 | 去哪里 |
| --- | --- |
| 查看 Spec 规划目录 | [spec/](./spec) |
| 快速理解系统结构 | [架构总览](./gientech-harness/架构总览.md) |
| 查看设计入口 | [设计入口](./gientech-harness/docs/设计入口.md) |
| 查看执行计划入口 | [执行计划入口](./gientech-harness/docs/执行计划入口.md) |
| 理解产品与用户价值 | [产品与用户价值](./gientech-harness/docs/产品与用户价值.md) |
| 查看工程质量评分 | [质量评分](./gientech-harness/docs/质量评分.md) |
| 关注可靠性 | [可靠性说明](./gientech-harness/docs/可靠性说明.md) |
| 关注安全性 | [安全说明](./gientech-harness/docs/安全说明.md) |
| 查看生成式工程 Wiki | [wiki/](./gientech-harness/wiki) |

## 知识库结构

```
.gientech/
├── gientech-harness/
│   ├── docs/
│   │   ├── 产品规格/
│   │   │   └── index.md
│   │   ├── 参考资料/
│   │   │   ├── frontend-surface-llms.txt
│   │   │   └── index.md
│   │   ├── 执行计划/
│   │   │   ├── active/
│   │   │   │   └── index.md
│   │   │   ├── completed/
│   │   │   │   └── index.md
│   │   │   ├── index.md
│   │   │   └── tech-debt-tracker.md
│   │   ├── 生成内容/
│   │   │   └── db-schema.md
│   │   ├── 设计文档/
│   │   │   ├── core-beliefs.md
│   │   │   └── index.md
│   │   ├── 产品与用户价值.md
│   │   ├── 可靠性说明.md
│   │   ├── 安全说明.md
│   │   ├── 执行计划入口.md
│   │   ├── 设计入口.md
│   │   └── 质量评分.md
│   ├── wiki/
│   │   ├── API参考/
│   │   │   └── API参考.md
│   │   ├── 业务逻辑层/
│   │   │   └── 业务逻辑层.md
│   │   ├── 基础设施与中间件/
│   │   │   └── 基础设施与中间件.md
│   │   ├── 外部系统集成/
│   │   │   └── 外部系统集成.md
│   │   ├── 数据模型/
│   │   │   └── 数据模型.md
│   │   ├── 快速开始.md
│   │   ├── 技术栈与依赖.md
│   │   ├── 测试策略.md
│   │   ├── 编码指引.md
│   │   ├── 部署与运维.md
│   │   └── 项目概述.md
│   └── 架构总览.md
├── spec/
│   └── <taskId-任务描述>/
│       ├── actions.md
│       ├── design.md
│       └── specification.md
└── AGENTS.md
```

## 核心入口

- `spec/<taskId-任务描述>/specification.md`：Spec 模式会按任务写入需求规格。
- `spec/<taskId-任务描述>/design.md`：Spec 模式会按任务写入技术设计。
- `spec/<taskId-任务描述>/actions.md`：Spec 模式会按任务写入实施计划。
- [wiki/](./gientech-harness/wiki)：按主题查看生成式工程 Wiki 的完整目录。
- [wiki/项目概述.md](./gientech-harness/wiki/项目概述.md)：先快速建立对当前工程的整体认知。
- [wiki/快速开始.md](./gientech-harness/wiki/快速开始.md)：查看本地运行、调试与上手入口。
- [wiki/技术栈与依赖.md](./gientech-harness/wiki/技术栈与依赖.md)：补充依赖栈、基础设施与关键约束。
- [wiki/API参考/API参考.md](./gientech-harness/wiki/API参考/API参考.md)：继续深入接口与协议侧细节。
- [设计文档索引](./gientech-harness/docs/设计文档/index.md)
- [执行计划索引](./gientech-harness/docs/执行计划/index.md)
- [产品规格索引](./gientech-harness/docs/产品规格/index.md)
- [参考资料索引](./gientech-harness/docs/参考资料/index.md)

## 仓库地图

### 已确认的职责分区

- [核心源码](../backend)：包含项目主要实现代码，是理解业务能力与系统结构的第一入口。
- [前端与交互层](../frontend)：承载页面、界面组件或用户交互逻辑，适合从这里理解可见功能入口。

### 顶层目录快照

- [admin](../admin)：下一级包含 `admin.css`、`admin.js`、`index.html`
- [backend](../backend)：下一级包含 `common/`、`config/`、`controllers/`、`middleware/`、`models/`、`routes/` 等
- [frontend](../frontend)：下一级包含 `components/`、`css/`、`js/`
- [uploads](../uploads)：下一级包含 `1778806950150892000_IMG_4329_compressed.jpg`、`1778806950150892000_IMG_4329.png`、`1778806950157011000_IMG_4330_compressed.jpg`、`1778806950157011000_IMG_4330.png`、`1778806978916330000_IMG_4329_compressed.jpg`、`1778806978916330000_IMG_4329.png` 等

## Agent 协作约定

- 先把代码、配置与已存在文档视为第一事实来源，再扩展到 Harness 文档。
- `AGENTS.md` 只负责导航，不承载过长的细节正文。
- `spec/` 是 Spec 模式的真实规划目录，通常按 `spec/<taskId-任务描述>/` 组织任务级 specification.md / design.md / actions.md。
- `wiki/` 是生成型参考文档，适合快速理解模块与技术细节。
- `docs/` 下的文档更适合沉淀结构化结论、设计决策、执行计划与风险说明。
