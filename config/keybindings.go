package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jroimartin/gocui"
	"gopkg.in/yaml.v3"
)

type KeybindingConfig struct {
	Keybindings map[string]string `yaml:"keybindings"`
}

func LoadKeybindings() (*KeybindingConfig, error) {
	defaults := DefaultKeybindings()

	configPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".config", "lazypost", "keybindings.yaml"),
		"./keybindings.yaml",
	}

	for _, path := range configPaths {
		data, err := os.ReadFile(path)
		if err == nil {
			var cfg KeybindingConfig
			if err := yaml.Unmarshal(data, &cfg); err == nil {
				return &cfg, nil
			}
		}
	}

	return defaults, nil
}

func DefaultKeybindings() *KeybindingConfig {
	return &KeybindingConfig{
		Keybindings: map[string]string{
			"quit":          "ctrl-c",
			"quit-alt":      "q",
			"next-view":     "tab",
			"prev-view":     "[",
			"send-request":  "ctrl-s",
			"start-editor":  ":",
			"switch-method": "ctrl-m",
		},
	}
}

func ParseKey(keyStr string) (gocui.Key, gocui.Modifier, error) {
	keyStr = normalizeKeyString(keyStr)

	switch keyStr {
	case "ctrl-c":
		return gocui.KeyCtrlC, gocui.ModNone, nil
	case "ctrl-s":
		return gocui.KeyCtrlS, gocui.ModNone, nil
	case "ctrl-m":
		return gocui.KeyCtrlM, gocui.ModNone, nil
	case "tab":
		return gocui.KeyTab, gocui.ModNone, nil
	case "enter":
		return gocui.KeyEnter, gocui.ModNone, nil
	case "esc", "escape":
		return gocui.KeyEsc, gocui.ModNone, nil
	case "q":
		return 'q', gocui.ModNone, nil
	case ":":
		return ':', gocui.ModNone, nil
	case "[":
		return '[', gocui.ModNone, nil
	case "]":
		return ']', gocui.ModNone, nil
	default:
		return 0, 0, fmt.Errorf("unknown key: %s", keyStr)
	}
}

func normalizeKeyString(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

var validActions = map[string]bool{
	"quit":          true,
	"quit-alt":      true,
	"next-view":     true,
	"prev-view":     true,
	"send-request":  true,
	"start-editor":  true,
	"switch-method": true,
}

func (c *KeybindingConfig) Validate() []string {
	var warnings []string
	for action, keyStr := range c.Keybindings {
		if !validActions[action] {
			warnings = append(warnings, fmt.Sprintf("unknown action '%s', skipping", action))
			continue
		}
		if _, _, err := ParseKey(keyStr); err != nil {
			warnings = append(warnings, fmt.Sprintf("invalid key '%s' for action '%s': %v", keyStr, action, err))
		}
	}
	return warnings
}
