package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/fzxs8/duolasdk/core/ai"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const AdapterConfigKey = "openai_adapter_config"

// WailsLogger 实现了 io.Writer 接口, 用于解析日志级别并将其写入 Wails 事件
type WailsLogger struct {
	ctx context.Context
}

// Write 会解析 p []byte 中由 gommon/log 生成的日志字符串, 提取级别, 并发送到前端
func (w *WailsLogger) Write(p []byte) (n int, err error) {
	if w.ctx == nil {
		return 0, errors.New("WailsLogger: context is nil")
	}

	message := string(p)
	level := "INFO" // 默认级别

	// 基于 gommon/log 的默认头格式, 从消息中解析日志级别
	if strings.Contains(message, " DEBUG ") {
		level = "DEBUG"
	} else if strings.Contains(message, " INFO ") {
		level = "INFO"
	} else if strings.Contains(message, " WARN ") {
		level = "WARN"
	} else if strings.Contains(message, " ERROR ") {
		level = "ERROR"
	}

	// 为了不在前端重复显示, 我们从原始消息中移除日志头
	// gommon/log 的格式是: `[PREFIX]- TIMESTAMP LEVEL FILE:LINE MSG`
	var cleanMessage string
	// 找到日志消息主体开始的位置
	parts := strings.SplitN(message, " ", 5) // 分割成5部分: [PREFIX]-, TIMESTAMP, LEVEL, FILE:LINE, MSG
	if len(parts) == 5 {
		cleanMessage = parts[4]
	} else {
		cleanMessage = message // 如果格式不匹配, 返回原始消息
	}

	entry := types.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   strings.TrimSpace(cleanMessage),
	}
	runtime.EventsEmit(w.ctx, "openai-adapter-log", entry)
	return len(p), nil
}

// OpenAIAdapterManager 负责管理适配器服务的生命周期和配置
type OpenAIAdapterManager struct {
	// log        *core.AppLog // 移除内部log, 避免潜在的nil问题
	store      *duolasdk.AppStore
	config     types.OpenAIAdapterConfig
	runtimeCtx context.Context // Wails runtime context for eventing

	server *http.Server
	status types.OpenAIAdapterStatus
	// mu         sync.Mutex // 移除互斥锁, 简化逻辑
	configMgr *OllamaConfigManager
}

// NewOpenAIAdapterManager 创建一个新的管理器实例
func NewOpenAIAdapterManager(log *core.AppLog, store *duolasdk.AppStore, configMgr *OllamaConfigManager) *OpenAIAdapterManager {
	m := &OpenAIAdapterManager{
		// log:       log.WithPrefix("AdapterManager"),
		store:     store,
		configMgr: configMgr,
	}

	defaultConfig := types.OpenAIAdapterConfig{
		ListenIP:   "127.0.0.1",
		ListenPort: 11223,
	}

	configStr, err := store.Get(AdapterConfigKey)
	if err != nil {
		log.Warn("Failed to load adapter config, using defaults.", "error", err)
		m.config = defaultConfig
	} else if configStr != "" {
		if err := json.Unmarshal([]byte(configStr), &m.config); err != nil {
			log.Error("Failed to unmarshal adapter config, using defaults.", "error", err)
			m.config = defaultConfig
		}
	} else {
		m.config = defaultConfig
	}

	return m
}

// SetContext sets the runtime context for the manager.
func (m *OpenAIAdapterManager) SetContext(ctx context.Context) {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.runtimeCtx = ctx
}

// emitStatusUpdate sends the current status to the frontend.
func (m *OpenAIAdapterManager) emitStatusUpdate() {
	if m.runtimeCtx == nil {
		return
	}
	// We send a copy to avoid any race conditions with the frontend
	statusCopy := m.status
	runtime.EventsEmit(m.runtimeCtx, "openai-adapter-status-changed", statusCopy)
}

