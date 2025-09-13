package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk/core"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ModelManager 模型管理器
type ModelManager struct {
	ctx       context.Context
	app       *App // 保留对App的引用以访问全局状态和方法
	logger    *core.AppLog
	configMgr *OllamaConfigManager
}

// 模型运行状态管理 (在内存中)
var runningModels = make(map[string]*types.RunningModel)

// NewModelManager 创建新的模型管理器
func NewModelManager(app *App, configMgr *OllamaConfigManager, logger *core.AppLog) *ModelManager {
	return &ModelManager{
		app:       app,
		logger:    logger.WithPrefix("ModelManager"),
		configMgr: configMgr,
	}
}

// SetContext 设置上下文
func (m *ModelManager) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// ListModelsByServer 根据服务器获取模型列表
func (m *ModelManager) ListModelsByServer(serverID string) ([]types.Model, error) {
	m.logger.Debug("开始获取模型列表", "serverID", serverID)
	serverConfig, err := m.configMgr.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	// 为特定服务器创建临时的HTTP客户端
	m.logger.Debug("开始获取模型列表", "serverConfig", serverConfig)
	tempClient := core.NewHttp(m.logger.WithPrefix("TempClient"))
	tempClient.Create(&core.Config{
		BaseURL: serverConfig.BaseURL,
	})

	response, err := tempClient.Get("/api/tags", core.Options{})
	if err != nil {
		m.logger.Error("请求Ollama API [/api/tags] 失败", "serverID", serverID, "error", err)
		return nil, err
	}

	var result types.ListModelsResponse
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		m.logger.Error("反序列化模型列表响应失败", "error", err)
		return nil, err
	}

	// 检查并设置每个模型的运行状态
	for i := range result.Models {
		result.Models[i].IsRunning = m.IsModelRunning(result.Models[i].Name)
	}

	m.logger.Debug("成功获取并处理了 %d 个模型", len(result.Models))
	return result.Models, nil
}

// IsModelRunning 检查模型是否在运行
func (m *ModelManager) IsModelRunning(modelName string) bool {
	_, exists := runningModels[modelName]
	return exists
}

// RunModel 运行模型
func (m *ModelManager) RunModel(modelName string, params types.ModelParams) error {
	m.logger.Info("准备运行模型", "modelName", modelName)
	if _, exists := runningModels[modelName]; exists {
		return fmt.Errorf("模型 %s 已经在运行", modelName)
	}
	requestBody := map[string]interface{}{
		"model":      modelName,
		"keep_alive": -1, // 设置为-1以保持模型加载
	}
	response, err := m.app.httpClient.Post("/api/generate", core.Options{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    requestBody,
	})
	if err != nil {
		m.logger.Error("启动模型失败", "modelName", modelName, "error", err)
		return fmt.Errorf("启动模型失败: %v", err)
	}
	if response.StatusCode >= 400 {
		m.logger.Error("启动模型失败，状态码非200", "modelName", modelName, "statusCode", response.StatusCode, "body", response.Body)
		return fmt.Errorf("启动模型失败，状态码: %d, 响应: %s", response.StatusCode, response.Body)
	}

	// 更新内存中的运行状态
	runningModels[modelName] = &types.RunningModel{
		Name:   modelName,
		Params: params,
	}

	m.logger.Info("模型成功启动", "modelName", modelName)
	runtime.EventsEmit(m.ctx, "model:started", modelName)
	return nil
}

// StopModel 停止模型
func (m *ModelManager) StopModel(modelName string) error {
	m.logger.Info("准备停止模型", "modelName", modelName)
	if _, exists := runningModels[modelName]; !exists {
		return fmt.Errorf("模型 %s 未在运行", modelName)
	}

	// 在Ollama中，没有直接的“停止”API，通常是通过加载另一个模型或发送一个keep_alive:0的请求来卸载它
	// 这里我们只从内存中移除，前端会看到状态更新。实际的模型卸载由Ollama的策略决定。
	delete(runningModels, modelName)
	m.logger.Info("模型已从管理器中移除（停止）", "modelName", modelName)
	runtime.EventsEmit(m.ctx, "model:stopped", modelName)
	return nil
}

