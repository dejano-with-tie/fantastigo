#!/bin/bash
set -e

readonly app_path="$1"
readonly output_app_path="$2"

echo "Building '$app_path' app with output '$output_app_path'..."

go build -o "$output_app_path" "$app_path"

echo "DONE"