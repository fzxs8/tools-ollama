package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings" // 导入 strings 包
	"time"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/fzxs8/duolasdk/core/ai" // 导入 ai 包
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// PromptEngineering 管理提示词工程功能
type PromptEngineering struct {
	ctx       context.Context
	store     *duolasdk.AppStore
	configMgr *OllamaConfigManager
	logger    *core.AppLog
}

// NewPromptPilot 创建一个新的 PromptEngineering 实例
func NewPromptPilot(store *duolasdk.AppStore, configMgr *OllamaConfigManager, logger *core.AppLog) *PromptEngineering {
	return &PromptEngineering{
		store:     store,
		configMgr: configMgr,
		logger:    logger.WithPrefix("PromptEngineering"),
	}
}

// Startup 在应用启动时调用
func (p *PromptEngineering) Startup(ctx context.Context) {
	p.ctx = ctx
	p.logger.Info("提示词大师模块已启动")
}

// GeneratePromptStream 流式生成一个提示词
func (p *PromptEngineering) GeneratePromptStream(idea string, model string, serverId string) {
	p.logger.Debug("开始生成提示词流", "idea", idea, "model", model, "serverId", serverId)

	go func() {
		if idea == "" || model == "" || serverId == "" {
			p.logger.Error("GeneratePromptStream 参数无效", "idea", idea, "model", model, "serverId", serverId)
			runtime.EventsEmit(p.ctx, "prompt_pilot_stream_error", map[string]string{"model": model, "error": "参数无效"})
			return
		}

		serverConfig, err := p.configMgr.GetServerByID(serverId)
		if err != nil {
			p.logger.Error("获取服务器配置失败", "serverId", serverId, "error", err)
			runtime.EventsEmit(p.ctx, "prompt_pilot_stream_error", map[string]string{"model": model, "error": err.Error()})
			return
		}

		p.logger.Debug("使用服务器", "serverName", serverConfig.Name, "baseURL", serverConfig.BaseURL)

		cleanBaseURL := strings.TrimSuffix(serverConfig.BaseURL, "/")
		finalOllamaURL := cleanBaseURL + "/api"
		p.logger.Debug("prompt_pilot.go: Final Ollama Base URL being passed to NewOllamaProvider", "url", finalOllamaURL) // Added debug log
		provider := ai.NewOllamaProvider(p.logger, finalOllamaURL)
		messages := []types.Message{
			{
				Role:    "system",
				Content: "你是一位专业的提示词工程大师。你的任务是根据用户提供的想法，创作出一个清晰、具体、高效的提示词，以便用于大型语言模型。",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("这是我的想法：%s", idea),
			},
		}

		callback := func(chunk string) {
			runtime.EventsEmit(p.ctx, "prompt_pilot_stream", map[string]string{"model": model, "chunk": chunk})
		}

		coreMessages := ToCoreMessages(messages)
		if err := provider.ChatStream(model, coreMessages, callback); err != nil {
			p.logger.Error("流式生成提示词失败", "model", model, "error", err)
			runtime.EventsEmit(p.ctx, "prompt_pilot_stream_error", map[string]string{"model": model, "error": err.Error()})
			return
		}

		runtime.EventsEmit(p.ctx, "prompt_pilot_stream_done", map[string]string{"model": model})
		p.logger.Debug("提示词流式生成完成", "model", model)
	}()
}

// OptimizePrompt 优化一个提示词
func (p *PromptEngineering) OptimizePrompt(content string, feedback string, model string, serverId string) (string, error) {
	p.logger.Debug("开始优化提示词", "model", model, "serverId", serverId)

	if content == "" || feedback == "" || model == "" || serverId == "" {
		return "", fmt.Errorf("参数不能为空")
	}

	// 此处为简化实现，实际可调用AI进行优化
	optimizedPrompt := fmt.Sprintf("%s\n\n--- 优化建议 ---\n%s\n\n--- 基于建议优化后的版本 ---", content, feedback)
	p.logger.Debug("提示词优化完成")

	return optimizedPrompt, nil
}

// SavePrompt 保存一个提示词到存储
func (p *PromptEngineering) SavePrompt(prompt types.Prompt) error {
	p.logger.Debug("准备保存提示词", "promptName", prompt.Name)

	if prompt.ID == "" {
		prompt.ID = uuid.New().String()
		now := time.Now().UnixNano() / int64(time.Millisecond)
		prompt.CreatedAt = now
		prompt.UpdatedAt = now
		prompt.Version = 1
		p.logger.Debug("创建新提示词", "id", prompt.ID)
	} else {
		prompt.UpdatedAt = time.Now().UnixNano() / int64(time.Millisecond)
		prompt.Version++
		p.logger.Debug("更新提示词", "id", prompt.ID, "newVersion", prompt.Version)
	}

	promptJSON, err := json.Marshal(prompt)
	if err != nil {
		p.logger.Error("序列化提示词失败", "error", err)
		return err
	}

	if err := p.store.HSet("prompts", prompt.ID, string(promptJSON)); err != nil {
		p.logger.Error("保存提示词到数据库失败", "error", err)
		return err
	}

	p.logger.Info("提示词保存成功", "promptName", prompt.Name, "id", prompt.ID)
	return nil
}

// ListPrompts 返回所有已保存的提示词
func (p *PromptEngineering) ListPrompts() ([]types.Prompt, error) {
	p.logger.Debug("开始列出所有已保存的提示词")

	promptsMap, err := p.store.HGetAll("prompts")
	if err != nil {
		p.logger.Error("从数据库获取提示词列表失败", "error", err)
		return nil, err
	}

	prompts := make([]types.Prompt, 0, len(promptsMap))
	for _, promptJSON := range promptsMap {
		var prompt types.Prompt
		if err := json.Unmarshal([]byte(promptJSON), &prompt); err != nil {
			p.logger.Warn("解析单个提示词失败", "error", err, "jsonData", promptJSON)
			continue
		}
		prompts = append(prompts, prompt)
	}

	p.logger.Debug("成功获取 %d 个提示词", len(prompts))
	return prompts, nil
}

// GetPrompt 根据ID返回一个特定的提示词
func (p *PromptEngineering) GetPrompt(id string) (types.Prompt, error) {
	p.logger.Debug("开始获取提示词", "id", id)

	promptJSON, err := p.store.HGet("prompts", id)
	if err != nil {
		p.logger.Error("从数据库获取提示词失败", "id", id, "error", err)
		return types.Prompt{}, err
	}

	var prompt types.Prompt
	if err := json.Unmarshal([]byte(promptJSON), &prompt); err != nil {
		p.logger.Error("解析提示词JSON失败", "id", id, "error", err)
		return types.Prompt{}, err
	}

	p.logger.Debug("成功获取提示词", "promptName", prompt.Name)
	return prompt, nil
}

// DeletePrompt 根据ID删除一个提示词
func (p *PromptEngineering) DeletePrompt(id string) error {
	p.logger.Debug("开始删除提示词", "id", id)

	if err := p.store.HDel("prompts", id); err != nil {
		p.logger.Error("从数据库删除提示词失败", "id", id, "error", err)
		return err
	}

	p.logger.Info("提示词删除成功", "id", id)
	return nil
}
