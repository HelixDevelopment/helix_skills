#!/usr/bin/env bash
# pre_build_verification.sh — preflight gate wired ahead of build/merge.
# Fails fast if constitution inheritance is not real. The paired mutation
# that proves this gate is not a bluff gate lives in
# tests/meta_test_constitution_inheritance.sh (Constitution §1.1).
set -uo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec bash "$SCRIPT_DIR/test_constitution_inheritance.sh"
