#!/bin/bash
set -e

declare -a arr=("github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest")

echo "Installing dependencies..."

for i in "${arr[@]}"
do
   echo "Installing $i"
   go install "$i"
#   go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
done

echo "DONE"