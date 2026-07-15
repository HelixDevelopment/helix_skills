# G41 — PreToolUse Guard Hook Wiring Design (§11.4.109 Anti-Forgetting Enforcement)

**Revision:** 2
**Last modified:** 2026-07-15T21:56:03Z

## As-landed addendum (read this FIRST — supersedes the present-tense sections below)

**Verified against HEAD, 2026-07-15T21:33:52Z (§11.4.6 — read, not guessed).**
Everything from "## Goal + §11.4.109" through "## Hermetic test suite design"
below was written against the tree **before** commit `0438d7e` landed. At
HEAD, the repo state is materially different in four ways:

1. **A `.claude/settings.json` now exists — but not where this design
   proposed.** The "Current settings (cited)" section below examined this
   project's repo-root `.claude/` directory, which held (and still holds,
   unchanged, 286 bytes) only `settings.local.json` with no `hooks` key —
   confirmed still true at HEAD for that specific file/path. What actually
   landed is a **separate, new** file at
   `docs/research/mvp/Agent_AI_Skill_Tree_Development/.claude/settings.json`
   (the MVP-project subdirectory, verified present and read in full at HEAD)
   — not an edit to the repo-root `settings.local.json` this design's
   "`.claude/settings.local.json` hook wiring" section (below) proposed
   wiring into. This supersedes the "There is no `.claude/settings.json` in
   this repo" claim (below, originally line 44) and the "No `hooks` key at
   all … zero PreToolUse enforcement exists" claim (below, originally lines
   64-65) — both were accurate only of the pre-`0438d7e` tree; at HEAD a
   `hooks.PreToolUse` key is wired (at the MVP-subdir path, not the repo
   root).
2. **Path expression differs from the design.** The design's JSON (in the
   wiring-design JSON below) specifies `bash "$CLAUDE_PROJECT_DIR/constitution/scripts/hooks/
   guard-forbidden-commands.sh"`. The landed hook command
   (`docs/research/mvp/Agent_AI_Skill_Tree_Development/.claude/settings.json`,
   line 9, verified) is instead `bash "$(git -C "$CLAUDE_PROJECT_DIR"
   rev-parse --show-toplevel)/constitution/scripts/hooks/
   guard-forbidden-commands.sh"`. The landed hermetic suite's own header
   comment explains why
   (`project/scripts/test_guard_forbidden_commands.sh:65-71`, verified): the
   wired settings.json lives four directories into the repo, so
   `$CLAUDE_PROJECT_DIR` at hook-invocation time may resolve to that MVP
   subdirectory rather than the true repo root, and a bare
   `"$CLAUDE_PROJECT_DIR/constitution/..."` path (as this design specifies)
   would then resolve incorrectly; `git rev-parse --show-toplevel` resolves
   the actual repo root regardless of which directory is treated as the
   project root.
3. **Hook set is narrower than designed.** The design's Phase-1 `Bash`-matcher
   group (in the Bash-matcher hook-group JSON below, and in the Phase-1
   rollout passage below) wires **two** hooks —
   `guard-forbidden-commands.sh` AND `guard-work-track-binding.sh`. The
   landed settings.json wires **only** `guard-forbidden-commands.sh`
   (verified — the file's sole `PreToolUse` entry has a single hook in its
   `hooks` array). `guard-work-track-binding.sh` (§11.4.191) is NOT wired at
   HEAD; the Phase-1/Phase-2 rollout choice this design surfaced (below) was
   resolved narrower than either phase as designed — only one of the two
   Phase-1 hooks landed.
