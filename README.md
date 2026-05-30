# PALET
A terminal command palette. Search commands, fill in args, and run — no more memorising flags.

<img width="1470" height="956" alt="Screenshot 2569-05-30 at 23 03 31" src="https://github.com/user-attachments/assets/1dd9fb5b-75f8-46d6-8de6-f578e6bf430f" />

## Install

```bash
go install github.com/Thitipong-PP/palet@latest
```

Or build from source:

```bash
git clone https://github.com/Thitipong-PP/palet
cd palet
go build -o palet .
```

Then make sure `$GOPATH/bin` is in your PATH:

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH="$PATH:$(go env GOPATH)/bin"

# Then reload
source ~/.zshrc   # or source ~/.bashrc
```

Already done? just run `palet` directly.

## Usage
Run palet to open the command palette.
```bash
palet
```

| Key | Action |
|-----|--------|
| type | search commands |
| `↑` `↓` or `j` `k` | navigate |
| `enter` | select command / next arg field |
| `tab` | next arg field |
| `shift+tab` | previous arg field |
| `e` | execute command |
| `c` | copy to clipboard |
| `esc` | go back |
| `q` | quit |

## Plugins

Commands live in YAML files. palet loads plugins from two locations:
 
| Path | Purpose |
|------|---------|
| `./plugins/` | project-local plugins |
| `~/.config/palet/plugins/` | user-global plugins |
 
Both directories are scanned on startup. Only files that have changed since the last run are re-parsed — unchanged files are loaded from cache at `~/.cache/palet/index.json`.

### Writing a plugin

```yaml
name: git
description: Git version control commands

commands:
  - name: git add
    description: Add stage to commit
    template: 'git add "{{.directory}}"'
    args:
      - name: directory
        description: Directory name
        type: string
        required: true
  - name: git commit
    description: Commit staged changes with a message
    template: 'git commit -m "{{.message}}"'
    args:
      - name: message
        description: Commit message
        type: string
        required: true

  - name: git push
    description: Push commits to a remote branch
    template: "git push {{.remote}} {{.branch}}"
    args:
      - name: remote
        description: Remote name
        type: string
        required: true
        default: "origin"
      - name: branch
        description: Branch name
        type: string
        required: true
        default: "main"
```

### Arg types
 
| type | description |
|------|-------------|
| `string` | plain text |
| `bool` | true / false — renders as `{{if .name}}` in template |
| `enum` | pick from `choices` list |
| `file` | file path |
| `dir` | directory path |

### Template syntax
 
Templates use Go's `text/template`. Each arg is available as `{{.argname}}`.
 
```yaml
# simple
template: "go run {{.file}}"
 
# optional flag
template: "go test {{.flags}} {{.package}}"
 
# conditional
template: "git commit{{if .amend}} --amend{{end}} -m \"{{.message}}\""
```

## License
[Apache-2.0 license](https://github.com/Thitipong-PP/palet/blob/main/LICENSE)
