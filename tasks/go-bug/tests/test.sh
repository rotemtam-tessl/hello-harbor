#!/bin/bash
set -e

cd /workspace

# Expected hashes for the original migrations (must not be modified)
EXPECTED_HASH_1="20260114131741.sql h1:FogWQlVt2p7AeEsYreAlFOOdSOvgcsSf9D+pUShdAS4="
EXPECTED_HASH_2="20260114133748_add_categories.sql h1:pgqzTNd3JvZJDwR24vtCpSgjokg/u9G0mtqvXBnPJr4="
EXPECTED_HASH_3="20260114133824_seed_categories.sql h1:1H3sphk8Y0mul4zOiG2WQhhkD6MQ4Q0+Hg3o9Ms3t14="

# Check that original migrations were not modified (linear history preserved)
echo "Checking migration history integrity..."
if ! grep -qF "$EXPECTED_HASH_1" migrations/atlas.sum; then
    echo "Error: Original migration 20260114131741.sql was modified"
    echo "0" > /logs/verifier/reward.txt
    exit 0
fi
if ! grep -qF "$EXPECTED_HASH_2" migrations/atlas.sum; then
    echo "Error: Original migration 20260114133748_add_categories.sql was modified"
    echo "0" > /logs/verifier/reward.txt
    exit 0
fi
if ! grep -qF "$EXPECTED_HASH_3" migrations/atlas.sum; then
    echo "Error: Original migration 20260114133824_seed_categories.sql was modified"
    echo "0" > /logs/verifier/reward.txt
    exit 0
fi
echo "Migration history intact."

# Validate migrations with Atlas
echo "Validating migrations..."
if ! atlas migrate validate --env gorm; then
    echo "Error: Migration validation failed"
    echo "0" > /logs/verifier/reward.txt
    exit 0
fi
echo "Migrations validated."

# Run Go tests
echo "Running tests..."
if go test ./... -v; then
    echo "Success: All tests passed"
    echo "1" > /logs/verifier/reward.txt
else
    echo "Error: Tests failed"
    echo "0" > /logs/verifier/reward.txt
fi
