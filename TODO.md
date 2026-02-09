# Langz TODO

## Language Features

- [ ] **else if / elif** — chained conditionals, codegen to Bash `elif`
- [ ] **String methods** — `.len()`, `.contains()`, `.split()`, `.trim()`, `.replace()` via dot syntax
- [ ] **Array indexing** — `items[0]`, `items[-1]`, bracket access on lists/maps

## Tooling

- [ ] **Better error messages** — include source line context, column numbers, hints for common mistakes
- [ ] **`langz fmt` subcommand** — auto-format `.lz` files (indentation, spacing)
- [ ] **LSP basics** — language server for editor support (diagnostics, go-to-definition)

## More Builtins

- [ ] **Networking** — `http_get()`, `http_post()`, `ping()`, `curl` wrappers
- [ ] **JSON parsing** — `json_get(data, "key")` via `jq`, parse/extract fields
- [ ] **Process management** — `pid()`, `kill()`, `ps()`, `bg()` for background processes
