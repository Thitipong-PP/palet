package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/Thitipong-PP/palet/internal/index"
	"github.com/Thitipong-PP/palet/internal/plugin"
	"github.com/Thitipong-PP/palet/internal/tui"
)

//go:embed plugins/*.yaml
var embeddedPlugins embed.FS

func main() {
	// Set embedded plugins in the loader
	plugin.SetEmbeddedPlugins(embeddedPlugins)
	
	entries := index.LoadCached()
	if err := tui.Start(entries); err != nil {
		fmt.Fprintln(os.Stderr, "palet:", err)
		os.Exit(1)
	}
}
