#!/usr/bin/env bash

set -e
set -u
set -o pipefail

echo
echo "Checking .env file"
ENV_FILE_TPT="./.env.template"
ENV_FILE_OUT="./build/.env"

if [ ! -e $ENV_FILE_OUT ]; then
  echo "ERROR: No .env file found. Making new one."
  cp $ENV_FILE_TPT $ENV_FILE_OUT
  echo "Please ensure $ENV_FILE_OUT is correct before continuing"
  exit 1
else
  echo ".env file found"
fi

ENV_LC_SRC=$(wc -l < $ENV_FILE_TPT)
ENV_LC_OUT=$(wc -l < $ENV_FILE_OUT)

echo "$ENV_FILE_TPT lines: $ENV_LC_SRC | $ENV_FILE_OUT lines: $ENV_LC_OUT"

if [ "$ENV_LC_SRC" -eq "$ENV_LC_OUT" ]; then
  echo "No change in env line count."
else
  echo "ERROR: .env file lc change detected. Update $ENV_FILE_OUT first."
  echo "Diff:"
  echo
  diff $ENV_FILE_TPT $ENV_FILE_OUT
  echo
  exit 1
fi

echo
echo "Building Go"
GO_OUT_FILE="./build/main"
time go build -v -o $GO_OUT_FILE ./cmd
echo "Wrote to $GO_OUT_FILE"

echo
echo "Building CSS"
CSS_OUT_FILE="./build/static/style.css"
time ./tailwindcss -i ./css/style.css -o $CSS_OUT_FILE --minify
echo "Wrote to $CSS_OUT_FILE"

echo
echo "Copying views"
cp -r ./views ./build

echo
echo "Checking log file"
LOGFILE_K=GGTS_LOGFILE
LOGFILE_V=$(grep "$LOGFILE_K" "$ENV_FILE_OUT" | sed -E "s/$LOGFILE_K=//")
LOGFILE_P=./build/$LOGFILE_V

if [ ! -n "$LOGFILE_V" ]; then
  echo "ERROR: no value for .env logfile key $LOGFILE_K"
  exit 1
fi

if [ ! -e "$LOGFILE_P" ]; then
  echo "No file $LOGFILE_P. Creating it..."
  touch $LOGFILE_P
  chmod 666 $LOGFILE_P
else
  echo "Existing log file found."
fi

echo
echo "Checking env"
ENV_K=GGTS_ENV
ENV_V=$(grep "$ENV_K" "$ENV_FILE_OUT" | sed -E "s/$ENV_K=//")
echo "env: $ENV_V"
if [ "$ENV_V" == "production" ]; then
  echo "Restarting ggts.service"
  systemctl restart ggts.service
else
  echo "No process to restart"
fi

echo
echo "Finished"
echo
