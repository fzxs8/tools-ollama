package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/google/uuid"
)

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

// Server 代表一个配置好的Ollama服务器
type Server struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	BaseURL    string `json:"base_url"`
	APIKey     string `json:"api_key"`
	IsActive   bool   `json:"is_active"`
	TestStatus string `json:"test_status"`
	Type       string `json:"type"`
}

// PromptPilot 管理提示词工程功能
type PromptPilot struct {
	ctx       context.Context
	store     *duolasdk.AppStore
	configMgr *OllamaConfigManager
	logger    *core.AppLog
}

// NewPromptPilot 创建一个新的 PromptPilot 实例
func NewPromptPilot(store *duolasdk.AppStore, configMgr *OllamaConfigManager) *PromptPilot {
	logger := core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "提示词大师", // 使用中文前缀
	})

	return &PromptPilot{
		store:     store,
		configMgr: configMgr,
		logger:    logger,
	}
}

// Startup 在应用启动时调用
func (p *PromptPilot) Startup(ctx context.Context) {
	p.ctx = ctx
	p.logger.Info("提示词大师模块已启动")
}

// toCoreMessages 将 main.Message 转换为 core.Message
// 因为它们结构相同但在不同包中，所以需要转换
func (p *PromptPilot) toCoreMessages(messages []Message) []core.Message {
	coreMsgs := make([]core.Message, len(messages))
	for i, msg := range messages {
		coreMsgs[i] = core.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return coreMsgs
}

// GeneratePromptStream 流式生成一个提示词
func (p *PromptPilot) GeneratePromptStream(idea string, model string, serverId string) error {
	p.logger.Debug("开始生成提示词流, 想法: %s, 模型: %s, 服务器ID: %s", idea, model, serverId)

	if idea == "" {
		return fmt.Errorf("想法不能为空")
	}
	if model == "" {
		return fmt.Errorf("必须指定模型")
	}
	if serverId == "" {
		return fmt.Errorf("必须指定服务器ID")
	}

	// 1. 根据 serverId 获取服务器配置
	var serverConfig *OllamaServerConfig
	var err error

	//if serverId == "local" {
	//	localConfig, err := p.configMgr.GetLocalConfig()
	//	if err != nil {
	//		p.logger.Error("获取本地服务器配置失败: %v", err)
	//		return fmt.Errorf("获取本地服务器配置失败: %w", err)
	//	}
	//	serverConfig = &localConfig
	//} else {
	serverConfig, err = p.configMgr.GetServerByID(serverId)
	if err != nil {
		p.logger.Error("获取远程服务器配置失败 (ID: %s): %v", serverId, err)
		return fmt.Errorf("找不到ID为 %s 的服务器: %w", serverId, err)
	}
	//}

	p.logger.Debug("使用服务器: %s (%s)", serverConfig.Name, serverConfig.BaseURL)

	// 2. 创建 Ollama 提供者
	provider := core.NewOllamaProvider(serverConfig.BaseURL)

	// 3. 构建发送给 AI 的消息
	messages := []Message{
		{
			Role:    "system",
			Content: "你是一位专业的提示词工程大师。你的任务是根据用户提供的想法，创作出一个清晰、具体、高效的提示词，以便用于大型语言模型。",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("这是我的想法：%s", idea),
		},
	}
	p.logger.Debug("构建的消息体: %+v", messages)

	// 4. 定义流式回调函数
	callback := func(chunk string) {
		p.logger.Debug("收到流式数据块: %s (模型: %s)", chunk, model)
		runtime.EventsEmit(p.ctx, "prompt_pilot_stream", map[string]string{"model": model, "chunk": chunk})
	}

	// 5. 调用核心库的 ChatStream 方法
	coreMessages := p.toCoreMessages(messages)
	err = provider.ChatStream(model, coreMessages, callback)
	if err != nil {
		p.logger.Error("流式生成提示词失败: %v", err)
		runtime.EventsEmit(p.ctx, "prompt_pilot_stream_error", map[string]string{"model": model, "error": err.Error()})
		return fmt.Errorf("调用Ollama服务失败: %w", err)
	}

	runtime.EventsEmit(p.ctx, "prompt_pilot_stream_done", map[string]string{"model": model})
	p.logger.Debug("提示词流式生成完成")
	return nil
}

// OptimizePrompt 优化一个提示词
func (p *PromptPilot) OptimizePrompt(content string, feedback string, model string, serverId string) (string, error) {
	p.logger.Debug("开始优化提示词, 模型: %s, 服务器: %s", model, serverId)

	if content == "" {
		return "", fmt.Errorf("提示词内容不能为空")
	}
	if feedback == "" {
		return "", fmt.Errorf("优化反馈不能为空")
	}

	// 模拟优化
	optimizedPrompt := fmt.Sprintf("%s\n\n--- 优化建议 ---\n%s\n\n--- 基于建议优化后的版本 ---", content, feedback)
	p.logger.Debug("提示词优化完成")

	return optimizedPrompt, nil
}

// SavePrompt 保存一个提示词到存储
func (p *PromptPilot) SavePrompt(prompt Prompt) error {
	p.logger.Debug("准备保存提示词: %s", prompt.Name)

	// 如果是新提示词, 生成ID和时间戳
	if prompt.ID == "" {
		prompt.ID = uuid.New().String()
		now := time.Now().UnixNano() / int64(time.Millisecond)
		prompt.CreatedAt = now
		prompt.UpdatedAt = now
		prompt.Version = 1
		p.logger.Debug("新提示词, 已生成ID: %s", prompt.ID)
	} else {
		// 更新已有提示词的时间戳和版本
		prompt.UpdatedAt = time.Now().UnixNano() / int64(time.Millisecond)
		prompt.Version++
		p.logger.Debug("更新已有提示词, ID: %s, 版本: %d", prompt.ID, prompt.Version)
	}

	// 序列化为JSON
	promptJSON, err := json.Marshal(prompt)
	if err != nil {
		p.logger.Error("序列化提示词失败: %v", err)
		return err
	}

	// 使用哈希存储保存
	err = p.store.HSet("prompts", prompt.ID, string(promptJSON))
	if err != nil {
		p.logger.Error("保存提示词到数据库失败: %v", err)
		return err
	}

	p.logger.Info("提示词保存成功: %s (ID: %s)", prompt.Name, prompt.ID)
	return nil
}

// ListPrompts 返回所有已保存的提示词
func (p *PromptPilot) ListPrompts() ([]Prompt, error) {
	p.logger.Debug("开始列出所有已保存的提示词")

	// 从存储中获取所有提示词
	promptsMap, err := p.store.HGetAll("prompts")
	if err != nil {
		p.logger.Error("从数据库获取提示词列表失败: %v", err)
		return nil, err
	}

	// 解析提示词
	prompts := make([]Prompt, 0, len(promptsMap))
	for _, promptJSON := range promptsMap {
		var prompt Prompt
		if err := json.Unmarshal([]byte(promptJSON), &prompt); err != nil {
			p.logger.Warn("解析单个提示词失败: %v, 内容: %s", err, promptJSON)
			continue
		}
		prompts = append(prompts, prompt)
	}

	p.logger.Debug("成功获取 %d 个提示词", len(prompts))
	return prompts, nil
}

// GetPrompt 根据ID返回一个特定的提示词
func (p *PromptPilot) GetPrompt(id string) (Prompt, error) {
	p.logger.Debug("开始获取提示词, ID: %s", id)

	// 从存储中获取
	promptJSON, err := p.store.HGet("prompts", id)
	if err != nil {
		p.logger.Error("从数据库获取提示词失败: %v", err)
		return Prompt{}, err
	}

	// 解析
	var prompt Prompt
	if err := json.Unmarshal([]byte(promptJSON), &prompt); err != nil {
		p.logger.Error("解析提示词JSON失败: %v", err)
		return Prompt{}, err
	}

	p.logger.Debug("成功获取提示词: %s", prompt.Name)
	return prompt, nil
}

// DeletePrompt 根据ID删除一个提示词
func (p *PromptPilot) DeletePrompt(id string) error {
	p.logger.Debug("开始删除提示词, ID: %s", id)

	// 从存储中删除
	err := p.store.HDel("prompts", id)
	if err != nil {
		p.logger.Error("从数据库删除提示词失败: %v", err)
		return err
	}

	p.logger.Info("提示词删除成功, ID: %s", id)
	return nil
}
