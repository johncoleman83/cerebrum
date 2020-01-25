#!/usr/bin/env bash

function execute() {
    set -e

    COVERAGE_TXT="./test/coverage.txt"
    echo "" > $COVERAGE_TXT

    for d in $(go list ../... | grep -v -e pkg/utl/mock); do
        go test -race -coverprofile=profile.out -covermode=atomic "$d"
        if [ -f profile.out ]; then
            cat profile.out >> $COVERAGE_TXT
            rm profile.out
        fi
    done
}

# don't execute if being sourced
if [[ "$0" = "$BASH_SOURCE" ]]; then
    execute
fi
