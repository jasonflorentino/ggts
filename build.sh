#!/usr/bin/env bash

go build -o ./out/main ./cmd

# You need tailwind CLI:
# https://tailwindcss.com/blog/standalone-cli
./tailwindcss -i ./css/style.css -o ./static/style.css --minify