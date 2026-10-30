#!/usr/bin/env bash
set -euo pipefail

if git diff --cached --name-only | grep -E '^backend/|^mobile/|^admin/|^docs/shared/' > /dev/null 2>&1; then
  if ! git diff --cached --name-only | grep -E '^CONTEXT.md$' > /dev/null 2>&1; then
    echo "ERROR: Atomic doc-commit rule violated. CONTEXT.md must be staged with code or spec changes."
    exit 1
  fi
fi

exit 0
