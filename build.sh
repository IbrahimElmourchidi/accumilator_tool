#!/bin/bash

# Build script for accumilator tool

echo "Building accumilator for multiple platforms..."

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/accumilator-windows-amd64.exe main.go
echo "Built: accumilator-windows-amd64.exe"

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/accumilator-linux-amd64 main.go
echo "Built: accumilator-linux-amd64"

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o bin/accumilator-darwin-amd64 main.go
echo "Built: accumilator-darwin-amd64"

# Build for ARM64 (Linux and macOS)
GOOS=linux GOARCH=arm64 go build -o bin/accumilator-linux-arm64 main.go
echo "Built: accumilator-linux-arm64"

GOOS=darwin GOARCH=arm64 go build -o bin/accumilator-darwin-arm64 main.go
echo "Built: accumilator-darwin-arm64"

echo "All builds completed!"