4. **Harness location + scope differ from the design's plan.** The design
   (in the hermetic-test-harness section below) sites the hermetic suite at
   `constitution/scripts/hooks/test_guard_forbidden_commands.sh` and states
   it "is written to include D4/D5 as RED today" (in the D4/D5
   false-positive-guard findings below). At
   HEAD the harness instead lives at
   `project/scripts/test_guard_forbidden_commands.sh` (185 lines, verified
   read in full) and contains **no case** exercising the D4 ("mention-only,
   unquoted, no trailing quote") or D5 ("mention-only, shell comment")
   scenarios from the "Hermetic test suite design" §D table below — it has a
   D1-equivalent case (`"tail -f push-failure log (false-positive
   regression, §11.4.180)"`, line 130, exit 0 expected) but nothing probing
   the D4/D5 unquoted-mention / bare-comment shapes. **The upstream §11.4.201
   false-positive gap this design documents (D4/D5 below) therefore remains
   OPEN and UNTESTED by the landed suite** — flagged by this design, not
   fixed upstream, and not (yet) held RED per §11.4.115 polarity in the
   landed harness as this design recommended.

Nothing in §3-§8 of this document (the inherited-guard contracts, the
Phase-1/Phase-2 rollout rationale, the `docs/AGENT_GUARDRAILS.md` plan, the
gate mapping, and the honest boundary) is invalidated by the above — those
sections describe mechanisms and reasoning that remain accurate; only the
"is this wired, where, and exactly how" present-tense claims in the "Current
settings" and ".claude/settings.local.json hook wiring" sections needed this
correction.

## Goal + §11.4.109

R23's compliance audit filed **G41** as a HIGH violation: this project has no
`PreToolUse` guard hook wired, so the constitutional command-classes that
MUST be mechanically blocked (§11.4.109) — host-direct emulator/adb, force-push
and verification bypass (§11.4.113), sudo/su, host-power — are enforceable
today only by an agent *remembering* to refuse them. §11.4.109 exists exactly
because "a rule the orchestrator forgets to paste is not enforcement."

This document designs — but does **not** apply — the wiring of the three
guard hooks the constitution submodule ships at
`constitution/scripts/hooks/`, plus the project-side `docs/AGENT_GUARDRAILS.md`
companion and a hermetic test suite. §11.4.109 clause (1) requires the hook
wired "in `.claude/settings.json` (or equivalent runtime settings)"; clause
(2) requires the guardrails doc with headings `SUBAGENT CONSTITUTIONAL
PREAMBLE` + `ORCHESTRATOR PRE-ACTION CHECKLIST` and the anchor literal
`11.4.109`; clause (3) requires a hermetic hook test suite with ≥20 cases.

Composition this design must satisfy (per task brief): §11.4.109 (the hook
mechanism itself) · §11.4.113 (force-push is one of the blocked classes) ·
§11.4.182 (track+branch+alias label on Agent dispatches) · §11.4.191
(work-to-track binding) · §11.4.201 (a guard MUST assert the REAL condition —
golden-TRUE + golden-FALSE-with-carrier fixtures, no false-positive
refusals).

**All findings below were verified by actually executing the three shipped
hooks against synthetic PreToolUse JSON payloads on this host** (§11.4.6 —
no guessing about exit codes; every exit code cited was captured, not
assumed).

## Current settings (cited) — pre-0438d7e snapshot; superseded by the As-landed addendum above

`ls .claude/` (repo root) showed exactly one file, as of this writing:

```
.claude/settings.local.json   (286 bytes, tracked)
```

There was **no** `.claude/settings.json` at the repo root at the time this
was written (this claim is superseded — see the As-landed addendum above; a
`.claude/settings.json` now exists, but at the MVP-subdirectory path, not the
repo root). This section's repo-root `settings.local.json` content, as of
this writing, was:

```json
{
  "permissions": {
    "allow": [
      "WebSearch",
      "WebFetch(domain:arxiv.org)",
      "Read(//tmp/**)",
      "Bash(echo \"exit=$?\")"
    ],
    "defaultMode": "bypassPermissions"
  },
  "enableAllProjectMcpServers": true,
  "enabledMcpjsonServers": [
    "codegraph"
  ]
}
```

No `hooks` key at all in this repo-root file — confirming G41's HIGH finding
as of this writing: zero `PreToolUse` enforcement existed in this project at
that time (superseded — see the As-landed addendum above; a `hooks.PreToolUse`
key is now wired, at the MVP-subdirectory settings.json, not this repo-root
file).

**Naming note, verified, not a defect to fix here.** `.claude/settings.local.json`
is normally Claude Code's *personal, untracked* override file. This project's
`.gitignore` line 1 is `!.claude/settings.local.json` — an explicit
un-ignore that makes this project **track** `settings.local.json` as its
shared/project config (`git ls-files .claude` returns only this one path).
There was no separate tracked `settings.json` at the repo root, as of this
writing (superseded — see the As-landed addendum above; a separate, tracked
`settings.json` now exists, but at the MVP-subdirectory path, not the repo
root). §11.4.109 clause (1)'s "(or
equivalent runtime settings)" phrasing accommodates this: Claude Code reads
whichever settings files are present regardless of git-tracking status, so
wiring the hooks into the existing tracked `settings.local.json` satisfies
the mandate without introducing a second settings file. (An alternative —
introducing a genuine tracked `settings.json` and reserving `settings.local.json`
for untracked personal overrides — is a separate, larger repo-hygiene change
and is **not** designed here; it is out of scope for closing G41.)

