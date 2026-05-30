package plugin

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// embeddedPlugins is injected from main.go
var embeddedPlugins embed.FS

// SetEmbeddedPlugins sets the embedded file system to use for plugin loading
func SetEmbeddedPlugins(fs embed.FS) {
	embeddedPlugins = fs
}

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

// LoadAll loads every plugin from every discovery dir plus embedded plugins.
// Dirs that don't exist are silently skipped.
func LoadAll() []Plugin {
	var out []Plugin
	
	// First load from filesystem directories
	for _, dir := range Dirs() {
		out = append(out, loadDir(dir)...)
	}
	
	// Then load from embedded plugins
	out = append(out, LoadEmbedded()...)
	
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
		p, err := ParseFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		out = append(out, p)
	}
	return out
}

// LoadEmbedded loads all .yaml files from the embedded plugins directory.
// This is exported so cache.go can call it.
func LoadEmbedded() []Plugin {
	entries, err := embeddedPlugins.ReadDir("plugins")
	if err != nil {
		// If embedded plugins are not set or empty, return empty slice
		return nil
	}
	
	var out []Plugin
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		
		data, err := embeddedPlugins.ReadFile(path.Join("plugins", e.Name()))
		if err != nil {
			continue
		}
		
		var p Plugin
		if err := yaml.Unmarshal(data, &p); err != nil {
			continue
		}
		if err := p.Validate(); err != nil {
			continue
		}
		out = append(out, p)
	}
	
	return out
}

// ParseFile parses and validates a single YAML plugin file.
func ParseFile(path string) (Plugin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Plugin{}, err
	}
	var p Plugin
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Plugin{}, err
	}
	if err := p.Validate(); err != nil {
		return Plugin{}, fmt.Errorf("%s: %w", path, err)
	}
	return p, nil
}