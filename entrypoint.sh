#!/bin/bash

mkdir -p $(dirname "$PRIVATE_KEY_PATH")

if [ ! -f "$PRIVATE_KEY_PATH" ] || [ ! -f "$PUBLIC_KEY_PATH" ]; then
    echo "Generating RSA keys..."
    openssl genrsa -out "$PRIVATE_KEY_PATH" 2048
    openssl rsa -in "$PRIVATE_KEY_PATH" -pubout -out "$PUBLIC_KEY_PATH"
fi

exec "$@"
