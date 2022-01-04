#!/bin/bash

# Exit on error
set -Eeu

# Ignore generated & testutils files
cat $1 | grep -Fv -e "/deps" -e "/tests" -e "/public" > build/coverage/cover.out

# Generate coverage report in html format
go tool cover -html=build/coverage/cover.out -o $2

# Remove temporary file
rm $1 build/coverage/cover.out || true
