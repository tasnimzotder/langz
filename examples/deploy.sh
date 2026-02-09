#!/bin/bash
set -euo pipefail

app="webapp"
version="2.1.0"
platform=$(uname -s | tr '[:upper:]' '[:lower:]')
host=$(hostname)
user=$(whoami)
log() {
  local msg="$1"
  echo "[deploy] ${msg}"
}
log "Starting deployment of ${app} v${version}"
log "Platform: ${platform}, Host: ${host}, User: ${user}"
rm -rf "dist"
mkdir -p "dist"
mkdir -p "dist/assets"
echo "app=${app}" > "dist/manifest.txt"
echo "version=${version}" >> "dist/manifest.txt"
echo "built_by=${user}" >> "dist/manifest.txt"
if [ -e "config.json" ]; then
  log "Config found, copying"
  cp "config.json" "dist/config.json"
else
  log "No config, using defaults"
  echo "{}" > "dist/config.json"
fi
assets=("style.css" "app.js" "index.html")
for asset in "${assets[@]}"; do
  log "Bundling ${asset}"
done
env_name="${DEPLOY_ENV:-staging}"
case "$env_name" in
  production)
    log "PRODUCTION deploy"
    ;;
  staging)
    log "Staging deploy"
    ;;
  *)
    log "Unknown environment: ${env_name}"
    ;;
esac
if [ -d "dist" ]; then
  log "Build directory ready"
else
  log "Build failed!"
  exit 1
fi
log "Deployment complete!"
