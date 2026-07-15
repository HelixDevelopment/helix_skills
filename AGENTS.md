# Helix Skills — AGENTS.md

> Base agent rules live at `constitution/AGENTS.md` and the
> `constitution/Constitution.md` it references. **READ THOSE FIRST.**
> The base files are authoritative for any topic not covered here.
> Project-specific rules below extend them; they never weaken them.

## Critical base rules restated (for agents that don't follow @imports)

- **No bluffing.** Every PASS carries positive evidence. Constitution §11.4.
- **Mutation-paired gates.** Every new gate has a paired mutation
  proving it catches regressions. Constitution §1.1.
- **No guessing language.** `likely`, `probably`, `maybe`, `seems`,
  `appears` etc. are forbidden when reporting causes. Constitution §11.4.6.
- **Credentials never tracked.** Runtime-load only. Constitution §11.4.10.
- **Never force-push / never bypass hooks** without explicit,
  in-session operator authorization.
- **Hardlinked `.git` backup before any destructive op.** Constitution §9.
- **Multi-upstream push is the norm.** Constitution §2.1.

## Project-specific agent rules

_None beyond the inherited base at this time._

## Constitution inheritance

This project vendors the Helix Constitution at `constitution/`. Locate it
from any nested depth with `constitution/find_constitution.sh`. The
inheritance gate lives at `tests/constitution_inheritance_gate.sh` and its
false-positive-immunity meta-test at
`tests/meta_test_constitution_inheritance.sh`.
