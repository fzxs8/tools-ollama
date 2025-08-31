package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fzxs8/duolasdk/core"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ModelManager 模型管理器
type ModelManager struct {
	ctx context.Context
	app *App
}

// NewModelManager 创建新的模型管理器
func NewModelManager(app *App) *ModelManager {
	return &ModelManager{
		app: app,
	}
}

// SetContext 设置上下文
func (m *ModelManager) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// ModelParams 模型参数
type ModelParams struct {
	Temperature   float64 `json:"temperature"`
	TopP          float64 `json:"top_p"`
	TopK          int     `json:"top_k"`
	Context       int     `json:"context"`
	NumPredict    int     `json:"num_predict"`
	RepeatPenalty float64 `json:"repeat_penalty"`
}

// ModelSearchParams 模型搜索参数
type ModelSearchParams struct {
	Query    string   `json:"query"`
	Families []string `json:"families"`
	Tags     []string `json:"tags"`
}

// RunningModel 运行中的模型
type RunningModel struct {
	Name      string      `json:"name"`
	Params    ModelParams `json:"params"`
	StartTime time.Time   `json:"start_time"`
	IsActive  bool        `json:"is_active"`
}

// 模型运行状态管理
var runningModels = make(map[string]*RunningModel)

// ListModels 获取模型列表
func (m *ModelManager) ListModels() ([]Model, error) {
	return m.app.ListModels()
}

// RunModel 运行模型
func (m *ModelManager) RunModel(modelName string, params ModelParams) error {
	// 检查模型是否已经在运行
	if _, exists := runningModels[modelName]; exists {
		return fmt.Errorf("模型 %s 已经在运行", modelName)
	}

	// 发送请求到Ollama生成端点来预热模型
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "hello",
		"stream": false,
	}

	response, err := m.app.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})

	if err != nil {
		return fmt.Errorf("启动模型失败: %v", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("启动模型失败，状态码: %d", response.StatusCode)
	}

	// 安全地解析响应
	var result map[string]interface{}
	if response.Body != "" {
		if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
			// 即使解析失败，我们也认为模型启动成功，因为HTTP请求是成功的
			fmt.Printf("警告: 解析Ollama响应失败: %v\n", err)
		}
	}

	// 创建运行中的模型记录
	runningModels[modelName] = &RunningModel{
		Name:      modelName,
		Params:    params,
		StartTime: time.Now(),
		IsActive:  true,
	}

	// 发送通知
	runtime.EventsEmit(m.ctx, "model:started", map[string]interface{}{
		"name": modelName,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})

	return nil
}

// StopModel 停止模型
func (m *ModelManager) StopModel(modelName string) error {
	// 检查模型是否在运行
	if _, exists := runningModels[modelName]; !exists {
		return fmt.Errorf("模型 %s 未在运行", modelName)
	}

	// 发送请求到Ollama生成端点来卸载模型
	requestBody := map[string]interface{}{
		"model":   modelName,
		"prompt":  "",
		"stream":  false,
		"options": map[string]interface{}{"num_predict": 1},
	}

	// 使用defer recover()防止panic导致应用崩溃
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("恢复 StopModel panic: %v\n", r)
			// 即使出现panic，也从运行列表中移除模型
			delete(runningModels, modelName)
			// 发送停止通知
			runtime.EventsEmit(m.ctx, "model:stopped", map[string]interface{}{
				"name": modelName,
				"time": time.Now().Format("2006-01-02 15:04:05"),
			})
		}
	}()

	_, err := m.app.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})

	if err != nil {
		// 即使API调用失败，我们仍然从运行列表中移除模型
		delete(runningModels, modelName)
		// 发送通知
		runtime.EventsEmit(m.ctx, "model:stopped", map[string]interface{}{
			"name": modelName,
			"time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return fmt.Errorf("停止模型时出现警告: %v", err)
	}

	// 删除运行记录
	delete(runningModels, modelName)

	// 发送通知
	runtime.EventsEmit(m.ctx, "model:stopped", map[string]interface{}{
		"name": modelName,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})

	return nil
}

// GetModelStatus 获取模型状态
func (m *ModelManager) GetModelStatus(modelName string) (*RunningModel, error) {
	if model, exists := runningModels[modelName]; exists {
		return model, nil
	}
	return nil, fmt.Errorf("模型 %s 未在运行", modelName)
}

// ListRunningModels 获取运行中的模型列表
func (m *ModelManager) ListRunningModels() []RunningModel {
	var models []RunningModel
	for _, model := range runningModels {
		models = append(models, *model)
	}
	return models
}

// TestModel 测试模型
func (m *ModelManager) TestModel(modelName string) (string, error) {
	// 发送测试请求到Ollama
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "你好，请简单介绍一下自己",
		"stream": false,
	}

	response, err := m.app.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})

	if err != nil {
		return "", fmt.Errorf("测试模型失败: %v", err)
	}

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("测试模型失败，状态码: %d", response.StatusCode)
	}

	// 安全地解析响应
	var result map[string]interface{}
	if response.Body != "" {
		if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
			return "", fmt.Errorf("解析响应失败: %v", err)
		}
	} else {
		return "测试完成，但未收到有效响应", nil
	}

	if responseText, ok := result["response"].(string); ok {
		return responseText, nil
	}

	// 如果没有response字段，返回完整响应供调试
	responseStr, _ := json.Marshal(result)
	return fmt.Sprintf("测试完成，响应内容: %s", string(responseStr)), nil
}

