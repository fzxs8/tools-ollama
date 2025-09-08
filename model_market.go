package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fzxs8/duolasdk/core"
)

// ModelMarket 模型市场
type ModelMarket struct {
	ctx context.Context
	app *App
}

// OnlineModel 在线模型信息
type OnlineModel struct {
	Name        string `json:"name"`
	PullCount   int64  `json:"pull_count"`
	UpdatedAt   string `json:"updated_at"`
	Description string `json:"description"`
}

// NewModelMarket 创建新的模型市场
func NewModelMarket(app *App) *ModelMarket {
	return &ModelMarket{
		app: app,
	}
}

// SetContext 设置上下文
func (m *ModelMarket) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// SearchOnlineModels 在 ollamadb.dev 上搜索模型
func (m *ModelMarket) SearchOnlineModels(query string) ([]interface{}, error) {
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
