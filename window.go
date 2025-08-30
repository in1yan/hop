package main

import (
	"context"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
)

type Window struct {
	Handle uintptr
	Title  string
}
type KBDLLHOOKSTRUCT struct {
	VkCode    uint32
	ScanCode  uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

var (
	user32          = windows.NewLazyDLL("user32.dll")
	getWindowText   = user32.NewProc("GetWindowTextW")
	setWindow       = user32.NewProc("SetForegroundWindow")
	getWindowRect   = user32.NewProc("GetWindowRect")
	getAncestor     = user32.NewProc("GetAncestor")
	setHookEx       = user32.NewProc("SetWindowsHookExW")
	unhookEx        = user32.NewProc("UnhookWindowsHookEx")
	getMessage      = user32.NewProc("GetMessageA")
	callNextHookEx  = user32.NewProc("CallNextHookEx")
	hookHandle      windows.Handle
	IsWindowVisible bool = true
)

const (
	GA_ROOT        = 2
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WH_KEYBOARD_LL = 13
)

func GetOpenWindows() []Window {
	var openWindows []Window
	cb := syscall.NewCallback(func(hwnd windows.HWND, lparam uintptr) uintptr {
		buf := make([]uint16, 255)
		getWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		title := syscall.UTF16ToString(buf)
		var rect = windows.Rect{}
		root, _, _ := getAncestor.Call(uintptr(hwnd), GA_ROOT)
		getWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
		if root != uintptr(hwnd) {
			return 1
		}
		if title != "" && windows.IsWindowVisible(hwnd) && rect.Right-rect.Left > 0 && rect.Bottom-rect.Top > 0 {
			openWindows = append(openWindows, Window{Handle: uintptr(hwnd), Title: title})
		}
		return 1
	})

	err := windows.EnumWindows(cb, unsafe.Pointer(nil))
	if err != nil {
		fmt.Println("Error:", err)
	}
	openWindows = openWindows[1 : len(openWindows)-4]
	return openWindows
}

func SetForegroundWindow(hwnd uintptr) error {
	ret, _, err := setWindow.Call(hwnd)
	if ret == 0 {
		return fmt.Errorf("SetForegroundWindow failed: %v", err)
	}
	return nil
}
func InstallHook(ctx context.Context) {
	var cb = syscall.NewCallback(func(nCode int, wparam uintptr, lparam uintptr) uintptr {
		if nCode >= 0 && (wparam == WM_KEYDOWN) {
			kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
			if kbdStruct.VkCode == windows.VK_RSHIFT {
				if !IsWindowVisible {
					// runtime.WindowReload(ctx)
					runtime.Show(ctx)
					fmt.Println("Root:", root)
					IsWindowVisible = true
					runtime.EventsEmit(ctx, "windows:update", GetOpenWindows())
				} else if IsWindowVisible {
					runtime.Hide(ctx)
					IsWindowVisible = false
				}
				fmt.Println("Right Shift Pressed")
			}
		}
		ret, _, _ := callNextHookEx.Call(0, uintptr(nCode), wparam, lparam)
		return ret
	})
	hook, _, _ := setHookEx.Call(
		WH_KEYBOARD_LL,
		cb,
		0,
		0,
	)
	hookHandle = windows.Handle(hook)
	var msg struct {
		hwnd    uintptr
		message uint32
		wparam  uintptr
		lparam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}
	for {
		ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
	}

	if hookHandle != 0 {
		unhookEx.Call(uintptr(hookHandle))
	}
}
