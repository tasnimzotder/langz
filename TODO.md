# Langz TODO

## Language Features

- [x] **else if / elif** — chained conditionals, codegen to Bash `elif`
- [x] **String methods** — `.replace()`, `.contains()`, `.starts_with()`, `.ends_with()`, `.split()`, `.join()`, `.length()`
- [x] **Array indexing** — `items[0]`, bracket access on lists/maps
- [x] **Compound assignment** — `+=`, `-=`, `*=`, `/=`
- [x] **Default parameters** — `fn greet(name: str = "world")`
- [x] **Pipe operator** — `x |> upper()` for chaining builtins
- [ ] **Imports/modules** — split code across files
- [ ] **Floating point** — decimal number support
- [ ] **Multi-line strings** — heredoc or triple-quote syntax

## Tooling

- [x] **LSP server** — diagnostics, hover, completion, signature help, go-to-definition, formatting
- [x] **VS Code extension** — syntax highlighting + LSP integration
- [x] **Vim/Neovim plugin** — syntax highlighting + LSP setup
- [x] **`langz fmt` subcommand** — auto-format `.lz` files
- [x] **Error detection** — ILLEGAL tokens for unterminated strings/unknown chars, parser error reporting, codegen error markers
- [ ] **Better error messages** — include source line context with `^` pointer in more places
- [ ] **Treesitter grammar** — for better syntax highlighting in editors

## More Builtins

- [x] **JSON parsing** — `json_get(data, "key")` via `jq`
- [ ] **Process management** — `pid()`, `kill()`, `ps()`, `bg()` for background processes
- [ ] **Regex** — `match()`, `replace_regex()` for pattern matching
