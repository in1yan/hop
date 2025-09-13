package main

import (
	"context"
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed icons/icon.ico
var iconData []byte

func onReady(ctx context.Context) {
	systray.SetIcon(iconData)
	systray.SetTitle("Hop")
	systray.SetTooltip("Ready to Hop")
	mQuit := systray.AddMenuItem("Quit", "Quit Hopping")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		runtime.Quit(ctx)
	}()
}
func StartTray(ctx context.Context) {
	systray.Run(func() { onReady(ctx) }, nil)
}
