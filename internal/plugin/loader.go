package plugin

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Dirs returns all directories to scan, in priority order:
//  1. ./plugins  (project-local)
//  2. ~/.config/palet/plugins  (user-global)
func Dirs() []string {
	var dirs []string
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(wd, "plugins"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, ".config", "palet", "plugins"))
	}
	return dirs
}

// LoadAll loads every plugin from every discovery dir.
// Dirs that don't exist are silently skipped.
func LoadAll() []Plugin {
	var out []Plugin
	for _, dir := range Dirs() {
		out = append(out, loadDir(dir)...)
	}
	return out
}

// loadDir reads all .yaml / .yml files in a directory.
func loadDir(dir string) []Plugin {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	out := make([]Plugin, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		p, err := parseFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue // skip malformed files
		}
		out = append(out, p)
	}
	return out
}

func parseFile(path string) (Plugin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Plugin{}, err
	}
	var p Plugin
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Plugin{}, err
	}
	return p, nil
}