# Internals

How LangZ works under the hood.

## Compilation Pipeline

```
Source (.lz) → Lexer → Parser → AST → Codegen → Bash (.sh)
```

| Stage | Description |
|-------|-------------|
| **Lexer** | Tokenizes source into tokens (identifiers, strings, operators, keywords) |
| **Parser** | Recursive descent parser builds an Abstract Syntax Tree |
| **Codegen** | Walks the AST and emits Bash code |

## Project Structure

```
langz/
├── cmd/langz/              CLI entry point
├── internal/
│   ├── ast/                AST node types
│   ├── lexer/              Tokenizer (tokens.go, lexer.go)
│   ├── parser/             Recursive descent parser
│   │   ├── parser.go       Core parser, entry points
│   │   ├── expressions.go  Expression parsing
│   │   └── statements.go   Statement parsing
│   ├── codegen/            Bash code generator
│   │   ├── codegen.go      Core generator
│   │   ├── expressions.go  Expression codegen
│   │   ├── statements.go   Statement codegen
│   │   ├── fetch.go        fetch() codegen (multi-line curl)
│   │   └── builtins/       Built-in function registry
│   └── lsp/                Language Server Protocol
├── editors/vscode/         VS Code extension
├── test/integration/       End-to-end tests
├── examples/               Example .lz scripts
└── docs/                   Documentation (MkDocs)
```

## Key Design Decisions

### No Parser Generator

The parser is hand-written recursive descent. This keeps dependencies minimal and makes the parser easy to extend with new syntax.

### Two Codegen Function Types

- **`ExprGen`** -- generates quoted expressions (e.g. `"$var"`)
- **`RawValueGen`** -- generates unquoted values (e.g. `var` for arithmetic)

This distinction is critical for Bash correctness. Arithmetic contexts like `$((...))` need unquoted variable names, while most other contexts need quoted values.

### Builtin Registry Pattern

Built-in functions are registered in maps (`stmtBuiltins`, `exprBuiltins`) with a common handler signature. This makes adding new builtins a one-line change.

### Statement-Level Fetch

`fetch()` can't be a simple expression builtin because it needs multi-line output (mktemp, curl, parse status, cleanup). It's intercepted at the codegen level in `genAssignment()` and `genFuncCallStmt()` before the normal builtin dispatch.

### Convention Variables

`fetch()` sets `_status`, `_body`, `_headers` as global variables. This avoids the need for structured return types while keeping the API simple.

### Token-Based LSP

The LSP server uses token streams (not AST) for hover, completion, and definition. This avoids needing position tracking in the AST and keeps the LSP implementation simple and robust.

## Testing

The project uses [testify](https://github.com/stretchr/testify) for assertions and [gotestsum](https://github.com/gotestyourself/gotestsum) as test runner.

```bash
# Run all tests
gotestsum -- ./...

# Run specific test
gotestsum -- -run TestName ./path/

# Run integration tests only
gotestsum -- ./test/integration/
```

Test categories:

- **Parser tests** -- verify AST construction from source
- **Codegen tests** -- verify generated Bash from AST
- **Integration tests** -- compile LangZ to Bash, execute it, check output
- **LSP tests** -- verify diagnostics, hover, completion, signature help
