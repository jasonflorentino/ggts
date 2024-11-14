#!/usr/bin/env bash

echo "Building Go"
GO_OUT_FILE="./build/main"
time go build -v -o $GO_OUT_FILE ./cmd
echo "Wrote to $GO_OUT_FILE"

echo ""
echo "Building CSS"
CSS_OUT_FILE="./build/static/style.css"
time ./tailwindcss -i ./css/style.css -o $CSS_OUT_FILE --minify
echo "Wrote to $CSS_OUT_FILE"

echo ""
echo "Copying views"
cp -r ./views ./build

echo ""
echo "Checking .env"
ENV_FILE_TPT="./.env.template"
ENV_FILE_OUT="./build/.env"

if [ ! -e $ENV_FILE_OUT ]; then
  echo "ERROR: No .env file found. Making new one."
  cp $ENV_FILE_TPT $ENV_FILE_OUT
  echo "Please ensure $ENV_FILE_OUT is correct before continuing"
  exit 1
fi

ENV_LC_SRC=$(wc -l < $ENV_FILE_TPT)
ENV_LC_OUT=$(wc -l < $ENV_FILE_OUT)

echo "$ENV_FILE_TPT lines: $ENV_LC_SRC | $ENV_FILE_OUT lines: $ENV_LC_OUT"

if [ "$ENV_LC_SRC" -eq "$ENV_LC_OUT" ]; then
  echo "No change in env line count."
else
  echo "ERROR: .env file lc change detected. Update $ENV_FILE_OUT first."
  echo "Diff:"
  echo ""
  diff $ENV_FILE_TPT $ENV_FILE_OUT 
  echo ""
  exit 1
fi

echo ""
echo "Restarting nginx"
systemctl restart nginx

echo ""
echo "Finished"