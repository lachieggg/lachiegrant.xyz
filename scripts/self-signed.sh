#!/bin/sh

current_dir="$(basename "$PWD")"

# Check if we are in the root project folder
if [ "$current_dir" = "scripts" ] || [ ! -f "Makefile" ]; then
    echo "Must be run from the root project folder."
    exit;
fi

# Generate a new SSL certificate and key 
openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
    -keyout ./tls/nginx/keys/privkey.pem \
    -out ./tls/nginx/certs/fullchain.pem
