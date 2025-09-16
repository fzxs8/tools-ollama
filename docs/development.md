# 哆啦桌面端 - 开发指南

本指南为希望参与项目开发的贡献者提供指导。

## 1. 技术栈

-   **后端**: Go
-   **前端**: Vue 3 (Composition API) + TypeScript
-   **UI框架**: Element Plus
-   **桌面应用框架**: Wails v2
-   **数据存储**: 内嵌 `duolasdk` 存储功能 (基于SQLite)

## 2. 架构概览

项目后端采用模块化设计，核心业务逻辑被下沉到各自的管理器 (`Manager`) 中，例如 `ModelManager`, `ChatManager` 等。主入口 `app.go` 负责依赖注入和路由，将 `store`, `logger` 等核心服务注入到各个管理器中，以实现低耦合。

前端采用基于Vue 3的组件化架构，页面与模块一一对应。状态管理主要使用Vue 3的响应式API (`ref`, `reactive`)。对于跨组件或全局状态，未来计划引入Pinia。

前后端通过Wails v2框架进行通信，后端Go方法被绑定到前端，前端可以直接调用。对于实时性要求高的场景（如流式输出、日志），则使用Wails的事件系统 (`runtime.EventsEmit` 和 `runtime.EventsOn`)。

## 3. 开发规范

请严格遵守项目根目录 `docs/设计/通用文档/开发规范.md` 中定义的各项规范。以下为关键摘要：

### 代码规范
-   **Go**:
    -   所有对外暴露的结构体字段必须有 `json` 标签。
    -   所有可能出错的函数都应返回 `error`。
    -   遵循Go社区的通用编码风格。
-   **Vue/TypeScript**:
    -   使用 `<script setup lang="ts">` 语法。
    -   组件内代码顺序：`imports` -> `interfaces` -> `consts` -> `reactive data` -> `computed` -> `methods` -> `lifecycle hooks`。
    -   命名：组件用 `PascalCase`，变量和方法用 `camelCase`。
    -   强类型：尽可能使用明确的TypeScript类型，避免 `any`。

### 界面开发规范
-   **组件库**: 优先使用 `Element Plus` 提供的组件。
-   **样式**: 使用 Scoped CSS 或 CSS Modules 避免样式污染。
-   **响应式**: 布局应考虑不同窗口尺寸的适配。
-   **反馈**: 所有异步操作都必须有明确的加载状态 (`loading`)，操作结果需通过 `ElMessage` 给予用户反馈。

### 错误处理
-   **前端**: 所有对后端的API调用都必须用 `try...catch` 包裹，并在 `catch` 块中使用 `ElMessage.error` 向用户显示清晰的错误信息。
-   **后端**: 返回给前端的错误信息应清晰、具体，避免暴露内部实现细节。

## 4. 关键问题与解决方案参考

在开发过程中，我们积累了一些常见问题的解决方案，记录在各个模块的 `06.待办事项.md` 或特定的修复记录文档中，例如：

-   **Wails调试配置**: 如何配置 `wails.json` 以便使用 `dlv` 进行后端调试。
-   **状态同步问题**: 如何确保跨页面、前后端的状态一致性（例如，模型运行状态、默认服务配置）。
-   **CORS问题**: 如何在Go后端正确处理浏览器的CORS预检请求。
-   **数据迁移**: 如何在升级存储结构时，为用户提供无感的旧数据迁移。

在开始开发新功能或修复Bug前，请查阅相关模块的设计文档，特别是“待办事项”或“修复记录”部分，以避免重复踩坑。

## 5. 贡献流程

1.  Fork 本仓库。
2.  创建一个新的功能分支 (`git checkout -b feature/xxx`) 或修复分支 (`git checkout -b fix/xxx`)。
3.  进行代码开发，确保遵循开发规范。
4.  为新增的核心逻辑编写单元测试。
5.  提交代码并创建一个 Pull Request。
6.  在PR描述中清晰地说明本次变更的内容和目的。
7.  等待Code Review并根据反馈进行修改。
