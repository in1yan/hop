# Hop 🚀

A fast, lightweight window switcher for Windows that lets you quickly switch between open applications using keyboard shortcuts and visual hints.

<img width="1919" height="1079" alt="image" src="https://github.com/user-attachments/assets/3a3c671f-87a5-41c8-b133-17a9c565db61" />


## Features

- 🔥 **Lightning Fast**: Instantly switch between windows with keyboard shortcuts
- 🎯 **Visual Hints**: Type hint letters to quickly select any window
- 🌟 **Modern UI**: Clean, translucent interface with Acrylic blur effects
- ⚡ **Global Hotkey**: Toggle with Right Shift key from anywhere
- 🪶 **Lightweight**: Minimal resource usage, runs in the background
- 🎨 **Customizable**: Rose Pine color scheme with transparency support

## How It Works

1. **Activate**: Press the **Right Shift** key to open the window switcher
2. **Navigate**: See all open windows with letter hints (a, b, c, etc.)
3. **Select**: Type the hint letters to instantly switch to that window
4. **Toggle**:  Press **Right Shift** again to toggle
5. **Close**: Press **Escape** to close (Exit the Process)

## Installation

### Download Pre-built Binary
1. Download the latest release from the [Releases](../../releases) page
2. Extract the executable to your preferred location
3. Run `hop.exe` - it will start minimized in the background
4. Use **Right Shift** to activate the window switcher

### Build from Source

#### Prerequisites
- [Go](https://golang.org/dl/) (version 1.23 or later)
- [Node.js](https://nodejs.org/) (version 16 or later)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

#### Build Steps
```bash
# Clone the repository
git clone https://github.com/in1yan/hop.git
cd hop

# Install Wails (if not already installed)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build the application
wails build
```

The built executable will be available in the `build/bin/` directory.

## Development

### Live Development
```bash
# Start the development server with hot reload
wails dev
```

This will start a Vite development server for fast frontend development and a development version of the app.

### Project Structure
```
hop/
├── app.go              # Main application logic
├── window.go           # Windows API integration and hotkey handling
├── main.go            # Application entry point and configuration
├── frontend/          # React frontend
│   ├── src/
│   │   ├── App.jsx    # Main UI component
│   │   └── App.css    # Styling
│   └── package.json   # Frontend dependencies
├── build/             # Build assets and output
└── wails.json         # Wails configuration
```

### Key Components

- **Global Hotkey**: Right Shift key detection using Windows API hooks
- **Window Enumeration**: Fetches all visible windows using Windows API
- **Hint Generation**: Creates alphabetical hints for quick selection
- **Transparent UI**: Frameless window with Acrylic blur effects

## Configuration

The application can be configured by editing `wails.json`:

```json
{
  "name": "hop",
  "outputfilename": "hop",
  "author": {
    "name": "iniyan",
    "email": "viniyan563@gmail.com"
  }
}
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| **Right Shift** | Toggle window switcher |
| **a-z** | Type hint letters to select window |
| **Escape** | Close the application |

## Technical Details

- **Backend**: Go with Windows API integration
- **Frontend**: React with Vite for fast development
- **Framework**: Wails v2 for native desktop app development
- **Platform**: Windows (uses Windows-specific APIs)
- **UI Effects**: Translucent window with Acrylic backdrop

## System Requirements

- Windows 10 or later
- 64-bit architecture
- Minimal RAM and CPU usage

## Testing

The project includes a comprehensive test suite that works across different platforms. For detailed information about running tests, see [TESTING.md](TESTING.md).

```bash
# Run tests on any platform
go test -v .

# Run tests with coverage
go test -v -cover .
```

The test suite includes:
- Platform-agnostic unit tests (run on any OS)
- Windows-specific integration tests (run only on Windows)
- Stub implementations for non-Windows platforms

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

When contributing, please ensure that:
- Your code includes appropriate tests
- All existing tests pass
- You follow the existing code style

## Author

Created by **iniyan**

---

*Built with ❤️ using [Wails](https://wails.io/)*
