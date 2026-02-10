# Functions

## Defining Functions

Use `fn` to define functions with typed parameters:

```
fn greet(name: str) {
    print("Hello {name}")
}

greet("world")
```

**Generated Bash:**
```bash
greet() {
  local name="$1"
  echo "Hello ${name}"
}
greet "world"
```

## Return Values

Functions can return values with `-> type` and `return`:

```
fn add(a: int, b: int) -> int {
    return a + b
}

result = add(3, 5)
print(result)
```

## Multiple Parameters

```
fn deploy(app: str, env: str, version: str) {
    print("Deploying {app} v{version} to {env}")
    mkdir("dist/{app}")
    write("dist/{app}/version.txt", version)
}
```

## Script Arguments

Access command-line arguments with `args()`:

```
for arg in args() {
    print("arg: {arg}")
}
```

**Generated Bash:**
```bash
for arg in "$@"; do
  echo "arg: ${arg}"
done
```
