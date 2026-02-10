# Editor Support

LangZ includes a built-in Language Server Protocol (LSP) server for rich editor integration.

## Features

| Feature | Description |
|---------|-------------|
| **Diagnostics** | Real-time parse error highlighting |
| **Hover** | Documentation for builtins and keyword arguments |
| **Completion** | Builtins, keywords, user symbols, and context-aware kwargs |
| **Signature Help** | Parameter hints when typing function calls |
| **Go-to-Definition** | Jump to variable and function definitions |
| **Document Symbols** | Outline view of variables and functions |
| **Formatting** | Auto-format `.lz` files |

## VS Code Extension

A VS Code extension is included in the `editors/vscode/` directory.

### Install from Source

```bash
cd editors/vscode
bun install
bun run build
```

Then install the `.vsix` file via VS Code: **Extensions > ... > Install from VSIX**.

### Features in VS Code

- Syntax highlighting for `.lz` files
- All LSP features listed above
- Context-aware keyword argument completion inside `fetch()` and other builtins
- Hover tooltips on kwargs like `timeout:`, `method:`, etc.

## Other Editors

The LSP server can be used with any editor that supports the Language Server Protocol. Start the server with:

```bash
langz lsp
```

Configure your editor to use `langz lsp` as the language server for `.lz` files.
