#!/usr/bin/env bash
# meta_test_constitution_inheritance.sh — proves the inheritance gate is
# NOT a bluff gate (Constitution §1.1 false-positive immunity).
#
# Delegates to the constitution-side generic harness
# constitution/meta_test_inheritance.sh, which:
#   1. snapshots constitution/Constitution.md
#   2. strips the §11.4 forensic-anchor heading (the regression)
#   3. runs our project gate
#   4. asserts the gate now FAILS
#   5. restores Constitution.md
#
# Exit 0 here == the gate correctly caught the mutation.
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel 2>/dev/null || echo "$(dirname "$SCRIPT_DIR")")"

META="$REPO_ROOT/constitution/meta_test_inheritance.sh"
GATE_CMD="bash $REPO_ROOT/tests/constitution_inheritance_gate.sh"

if [ ! -f "$META" ]; then
    echo "FAIL: constitution/meta_test_inheritance.sh not found at $META" >&2
    exit 2
fi

echo "Meta-test: mutating Constitution.md and asserting the gate FAILS..."
bash "$META" "$GATE_CMD"
