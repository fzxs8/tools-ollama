package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// RunningModel 运行中的模型
type RunningModel struct {
	Name      string      `json:"name"`
	Params    ModelParams `json:"params"`
	StartTime time.Time   `json:"start_time"`
	IsActive  bool        `json:"is_active"`
}

// 模型运行状态管理
var runningModels = make(map[string]*RunningModel)

// IsModelRunning 检查模型是否在运行
func (m *ModelManager) IsModelRunning(modelName string) bool {
	_, exists := runningModels[modelName]
	return exists
}

// RunModel 运行模型
func (m *ModelManager) RunModel(modelName string, params ModelParams) error {
	if _, exists := runningModels[modelName]; exists {
		return fmt.Errorf("模型 %s 已经在运行", modelName)
	}
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "hello",
		"stream": false,
	}
	response, err := m.app.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    requestBody,
	})
	if err != nil {
		return fmt.Errorf("启动模型失败: %v", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("启动模型失败，状态码: %d", response.StatusCode)
	}
	runningModels[modelName] = &RunningModel{
		Name:      modelName,
		Params:    params,
		StartTime: time.Now(),
		IsActive:  true,
	}
	runtime.EventsEmit(m.ctx, "model:started", map[string]interface{}{
		"name": modelName,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
	return nil
}

// StopModel 停止模型
func (m *ModelManager) StopModel(modelName string) error {
	if _, exists := runningModels[modelName]; !exists {
		return fmt.Errorf("模型 %s 未在运行", modelName)
	}
	requestBody := map[string]interface{}{
		"model":      modelName,
		"keep_alive": 0,
	}
	_, err := m.app.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    requestBody,
	})
	if err != nil {
		delete(runningModels, modelName)
		runtime.EventsEmit(m.ctx, "model:stopped", map[string]interface{}{"name": modelName})
		return fmt.Errorf("停止模型时出现警告(但UI已更新): %v", err)
	}
	delete(runningModels, modelName)
	runtime.EventsEmit(m.ctx, "model:stopped", map[string]interface{}{"name": modelName})
	return nil
}

// TestModel 测试模型
func (m *ModelManager) TestModel(modelName string, prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": prompt,
		"stream": false,
	}
	response, err := m.app.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    requestBody,
	})
	if err != nil {
		return "", fmt.Errorf("测试模型失败: %v", err)
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("测试模型失败，状态码: %d", response.StatusCode)
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	if responseText, ok := result["response"].(string); ok {
		return responseText, nil
	}
	return "", fmt.Errorf("未在响应中找到 'response' 字段")
}

// DownloadModel 下载模型
func (m *ModelManager) DownloadModel(serverID string, modelName string) {
	logger := core.NewLogger(&core.LoggerOption{Type: "console", Level: "debug", Prefix: "DownloadClient"})
	logger.Infof("开始下载模型: %s, 服务器ID: %s", modelName, serverID)

	var serverConfig *OllamaServerConfig
	if serverID == "local" {
		localConfig, err := m.app.configMgr.GetLocalConfig()
		if err != nil {
			// 如果获取配置失败，使用默认的本地服务器配置
			localConfig = OllamaServerConfig{BaseURL: "http://localhost:11434"}
		}
		serverConfig = &localConfig
	} else {
		servers, err := m.app.configMgr.GetRemoteServers()
		if err != nil {
			runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": "获取远程服务器列表失败: " + err.Error()})
			return
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
			err := fmt.Errorf("找不到指定的服务器: %s", serverID)
			runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": err.Error()})
			return
		}
	}

	downloadClient := core.NewHttp(logger)
	downloadClient.Create(&core.Config{BaseURL: serverConfig.BaseURL})

	requestBody := map[string]interface{}{
		"name":   modelName,
		"stream": true,
	}

	resp, err := downloadClient.PostStream("/api/pull", core.Options{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    requestBody,
	})
	if err != nil {
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": "创建下载请求失败: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("下载失败，服务器响应: %s (状态码 %d)", string(bodyBytes), resp.StatusCode)
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": errMsg})
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var progressInfo map[string]interface{}
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if err := json.Unmarshal(line, &progressInfo); err != nil {
			logger.Warnf("无法解析下载进度JSON: %s, 错误: %v", string(line), err)
			continue
		}
		progressInfo["model"] = modelName
		runtime.EventsEmit(m.ctx, "model:download:progress", progressInfo)
	}

	if err := scanner.Err(); err != nil {
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": "读取下载流失败: " + err.Error()})
		return
	}

	runtime.EventsEmit(m.ctx, "model:download:done", map[string]interface{}{"model": modelName})
}

// SetModelParams 设置模型参数
func (m *ModelManager) SetModelParams(modelName string, params ModelParams) error {
	if model, exists := runningModels[modelName]; exists {
		model.Params = params
	}
	return nil
}

// GetModelParams 获取模型参数
func (m *ModelManager) GetModelParams(modelName string) (ModelParams, error) {
	if model, exists := runningModels[modelName]; exists {
		return model.Params, nil
	}
	return ModelParams{
		Temperature:   0.8,
		TopP:          0.9,
		TopK:          40,
		Context:       2048,
		NumPredict:    512,
		RepeatPenalty: 1.1,
	}, nil
}

// SearchModels 在 ollamadb.dev 上搜索模型
func (m *ModelManager) SearchModels(query string) ([]interface{}, error) {
	logger := core.NewLogger(&core.LoggerOption{Type: "console", Level: "debug", Prefix: "SearchClient"})
	searchClient := core.NewHttp(logger)

	resp, err := searchClient.Get("https://ollamadb.dev/api/v1/models", core.Options{
		Query: map[string]string{
			"search": query,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("搜索模型失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("搜索模型失败: 状态码 %d", resp.StatusCode)
	}

	var searchResult map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &searchResult)
	if err != nil {
		return nil, fmt.Errorf("解析搜索结果失败: %w", err)
	}

	models, ok := searchResult["models"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("在搜索结果中找不到模型")
	}

	return models, nil
}
