#!/bin/bash

# Exit on error
set -Eeu

# Ignore generated & testutils files
cat $1 | grep -Fv -e "/tests" > build/coverage/cover.out

go tool cover -func=build/coverage/cover.out | grep total:
# Generate coverage report in html format
go tool cover -html=build/coverage/cover.out -o $2

# Remove temporary file
rm build/coverage/cover.out || true
