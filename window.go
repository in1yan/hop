package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Window struct {
	Handle uintptr
	Title  string
}

var (
	user32        = windows.NewLazyDLL("user32.dll")
	getWindowText = user32.NewProc("GetWindowTextW")
	setWindow     = user32.NewProc("SetForegroundWindow")
	getWindowRect = user32.NewProc("GetWindowRect")
	getAncestor   = user32.NewProc("GetAncestor")
)

const (
	GA_ROOT = 2
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
	return openWindows
}

func SetForegroundWindow(hwnd uintptr) error {
	ret, _, err := setWindow.Call(hwnd)
	if ret == 0 {
		return fmt.Errorf("SetForegroundWindow failed: %v", err)
	}
	return nil
}
