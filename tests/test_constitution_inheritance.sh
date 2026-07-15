#!/usr/bin/env bash
# test_constitution_inheritance.sh — comprehensive host-side inheritance
# test. Asserts all gate invariants, the recursive child-submodule
# inheritance pointers, and the false-positive-immunity meta-test.
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel 2>/dev/null || echo "$(dirname "$SCRIPT_DIR")")"
cd "$REPO_ROOT" || exit 2
overall=0

echo "== [1/3] Inheritance gate (invariants 1-5) =="
if bash tests/constitution_inheritance_gate.sh; then echo "  -> PASS"; else echo "  -> FAIL"; overall=1; fi
echo

echo "== [2/3] Recursive nested-submodule inheritance pointers =="
# Enumerate owned consuming submodules. Exclude the constitution source
# submodule and anything the constitution itself owns (constitution/...).
mapfile -t subs < <(git submodule status --recursive 2>/dev/null \
                    | awk '{print $2}' \
                    | grep -vE '^constitution(/|$)' || true)
if [ "${#subs[@]}" -eq 0 ]; then
    echo "  0 owned consuming submodules found (only the 'constitution' source"
    echo "  submodule is present) -> PASS (vacuously true, reported explicitly)"
else
    for s in "${subs[@]}"; do
        ok=1
        for f in CLAUDE.md AGENTS.md; do
            if [ -f "$s/$f" ] && grep -qiF 'constitution' "$s/$f"; then :; else
                echo "  MISSING inheritance pointer: $s/$f"; ok=0
            fi
        done
        if [ "$ok" -eq 1 ]; then echo "  $s -> PASS"; else echo "  $s -> FAIL"; overall=1; fi
    done
fi
echo

echo "== [3/3] False-positive-immunity meta-test (mutation) =="
if bash tests/meta_test_constitution_inheritance.sh; then echo "  -> PASS"; else echo "  -> FAIL"; overall=1; fi
echo

if [ "$overall" -eq 0 ]; then
    echo "RESULT: PASS — constitution inheritance verified (invariants + recursion + anti-bluff)"
else
    echo "RESULT: FAIL — see failures above"
fi
exit $overall
