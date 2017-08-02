#!/usr/bin/env bash

# go test doesn't support colleting coverage across multiple packages,following script is workaround for that 
# source:  https://github.com/codecov/example-go

set -e

echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done