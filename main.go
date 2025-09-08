package main

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 设置信号处理，避免与GTK冲突
	setupSignalHandling()

	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Ollama 客户端",
		Width:  1366,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
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
