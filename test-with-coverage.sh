#!/usr/bin/env bash

# Reference: https://github.com/codecov/example-go

set -e
echo "" > cover.out

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> cover.out
        rm profile.out
    fi
done