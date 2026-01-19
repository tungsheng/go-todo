# Developer Guide

## Structure

```
go-todo/
├── cmd/todo/main.go          # Entry point
├── internal/
│   ├── model/todo.go         # Todo, Status, TimeTag
│   ├── storage/libsql.go     # Database operations
│   └── tui/
│       ├── app.go            # TUI (bubbletea)
│       └── styles.go         # Styles (lipgloss)
├── Makefile
└── README.md
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `charmbracelet/bubbletea` | TUI framework |
| `charmbracelet/bubbles` | Input components |
| `charmbracelet/lipgloss` | Styling |
| `tursodatabase/go-libsql` | Database |

## Database

```sql
CREATE TABLE todos (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT NOT NULL,
    status     TEXT DEFAULT 'pending',
    time_tag   TEXT DEFAULT '',
    created_at INTEGER,
    updated_at INTEGER
);
```

## Build

```bash
make build   # Build binary
make run     # Build and run
make clean   # Remove binary
make tidy    # Tidy dependencies
```

## Architecture

### Model

- `Status`: pending, in_progress, done, closed
  - `Icon()` → ○ ◐ ● ✕
  - `Next()` → cycles pending/in_progress/done
  - `ToggleClosed()` → toggles closed state
- `TimeTag`: today, week, month
  - `Label()` → display text

### Storage

- `New()` → init database
- `ListFiltered(tag)` → list todos by time tag
- `Create(title, tag)` → create todo
- `Update(todo)` → update todo
- `Delete(id)` → delete todo

### TUI

Elm architecture: Model → Update → View

Modes: list, new, edit, confirm delete
