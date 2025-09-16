package main

import (
	"context"
	"fmt"

	"github.com/fzxs8/duolasdk/core"
)

// ModelMarket 模型市场
type ModelMarket struct {
	ctx    context.Context
	app    *App
	logger *core.AppLog
}

// NewModelMarket 创建新的模型市场
func NewModelMarket(app *App, logger *core.AppLog) *ModelMarket {
	return &ModelMarket{
		app:    app,
		logger: logger.WithPrefix("ModelMarket"),
	}
}

// SetContext 设置上下文
func (m *ModelMarket) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// SearchOnlineModels 在 ollamadb.dev 上搜索模型
func (m *ModelMarket) SearchOnlineModels(query string) ([]interface{}, error) {
	m.logger.Debug("开始在线搜索模型", "query", query)
	searchClient := core.NewHttp(m.logger.WithPrefix("SearchClient"))
	//searchClient.

	// 确保使用完整的绝对URL
	baseURL := "https://ollamadb.dev/api/v1/models"
	resp, err := searchClient.Get(baseURL, core.Options{
		Query: map[string]string{
			"search": query,
		},
	})
	if err != nil {
		m.logger.Error("在线搜索模型失败", "query", query, "error", err)
		return nil, fmt.Errorf("搜索模型失败: %w", err)
	}

	if err := HandleHTTPError(resp.StatusCode, resp.Body, m.logger, "在线搜索模型"); err != nil {
		return nil, err
	}

	var searchResult map[string]interface{}
	if err := UnmarshalJSONWithError([]byte(resp.Body), &searchResult, m.logger, "解析在线模型搜索结果"); err != nil {
		return nil, err
	}

	modelsField, err := ExtractFieldFromResponse(searchResult, "models")
	if err != nil {
		m.logger.Error("在线模型搜索结果中缺少 'models' 字段")
		return nil, fmt.Errorf("在搜索结果中找不到模型")
	}

	models, ok := modelsField.([]interface{})
	if !ok {
		m.logger.Error("在线模型搜索结果中 'models' 字段类型错误")
		return nil, fmt.Errorf("模型数据格式错误")
	}

	m.logger.Debug("成功发现在线模型", "count", len(models))
	return models, nil
}
