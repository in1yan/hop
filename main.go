package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	// ctx, cancel := context.WithCancel((context.Background()))
	// defer cancel()
	// Create application with options
	err := wails.Run(&options.App{
		Title:       "hop",
		Width:       500,
		Height:      200,
		Frameless:   true,
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 25, G: 23, B: 36, A: 180}, // Rose Pine base with transparency
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			go InstallHook(ctx)
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic, // Acrylic for frosty effect
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
