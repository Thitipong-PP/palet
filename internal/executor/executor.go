package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
)

// BuildCommand renders a Go text/template against the collected arg values.
//
// Values that look boolean ("true"/"false") are passed to the template as real
// bools so conditional templates such as `{{if .new}}-b {{end}}` behave as
// expected instead of treating the literal string "false" as truthy.
func BuildCommand(tmpl string, vals map[string]string) (string, error) {
	t, err := template.New("cmd").Option("missingkey=zero").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	data := make(map[string]any, len(vals))
	for k, v := range vals {
		switch strings.TrimSpace(v) {
		case "true":
			data[k] = true
		case "false":
			data[k] = false
		default:
			data[k] = v
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}

	// Collapse stray whitespace produced by conditional blocks.
	out := strings.Join(strings.Fields(buf.String()), " ")
	return out, nil
}

// Run executes the built command through the user's shell, wiring stdio
// straight through so the caller's terminal receives the output.
func Run(cmd string) error {
    fields := strings.Fields(cmd)
    if len(fields) == 0 {
        return fmt.Errorf("nothing to run")
    }
    c := exec.Command(fields[0], fields[1:]...)
    c.Stdin = os.Stdin
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr
    return c.Run()
}

// Copy places the built command on the system clipboard.
func Copy(s string) error {
	name, args := clipboardCmd()
	if name == "" {
		return fmt.Errorf("no clipboard utility found")
	}
	c := exec.Command(name, args...)
	c.Stdin = strings.NewReader(s)
	if err := c.Run(); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}
	return nil
}

func shell() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	}
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh
	}
	return "sh"
}

// clipboardCmd picks the first available clipboard utility for the platform.
func clipboardCmd() (string, []string) {
	switch runtime.GOOS {
	case "darwin":
		return "pbcopy", nil
	case "windows":
		return "clip", nil
	default:
		for _, candidate := range [][]string{
			{"wl-copy"},
			{"xclip", "-selection", "clipboard"},
			{"xsel", "--clipboard", "--input"},
		} {
			if path, err := exec.LookPath(candidate[0]); err == nil {
				return path, candidate[1:]
			}
		}
		return "", nil
	}
}
