#!/bin/bash

# Exit script if any command fails
set -e

# Clean previous builds
echo "Cleaning previous builds..."
rm -f squirrel

# Build the project
echo "Building the project..."
go build -o squirrel .

# Run tests
echo "Running tests..."
go test ./...

# Output build status
echo "Build completed successfully!"