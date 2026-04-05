# pomodoro-tui

A terminal UI Pomodoro / focus-session tracker built in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea). It lets you browse your projects and tasks, run timed focus sessions, and persist every session to an external API (designed for [n8n](https://n8n.io/) webhooks).

---

## Features

- **Project browser** — lists your projects sorted by the time of the most recent task update.
- **Task browser** — shows the current week's tasks for the selected project, sorted by completion status, priority, and deadline. Deadlines are displayed in human-readable relative time ("in 2 days", "yesterday").
- **Task creation** — add a new task directly from the terminal with name, priority (`!` / `!!` / `!!!`), deadline, estimated duration, and a free-text description.
- **Toggle task completion** — mark tasks done/undone without leaving the TUI.
- **Focus session selection** — choose a session type (e.g. *work 45 min*, *break 5 min*) from a configurable list. A side panel shows the task description and estimated duration.
- **Countdown timer** — counts down the selected session duration. The session is created in the API the moment the timer starts.
- **Auto-stop & persist** — when the timer reaches zero (or you press `b` to exit early), the session end time is recorded and pushed to the API. The elapsed time is capped at the session duration.
- **Desktop notifications** — a system notification and a beep sound are triggered when a session finishes.

---

## Navigation

| Screen | Key | Action |
|---|---|---|
| Any | `ctrl+c` | Quit |
| Projects | `↑ / ↓` | Move cursor |
| Projects | `enter` | Open tasks for selected project |
| Tasks | `↑ / ↓` | Move cursor |
| Tasks | `enter` | Open session picker for selected task |
| Tasks | `a` | Open new-task form |
| Tasks | `c` or `tab` | Toggle completion of selected task |
| Tasks | `b` or `backspace` | Back to projects |
| Sessions | `↑ / ↓` | Move cursor |
| Sessions | `enter` | Start timer |
| Sessions | `b` or `backspace` | Back to tasks |
| Timer | `b` or `backspace` | Stop timer and save session |
| New task | `tab` / `shift+tab` | Next / previous field |
| New task | `alt+enter` | Save task and return to task list |
| New task | `esc` | Cancel and return to task list |

---

## Requirements

- **Go 1.22+** (the module uses `go 1.25.7` in `go.mod`; any recent toolchain works)
- A running **n8n** instance (or any HTTP API) that implements the expected webhook endpoints
- Linux: no extra dependencies (desktop notifications use DBus / libnotify)
- Windows: no extra dependencies (notifications use the Windows Toast API)

---

## Configuration

Copy the example config and fill in your values:

```bash
cp example-config.json config.json
```

`config.json` is read from the working directory at startup. It is excluded from version control — never commit it, as it contains your API key.

### Fields

| Field | Type | Description |
|---|---|---|
| `api_url` | string | Webhook API base URL |
| `api_key` | string | API authentication key |
| `device_name` | string | Label identifying this machine in session records |
| `session_types` | array | List of focus session types shown in the session picker |
| `session_types[].name` | string | Display name of the session type |
| `session_types[].duration_minutes` | integer | Session length in minutes (must be > 0) |

### Example

```json
{
  "api_url": "https://your-n8n-instance.example/webhook/projects",
  "api_key": "your-secret-key",
  "device_name": "my-laptop",
  "session_types": [
    {"name": "work",        "duration_minutes": 45},
    {"name": "short work",  "duration_minutes": 25},
    {"name": "break",       "duration_minutes": 5},
    {"name": "long break",  "duration_minutes": 15}
  ]
}
```

All fields are required. The app exits with an error if any field is missing, the file is absent, or the JSON is malformed.

---

## Installation & Running

### Linux

```bash
# Clone the repo
git clone https://github.com/your-user/pomodoro-tui.git
cd pomodoro-tui

# Install dependencies
go mod download

# Copy and edit the config
cp example-config.json config.json
$EDITOR config.json

# Run directly
go run .

# Or build a binary and run it
go build -o pomodoro-tui .
./pomodoro-tui
```

### Windows

```powershell
# Clone the repo
git clone https://github.com/your-user/pomodoro-tui.git
cd pomodoro-tui

# Install dependencies
go mod download

# Copy and edit the config
copy example-config.json config.json
notepad config.json

# Run directly
go run .

# Or build a binary and run it
go build -o pomodoro-tui.exe .
.\pomodoro-tui.exe
```

> **Note for Windows users:** run the binary inside **Windows Terminal** or any modern terminal emulator that supports ANSI escape codes. The legacy `cmd.exe` window may not render colours and box-drawing characters correctly.

---

## Project Structure

```
pomodoro-tui/
├── main.go            # Entry point — loads config, wires client, starts Bubble Tea
├── model.go           # Central model struct and Init/Update/View dispatch
├── projects.go        # Projects screen
├── tasks.go           # Tasks screen
├── sessions.go        # Session-type picker screen
├── timer.go           # Countdown timer screen
├── newTask.go         # New-task form screen
├── newTaskForm.go     # Form struct and field definitions
├── view.go            # Shared styles (lipgloss)
├── item.go            # list.Item adapter
├── utils.go           # Generic Map helper
├── apiclient/         # HTTP client, domain types, request/response DTOs
└── config/            # Config file loading and validation
```

---

## License

See [LICENSE](LICENSE).

## LLM Usage
LLMs and code-assist tools have been used to build this project.
- Claude Haiku/Sonnet 4.5
- Charm's Crush
