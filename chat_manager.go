package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/google/uuid"
)

// Message 聊天消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatSystemPrompt 系统提示词结构
type ChatSystemPrompt struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Prompt    string `json:"prompt"`
	CreatedAt int64  `json:"createdAt"`
}

// ChatManager 聊天管理器
type ChatManager struct {
	ctx        context.Context
	store      *duolasdk.AppStore
	aiProvider interface {
		Chat(model string, messages []Message) (string, error)
		ChatStream(model string, messages []Message, callback func(string)) error
	}
	logger *core.AppLog
}

// NewChatManager 创建聊天管理器实例
func NewChatManager(ctx context.Context, store *duolasdk.AppStore) *ChatManager {
	logger := core.NewLogger(&core.LoggerOption{Type: "console", Level: "debug", Prefix: "ChatManager"})
	return &ChatManager{
		ctx:    ctx,
		store:  store,
		logger: logger,
	}
}

// SetContext 设置上下文
func (cm *ChatManager) SetContext(ctx context.Context) {
	cm.ctx = ctx
}

// SetAIProvider 设置AI提供者
func (cm *ChatManager) SetAIProvider(provider interface {
	Chat(model string, messages []Message) (string, error)
	ChatStream(model string, messages []Message, callback func(string)) error
}) {
	cm.aiProvider = provider
}

// KVSet 设置键值对
func (cm *ChatManager) KVSet(key, value string) error {
	return cm.store.Set(key, value, true) // 使用持久化存储
}

// KVGet 获取键值对
func (cm *ChatManager) KVGet(key string) (string, error) {
	value, err := cm.store.Get(key, true) // 使用持久化存储
	if err != nil {
		return "", err
	}
	return value, nil
}

// KVList 获取键值对列表
func (cm *ChatManager) KVList(key string) (string, error) {
	// 这里需要实现获取键值对列表的逻辑
	// 由于当前store接口没有提供List方法，暂时返回空字符串
	// 注意：这个方法名与功能不匹配，应该重构
	return "[]", nil
}

// KVDelete 删除键值对
func (cm *ChatManager) KVDelete(key string) error {
	return cm.store.Delete(key, true) // 使用持久化存储
}

// SaveSystemPrompt 保存系统提示词
func (cm *ChatManager) SaveSystemPrompt(title, prompt string) error {
	// 获取现有的提示词列表
	promptListStr, err := cm.KVGet("system_prompts")
	if err != nil {
		return fmt.Errorf("获取系统提示词列表失败: %w", err)
	}

	var promptList []ChatSystemPrompt
	if promptListStr != "" {
		if err := json.Unmarshal([]byte(promptListStr), &promptList); err != nil {
			return fmt.Errorf("解析系统提示词列表失败: %w", err)
		}
	}

	// 创建新的提示词
	newPrompt := ChatSystemPrompt{
		ID:        fmt.Sprintf("%d", len(promptList)+1),
		Title:     title,
		Prompt:    prompt,
		CreatedAt: getCurrentTimestamp(),
	}

	// 添加到列表
	promptList = append(promptList, newPrompt)

	// 保存到存储
	promptListBytes, err := json.Marshal(promptList)
	if err != nil {
		return fmt.Errorf("序列化系统提示词列表失败: %w", err)
	}

	return cm.KVSet("system_prompts", string(promptListBytes))
}

// GetSystemPrompts 获取所有系统提示词
func (cm *ChatManager) GetSystemPrompts() ([]ChatSystemPrompt, error) {
	promptListStr, err := cm.KVGet("system_prompts")
	if err != nil {
		return nil, fmt.Errorf("获取系统提示词列表失败: %w", err)
	}

	if promptListStr == "" {
		return []ChatSystemPrompt{}, nil
	}

	var promptList []ChatSystemPrompt
	if err := json.Unmarshal([]byte(promptListStr), &promptList); err != nil {
		return nil, fmt.Errorf("解析系统提示词列表失败: %w", err)
	}

	return promptList, nil
}

