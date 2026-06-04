package plugin

import (
	"runtime"
	"strings"
	"testing"
)

// ── Plugin.Validate ──────────────────────────────────────────────────────────

func TestValidate_ValidPlugin(t *testing.T) {
	p := Plugin{
		Name: "git",
		Commands: []Command{
			{Name: "commit", Template: "git commit -m {{.message}}",
				Args: []Arg{{Name: "message", Type: "string", Required: true}}},
		},
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_EmptyName(t *testing.T) {
	p := Plugin{Name: "", Commands: []Command{{Name: "cmd", Template: "echo"}}}
	err := p.Validate()
	assertError(t, err, "non-empty name")
}

func TestValidate_WhitespaceName(t *testing.T) {
	p := Plugin{Name: "   ", Commands: []Command{{Name: "cmd", Template: "echo"}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for whitespace-only name")
	}
}

func TestValidate_NoCommands(t *testing.T) {
	p := Plugin{Name: "git", Commands: nil}
	assertError(t, p.Validate(), "no commands")
}

func TestValidate_CommandEmptyName(t *testing.T) {
	p := Plugin{Name: "git", Commands: []Command{{Name: "", Template: "git status"}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty command name")
	}
}

func TestValidate_CommandEmptyTemplate(t *testing.T) {
	p := Plugin{Name: "git", Commands: []Command{{Name: "status", Template: ""}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty command template")
	}
}

func TestValidate_DuplicateArgName(t *testing.T) {
	p := Plugin{Name: "git", Commands: []Command{{
		Name:     "push",
		Template: "git push",
		Args: []Arg{
			{Name: "remote", Type: "string"},
			{Name: "remote", Type: "string"},
		},
	}}}
	assertError(t, p.Validate(), "duplicate arg name")
}

func TestValidate_UnknownArgType(t *testing.T) {
	p := Plugin{Name: "git", Commands: []Command{{
		Name:     "push",
		Template: "git push",
		Args:     []Arg{{Name: "remote", Type: "badtype"}},
	}}}
	assertError(t, p.Validate(), "unknown type")
}

func TestValidate_EnumWithoutChoices(t *testing.T) {
	p := Plugin{Name: "git", Commands: []Command{{
		Name:     "push",
		Template: "git push",
		Args:     []Arg{{Name: "strategy", Type: "enum"}},
	}}}
	assertError(t, p.Validate(), "at least one choice")
}

func TestValidate_EnumWithChoices(t *testing.T) {
	p := Plugin{Name: "docker", Commands: []Command{{
		Name:     "run",
		Template: "docker run",
		Args:     []Arg{{Name: "restart", Type: "enum", Choices: []string{"always", "no"}}},
	}}}
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_BoolInvalidDefault(t *testing.T) {
	p := Plugin{Name: "git", Commands: []Command{{
		Name:     "push",
		Template: "git push",
		Args:     []Arg{{Name: "force", Type: "bool", Default: "yes"}},
	}}}
	assertError(t, p.Validate(), "bool arg default")
}

func TestValidate_BoolValidDefaults(t *testing.T) {
	for _, def := range []string{"true", "false", ""} {
		p := Plugin{Name: "git", Commands: []Command{{
			Name:     "push",
			Template: "git push",
			Args:     []Arg{{Name: "force", Type: "bool", Default: def}},
		}}}
		if err := p.Validate(); err != nil {
			t.Errorf("default %q: unexpected error: %v", def, err)
		}
	}
}

func TestValidate_ArgEmptyName(t *testing.T) {
	p := Plugin{Name: "git", Commands: []Command{{
		Name:     "push",
		Template: "git push",
		Args:     []Arg{{Name: "", Type: "string"}},
	}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty arg name")
	}
}

func TestValidate_AllArgTypes(t *testing.T) {
	for _, typ := range []string{"string", "bool", "file", "dir"} {
		p := Plugin{Name: "x", Commands: []Command{{
			Name:     "cmd",
			Template: "echo",
			Args:     []Arg{{Name: "a", Type: typ}},
		}}}
		if err := p.Validate(); err != nil {
			t.Errorf("type %q: unexpected error: %v", typ, err)
		}
	}
}

// ── Plugin.MatchesOS ─────────────────────────────────────────────────────────

func TestMatchesOS_EmptyAlwaysMatches(t *testing.T) {
	p := Plugin{}
	if !p.MatchesOS() {
		t.Error("empty OS list should match all platforms")
	}
}

func TestMatchesOS_CurrentOS(t *testing.T) {
	p := Plugin{OS: []string{runtime.GOOS}}
	if !p.MatchesOS() {
		t.Errorf("plugin with OS=%q should match current OS %q", runtime.GOOS, runtime.GOOS)
	}
}

func TestMatchesOS_UnknownOS(t *testing.T) {
	p := Plugin{OS: []string{"__nonexistent_os__"}}
	if p.MatchesOS() {
		t.Error("plugin should not match current OS")
	}
}

func TestMatchesOS_MultipleOSCurrentIncluded(t *testing.T) {
	p := Plugin{OS: []string{"__nonexistent_os__", runtime.GOOS}}
	if !p.MatchesOS() {
		t.Error("plugin should match when current OS is in the list")
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func assertError(t *testing.T, err error, contains string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", contains)
	}
	if !strings.Contains(err.Error(), contains) {
		t.Errorf("error %q should contain %q", err.Error(), contains)
	}
}