#!/bin/bash

set -e

go get -u github.com/golang/dep/cmd/dep
cd "$GOPATH/src/github.com/IngCr3at1on/glas-cli"
dep ensure