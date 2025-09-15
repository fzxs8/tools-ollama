package types

import "time"

// Model 模型信息
type Model struct {
	Name       string                 `json:"name"`
	Model      string                 `json:"model"`
	ModifiedAt string                 `json:"modified_at"`
	Size       int64                  `json:"size"`
	Digest     string                 `json:"digest"`
	Details    map[string]interface{} `json:"details"`
	IsRunning  bool                   `json:"isRunning"`
}

// Conversation 定义了一个完整的对话会话
type Conversation struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Messages     []Message `json:"messages"`
	ModelName    string    `json:"modelName"`
	SystemPrompt string    `json:"systemPrompt"` // JSON string of the active system prompt
	ModelParams  string    `json:"modelParams"`  // JSON string of the model parameters
	Timestamp    int64     `json:"timestamp"`
}

// Message 聊天消息结构
type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ListModelsResponse 模型列表响应
type ListModelsResponse struct {
	Models []Model `json:"models"`
}

// OllamaServerConfig Ollama服务器配置
type OllamaServerConfig struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	BaseURL    string `json:"baseUrl"`
	APIKey     string `json:"apiKey"`
	IsActive   bool   `json:"isActive"`
	TestStatus string `json:"testStatus"`
}

// Prompt 代表一个已保存的提示词
type Prompt struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	CreatedAt   int64    `json:"createdAt"`
	UpdatedAt   int64    `json:"updatedAt"`
	Models      []string `json:"models"`
	Version     int      `json:"version"`
	Tags        []string `json:"tags"`
	CreatedBy   string   `json:"createdBy"`
}

// ModelParams 模型参数
type ModelParams struct {
	Temperature   float64 `json:"temperature"`
	TopP          float64 `json:"topP"`
	TopK          int     `json:"topK"`
	Context       int     `json:"context"`
	NumPredict    int     `json:"numPredict"`
	RepeatPenalty float64 `json:"repeatPenalty"`
}

// RunningModel 运行中的模型
type RunningModel struct {
	Name      string      `json:"name"`
	Params    ModelParams `json:"params"`
	StartTime time.Time   `json:"startTime"`
	IsActive  bool        `json:"isActive"`
}

// OnlineModel 在线模型信息
type OnlineModel struct {
	Name        string `json:"name"`
	PullCount   int64  `json:"pullCount"`
	UpdatedAt   string `json:"updatedAt"`
	Description string `json:"description"`
}

// --- Ollama API Debugger Types ---

// RequestHeader 请求头结构
type RequestHeader struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

// QueryParam 查询参数结构
type QueryParam struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

// RequestBodyType 请求体类型枚举
type RequestBodyType string

const (
	RequestBodyTypeNone     RequestBodyType = "none"
	RequestBodyTypeRaw      RequestBodyType = "raw"
	RequestBodyTypeFormData RequestBodyType = "formData"
)

// RawBodyContentType 原始请求体内容类型枚举
type RawBodyContentType string

const (
	RawBodyContentTypeJson RawBodyContentType = "application/json"
	RawBodyContentTypeText RawBodyContentType = "text/plain"
	RawBodyContentTypeHtml RawBodyContentType = "text/html"
	RawBodyContentTypeXml  RawBodyContentType = "application/xml"
)

// FormDataItem 用于表示表单数据中的键值对
type FormDataItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RequestBody 请求体结构
type RequestBody struct {
	Type           RequestBodyType    `json:"type"`
	RawContent     string             `json:"rawContent,omitempty"`     // For Raw type
	RawContentType RawBodyContentType `json:"rawContentType,omitempty"` // For Raw type
	FormData       []FormDataItem     `json:"formData,omitempty"`       // For FormData type
}

// ApiRequest API请求结构
type ApiRequest struct {
	Method           string          `json:"method"`           // HTTP 方法 (GET, POST, etc.)
	SelectedServerID string          `json:"selectedServerId"` // 选中的 Ollama 服务 ID
	Path             string          `json:"path"`             // API 路径 (不包含 Base URL)
	QueryParams      []QueryParam    `json:"queryParams"`
	Headers          []RequestHeader `json:"headers"`
	Body             RequestBody     `json:"body"`
}

// ApiResponse API响应结构
type ApiResponse struct {
	StatusCode        int             `json:"statusCode"`
	StatusText        string          `json:"statusText"`
	Headers           []RequestHeader `json:"headers"` // Changed to RequestHeader to match frontend
	Body              string          `json:"body"`
	RequestDurationMs int64           `json:"requestDurationMs"` // 请求耗时 (毫秒)
	Error             string          `json:"error,omitempty"`   // 错误信息
}

// --- OpenAI Adapter Types ---

// OpenAIAdapterConfig holds the user-configurable settings for the adapter.
type OpenAIAdapterConfig struct {
	ListenIP             string `json:"listenIp"`
	ListenPort           int    `json:"listenPort"`
	TargetOllamaServerID string `json:"targetOllamaServerId"`
}

// OpenAIAdapterStatus represents the current runtime status of the adapter service.
type OpenAIAdapterStatus struct {
	IsRunning bool   `json:"isRunning"`
	Error     string `json:"error,omitempty"` // Holds error message if startup failed
}

// LogEntry represents a single log message sent to the frontend.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // e.g., "INFO", "ERROR", "DEBUG"
	Message   string    `json:"message"`
}
