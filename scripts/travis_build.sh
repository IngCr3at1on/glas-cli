#!/bin/bash

set -e

pkg='github.com/IngCr3at1on/glas-cli'

cd "$GOPATH/src/$pkg"

env GOOS=linux GOARCH=amd64 go build -o glas-cli_linux64 "$pkg"
env GOOS=linux GOARCH=386 go build -o glas-cli_linux386 "$pkg"
env GOOS=darwin GOARCH=amd64 go build -o glas-cli_darwin64 "$pkg"
env GOOS=darwin GOARCH=386 go build -o glas-cli_darwin386 "$pkg"
#TODO: figure out why thse won't build...
#env GOOS=windows GOARCH=amd64 go build -o glas-cli_windows64 "$pkg"
#env GOOS=windows GOARCH=386 go build -o glas-cli_windows386 "$pkg"
