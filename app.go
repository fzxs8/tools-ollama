package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// ListModelsResponse 模型列表响应1
type ListModelsResponse struct {
	Models []Model `json:"models"`
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
	configMgr    *OllamaConfigManager
	chatManager  *ChatManager
	modelManager *ModelManager
	modelMarket  *ModelMarket
	promptPilot  *PromptPilot // 添加 PromptPilot
	httpClient   *core.HttpCli
	store        *duolasdk.AppStore
	aiProvider   interface {
		Chat(model string, messages []core.Message) (string, error)
		ChatStream(model string, messages []core.Message, callback func(string)) error
	}
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

	// 初始化聊天管理器
	app.chatManager = NewChatManager(context.Background(), store)

	// 初始化 PromptPilot
	app.promptPilot = NewPromptPilot(store, app.configMgr)

	// 获取活动服务器配置
	activeServer, err := app.configMgr.GetActiveServer()
	if err != nil {
		//fmt.Errorf("找不到服务器: %s", err)
		return nil
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

	// 初始化AI提供者
	ollamaProvider := core.NewOllamaProvider(activeServer.BaseURL)

	// 创建适配器以解决Message类型不匹配问题
	aiProviderAdapter := NewAIProviderAdapter(ollamaProvider)

	app.aiProvider = ollamaProvider
	app.chatManager.SetAIProvider(aiProviderAdapter)

	// 初始化模型管理器和模型市场
	app.modelManager = NewModelManager(app)
	app.modelMarket = NewModelMarket(app)

	return app
}

// startup 在应用启动时调用，保存上下文
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.modelManager.SetContext(ctx)
	a.modelMarket.SetContext(ctx)
	a.chatManager.SetContext(ctx)
	a.promptPilot.Startup(ctx) // 启动 PromptPilot
}

// --- PromptPilot Methods ---

// GeneratePromptStream 异步流式生成一个提示词
func (a *App) GeneratePromptStream(idea string, model string, serverId string) error {
	return a.promptPilot.GeneratePromptStream(idea, model, serverId)
}

// OptimizePrompt 优化一个提示词
func (a *App) OptimizePrompt(content string, feedback string, model string, serverId string) (string, error) {
	return a.promptPilot.OptimizePrompt(content, feedback, model, serverId)
}

// SavePrompt 保存一个提示词到存储
func (a *App) SavePrompt(prompt Prompt) error {
	return a.promptPilot.SavePrompt(prompt)
}

// ListPrompts 返回所有已保存的提示词
func (a *App) ListPrompts() ([]Prompt, error) {
	return a.promptPilot.ListPrompts()
}

// GetPrompt 根据ID返回一个特定的提示词
func (a *App) GetPrompt(id string) (Prompt, error) {
	return a.promptPilot.GetPrompt(id)
}

// DeletePrompt 根据ID删除一个提示词
func (a *App) DeletePrompt(id string) error {
	return a.promptPilot.DeletePrompt(id)
}

// --- End PromptPilot Methods ---

// KVSet 设置键值对
func (a *App) KVSet(key, value string) error {
	return a.chatManager.KVSet(key, value)
}

// KVGet 获取键值对
func (a *App) KVGet(key string) (string, error) {
	return a.chatManager.KVGet(key)
}

// KVList 获取键值对列表
func (a *App) KVList(key string) (string, error) {
	return a.chatManager.KVList(key)
}

// KVDelete 删除键值对
func (a *App) KVDelete(key string) error {
	return a.chatManager.KVDelete(key)
}

