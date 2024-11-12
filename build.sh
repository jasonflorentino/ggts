#!/usr/bin/env bash

echo "Building Go"
GO_OUT_FILE="./out/main"
time GOOS=linux GOARCH=amd64 go build -v -o $GO_OUT_FILE ./cmd
echo "Wrote to $GO_OUT_FILE"

echo ""
echo "Building CSS"
CSS_OUT_FILE="./static/style.css"
# You need tailwind CLI:
# https://tailwindcss.com/blog/standalone-cli
time ./tailwindcss -i ./css/style.css -o $CSS_OUT_FILE --minify
echo "Wrote to $CSS_OUT_FILE"

echo ""
echo "Finished"