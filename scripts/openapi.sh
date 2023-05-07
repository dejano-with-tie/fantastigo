#!/bin/bash
set -e

readonly spec_name="$1"

readonly srv_output_dir="internal/$spec_name/server"
readonly client_output_dir="pkg/http/$spec_name"

mkdir -p "$srv_output_dir"
mkdir -p "$client_output_dir"

echo "Generating echo server stub and a client from openapi spec..."
# generate server types & echo server stub
## This can be defined in a config yaml. Generally speaking, yaml is more user friendly & elegant
## TODO Use yaml for args
oapi-codegen -generate types -o "$srv_output_dir/openapi_types.gen.go" -package "server" "docs/openapi/$spec_name.yaml"
oapi-codegen -generate server -o "$srv_output_dir/openapi_api.gen.go" -package "server" "docs/openapi/$spec_name.yaml"
#oapi-codegen -generate spec -o "$srv_output_dir/openapi_spec.gen.go" -package "server" "docs/openapi/$spec_name.yaml"

# generate client types & client
oapi-codegen -generate types -o "$client_output_dir/openapi_types.gen.go" -package "$spec_name" "docs/openapi/$spec_name.yaml"
oapi-codegen -generate client -o "$client_output_dir/openapi_client.gen.go" -package "$spec_name" "docs/openapi/$spec_name.yaml"

echo "DONE"