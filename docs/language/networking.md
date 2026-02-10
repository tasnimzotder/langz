# Networking

LangZ provides HTTP support via `fetch()` (backed by `curl`) and JSON parsing via `json_get()` (backed by `jq`).

## fetch()

### Simple GET

```
data = fetch("https://api.example.com/health")
print(data)
```

### POST with Options

`fetch()` supports keyword arguments for full HTTP control:

```
resp = fetch("https://api.example.com/users",
    method: "POST",
    body: payload,
    headers: {"Content-Type": "application/json"},
    timeout: 10,
    retries: 3
)
```

### Keyword Arguments

| Kwarg | Description | Default |
|-------|-------------|---------|
| `method:` | HTTP method (GET, POST, PUT, PATCH, DELETE) | GET |
| `body:` | Request body data | none |
| `headers:` | Request headers as map | none |
| `timeout:` | Max seconds to wait for response | none |
| `retries:` | Number of retry attempts on failure | none |

### Convention Variables

Every `fetch()` call sets three convention variables:

| Variable | Description |
|----------|-------------|
| `_status` | HTTP status code (e.g. `200`, `404`) |
| `_body` | Response body |
| `_headers` | Response headers |

```
fetch("https://api.example.com/users")

if _status == 200 {
    print("OK: {_body}")
} else {
    print("Failed with status {_status}")
}
```

!!! warning
    A second `fetch()` call overwrites `_status`, `_body`, and `_headers`. Save values to named variables if you need them across multiple requests.

### Error Handling

Use `or` for fallback values on failure:

```
data = fetch("https://api.example.com/data") or "unavailable"
```

When `retries:` is set, the request is retried up to N times with a 1-second delay between attempts. Retries stop on a 2xx response.

## json_get()

Extract values from JSON strings using jq path expressions:

```
data = fetch("https://api.example.com/user/1")
name = json_get(_body, ".name")
city = json_get(_body, ".address.city")
print("User: {name} from {city}")
```

!!! note
    `json_get()` requires `jq` to be installed on the target system.

## Full Example

```
// Create a user via API
payload = read("user.json")
resp = fetch("https://api.example.com/users",
    method: "POST",
    body: payload,
    headers: {"Content-Type": "application/json"},
    timeout: 30,
    retries: 2
)

if _status == 201 {
    user_id = json_get(_body, ".id")
    print("Created user: {user_id}")
} else {
    print("Failed: HTTP {_status}")
    exit(1)
}
```
