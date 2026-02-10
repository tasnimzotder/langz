# Variables & Types

## Assignment

Variables are assigned with `=`. No `let`, `var`, or `$` needed.

```
name = "langz"
port = 8080
verbose = true
```

**Generated Bash:**
```bash
name="langz"
port=8080
verbose=true
```

## Types

LangZ supports these value types:

| Type | Example | Bash |
|------|---------|------|
| String | `"hello"` | `"hello"` |
| Number | `42` | `42` |
| Boolean | `true` / `false` | `true` / `false` |
| List | `["a", "b", "c"]` | `("a" "b" "c")` |
| Map | `{host: "localhost"}` | `varname_key=val` |

## String Interpolation

Use `{variable}` inside strings to interpolate values:

```
host = "localhost"
port = 3000
print("Server at {host}:{port}")
```

**Generated Bash:**
```bash
echo "Server at ${host}:${port}"
```

## Lists

```
items = ["alpha", "beta", "gamma"]

for item in items {
    print(item)
}
```

## Maps

Maps use `key: value` syntax. Keys can be identifiers or strings:

```
config = {host: "localhost", port: 8080}
headers = {"Content-Type": "application/json", "Accept": "text/html"}
```

String keys are useful when keys contain special characters like hyphens.

## Indexing

Access array elements by index and map values by string key:

```
items = ["alpha", "beta", "gamma"]
first = items[0]
print(first)

config = {host: "localhost", port: "8080"}
val = config["host"]
print(val)
```

Index assignment:

```
items[1] = "BETA"
```

## Compound Assignment

Shorthand operators `+=`, `-=`, `*=`, `/=` for updating variables:

```
count = 0
count += 1
total = 100
total -= 10
```

**Generated Bash:**
```bash
count=$((count + 1))
total=$((total - 10))
```

## Arithmetic

Standard arithmetic operators work on numeric values:

```
a = 10
b = 3
sum = a + b
product = a * b
remainder = a % b
complex = (a + b) * 2
```

**Generated Bash:**
```bash
sum=$((a + b))
product=$((a * b))
```
