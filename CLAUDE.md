# LangZ

Scripting language that transpiles to Bash. Go module: `github.com/tasnimzotder/langz`

## Commands

- `gotestsum -- ./...` — run all tests (170 total)
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
- `test/integration/` — E2E tests that compile and run Bash
- `examples/` — Sample .lz scripts

## Code Patterns

- Builtins use registry maps (`stmtBuiltins`, `exprBuiltins`) with `builtinHandler` function type
- Two codegen function types: `ExprGen` (quoted) and `RawValueGen` (unquoted)
- Parser methods split across files but share `*Parser` receiver
- Integration tests use `projectRoot()` helper to locate `go.mod` for CLI subprocess tests

## Known Issues

- `or` works in assignments but not yet in conditions (only `and` works in conditions)

## Style

- Testify (`assert`/`require`) for all assertions
- Table-driven tests where applicable
- No parser generator — hand-written recursive descent
