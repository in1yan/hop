//go:build !windows
// +build !windows

package main

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

const (
	GA_ROOT         = 2
	WM_KEYDOWN      = 0x0100
	WM_KEYUP        = 0x0101
	WH_KEYBOARD_LL  = 13
	ASFW_ANY        = ^uint32(0) // -1
	VK_MENU         = 0x12       // Alt key
	KEYEVENTF_KEYUP = 0x0002
	SW_RESTORE      = 9
)

var (
	IsWindowVisible bool = true
)

// Stub functions for non-Windows platforms
func GetOpenWindows() []Window {
	return []Window{}
}

func SetForegroundWindow(hwnd uintptr) error {
	return nil
}
