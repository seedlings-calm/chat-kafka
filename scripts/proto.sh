#!/usr/bin/env bash
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
    protoc \
    -I "proto" \
    -I "third_party" \
    --go_out=. \
      --go-grpc_out=. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')

done