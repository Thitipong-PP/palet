package tui

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	HiddenPlugins map[string]bool `json:"hidden_plugins"`
}

func settingsPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "palet", "settings.json")
	}
	return filepath.Join(os.TempDir(), "palet-settings.json")
}

func loadSettings() map[string]bool {
	data, err := os.ReadFile(settingsPath())
	if err != nil {
		return make(map[string]bool)
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil || s.HiddenPlugins == nil {
		return make(map[string]bool)
	}
	return s.HiddenPlugins
}

func saveSettings(hiddenPlugins map[string]bool) error {
	path := settingsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	
	s := Settings{HiddenPlugins: hiddenPlugins}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
