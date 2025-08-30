package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
)

// Model 模型信息
type Model struct {
	Name       string                 `json:"name"`
	Model      string                 `json:"model"`
	ModifiedAt string                 `json:"modified_at"`
	Size       int64                  `json:"size"`
	Digest     string                 `json:"digest"`
	Details    map[string]interface{} `json:"details"`
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
}

// App struct
type App struct {
	ctx          context.Context
	httpClient   *core.HttpCli
	store        *duolasdk.AppStore
	modelManager *ModelManager
	configMgr    *OllamaConfigManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	// 初始化存储
	store := duolasdk.NewStore(
		core.StoreOption{
			FileName: "ollama-client.db",
		})
	// 创建存储实例
	// 创建应用实例
	app := &App{
		store: store,
	}

	// 初始化配置管理器
	app.configMgr = NewOllamaConfigManager(store)

	// 初始化模型管理器
	app.modelManager = NewModelManager(app)

	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.modelManager.SetContext(ctx)
}

// ListModels 获取模型列表
func (a *App) ListModels() ([]Model, error) {
	// 首先尝试从存储中获取缓存的模型列表
	cached, err := a.store.Get("models_cache")
	if err == nil {
		var models []Model
		if json.Unmarshal([]byte(cached), &models) == nil {
			return models, nil
		}
	}

	// 如果缓存不存在或解析失败，则从API获取
	response, err := a.httpClient.Get("/api/tags", core.Options{})
	if err != nil {
		return nil, err
	}

	var result ListModelsResponse
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		return nil, err
	}

	// 将结果缓存到存储中
	data, _ := json.Marshal(result.Models)
	a.store.Set("models_cache", string(data))

	return result.Models, nil
}

// ListModelsByServer 根据服务器获取模型列表
func (a *App) ListModelsByServer(serverID string) ([]Model, error) {
	// 获取服务器配置
	var serverConfig *OllamaServerConfig

	if serverID == "local" {
		// 获取本地配置
		baseURL, err := a.configMgr.GetLocalConfig()
		if err != nil {
			baseURL = "http://localhost:11434"
		}
		serverConfig = &OllamaServerConfig{
			ID:      "local",
			Name:    "本地服务",
			BaseURL: baseURL,
			APIKey:  "",
		}
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
			return nil, fmt.Errorf("server not found: %s", serverID)
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

	return result.Models, nil
}

// ChatMessage 发送聊天消息
func (a *App) ChatMessage(model string, messages []Message) (string, error) {
	request := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// 将请求数据转换为map[string]interface{}
	var requestBody map[string]interface{}
	if err := json.Unmarshal(data, &requestBody); err != nil {
		return "", err
	}

	response, err := a.httpClient.Post("/api/chat", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("chat failed with status: %d, body: %s", response.StatusCode, response.Body)
	}

	// 解析响应
	var chatResponse ChatResponse
	if err := json.Unmarshal([]byte(response.Body), &chatResponse); err != nil {
		return "", err
	}

	return chatResponse.Message.Content, nil
}

// PullModel 拉取模型
func (a *App) PullModel(modelName string) error {
	request := PullModelRequest{
		Name: modelName,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// 将请求数据转换为map[string]interface{}
	var requestBody map[string]interface{}
	if err := json.Unmarshal(data, &requestBody); err != nil {
		return err
	}

	response, err := a.httpClient.Post("/api/pull", core.Options{
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
		return fmt.Errorf("pull model failed with status: %d, body: %s", response.StatusCode, response.Body)
	}

	// 清除模型列表缓存
	a.store.Delete("models_cache")

	return nil
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
		return fmt.Errorf("delete model failed with status: %d, body: %s", response.StatusCode, response.Body)
	}

	// 清除模型列表缓存
	a.store.Delete("models_cache")

	return nil
}

// LoadModel 加载模型
func (a *App) LoadModel(modelName string) error {
	// 在Ollama中，模型会在首次使用时自动加载
	// 这里我们可以通过发送一个简单的请求来触发加载
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "hello",
	}

	_, err := a.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})
	if err != nil {
		return err
	}

	return nil
}

// UnloadModel 卸载模型
func (a *App) UnloadModel(modelName string) error {
	// Ollama会自动管理模型的加载和卸载
	// 我们可以通过发送一个请求来触发卸载
	// 但实际的卸载由Ollama服务管理
	requestBody := map[string]interface{}{
		"model":      modelName,
		"keep_alive": "0",
	}

	_, err := a.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})
	return err
}

// SaveModelSettings 保存模型设置
func (a *App) SaveModelSettings(modelName string, options map[string]interface{}) error {
	key := fmt.Sprintf("model_options_%s", modelName)
	data, err := json.Marshal(options)
	if err != nil {
		return err
	}

	return a.store.Set(key, string(data))
}

