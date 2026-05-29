package main

import (
	"github.com/Thitipong-PP/palet/internal/index"
	"github.com/Thitipong-PP/palet/internal/plugin"
)

func main() {
	index.Build(plugin.LoadAll())
}