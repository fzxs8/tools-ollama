package types

import "time"

// Model 模型信息
type Model struct {
	Name       string                 `json:"name"`
	Model      string                 `json:"model"`
	ModifiedAt string                 `json:"modifiedAt"`
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
	TestStatus string `json:"testtatus"`
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