// GetModelSettings 获取模型设置
func (a *App) GetModelSettings(modelName string) (map[string]interface{}, error) {
	key := fmt.Sprintf("model_options_%s", modelName)
	data, err := a.store.Get(key)
	if err != nil {
		// 返回默认设置
		return map[string]interface{}{
			"temperature": 0.8,
			"top_p":       0.9,
			"context":     2048,
		}, nil
	}

	var options map[string]interface{}
	if err := json.Unmarshal([]byte(data), &options); err != nil {
		return nil, err
	}

	return options, nil
}

// SaveOllamaServerConfig 保存Ollama服务器配置
func (a *App) SaveOllamaServerConfig(baseUrl string) error {
	// 更新HTTP客户端的基础URL
	a.httpClient.Create(&core.Config{
		BaseURL: baseUrl,
	})

	// 保存到存储中
	return a.configMgr.SaveLocalConfig(baseUrl)
}

// GetOllamaServerConfig 获取Ollama服务器配置
func (a *App) GetOllamaServerConfig() (string, error) {
	url, err := a.configMgr.GetLocalConfig()
	if err != nil {
		// 返回默认URL
		return "http://localhost:11434", nil
	}
	return url, nil
}

// SaveRemoteServers 保存远程服务器列表
func (a *App) SaveRemoteServers(servers []OllamaServerConfig) error {
	return a.configMgr.SaveRemoteServers(servers)
}

// GetRemoteServers 获取远程服务器列表
func (a *App) GetRemoteServers() ([]OllamaServerConfig, error) {
	return a.configMgr.GetRemoteServers()
}

// AddRemoteServer 添加远程服务器
func (a *App) AddRemoteServer(server OllamaServerConfig) error {
	return a.configMgr.AddRemoteServer(server)
}

// UpdateRemoteServer 更新远程服务器
func (a *App) UpdateRemoteServer(server OllamaServerConfig) error {
	return a.configMgr.UpdateRemoteServer(server)
}

// DeleteRemoteServer 删除远程服务器
func (a *App) DeleteRemoteServer(serverID string) error {
	return a.configMgr.DeleteRemoteServer(serverID)
}

// SetActiveServer 设置活动服务器
func (a *App) SetActiveServer(serverID string) error {
	return a.configMgr.SetActiveServer(serverID)
}

// GetActiveServer 获取活动服务器
func (a *App) GetActiveServer() (*OllamaServerConfig, error) {
	return a.configMgr.GetActiveServer()
}

// TestOllamaServer 测试Ollama服务器是否可达
func (a *App) TestOllamaServer(baseURL string) error {
	// 为测试创建临时的HTTP客户端
	tempClient := core.NewHttp(core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "TestClient",
	}))

	tempClient.Create(&core.Config{
		BaseURL: baseURL,
	})

	// 尝试访问/tags端点来验证服务器是否可达
	_, err := tempClient.Get("/api/tags", core.Options{})
	if err != nil {
		return fmt.Errorf("无法连接到Ollama服务器: %w", err)
	}

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

// GetModelStatus 获取模型状态
func (a *App) GetModelStatus(modelName string) (map[string]interface{}, error) {
	status, err := a.modelManager.GetModelStatus(modelName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":       status.Name,
		"start_time": status.StartTime.Format("2006-01-02 15:04:05"),
		"is_active":  status.IsActive,
		"params": map[string]interface{}{
			"temperature":    status.Params.Temperature,
			"top_p":          status.Params.TopP,
			"top_k":          status.Params.TopK,
			"context":        status.Params.Context,
			"num_predict":    status.Params.NumPredict,
			"repeat_penalty": status.Params.RepeatPenalty,
		},
	}, nil
}

// ListRunningModels 获取运行中的模型列表
func (a *App) ListRunningModels() []map[string]interface{} {
	models := a.modelManager.ListRunningModels()

	var result []map[string]interface{}
	for _, model := range models {
		result = append(result, map[string]interface{}{
			"name":       model.Name,
			"start_time": model.StartTime.Format("2006-01-02 15:04:05"),
			"is_active":  model.IsActive,
		})
	}

	return result
}

// TestModel 测试模型
func (a *App) TestModel(modelName string) (string, error) {
	return a.modelManager.TestModel(modelName)
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

// SearchModels 搜索模型
func (a *App) SearchModels(params map[string]interface{}) ([]Model, error) {
	searchParams := ModelSearchParams{}

	if query, ok := params["query"].(string); ok {
		searchParams.Query = query
	}

	if families, ok := params["families"].([]interface{}); ok {
		for _, f := range families {
			if family, ok := f.(string); ok {
				searchParams.Families = append(searchParams.Families, family)
			}
		}
	}

	if tags, ok := params["tags"].([]interface{}); ok {
		for _, t := range tags {
			if tag, ok := t.(string); ok {
				searchParams.Tags = append(searchParams.Tags, tag)
			}
		}
	}

	return a.modelManager.SearchModels(searchParams)
}

// GetModelFamilies 获取模型家族列表
func (a *App) GetModelFamilies() []string {
	return a.modelManager.GetModelFamilies()
}

// GetModelTags 获取模型标签列表
func (a *App) GetModelTags() []string {
	return a.modelManager.GetModelTags()
}
