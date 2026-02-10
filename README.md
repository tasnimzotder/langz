# LangZ

A scripting language that transpiles to Bash, designed for DevOps, SRE, and system administration tasks.

Write clean, readable scripts in LangZ and compile them to portable Bash that runs anywhere.

> **Note:** This project was entirely created by an LLM ([Claude](https://claude.ai) by Anthropic), including the language design, compiler, LSP server, VS Code extension, tests, and documentation.

## Why LangZ?

Bash is powerful but hard to read and write correctly. LangZ gives you:

- **Clean syntax** -- no `$`, `fi`, `esac`, or quoting headaches
- **String interpolation** -- `"Hello {name}"` instead of `"Hello $name"`
- **Built-in DevOps functions** -- file ops, system info, HTTP requests, JSON parsing
- **Safe defaults** -- generates `set -euo pipefail` automatically
- **Zero runtime** -- compiles to plain Bash, nothing to install on target

## Quick Start

```bash
go install github.com/tasnimzotder/langz/cmd/langz@latest
```

Create `deploy.lz`:

```
app = "webapp"
env_name = env("DEPLOY_ENV") or "staging"

fn log(msg: str) {
    print("[deploy] {msg}")
}

log("Deploying {app} to {env_name}")

mkdir("dist")
write("dist/version.txt", "1.0.0")

if exists("config.json") {
    copy("config.json", "dist/config.json")
} else {
    log("No config, using defaults")
}
```

Compile or run directly:

```bash
langz build deploy.lz   # generates deploy.sh
langz run deploy.lz     # compile and execute
```

## Language Features

### Variables

```
name = "langz"
port = 8080
verbose = true
items = ["a", "b", "c"]
config = {host: "localhost", port: 8080}
```

### String Interpolation

```
host = "localhost"
port = 3000
print("Server at {host}:{port}")
```

### Functions

```
fn greet(name: str) {
    print("Hello {name}")
}

fn add(a: int, b: int) -> int {
    return a + b
}

greet("world")
```

### Control Flow

```
// if/else
if status == "ready" {
    print("go")
} else {
    print("wait")
}

// for loops
for item in items {
    print(item)
}

for i in range(1, 10) {
    print(i)
}

// while loops
while retries > 0 {
    print("trying...")
    break
}

// match (compiles to case/esac)
match platform {
    "darwin" => print("macOS")
    "linux"  => print("Linux")
    _        => print("unknown")
}
```

### Error Handling with `or`

```
// Default values
name = env("APP_NAME") or "myapp"

// Exit on failure
data = read("config.json") or exit(1)

// Skip in loops
content = read(f) or continue

// Block fallback
result = exec("deploy.sh") or {
    print("deploy failed")
    "fallback"
}
```

### Arithmetic

```
a = 10
b = 3
sum = a + b
product = a * b
remainder = a % b
complex = (a + b) * 2

if a + b > 10 {
    print("big")
}
```

## Built-in Functions

### I/O

| Function | Description | Bash output |
|----------|-------------|-------------|
| `print(msg)` | Print to stdout | `echo "msg"` |
| `read(path)` | Read file contents | `$(cat "path")` |
| `write(path, data)` | Write to file | `echo "data" > "path"` |
| `append(path, data)` | Append to file | `echo "data" >> "path"` |

### File Operations

| Function | Description | Bash output |
|----------|-------------|-------------|
| `exists(path)` | Check if path exists | `[ -e "path" ]` |
| `is_file(path)` | Check if file | `[ -f "path" ]` |
| `is_dir(path)` | Check if directory | `[ -d "path" ]` |
| `mkdir(path)` | Create directory | `mkdir -p "path"` |
| `rm(path)` | Remove file | `rm -f "path"` |
| `rmdir(path)` | Remove directory | `rm -rf "path"` |
| `copy(src, dst)` | Copy file | `cp "src" "dst"` |
| `move(src, dst)` | Move file | `mv "src" "dst"` |
| `chmod(path, mode)` | Change permissions | `chmod 755 "path"` |
| `glob(pattern)` | Glob files | `(*.log)` |

### System

| Function | Description | Bash output |
|----------|-------------|-------------|
| `exec(cmd)` | Run shell command | `$(cmd)` |
| `env(name)` | Get env variable | `"${NAME}"` |
| `os()` | Get OS name | `$(uname -s \| tr ...)` |
| `arch()` | Get architecture | `$(uname -m)` |
| `hostname()` | Get hostname | `$(hostname)` |
| `whoami()` | Get current user | `$(whoami)` |
| `sleep(n)` | Sleep n seconds | `sleep n` |
| `exit(code)` | Exit with code | `exit code` |

### String & Path

| Function | Description | Bash output |
|----------|-------------|-------------|
| `upper(s)` | Uppercase | `$(echo "s" \| tr ...)` |
| `lower(s)` | Lowercase | `$(echo "s" \| tr ...)` |
| `trim(s)` | Trim whitespace | `$(echo "s" \| xargs)` |
| `len(list)` | List length | `${#list[@]}` |
| `dirname(path)` | Directory name | `$(dirname "path")` |
| `basename(path)` | Base name | `$(basename "path")` |

### Networking

| Function | Description |
|----------|-------------|
| `fetch(url, ...)` | HTTP request via curl (see below) |
| `json_get(data, path)` | Extract JSON value via jq |

#### `fetch()` -- Full HTTP Support

```
// Simple GET
data = fetch("https://api.example.com/health")

// POST with headers, body, timeout, and retries
resp = fetch("https://api.example.com/users",
  method: "POST",
  body: payload,
  headers: {"Content-Type": "application/json"},
  timeout: 10,
  retries: 3
)

// Convention variables set after every fetch:
//   _status  -- HTTP status code (e.g. 200)
//   _body    -- response body
//   _headers -- response headers

if _status == 200 {
  name = json_get(_body, ".name")
  print("Created: {name}")
}

// Error fallback with `or`
data = fetch("https://api.example.com/data") or "unavailable"
```

Supported keyword arguments:

| Kwarg | Description | Default |
|-------|-------------|---------|
| `method:` | HTTP method (GET, POST, PUT, PATCH, DELETE) | GET |
| `body:` | Request body data | none |
| `headers:` | Request headers as map | none |
| `timeout:` | Max seconds to wait | none |
| `retries:` | Number of retry attempts | none |

### Date/Time

| Function | Description | Bash output |
|----------|-------------|-------------|
| `timestamp()` | Unix timestamp | `$(date +%s)` |
| `date()` | Current date (YYYY-MM-DD) | `$(date +"%Y-%m-%d")` |

## Example: Deployment Script

LangZ source (`deploy.lz`):

```
app = "webapp"
version = "2.1.0"
platform = os()
host = hostname()

fn log(msg: str) {
    print("[deploy] {msg}")
}

log("Deploying {app} v{version} on {host}")

mkdir("dist")
write("dist/manifest.txt", "app={app}")
append("dist/manifest.txt", "version={version}")

env_name = env("DEPLOY_ENV") or "staging"

match env_name {
    "production" => log("PRODUCTION deploy")
    "staging"    => log("Staging deploy")
    _            => log("Unknown: {env_name}")
}

log("Done!")
```

Generated Bash:

```bash
#!/bin/bash
set -euo pipefail

app="webapp"
version="2.1.0"
platform=$(uname -s | tr '[:upper:]' '[:lower:]')
host=$(hostname)
log() {
  local msg="$1"
  echo "[deploy] ${msg}"
}
log "Deploying ${app} v${version} on ${host}"
mkdir -p "dist"
echo "app=${app}" > "dist/manifest.txt"
echo "version=${version}" >> "dist/manifest.txt"
env_name="${DEPLOY_ENV:-staging}"
case "$env_name" in
  production)
    log "PRODUCTION deploy"
    ;;
  staging)
    log "Staging deploy"
    ;;
  *)
    log "Unknown: ${env_name}"
    ;;
esac
log "Done!"
```

## Editor Support

LangZ includes a built-in LSP server with:

- **Diagnostics** -- real-time parse error highlighting
- **Hover** -- documentation for all builtins and keyword arguments
- **Completion** -- builtins, keywords, user-defined symbols, and context-aware kwarg suggestions inside function calls
- **Signature help** -- parameter hints when typing function calls
- **Go-to-definition** -- jump to variable and function definitions
- **Document symbols** -- outline view of variables and functions
- **Formatting** -- auto-format `.lz` files

A VS Code extension is included in the `editors/vscode/` directory.

## Project Structure

```
langz/
├── cmd/langz/          CLI entry point
├── internal/
│   ├── ast/            Abstract syntax tree nodes
│   ├── lexer/          Tokenizer
│   ├── parser/         Recursive descent parser
│   ├── codegen/        Bash code generator
│   │   └── builtins/   Built-in function registry
│   └── lsp/            Language Server Protocol
├── editors/vscode/     VS Code extension
├── test/integration/   End-to-end tests
├── examples/           Example .lz scripts
└── go.mod
```

## Development

```bash
# Run all tests
gotestsum -- ./...

# Run specific test
gotestsum -- -run TestE2E_HelloWorld ./test/integration/

# Build and install
go install ./cmd/langz
```

## License

MIT
