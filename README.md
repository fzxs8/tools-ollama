# Ollama 模型管理工具

Ollama 模型管理工具是一个基于 Wails 框架开发的桌面应用程序，用于管理本地和远程 Ollama 服务器上的大语言模型。

## 功能特性

- **多服务器支持**：管理本地 Ollama 服务和多个远程 Ollama 服务器
- **模型管理**：查看、运行、停止和删除大语言模型
- **模型参数配置**：设置模型运行参数（温度、top_p、上下文大小等）
- **模型搜索**：根据名称、家族、标签等条件搜索模型
- **直观界面**：基于 Vue.js 和 Element Plus 构建的现代化用户界面

## 安装说明

### 环境要求

- Go 1.19+
- Node.js 16+
- Ollama 服务（本地或远程）

### 开发环境搭建

1. 克隆项目代码：

```bash
git clone <项目地址>
```

2. 安装前端依赖：

```bash
cd frontend
npm install
```

3. 运行开发环境：

```bash
wails dev
```

### 构建应用

```bash
wails build
```

docker 启动

```shell
sudo sysctl fs.inotify.max_user_watches=524288
sudo sysctl fs.inotify.max_user_instances=512
```
```shell
docker run --rm \
  -v "/home/fzxs/workspaces/demo/duola/duola-desktop":/project \
  -w /project/tools-ollama \
  -p 34115:34115 \
  -p 40000:40000 \
  -e DISPLAY=$DISPLAY \
  -e WEBKIT_DISABLE_COMPOSITING_MODE=1 \
  -e CHOKIDAR_USEPOLLING=1 \
  -v /tmp/.X11-unix:/tmp/.X11-unix:ro \
  -v /tmp/logs:/tmp/logs \
  wails-dev-env:1.0.7 \
  wails dev -debounce 2000 -reloaddirs ../duolasdk

``` 

```shell
docker run -it --rm wails-dev-env:1.0.6  ls -l
```

开放了对所有本地用户的访问权限 运行docker 访问 X server
```shell
sudo xhost +local:
```

## 项目结构

```
tools-ollama/
├── app.go              # 主应用逻辑
├── main.go             # 程序入口
├── model_manager.go    # 模型管理逻辑
├── ollama_config.go    # 配置管理
├── docs/               # 文档目录
├── frontend/           # 前端代码
├── build/              # 构建相关文件
└── wails.json          # Wails 配置文件
```

## 文档

详细文档请查看 [docs](./docs) 目录：

### 需求文档

1. [需求文档](docs/设计)
    - [通用文档](docs/设计/通用文档)
    - [ModelManager 页面](docs/设计/ModelManager)
    - [OllamaSettings 页面](docs/设计/OllamaSettings)

## 许可证

[MIT](./LICENSE)