## Inherited guards (cited, by-reference)

All three hooks live at `constitution/scripts/hooks/` and are consumed
**by reference** — never copied into this project's own `scripts/` tree
(§11.4.109/§11.4.177: a copy diverges silently from upstream fixes).

### 1. `guard-forbidden-commands.sh` (§11.4.109 primary enforcement + §11.4.113)

Contract: reads the PreToolUse JSON payload from stdin; for `tool_name !=
"Bash"` or an empty `.tool_input.command`, exits 0 (allow) immediately.
Otherwise scans the command for four blocked classes:

| Class | Constitutional clause | Detector |
|---|---|---|
| Host-direct emulator / `adb install` / `am instrument` | §6.X / §6.V / §6.AG (project-layer, Containers-submodule mandate §11.4.76) | word-bounded regex on `emulator -avd`, `adb [-s <serial>] install`, `am instrument` |
| `git push --force` / `-f` / `--force-with-lease` / `+<refspec>`, `--no-verify`, `--no-gpg-sign` | §11.4.113 (absolute no-force-push) / §6.T.3 | quote-aware clause-splitter (`;`/`&&`/`\|\|`/`\|`/newline, never splitting inside quotes) + a word-bounded `git...push` detector, requiring the git-push match AND the force-flag/`+`-prefix to land in the **same clause** — the 2026-07-11 false-positive fix (see below) |
| `sudo` / `su` | §6.U | word-bounded, standalone-token match (not `sudo`-as-substring like `subl`) |
| Host-power (`systemctl suspend/hibernate/poweroff/...`, `loginctl ...`, `pm-suspend`, `shutdown`) | Host Machine Stability Directive (project §12) | word-bounded regex |

Exit-code contract: **exit 0 = allow**, **exit 2 = BLOCK** (stderr text fed
back to the agent as the refusal reason), any other exit code never used.

Escape hatch: a command containing the literal marker
`# guardrails:allow <reason>` downgrades a match to a stderr **WARNING**
(exit 0, not blocked) for every class **except** host-power, which the
script marks `no-override` — the marker is ignored for host-power and the
block still fires at exit 2.

**Verified behaviour (executed, not guessed) — see the "Hermetic test suite"
section below for the full captured-evidence table.**

### 2. `guard-track-branch-label.sh` (§11.4.182)

Contract: only inspects `tool_name ∈ {Agent, Task, TaskCreate}`; every other
tool exits 0 untouched. For a matched dispatch, `.tool_input.description`
(falling back to `.tool_input.subagent`) MUST start with
`^\(T[0-9]+/[^)]+ - [^)]+\) ` — i.e. `(T<N>/<branch> - <alias>)` where `<N>`
is **mandatorily numeric**. It also cross-checks the dispatch's `<alias>`
field against the LIVE alias derived from `CLAUDE_CONFIG_DIR` (basename
`.claude-<alias>`) via the reference labeler
`constitution/scripts/multitrack/track_branch_label.sh` — a concrete-vs-concrete
mismatch blocks; a `?` on either side is honestly accepted (§11.4.6). No
`# guardrails:allow` escape hatch exists anywhere in this script (verified
by reading its full source — there is no marker-handling logic at all,
unlike the other two hooks).

**Critical, verified finding (see "Blocker" below):** the track-number
component `<N>` is **required to be numeric by the regex itself** — the
reference labeler's own fallback for an off-track cwd (any path not under
`/mnt/track<N>/...`) emits literal `?` for `<N>`, and `?` does **not** match
`[0-9]+`. So a dispatch issued from a checkout that is not itself an
`/mnt/track<N>/...` mount is blocked **unconditionally**, even when its
label is honestly and correctly formed. This is by design per the source
comment: *"Track MUST be numeric (an off-track '?' is surfaced by a BLOCK,
never mislabeled)."*

### 3. `guard-work-track-binding.sh` (§11.4.191)

Contract: handles **both** the dispatch path (`tool_name ∈ {Agent, Task,
TaskCreate}`) and the commit path (`tool_name == "Bash"` matching a
`git commit` / `git add -A|.|--all` / `commit_all.sh` / `commit_docs.sh`
pattern), everything else exits 0.

