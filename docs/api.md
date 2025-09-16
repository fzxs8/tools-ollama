# 哆啦桌面端 - API 文档

本文档描述了哆啦桌面端所提供的各类API接口，主要分为两部分：供外部应用调用的**公共API**和供内部前端页面使用的**管理API**。

## 1. 公共API (Public-Facing API)

这是通过 **OpenAI 适配器** 功能暴露的、与OpenAI `v1/chat/completions` 接口兼容的API。

### Chat Completions

-   **路径**: `/v1/chat/completions`
-   **方法**: `POST`
-   **描述**: 接收符合OpenAI规范的聊天请求，并将其转发给配置好的Ollama模型进行处理。支持流式（Server-Sent Events）和非流式响应。

#### cURL 示例：非流式请求

```sh
curl http://127.0.0.1:11223/v1/chat/completions -X POST \
-H "Content-Type: application/json" \
-d \
'{ 
  "model": "llama3",
  "messages": [
    {
      "role": "user",
      "content": "你好，介绍一下你自己"
    }
  ],
  "stream": false
}'
```

#### cURL 示例：流式请求

```sh
curl http://127.0.0.1:11223/v1/chat/completions -X POST \
-H "Content-Type: application/json" \
-d \
'{ 
  "model": "llama3",
  "messages": [
    {
      "role": "user",
      "content": "给我讲一个关于程序员的笑话"
    }
  ],
  "stream": true
}'
```

## 2. 内部管理API (Wails-Bound)

这些是后端Go方法，通过Wails框架绑定到前端，用于应用内部的管理和操作。前端可通过 `window.go.main.App` 对象调用。

### 会话管理 (ChatManager)

-   `ChatMessage(modelName string, messages []types.Message, stream bool) (string, error)`: 发送聊天消息。
-   `ListConversations() ([]*types.Conversation, error)`: 获取所有对话列表。
-   `SaveConversation(conv *types.Conversation) (*types.Conversation, error)`: 保存一个对话（新建或更新）。
-   `GetConversation(id string) (*types.Conversation, error)`: 获取单个对话的详细信息。
-   `DeleteConversation(id string) error`: 删除一个对话。

### 模型管理 (ModelManager)

-   `ListModelsByServer(serverID string) ([]types.Model, error)`: 获取指定服务器上的模型列表。
-   `RunModel(modelName string, params map[string]interface{}) error`: 运行一个模型。
-   `StopModel(modelName string) error`: 停止一个正在运行的模型。
-   `DeleteModel(modelName string) error`: 删除一个本地模型。
-   `TestModel(modelName, prompt string) (*types.TestModelResponse, error)`: 使用给定的提示词测试一个模型。
-   `DownloadModel(serverID, modelName string) (types.DownloadTask, error)`: 下载一个新模型。
-   `GetDownloadTasks() ([]types.DownloadTask, error)`: 获取所有下载任务的状态。
-   `CancelDownload(taskID string) error`: 取消一个正在进行的下载。
-   `SearchOnlineModels(query string) ([]interface{}, error)`: 从ollamadb.dev搜索在线模型。

### 提示词工程 (PromptEngineering)

-   `GeneratePromptStream(req GenerateRequest) error`: 触发一个并行的、流式的提示词生成任务。结果通过Wails事件返回。
-   `SavePrompt(prompt Prompt) (*Prompt, error)`: 保存一个新的提示词。
-   `ListPrompts() ([]Prompt, error)`: 获取所有已保存的提示词。
-   `UpdatePrompt(prompt Prompt) (*Prompt, error)`: 更新一个已存在的提示词。
-   `DeletePrompt(id string) error`: 删除一个提示词。

### 服务设置 (OllamaSettings)

-   `GetServers() ([]types.OllamaServerConfig, error)`: 获取所有已配置的Ollama服务列表。
-   `AddServer(server types.OllamaServerConfig) error`: 添加一个新的服务配置。
-   `UpdateServer(server types.OllamaServerConfig) error`: 更新一个已存在的服务配置。
-   `DeleteServer(serverID string) error`: 删除一个服务配置。
-   `SetActiveServer(serverID string) error`: 将指定服务设为活动（默认）服务。
-   `TestOllamaServer(baseURL string) (string, error)`: 测试与指定Ollama服务的连接。

### OpenAI 适配器 (OpenAIAdapter)

-   `GetOpenAIAdapterConfig() (*OpenAIAdapterConfig, error)`: 获取适配器的当前配置。
-   `SaveOpenAIAdapterConfig(config OpenAIAdapterConfig) error`: 保存适配器的新配置。
-   `StartAdapterServer() error`: 启动适配器服务。
-   `StopAdapterServer() error`: 停止适配器服务。
-   `GetOpenAIAdapterStatus() (*OpenAIAdapterStatus, error)`: 获取适配器服务的当前运行状态。
-   `GetAdapterAPIDocs() (map[string]string, error)`: 获取格式化后的API `curl` 示例。

### API 调试器 (ApiDebugger)

-   `SendHttpRequest(request types.ApiRequest) (types.ApiResponse, error)`: 发送一个完整的、自定义的HTTP请求到指定的Ollama服务。
