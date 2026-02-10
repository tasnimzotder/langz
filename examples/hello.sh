#!/bin/bash
set -euo pipefail

name="world"
echo "Hello, Langz!"
echo "$name"
_fetch_attempt=0
_fetch_max=2
while [ "$_fetch_attempt" -lt "$_fetch_max" ]; do
  _fetch_attempt=$((_fetch_attempt + 1))
  _tmp_headers=$(mktemp)
  _tmp_body=$(mktemp)
  _status=$(curl -s -w "%{http_code}" --max-time 5000 -D "$_tmp_headers" -o "$_tmp_body" "https://jsonplaceholder.typicode.com/todos/1") || true
  _body=$(cat "$_tmp_body")
  _headers=$(cat "$_tmp_headers")
  rm -f "$_tmp_headers" "$_tmp_body"
  if [ "$_status" -ge 200 ] && [ "$_status" -lt 300 ]; then
    break
  fi
  sleep 1
done
resp="$_body"
echo "$resp"
echo "$_status"
echo "$_headers"
echo "$_body"
json_get "$_body"
