#!/bin/bash
set -euo pipefail

app="myapp"
version="1.0"
greet() {
  local name="$1"
  echo "Hello ${name} from ${app} v${version}"
}
greet "world"
if [ -e "config.json" ]; then
  echo "Config found"
else
  echo "No config, using defaults"
fi
mkdir -p "build/output"
echo "$version" > "build/output/version.txt"
