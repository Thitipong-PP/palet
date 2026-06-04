package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

// ── ParseFile ────────────────────────────────────────────────────────────────

func TestParseFile_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	content := `
name: git
description: Git commands
commands:
  - name: status
    description: Show status
    template: git status
`
	path := filepath.Join(dir, "git.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	p, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "git" {
		t.Errorf("expected name %q, got %q", "git", p.Name)
	}
	if len(p.Commands) != 1 {
		t.Errorf("expected 1 command, got %d", len(p.Commands))
	}
}

func TestParseFile_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte(":::invalid: yaml:::"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseFile(path); err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	if _, err := ParseFile("/nonexistent/path/plugin.yaml"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseFile_FailsValidation(t *testing.T) {
	dir := t.TempDir()
	// Valid YAML but no commands → fails Validate()
	content := `
name: empty-plugin
description: No commands
`
	path := filepath.Join(dir, "empty.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseFile(path); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestParseFile_WithArgs(t *testing.T) {
	dir := t.TempDir()
	content := `
name: docker
commands:
  - name: run
    template: docker run {{.image}}
    args:
      - name: image
        type: string
        required: true
`
	path := filepath.Join(dir, "docker.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	p, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Commands[0].Args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(p.Commands[0].Args))
	}
	if !p.Commands[0].Args[0].Required {
		t.Error("expected arg to be required")
	}
}

// ── loadDir ──────────────────────────────────────────────────────────────────

func TestLoadDir_NonExistentDirReturnsEmpty(t *testing.T) {
	result := loadDir("/nonexistent/dir/12345")
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d plugins", len(result))
	}
}

func TestLoadDir_SkipsNonYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("not a plugin"), 0600)
	os.WriteFile(filepath.Join(dir, "data.json"), []byte(`{}`), 0600)

	result := loadDir(dir)
	if len(result) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(result))
	}
}

func TestLoadDir_LoadsYAMLFiles(t *testing.T) {
	dir := t.TempDir()
	yaml1 := `name: p1
commands:
  - name: cmd1
    template: echo 1`
	yaml2 := `name: p2
commands:
  - name: cmd2
    template: echo 2`

	os.WriteFile(filepath.Join(dir, "p1.yaml"), []byte(yaml1), 0600)
	os.WriteFile(filepath.Join(dir, "p2.yml"), []byte(yaml2), 0600)

	result := loadDir(dir)
	if len(result) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(result))
	}
}

func TestLoadDir_SkipsSubdirectories(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "subdir"), 0700)
	os.WriteFile(filepath.Join(dir, "subdir", "plugin.yaml"), []byte(`name: sub
commands:
  - name: cmd
    template: echo`), 0600)

	result := loadDir(dir)
	if len(result) != 0 {
		t.Errorf("expected 0 plugins (subdir should be skipped), got %d", len(result))
	}
}

func TestLoadDir_SkipsInvalidPlugins(t *testing.T) {
	dir := t.TempDir()
	// Valid
	os.WriteFile(filepath.Join(dir, "good.yaml"), []byte(`name: good
commands:
  - name: cmd
    template: echo`), 0600)
	// Invalid (no commands)
	os.WriteFile(filepath.Join(dir, "bad.yaml"), []byte(`name: bad`), 0600)

	result := loadDir(dir)
	if len(result) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(result))
	}
	if result[0].Name != "good" {
		t.Errorf("expected plugin %q, got %q", "good", result[0].Name)
	}
}