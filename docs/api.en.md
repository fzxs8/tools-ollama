# Duola Desktop - API Documentation

This document describes the various API endpoints provided by Duola Desktop, divided into two main parts: the **Public-Facing API** for external applications and the **Internal Management API** for the internal frontend pages.

## 1. Public-Facing API

This is the OpenAI `v1/chat/completions` compatible API exposed through the **OpenAI Adapter** feature.

### Chat Completions

-   **Path**: `/v1/chat/completions`
-   **Method**: `POST`
-   **Description**: Receives and processes chat requests compatible with the OpenAI specification, forwarding them to the configured Ollama model. It supports both streaming (Server-Sent Events) and non-streaming responses.

#### cURL Example: Non-Streaming Request

```sh
curl http://127.0.0.1:11223/v1/chat/completions -X POST \
-H "Content-Type: application/json" \
-d '{
  "model": "llama3",
  "messages": [
    {
      "role": "user",
      "content": "Hello, introduce yourself."
    }
  ],
  "stream": false
}'
```

#### cURL Example: Streaming Request

```sh
curl http://127.0.0.1:11223/v1/chat/completions -X POST \
-H "Content-Type: application/json" \
-d '{
  "model": "llama3",
  "messages": [
    {
      "role": "user",
      "content": "Tell me a joke about a programmer."
    }
  ],
  "stream": true
}'
```

## 2. Internal Management API (Wails-Bound)

These are Go methods bound to the frontend via the Wails framework for internal management and operations. The frontend can call them via the `window.go.main.App` object.

### Chat Manager

-   `ChatMessage(modelName string, messages []types.Message, stream bool) (string, error)`: Sends a chat message.
-   `ListConversations() ([]*types.Conversation, error)`: Gets the list of all conversations.
-   `SaveConversation(conv *types.Conversation) (*types.Conversation, error)`: Saves a conversation (creates or updates).
-   `GetConversation(id string) (*types.Conversation, error)`: Gets the details of a single conversation.
-   `DeleteConversation(id string) error`: Deletes a conversation.

### Model Manager

-   `ListModelsByServer(serverID string) ([]types.Model, error)`: Gets the list of models on a specific server.
-   `RunModel(modelName string, params map[string]interface{}) error`: Runs a model.
-   `StopModel(modelName string) error`: Stops a running model.
-   `DeleteModel(modelName string) error`: Deletes a local model.
-   `TestModel(modelName, prompt string) (*types.TestModelResponse, error)`: Tests a model with a given prompt.
-   `DownloadModel(serverID, modelName string) (types.DownloadTask, error)`: Downloads a new model.
-   `GetDownloadTasks() ([]types.DownloadTask, error)`: Gets the status of all download tasks.
-   `CancelDownload(taskID string) error`: Cancels an ongoing download.
-   `SearchOnlineModels(query string) ([]interface{}, error)`: Searches for online models from ollamadb.dev.

### Prompt Engineering

-   `GeneratePromptStream(req GenerateRequest) error`: Triggers a parallel, streaming prompt generation task. Results are returned via Wails events.
-   `SavePrompt(prompt Prompt) (*Prompt, error)`: Saves a new prompt.
-   `ListPrompts() ([]Prompt, error)`: Gets all saved prompts.
-   `UpdatePrompt(prompt Prompt) (*Prompt, error)`: Updates an existing prompt.
-   `DeletePrompt(id string) error`: Deletes a prompt.

### Ollama Settings

-   `GetServers() ([]types.OllamaServerConfig, error)`: Gets the list of all configured Ollama services.
-   `AddServer(server types.OllamaServerConfig) error`: Adds a new service configuration.
-   `UpdateServer(server types.OllamaServerConfig) error`: Updates an existing service configuration.
-   `DeleteServer(serverID string) error`: Deletes a service configuration.
-   `SetActiveServer(serverID string) error`: Sets the specified service as the active (default) one.
-   `TestOllamaServer(baseURL string) (string, error)`: Tests the connection to a specified Ollama service.

### OpenAI Adapter

-   `GetOpenAIAdapterConfig() (*OpenAIAdapterConfig, error)`: Gets the current adapter configuration.
-   `SaveOpenAIAdapterConfig(config OpenAIAdapterConfig) error`: Saves the new adapter configuration.
-   `StartAdapterServer() error`: Starts the adapter service.
-   `StopAdapterServer() error`: Stops the adapter service.
-   `GetOpenAIAdapterStatus() (*OpenAIAdapterStatus, error)`: Gets the current running status of the adapter service.
-   `GetAdapterAPIDocs() (map[string]string, error)`: Gets formatted API `curl` examples.

### API Debugger

-   `SendHttpRequest(request types.ApiRequest) (types.ApiResponse, error)`: Sends a complete, custom HTTP request to the specified Ollama service.
