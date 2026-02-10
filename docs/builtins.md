# Builtins Reference

All built-in functions available in LangZ.

## I/O

| Function | Description | Bash |
|----------|-------------|------|
| `print(args...)` | Print to stdout | `echo args` |
| `read(path)` | Read file contents | `$(cat path)` |
| `write(path, content)` | Write to file | `echo content > path` |
| `append(path, content)` | Append to file | `echo content >> path` |

## File Operations

| Function | Description | Bash |
|----------|-------------|------|
| `exists(path)` | Check if path exists | `[ -e path ]` |
| `is_file(path)` | Check if file | `[ -f path ]` |
| `is_dir(path)` | Check if directory | `[ -d path ]` |
| `mkdir(path)` | Create directory (with parents) | `mkdir -p path` |
| `rm(path)` | Remove file | `rm -f path` |
| `rmdir(path)` | Remove directory recursively | `rm -rf path` |
| `copy(src, dst)` | Copy file | `cp src dst` |
| `move(src, dst)` | Move/rename file | `mv src dst` |
| `chmod(path, mode)` | Change permissions | `chmod mode path` |
| `chown(path, owner)` | Change owner | `chown owner path` |
| `glob(pattern)` | Expand glob pattern | `(pattern)` |

## System

| Function | Description | Bash |
|----------|-------------|------|
| `exec(cmd)` | Run shell command | `$(cmd)` |
| `env(name)` | Get environment variable | `"${NAME}"` |
| `os()` | Get OS name (lowercase) | `$(uname -s \| tr ...)` |
| `arch()` | Get CPU architecture | `$(uname -m)` |
| `hostname()` | Get hostname | `$(hostname)` |
| `whoami()` | Get current user | `$(whoami)` |
| `sleep(n)` | Sleep n seconds | `sleep n` |
| `exit(code)` | Exit with status code | `exit code` |
| `args()` | Get script arguments | `("$@")` |

## String & Path

| Function | Description | Bash |
|----------|-------------|------|
| `upper(s)` | Convert to uppercase | `$(echo s \| tr ...)` |
| `lower(s)` | Convert to lowercase | `$(echo s \| tr ...)` |
| `trim(s)` | Trim whitespace | `$(echo s \| xargs)` |
| `len(list)` | Get list length | `${#list[@]}` |
| `dirname(path)` | Directory part of path | `$(dirname path)` |
| `basename(path)` | Filename part of path | `$(basename path)` |
| `range(start, end)` | Generate number sequence | `$(seq start end)` |

## String Methods

Methods called on string variables:

| Method | Description | Bash |
|--------|-------------|------|
| `s.replace(old, new)` | Replace all occurrences | `"${s//old/new}"` |
| `s.contains(substr)` | Check if string contains substr | `[[ "$s" == *"substr"* ]]` |
| `s.starts_with(prefix)` | Check if string starts with prefix | `[[ "$s" == "prefix"* ]]` |
| `s.ends_with(suffix)` | Check if string ends with suffix | `[[ "$s" == *"suffix" ]]` |
| `s.split(sep)` | Split string into array | `IFS='sep' read -ra arr <<< "$s"` |
| `s.join(sep)` | Join array elements | `$(IFS='sep'; echo "${s[*]}")` |
| `s.length()` | Get string length | `${#s}` |

`.contains()`, `.starts_with()`, and `.ends_with()` return conditions for use in `if`/`while`.
`.split()` produces an array that can be indexed:

```
name = "hello world"
result = name.replace("world", "langz")

if name.contains("hello") {
    print("found")
}

csv = "alice,bob,charlie"
parts = csv.split(",")
print(parts[1])  // "bob"

items = ["a", "b", "c"]
joined = items.join("-")  // "a-b-c"

n = name.length()  // 11
```

## Networking

| Function | Description |
|----------|-------------|
| `fetch(url, ...)` | HTTP request via curl ([details](language/networking.md)) |
| `json_get(data, path)` | Extract JSON value via jq |

### fetch() Keyword Arguments

| Kwarg | Description | Default |
|-------|-------------|---------|
| `method:` | HTTP method | GET |
| `body:` | Request body | none |
| `headers:` | Headers map | none |
| `timeout:` | Timeout in seconds | none |
| `retries:` | Retry count | none |

Sets convention variables: `_status`, `_body`, `_headers`.

## Date/Time

| Function | Description | Bash |
|----------|-------------|------|
| `timestamp()` | Unix timestamp | `$(date +%s)` |
| `date()` | Current date (YYYY-MM-DD) | `$(date +"%Y-%m-%d")` |
