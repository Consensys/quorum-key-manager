#!/bin/bash

# Exit on error
set -Eeu

# Ignore generated & testutils files
cat $1 | grep -Fv -e "_mock.go" -e "/tests" -e "/mock" -e "/testutils" -e "/testdata" > "$1.tmp"

# Print total coverage
go tool cover -func="$1.tmp" | grep total:

# Generate coverage report in html format
go tool cover -html="$1.tmp" -o $2

cat "$1.tmp" > $1

rm "$1.tmp"
