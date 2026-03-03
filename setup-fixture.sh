#!/usr/bin/env bash
# setup-fixture.sh — Clones fixture at pinned SHA, initializes claude-memory, indexes it.
# Usage: bash setup-fixture.sh [FIXTURE_DIR]
# Outputs the project ID on the last line.
set -euo pipefail

FIXTURE_REPO="https://github.com/Intenteon/claude-memory-fixture"
FIXTURE_SHA="e4c954ce4783ca37bb1c5f2e1a45825fbed8e93f"
FIXTURE_DIR="${1:-/tmp/claude-memory-e2e-fixture}"

echo "=== claude-memory E2E Fixture Setup ==="
echo "Repo:    $FIXTURE_REPO"
echo "SHA:     $FIXTURE_SHA"
echo "Dir:     $FIXTURE_DIR"
echo

# Check if already at correct SHA
if [ -d "$FIXTURE_DIR/.git" ]; then
  CURRENT_SHA=$(git -C "$FIXTURE_DIR" rev-parse HEAD 2>/dev/null || echo "")
  if [ "$CURRENT_SHA" = "$FIXTURE_SHA" ]; then
    echo "Fixture already at correct SHA."

    # Check if already indexed
    PROJECT_ID=$(claude-memory project list 2>/dev/null | grep -A1 -i "fixture" | grep "ID:" | awk '{print $2}' || echo "")
    if [ -n "$PROJECT_ID" ]; then
      echo "Already indexed: $PROJECT_ID"
      echo "$PROJECT_ID"
      exit 0
    fi
    echo "Not yet indexed, proceeding..."
  else
    echo "SHA mismatch (have: $CURRENT_SHA), re-cloning..."
    rm -rf "$FIXTURE_DIR"
  fi
fi

# Clone if needed
if [ ! -d "$FIXTURE_DIR" ]; then
  echo "Cloning fixture repo..."
  git clone "$FIXTURE_REPO" "$FIXTURE_DIR"
  git -C "$FIXTURE_DIR" checkout "$FIXTURE_SHA"
fi

# Initialize claude-memory
echo "Initializing claude-memory..."
cd "$FIXTURE_DIR"
claude-memory init --force -y

# Index (structure + call graph)
echo "Indexing (structure pass)..."
claude-memory index --no-embeddings --workers 4

# Index (embeddings pass)
echo "Indexing (embeddings pass)..."
claude-memory index --workers 4 || echo "Warning: embedding pass failed (ollama may not be running)"

# Discover project ID
echo
echo "Discovering project ID..."
PROJECT_ID=$(claude-memory project list 2>/dev/null | grep -A1 -i "fixture" | grep "ID:" | awk '{print $2}' || echo "")

if [ -z "$PROJECT_ID" ]; then
  echo "ERROR: Could not discover fixture project ID"
  echo "Run 'claude-memory project list' to check"
  exit 1
fi

echo
echo "=== Setup Complete ==="
echo "Project ID: $PROJECT_ID"
echo
echo "Run tests with:"
echo "  FIXTURE_PROJECT_ID=$PROJECT_ID go test ./tests/... -v"
echo
echo "$PROJECT_ID"
