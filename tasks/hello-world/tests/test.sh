#!/bin/bash
set -e

cd /workspace

# Check if hello.py exists
if [ ! -f "hello.py" ]; then
    echo "Error: hello.py not found"
    echo "0" > /logs/verifier/reward.txt
    exit 0
fi

# Run the script and capture output
OUTPUT=$(python hello.py 2>&1)

# Check if output matches exactly
if [ "$OUTPUT" = "Hello, World!" ]; then
    echo "Success: Output matches expected"
    echo "1" > /logs/verifier/reward.txt
else
    echo "Error: Expected 'Hello, World!' but got '$OUTPUT'"
    echo "0" > /logs/verifier/reward.txt
fi
