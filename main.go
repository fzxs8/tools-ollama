package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 设置信号处理，避免与GTK冲突
	setupSignalHandling()

	// Create an instance of the app structure
	app := NewApp()

	// Initialize storage
	store := duolasdk.NewStore(
		core.StoreOption{
			FileName: "ollama-client.db",
		})

	// Create instances of other components
	modelManager := NewModelManager(app)
	chatManager := NewChatManager(app.ctx, store)
	ollamaConfig := NewOllamaConfigManager(store)
	modelMarket := NewModelMarket(app)
	promptPilot := NewPromptPilot(store)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Ollama 客户端",
		Width:  1366,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			modelManager.SetContext(ctx)
			chatManager.SetContext(ctx)
			//ollamaConfig.SetContext(ctx)
			modelMarket.SetContext(ctx)
			promptPilot.Startup(ctx)

			// Set HTTP client for prompt pilot
			promptPilot.SetHTTPClient(app.httpClient)

			// 初始化日志
			logger := core.NewLogger(&core.LoggerOption{Type: "console", Level: "debug", Prefix: "main"})
			logger.Info("Application started successfully")
		},
		Bind: []interface{}{
			app,
			modelManager,
			chatManager,
			ollamaConfig,
			modelMarket,
			promptPilot,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// setupSignalHandling 设置信号处理以避免与GTK冲突
func setupSignalHandling() {
	// 忽略SIGPIPE信号，避免GTK和Go运行时之间的冲突
	signal.Ignore(syscall.SIGPIPE)

	// 设置一个简单的信号处理程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// 优雅地关闭应用
		os.Exit(0)
	}()
}
