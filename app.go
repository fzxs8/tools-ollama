package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"tools-ollama/types"

	"github.com/16chusi/duolasdk"
	"github.com/16chusi/duolasdk/core"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx               context.Context
	logger            *core.AppLog
	configMgr         *OllamaConfigManager
	chatManager       *ChatManager
	modelManager      *ModelManager
	modelMarket       *ModelMarket
	promptEngineering *PromptEngineering
	ollamaApiDebugger *OllamaApiDebugger
	httpClient        *core.HttpCli
	adapterManager    *OpenAIAdapterManager
}

// NewApp 创建一个新的 App 应用
func NewApp() *App {
	logger := core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "DuoLa",
	})

	// 创建数据目录
	os.MkdirAll("data", 0755)

	store := duolasdk.NewStore(core.StoreOption{
		FileName: "tools-ollama",
	})

	app := &App{
		logger: logger,
	}

	app.configMgr = NewOllamaConfigManager(store, logger)
	app.promptEngineering = NewPromptPilot(store, app.configMgr, logger)
	app.chatManager = NewChatManager(context.Background(), store, logger)
	app.modelManager = NewModelManager(app, app.configMgr, logger)
	app.modelMarket = NewModelMarket(app, logger)
	app.ollamaApiDebugger = NewOllamaApiDebugger(logger, app.configMgr)
	app.adapterManager = NewOpenAIAdapterManager(logger, store, app.configMgr)

	// 设置ChatManager的AIProvider
	app.chatManager.SetAIProvider(NewAIProviderAdapter(app.modelManager, logger))

	app.httpClient = core.NewHttp(logger.WithPrefix("HttpClient"))
	if err := app.rebuildDependencies(); err != nil {
		logger.Fatal("初始化应用依赖失败", "error", err)
		return nil
	}

	return app
}

// rebuildDependencies 根据当前活动服务器重建依赖
func (a *App) rebuildDependencies() error {
	a.logger.Debug("开始重建应用依赖")

	// 获取活动服务器配置
	activeServer, err := a.configMgr.GetActiveServer()
	if err != nil {
		a.logger.Warn("未找到活动服务器，使用默认配置", "error", err)
		// 如果没有活动服务器，创建一个默认的HTTP客户端
		a.httpClient.Create(&core.Config{
			BaseURL: "http://localhost:11434",
		})
		return nil
	}

	a.logger.Info("使用活动服务器配置初始化HTTP客户端", "serverName", activeServer.Name, "baseURL", activeServer.BaseURL)

	// 确保BaseURL包含协议前缀
	baseURL := EnsureHTTPPrefix(activeServer.BaseURL)
	if baseURL != activeServer.BaseURL {
		a.logger.Debug("为BaseURL添加http://前缀", "originalURL", activeServer.BaseURL, "newURL", baseURL)
	}

	// 重新配置HTTP客户端
	a.httpClient.Create(&core.Config{
		BaseURL: baseURL,
	})

	a.logger.Debug("HTTP客户端重建完成", "baseURL", baseURL)
	return nil
}

// startup 在应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logger.Info("应用启动，正在设置所有模块的上下文")
	a.modelManager.SetContext(ctx)
	a.modelMarket.SetContext(ctx)
	a.chatManager.SetContext(ctx)
	a.promptEngineering.Startup(ctx)
	a.ollamaApiDebugger.SetContext(ctx)
	a.adapterManager.SetContext(ctx) // 注入 Wails 上下文到适配器管理器
}

// --- OpenAI Adapter Manager Methods (Wails API) ---

func (a *App) GetOpenAIAdapterConfig() types.OpenAIAdapterConfig {
	a.logger.Info("[PANIC_DEBUG] 1. Entering GetOpenAIAdapterConfig")
	if a.adapterManager == nil {
		a.logger.Error("[PANIC_DEBUG] FATAL: a.adapterManager is NIL! This should not happen.")
		// 返回一个值以防止进一步的恐慌, 但根本问题需要解决
		return types.OpenAIAdapterConfig{}
	}
	a.logger.Info("[PANIC_DEBUG] 2. a.adapterManager is valid, proceeding to call GetConfig().")
	result := a.adapterManager.GetConfig()
	a.logger.Info("[PANIC_DEBUG] 5. Successfully returned from GetConfig().")
	return result
}

