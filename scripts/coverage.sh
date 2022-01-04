#!/bin/bash

# Exit on error
set -Eeu

# Ignore generated & testutils files
cat $1 | grep -Fv -e "/tests" -e "/mock" > $1

# Print total coverage
go tool cover -func=$1 | grep total:

# Generate coverage report in html format
go tool cover -html=$1 -o $2