// ListModelsByServer 根据服务器获取模型列表
func (a *App) ListModelsByServer(serverID string) ([]Model, error) {
	// 获取服务器配置
	var serverConfig *OllamaServerConfig

	// 获取远程服务器配置
	servers, err := a.configMgr.GetServers()
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

	if ct, ok := params["context"].(float64); ok {
		modelParams.Context = int(ct)
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
func (a *App) ChatMessage(modelName string, messages []Message, stream bool) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ChatMessage方法中发生恐慌: %v", r)
		}
	}()

	log.Printf("ChatMessage调用: 模型=%s, 消息数量=%d, 是否流式=%t", modelName, len(messages), stream)

	// 如果使用流式传输
	if stream {
		log.Printf("使用流式传输")
		// 创建通道用于接收流式数据
		chunkChan := make(chan string, 100)
		errChan := make(chan error, 1)

		// 启动goroutine处理流式请求
		go func() {
			defer close(chunkChan)
			defer close(errChan)

			err := a.chatManager.ChatStream(modelName, messages, func(content string) {
				chunkChan <- content
			})
			errChan <- err
		}()

		// 通过Wails Events发送数据块
		result := ""
		for {
			select {
			case chunk, ok := <-chunkChan:
				if !ok {
					// 通道关闭，流式传输结束
					return result, nil
				}
				result += chunk
				// 发送事件到前端
				runtime.EventsEmit(a.ctx, "chat_stream_chunk", chunk)
			case err := <-errChan:
				if err != nil {
					return "", err
				}
				// 正常结束
				return result, nil
			}
		}
	} else {
		log.Printf("使用阻塞式传输")
		// 使用阻塞式传输
		result, err := a.chatManager.Chat(modelName, messages)
		if err != nil {
			log.Printf("阻塞式传输错误: %v", err)
			return "", err
		}
		log.Printf("阻塞式传输成功，结果长度=%d", len(result))
		log.Printf("阻塞式传输返回结果前100个字符: %s", func() string {
			if len(result) > 100 {
				return result[:100] + "..."
			}
			return result
		}())
		return result, nil
	}
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

// GetServers 获取远程Ollama服务器列表
func (a *App) GetServers() ([]OllamaServerConfig, error) {
	return a.configMgr.GetServers()
}

// SearchOnlineModels 搜索在线模型
func (a *App) SearchOnlineModels(query string) ([]interface{}, error) {
	return a.modelMarket.SearchOnlineModels(query)
}

// TestOllamaServer 测试与Ollama服务器的连接1
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
		return "连接失败：", err
	}

	if resp.StatusCode == http.StatusOK && resp.Body == "Ollama is running" {
		return "连接成功", nil
	}
	log.Println("非预期的响应:", resp.Body)
	return "连接失败", fmt.Errorf("非预期的响应: %s", resp.Body)
}

// SaveOllamaServerConfig 保存本地Ollama服务器配置
//func (a *App) SaveOllamaServerConfig(config []OllamaServerConfig) error {
//	//config, err := a.configMgr.GetConfig()
//	//if err != nil {
//	//	config = OllamaServerConfig{
//	//		ID:   "local",
//	//		Name: "本地服务",
//	//		Type: "local",
//	//	}
//	//}
//	//config.BaseURL = baseURL
//	return a.configMgr.SaveServers(config)
//}

// AddServer 添加一个新的远程服务器
func (a *App) AddServer(server OllamaServerConfig) error {
	return a.configMgr.AddServer(server)
}
func (a *App) UpdateServerTestStatus(serverID string, status string) error {
	return a.configMgr.UpdateServerTestStatus(serverID, status)
}

// UpdateServer 更新一个已存在的远程服务器
func (a *App) UpdateServer(server OllamaServerConfig) error {
	return a.configMgr.UpdateServer(server)
}

// DeleteServer 删除一个远程服务器
func (a *App) DeleteServer(serverID string) error {
	return a.configMgr.DeleteServer(serverID)
}

// SetActiveServer 设置活动服务器
func (a *App) SetActiveServer(serverID string) error {
	err := a.configMgr.SetActiveServer(serverID)
	if err != nil {
		return err
	}

	// 获取新的活动服务器配置
	activeServer, err := a.configMgr.GetActiveServer()
	if err != nil {
		return err
	}

	// 更新HTTP客户端配置
	a.httpClient.Create(&core.Config{
		BaseURL: activeServer.BaseURL,
	})

	// 更新AI提供者
	ollamaProvider := core.NewOllamaProvider(activeServer.BaseURL)

	// 创建适配器以解决Message类型不匹配问题
	aiProviderAdapter := NewAIProviderAdapter(ollamaProvider)

	a.aiProvider = ollamaProvider
	a.chatManager.SetAIProvider(aiProviderAdapter)

	return nil
}

//// SaveLocalServerTestStatus 保存本地服务器的测试状态
//func (a *App) SaveLocalServerTestStatus(status string) error {
//	return a.configMgr.SaveLocalServerTestStatus(status)
//}

//// GetLocalServerTestStatus 获取本地服务器的测试状态
//func (a *App) GetLocalServerTestStatus() (string, error) {
//	return a.configMgr.GetLocalServerTestStatus()
//}

// GetActiveServer 获取活动服务器1
func (a *App) GetActiveServer() (*OllamaServerConfig, error) {
	return a.configMgr.GetActiveServer()
}