- **Dispatch path**: parses the same `(T<N>/<branch>[ - <alias>])` prefix;
  if the label does not match the numeric-track shape at all, it exits 0
  (defers to the sibling `guard-track-branch-label.sh`, avoiding a double
  report). If it matches AND the description contains one or more
  `ATM-NNN`/`SPK-NNN` ticket references, it resolves the ticket's
  logic-group's canonical `(branch, track)` via
  `constitution/scripts/multitrack/multitrack_work_binding.sh check ...`
  against the workable-items DB and blocks (exit 2) on mismatch.
- **Commit path**: resolves the to-be-committed file set from git plumbing
  (staged set, else the whole dirty worktree for a broad-stage command) and
  checks each file's owning logic-group's canonical `(branch, track)`
  against the current checkout via the same resolver. Fails closed (exit 2)
  on a `$(...)`/backtick dynamic pathspec it cannot verify. Verified on this
  project: an ordinary `git commit -m ...` with nothing staged and no `-a`
  flag exits 0 (FILESET stays empty → "nothing to commit → nothing to
  check"), regardless of the 21 currently-dirty files in this working tree.

Escape hatch: `# guardrails:allow <reason>` downgrades a commit-path block
to an audited WARN (mirrors `guard-forbidden-commands.sh`); the dispatch
path has no escape-hatch handling of its own (it defers entirely to the
resolver's verdict).

## `.claude/settings.local.json` hook wiring (design — NOT applied)

All three scripts self-filter on `tool_name` internally, so each is wired
against the narrowest matcher it actually acts on: `guard-forbidden-commands.sh`
→ `Bash` only; `guard-track-branch-label.sh` → `Agent|Task|TaskCreate` only;
`guard-work-track-binding.sh` → both (registered once per matcher group,
since a single tool call only ever matches one matcher, so it never double-fires).
Paths are referenced via `$CLAUDE_PROJECT_DIR` (the Claude Code hook
environment variable that resolves to the project root regardless of the
invoking shell's cwd) so the wiring is portable to any clone location —
never a hardcoded absolute path (§11.4.6/§11.4.177).

```json
{
  "permissions": {
    "allow": [
      "WebSearch",
      "WebFetch(domain:arxiv.org)",
      "Read(//tmp/**)",
      "Bash(echo \"exit=$?\")"
    ],
    "defaultMode": "bypassPermissions"
  },
  "enableAllProjectMcpServers": true,
  "enabledMcpjsonServers": [
    "codegraph"
  ],
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "bash \"$CLAUDE_PROJECT_DIR/constitution/scripts/hooks/guard-forbidden-commands.sh\"",
            "timeout": 10
          },
          {
            "type": "command",
            "command": "bash \"$CLAUDE_PROJECT_DIR/constitution/scripts/hooks/guard-work-track-binding.sh\"",
            "timeout": 10
          }
        ]
      },
      {
        "matcher": "Agent|Task|TaskCreate",
        "hooks": [
          {
            "type": "command",
            "command": "bash \"$CLAUDE_PROJECT_DIR/constitution/scripts/hooks/guard-track-branch-label.sh\"",
            "timeout": 10
          },
          {
            "type": "command",
            "command": "bash \"$CLAUDE_PROJECT_DIR/constitution/scripts/hooks/guard-work-track-binding.sh\"",
            "timeout": 10
          }
        ]
      }
    ]
  }
}
```

### BLOCKER — do not wire the `Agent|Task|TaskCreate` group unconditionally yet

Verified by direct execution (payloads + captured exit codes below): this
project's own conductor session runs from
`/home/milos/Factory/projects/tools_and_research/helix_skills`, which is
**not** an `/mnt/track<N>/...` mount (the host does have `/mnt/track1..4`
for the multi-track engine, but this checkout is the conductor's main
checkout, per §11.4.187 clause 4 — "the conductor... is NEVER bound into a
track worktree"). Because `guard-track-branch-label.sh`'s track field is
**mandatorily numeric** and has **no escape hatch**, wiring it now would
unconditionally BLOCK (exit 2) every `Agent`/`Task`/`TaskCreate` dispatch
issued from this checkout — including an *honestly, correctly* labeled one
— because the reference labeler can only ever emit `?` for `<N>` here, and
`?` never satisfies `[0-9]+`. The dispatch that produced *this very
document* did not carry a `(T<N>/branch - alias)` prefix, illustrating the
project's current, un-labeled Agent-dispatch convention.

**Recommended phased rollout** (a decision for the operator, per §11.4.66 /
§11.4.101 — this is not a reversible, zero-risk default, so it is not
resolved unilaterally here):

- **Phase 1 (safe, closes the HIGH G41 finding on its own):** wire only the
  `Bash` matcher group (`guard-forbidden-commands.sh` +
  `guard-work-track-binding.sh`'s commit path). This alone fully satisfies
  §11.4.109's core ask — the forbidden-command classes are mechanically
  blocked — and was verified to have zero impact on ordinary `git`/shell
  use in this project (see test cases below).
  - Note: `guard-work-track-binding.sh` in the `Bash` group is asserted safe
    here because this project currently has no `docs/workable_items.db` and
    no staged commits pass the broad-stage/staged-set precondition (verified
    above) — it is a no-op today, and only becomes load-bearing once/if this
    project adopts the §11.4.93 workable-items DB.
- **Phase 2 (operator-gated):** wire the `Agent|Task|TaskCreate` group only
  after the operator decides EITHER (a) this project adopts the
  `/mnt/track<N>/helix_skills/` multi-track layout (§11.4.178's
  multi-project-per-track clarification) so its own dispatches run
  track-numbered, OR (b) the project's dispatch convention is updated to
  always prefix Agent/Task/TaskCreate descriptions with a valid
  `(T<N>/<branch> - <alias>)` label before this hook is wired. Wiring Phase
  2 before that decision would break subagent dispatch outright with no
  override.

## Project `docs/AGENT_GUARDRAILS.md` plan

Per §11.4.35 (project-vs-universal split): the constitution's
`constitution/docs/AGENT_GUARDRAILS.md` classifies itself explicitly as
**universal** ("the preamble + checklist pattern is reusable across any
HelixConstitution-consuming project; the §-clause numbers are Lava-specific
examples"). Its ten numbered preamble rules cite a *different* consuming
project's clause numbers (e.g. "§6.X / §6.V / §6.AG" for emulators, "§6.Z"
for pre-distribute gating, "§6.R" for no-hardcoding) — these are that other
project's project-layer instantiations of universal anchors, not this
project's.

**Design decision: REFERENCE, do not adapt-by-copy.** This project
(`helix_skills`) is a skills/design-research repository with no emulator
fleet, no APK distribution pipeline, and (currently) no multi-track adoption
of its own — most of the Lava-specific project-layer clause citations in the
upstream preamble (rules 1, 7) do not have a project-local equivalent here
yet. Rewriting the preamble with `helix_skills`-specific clause numbers that
don't exist would be inventing structure (§11.4.6). The correct instantiation
is a **thin project-local `docs/AGENT_GUARDRAILS.md`** that:

1. Opens with the same `## INHERITED FROM constitution/CLAUDE.md`-style
   pointer this project's root `CLAUDE.md` already uses (§11.4.35).
2. States plainly: *"This project has no project-layer clause set of its
   own yet (no §6.X-style numbering exists in `helix_skills/CLAUDE.md`).
   The SUBAGENT CONSTITUTIONAL PREAMBLE below is the constitution's
   universal preamble with its project-layer-only rules (1, 7, 8) marked
   `N/A — this project` rather than silently dropped or force-fit onto
   invented clause numbers."*
3. Reproduces the `SUBAGENT CONSTITUTIONAL PREAMBLE` heading verbatim,
   preamble rules 2–6, 9–10 unchanged (universal, §11.4-numbered, apply
   as-is), rules 1/7/8 annotated `N/A — this project` with a one-line
   reason (no emulator fleet / no distribution pipeline / no approved-remote
   allowlist defined yet) rather than invented placeholders.
4. Reproduces the `ORCHESTRATOR PRE-ACTION CHECKLIST` heading verbatim with
   the same N/A annotation discipline on its Emulator and Distribute
   subsections.
5. Carries the anchor literal `11.4.109` at least once (satisfies
   §11.4.109 clause 2 / gate `CM-COVENANT-114-109-PROPAGATION`'s literal
   check) plus `11.4.182` and `11.4.191` since this doc also documents
   those two hooks' wiring per the task's composition requirement.
6. Points at the three hooks **by reference**
   (`constitution/scripts/hooks/guard-*.sh`) — never copies their source
   into `docs/` or `scripts/`.
7. States the Phase 1/Phase 2 rollout decision from this design (or its
   resolution, once the operator decides) so a future session reads the
   live wiring state from one place.

## Hermetic test suite design (≥20 cases, real evidence cited)

All cases below were executed against the live scripts on this host with
synthetic PreToolUse JSON on stdin (`{"tool_name": "...", "tool_input":
{...}}`), not merely designed on paper. Each row cites its **actual
captured exit code** (§11.4.6 — no guessing). A hermetic test harness
(`constitution/scripts/hooks/test_guard_forbidden_commands.sh`, a sibling of
the existing `test_guard_branch_consistency.sh` / `test_guard_track_branch_label.sh`
/ `test_guard_work_track_binding.sh`, which this design confirms are absent
for `guard-forbidden-commands.sh` — a gap this design flags for a follow-on
workable item) should assert exactly these cases mechanically.

### A. `guard-forbidden-commands.sh` (14 cases)

| # | Case | Command | Expected | **Captured** |
|---|---|---|---|---|
| A1 | Blocked: emulator | `emulator -avd test_avd` | exit 2 | **exit 2** ✓ |
| A2 | Blocked: adb install | `adb install app.apk` | exit 2 | **exit 2** ✓ |
| A3 | Blocked: adb -s install | `adb -s emulator-5554 install app.apk` | exit 2 | **exit 2** ✓ |
| A4 | Blocked: am instrument | `am instrument -w com.foo/androidx.test.runner.AndroidJUnitRunner` | exit 2 | **exit 2** ✓ |
| A5 | Blocked: force-push `--force` | `git push --force origin main` | exit 2 | **exit 2** ✓ |
| A6 | Blocked: force-push `-f` | `git push -f origin main` | exit 2 | **exit 2** ✓ |
| A7 | Blocked: force-push `--force-with-lease` | `git push --force-with-lease origin main` | exit 2 | **exit 2** ✓ |
| A8 | Blocked: `--no-verify` | `git commit --no-verify -m x` | exit 2 | **exit 2** ✓ |
| A9 | Blocked: `--no-gpg-sign` | `git commit --no-gpg-sign -m x` | exit 2 | **exit 2** ✓ |
| A10 | Blocked: sudo | `sudo apt-get install foo` | exit 2 | **exit 2** ✓ |
| A11 | Blocked: su | `su -c whoami` | exit 2 | **exit 2** ✓ |
| A12 | Blocked, NO-OVERRIDE: `systemctl suspend` | `systemctl suspend` | exit 2 | **exit 2** ✓ |
| A13 | Blocked, NO-OVERRIDE: `shutdown` | `shutdown -h now` | exit 2 | **exit 2** ✓ |
| A14 | Non-Bash tool passthrough | `tool_name="Read"` | exit 0 | **exit 0** ✓ |

### B. Allowed / benign (§11.4.201 golden-FALSE) (3 cases)

| # | Case | Command | Expected | **Captured** |
|---|---|---|---|---|
| B1 | Ordinary status | `git status` | exit 0 | **exit 0** ✓ |
| B2 | Ordinary push, no force | `git push origin main` | exit 0 | **exit 0** ✓ |
| B3 | Unrelated command | `ls -la` | exit 0 | **exit 0** ✓ |

### C. Escape hatch (2 cases)

| # | Case | Command | Expected | **Captured** |
|---|---|---|---|---|
| C1 | Escape hatch fires (non-power class) | `git push --force origin main # guardrails:allow operator-approved-mirror-reconciliation` | exit 0 + stderr WARNING | **exit 0**, stderr = `guardrails: WARNING — §6.T.3 force-push: ...` ✓ |
| C2 | Escape hatch REJECTED for host-power (no-override) | `systemctl suspend # guardrails:allow operator-said-ok` | exit 2 (marker ignored) | **exit 2** ✓ |

### D. §11.4.201 false-positive-guard cases — carrier fixtures (5 cases, the load-bearing set)

| # | Case | Command | Expected (no false positive) | **Captured** |
|---|---|---|---|---|
| D1 | Carrier: `-f` on an unrelated command chained after a real `git fetch` | `git fetch --all && tail -f qa-results/push_failures/x.log` | exit 0 | **exit 0** ✓ (this is the exact 2026-07-11 forensic regression cited in the script's own header comment — re-verified live here) |
| D2 | Carrier: path literally named `push_failures` | `ls qa-results/push_failures/` | exit 0 | **exit 0** ✓ |
| D3 | Mention-only, quoted, trailing-quote artifact | `echo "we must never run: git push --force"` | exit 0 (intended) | **exit 0** — see finding below |
| D4 | Mention-only, unquoted, no trailing quote | `echo NOTE never run git push --force` | exit 0 (intended) | **exit 2** ✗ — **VERIFIED GAP** |
| D5 | Mention-only, shell comment | `# reminder: do not git push --force here` | exit 0 (intended) | **exit 2** ✗ — **VERIFIED GAP** |

**§11.4.201 finding (verified by execution, not guessed) — a real,
currently-shipping false-positive class in `guard-forbidden-commands.sh`.**
D1/D2 (the documented 2026-07-11 fix) genuinely work: an unrelated `-f`
flag or a path merely containing the string `push_failures` does not
trigger a block, because the quote-aware clause splitter correctly keeps
the git-push detector and the force-flag detector scoped to the same
executed clause. D3 passes, but **not** for a principled content-vs-command
reason — direct regex testing (`[[ "$CMD" =~ ... ]]` executed on this host)
shows the git-push detector (`GIT_PUSH_RE`) **does** match inside the quoted
echo string (`' git push '` captured), and the block is only avoided because
the trailing `"` immediately after `--force` defeats the force-flag
regex's own end-anchor (`([[:space:]]|=|$)` — the character after `--force`
is a literal quote, not space/`=`/true-end-of-string). D4 and D5 — the same
"merely discussing/mentioning the forbidden command, never invoking it"
intent, phrased without that accidental trailing-quote artifact — are
**incorrectly BLOCKED** (exit 2), a textbook §11.4.201 false-positive
refusal: the guard fires on the proxy signal "these words appear in
sequence somewhere in the command string" rather than the REAL condition
"a `git` subcommand `push` is actually about to execute with a force flag."
**This is a genuine, currently-present defect in the upstream shared hook**,
newly discovered by this design's hermetic-suite exercise — its fix is a
constitution-submodule change (§11.4.26 workflow: fetch/pull → classify →
validate → commit+push to all upstreams) and is explicitly **out of scope**
for this read-only G41 wiring design. Recorded here as a tracked follow-on
finding; the hermetic suite designed above is written to include D4/D5 as
RED today (§11.4.115 polarity) so the fix, once landed upstream, is proven
by the same suite flipping to GREEN.

### E. `guard-track-branch-label.sh` (§11.4.182) (4 cases)

| # | Case | Payload | Expected | **Captured** |
|---|---|---|---|---|
| E1 | Blocked: no label at all | `Agent`, description = `"no label here just a task description"` | exit 2 | **exit 2** ✓ |
| E2 | Allowed: well-formed label | `Agent`, description = `"(T1/main - claude1) ATM-312 do the thing"` | exit 0 | **exit 0** ✓ |
| E3 | Non-agent tool passthrough | `Bash`, same description text | exit 0 | **exit 0** ✓ |
| E4 | BLOCKER case: honest off-track label (`T?`) still blocked | `Agent`, description = `"(T?/main - ?) ATM-312 honest unknowns"` | **exit 2** (documented: track MUST be numeric) | **exit 2** ✓ — confirms the BLOCKER above |

### F. `guard-work-track-binding.sh` (§11.4.191) (3 cases)

| # | Case | Payload | Expected | **Captured** |
|---|---|---|---|---|
| F1 | Dispatch, no ticket referenced → honest passthrough | `Agent`, description = `"(T1/main - claude1) no ticket mentioned here"` | exit 0 | **exit 0** ✓ |
| F2 | Commit path, non-commit command → passthrough | `Bash`, command = `"ls -la"` | exit 0 | **exit 0** ✓ |
| F3 | Commit path, real `git commit`, nothing staged, this project's actual dirty tree (21 files, none staged) → passthrough | `Bash`, command = `"git commit -m 'test commit message'"` | exit 0 | **exit 0** ✓ |

**Total designed cases: 14 (A) + 3 (B) + 2 (C) + 5 (D) + 4 (E) + 3 (F) = 31**,
exceeding the §11.4.109 clause-3 floor of ≥20. 29 of 31 pass today
(exit code matches "expected/intended" behaviour); 2 (D4, D5) are honest
RED findings against the *intended* no-false-positive behaviour, tracked
per §11.4.115 polarity rather than silently omitted or weakened to match
the current (defective) output.

## Gate mapping + §1.1 mutations

| Gate | Anchor | Asserts | Paired §1.1 mutation |
|---|---|---|---|
| `CM-ANTI-FORGETTING-ENFORCEMENT` | §11.4.109 | hook present at `constitution/scripts/hooks/guard-forbidden-commands.sh` + wired as a `PreToolUse` entry in `.claude/settings.local.json` + `docs/AGENT_GUARDRAILS.md` present + hermetic test suite present | Remove the `Bash`-matcher hook entry from `settings.local.json` → gate FAILs ("hook not wired"); separately, move/rename `guard-forbidden-commands.sh` aside → gate FAILs ("hook not present") |
| `CM-COVENANT-114-109-PROPAGATION` | §11.4.109 | literal `11.4.109` present in this project's `CLAUDE.md`/`AGENTS.md`/`QWEN.md` mirrors AND in the new project `docs/AGENT_GUARDRAILS.md` | Strip the literal from `docs/AGENT_GUARDRAILS.md` → gate FAILs |
| `CM-TRACK-BRANCH-LABEL` | §11.4.182 | labeler (`track_branch_label.sh`) + hook (`guard-track-branch-label.sh`) exist/executable/`bash -n`-clean; hook wired in settings (Phase 2, once applied); convention doc exists | Remove the `Agent\|Task\|TaskCreate`-matcher `guard-track-branch-label.sh` entry → "hook wired in settings" sub-check FAILs |
| `CM-COVENANT-114-182-PROPAGATION` | §11.4.182 | literal `11.4.182` present in the project's governance docs | Strip the literal → gate FAILs |
| `CM-WORK-TRACK-BINDING-ENFORCED` | §11.4.191 | on `merge-base(HEAD,main)..HEAD`, every `group_paths`-owned file's registered group destination == current branch (detective layer — fires regardless of whether the preventive hook is wired) | Commit a file belonging to a registered group onto the wrong branch (once this project adopts §11.4.93 workable-items DB + logic groups) → gate FAILs |
| `CM-COVENANT-114-191-PROPAGATION` | §11.4.191 | literal `11.4.191` present in the project's governance docs | Strip the literal → gate FAILs |
| `CM-GUARD-ASSERTS-REAL-CONDITION` | §11.4.201 | every wired guard ships a golden-TRUE fixture (condition present → fires) AND a golden-FALSE-with-carrier fixture (condition absent, proxy signal present → must NOT fire) | Cases D4/D5 above ARE this mutation, already caught: `guard-forbidden-commands.sh` currently FAILS its own golden-FALSE-with-carrier fixture for the mention-only phrasing — a real, tracked finding, not a hypothetical mutation |

## Honest boundary (§11.4.6)

**What the hook (the floor) covers:** it makes the four §11.4.109 command
classes, the §11.4.182 label format, and the §11.4.191 file/dispatch
binding **mechanically impossible to skip by forgetting** — enforced at
the literal tool-call boundary, independent of what any agent (subagent or
conductor) recalls from a dispatch prompt. It does **not** understand
intent, does not know whether a force-push is operator-approved (only that
an in-band `# guardrails:allow <reason>` marker exists), and — per the D3–D5
finding above — its force-push detector currently has a real,
narrowly-scoped false-positive gap on mention-only/discussion text that
happens not to end in a specific quote-boundary artifact.

**What the preamble (the ceiling) covers, that the hook structurally
cannot:** anti-bluff intent (a test that passes without exercising the
real behaviour), resource-cap discipline, hardcoding avoidance,
pre-distribute gate execution, CONTINUATION-maintenance discipline, and
"real captured evidence, no guessing" — none of these are pattern-matchable
from a bash command string; they require the agent to actually apply
judgement per the pasted preamble. The hook and the preamble are
**complementary, not substitutable**: the hook guarantees the four
command-classes are blocked even when the preamble is never read; the
preamble covers everything the hook cannot see.

**What this design explicitly does NOT resolve:** whether/when Phase 2
(`Agent|Task|TaskCreate` wiring) should be applied to this project is an
operator decision this design surfaces but does not make (§11.4.66/§11.4.101
— high-blast-radius, not obviously reversible without an escape hatch, since
`guard-track-branch-label.sh` has none). The D4/D5 false-positive gap in the
upstream shared hook is flagged, evidenced, and left for a dedicated
constitution-submodule follow-on change — this design neither patches it
nor silently omits it from the test plan.
