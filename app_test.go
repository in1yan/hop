package main

import (
	"testing"
	"context"
)

// Test that NewApp creates a valid App instance
func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
}

// Test that App.startup properly sets the context
func TestAppStartup(t *testing.T) {
	app := NewApp()
	ctx := context.Background()
	
	// This should not panic
	app.startup(ctx)
	
	if app.ctx != ctx {
		t.Error("startup() did not properly set context")
	}
}

// Test thread-safe window visibility functions
func TestWindowVisibility(t *testing.T) {
	// Test initial state
	if getWindowVisibility() != false {
		t.Error("Initial window visibility should be false")
	}
	
	// Test setting visibility
	setWindowVisibility(true)
	if !getWindowVisibility() {
		t.Error("Window visibility should be true after setting to true")
	}
	
	setWindowVisibility(false)
	if getWindowVisibility() {
		t.Error("Window visibility should be false after setting to false")
	}
}

// Test concurrent access to window visibility (basic race condition test)
func TestWindowVisibilityConcurrency(t *testing.T) {
	done := make(chan bool)
	
	// Start multiple goroutines that toggle visibility
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				setWindowVisibility(j%2 == 0)
				getWindowVisibility()
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// If we get here without data races, the test passes
	t.Log("Concurrent window visibility operations completed successfully")
}