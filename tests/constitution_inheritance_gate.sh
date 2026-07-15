#!/usr/bin/env bash
# constitution_inheritance_gate.sh — pre-build / pre-merge gate that
# verifies the Helix Constitution inheritance is REAL, not decorative.
#
# Invariants (all must hold; each emits positive evidence on PASS per
# Constitution §7.1 "no false-success"):
#   Inv1  constitution/ present (Constitution.md readable)
#   Inv2  Constitution.md carries the §11.4 forensic-anchor HEADING line
#         (the full heading, not the bare substring — the TOC also
#         contains the phrase, so matching only the substring would make
#         this a bluff gate that the paired mutation could not flip).
#   Inv3  constitution/CLAUDE.md carries the MANDATORY ANTI-BLUFF COVENANT
#   Inv4  constitution/AGENTS.md carries an anti-bluff covenant reference
#   Inv5  parent CLAUDE.md / AGENTS.md / project constitution each point
#         at constitution/
#
# Paired mutation (Constitution §1.1): tests/meta_test_constitution_inheritance.sh
# Exit: 0 = all invariants verified; non-zero = at least one FAIL.
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel 2>/dev/null || echo "$(dirname "$SCRIPT_DIR")")"

# Locate the constitution robustly (works from any nested depth).
CONSTITUTION_DIR=""
if [ -x "$REPO_ROOT/constitution/find_constitution.sh" ]; then
    CONSTITUTION_DIR="$(bash "$REPO_ROOT/constitution/find_constitution.sh" 2>/dev/null || true)"
fi
[ -z "$CONSTITUTION_DIR" ] && CONSTITUTION_DIR="$REPO_ROOT/constitution"

# Exact sentinel heading — MUST match constitution/meta_test_inheritance.sh
SENTINEL='### §11.4 End-user quality guarantee — forensic anchor (User mandate, 2026-04-28)'

fail=0
pass() { echo "PASS: $1"; }
bad()  { echo "FAIL: $1" >&2; fail=1; }

# Inv1
if [ -f "$CONSTITUTION_DIR/Constitution.md" ]; then
    pass "Inv1 constitution present ($CONSTITUTION_DIR)"
else
    bad  "Inv1 constitution/Constitution.md not found (looked in $CONSTITUTION_DIR)"
fi

# Inv2 — full heading line, not bare substring
if grep -qF -- "$SENTINEL" "$CONSTITUTION_DIR/Constitution.md" 2>/dev/null; then
    pass "Inv2 §11.4 forensic-anchor heading present"
else
    bad  "Inv2 §11.4 forensic-anchor heading missing from Constitution.md"
fi

# Inv3
if grep -qF -- 'MANDATORY ANTI-BLUFF COVENANT' "$CONSTITUTION_DIR/CLAUDE.md" 2>/dev/null; then
    pass "Inv3 CLAUDE.md MANDATORY ANTI-BLUFF COVENANT anchor present"
else
    bad  "Inv3 CLAUDE.md MANDATORY ANTI-BLUFF COVENANT anchor missing"
fi

# Inv4 (case-insensitive — the anchor appears in mixed case)
if grep -qiF -- 'anti-bluff covenant' "$CONSTITUTION_DIR/AGENTS.md" 2>/dev/null; then
    pass "Inv4 AGENTS.md anti-bluff covenant anchor present"
else
    bad  "Inv4 AGENTS.md anti-bluff covenant anchor missing"
fi

# Inv5 — parent files reference the submodule
for f in CLAUDE.md AGENTS.md docs/guides/HELIX_SKILLS_CONSTITUTION.md; do
    if grep -qF -- 'constitution/' "$REPO_ROOT/$f" 2>/dev/null; then
        pass "Inv5 $f references constitution/"
    else
        bad  "Inv5 $f does not reference constitution/ (or is missing)"
    fi
done

echo
if [ "$fail" -ne 0 ]; then
    echo "GATE: FAIL — constitution inheritance is NOT satisfied"
    exit 1
fi
echo "GATE: PASS — all 7 invariant checks verified with positive evidence"
exit 0
