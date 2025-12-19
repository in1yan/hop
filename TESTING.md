# Testing Documentation

## Overview

This repository includes comprehensive test coverage for the Hop window switcher application. Due to the Windows-specific nature of the application, the test suite is designed to work across different platforms using build tags.

## Test Files

### Platform-Agnostic Tests

- **`common_test.go`**: Contains tests that can run on any platform
  - Tests for the `Greet()` function's string formatting
  - Tests for `Window` struct field operations
  - Tests for `App` struct creation
  - Tests for `KBDLLHOOKSTRUCT` struct fields
  - Tests for Windows API constants

### Windows-Specific Tests

- **`app_test.go`**: Contains integration tests for the App struct (Windows only)
- **`window_test.go`**: Contains integration tests for Window operations (Windows only)

### Stub Files

To enable testing on non-Windows platforms, stub implementations are provided:

- **`app_stub.go`**: Stub implementation of App methods for non-Windows platforms
- **`window_stub.go`**: Stub implementation of Window functions for non-Windows platforms

## Running Tests

### On Linux/macOS

```bash
go test -v .
```

This will run only the platform-agnostic tests from `common_test.go`.

### On Windows

```bash
go test -v .
```

This will run all tests including Windows-specific tests from `app_test.go` and `window_test.go`.

### Run with Coverage

```bash
# On any platform
go test -v -cover .

# Generate coverage report
go test -coverprofile=coverage.out .
go tool cover -html=coverage.out -o coverage.html
```

## Test Coverage

The test suite covers:

- ✅ String formatting in the `Greet()` function
- ✅ Window struct field access and manipulation
- ✅ App struct creation and initialization
- ✅ KBDLLHOOKSTRUCT struct fields
- ✅ Windows API constants (GA_ROOT, WM_KEYDOWN, etc.)
- ✅ Slice operations with Window structs
- ✅ Unicode character handling

### Windows-Specific Coverage (when running on Windows)

- Window enumeration
- Window focus management
- Keyboard hook installation
- Tray icon functionality

## Build Tags

The project uses Go build tags to ensure platform-specific code only compiles on appropriate platforms:

- `//go:build windows` - Code that only builds on Windows
- `//go:build !windows` - Code that builds on non-Windows platforms (stubs)

## Adding New Tests

When adding new tests:

1. **For platform-agnostic logic**: Add tests to `common_test.go`
2. **For Windows-specific features**: Add tests to the appropriate Windows-specific test file with the `//go:build windows` tag
3. **For new stub functions**: Update the stub files (`*_stub.go`) to match the Windows implementation signatures

## Continuous Integration

For CI/CD pipelines:

- **Linux/macOS runners**: Will run platform-agnostic tests
- **Windows runners**: Will run the full test suite including Windows-specific tests

Example GitHub Actions configuration:

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - run: go test -v -cover ./...
```

## Notes

- The application is primarily designed for Windows and uses Windows-specific APIs
- Stub implementations allow tests to compile and run on any platform
- Full integration testing requires a Windows environment
- The test suite focuses on testable business logic and data structures
