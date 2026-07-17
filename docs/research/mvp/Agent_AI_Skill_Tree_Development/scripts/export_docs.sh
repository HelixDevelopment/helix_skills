#!/bin/bash
# scripts/export_docs.sh — Generate HTML + PDF exports of tracked documents
# Constitution §11.4.12 (auto-generated docs sync) + §11.4.65 (Markdown export)
#
# Usage: bash scripts/export_docs.sh
# Output: docs/exports/<basename>.html + docs/exports/<basename>.pdf
# Dependencies: pandoc, weasyprint (optional for PDF)
set -euo pipefail

PROJ="$(cd "$(dirname "$0")/.." && pwd)"
EXPORT_DIR="${PROJ}/docs/exports"
mkdir -p "$EXPORT_DIR"

# Tracked documents to export (§11.4.65 scope)
DOCS=(
  "GAPS_AND_RISKS_REGISTER.md"
  "CONTINUATION.md"
  "REQUIREMENTS.md"
  "IMPLEMENTATION_PLAN.md"
  "SPEC.md"
)

# HTML template with basic styling
HTML_HEAD='<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 900px; margin: 2em auto; padding: 0 1em; line-height: 1.6; color: #333; }
h1 { border-bottom: 2px solid #333; padding-bottom: 0.3em; }
h2 { border-bottom: 1px solid #ccc; padding-bottom: 0.2em; margin-top: 1.5em; }
table { border-collapse: collapse; width: 100%; margin: 1em 0; }
th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
th { background: #f5f5f5; }
code { background: #f0f0f0; padding: 2px 4px; border-radius: 3px; font-size: 0.9em; }
pre { background: #f5f5f5; padding: 1em; overflow-x: auto; border-radius: 4px; }
a { color: #0366d6; }
.status-fixed { color: #28a745; font-weight: bold; }
.status-queued { color: #6c757d; }
.status-blocked { color: #dc3545; font-weight: bold; }
</style></head><body>'

export_count=0
fail_count=0

for doc in "${DOCS[@]}"; do
  src="${PROJ}/${doc}"
  if [ ! -f "$src" ]; then
    echo "SKIP: ${doc} (not found)"
    continue
  fi

  base="${doc%.md}"
  html_out="${EXPORT_DIR}/${base}.html"
  pdf_out="${EXPORT_DIR}/${base}.pdf"

  # Generate HTML
  echo -n "HTML: ${doc} -> ${base}.html ... "
  if pandoc "$src" \
    --from markdown --to html5 \
    --standalone \
    --metadata title="${base}" \
    --output "$html_out" 2>/dev/null; then
    echo "OK"
    export_count=$((export_count + 1))
  else
    echo "FAIL"
    fail_count=$((fail_count + 1))
    continue
  fi

  # Generate PDF (via weasyprint from HTML)
  if command -v weasyprint >/dev/null 2>&1; then
    echo -n "PDF: ${doc} -> ${base}.pdf ... "
    if weasyprint "$html_out" "$pdf_out" 2>/dev/null; then
      echo "OK"
      export_count=$((export_count + 1))
    else
      echo "FAIL (weasyprint)"
      fail_count=$((fail_count + 1))
    fi
  else
    echo "PDF: SKIP (weasyprint not installed)"
  fi
done

echo ""
echo "=== Export complete: ${export_count} files generated, ${fail_count} failures ==="
echo "Output: ${EXPORT_DIR}/"
ls -la "$EXPORT_DIR/" 2>/dev/null

exit $fail_count
