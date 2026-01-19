# go-todo

Minimal TUI todo app with vim-style navigation.

## Install

```bash
go install github.com/tungsheng/go-todo/cmd/todo@latest
```

Or build from source:

```bash
make build
./bin/todo
```

## Usage

```
╭──────────╮
│ go-todo  │ today
╰──────────╯

  ○ Buy groceries
> ◐ Finish report
  ● Call mom
  ── closed items ──
  ✕ Old task

n:new  e:edit  d:delete  s:status  t:tag  x:close  q:quit
```

## Keys

| Key | Action |
|-----|--------|
| `j/k` | Navigate |
| `n` | New todo |
| `e` | Edit |
| `d` | Delete |
| `s/space` | Cycle status (○ → ◐ → ●) |
| `x` | Toggle closed |
| `t` | Cycle tag (today/week/month) |
| `q` | Quit |

## Storage

`~/.go-todo/todos.db` (libsql/SQLite)

## License

MIT
