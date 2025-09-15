package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk/core"
	"github.com/google/uuid"
)

// ToCoreMessages 将types.Message转换为core.Message
// 这是一个公共工具函数，用于在不同模块间转换消息类型
func ToCoreMessages(messages []types.Message) []core.Message {
	coreMessages := make([]core.Message, len(messages))
	for i, msg := range messages {
		coreMessages[i] = core.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return coreMessages
}

// EnsureHTTPPrefix 确保URL包含协议前缀
func EnsureHTTPPrefix(url string) string {
	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

// ConvertToModelParams 将map参数转换为ModelParams结构体
func ConvertToModelParams(params map[string]interface{}) types.ModelParams {
	modelParams := types.ModelParams{
		Temperature: 0.8, TopP: 0.9, TopK: 40,
		Context: 2048, NumPredict: 512, RepeatPenalty: 1.1,
	}

	if temp, ok := params["temperature"].(float64); ok {
		modelParams.Temperature = temp
	}
	if topP, ok := params["topP"].(float64); ok {
		modelParams.TopP = topP
	}
	if topK, ok := params["topK"].(float64); ok {
		modelParams.TopK = int(topK)
	}
	if ct, ok := params["context"].(float64); ok {
		modelParams.Context = int(ct)
	}
	if numPredict, ok := params["numPredict"].(float64); ok {
		modelParams.NumPredict = int(numPredict)
	}
	if repeatPenalty, ok := params["repeatPenalty"].(float64); ok {
		modelParams.RepeatPenalty = repeatPenalty
	}

	return modelParams
}

// ConvertFromModelParams 将ModelParams结构体转换为map
func ConvertFromModelParams(params types.ModelParams) map[string]interface{} {
	return map[string]interface{}{
		"temperature":   params.Temperature,
		"topP":          params.TopP,
		"topK":          params.TopK,
		"context":       params.Context,
		"numPredict":    params.NumPredict,
		"repeatPenalty": params.RepeatPenalty,
	}
}

// GenerateAPIDocs 生成API文档示例
func GenerateAPIDocs(ip string, port int) map[string]string {
	if ip == "0.0.0.0" {
		ip = "127.0.0.1" // 为方便用户复制，将 0.0.0.0 显示为 127.0.0.1
	}

	nonStreamingCurl := fmt.Sprintf(
		`curl http://%s:%d/v1/chat/completions -X POST \
-H "Content-Type: application/json" \
-d '{
 "model": "llama3",
 "messages": [
   {
     "role": "user",
     "content": "你好，介绍一下你自己"
   }
 ],
 "stream": false
}'`,
		ip, port,
	)
	streamingCurl := fmt.Sprintf(
		`curl http://%s:%d/v1/chat/completions -X POST \
-H "Content-Type: application/json" \
-d '{
 "model": "llama3",
 "messages": [
   {
     "role": "user",
     "content": "给我讲一个关于程序员的笑话"
   }
 ],
 "stream": true
}'`,
		ip, port,
	)

	return map[string]string{
		"非流式请求": nonStreamingCurl,
		"流式请求":  streamingCurl,
	}
}

// UnmarshalJSONWithError 通用JSON反序列化函数，带错误处理
func UnmarshalJSONWithError(data []byte, v interface{}, logger *core.AppLog, context string) error {
	if err := json.Unmarshal(data, v); err != nil {
		logger.Error(fmt.Sprintf("%s JSON反序列化失败", context), "error", err, "data", string(data))
		return fmt.Errorf("%s JSON反序列化失败: %w", context, err)
	}
	return nil
}

// MarshalJSONWithError 通用JSON序列化函数，带错误处理
func MarshalJSONWithError(v interface{}, logger *core.AppLog, context string) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		logger.Error(fmt.Sprintf("%s JSON序列化失败", context), "error", err)
		return nil, fmt.Errorf("%s JSON序列化失败: %w", context, err)
	}
	return data, nil
}

// GenerateUniqueID 生成唯一ID
func GenerateUniqueID() string {
	return uuid.New().String()
}

// GetCurrentTimestamp 获取当前时间戳（毫秒）
func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// ReverseSlice 反转切片（用于时间倒序排列）
func ReverseSlice[T any](slice []T) {
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-1-i] = slice[len(slice)-1-i], slice[i]
	}
}

// HandleHTTPError 处理HTTP错误响应
func HandleHTTPError(statusCode int, body string, logger *core.AppLog, operation string) error {
	if statusCode >= 400 {
		errMsg := fmt.Sprintf("%s失败，状态码: %d, 响应: %s", operation, statusCode, body)
		logger.Error(errMsg)
		return fmt.Errorf(errMsg)
	}
	return nil
}

// ExtractMessageContent 从响应中提取消息内容
func ExtractMessageContent(response map[string]interface{}) (string, error) {
	if message, ok := response["message"].(map[string]interface{}); ok {
		if content, ok := message["content"].(string); ok {
			return content, nil
		}
	}
	return "", fmt.Errorf("未在响应中找到消息内容")
}

// ExtractResponseContent 从生成响应中提取内容
func ExtractResponseContent(response map[string]interface{}) (string, error) {
	if responseText, ok := response["response"].(string); ok {
		return responseText, nil
	}
	return "", fmt.Errorf("未在响应中找到 'response' 字段")
}

// CreateHTTPClientWithTimeout 创建带超时的HTTP客户端
func CreateHTTPClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

// BuildURLWithQuery 构建带查询参数的URL
func BuildURLWithQuery(baseURL, path string, queryParams map[string]string) (string, error) {
	fullURL, err := url.JoinPath(baseURL, path)
	if err != nil {
		return "", fmt.Errorf("URL拼接失败: %w", err)
	}

	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("URL解析失败: %w", err)
	}

	q := parsedURL.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	parsedURL.RawQuery = q.Encode()
	return parsedURL.String(), nil
}

// ExtractFieldFromResponse 从响应中提取指定字段
func ExtractFieldFromResponse(response map[string]interface{}, fieldName string) (interface{}, error) {
	if field, ok := response[fieldName]; ok {
		return field, nil
	}
	return nil, fmt.Errorf("在响应中找不到字段 '%s'", fieldName)
}

// ValidateServerConfig 验证服务器配置
func ValidateServerConfig(config *types.OllamaServerConfig) error {
	if config.ID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if config.Name == "" {
		return fmt.Errorf("服务器名称不能为空")
	}
	if config.BaseURL == "" {
		return fmt.Errorf("服务器URL不能为空")
	}
	return nil
}

// FindActiveServer 从服务器列表中查找活跃服务器
func FindActiveServer(servers []types.OllamaServerConfig) *types.OllamaServerConfig {
	for i, server := range servers {
		if server.IsActive {
			return &servers[i]
		}
	}
	return nil
}

// SetServerActiveStatus 设置服务器活跃状态
func SetServerActiveStatus(servers []types.OllamaServerConfig, targetServerID string) ([]types.OllamaServerConfig, bool) {
	found := false
	for i := range servers {
		if servers[i].ID == targetServerID {
			servers[i].IsActive = true
			found = true
		} else {
			servers[i].IsActive = false
		}
	}
	return servers, found
}

// ReadResponseBody 读取HTTP响应体
func ReadResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}
	return string(body), nil
}
