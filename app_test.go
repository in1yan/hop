//go:build windows
// +build windows

package main

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
	if app.ctx != nil {
		t.Error("Expected ctx to be nil before startup")
	}
}

func TestGreet(t *testing.T) {
	app := NewApp()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple name",
			input:    "World",
			expected: "Hello World, It's show time!",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "Hello , It's show time!",
		},
		{
			name:     "Name with spaces",
			input:    "John Doe",
			expected: "Hello John Doe, It's show time!",
		},
		{
			name:     "Special characters",
			input:    "Test@123",
			expected: "Hello Test@123, It's show time!",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.Greet(tt.input)
			if result != tt.expected {
				t.Errorf("Greet(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
