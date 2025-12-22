#!/bin/bash

protoRoot="$(pwd)/proto"
serviceProtoRoot="$protoRoot/user-service"
protoFiles=$(find "$serviceProtoRoot" -name "*.proto")
for file in $protoFiles; do
    relativePath="${file#$protoRoot/}"
    protoc -I="$protoRoot" \
           --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative "$relativePath"
done
