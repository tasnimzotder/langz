# Error Handling

LangZ uses the `or` keyword for error handling. When an expression fails, the fallback value is used instead.

## Default Values

```
name = env("APP_NAME") or "myapp"
config = read("config.json") or "{}"
```

**Generated Bash:**
```bash
name="${APP_NAME:-myapp}"
```

## Exit on Failure

```
data = read("required.json") or exit(1)
```

## Skip in Loops

```
for f in glob("*.json") {
    content = read(f) or continue
    print(content)
}
```

## Block Fallback

Use a block for multi-statement fallbacks:

```
result = exec("deploy.sh") or {
    print("deploy failed, using fallback")
    "fallback_value"
}
```

## Fetch Error Handling

`fetch()` sets the `_status` convention variable, and supports `or` fallback:

```
// Fallback on HTTP error
data = fetch("https://api.example.com/data") or "unavailable"

// Check status explicitly
fetch("https://api.example.com/health")
if _status != 200 {
    print("Health check failed: {_status}")
    exit(1)
}
```

## How It Works

LangZ generates `set -euo pipefail` by default, which means any command failure exits the script. The `or` keyword wraps expressions in error-handling patterns:

- **env()** uses Bash parameter defaults: `${VAR:-default}`
- **General expressions** use `if cmd 2>/dev/null; then ... else ... fi`
- **fetch()** uses `|| true` to prevent `set -e` from killing the script, then checks `_status`
