package search

import (
	"strings"

	"github.com/Thitipong-PP/palet/internal/index"
)

const (
	scoreExactPrefix  = 100
	scorePluginPrefix = 90
	scoreNameContains = 70
	scoreDescContains = 50
	scorePluginName   = 40
	scoreFuzzy        = 20
)

// Query returns entries matching query, ranked by relevance.
// An empty query returns all entries unchanged.
func Query(entries []index.Entry, query string) []index.Entry {
	if query == "" {
		return entries
	}
	q := strings.ToLower(strings.TrimSpace(query))

	type hit struct {
		entry index.Entry
		score int
	}

	hits := make([]hit, 0, len(entries))
	for _, e := range entries {
		if s := score(e, q); s > 0 {
			hits = append(hits, hit{e, s})
		}
	}

	// insertion sort — list is small (hundreds of commands at most)
	for i := 1; i < len(hits); i++ {
		for j := i; j > 0 && hits[j].score > hits[j-1].score; j-- {
			hits[j], hits[j-1] = hits[j-1], hits[j]
		}
	}

	out := make([]index.Entry, len(hits))
	for i, h := range hits {
		out[i] = h.entry
	}
	return out
}

func score(e index.Entry, q string) int {
	name   := strings.ToLower(e.Command.Name)
	desc   := strings.ToLower(e.Command.Description)
	plugin := strings.ToLower(e.Plugin.Name)

	switch {
	case strings.HasPrefix(name, q):
		return scoreExactPrefix
	case strings.HasPrefix(plugin+name, q):
		return scorePluginPrefix
	case strings.Contains(name, q):
		return scoreNameContains
	case strings.Contains(desc, q):
		return scoreDescContains
	case strings.Contains(plugin, q):
		return scorePluginName
	case fuzzy(name, q) || fuzzy(desc, q):
		return scoreFuzzy
	}
	return 0
}

// fuzzy returns true if every rune of pattern appears in s in order.
func fuzzy(s, pattern string) bool {
	pi := 0
	for _, ch := range s {
		if pi < len(pattern) && ch == rune(pattern[pi]) {
			pi++
		}
	}
	return pi == len(pattern)
}