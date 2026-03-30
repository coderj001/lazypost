package config

import (
	"os"
	"path/filepath"

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
			"send-request":  "ctrl-s",
			"start-editor":  ":",
			"switch-method": "ctrl-m",
		},
	}
}
