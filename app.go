package main

import (
	"context"
	"fmt"
	"net/http"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	logger       *core.AppLog
	configMgr    *OllamaConfigManager
	chatManager  *ChatManager
	modelManager *ModelManager
	modelMarket  *ModelMarket
	promptPilot  *PromptPilot
	httpClient   *core.HttpCli // 全局httpClient，主要给ModelManager用
}

// NewApp 创建一个新的 App 应用
func NewApp() *App {
	// 1. 初始化核心依赖
	logger := core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "DuoLa",
	})

	store := duolasdk.NewStore(core.StoreOption{
		FileName: "ollama-client.db",
	})

	app := &App{
		logger: logger,
	}

	// 2. 初始化所有管理器，并注入依赖
	app.configMgr = NewOllamaConfigManager(store, logger)
	app.promptPilot = NewPromptPilot(store, app.configMgr, logger)
	app.chatManager = NewChatManager(context.Background(), store, logger)
	app.modelManager = NewModelManager(app, app.configMgr, logger)
	app.modelMarket = NewModelMarket(app, logger)

	// 3. 初始化全局HTTP客户端和AIProvider
	app.httpClient = core.NewHttp(logger.WithPrefix("HttpClient"))
	if err := app.rebuildDependencies(); err != nil {
		logger.Fatal("初始化应用依赖失败", "error", err)
		return nil
	}

	return app
}

// rebuildDependencies 根据当前活动服务器重建依赖（如httpClient和AIProvider）
func (a *App) rebuildDependencies() error {
	a.logger.Debug("正在根据活动服务器重建依赖")
	activeServer, err := a.configMgr.GetActiveServer()
	if err != nil {
		// 如果没有配置，这可能是首次启动，是正常情况
		a.logger.Warn("获取活动服务器失败，可能是首次启动", "error", err)
		// 即使没有活动服务器，也创建一个空的httpClient，避免nil panic
		a.httpClient.Create(&core.Config{})
		return nil
	}

	a.logger.Info("检测到活动服务器", "serverName", activeServer.Name, "serverURL", activeServer.BaseURL)

	// 更新HTTP客户端配置
	a.httpClient.Create(&core.Config{
		BaseURL: activeServer.BaseURL,
	})

	// 更新AI提供者
	ollamaProvider := core.NewOllamaProvider(activeServer.BaseURL)
	aiProviderAdapter := NewAIProviderAdapter(ollamaProvider, a.logger)
	a.chatManager.SetAIProvider(aiProviderAdapter)

	a.logger.Debug("依赖重建完成")
	return nil
}

// startup 在应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logger.Info("应用启动，正在设置所有模块的上下文")
	a.modelManager.SetContext(ctx)
	a.modelMarket.SetContext(ctx)
	a.chatManager.SetContext(ctx)
	a.promptPilot.Startup(ctx)
}

// --- PromptPilot Methods ---
func (a *App) GeneratePromptStream(idea string, model string, serverId string) {
	a.promptPilot.GeneratePromptStream(idea, model, serverId)
}
func (a *App) OptimizePrompt(content string, feedback string, model string, serverId string) (string, error) {
	return a.promptPilot.OptimizePrompt(content, feedback, model, serverId)
}
func (a *App) SavePrompt(prompt types.Prompt) error {
	return a.promptPilot.SavePrompt(prompt)
}
func (a *App) ListPrompts() ([]types.Prompt, error) {
	return a.promptPilot.ListPrompts()
}
func (a *App) GetPrompt(id string) (types.Prompt, error) {
	return a.promptPilot.GetPrompt(id)
}
func (a *App) DeletePrompt(id string) error {
	return a.promptPilot.DeletePrompt(id)
}

// --- ModelManager Methods ---
func (a *App) ListModelsByServer(serverID string) ([]types.Model, error) {
	return a.modelManager.ListModelsByServer(serverID)
}
func (a *App) DownloadModel(serverID string, modelName string) {
	go a.modelManager.DownloadModel(serverID, modelName)
}
func (a *App) DeleteModel(modelName string) error {
	return a.modelManager.DeleteModel(modelName)
}
func (a *App) RunModel(modelName string, params map[string]interface{}) error {
	// 参数转换逻辑保留在此处或移至ModelManager
	modelParams := types.ModelParams{
		Temperature: 0.8, TopP: 0.9, TopK: 40,
		Context: 2048, NumPredict: 512, RepeatPenalty: 1.1,
	}
	if temp, ok := params["temperature"].(float64); ok {
		modelParams.Temperature = temp
	}
	if topP, ok := params["topP"].(float64); ok {
		modelParams.TopP = topP
	}
	if topK, ok := params["topK"].(float64); ok {
		modelParams.TopK = int(topK)
	}
	if ct, ok := params["context"].(float64); ok {
		modelParams.Context = int(ct)
	}
	if numPredict, ok := params["numPredict"].(float64); ok {
		modelParams.NumPredict = int(numPredict)
	}
	if repeatPenalty, ok := params["repeatPenalty"].(float64); ok {
		modelParams.RepeatPenalty = repeatPenalty
	}

	return a.modelManager.RunModel(modelName, modelParams)
}
func (a *App) StopModel(modelName string) error {
	return a.modelManager.StopModel(modelName)
}
func (a *App) TestModel(modelName string, prompt string) (string, error) {
	return a.modelManager.TestModel(modelName, prompt)
}
func (a *App) SetModelParams(modelName string, params map[string]interface{}) error {
	// 参数转换逻辑
	modelParams := types.ModelParams{}
	if temp, ok := params["temperature"].(float64); ok {
		modelParams.Temperature = temp
	}
	// ... 其他参数 ...
	return a.modelManager.SetModelParams(modelName, modelParams)
}
func (a *App) GetModelParams(modelName string) (map[string]interface{}, error) {
	params, err := a.modelManager.GetModelParams(modelName)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"temperature":   params.Temperature,
		"topP":          params.TopP,
		"topK":          params.TopK,
		"context":       params.Context,
		"numPredict":    params.NumPredict,
		"repeatPenalty": params.RepeatPenalty,
	}, nil
}

