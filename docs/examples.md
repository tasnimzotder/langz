# Examples

## Deployment Script

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

## Log File Cleanup

```
max_age = 7
count = 0

for f in glob("/var/log/app/*.log") {
    age = exec("find {f} -mtime +{max_age} -print")
    if age != "" {
        rm(f)
        count = count + 1
    }
}

print("Cleaned {count} old log files")
```

## Health Check with Retry

```
fn check_health(url: str, max_retries: int) {
    data = fetch(url, timeout: 5, retries: max_retries) or "failed"

    if _status == 200 {
        print("OK: {url}")
    } else {
        print("FAIL: {url} (status {_status})")
        exit(1)
    }
}

check_health("https://api.example.com/health", 3)
check_health("https://cdn.example.com/ping", 2)
```

## API Integration

```
// Fetch user data and extract fields
fetch("https://api.example.com/user/1", timeout: 10)

if _status == 200 {
    name = json_get(_body, ".name")
    email = json_get(_body, ".email")
    print("User: {name} ({email})")
} else {
    print("API error: {_status}")
}
```

## System Info Report

```
fn report() {
    _os = os()
    _arch = arch()
    host = hostname()
    user = whoami()
    ts = timestamp()

    print("=== System Report ===")
    print("OS:       {_os}")
    print("Arch:     {_arch}")
    print("Host:     {host}")
    print("User:     {user}")
    print("Time:     {ts}")

    mkdir("reports")
    write("reports/system.txt", "os={_os} arch={_arch} host={host}")
}

report()
```

## File Backup

```
src = env("BACKUP_SRC") or "/etc/nginx"
dst = env("BACKUP_DST") or "/tmp/backups"
today = date()

if is_dir(src) {
    backup_dir = "{dst}/{today}"
    mkdir(backup_dir)
    exec("cp -r {src}/* {backup_dir}/")
    print("Backed up {src} to {backup_dir}")
} else {
    print("Source not found: {src}")
    exit(1)
}
```
