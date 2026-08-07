#!/bin/bash
set -e
echo "Running pre-commit checks..."

# Go vet
echo "→ Go vet..."
cd apps/api && go vet ./... && cd ../..

# Go build
echo "→ Go build..."
cd apps/api && go build ./... && cd ../..

# Web build
echo "→ Web build..."
pnpm --filter web build

echo "✅ All checks passed!"