// --- ChatManager Methods ---
func (a *App) ChatMessage(modelName string, messages []types.Message, stream bool) (string, error) {
	return a.chatManager.ChatMessage(modelName, messages, stream)
}
func (a *App) ListConversations() ([]*types.Conversation, error) {
	return a.chatManager.ListConversations()
}
func (a *App) SaveConversation(conv *types.Conversation) (*types.Conversation, error) {
	return a.chatManager.SaveConversation(conv)
}
func (a *App) GetConversation(id string) (*types.Conversation, error) {
	return a.chatManager.GetConversation(id)
}
func (a *App) DeleteConversation(id string) error {
	return a.chatManager.DeleteConversation(id)
}

// --- ConfigManager Methods ---
func (a *App) GetServers() ([]types.OllamaServerConfig, error) {
	return a.configMgr.GetServers()
}
func (a *App) AddServer(server types.OllamaServerConfig) error {
	return a.configMgr.AddServer(server)
}
func (a *App) UpdateServer(server types.OllamaServerConfig) error {
	return a.configMgr.UpdateServer(server)
}
func (a *App) DeleteServer(serverID string) error {
	return a.configMgr.DeleteServer(serverID)
}
func (a *App) SetActiveServer(serverID string) error {
	if err := a.configMgr.SetActiveServer(serverID); err != nil {
		return err
	}
	// 重新加载依赖项以反映服务器变化
	return a.rebuildDependencies()
}
func (a *App) GetActiveServer() (*types.OllamaServerConfig, error) {
	return a.configMgr.GetActiveServer()
}
func (a *App) UpdateServerTestStatus(serverID string, status string) error {
	return a.configMgr.UpdateServerTestStatus(serverID, status)
}
func (a *App) TestOllamaServer(baseURL string) (string, error) {
	client := core.NewHttp(a.logger.WithPrefix("TestClient"))
	client.Create(&core.Config{BaseURL: baseURL})
	resp, err := client.Get("/", core.Options{})
	if err != nil {
		return "连接失败", err
	}
	if resp.StatusCode == http.StatusOK && resp.Body == "Ollama is running" {
		return "连接成功", nil
	}
	return "连接失败", fmt.Errorf("非预期的响应: %s", resp.Body)
}

// --- ModelMarket Methods ---
func (a *App) SearchOnlineModels(query string) ([]interface{}, error) {
	return a.modelMarket.SearchOnlineModels(query)
}

// --- Browser/System Methods ---
func (a *App) OpenInBrowser(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// AIProviderAdapter 适配器，用于解决core.Message和本地Message类型不匹配的问题
type AIProviderAdapter struct {
	provider *core.OllamaProvider
	logger   *core.AppLog
}

// NewAIProviderAdapter 创建AI提供者适配器
func NewAIProviderAdapter(provider *core.OllamaProvider, logger *core.AppLog) *AIProviderAdapter {
	return &AIProviderAdapter{
		provider: provider,
		logger:   logger.WithPrefix("AIAdapter"),
	}
}

// Chat 适配Chat方法
func (a *AIProviderAdapter) Chat(model string, messages []core.Message) (string, error) {
	a.logger.Debug("Adapter: 开始阻塞式聊天请求", "model", model, "messageCount", len(messages))
	response, err := a.provider.Chat(model, messages)
	if err != nil {
		a.logger.Error("Adapter: 阻塞式聊天请求失败", "error", err)
		return "", err
	}
	return response, nil
}

// ChatStream 适配ChatStream方法
func (a *AIProviderAdapter) ChatStream(model string, messages []core.Message, callback func(string)) error {
	a.logger.Debug("Adapter: 开始流式聊天请求", "model", model, "messageCount", len(messages))
	err := a.provider.ChatStream(model, messages, callback)
	if err != nil {
		a.logger.Error("Adapter: 流式聊天请求失败", "error", err)
	}
	return err
}
