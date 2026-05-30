package index

import (
	"github.com/Thitipong-PP/palet/internal/plugin"
)

// Entry is one searchable, flat command record.
type Entry struct {
	Plugin  plugin.Plugin
	Command plugin.Command
}

// Build flattens a slice of plugins into a single Entry slice.
// Allocation is done once up front.
func Build(plugins []plugin.Plugin) []Entry {
	total := 0
	for i := range plugins {
		total += len(plugins[i].Commands)
	}
	out := make([]Entry, 0, total)
	for _, p := range plugins {
		for _, c := range p.Commands {
			out = append(out, Entry{Plugin: p, Command: c})
		}
	}
	return out
}