// Start 启动服务
func (m *OpenAIAdapterManager) Start() error {
	// m.mu.Lock()
	// defer m.mu.Unlock()

	if m.status.IsRunning {
		return errors.New("adapter server is already running")
	}

	if m.runtimeCtx == nil {
		return errors.New("runtime context is not set, cannot start server")
	}

	// 为适配器服务创建专用的日志记录器
	adapterLogger := core.NewLogger(&core.LoggerOption{
		Level:  "debug",
		Prefix: "Adapter",
		Type:   "console", // 类型仅用于初始化, 我们将覆盖其输出
	})
	// 将日志输出重定向到我们的 Wails 事件写入器
	adapterLogger.SetOutput(&WailsLogger{ctx: m.runtimeCtx})

	if m.config.TargetOllamaServerID == "" {
		adapterLogger.Error("Target Ollama server is not configured")
		return errors.New("target Ollama server is not configured")
	}

	targetServer, err := m.configMgr.GetServerByID(m.config.TargetOllamaServerID)
	if err != nil {
		adapterLogger.Error(fmt.Sprintf("Failed to get target ollama server config: %v", err))
		return fmt.Errorf("failed to get target ollama server config: %w", err)
	}

	cleanBaseURL := strings.TrimSuffix(targetServer.BaseURL, "/")
	finalOllamaURL := cleanBaseURL + "/api/"

	// 使用新的日志记录器创建客户端和适配器
	ollamaClient := ai.NewOllamaClient(adapterLogger.WithPrefix("OllamaClient"), finalOllamaURL)
	adapter := ai.NewOpenAIAdapter(adapterLogger.WithPrefix("Handler"), ollamaClient, m.runtimeCtx)

	mux := http.NewServeMux()
	mux.Handle("/v1/chat/completions", adapter)

	addr := fmt.Sprintf("%s:%d", m.config.ListenIP, m.config.ListenPort)
	m.server = &http.Server{Addr: addr, Handler: mux}

	go func() {
		adapterLogger.Info(fmt.Sprintf("Starting OpenAI compatible API server on %s, targeting %s", addr, finalOllamaURL))
		if err := m.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// 使用新的 logger 记录错误, 它会自动推送到前端
			adapterLogger.Error(fmt.Sprintf("Server failed to start: %v", err))

			// m.mu.Lock()
			m.status.IsRunning = false
			m.status.Error = err.Error()
			m.emitStatusUpdate() // Send status update on failure
			// m.mu.Unlock()
		}
	}()

	m.status.IsRunning = true
	m.status.Error = ""
	m.emitStatusUpdate() // Send status update on success
	// m.log.Info("OpenAI adapter server has been started.")

	return nil
}

// Stop 停止服务
func (m *OpenAIAdapterManager) Stop(ctx context.Context) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()

	if !m.status.IsRunning || m.server == nil {
		return errors.New("adapter server is not running")
	}

	// 创建一个临时的日志记录器用于停止过程的日志输出
	adapterLogger := core.NewLogger(&core.LoggerOption{Level: "debug", Prefix: "Adapter", Type: "console"})
	adapterLogger.SetOutput(&WailsLogger{ctx: m.runtimeCtx})

	adapterLogger.Info("Stopping OpenAI adapter server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := m.server.Shutdown(shutdownCtx); err != nil {
		adapterLogger.Error(fmt.Sprintf("Server shutdown failed: %v", err))
		// Even if shutdown fails, we consider it stopped
	}

	m.server = nil
	m.status.IsRunning = false
	m.status.Error = ""
	m.emitStatusUpdate() // Send status update after stopping
	adapterLogger.Info("OpenAI adapter server has been stopped successfully.")

	return nil
}

// GetConfig 获取配置
func (m *OpenAIAdapterManager) GetConfig() types.OpenAIAdapterConfig {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return m.config
}

// SaveConfig 保存配置
func (m *OpenAIAdapterManager) SaveConfig(cfg types.OpenAIAdapterConfig) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()

	if m.status.IsRunning {
		return errors.New("cannot save config while server is running")
	}

	configBytes, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := m.store.Set(AdapterConfigKey, string(configBytes)); err != nil {
		return err
	}
	m.config = cfg
	return nil
}

// GetStatus 获取状态
func (m *OpenAIAdapterManager) GetStatus() types.OpenAIAdapterStatus {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return m.status
}
