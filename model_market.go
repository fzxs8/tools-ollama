package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fzxs8/duolasdk/core"
	"net/http"
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

	resp, err := searchClient.Get("https://ollamadb.dev/api/v1/models", core.Options{
		Query: map[string]string{
			"search": query,
		},
	})
	if err != nil {
		m.logger.Error("在线搜索模型失败", "query", query, "error", err)
		return nil, fmt.Errorf("搜索模型失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		m.logger.Error("在线搜索模型API返回非200状态码", "query", query, "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("搜索模型失败: 状态码 %d", resp.StatusCode)
	}

	var searchResult map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &searchResult)
	if err != nil {
		m.logger.Error("解析在线模型搜索结果失败", "error", err)
		return nil, fmt.Errorf("解析搜索结果失败: %w", err)
	}

	models, ok := searchResult["models"].([]interface{})
	if !ok {
		m.logger.Error("在线模型搜索结果中缺少 'models' 字段")
		return nil, fmt.Errorf("在搜索结果中找不到模型")
	}

	m.logger.Debug("成功发现在线模型", "count", len(models))
	return models, nil
}
