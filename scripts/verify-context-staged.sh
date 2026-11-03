#!/usr/bin/env bash
set -euo pipefail

if git diff --cached --name-only | grep -E '^backend/|^mobile/|^admin/|^docs/shared/' > /dev/null 2>&1; then
  if ! git diff --cached --name-only | grep -E '^CONTEXT.md$' > /dev/null 2>&1; then
    echo "ERROR: Atomic doc-commit rule violated. CONTEXT.md must be staged with code or spec changes."
    exit 1
  fi
  if ! git diff --cached --name-only | grep -E '^REQUIREMENTS_TRACKER.md$' > /dev/null 2>&1; then
    echo "ERROR: Atomic tracker rule violated. REQUIREMENTS_TRACKER.md must be staged with code or spec changes."
    exit 1
  fi
fi

exit 0
