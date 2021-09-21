#!/bin/bash

# Exit on error
set -Eeu

mkdir -p build/coverage
go test -covermode=count -coverprofile build/coverage/profile.out "$@"

# Ignore generated & testutils files
cat build/coverage/profile.out | grep -Fv -e "/testutils" -e "/integration-tests" >build/coverage/cover.out

# Generate coverage report in html format
go tool cover -func=build/coverage/cover.out | grep total:
go tool cover -html=build/coverage/cover.out -o build/coverage/coverage.html

# Remove temporary file
rm build/coverage/profile.out build/coverage/cover.out || true
