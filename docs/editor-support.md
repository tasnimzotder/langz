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

## Neovim

A Vim/Neovim plugin is included in the `editors/nvim/` directory with syntax highlighting, filetype detection, and comment settings.

### Install the Plugin

**With [lazy.nvim](https://github.com/folke/lazy.nvim):**

```lua
{
  dir = "~/path/to/langz/editors/nvim",
  ft = "langz",
}
```

**With manual symlink:**

```bash
# Neovim
ln -s /path/to/langz/editors/nvim ~/.local/share/nvim/site/pack/langz/start/langz

# Vim
ln -s /path/to/langz/editors/nvim ~/.vim/pack/langz/start/langz
```

### LSP Setup (Neovim)

Add this to your Neovim config (e.g. `~/.config/nvim/init.lua` or in `after/ftplugin/langz.lua`):

**With `nvim-lspconfig` (custom server):**

```lua
local lspconfig = require("lspconfig")
local configs = require("lspconfig.configs")

if not configs.langz then
  configs.langz = {
    default_config = {
      cmd = { "langz", "lsp" },
      filetypes = { "langz" },
      root_dir = lspconfig.util.find_git_ancestor,
      single_file_support = true,
    },
  }
end

lspconfig.langz.setup({})
```

**Without `nvim-lspconfig`:**

```lua
vim.api.nvim_create_autocmd("FileType", {
  pattern = "langz",
  callback = function()
    vim.lsp.start({
      name = "langz",
      cmd = { "langz", "lsp" },
      root_dir = vim.fs.root(0, { ".git" }),
    })
  end,
})
```

### LSP Setup (Vim)

With [vim-lsp](https://github.com/prabirshrestha/vim-lsp):

```vim
if executable('langz')
  au User lsp_setup call lsp#register_server(#{
    \ name: 'langz',
    \ cmd: ['langz', 'lsp'],
    \ allowlist: ['langz'],
    \ })
endif
```

## Other Editors

The LSP server can be used with any editor that supports the Language Server Protocol. Start the server with:

```bash
langz lsp
```

Configure your editor to use `langz lsp` as the language server for `.lz` files.
