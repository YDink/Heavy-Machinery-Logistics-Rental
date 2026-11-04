#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# scripts/fixtures/apply-fixtures.sh
# Applies all idempotent database test fixtures for multi-machine parity (macOS/Linux)
# ==============================================================================

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-dink_machinery}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "==> Synchronizing Idempotent Database Test Fixtures [Target: ${DB_NAME}@${DB_HOST}:${DB_PORT}]..."

export PGPASSWORD="${DB_PASSWORD}"

for fixture_file in "${SCRIPT_DIR}"/*.sql; do
  [ -f "${fixture_file}" ] || continue
  echo "    --> Applying $(basename "${fixture_file}")..."
  psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f "${fixture_file}" > /dev/null
done

echo "==> Database Test Fixtures Successfully Synchronized (Multi-Machine Parity Guaranteed)."
