# Ollama Desktop Client

Ollama桌面客户端，提供图形化界面来管理本地Ollama模型。

## 功能特性

- 模型管理：查看、下载、删除本地模型
- 模型市场：浏览和下载Ollama模型库中的模型
- 聊天测试：与模型进行交互式对话
- 系统监控：监控Ollama服务状态和系统资源使用情况
- 参数配置：自定义模型运行参数
- 服务设置：配置本地和远程Ollama服务

## 界面预览

![Ollama Client](./screenshot.png)

## 安装要求

- Ollama服务已安装并运行
- 支持的操作系统：Windows、macOS、Linux

## 快速开始

1. 确保Ollama服务正在运行：
   ```bash
   ollama serve
   ```

2. 安装前端依赖（可能需要管理员权限）：
   ```bash
   cd frontend
   npm install --legacy-peer-deps
   ```
   
   如果遇到权限问题，请尝试：
   ```bash
   sudo chown -R $(whoami) .
   npm install --legacy-peer-deps
   ```

3. 构建并运行客户端：
   ```bash
   wails build
   ./build/bin/ollama-desktop
   ```

## 使用说明

### 模型管理
在"模型管理"页面，您可以：
- 查看本地已下载的模型列表
- 查看模型详细信息
- 运行、加载、卸载和删除模型

### 模型市场
在"模型市场"页面，您可以：
- 浏览Ollama官方模型库
- 搜索特定模型
- 查看模型详情并下载

### 聊天测试
在"聊天测试"页面，您可以：
- 选择要测试的模型
- 与模型进行对话
- 调整模型参数（温度、Top-P等）

### 系统监控
在"系统监控"页面，您可以：
- 查看Ollama服务状态
- 监控系统资源使用情况
- 查看模型加载状态
- 查看最近的活动日志

### 服务设置
在"服务设置"页面，您可以：
- 配置本地Ollama服务参数
- 添加和管理多个远程Ollama服务
- 测试服务连接状态
- 启动或停止本地服务

## 技术栈

- 前端：Vue 3 + TypeScript + Element Plus
- 后端：Go + Wails框架
- 通信：通过Ollama REST API与服务交互

## 项目结构

```
frontend/
  ├── src/
  │   ├── views/           # 页面组件
  │   │   ├── ModelManager.vue    # 模型管理页面
  │   │   ├── ModelMarket.vue     # 模型市场页面
  │   │   ├── ChatInterface.vue   # 聊天测试页面
  │   │   ├── SystemMonitor.vue   # 系统监控页面
  │   │   └── OllamaSettings.vue  # 服务设置页面
  │   ├── router/          # 路由配置
  │   ├── App.vue          # 主应用组件
  │   └── main.ts          # 应用入口
  ├── vite.config.ts       # Vite配置
  └── package.json         # 项目依赖
```

## Vue项目结构

本项目使用Vue 3 Composition API和TypeScript构建，包含以下主要文件：
- `src/App.vue`：主Vue组件
- `src/main.ts`：应用入口文件
- `src/router/index.ts`：路由配置
- `src/views/*.vue`：各个功能页面

## 依赖说明

本项目使用以下主要依赖：
- Vue 3：前端框架
- Element Plus：UI组件库
- Vue Router：路由管理
- TypeScript：类型检查
- Vite：构建工具

## 故障排除

### 构建错误
如果遇到构建错误，请尝试以下步骤：

1. 清理node_modules并重新安装：
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install --legacy-peer-deps
   ```

2. 检查Vite配置是否正确：
   确保[vite.config.ts](file:///home/fzxs/workspaces/demo/duola/duola-desktop/duola-tools/frontend/vite.config.ts)文件中正确配置了Element Plus和Vue插件

3. 检查Element Plus导入：
   确保在[main.ts](file:///home/fzxs/workspaces/demo/duola/duola-desktop/duola-tools/frontend/src/main.ts)中正确导入了Element Plus及其图标

### 权限问题
如果遇到权限问题，请尝试：
```bash
sudo chown -R $(whoami) .
```

### 依赖冲突
如果遇到依赖版本冲突，请使用`--legacy-peer-deps`参数：
```bash
npm install --legacy-peer-deps
```

## 开发

```bash
# 开发模式运行
wails dev

# 构建生产版本
wails build
```

## 许可证

MIT