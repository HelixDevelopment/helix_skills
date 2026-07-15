# CONTINUATION — HelixKnowledge Skill Graph System (MVP)

**Revision:** 1
**Last modified:** 2026-07-15T21:00:00Z
**Purpose:** §12.10 / §11.4.131 standing session-resumption file. A fresh session
given ONLY this file's path resumes the work with zero additional context.
Keep in sync on every material state change (new HEAD, phase, in-flight agent,
blocking decision).

---

## SHORT resume sentence (§11.4.127)

Read this file + `REQUIREMENTS.md` + `IMPLEMENTATION_PLAN.md` + `GAPS_AND_RISKS_REGISTER.md`,
`git fetch --all`, then continue the **P0.5 critical-remediation spine** (security
gaps from the R17 register) as a single serialized Go-mutator lane with a mandatory
**§11.4.209 Fable-xhigh review before every commit**, keeping ≤1 design-research
stream one step ahead in parallel.

## Where this work lives

- **Repo:** `helix_skills` (top-level git repo; only submodule = `constitution`).
  This project is the subtree `docs/research/mvp/Agent_AI_Skill_Tree_Development/`.
- **Go backend:** `project/` (module `github.com/helixdevelopment/skill-system`).
- **Upstreams (4, push via `git push origin HEAD:main` — fans out):** gitflic,
  github, gitlab, gitverse. **Absolute no force-push (§11.4.113); fast-forward only.**
- **Current HEAD:** `41f926a` (all 4 upstreams in sync as of this revision).

## Authoritative docs (read these; do NOT duplicate their content here)

| Doc | Role |
|---|---|
| `REQUIREMENTS.md` | Living source of truth for R1–R17 + founding blueprint + the TOON correction (real TOON wire format, NOT TOML). |
| `IMPLEMENTATION_PLAN.md` | Phases P0–P13 + **P0.5 critical remediation** + X1–X4 + traceability matrix. |
| `GAPS_AND_RISKS_REGISTER.md` | The R17 adversarial audit — **27 findings G01–G27**; this IS the P0.5 work queue. |
| `SPEC.md` | Technical spec (skill model, DAG, API, MCP). |
| `research/*.md` | Landed design research: granularity, toon_go_codec, opendesign_incorporation (R12), helix_interop_incorporation (R4), testing_infrastructure_plan (R8/R17), docs_chain_incorporation (R10), g02_sandbox_faildesign (R17, in flight). |

## PHASE / NEXT / terminal goal

- **PHASE:** P0.5 critical remediation (security + correctness gaps) before feature phases.
- **Terminal goal:** the full self-growing HelixKnowledge Skill Graph System per R1–R17,
  built BY the system, ~100% test-covered, zero bluff.
- **Immediate NEXT:** finish G01 attempt-2 → Fable-xhigh re-review → GO → commit → G02
  Go fix (against `research/g02_sandbox_faildesign.md`) → G03 → G05 → G06/G07 → G11.

## Live-state anchors (moment-valid — update on change)

- **In flight (background agents):**
  - **G01 attempt-2 fix** (Go, SOLE mutator): collapse the double-bound HTTP listeners
    to ONE hardened listener (mount `/mcp/v1` routes behind the same fail-closed
    CORS+auth via `HTTPTransport.RegisterRoutes`), fail-hard on bind error, fix the
    config `${VAR}` interpolation trap on `api_keys`/`allowed_origins`, add config
    tests + ops docs. Touches: `cmd/server/main.go`, `internal/mcp/http_transport.go`,
    `internal/mcp/server.go`, `internal/config/config.go` + `_test.go`, `config/config.toml`,
    `.env.example`, new `cmd/server/security_test.go`. **Not yet committed.**
  - **G02 sandbox fail-closed design** (research): vetted static-first / rootless-Podman-or-
    refuse design for the RCE sandbox. Writes `research/g02_sandbox_faildesign.md`.
- **Uncommitted working tree:** the G01-fix Go/config files above (awaiting re-review + GO).
- **G01 status:** attempt 1 was SOURCE-only; the §11.4.209 Fable-xhigh review returned
  **NO-GO** on a confirmed §11.4.108 SOURCE≠RUNTIME defect (a second wildcard-CORS +
  zero-auth MCP listener races the hardened one on the same port). Attempt-2 remediation
  is the in-flight fix. See the G01 STATUS block in the register for the full forensic trail.

## Binding constraints (do NOT violate)

- **Anti-bluff §11.4** — every claim carries captured evidence; I independently re-verify
  every subagent claim (re-run build/vet/test, grep residue, verify reachability) before commit.
- **§11.4.209** — the mandatory independent code review runs on **Fable at xhigh effort**
  (Opus xhigh only if Fable genuinely unavailable). Applies before EVERY commit/build.
- **§11.4.84** — never two Go-source mutators on this one tree at once; grep for mutation
  residue + account for every staged file before commit; narrow-stage docs (never `git add -A`
  while a mutator runs).
- **§11.4.113** — force-push strictly forbidden; merge-onto-latest-main, ff-only.
- **§11.4.108** — SOURCE green ≠ RUNTIME correct; prove the fix on the running binary.
- **TOON not TOML** — API wire = real TOON (`toon-format/toon-go`, MIT, vendor) + JSON
  fallback; TOML retained only for on-disk skill files + `config.toml`.
- **No `--no-verify`/`--force`/bypass flags. No secrets in git (§11.4.10).**

## OPEN operator decisions (surface before the dependent lane goes active)

- **G14 / X1 — submodule policy:** §11.4.28C single-canonical vs. the operator's stated
  parent-priority + both-synced. UNRESOLVED. Blocks vendoring OpenDesign / docs_chain /
  containers / the 7 helix-deps. Surface via §11.4.66 before starting the vendoring lane.

## Fleet discipline

Go remediation is one serialized mutator lane (§11.4.84). Parallelize only genuinely
independent work: keep ≤1 design-research stream one step ahead of the Go lane (avoids an
un-wired research backlog, §11.4.197). Every fix → Fable-xhigh review → commit → push 4
upstreams → next link.
