#!/bin/bash
set -e

cd /workspace

# Verify migrations are in sync with schema
atlas migrate diff --env gorm sync

# Run tests to verify the fix
go test ./... -v