// DeleteModel 删除模型
func (m *ModelManager) DeleteModel(modelName string) error {
	return m.app.DeleteModel(modelName)
}

// SetModelParams 设置模型参数
func (m *ModelManager) SetModelParams(modelName string, params ModelParams) error {
	// 如果模型正在运行，更新运行参数
	if model, exists := runningModels[modelName]; exists {
		model.Params = params
	}

	// 这里可以将参数保存到本地存储或配置文件中
	// 暂时只在内存中保存
	return nil
}

// GetModelParams 获取模型参数
func (m *ModelManager) GetModelParams(modelName string) (ModelParams, error) {
	// 如果模型正在运行，返回运行参数
	if model, exists := runningModels[modelName]; exists {
		return model.Params, nil
	}

	// 默认参数
	defaultParams := ModelParams{
		Temperature:   0.8,
		TopP:          0.9,
		TopK:          40,
		Context:       2048,
		NumPredict:    512,
		RepeatPenalty: 1.1,
	}

	return defaultParams, nil
}

// SearchModels 搜索模型
func (m *ModelManager) SearchModels(params ModelSearchParams) ([]Model, error) {
	// 这里应该调用Ollama的模型库API或者维护一个本地模型库
	// 目前使用模拟数据
	allModels := []Model{
		{
			Name:       "llama3:8b",
			Model:      "llama3:8b",
			ModifiedAt: "2024-01-01T12:00:00Z",
			Size:       4670000000,
			Digest:     "sha256:abc123",
		},
		{
			Name:       "llama3:70b",
			Model:      "llama3:70b",
			ModifiedAt: "2024-01-02T12:00:00Z",
			Size:       40000000000,
			Digest:     "sha256:def456",
		},
		{
			Name:       "mistral:7b",
			Model:      "mistral:7b",
			ModifiedAt: "2024-01-03T12:00:00Z",
			Size:       4100000000,
			Digest:     "sha256:ghi789",
		},
		{
			Name:       "mixtral:8x7b",
			Model:      "mixtral:8x7b",
			ModifiedAt: "2024-01-04T12:00:00Z",
			Size:       46700000000,
			Digest:     "sha256:jkl012",
		},
		{
			Name:       "gemma:7b",
			Model:      "gemma:7b",
			ModifiedAt: "2024-01-05T12:00:00Z",
			Size:       5000000000,
			Digest:     "sha256:mno345",
		},
		{
			Name:       "phi3:3.8b",
			Model:      "phi3:3.8b",
			ModifiedAt: "2024-01-06T12:00:00Z",
			Size:       2500000000,
			Digest:     "sha256:pqr678",
		},
	}

	var filteredModels []Model

	// 如果没有查询条件，返回所有模型
	if params.Query == "" && len(params.Families) == 0 && len(params.Tags) == 0 {
		return allModels, nil
	}

	// 根据查询条件过滤模型
	for _, model := range allModels {
		match := true

		// 根据查询关键词过滤
		if params.Query != "" {
			lowerQuery := strings.ToLower(params.Query)
			lowerName := strings.ToLower(model.Name)

			// 支持正则表达式匹配
			if matched, _ := regexp.MatchString(lowerQuery, lowerName); !matched {
				// 普通字符串匹配
				if !strings.Contains(lowerName, lowerQuery) {
					match = false
				}
			}
		}

		// 根据家族过滤
		if len(params.Families) > 0 {
			familyMatch := false
			modelFamily := strings.Split(model.Name, ":")[0]
			for _, family := range params.Families {
				if strings.ToLower(modelFamily) == strings.ToLower(family) {
					familyMatch = true
					break
				}
			}
			if !familyMatch {
				match = false
			}
		}

		// 根据标签过滤
		if len(params.Tags) > 0 {
			tagMatch := false
			// 简化的标签匹配逻辑
			for _, tag := range params.Tags {
				lowerTag := strings.ToLower(tag)
				if strings.Contains(strings.ToLower(model.Name), lowerTag) {
					tagMatch = true
					break
				}
			}
			if !tagMatch {
				match = false
			}
		}

		if match {
			filteredModels = append(filteredModels, model)
		}
	}

	// 按名称排序
	sort.Slice(filteredModels, func(i, j int) bool {
		return filteredModels[i].Name < filteredModels[j].Name
	})

	return filteredModels, nil
}

// GetModelFamilies 获取模型家族列表
func (m *ModelManager) GetModelFamilies() []string {
	families := []string{"llama3", "mistral", "mixtral", "gemma", "phi3"}
	return families
}

// GetModelTags 获取模型标签列表
func (m *ModelManager) GetModelTags() []string {
	tags := []string{"embedding", "vision", "tools", "thinking", "8b", "70b", "7b", "3.8b"}
	return tags
}
