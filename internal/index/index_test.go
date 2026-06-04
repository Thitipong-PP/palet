package index

import (
	"testing"

	"github.com/Thitipong-PP/palet/internal/plugin"
)

func TestBuild_Empty(t *testing.T) {
	entries := Build(nil)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestBuild_SinglePlugin(t *testing.T) {
	plugins := []plugin.Plugin{
		{
			Name: "git",
			Commands: []plugin.Command{
				{Name: "status", Template: "git status"},
				{Name: "commit", Template: "git commit"},
			},
		},
	}
	entries := Build(plugins)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Plugin.Name != "git" {
			t.Errorf("expected plugin name %q, got %q", "git", e.Plugin.Name)
		}
	}
	if entries[0].Command.Name != "status" {
		t.Errorf("expected first command %q, got %q", "status", entries[0].Command.Name)
	}
}

func TestBuild_MultiplePlugins(t *testing.T) {
	plugins := []plugin.Plugin{
		{Name: "git", Commands: []plugin.Command{
			{Name: "status", Template: "git status"},
			{Name: "push", Template: "git push"},
		}},
		{Name: "docker", Commands: []plugin.Command{
			{Name: "ps", Template: "docker ps"},
		}},
		{Name: "npm", Commands: []plugin.Command{
			{Name: "install", Template: "npm install"},
			{Name: "run", Template: "npm run {{.script}}"},
			{Name: "test", Template: "npm test"},
		}},
	}

	entries := Build(plugins)
	if len(entries) != 6 {
		t.Fatalf("expected 6 entries, got %d", len(entries))
	}
}

func TestBuild_PreservesPluginReference(t *testing.T) {
	plugins := []plugin.Plugin{
		{
			Name:        "git",
			Description: "Git version control",
			Commands:    []plugin.Command{{Name: "log", Template: "git log"}},
		},
	}
	entries := Build(plugins)
	if entries[0].Plugin.Description != "Git version control" {
		t.Error("plugin description should be preserved in entry")
	}
}

func TestBuild_PreservesCommandArgs(t *testing.T) {
	plugins := []plugin.Plugin{{
		Name: "docker",
		Commands: []plugin.Command{{
			Name:     "run",
			Template: "docker run {{.image}}",
			Args:     []plugin.Arg{{Name: "image", Type: "string", Required: true}},
		}},
	}}
	entries := Build(plugins)
	if len(entries[0].Command.Args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(entries[0].Command.Args))
	}
}

func TestBuild_PluginWithNoCommands(t *testing.T) {
	plugins := []plugin.Plugin{
		{Name: "empty", Commands: nil},
		{Name: "git", Commands: []plugin.Command{{Name: "status", Template: "git status"}}},
	}
	entries := Build(plugins)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}