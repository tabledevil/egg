#!/bin/bash
echo "Packing game data..."
go run cmd/packer/main.go

echo "Building binary..."
go build -o ctf-tool main.go

echo "Done! Run ./ctf-tool to play."
