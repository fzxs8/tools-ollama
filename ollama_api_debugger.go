package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk/core"
)

// OllamaApiDebugger 模块负责处理所有与 API 调试器相关的功能
type OllamaApiDebugger struct {
	ctx       context.Context
	logger    *core.AppLog
	configMgr *OllamaConfigManager
}

// NewOllamaApiDebugger 创建一个新的 OllamaApiDebugger 实例
func NewOllamaApiDebugger(logger *core.AppLog, configMgr *OllamaConfigManager) *OllamaApiDebugger {
	return &OllamaApiDebugger{
		logger:    logger.WithPrefix("ApiDebugger"),
		configMgr: configMgr,
	}
}

// SetContext 在应用启动时设置上下文
func (d *OllamaApiDebugger) SetContext(ctx context.Context) {
	d.ctx = ctx
}

// SendHttpRequest 使用 Go 标准库 net/http 直接处理前端发送的 HTTP 请求
func (d *OllamaApiDebugger) SendHttpRequest(request types.ApiRequest) (types.ApiResponse, error) {
	d.logger.Debug("Received API Debugger Request", "method", request.Method, "path", request.Path, "serverID", request.SelectedServerID)

	var apiResponse types.ApiResponse
	startTime := time.Now()

	// 1. 获取 Base URL
	var baseURL string
	if request.SelectedServerID != "" {
		serverConfig, err := d.configMgr.GetServerByID(request.SelectedServerID)
		if err != nil {
			d.logger.Errorf("Failed to get server config for ID %s: %v", request.SelectedServerID, err)
			apiResponse.Error = fmt.Sprintf("无法获取服务器配置: %v", err)
			return apiResponse, nil
		}
		baseURL = serverConfig.BaseURL
	} else {
		apiResponse.Error = "未选择Ollama服务或服务配置无效"
		return apiResponse, nil
	}

	// 2. 构建完整的 URL (BaseURL + Path + Query Params)
	queryParams := make(map[string]string)
	for _, param := range request.QueryParams {
		if param.Enabled {
			queryParams[param.Key] = param.Value
		}
	}

	finalURL, err := BuildURLWithQuery(baseURL, request.Path, queryParams)
	if err != nil {
		d.logger.Errorf("Failed to build URL: %v", err)
		apiResponse.Error = fmt.Sprintf("URL构建失败: %v", err)
		return apiResponse, nil
	}

	// 3. 准备请求体 (io.Reader)
	var bodyReader io.Reader
	var contentType string

	switch request.Body.Type {
	case types.RequestBodyTypeRaw:
		bodyReader = strings.NewReader(request.Body.RawContent)
		contentType = string(request.Body.RawContentType)
	case types.RequestBodyTypeFormData:
		formData := url.Values{}
		for _, field := range request.Body.FormData {
			formData.Set(field.Key, field.Value)
		}
		bodyReader = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	case types.RequestBodyTypeNone:
		bodyReader = nil
	}

	// 4. 使用 net/http 创建请求
	req, err := http.NewRequestWithContext(d.ctx, strings.ToUpper(request.Method), finalURL, bodyReader)
	if err != nil {
		d.logger.Errorf("Failed to create HTTP request: %v", err)
		apiResponse.Error = fmt.Sprintf("创建请求失败: %v", err)
		return apiResponse, nil
	}

	// 5. 手动设置所有请求头
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for _, header := range request.Headers {
		if header.Enabled {
			if strings.EqualFold(header.Key, "Content-Type") {
				req.Header.Set(header.Key, header.Value)
			} else {
				req.Header.Add(header.Key, header.Value)
			}
		}
	}

	// 6. 发送请求
	d.logger.Debug("Sending HTTP request via stdlib", "method", req.Method, "url", req.URL, "headers", req.Header)
	client := CreateHTTPClientWithTimeout(30 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		d.logger.Errorf("Failed to send HTTP request: %v", err)
		apiResponse.Error = fmt.Sprintf("发送请求失败: %v", err)
		return apiResponse, nil
	}
	defer resp.Body.Close()

	// 7. 读取并处理响应
	respBody, err := ReadResponseBody(resp)
	if err != nil {
		d.logger.Errorf("Failed to read response body: %v", err)
		apiResponse.Error = fmt.Sprintf("读取响应失败: %v", err)
		return apiResponse, nil
	}

	apiResponse.StatusCode = resp.StatusCode
	apiResponse.StatusText = http.StatusText(resp.StatusCode)
	apiResponse.Body = respBody

	for k, v := range resp.Header {
		apiResponse.Headers = append(apiResponse.Headers, types.RequestHeader{Key: k, Value: strings.Join(v, ", "), Enabled: true})
	}

	apiResponse.RequestDurationMs = time.Since(startTime).Milliseconds()
	d.logger.Debug("API Debugger Request completed", "statusCode", apiResponse.StatusCode, "duration", apiResponse.RequestDurationMs)

	return apiResponse, nil
}

// GetOllamaServers 用于前端获取Ollama服务器配置列表
func (d *OllamaApiDebugger) GetOllamaServers() ([]types.OllamaServerConfig, error) {
	return d.configMgr.GetServers()
}
