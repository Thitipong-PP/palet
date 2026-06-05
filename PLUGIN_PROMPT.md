# Task
Write a complete palet plugin YAML file for **[TOOL NAME]**.

---

# Plugin YAML schema

```yaml
name: string                  # plugin name (lowercase)
description: string           # one-line description
hidden: bool                  # true if it a special command

commands:
  - name: string              # e.g. "go run"
    description: string       # shown in command list
    template: string          # Go text/template, args as {{.argname}}
    args:
      - name: string
        description: string
        type: string          # string | bool | enum | file | dir
        required: bool
        default: string       # optional
        choices: []           # only for type: enum
        flag: string          # optional, e.g. "-o"
```

---

# Template rules

- Each arg is referenced as `{{.argname}}` in the template
- Optional args use `{{if .argname}} --flag {{.argname}}{{end}}`
- Bool args use `{{if .argname}} --flag{{end}}`
- String args that need quoting: `"{{.argname}}"`

Examples:
```yaml
# simple
template: "go run {{.file}}"

# optional flag
template: "go build{{if .output}} -o {{.output}}{{end}} {{.package}}"

# bool flag
template: "docker run{{if .detach}} -d{{end}} {{.image}}"

# quoted string
template: "git commit -m \"{{.message}}\""
```

---

# Requirements

- Cover **all commonly used subcommands** of [TOOL NAME]
- Group commands with comments (e.g. `# ── build ───`)
- Every arg must have a clear `description`
- Use `default` for args that have sensible defaults
- Use `type: enum` with `choices` when the arg has a fixed set of values
- Required args that have no default should have `required: true`
- Optional args should have `required: false`
- `args: []` for commands that take no arguments

---

# Output

Return only the YAML file content, no explanation.

---

# Warning

If some field you using coron(:) you need to using quotes (yaml structure will not running if you forgot)