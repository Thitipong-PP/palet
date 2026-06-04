---
name: Plugin request
about: Suggest an plugin for this project
title: "[Plugins]"
labels: enhancement, good first issue
assignees: Thitipong-PP

---

---

**name:** New Plugin Request

**about:** Request or contribute a new plugin for palet

**title:** "feat(plugin): add <tool-name> plugin"

**labels:** plugin, enhancement

**assignees:** ""

---

## Plugin info

| Field | Value |
|-------|-------|
| Tool name | <!-- e.g. kubectl --> |
| Tool description | <!-- e.g. Kubernetes command-line tool --> |
| Official docs | <!-- e.g. https://kubernetes.io/docs/reference/kubectl/ --> |
| hidden | <!-- true (specialized) or false (everyone uses it) --> |
| os | <!-- leave empty for all platforms, or e.g. [linux, darwin] --> |

## Commands to cover

List the commands you want to include:

- [ ] `tool subcommand` — description
- [ ] `tool subcommand` — description
- [ ] `tool subcommand` — description

## Contributing

If you'd like to implement this yourself:

1. Use `PLUGIN_PROMPT.md` as the prompt template
2. Create `plugins/<tool-name>.yaml` following the existing schema
3. Test locally by placing the file in `~/.config/palet/plugins/`
4. Open a PR referencing this issue
