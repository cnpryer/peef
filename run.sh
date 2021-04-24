#!/bin/bash

# Abort if error
set -e

cd "$(dirname "$0")"

trap 'kill $(jobs -p)' SIGINT SIGTERM EXIT

# echo "Installing dependencies"

while true; do

    go run main.go -token="$1" -guild="$2" -rmcmd -debug &
    
    inotifywait -r -e modify --exclude '(peef-bot.log|.git/)' .

    echo 'inotifywait ended'

    kill $(jobs -p) || true

    echo 'restarting...'

    sleep 1
done