package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Model 模型信息
type Model struct {
	Name       string                 `json:"name"`
	Model      string                 `json:"model"`
	ModifiedAt string                 `json:"modified_at"`
	Size       int64                  `json:"size"`
	Digest     string                 `json:"digest"`
	Details    map[string]interface{} `json:"details"`
	IsRunning  bool                   `json:"is_running"`
}

// ListModelsResponse 模型列表响应
type ListModelsResponse struct {
	Models []Model `json:"models"`
}

// Message 聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
}

// PullModelRequest 拉取模型请求
type PullModelRequest struct {
	Name string `json:"name"`
}

// OllamaServerConfig Ollama服务器配置
type OllamaServerConfig struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	BaseURL    string `json:"base_url"`
	APIKey     string `json:"api_key"`
	IsActive   bool   `json:"is_active"`
	TestStatus string `json:"test_status"`
	Type       string `json:"type"`
}

// App 应用结构体
type App struct {
	ctx          context.Context
	httpClient   *core.HttpCli
	store        *duolasdk.AppStore
	modelManager *ModelManager
	configMgr    *OllamaConfigManager
}

// NewApp 创建一个新的 App 应用
func NewApp() *App {
	// 初始化存储
	store := duolasdk.NewStore(
		core.StoreOption{
			FileName: "ollama-client.db",
		})
	// 创建应用实例
	app := &App{
		store: store,
	}

	// 初始化配置管理器
	app.configMgr = NewOllamaConfigManager(store)

	// 获取活动服务器配置
	activeServer, err := app.configMgr.GetActiveServer()
	if err != nil {
		// 如果没有活动服务器，尝试获取本地配置
		localConfig, localErr := app.configMgr.GetLocalConfig()
		if localErr != nil {
			// 如果本地配置也没有，使用默认值
			activeServer = &OllamaServerConfig{BaseURL: "http://localhost:11434"}
		} else {
			activeServer = &localConfig
		}
	}

	// 初始化HTTP客户端
	app.httpClient = core.NewHttp(core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "HttpClient",
	}))
	app.httpClient.Create(&core.Config{
		BaseURL: activeServer.BaseURL,
	})

	// 初始化模型管理器
	app.modelManager = NewModelManager(app)

	return app
}

// startup 在应用启动时调用，保存上下文
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.modelManager.SetContext(ctx)
}

// ListModelsByServer 根据服务器获取模型列表
func (a *App) ListModelsByServer(serverID string) ([]Model, error) {
	// 获取服务器配置
	var serverConfig *OllamaServerConfig

	if serverID == "local" {
		// 获取本地配置
		localConfig, err := a.configMgr.GetLocalConfig()
		if err != nil {
			return nil, err
		}
		serverConfig = &localConfig
	} else {
		// 获取远程服务器配置
		servers, err := a.configMgr.GetRemoteServers()
		if err != nil {
			return nil, err
		}

		found := false
		for _, server := range servers {
			if server.ID == serverID {
				serverConfig = &server
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("找不到服务器: %s", serverID)
		}
	}

	// 为特定服务器创建临时的HTTP客户端
	tempClient := core.NewHttp(core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "TempClient",
	}))

	tempClient.Create(&core.Config{
		BaseURL: serverConfig.BaseURL,
	})

	// 获取模型列表
	response, err := tempClient.Get("/api/tags", core.Options{})
	if err != nil {
		return nil, err
	}

	var result ListModelsResponse
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		return nil, err
	}

	// 检查并设置每个模型的运行状态
	for i := range result.Models {
		result.Models[i].IsRunning = a.modelManager.IsModelRunning(result.Models[i].Name)
	}

	return result.Models, nil
}

// DownloadModel 下载模型
func (a *App) DownloadModel(serverID string, modelName string) {
	a.modelManager.SetContext(a.ctx)
	go a.modelManager.DownloadModel(serverID, modelName)
}

// DeleteModel 删除模型
func (a *App) DeleteModel(modelName string) error {
	requestBody := map[string]interface{}{
		"name": modelName,
	}

	response, err := a.httpClient.Do("DELETE", "/api/delete", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})
	if err != nil {
		return err
	}

	// 检查响应状态
	if response.StatusCode >= 400 {
		return fmt.Errorf("删除模型失败，状态码: %d, 内容: %s", response.StatusCode, response.Body)
	}

	// 清除模型列表缓存
	a.store.Delete("models_cache")

	return nil
}

// RunModel 运行模型
func (a *App) RunModel(modelName string, params map[string]interface{}) error {
	modelParams := ModelParams{
		Temperature:   0.8,
		TopP:          0.9,
		TopK:          40,
		Context:       2048,
		NumPredict:    512,
		RepeatPenalty: 1.1,
	}

	// 从参数中提取值
	if temp, ok := params["temperature"].(float64); ok {
		modelParams.Temperature = temp
	}

	if topP, ok := params["top_p"].(float64); ok {
		modelParams.TopP = topP
	}

	if topK, ok := params["top_k"].(float64); ok {
		modelParams.TopK = int(topK)
	}

	if context, ok := params["context"].(float64); ok {
		modelParams.Context = int(context)
	}

	if numPredict, ok := params["num_predict"].(float64); ok {
		modelParams.NumPredict = int(numPredict)
	}

	if repeatPenalty, ok := params["repeat_penalty"].(float64); ok {
		modelParams.RepeatPenalty = repeatPenalty
	}

	return a.modelManager.RunModel(modelName, modelParams)
}

// StopModel 停止模型
func (a *App) StopModel(modelName string) error {
	return a.modelManager.StopModel(modelName)
}

// TestModel 测试模型
func (a *App) TestModel(modelName string, prompt string) (string, error) {
	return a.modelManager.TestModel(modelName, prompt)
}

