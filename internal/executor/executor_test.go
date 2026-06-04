package executor

import (
	"strings"
	"testing"
)

// ── BuildCommand ──────────────────────────────────────────────────────────────

func TestBuildCommand_BasicSubstitution(t *testing.T) {
	out, err := BuildCommand("git commit -m {{.message}}", map[string]string{
		"message": "fix bug",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "git commit -m fix bug" {
		t.Errorf("got %q", out)
	}
}

func TestBuildCommand_MissingKeyYieldsZero(t *testing.T) {
	// missingkey=zero → missing key produces empty string, not an error
	out, err := BuildCommand("git push {{.remote}}", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "git push") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestBuildCommand_BoolTrue(t *testing.T) {
	out, err := BuildCommand(`git push{{if .force}} --force{{end}}`, map[string]string{
		"force": "true",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "--force") {
		t.Errorf("expected --force in output, got %q", out)
	}
}

func TestBuildCommand_BoolFalse(t *testing.T) {
	out, err := BuildCommand(`git push{{if .force}} --force{{end}}`, map[string]string{
		"force": "false",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "--force") {
		t.Errorf("--force should not appear when force=false, got %q", out)
	}
}

func TestBuildCommand_StringFalseIsNotBool(t *testing.T) {
	// "false" string → real bool false, not literal string "false"
	out, err := BuildCommand(`{{.val}}`, map[string]string{"val": "false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Go template renders bool false as "false"
	if out != "false" {
		t.Errorf("expected %q, got %q", "false", out)
	}
}

func TestBuildCommand_CollapseWhitespace(t *testing.T) {
	// Conditional that leaves extra spaces should be collapsed
	out, err := BuildCommand(`git commit  {{if .scope}}-s {{.scope}}{{end}}  -m {{.msg}}`, map[string]string{
		"scope": "",
		"msg":   "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not have double spaces
	if strings.Contains(out, "  ") {
		t.Errorf("output has double spaces: %q", out)
	}
}

func TestBuildCommand_InvalidTemplate(t *testing.T) {
	_, err := BuildCommand("git {{.unclosed", map[string]string{})
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestBuildCommand_MultipleArgs(t *testing.T) {
	out, err := BuildCommand("docker run -p {{.port}}:{{.port}} {{.image}}", map[string]string{
		"port":  "8080",
		"image": "nginx",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "docker run -p 8080:8080 nginx" {
		t.Errorf("got %q", out)
	}
}

// ── splitArgs ─────────────────────────────────────────────────────────────────

func TestSplitArgs_Simple(t *testing.T) {
	args, err := splitArgs("git commit -m hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"git", "commit", "-m", "hello"}
	assertStringSlice(t, args, want)
}

func TestSplitArgs_DoubleQuotedArg(t *testing.T) {
	args, err := splitArgs(`git commit -m "fix the bug"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d: %v", len(args), args)
	}
	if args[3] != "fix the bug" {
		t.Errorf("expected %q, got %q", "fix the bug", args[3])
	}
}

func TestSplitArgs_SingleQuotedArg(t *testing.T) {
	args, err := splitArgs("git commit -m 'fix the bug'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d: %v", len(args), args)
	}
	if args[3] != "fix the bug" {
		t.Errorf("expected %q, got %q", "fix the bug", args[3])
	}
}

func TestSplitArgs_UnclosedDoubleQuote(t *testing.T) {
	_, err := splitArgs(`git commit -m "unclosed`)
	if err == nil {
		t.Fatal("expected error for unclosed double quote")
	}
}

func TestSplitArgs_UnclosedSingleQuote(t *testing.T) {
	_, err := splitArgs("git commit -m 'unclosed")
	if err == nil {
		t.Fatal("expected error for unclosed single quote")
	}
}

func TestSplitArgs_Empty(t *testing.T) {
	args, err := splitArgs("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected empty args, got %v", args)
	}
}

func TestSplitArgs_OnlyWhitespace(t *testing.T) {
	args, err := splitArgs("   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected empty args, got %v", args)
	}
}

func TestSplitArgs_MultipleSpacesBetweenArgs(t *testing.T) {
	args, err := splitArgs("git   status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"git", "status"}
	assertStringSlice(t, args, want)
}

func TestSplitArgs_QuotedArgWithSpacesAtEnd(t *testing.T) {
	args, err := splitArgs(`echo "hello world" extra`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"echo", "hello world", "extra"}
	assertStringSlice(t, args, want)
}

func TestSplitArgs_FlagWithEquals(t *testing.T) {
	args, err := splitArgs("docker run --name=myapp nginx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"docker", "run", "--name=myapp", "nginx"}
	assertStringSlice(t, args, want)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func assertStringSlice(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("len mismatch: got %v, want %v", got, want)
		return
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] got %q, want %q", i, got[i], want[i])
		}
	}
}