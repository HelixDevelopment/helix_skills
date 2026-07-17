#!/bin/bash
# scripts/regenerate_from_db.sh — Regenerate tracked documents from workable_items.db
# Constitution §11.4.12 (auto-generated docs sync) + §11.4.93 (DB as SSoT)
#
# This script proves the DB → docs pipeline works:
#   1. workable-items db-to-md regenerates the item listings from DB
#   2. export_docs.sh generates HTML + PDF from the regenerated markdown
#
# Usage: bash scripts/regenerate_from_db.sh
# Dependencies: workable-items binary, pandoc, weasyprint
set -euo pipefail

PROJ="$(cd "$(dirname "$0")/.." && pwd)"
DB="${PROJ}/workable_items.db"
BIN="${PROJ}/../../constitution/scripts/workable-items/cmd/workable-items/workable-items"

# Build binary if not present
if [ ! -x "$BIN" ]; then
  echo "Building workable-items binary..."
  (cd "$(dirname "$BIN")" && go build -o workable-items . 2>/dev/null) || {
    # Fallback: use the one in /tmp
    BIN="/tmp/workable-items"
    if [ ! -x "$BIN" ]; then
      echo "ERROR: workable-items binary not found"
      exit 1
    fi
  }
fi

DB="${PROJ}/workable_items.db"
if [ ! -f "$DB" ]; then
  echo "ERROR: ${DB} not found"
  exit 1
fi

echo "=== Regenerating from DB: ${DB} ==="
echo "Items: $(sqlite3 "$DB" "SELECT COUNT(*) FROM items;")"
echo ""

# Step 1: Validate DB
echo "--- Step 1: Validate DB ---"
$BIN validate --db "$DB" 2>&1
echo ""

# Step 2: Generate summary
echo "--- Step 2: Status summary ---"
sqlite3 "$DB" "SELECT status, COUNT(*) as cnt FROM items GROUP BY status ORDER BY cnt DESC;"
echo ""
sqlite3 "$DB" "SELECT type, COUNT(*) as cnt FROM items GROUP BY type ORDER BY cnt DESC;"
echo ""

# Step 3: Export to markdown (proves DB → md pipeline)
echo "--- Step 3: Export to markdown ---"
EXPORT_MD="${PROJ}/docs/exports/items_from_db.md"
mkdir -p "$(dirname "$EXPORT_MD")"
$BIN sync db-to-md --db "$DB" --out-issues "$EXPORT_MD" 2>&1
echo "Exported: ${EXPORT_MD}"
echo ""

# Step 4: Generate HTML + PDF from exported markdown
echo "--- Step 4: Generate HTML + PDF ---"
bash "${PROJ}/scripts/export_docs.sh" 2>&1 | tail -5
echo ""

echo "=== Regeneration complete ==="