// DeleteModel 删除模型
func (m *ModelManager) DeleteModel(modelName string) error {
	m.logger.Info("准备删除模型", "modelName", modelName)
	requestBody := map[string]interface{}{
		"name": modelName,
	}

	response, err := m.app.httpClient.Do("DELETE", "/api/delete", core.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
	})
	if err != nil {
		m.logger.Error("删除模型请求失败", "modelName", modelName, "error", err)
		return err
	}

	if response.StatusCode >= 400 {
		m.logger.Error("删除模型失败，状态码非200", "modelName", modelName, "statusCode", response.StatusCode, "body", response.Body)
		return fmt.Errorf("删除模型失败，状态码: %d, 内容: %s", response.StatusCode, response.Body)
	}

	m.logger.Info("模型删除成功", "modelName", modelName)
	return nil
}

// TestModel 测试模型
func (m *ModelManager) TestModel(modelName string, prompt string) (string, error) {
	m.logger.Debug("开始测试模型", "modelName", modelName)
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
	logger := m.logger.WithPrefix("DownloadClient")
	logger.Infof("开始下载模型: %s, 服务器ID: %s", modelName, serverID)

	serverConfig, err := m.configMgr.GetServerByID(serverID)
	if err != nil {
		logger.Errorf("获取服务器配置失败: %v", err)
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": err.Error()})
		return
	}

	logger.Infof("使用服务器配置: %+v", serverConfig)

	downloadClient := core.NewHttp(logger)
	downloadClient.Create(&core.Config{
		BaseURL: serverConfig.BaseURL,
	})

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
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var progressInfo map[string]interface{}
		if err := json.Unmarshal(line, &progressInfo); err != nil {
			logger.Warnf("无法解析下载进度JSON: %s, 错误: %v", string(line), err)
			continue
		}

		if errorMsg, ok := progressInfo["error"]; ok {
			logger.Errorf("下载过程中出现错误: %s", errorMsg)
			runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": errorMsg})
			return // 中断下载
		}

		progressInfo["model"] = modelName
		runtime.EventsEmit(m.ctx, "model:download:progress", progressInfo)
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf("读取下载流失败: %v", err)
		runtime.EventsEmit(m.ctx, "model:download:error", map[string]interface{}{"model": modelName, "error": "读取下载流失败: " + err.Error()})
		return
	}

	logger.Infof("模型 %s 下载完成", modelName)
	runtime.EventsEmit(m.ctx, "model:download:done", map[string]interface{}{"model": modelName})
}

// SetModelParams 设置模型参数
func (m *ModelManager) SetModelParams(modelName string, params types.ModelParams) error {
	if model, exists := runningModels[modelName]; exists {
		m.logger.Debug("更新正在运行模型的参数", "modelName", modelName)
		model.Params = params
		return nil
	}
	m.logger.Warn("尝试为未在运行的模型设置参数", "modelName", modelName)
	// 即使模型未运行，也可以考虑将其参数保存到store中，以便下次运行
	return nil
}

// GetModelParams 获取模型参数
func (m *ModelManager) GetModelParams(modelName string) (types.ModelParams, error) {
	if model, exists := runningModels[modelName]; exists {
		m.logger.Debug("获取正在运行模型的参数", "modelName", modelName)
		return model.Params, nil
	}
	m.logger.Debug("模型未在运行，返回默认参数", "modelName", modelName)
	// 如果模型未运行，返回默认参数
	return types.ModelParams{
		Temperature:   0.8,
		TopP:          0.9,
		TopK:          40,
		Context:       2048,
		NumPredict:    512,
		RepeatPenalty: 1.1,
	}, nil
}
