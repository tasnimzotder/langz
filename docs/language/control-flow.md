# Control Flow

## If/Else

```
if status == "ready" {
    print("go")
} else {
    print("wait")
}
```

Conditions support `==`, `!=`, `>`, `<`, `>=`, `<=` and logical operators `and`, `or`:

```
if port > 1024 and is_file("config.json") {
    print("ready")
}
```

## Elif (Else-If Chaining)

Chain multiple conditions without nesting:

```
env_name = env("ENV") or "dev"
if env_name == "production" {
    print("PROD")
} else if env_name == "staging" {
    print("STAGING")
} else {
    print("DEV")
}
```

**Generated Bash:**
```bash
if [ "$env_name" = "production" ]; then
  echo "PROD"
elif [ "$env_name" = "staging" ]; then
  echo "STAGING"
else
  echo "DEV"
fi
```

## For Loops

### Iterate over a list

```
items = ["alpha", "beta", "gamma"]
for item in items {
    print(item)
}
```

### Iterate over a range

```
for i in range(1, 10) {
    print(i)
}
```

### Iterate over files

```
for f in glob("*.log") {
    print("Found: {f}")
}
```

## While Loops

```
retries = 3
while retries > 0 {
    result = exec("./deploy.sh") or ""
    if result != "" {
        break
    }
    retries -= 1
    sleep(1)
}
```

## Match (Pattern Matching)

`match` compiles to Bash `case`/`esac`:

```
platform = os()
match platform {
    "darwin" => print("macOS")
    "linux"  => print("Linux")
    _        => print("unknown: {platform}")
}
```

Match arms can have block bodies:

```
match env("DEPLOY_ENV") or "dev" {
    "production" => {
        print("PRODUCTION")
        exec("notify-slack.sh")
    }
    "staging" => print("staging")
    _ => print("dev")
}
```

## Break and Continue

`break` exits a loop, `continue` skips to the next iteration:

```
for f in glob("*.txt") {
    content = read(f) or continue
    print(content)
}
```

## Raw Bash Blocks

For shell-specific logic without a LangZ equivalent, use `bash { }` to embed raw Bash:

```
bash {
    set -euo pipefail
    trap 'cleanup' EXIT
}
```

Content is emitted verbatim into the generated script. Nested braces are handled correctly:

```
bash {
    if command -v docker &>/dev/null; then
        echo "Docker is installed"
    else
        echo "Docker not found"
    fi
}
```

You can mix bash blocks with regular LangZ code:

```
name = "deploy"
bash { echo "Running: $name" }
print("Done with {name}")
```
