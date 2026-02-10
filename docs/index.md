# LangZ

A scripting language that transpiles to Bash, designed for DevOps, SRE, and system administration tasks.

Write clean, readable scripts and compile them to portable Bash that runs anywhere.

!!! note
    This project was entirely created by an LLM ([Claude](https://claude.ai) by Anthropic), including the language design, compiler, LSP server, VS Code extension, tests, and documentation.

## Why LangZ?

Bash is powerful but hard to read and write correctly. LangZ gives you:

- **Clean syntax** -- no `$`, `fi`, `esac`, or quoting headaches
- **String interpolation** -- `"Hello {name}"` instead of `"Hello $name"`
- **Built-in DevOps functions** -- file ops, system info, HTTP requests, JSON parsing
- **Safe defaults** -- generates `set -euo pipefail` automatically
- **Zero runtime** -- compiles to plain Bash, nothing to install on target
- **Imports** -- split code across files with `import "lib.lz"`
- **Bash escape hatch** -- embed raw shell with `bash { ... }` when needed
- **Shebang support** -- `#!/usr/bin/env langz` for directly executable scripts

## Quick Example

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

This compiles to clean, portable Bash with `set -euo pipefail`, proper quoting, and standard shell idioms.

## Next Steps

- [Getting Started](getting-started.md) -- install and write your first script
- [Language Guide](language/variables.md) -- learn the syntax
- [Builtins Reference](builtins.md) -- all built-in functions
- [Examples](examples.md) -- real-world scripts
