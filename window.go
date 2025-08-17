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
	openWindows   []Window
)

func main() {

	cb := syscall.NewCallback(func(hwnd windows.HWND, lparam uintptr) uintptr {
		buf := make([]uint16, 255)
		getWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		title := syscall.UTF16ToString(buf)
		if title != "" && windows.IsWindowVisible(hwnd) {
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
	fmt.Println("Enter the window to focus: ")
	var hwnd uintptr
	fmt.Scanf("%d", &hwnd)
	setWindow.Call(hwnd)
}
