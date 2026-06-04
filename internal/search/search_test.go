package search

import (
	"testing"

	"github.com/Thitipong-PP/palet/internal/index"
	"github.com/Thitipong-PP/palet/internal/plugin"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func makeEntry(pluginName, cmdName, cmdDesc string) index.Entry {
	return index.Entry{
		Plugin:  plugin.Plugin{Name: pluginName},
		Command: plugin.Command{Name: cmdName, Description: cmdDesc},
	}
}

func entryNames(entries []index.Entry) []string {
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Command.Name
	}
	return names
}

// ── Query ─────────────────────────────────────────────────────────────────────

func TestQuery_EmptyQueryReturnsAll(t *testing.T) {
	entries := []index.Entry{
		makeEntry("git", "status", "Show status"),
		makeEntry("git", "commit", "Commit changes"),
		makeEntry("docker", "ps", "List containers"),
	}
	result := Query(entries, "")
	if len(result) != 3 {
		t.Errorf("empty query should return all %d entries, got %d", len(entries), len(result))
	}
}

func TestQuery_NoMatch(t *testing.T) {
	entries := []index.Entry{
		makeEntry("git", "status", "Show status"),
		makeEntry("git", "commit", "Commit changes"),
	}
	result := Query(entries, "xxxxxxxxxx")
	if len(result) != 0 {
		t.Errorf("expected 0 results, got %d", len(result))
	}
}

func TestQuery_ExactPrefixMatchRanksFirst(t *testing.T) {
	entries := []index.Entry{
		makeEntry("git", "status", "Show working tree status"),
		makeEntry("git", "stash", "Stash changes"),        // contains "sta" in name too
		makeEntry("git", "log", "Show commit log"),
	}
	result := Query(entries, "sta")
	if len(result) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(result))
	}
	// "stash" starts with "sta" → exact prefix → highest score
	if result[0].Command.Name != "stash" && result[0].Command.Name != "status" {
		t.Errorf("first result should be exact-prefix match, got %q", result[0].Command.Name)
	}
}

func TestQuery_CaseInsensitive(t *testing.T) {
	entries := []index.Entry{makeEntry("git", "STATUS", "Show status")}
	result := Query(entries, "status")
	if len(result) != 1 {
		t.Errorf("expected 1 result for case-insensitive match, got %d", len(result))
	}
}

func TestQuery_DescriptionMatch(t *testing.T) {
	entries := []index.Entry{
		makeEntry("git", "xyz123", "Show working tree status"),
	}
	result := Query(entries, "working tree")
	if len(result) != 1 {
		t.Errorf("expected 1 result for description match, got %d", len(result))
	}
}

func TestQuery_PluginNameMatch(t *testing.T) {
	entries := []index.Entry{
		makeEntry("docker", "ps", "List containers"),
		makeEntry("git", "status", "Show status"),
	}
	result := Query(entries, "docker")
	if len(result) != 1 {
		t.Errorf("expected 1 result for plugin name match, got %d", len(result))
	}
	if result[0].Plugin.Name != "docker" {
		t.Errorf("expected docker plugin, got %q", result[0].Plugin.Name)
	}
}

func TestQuery_FuzzyMatch(t *testing.T) {
	entries := []index.Entry{
		makeEntry("git", "commit", "Commit staged changes"),
	}
	// "cmt" is a fuzzy match for "commit"
	result := Query(entries, "cmt")
	if len(result) != 1 {
		t.Errorf("expected fuzzy match, got %d results", len(result))
	}
}

func TestQuery_ScoreOrdering(t *testing.T) {
	// "push" exactly starts with "push" → scoreExactPrefix
	// something that only contains "push" in desc → scoreDescContains
	entries := []index.Entry{
		makeEntry("git", "log", "push changes upstream"),      // desc contains "push"
		makeEntry("git", "push", "Upload commits to remote"),  // name exact prefix
	}
	result := Query(entries, "push")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0].Command.Name != "push" {
		t.Errorf("exact-prefix match 'push' should rank first, got %q", result[0].Command.Name)
	}
}

func TestQuery_WhitespaceOnlyQuery(t *testing.T) {
	entries := []index.Entry{makeEntry("git", "status", "Show status")}
	// after TrimSpace, empty → returns all
	result := Query(entries, "   ")
	if len(result) != 1 {
		t.Errorf("whitespace-only query should return all entries, got %d", len(result))
	}
}

// ── fuzzy ────────────────────────────────────────────────────────────────────

func TestFuzzy(t *testing.T) {
	tests := []struct {
		s, pattern string
		want       bool
	}{
		{"commit", "cmt", true},
		{"commit", "commit", true},
		{"commit", "commits", false},  // pattern longer than s
		{"", "", true},                // empty pattern always matches
		{"abc", "", true},
		{"abc", "abcd", false},
		{"status", "sts", true},
		{"status", "xyz", false},
	}
	for _, tt := range tests {
		got := fuzzy(tt.s, tt.pattern)
		if got != tt.want {
			t.Errorf("fuzzy(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
		}
	}
}