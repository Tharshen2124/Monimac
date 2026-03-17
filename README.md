# Monimac

Monimac is a terminal user interface (TUI) for monitoring and managing a Mac mini from a MacBook. It shows live CPU and memory usage alongside running Docker containers, and lets you stop containers directly from the terminal.

## Requirements

- Go 1.26 ([mise](https://mise.jdx.dev) is configured — run `mise install` to get the right version)
- Docker (optional — system metrics still work if Docker is not running)

## Setup

```bash
# Clone and enter the repo
git clone <repo-url>
cd monimac

# Install the correct Go version via mise
mise install

# Download dependencies
go mod download
```

## Commands

| Command | Description |
|---------|-------------|
| `go run .` | Run the TUI directly |
| `go build -o monimac .` | Compile a binary |
| `./monimac` | Run the compiled binary |
| `go mod tidy` | Sync dependencies after code changes |

## Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Select previous container |
| `↓` / `j` | Select next container |
| `enter` / `s` | Stop selected container (shows confirmation) |
| `y` | Confirm stop |
| `n` / `esc` | Cancel |
| `r` | Force refresh |
| `q` / `ctrl+c` | Quit |

## Features

- **CPU usage** — live percentage with a visual progress bar, refreshed every 3 seconds
- **Memory usage** — used / total with percentage bar
- **Docker containers** — name, image, status, CPU%, and memory usage per container
- **Stop containers** — select a container and press `s` or `enter`; a confirmation prompt prevents accidental stops
- **Docker-optional** — if Docker is not running or not installed, system metrics still display and the Docker section shows a clear error message

## Tech Stack

- Language: Go 1.26
- UI Framework: [Bubble Tea](https://charm.sh) by Charm
- Styling: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- System metrics: [gopsutil](https://github.com/shirou/gopsutil)

## Status

Initial working version. Currently monitors the local machine. SSH support for connecting to a remote Mac mini is planned.

## Notes

- This project focuses on operational visibility and fast incident response.
- Keep interactions simple, keyboard-first, and safe by default.