func (a *App) SaveOpenAIAdapterConfig(config types.OpenAIAdapterConfig) error {
	return a.adapterManager.SaveConfig(config)
}

func (a *App) StartAdapterServer() error {
	return a.adapterManager.Start()
}

func (a *App) StopAdapterServer() error {
	return a.adapterManager.Stop(a.ctx)
}

func (a *App) GetOpenAIAdapterStatus() types.OpenAIAdapterStatus {
	return a.adapterManager.GetStatus()
}

// GetAdapterAPIDocs 获取 API 文档
func (a *App) GetAdapterAPIDocs() (map[string]string, error) {
	cfg := a.adapterManager.GetConfig()
	return GenerateAPIDocs(cfg.ListenIP, cfg.ListenPort), nil
}

// ... (其他 App 方法保持不变) ...

// --- OllamaApiDebugger Methods ---
func (a *App) SendHttpRequest(request types.ApiRequest) (types.ApiResponse, error) {
	return a.ollamaApiDebugger.SendHttpRequest(request)
}

func (a *App) GetOllamaServers() ([]types.OllamaServerConfig, error) {
	return a.ollamaApiDebugger.GetOllamaServers()
}

// --- PromptEngineering Methods ---
func (a *App) GeneratePromptStream(idea string, model string, serverId string) {
	a.promptEngineering.GeneratePromptStream(idea, model, serverId)
}
func (a *App) OptimizePrompt(content string, feedback string, model string, serverId string) (string, error) {
	return a.promptEngineering.OptimizePrompt(content, feedback, model, serverId)
}
func (a *App) SavePrompt(prompt types.Prompt) error {
	return a.promptEngineering.SavePrompt(prompt)
}
func (a *App) ListPrompts() ([]types.Prompt, error) {
	return a.promptEngineering.ListPrompts()
}
func (a *App) GetPrompt(id string) (types.Prompt, error) {
	return a.promptEngineering.GetPrompt(id)
}
func (a *App) DeletePrompt(id string) error {
	return a.promptEngineering.DeletePrompt(id)
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
	modelParams := ConvertToModelParams(params)
	return a.modelManager.RunModel(modelName, modelParams)
}
func (a *App) StopModel(modelName string) error {
	return a.modelManager.StopModel(modelName)
}
func (a *App) TestModel(modelName string, prompt string) (string, error) {
	return a.modelManager.TestModel(modelName, prompt)
}
func (a *App) SetModelParams(modelName string, params map[string]interface{}) error {
	modelParams := ConvertToModelParams(params)
	return a.modelManager.SetModelParams(modelName, modelParams)
}
func (a *App) GetModelParams(modelName string) (map[string]interface{}, error) {
	params, err := a.modelManager.GetModelParams(modelName)
	if err != nil {
		return nil, err
	}
	return ConvertFromModelParams(params), nil
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
	modelManager *ModelManager
	logger       *core.AppLog
}

// NewAIProviderAdapter 创建AI提供者适配器
func NewAIProviderAdapter(modelManager *ModelManager, logger *core.AppLog) *AIProviderAdapter {
	return &AIProviderAdapter{
		modelManager: modelManager,
		logger:       logger.WithPrefix("AIAdapter"),
	}
}

// Chat 适配Chat方法
func (a *AIProviderAdapter) Chat(model string, messages []core.Message) (string, error) {
	a.logger.Debug("Adapter: 开始阻塞式聊天请求", "model", model, "messageCount", len(messages))
	response, err := a.modelManager.Chat(model, messages)
	if err != nil {
		a.logger.Error("Adapter: 阻塞式聊天请求失败", "error", err)
		return "", err
	}
	return response, nil
}

// ChatStream 适配ChatStream方法
func (a *AIProviderAdapter) ChatStream(model string, messages []core.Message, callback func(string)) error {
	a.logger.Debug("Adapter: 开始流式聊天请求", "model", model, "messageCount", len(messages))
	err := a.modelManager.ChatStream(model, messages, callback)
	if err != nil {
		a.logger.Error("Adapter: 流式聊天请求失败", "error", err)
	}
	return err
}
