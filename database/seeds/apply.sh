#!/usr/bin/env bash
# apply.sh - Apply seed data to the development/test database.
#
# Usage:
#   ./database/seeds/apply.sh
#   DATABASE_URL="postgres://..." ./database/seeds/apply.sh
#
# The DATABASE_URL environment variable takes precedence over the default value.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SEED_SQL="${SCRIPT_DIR}/seed.sql"

# Default connection settings (same as Makefile)
: "${DATABASE_URL:=postgres://admin:admin123@localhost:5432/project_visualization?sslmode=disable}"

echo "Applying seed data..."
echo "Database: ${DATABASE_URL}"
echo ""

if ! command -v psql &>/dev/null; then
  echo "Error: psql is not installed. Please install PostgreSQL client tools." >&2
  exit 1
fi

psql "${DATABASE_URL}" -f "${SEED_SQL}"

echo ""
echo "Seed data applied successfully."