// UpdateSystemPrompt 更新系统提示词
func (cm *ChatManager) UpdateSystemPrompt(id, title, prompt string) error {
	// 获取现有的提示词列表
	promptListStr, err := cm.KVGet("system_prompts")
	if err != nil {
		return fmt.Errorf("获取系统提示词列表失败: %w", err)
	}

	if promptListStr == "" {
		return fmt.Errorf("系统提示词列表为空")
	}

	var promptList []ChatSystemPrompt
	if err := json.Unmarshal([]byte(promptListStr), &promptList); err != nil {
		return fmt.Errorf("解析系统提示词列表失败: %w", err)
	}

	// 查找并更新提示词
	found := false
	for i, p := range promptList {
		if p.ID == id {
			promptList[i].Title = title
			promptList[i].Prompt = prompt
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("未找到ID为 %s 的系统提示词", id)
	}

	// 保存到存储
	promptListBytes, err := json.Marshal(promptList)
	if err != nil {
		return fmt.Errorf("序列化系统提示词列表失败: %w", err)
	}

	return cm.KVSet("system_prompts", string(promptListBytes))
}

// DeleteSystemPrompt 删除系统提示词
func (cm *ChatManager) DeleteSystemPrompt(id string) error {
	// 获取现有的提示词列表
	promptListStr, err := cm.KVGet("system_prompts")
	if err != nil {
		return fmt.Errorf("获取系统提示词列表失败: %w", err)
	}

	if promptListStr == "" {
		return nil
	}

	var promptList []ChatSystemPrompt
	if err := json.Unmarshal([]byte(promptListStr), &promptList); err != nil {
		return fmt.Errorf("解析系统提示词列表失败: %w", err)
	}

	// 查找并删除提示词
	newPromptList := []ChatSystemPrompt{}
	for _, p := range promptList {
		if p.ID != id {
			newPromptList = append(newPromptList, p)
		}
	}

	// 保存到存储
	promptListBytes, err := json.Marshal(newPromptList)
	if err != nil {
		return fmt.Errorf("序列化系统提示词列表失败: %w", err)
	}

	return cm.KVSet("system_prompts", string(promptListBytes))
}

// SetActiveSystemPrompt 设置激活的系统提示词
func (cm *ChatManager) SetActiveSystemPrompt(prompt ChatSystemPrompt) error {
	promptBytes, err := json.Marshal(prompt)
	if err != nil {
		return fmt.Errorf("序列化激活的系统提示词失败: %w", err)
	}

	return cm.KVSet("active_system_prompt", string(promptBytes))
}

// GetActiveSystemPrompt 获取激活的系统提示词
func (cm *ChatManager) GetActiveSystemPrompt() (*ChatSystemPrompt, error) {
	promptStr, err := cm.KVGet("active_system_prompt")
	if err != nil {
		return nil, fmt.Errorf("获取激活的系统提示词失败: %w", err)
	}

	if promptStr == "" {
		return nil, nil
	}

	var prompt ChatSystemPrompt
	if err := json.Unmarshal([]byte(promptStr), &prompt); err != nil {
		return nil, fmt.Errorf("解析激活的系统提示词失败: %w", err)
	}

	return &prompt, nil
}

// DeleteActiveSystemPrompt 删除激活的系统提示词
func (cm *ChatManager) DeleteActiveSystemPrompt() error {
	return cm.KVDelete("active_system_prompt")
}

// ListConversations 获取所有已保存的对话列表，按时间倒序排列
func (cm *ChatManager) ListConversations() ([]*Conversation, error) {
	// 使用HGetAll获取所有对话
	conversationsMap, err := cm.store.HGetAll("conversations")

	if err != nil {
		cm.logger.Error("获取对话列表失败: %v", err)
		return nil, fmt.Errorf("获取对话列表失败: %w", err)
	}

	// 创建对话列表
	conversations := make([]*Conversation, 0, len(conversationsMap))

	// 解析每个对话
	for _, convJSON := range conversationsMap {
		var conv Conversation
		if err := json.Unmarshal([]byte(convJSON), &conv); err != nil {
			cm.logger.Warn("解析对话失败: %v", err)
			continue // 跳过无效的对话
		}
		conversations = append(conversations, &conv)
	}

	// 按时间倒序排列（最新的在前面）
	for i := 0; i < len(conversations)/2; i++ {
		conversations[i], conversations[len(conversations)-1-i] = conversations[len(conversations)-1-i], conversations[i]
	}

	cm.logger.Debug("成功获取对话列表，共 %d 个对话", len(conversations))
	return conversations, nil
}

// SaveConversation 创建或更新一个对话
func (cm *ChatManager) SaveConversation(conv *Conversation) (*Conversation, error) {
	// 如果ID为空，则为新创建
	if conv.ID == "" {
		conv.ID = uuid.New().String()
		conv.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}

	// 序列化对话对象
	convJSON, err := json.Marshal(conv)
	if err != nil {
		cm.logger.Error("序列化对话失败: %v", err)
		return nil, fmt.Errorf("序列化对话失败: %w", err)
	}

	// 保存到存储中
	err = cm.store.HSet("conversations", conv.ID, string(convJSON))

	if err != nil {
		cm.logger.Error("保存对话失败: %v", err)
		return nil, fmt.Errorf("保存对话失败: %w", err)
	}

	cm.logger.Debug("对话保存成功，ID: %s", conv.ID)
	return conv, nil
}

// GetConversation 获取指定ID的单个对话的完整内容
func (cm *ChatManager) GetConversation(id string) (*Conversation, error) {
	// 从存储中获取对话
	convJSON, err := cm.store.HGet("conversations", id)

	if err != nil {
		cm.logger.Error("获取对话失败: %v", err)
		return nil, fmt.Errorf("获取对话失败: %w", err)
	}

	// 解析对话对象
	var conv Conversation
	if err := json.Unmarshal([]byte(convJSON), &conv); err != nil {
		cm.logger.Error("解析对话失败: %v", err)
		return nil, fmt.Errorf("解析对话失败: %w", err)
	}

	cm.logger.Debug("成功获取对话，ID: %s", id)
	return &conv, nil
}

// DeleteConversation 删除指定ID的对话
func (cm *ChatManager) DeleteConversation(id string) error {
	// 从存储中删除对话
	err := cm.store.HDel("conversations", id)

	if err != nil {
		cm.logger.Error("删除对话失败: %v", err)
		return fmt.Errorf("删除对话失败: %w", err)
	}

	cm.logger.Debug("对话删除成功，ID: %s", id)
	return nil
}

// Chat 发送聊天消息
func (cm *ChatManager) Chat(modelName string, messages []Message) (string, error) {
	if cm.aiProvider == nil {
		err := fmt.Errorf("AI提供者未设置")
		cm.logger.Error(err.Error())
		return "", err
	}

	cm.logger.Debug("开始阻塞式聊天: 模型=%s, 消息数量=%d", modelName, len(messages))

	defer func() {
		if r := recover(); r != nil {
			cm.logger.Error("Chat方法中发生恐慌: %v", r)
		}
	}()

	result, err := cm.aiProvider.Chat(modelName, messages)
	if err != nil {
		cm.logger.Error("聊天过程中发生错误: %v", err)
		return "", err
	}

	cm.logger.Debug("阻塞式聊天完成，结果长度=%d", len(result))
	cm.logger.Debug("阻塞式聊天返回结果前100个字符: %s", func() string {
		if len(result) > 100 {
			// 使用 rune 转换来正确处理 Unicode 字符
			runes := []rune(result)
			if len(runes) > 100 {
				return string(runes[:100]) + "..."
			}
			return result
		}
		return result
	}())
	return result, nil
}

// ChatStream 发送聊天消息并流式返回结果
func (cm *ChatManager) ChatStream(modelName string, messages []Message, callback func(string)) error {
	if cm.aiProvider == nil {
		err := fmt.Errorf("AI提供者未设置")
		cm.logger.Error(err.Error())
		return err
	}

	cm.logger.Debug("开始流式聊天: 模型=%s, 消息数量=%d", modelName, len(messages))

	defer func() {
		if r := recover(); r != nil {
			cm.logger.Error("ChatStream方法中发生恐慌: %v", r)
		}
	}()

	// 创建一个安全的回调函数包装器
	safeCallback := func(content string) {
		// recover任何可能的恐慌
		defer func() {
			if r := recover(); r != nil {
				cm.logger.Error("流式回调函数中发生恐慌: %v", r)
			}
		}()

		cm.logger.Debug("ChatStream回调接收到内容，长度=%d", len(content))
		// 调用原始回调函数
		if callback != nil {
			callback(content)
		}
	}

	return cm.aiProvider.ChatStream(modelName, messages, safeCallback)
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// safeTruncate 安全地截取字符串，避免在多字节字符中间截断
func safeTruncate(s string, length int) string {
	if len(s) <= length {
		return s
	}

	runes := []rune(s)
	if len(runes) <= length {
		return s
	}

	return string(runes[:length]) + "..."
}
