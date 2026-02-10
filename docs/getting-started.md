# Getting Started

## Installation

### From Source (Go)

```bash
go install github.com/tasnimzotder/langz/cmd/langz@latest
```

### From GitHub Releases

Download the latest binary from [GitHub Releases](https://github.com/tasnimzotder/langz/releases) for your platform (linux/darwin, amd64/arm64).

## Your First Script

Create a file called `hello.lz`:

```
name = "World"
print("Hello {name}!")
```

### Compile to Bash

```bash
langz build hello.lz
```

This generates `hello.sh`:

```bash
#!/bin/bash
set -euo pipefail

name="World"
echo "Hello ${name}!"
```

### Run Directly

```bash
langz run hello.lz
```

This compiles and executes in one step. You can also omit the `run` subcommand:

```bash
langz hello.lz
```

### Shebang Support

Add a shebang line to make `.lz` files directly executable:

```
#!/usr/bin/env langz
name = "World"
print("Hello {name}!")
```

```bash
chmod +x hello.lz
./hello.lz
```

The shebang line is automatically skipped by the compiler.

## Imports

Split code across files using `import`:

```
// lib/helpers.lz
fn greet(name: str) {
    print("Hello {name}")
}
```

```
// main.lz
import "lib/helpers.lz"
greet("world")
```

Import paths are resolved relative to the importing file. Circular imports are detected and reported as errors.

## File Extension

LangZ files use the `.lz` extension.

## What Gets Generated

Every LangZ script compiles to a Bash script with:

- `#!/bin/bash` shebang
- `set -euo pipefail` for safe defaults
- Proper variable quoting
- Standard shell idioms

The generated Bash is readable and portable -- no runtime dependencies.
