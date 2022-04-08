#!/usr/bin/env bash

# Checking go format
echo "==> Checking Lumen provider code compiles with gofmt requirements..."
gofmt_files=$(gofmt -l `find ./lumen -name '*.go'`)
if [[ -n ${gofmt_files} ]]; then
    echo 'gofmt running on following files:'
    echo "${gofmt_files}"
    echo "Use command: 'make fmt' to reformat code"
    exit 1
fi

exit 0