// ChatMessage 发送聊天消息到Ollama API
func (a *App) ChatMessage(modelName string, messages []Message) (string, error) {
	requestBody := ChatRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   false,
	}

	response, err := a.httpClient.Post("/api/chat", core.Options{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    requestBody,
	})
	if err != nil {
		return "", fmt.Errorf("发送聊天消息失败: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("发送聊天消息失败: 状态码 %d", response.StatusCode)
	}

	var chatResponse ChatResponse
	err = json.Unmarshal([]byte(response.Body), &chatResponse)
	if err != nil {
		return "", fmt.Errorf("解析聊天响应失败: %w", err)
	}

	return chatResponse.Message.Content, nil
}

// SetModelParams 设置模型参数
func (a *App) SetModelParams(modelName string, params map[string]interface{}) error {
	modelParams := ModelParams{
		Temperature:   0.8,
		TopP:          0.9,
		TopK:          40,
		Context:       2048,
		NumPredict:    512,
		RepeatPenalty: 1.1,
	}

	// 从参数中提取值
	if temp, ok := params["temperature"].(float64); ok {
		modelParams.Temperature = temp
	}

	if topP, ok := params["top_p"].(float64); ok {
		modelParams.TopP = topP
	}

	if topK, ok := params["top_k"].(float64); ok {
		modelParams.TopK = int(topK)
	}

	if context, ok := params["context"].(float64); ok {
		modelParams.Context = int(context)
	}

	if numPredict, ok := params["num_predict"].(float64); ok {
		modelParams.NumPredict = int(numPredict)
	}

	if repeatPenalty, ok := params["repeat_penalty"].(float64); ok {
		modelParams.RepeatPenalty = repeatPenalty
	}

	return a.modelManager.SetModelParams(modelName, modelParams)
}

// GetModelParams 获取模型参数
func (a *App) GetModelParams(modelName string) (map[string]interface{}, error) {
	params, err := a.modelManager.GetModelParams(modelName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"temperature":    params.Temperature,
		"top_p":          params.TopP,
		"top_k":          params.TopK,
		"context":        params.Context,
		"num_predict":    params.NumPredict,
		"repeat_penalty": params.RepeatPenalty,
	}, nil
}

// GetOllamaServerConfig 获取本地Ollama服务器配置
func (a *App) GetOllamaServerConfig() (string, error) {
	config, err := a.configMgr.GetLocalConfig()
	if err != nil {
		return "http://localhost:11434", nil
	}
	return config.BaseURL, nil
}

// GetRemoteServers 获取远程Ollama服务器列表
func (a *App) GetRemoteServers() ([]OllamaServerConfig, error) {
	return a.configMgr.GetRemoteServers()
}

// SearchModels 在ollamadb.dev上搜索模型
func (a *App) SearchModels(query string) ([]interface{}, error) {
	return a.modelManager.SearchModels(query)
}

// TestOllamaServer 测试与Ollama服务器的连接
func (a *App) TestOllamaServer(baseURL string) (string, error) {
	client := core.NewHttp(core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "TestClient",
	}))
	client.Create(&core.Config{
		BaseURL: baseURL,
	})

	resp, err := client.Get("/", core.Options{})
	if err != nil {
		return "连接失败", err
	}

	if resp.StatusCode == http.StatusOK && resp.Body == "Ollama is running" {
		return "连接成功", nil
	}

	return "连接失败", fmt.Errorf("非预期的响应: %s", resp.Body)
}

// SaveOllamaServerConfig 保存本地Ollama服务器配置
func (a *App) SaveOllamaServerConfig(baseURL string) error {
	config, err := a.configMgr.GetLocalConfig()
	if err != nil {
		config = OllamaServerConfig{
			ID:   "local",
			Name: "本地服务",
			Type: "local",
		}
	}
	config.BaseURL = baseURL
	return a.configMgr.SaveLocalConfig(config)
}

// AddRemoteServer 添加一个新的远程服务器
func (a *App) AddRemoteServer(server OllamaServerConfig) error {
	return a.configMgr.AddRemoteServer(server)
}

// UpdateRemoteServer 更新一个已存在的远程服务器
func (a *App) UpdateRemoteServer(server OllamaServerConfig) error {
	return a.configMgr.UpdateRemoteServer(server)
}

// DeleteRemoteServer 删除一个远程服务器
func (a *App) DeleteRemoteServer(serverID string) error {
	return a.configMgr.DeleteRemoteServer(serverID)
}

// SetActiveServer 设置活动服务器
func (a *App) SetActiveServer(serverID string) error {
	return a.configMgr.SetActiveServer(serverID)
}

// SaveLocalServerTestStatus 保存本地服务器的测试状态
func (a *App) SaveLocalServerTestStatus(status string) error {
	return a.configMgr.SaveLocalServerTestStatus(status)
}

// GetLocalServerTestStatus 获取本地服务器的测试状态
func (a *App) GetLocalServerTestStatus() (string, error) {
	return a.configMgr.GetLocalServerTestStatus()
}

// GetActiveServer 获取活动服务器
func (a *App) GetActiveServer() (*OllamaServerConfig, error) {
	return a.configMgr.GetActiveServer()
}

// OpenInBrowser opens the given URL in the default browser.
func (a *App) OpenInBrowser(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// TestOllamaConnection 测试Ollama连接
func (a *App) TestOllamaConnection() (bool, error) {
	response, err := a.httpClient.Get("/api/tags", core.Options{})
	if err != nil {
		return false, err
	}

	if response.StatusCode != 200 {
		return false, fmt.Errorf("状态码错误: %d", response.StatusCode)
	}

	return true, nil
}
