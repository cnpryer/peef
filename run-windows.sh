#!/bin/bash

# Abort if error
set -e

go run main.go -token="$TOKEN" -guild="$GUILD" -rmcmd -debug