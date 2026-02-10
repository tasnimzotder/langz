# Internals

How LangZ works under the hood.

## Compilation Pipeline

```
Source (.lz) → Lexer → Parser → AST → Codegen → Bash (.sh)
```

| Stage | Description |
|-------|-------------|
| **Lexer** | Tokenizes source into tokens (identifiers, strings, operators, keywords). Emits `ILLEGAL` tokens for malformed input. Supports unicode identifiers |
| **Parser** | Recursive descent parser builds an Abstract Syntax Tree. Reports structured errors for invalid tokens |
| **Codegen** | Walks the AST and emits Bash code. Marks unhandled nodes with `# error:` comments |

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

The server caches tokenized output per document (`Server.tokens` map), populated on `didOpen`/`didChange` and evicted on `didClose`. This avoids re-tokenizing on every hover, completion, or diagnostic request.

## Error Handling Architecture

Errors are handled at each pipeline stage, forming layered protection:

### Lexer: ILLEGAL Tokens

The lexer never silently skips input. Unknown characters and unterminated strings produce `ILLEGAL` tokens with descriptive values:

- `"unterminated string"` -- when a string literal reaches EOF without a closing `"`
- The character itself (e.g. `"@"`) -- when the lexer encounters an unrecognized character

### Parser: Error Reporting

The parser converts `ILLEGAL` tokens and unexpected token types into structured `ParseError` values:

- `ILLEGAL` tokens: the token's value becomes the error message (e.g. "unterminated string")
- Unknown tokens: `"unexpected token: <type>"` is reported

Errors accumulate in `Parser.errors` and are returned by `ParseAllErrors()`. The parser continues past errors to report as many as possible.

### Codegen: Error Markers

Unhandled AST node types in `genStatement()` and `genExpr()` emit Bash comments like `# error: unhandled statement type *ast.ExitCall`. The existing `findCodegenErrors()` function scans the generated output for these markers and returns them as structured errors.

### Panic Recovery

All public API entry points use `defer/recover`:

- **`Generate()`** -- returns empty output with `"internal error: ..."` on panic
- **`ParseWithErrors()`** / **`ParseAllErrors()`** -- returns partial AST with error on panic
- **Every LSP handler** -- uses a shared `recoverErr()` helper to convert panics to LSP errors, preventing server crashes

### Unicode Support

The lexer operates on runes (not bytes) using `utf8.DecodeRuneInString()`. Identifiers can contain any Unicode letter (`unicode.IsLetter()`) or digit (`unicode.IsDigit()`), supporting non-ASCII variable names like `名前` or `café`.

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
