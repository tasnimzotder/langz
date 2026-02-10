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

This compiles and executes in one step.

## File Extension

LangZ files use the `.lz` extension.

## What Gets Generated

Every LangZ script compiles to a Bash script with:

- `#!/bin/bash` shebang
- `set -euo pipefail` for safe defaults
- Proper variable quoting
- Standard shell idioms

The generated Bash is readable and portable -- no runtime dependencies.
