#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 go build -v -o ./out/main ./cmd

# You need tailwind CLI:
# https://tailwindcss.com/blog/standalone-cli
./tailwindcss -i ./css/style.css -o ./static/style.css --minify