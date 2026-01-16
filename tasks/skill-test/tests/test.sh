#!/bin/bash
set -e

cd /workspace

# Check if greeting file was created
if [ -f "greeting.txt" ]; then
    echo "Success: greeting.txt was created"
    cat greeting.txt
    echo "1" > /logs/verifier/reward.txt
else
    echo "Error: greeting.txt not found"
    echo "0" > /logs/verifier/reward.txt
fi