// OpenInBrowser opens the given URL in the default browser.
func (a *App) OpenInBrowser(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// TestOllamaConnection 测试Ollama连接1
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

// ListConversations 获取所有已保存的对话列表，按时间倒序排列
func (a *App) ListConversations() ([]*Conversation, error) {
	return a.chatManager.ListConversations()
}

// SaveConversation 创建或更新一个对话
func (a *App) SaveConversation(conv *Conversation) (*Conversation, error) {
	return a.chatManager.SaveConversation(conv)
}

// GetConversation 获取指定ID的单个对话的完整内容
func (a *App) GetConversation(id string) (*Conversation, error) {
	return a.chatManager.GetConversation(id)
}

// DeleteConversation 删除指定ID的对话
func (a *App) DeleteConversation(id string) error {
	return a.chatManager.DeleteConversation(id)
}

// AIProviderAdapter 适配器，用于解决core.Message和本地Message类型不匹配的问题
type AIProviderAdapter struct {
	provider *core.OllamaProvider
	logger   *core.AppLog
}

// NewAIProviderAdapter 创建AI提供者适配器
func NewAIProviderAdapter(provider *core.OllamaProvider) *AIProviderAdapter {
	logger := core.NewLogger(&core.LoggerOption{Type: "console", Level: "debug", Prefix: "AIProviderAdapter"})
	return &AIProviderAdapter{
		provider: provider,
		logger:   logger,
	}
}

// Chat 适配Chat方法
func (a *AIProviderAdapter) Chat(model string, messages []Message) (string, error) {
	// 验证模型名称
	if model == "" {
		a.logger.Warn("模型名称不能为空")
		return "", fmt.Errorf("模型名称不能为空")
	}

	// 验证消息数组
	if len(messages) == 0 {
		a.logger.Warn("消息列表不能为空")
		return "", fmt.Errorf("消息列表不能为空")
	}

	// 将本地Message类型转换为core.Message类型
	coreMessages := make([]core.Message, len(messages))
	for i, msg := range messages {
		if msg.Role == "" || msg.Content == "" {
			a.logger.Warn("消息 %d 缺少必要字段: role 或 content", i)
			return "", fmt.Errorf("消息 %d 缺少必要字段: role 或 content", i)
		}
		coreMessages[i] = core.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	a.logger.Debug("开始阻塞式聊天请求: 模型=%s, 消息数量=%d", model, len(messages))

	// 调用核心库的Chat方法，并处理可能的错误
	response, err := a.provider.Chat(model, coreMessages)
	if err != nil {
		a.logger.Error("阻塞式聊天请求失败: %v", err)
		return "", fmt.Errorf("调用Ollama服务失败: %w", err)
	}

	a.logger.Debug("阻塞式聊天请求成功，结果长度=%d", len(response))

	// 安全获取前100个字符
	var preview string
	if len(response) > 100 {
		// 确保在字符串边界截断，避免截断多字节字符
		preview = string([]rune(response)[:100]) + "..."
	} else {
		preview = response
	}
	a.logger.Debug("阻塞式聊天请求返回结果前100个字符: %s", preview)

	return response, nil
}

// ChatStream 适配ChatStream方法
func (a *AIProviderAdapter) ChatStream(model string, messages []Message, callback func(string)) error {
	// 将本地Message类型转换为core.Message类型
	coreMessages := make([]core.Message, len(messages))
	for i, msg := range messages {
		coreMessages[i] = core.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	a.logger.Debug("开始流式聊天请求: 模型=%s, 消息数量=%d", model, len(messages))

	// 创建一个安全的回调函数包装器
	safeCallback := func(content string) {
		// recover任何可能的恐慌
		defer func() {
			if r := recover(); r != nil {
				a.logger.Error("回调函数中发生恐慌: %v", r)
			}
		}()

		// 记录接收到的内容长度和内容预览
		a.logger.Debug("流式回调接收到内容，长度=%d，内容=%s", len(content), func() string {
			if len(content) > 50 {
				return content[:50] + "..."
			}
			return content
		}())

		// 调用原始回调函数
		if callback != nil {
			callback(content)
		}
	}

	err := a.provider.ChatStream(model, coreMessages, safeCallback)
	if err != nil {
		a.logger.Error("流式聊天请求失败: %v", err)
		return err
	}
	a.logger.Debug("流式聊天请求完成")
	return nil
}
