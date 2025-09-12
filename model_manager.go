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
	//if serverID == "local" {
	//	localConfig, err := m.app.configMgr.GetLocalConfig()
	//	if err != nil {
	//		// 如果获取配置失败，使用默认的本地服务器配置
	//		localConfig = OllamaServerConfig{BaseURL: "http://localhost:11434"}
	//		logger.Warnf("获取本地配置失败，使用默认配置: %v", err)
	//	}
	//	serverConfig = &localConfig
	//} else {
	servers, err := m.app.configMgr.GetServers()
	if err != nil {
		logger.Errorf("获取远程服务器列表失败: %v", err)
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
		logger.Errorf("找不到指定的服务器: %s", serverID)
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": err.Error()})
		return
	}
	//}

	logger.Infof("使用服务器配置: %+v", serverConfig)

	// 为特定服务器创建临时的HTTP客户端
	downloadClient := core.NewHttp(logger)
	downloadClient.Create(&core.Config{
		BaseURL: serverConfig.BaseURL,
	})

	// 修复：使用正确的字段名 "name" 而不是 "model"
	requestBody := map[string]interface{}{
		"name":   modelName,
		"stream": true,
	}

	logger.Infof("发送下载请求到 %s/api/pull，请求体: %+v", serverConfig.BaseURL, requestBody)

	resp, err := downloadClient.PostStream("/api/pull", core.Options{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    requestBody,
	})
	if err != nil {
		logger.Errorf("创建下载请求失败: %v", err)
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": "创建下载请求失败: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	logger.Infof("收到响应，状态码: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("下载失败，服务器响应: %s (状态码 %d)", string(bodyBytes), resp.StatusCode)
		logger.Errorf(errMsg)
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": errMsg})
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	errorOccurred := false
	var lastError string

	for scanner.Scan() {
		line := scanner.Bytes()
		logger.Debugf("收到进度行: %s", string(line))

		if len(line) == 0 {
			continue
		}

		var progressInfo map[string]interface{}
		if err := json.Unmarshal(line, &progressInfo); err != nil {
			logger.Warnf("无法解析下载进度JSON: %s, 错误: %v", string(line), err)
			continue
		}

		// 检查是否有错误信息
		if errorMsg, ok := progressInfo["error"]; ok && errorMsg != nil && errorMsg != "" {
			errorOccurred = true
			lastError = fmt.Sprintf("%v", errorMsg)
			logger.Errorf("下载过程中出现错误: %s", lastError)
			runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": lastError})
			// 出现错误时中断下载
			break
		}

		// 确保status字段存在
		if _, ok := progressInfo["status"]; !ok {
			progressInfo["status"] = "下载中"
		}

		// 添加默认的completed和total字段，以防它们不存在
		if _, ok := progressInfo["completed"]; !ok {
			progressInfo["completed"] = 0
		}
		if _, ok := progressInfo["total"]; !ok {
			progressInfo["total"] = 0
		}

		progressInfo["model"] = modelName
		runtime.EventsEmit(m.ctx, "model:download:progress", progressInfo)
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf("读取下载流失败: %v", err)
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": "读取下载流失败: " + err.Error()})
		return
	}

	// 只有在没有错误的情况下才发送下载完成事件
	if !errorOccurred {
		logger.Infof("模型 %s 下载完成", modelName)
		runtime.EventsEmit(m.ctx, "model:download:done", map[string]interface{}{"model": modelName})
	} else {
		logger.Infof("模型 %s 下载失败: %s", modelName, lastError)
		// 错误事件已经发送，无需再次发送
	}
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
