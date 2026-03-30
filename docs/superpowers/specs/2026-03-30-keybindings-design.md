# Keybinding Configuration Design

## Overview

Add YAML-based keybinding configuration to LazyPost, allowing users to customize keyboard shortcuts via a config file.

## Config File Location & Priority

1. `~/.config/lazypost/keybindings.yaml` (user-specific)
2. `./keybindings.yaml` (project/local override)

First found wins. If neither exists, use built-in defaults.

## YAML Format

```yaml
keybindings:
  quit: ctrl-c
  quit-alt: q
  next-view: tab
  send-request: ctrl-s
  start-editor: ':'
  switch-method: ctrl-m
```

## Supported Actions

| Action | Description |
|--------|-------------|
| `quit` | Exit application (primary) |
| `quit-alt` | Exit application (alternate) |
| `next-view` | Cycle to next view |
| `send-request` | Execute HTTP request |
| `start-editor` | Open floating editor |
| `switch-method` | Cycle through HTTP methods |

## Validation

- On startup, validate all keybindings in config
- Invalid key syntax or unknown action: log warning with file:line, skip that entry, continue loading
- Missing required actions: use built-in defaults

## Implementation Notes

- Use `gopkg.in/yaml.v3` for YAML parsing
- Key string format must match gocui key notation (e.g., `ctrl-c`, `tab`, `q`)
- Config is loaded once at startup; no hot-reload

## Built-in Defaults

```yaml
quit: ctrl-c
quit-alt: q
next-view: tab
send-request: ctrl-s
start-editor: ':'
switch-method: ctrl-m
```
