#!/bin/bash
# Performance test runner for the project visualization API.
# Usage: ./run.sh [smoke|load|stress|all]

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
RESULTS_DIR="$(dirname "$0")/results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT="$RESULTS_DIR/report_${TIMESTAMP}.txt"

mkdir -p "$RESULTS_DIR"

check_backend() {
  echo "Checking backend at $BASE_URL ..."
  if ! curl -sf "$BASE_URL/health" > /dev/null 2>&1; then
    echo "ERROR: Backend is not running at $BASE_URL"
    echo "Run: make up && make db-migrate"
    exit 1
  fi
  echo "Backend is healthy."
}

run_test() {
  local name="$1"
  local script="$2"
  local out="$RESULTS_DIR/${name}_${TIMESTAMP}.json"

  echo ""
  echo "======================================"
  echo " Running: $name"
  echo "======================================"
  k6 run \
    --out "json=$out" \
    -e "BASE_URL=$BASE_URL" \
    "$script" 2>&1 | tee -a "$REPORT"
}

TARGET="${1:-all}"

check_backend

{
  echo "Performance Test Report"
  echo "Date: $(date)"
  echo "Target: $BASE_URL"
  echo "Test: $TARGET"
  echo "======================================="
} > "$REPORT"

case "$TARGET" in
  smoke)
    run_test "smoke" "$(dirname "$0")/scripts/smoke.js"
    ;;
  load)
    run_test "load" "$(dirname "$0")/scripts/load.js"
    ;;
  stress)
    run_test "stress" "$(dirname "$0")/scripts/stress.js"
    ;;
  all)
    run_test "smoke"  "$(dirname "$0")/scripts/smoke.js"
    run_test "load"   "$(dirname "$0")/scripts/load.js"
    run_test "stress" "$(dirname "$0")/scripts/stress.js"
    ;;
  *)
    echo "Usage: $0 [smoke|load|stress|all]"
    exit 1
    ;;
esac

echo ""
echo "======================================"
echo " All tests complete. Report: $REPORT"
echo "======================================"
