package main

import (
	"context"
	"fmt"
	"sync"
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
	user32                  = windows.NewLazyDLL("user32.dll")
	getWindowText           = user32.NewProc("GetWindowTextW")
	setWindow               = user32.NewProc("SetForegroundWindow")
	getWindowRect           = user32.NewProc("GetWindowRect")
	getAncestor             = user32.NewProc("GetAncestor")
	setHookEx               = user32.NewProc("SetWindowsHookExW")
	unhookEx                = user32.NewProc("UnhookWindowsHookEx")
	getMessage              = user32.NewProc("GetMessageA")
	callNextHookEx          = user32.NewProc("CallNextHookEx")
	getWindow               = user32.NewProc("FindWindowW")
	setActiveWindow         = user32.NewProc("SetActiveWindow")
	procAllowSetForeground  = user32.NewProc("AllowSetForegroundWindow")
	procKeybdEvent          = user32.NewProc("keybd_event")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	isMinimized             = user32.NewProc("IsIconic")
	procShowWindow          = user32.NewProc("ShowWindow")
	hookHandle              windows.Handle
	windowVisibilityMutex   sync.RWMutex
	isWindowVisible         bool = false
)

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

// Thread-safe functions for window visibility state
func getWindowVisibility() bool {
	windowVisibilityMutex.RLock()
	defer windowVisibilityMutex.RUnlock()
	return isWindowVisible
}

func setWindowVisibility(visible bool) {
	windowVisibilityMutex.Lock()
	defer windowVisibilityMutex.Unlock()
	isWindowVisible = visible
}

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
	fmt.Println("Open windows: ", openWindows)
	return openWindows
}

func SetForegroundWindow(hwnd uintptr) error {
	ismin, _, _ := isMinimized.Call(hwnd)
	if ismin != 0 {
		procShowWindow.Call(hwnd, SW_RESTORE)
	}
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
				if !getWindowVisibility() {
					runtime.EventsEmit(ctx, "windows:update", GetOpenWindows())

					hwnd, _, _ := getWindow.Call(
						0,
						uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("hop"))),
					)

					runtime.Show(ctx)

					if hwnd == 0 {
						hwnd, _, _ = procGetForegroundWindow.Call()
					}

					procAllowSetForeground.Call(uintptr(ASFW_ANY))

					ret := SetForegroundWindow(hwnd)
					if ret != nil {
						procKeybdEvent.Call(uintptr(VK_MENU), 0, 0, 0)
						SetForegroundWindow(hwnd)
						procKeybdEvent.Call(uintptr(VK_MENU), 0, KEYEVENTF_KEYUP, 0)
					}

					setActiveWindow.Call(hwnd)

					setWindowVisibility(true)
					runtime.EventsEmit(ctx, "focus:input")
				} else {
					runtime.Hide(ctx)
					setWindowVisibility(false)
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

	// Message loop with context cancellation support
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, cleanup and exit
			if hookHandle != 0 {
				unhookEx.Call(uintptr(hookHandle))
			}
			return
		default:
			ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if ret == 0 {
				// WM_QUIT message received
				break
			}
		}
	}

	// Cleanup hook on exit
	if hookHandle != 0 {
		unhookEx.Call(uintptr(hookHandle))
	}
}
