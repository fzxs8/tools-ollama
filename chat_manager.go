package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// ChatSystemPrompt 系统提示词结构
type ChatSystemPrompt struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Prompt    string `json:"prompt"`
	CreatedAt int64  `json:"createdAt"`
}

// ChatManager 聊天管理器
type ChatManager struct {
	ctx   context.Context
	store interface {
		Set(key, value string, persistent ...bool) error
		Get(key string, persistent ...bool) (string, error)
		Delete(key string, persistent ...bool) error
	}
}

// NewChatManager 创建聊天管理器实例
func NewChatManager(ctx context.Context, store interface {
	Set(key, value string, persistent ...bool) error
	Get(key string, persistent ...bool) (string, error)
	Delete(key string, persistent ...bool) error
}) *ChatManager {
	return &ChatManager{
		ctx:   ctx,
		store: store,
	}
}

// SetContext 设置上下文
func (cm *ChatManager) SetContext(ctx context.Context) {
	cm.ctx = ctx
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
	value, err := cm.store.Get(key, true) // 使用持久化存储
	if err != nil {
		return "", err
	}
	return value, nil
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

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
