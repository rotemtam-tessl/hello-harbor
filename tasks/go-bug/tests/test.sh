#!/bin/bash
set -e

cd /workspace

# Run Go tests
if go test ./... -v; then
    echo "Success: All tests passed"
    echo "1" > /logs/verifier/reward.txt
else
    echo "Error: Tests failed"
    echo "0" > /logs/verifier/reward.txt
fi
