package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Window struct {
	Handle uintptr
	title  string
}

var (
	user32        = windows.NewLazyDLL("user32.dll")
	getWindowText = user32.NewProc("GetWindowTextW")
	setWindow     = user32.NewProc("SetForegroundWindow")
	getWindowRect = user32.NewProc("GetWindowRect")
	getAncestor   = user32.NewProc("GetAncestor")
	getWindowLong = user32.NewProc("GetWindowLongW")
	openWindows   []Window
)

const (
	GA_ROOT = 2
)

func main() {

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
			openWindows = append(openWindows, Window{Handle: uintptr(hwnd), title: title})
		}
		return 1
	})

	err := windows.EnumWindows(cb, unsafe.Pointer(nil))
	if err != nil {
		fmt.Println("Error:", err)
	}
	if openWindows != nil {
		for _, window := range openWindows {
			fmt.Println("HWND: ", window.Handle, "Title: ", window.title)
		}
	}
	// fmt.Println("Enter the window to focus: ")
	// var hwnd uintptr
	// fmt.Scanf("%d", &hwnd)
	// setWindow.Call(hwnd)
}
