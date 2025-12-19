//go:build windows
// +build windows

package main

import (
	"testing"
)

func TestWindowStruct(t *testing.T) {
	tests := []struct {
		name             string
		window           Window
		expectZeroHandle bool
	}{
		{
			name: "Basic window",
			window: Window{
				Handle: 12345,
				Title:  "Test Window",
			},
			expectZeroHandle: false,
		},
		{
			name: "Window with zero handle",
			window: Window{
				Handle: 0,
				Title:  "Zero Handle Window",
			},
			expectZeroHandle: true,
		},
		{
			name: "Window with empty title",
			window: Window{
				Handle: 54321,
				Title:  "",
			},
			expectZeroHandle: false,
		},
		{
			name: "Window with long title",
			window: Window{
				Handle: 99999,
				Title:  "This is a very long window title that might be used in real applications",
			},
			expectZeroHandle: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the struct can be created and accessed
			if tt.expectZeroHandle && tt.window.Handle != 0 {
				t.Errorf("Expected zero handle, got %d", tt.window.Handle)
			}
			if !tt.expectZeroHandle && tt.window.Handle == 0 {
				t.Errorf("Expected non-zero handle for %s", tt.name)
			}
			
			// Test that title is accessible
			_ = tt.window.Title
		})
	}
}

func TestWindowSlice(t *testing.T) {
	windows := []Window{
		{Handle: 1, Title: "Window 1"},
		{Handle: 2, Title: "Window 2"},
		{Handle: 3, Title: "Window 3"},
	}
	
	if len(windows) != 3 {
		t.Errorf("Expected 3 windows, got %d", len(windows))
	}
	
	// Test accessing individual windows
	if windows[0].Title != "Window 1" {
		t.Errorf("Expected 'Window 1', got %s", windows[0].Title)
	}
	
	if windows[1].Handle != 2 {
		t.Errorf("Expected handle 2, got %d", windows[1].Handle)
	}
}
