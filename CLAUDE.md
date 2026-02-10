# LangZ

Scripting language that transpiles to Bash. Go module: `github.com/tasnimzotder/langz`

## Commands

- `gotestsum -- ./...` — run all tests (367 total)
- `gotestsum -- -run TestName ./path/` — run specific test
- `go build ./...` — verify compilation
- `go install ./cmd/langz` — install CLI binary

## Architecture

Pipeline: Source (.lz) → Lexer → Parser (recursive descent) → Codegen → Bash

## Project Layout

- `cmd/langz/` — CLI entry point
- `internal/ast/` — AST node types
- `internal/lexer/` — Tokenizer (tokens.go for types, lexer.go for scanner)
- `internal/parser/` — Parser (parser.go core, expressions.go, statements.go)
- `internal/codegen/` — Bash generator (codegen.go core, expressions.go, statements.go)
- `internal/codegen/builtins/` — Built-in function registry (builtins.go types, exprs.go, stmts.go)
- `internal/lsp/` — Language Server Protocol (server.go core, per-feature files)
- `test/integration/` — E2E tests that compile and run Bash
- `examples/` — Sample .lz scripts

## Code Patterns

- Builtins use registry maps (`stmtBuiltins`, `exprBuiltins`) with `builtinHandler` function type
- Two codegen function types: `ExprGen` (quoted) and `RawValueGen` (unquoted)
- Parser methods split across files but share `*Parser` receiver
- Integration tests use `projectRoot()` helper to locate `go.mod` for CLI subprocess tests

## Error Handling Architecture

- **Lexer** emits `ILLEGAL` tokens for unterminated strings and unknown characters (never silently skips)
- **Parser** reports errors for `ILLEGAL` tokens and unexpected token types via `p.addError()`
- **Codegen** emits `# error:` comment markers for unhandled AST nodes; `findCodegenErrors()` scans output for these
- **Panic recovery** via `defer/recover` in all public APIs: `Generate()`, `ParseWithErrors()`, `ParseAllErrors()`, and every LSP handler
- **LSP token cache** (`Server.tokens` map) avoids re-tokenizing on every request; populated on `didOpen`/`didChange`, evicted on `didClose`
- **LSP document guards** check map existence before accessing `s.documents[uri]` in all handlers

## Style

- Testify (`assert`/`require`) for all assertions
- Table-driven tests where applicable
- No parser generator — hand-written recursive descent
- Unicode identifiers supported via `unicode.IsLetter()` / `utf8.DecodeRuneInString()`
