package plugin

import (
	"fmt"
	"runtime"
	"strings"
)

type Plugin struct {
	Name        string    `yaml:"name"        json:"name"`
	Description string    `yaml:"description" json:"description"`
	Commands    []Command `yaml:"commands"    json:"commands"`
	OS          []string  `yaml:"os"          json:"os"` // e.g. [linux, darwin, windows] — empty means all
	Hidden      bool      `yaml:"hidden"      json:"hidden"`
}

type Command struct {
	Name        string `yaml:"name"        json:"name"`
	Description string `yaml:"description" json:"description"`
	Template    string `yaml:"template"    json:"template"`
	Args        []Arg  `yaml:"args"        json:"args"`
}

type Arg struct {
	Name        string   `yaml:"name"        json:"name"`
	Description string   `yaml:"description" json:"description"`
	Type        string   `yaml:"type"        json:"type"`
	Required    bool     `yaml:"required"    json:"required"`
	Default     string   `yaml:"default"     json:"default"`
	Choices     []string `yaml:"choices"     json:"choices"`
	Flag        string   `yaml:"flag"        json:"flag"`
}

type ArgType string

const (
	ArgTypeString ArgType = "string"
	ArgTypeBool   ArgType = "bool"
	ArgTypeEnum   ArgType = "enum"
	ArgTypeFile   ArgType = "file"
	ArgTypeDir    ArgType = "dir"
)

var validArgTypes = map[string]bool{
	string(ArgTypeString): true,
	string(ArgTypeBool):   true,
	string(ArgTypeEnum):   true,
	string(ArgTypeFile):   true,
	string(ArgTypeDir):    true,
}

// Validate checks that a Plugin is structurally sound. It returns the first
// error found so the caller can report it alongside the file name.
func (p Plugin) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("plugin must have a non-empty name")
	}
	if len(p.Commands) == 0 {
		return fmt.Errorf("plugin %q has no commands", p.Name)
	}
	for ci, cmd := range p.Commands {
		if err := validateCommand(p.Name, ci, cmd); err != nil {
			return err
		}
	}
	return nil
}

func validateCommand(pluginName string, idx int, cmd Command) error {
	ref := fmt.Sprintf("plugin %q command[%d]", pluginName, idx)

	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%s: name must not be empty", ref)
	}
	if strings.TrimSpace(cmd.Template) == "" {
		return fmt.Errorf("%s %q: template must not be empty", ref, cmd.Name)
	}

	seen := make(map[string]bool, len(cmd.Args))
	for ai, arg := range cmd.Args {
		if err := validateArg(pluginName, cmd.Name, ai, arg); err != nil {
			return err
		}
		if seen[arg.Name] {
			return fmt.Errorf("plugin %q command %q: duplicate arg name %q", pluginName, cmd.Name, arg.Name)
		}
		seen[arg.Name] = true
	}
	return nil
}

func validateArg(pluginName, cmdName string, idx int, arg Arg) error {
	ref := fmt.Sprintf("plugin %q command %q arg[%d]", pluginName, cmdName, idx)

	if strings.TrimSpace(arg.Name) == "" {
		return fmt.Errorf("%s: name must not be empty", ref)
	}
	if arg.Type != "" && !validArgTypes[arg.Type] {
		return fmt.Errorf("%s %q: unknown type %q (valid: string, bool, enum, file, dir)", ref, arg.Name, arg.Type)
	}
	if arg.Type == string(ArgTypeEnum) && len(arg.Choices) == 0 {
		return fmt.Errorf("%s %q: enum arg must have at least one choice", ref, arg.Name)
	}
	if arg.Type == string(ArgTypeBool) {
		switch strings.TrimSpace(arg.Default) {
		case "true", "false", "":
		default:
			return fmt.Errorf("%s %q: bool arg default must be \"true\" or \"false\", got %q", ref, arg.Name, arg.Default)
		}
	}
	return nil
}

// MatchesOS returns true if this plugin should be loaded on the current OS.
// An empty OS list means the plugin is cross-platform and always loaded.
func (p Plugin) MatchesOS() bool {
	if len(p.OS) == 0 {
		return true
	}
	for _, os := range p.OS {
		if os == runtime.GOOS {
			return true
		}
	}
	return false
}