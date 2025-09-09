package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetWindows() []Window {
	return GetOpenWindows()
}

func (a *App) SetFocus(hwnd uintptr) error {
	runtime.Hide(a.ctx)
	return SetForegroundWindow(hwnd)
}

func (a *App) FocusHiddenInput() {
	runtime.EventsEmit(a.ctx, "focus:input")
}
