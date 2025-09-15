package main

import (
	"context"
	"fmt"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ChatManager 聊天管理器
type ChatManager struct {
	ctx        context.Context
	store      *duolasdk.AppStore
	aiProvider AIProvider
	logger     *core.AppLog
}

// AIProvider 定义了AI聊天能力的接口
type AIProvider interface {
	Chat(model string, messages []core.Message) (string, error)
	ChatStream(model string, messages []core.Message, callback func(string)) error
}

// NewChatManager 创建聊天管理器实例
func NewChatManager(ctx context.Context, store *duolasdk.AppStore, logger *core.AppLog) *ChatManager {
	return &ChatManager{
		ctx:    ctx,
		store:  store,
		logger: logger.WithPrefix("ChatManager"),
	}
}

// SetContext 设置上下文
func (cm *ChatManager) SetContext(ctx context.Context) {
	cm.ctx = ctx
}

// SetAIProvider 设置AI提供者
func (cm *ChatManager) SetAIProvider(provider AIProvider) {
	cm.logger.Debug("设置新的AI Provider")
	cm.aiProvider = provider
}

// ListConversations 获取所有已保存的对话列表，按时间倒序排列
func (cm *ChatManager) ListConversations() ([]*types.Conversation, error) {
	cm.logger.Debug("获取所有对话列表")
	conversationsMap, err := cm.store.HGetAll("conversations")
	if err != nil {
		cm.logger.Error("获取对话列表失败", "error", err)
		return nil, fmt.Errorf("获取对话列表失败: %w", err)
	}

	conversations := make([]*types.Conversation, 0, len(conversationsMap))
	for _, convJSON := range conversationsMap {
		var conv types.Conversation
		if err := UnmarshalJSONWithError([]byte(convJSON), &conv, cm.logger, "解析对话"); err != nil {
			continue
		}
		conversations = append(conversations, &conv)
	}

	// 按时间倒序排列（最新的在前面）
	ReverseSlice(conversations)

	cm.logger.Debug("成功获取对话列表", "count", len(conversations))
	return conversations, nil
}

// SaveConversation 创建或更新一个对话
func (cm *ChatManager) SaveConversation(conv *types.Conversation) (*types.Conversation, error) {
	if conv.ID == "" {
		conv.ID = GenerateUniqueID()
		conv.Timestamp = GetCurrentTimestamp()
		cm.logger.Debug("创建新对话", "id", conv.ID)
	} else {
		cm.logger.Debug("更新对话", "id", conv.ID)
	}

	convJSON, err := MarshalJSONWithError(conv, cm.logger, "序列化对话")
	if err != nil {
		return nil, err
	}

	if err := cm.store.HSet("conversations", conv.ID, string(convJSON)); err != nil {
		cm.logger.Error("保存对话到存储失败", "id", conv.ID, "error", err)
		return nil, fmt.Errorf("保存对话失败: %w", err)
	}

	cm.logger.Info("对话保存成功", "id", conv.ID)
	return conv, nil
}

// GetConversation 获取指定ID的单个对话的完整内容
func (cm *ChatManager) GetConversation(id string) (*types.Conversation, error) {
	cm.logger.Debug("获取对话", "id", id)
	convJSON, err := cm.store.HGet("conversations", id)
	if err != nil {
		cm.logger.Error("从存储获取对话失败", "id", id, "error", err)
		return nil, fmt.Errorf("获取对话失败: %w", err)
	}

	var conv types.Conversation
	if err := UnmarshalJSONWithError([]byte(convJSON), &conv, cm.logger, "解析对话"); err != nil {
		return nil, err
	}

	cm.logger.Debug("成功获取对话", "id", id)
	return &conv, nil
}

// DeleteConversation 删除指定ID的对话
func (cm *ChatManager) DeleteConversation(id string) error {
	cm.logger.Info("删除对话", "id", id)
	if err := cm.store.HDel("conversations", id); err != nil {
		cm.logger.Error("从存储删除对话失败", "id", id, "error", err)
		return fmt.Errorf("删除对话失败: %w", err)
	}
	return nil
}

// ChatMessage 发送聊天消息到Ollama API
func (cm *ChatManager) ChatMessage(modelName string, messages []types.Message, stream bool) (string, error) {
	cm.logger.Debug("收到聊天消息请求", "model", modelName, "messageCount", len(messages), "stream", stream)

	if cm.aiProvider == nil {
		cm.logger.Error("AI provider未设置")
		return "", fmt.Errorf("AI provider not set")
	}

	coreMessages := ToCoreMessages(messages)

	if stream {
		cm.logger.Debug("使用流式传输")
		go func() {
			defer func() {
				if r := recover(); r != nil {
					cm.logger.Error("流式聊天goroutine发生恐慌", "panic", r)
					runtime.EventsEmit(cm.ctx, "chat_stream_error", fmt.Sprintf("内部错误: %v", r))
				}
				runtime.EventsEmit(cm.ctx, "chat_stream_done")
			}()
			err := cm.aiProvider.ChatStream(modelName, coreMessages, func(content string) {
				runtime.EventsEmit(cm.ctx, "chat_stream_chunk", content)
			})
			if err != nil {
				cm.logger.Error("流式聊天失败", "error", err)
				runtime.EventsEmit(cm.ctx, "chat_stream_error", err.Error())
			}
		}()
		return "", nil // For stream, the main function returns immediately
	} else {
		cm.logger.Debug("使用阻塞式传输")
		result, err := cm.aiProvider.Chat(modelName, coreMessages)
		if err != nil {
			cm.logger.Error("阻塞式聊天失败", "error", err)
			return "", err
		}
		cm.logger.Debug("阻塞式传输成功", "resultLength", len(result))
		return result, nil
	}
}
