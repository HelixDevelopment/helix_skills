# Constitutional Compliance Audit — HelixKnowledge Skill Graph System

**Revision:** 1
**Last modified:** 2026-07-15T18:05:00Z
**Status:** compliance-audit, no code
**Scope:** `docs/research/mvp/Agent_AI_Skill_Tree_Development/` (Go backend under
`project/`, module `github.com/helixdevelopment/skill-system`).
**Authority audited against:** `constitution/Constitution.md` +
`constitution/CLAUDE.md` + `constitution/AGENTS.md` (canonical root, §11.4.35) +
gate scripts under `constitution/scripts/gates/`.
**Method:** §11.4.32-style enforcement sweep. Every verdict below is either
(a) **COMPLIANT** with cited file/line/gate evidence, (b) **VIOLATION** with the
exact gap + severity + fix direction, (c) **PENDING-AT-COMPLETION** — genuinely
un-satisfiable until the system is further built/deployed, with what will
satisfy it stated, or (d) **N-A** with a one-line reason. Forbidden vocabulary
(likely/probably/maybe/seems/appears/guess/apparently/perhaps) is not used
below except where quoting the Constitution's own text.

---

## 0. Executive summary

| Verdict | Count (of ~230 enumerated anchors incl. §1–§10, §12.x, §11.4.1–§11.4.209) |
|---|---|
| **COMPLIANT** | 47 |
| **VIOLATION** | 34 |
| **PENDING-AT-COMPLETION** | 26 |
| **N-A** | 108 |
| **UNCONFIRMED (honest audit-coverage gap)** | 15 |

The project is in **P0.5 critical-remediation** (pre-feature-complete,
pre-release) phase. Most PENDING/N-A verdicts are genuinely earned by that
phase (no release tag, no deployed target, no UI, no messenger surface, no
physical device) — not evasions. The VIOLATIONs cluster into five families,
all mechanically fixable without new engineering: (A) doc-governance
mechanics (revision headers, four-format exports, README doc-link section,
script companion docs), (B) the workable-item tracking constellation
(SQLite DB / ATM-NNN ids / Fixed.md-Summary pair — the project uses a
different, internally-consistent G0x markdown register instead), (C)
container-runtime default (rootful `docker` instead of rootless `podman`),
(D) the top-level governance-carrier set (missing QWEN.md/GEMINI.md; the
existing CLAUDE.md/AGENTS.md use the sanctioned `@import` pointer form, which
the *mechanical* propagation gates do not resolve — see §5), and (E)
project-specific request-history + PreToolUse guard-hook wiring.

---

## 1. Compliance ledger

### 1.1 Individually-assessed anchors (concrete, project-specific evidence)

| Anchor | Verdict | Evidence / gap |
|---|---|---|
| §1–§10 core covenant (no-bluff development principles) | COMPLIANT | `GAPS_AND_RISKS_REGISTER.md` is itself an anti-bluff artifact: every finding cites `file:line`, labels unverified claims `UNCONFIRMED`, and the register's own header states "positive-evidence-only, R11". |
| §11.4.1 FAIL-bluffs forbidden | COMPLIANT | G01 register entry documents attempt-1 producing a NO-GO from the Fable-xhigh reviewer rather than a false PASS; the reviewer caught a real SOURCE≠RUNTIME defect (two servers bound one port) that a shallower review would have missed. |
| §11.4.2 recorded-evidence requirement | PENDING-AT-COMPLETION | No user-visible feature is deployed/running yet (G03: validation + autoexpand pipelines are still dead code). Will be satisfied once `internal/api` is the live server and a captured request/response transcript exists under `docs/qa/`. |
| §11.4.3 per-environment-topology dispatch | N-A | Single-topology Go service (no device/OS-topology variance in the current MVP scope). |
| §11.4.4 test-interrupt-on-discovery | COMPLIANT | G01: "Attempt 1 hardened only the ad-hoc router... the §11.4.209 Fable-xhigh review returned NO-GO" — work stopped and was fixed before commit, not shipped past the discovery. |
| §11.4.5 captured-evidence quality analysis (audio/video) | N-A | No audio/video-producing feature exists in this project. |
| §11.4.6 no-guessing mandate | COMPLIANT | Forbidden-vocabulary scan of `IMPLEMENTATION_PLAN.md`, `REQUIREMENTS.md`, `SPEC.md`, `GAPS_AND_RISKS_REGISTER.md`, `CONTINUATION.md`, `plan.md` found zero live hedge-words; the sole hits (`IMPLEMENTATION_PLAN.md:22`) are the project's own restatement of the *rule itself* ("No `likely/should/probably`"), not a hedge. UNCONFIRMED items are explicitly labeled as such (e.g. G-register method note on `go build`/`go vet`). |
| §11.4.7 demotion-evidence rule | COMPLIANT | G01's STATUS block documents the NO-GO→fix→re-verify cycle citing the specific defect (double-bound listener) before any downgrade of severity. |
| §11.4.8 deep-web-research-before-implementation | COMPLIANT | `REQUIREMENTS.md` cites a live `git ls-remote` verification of `github.com/toon-format/toon` (HEAD `a19a117`) before adopting TOON as the wire format; 20+ `research/*.md` documents exist for granularity, embeddings, tree-sitter, sandboxing, etc. |
| §11.4.9 batch-source-fixes-before-rebuild | COMPLIANT | G02/G03/G05/G16/G21 were fixed and re-verified as one batch (CONTINUATION.md: "G02/G03/G05/G16/G21 Go impl DONE + Fable review-1 NO-GO (2 warnings) → FIXED → re-verified"). |
| §11.4.10 credentials-handling mandate | COMPLIANT | `project/.gitignore` ignores `.env`/`.env.local`; only `.env.example` is tracked; a secret-pattern scan (`sk-…`, AKIA…, PEM private-key headers) across tracked `.go`/`.toml`/`.yaml` returned zero hits. |
| §11.4.10.A pre-store credential leak audit | N-A | No new credential was stored by this project this session. |
| §11.4.11 file-layout discipline | COMPLIANT (minor caveat) | Clean `internal/`, `cmd/{cli,server,tui,worker}`, `migrations/`, `seed/`, `deploy/systemd/`, `docs/diagrams/` layout. Caveat: the MVP directory name itself is not snake_case — see §11.4.29. |
| §11.4.12 auto-generated docs sync | VIOLATION | Zero `.html`/`.pdf` siblings exist for any MVP-level doc (`REQUIREMENTS.md`, `IMPLEMENTATION_PLAN.md`, `SPEC.md`, `GAPS_AND_RISKS_REGISTER.md`, `CONTINUATION.md`, `plan.md`). No Issues_Summary equivalent regenerated. → **G38**. |
| §11.4.13 out-of-band sink-side evidence | N-A | No downstream network-introspection consumer is reachable yet (dead REST/MCP surfaces per G01/G03). |
| §11.4.14 test playback cleanup | N-A | No live-device test cycle has run yet. |
| §11.4.15 item-status tracking (6-state vocabulary) | VIOLATION | `GAPS_AND_RISKS_REGISTER.md` uses free-form `**STATUS (date) — <prose>**` blocks, not the closed vocabulary `{Queued\|In progress\|Ready for testing\|In testing\|Reopened\|Operator-blocked\|Fixed}`. Functionally informative but not mechanically sortable. → **G39**. |
| §11.4.16 item-type tracking ({Bug\|Feature\|Task}) | VIOLATION | Register uses a `**Category:**` field (`inconsistency/security`, `gap`, `existing-bug`, `test-coverage`, `danger-zone`) instead of the closed `{Bug\|Feature\|Task}` set. → folded into **G39**. |
| §11.4.17 universal-vs-project classification | N-A | This project proposed no new constitution-level rule this session. |
| §11.4.18 script documentation mandate | VIOLATION | `project/scripts/` has 13 scripts (`backup.sh`, `install.sh`, `_lib.sh`, `logs.sh`, `migrate.sh`, `package.sh`, `restart.sh`, `restore.sh`, `start.sh`, `status.sh`, `stop.sh`, `sync_submodules.sh`, `uninstall.sh`) with **zero** companion `docs/scripts/<name>.md` guides found anywhere under the MVP tree. → **G40**. |
| §11.4.19 Fixed-document column-alignment | N-A (constellation not adopted) | No `Fixed.md`/`Fixed_Summary.md` exist; the project uses the G0x register instead — see §11.4.93 discussion in §3. |
| §11.4.20 subagent-driven-by-default | COMPLIANT | CONTINUATION.md's "Fleet target 3–4 parallel" model plus the structurally-separated Fable-xhigh reviewer (never the same agent as the author) is exactly this discipline. |
| §11.4.21 operator-blocked status + self-resolution exhaustion | COMPLIANT | The **G14/X1 submodule-policy decision** is explicitly logged in CONTINUATION.md §"OPEN operator decisions" as blocked pending operator input, rather than auto-decided or silently dropped — the correct escalation path. |
| §11.4.22 document-sync commit discipline (lightweight docs-only wrapper) | VIOLATION | No `commit_docs.sh`-equivalent / `--docs-only` flag found; docs are committed via the normal commit path (per CONTINUATION's own §11.4.84 narrow-staging discipline, which is a partial substitute but not the dedicated wrapper). → folded into **G38**. |
| §11.4.24 build-resource stats tracking | N-A | `go build` for this module is a sub-1-minute build; the mandate applies to builds exceeding 1 minute wall-clock. |
| §11.4.25 full-automation-coverage mandate | PENDING-AT-COMPLETION | Feature coverage ledger cannot exist meaningfully while flagship pipelines are dead code (G03). Will be satisfied incrementally as G0x items close with the six invariants. |
| §11.4.26 constitution-submodule update workflow | COMPLIANT (N-A this session) | Constitution submodule present at expected path; not modified by this project this session, so the pull/validate/push pipeline was not triggered. |
| §11.4.27 no-fakes-beyond-unit + 100% test-type coverage | VIOLATION | 10 `_test.go` files exist (`security_test.go`, `skill_create_draft_test.go`, `skills_validation_test.go`, `middleware_test.go`, `auth_wiring_test.go`, `config_test.go`, `graph_test.go`, `sandbox_test.go`, `pipeline_test.go`, `analyzer_test.go`) — a real improvement over the register's original "0 tests" (G04) finding, but nowhere near the 13-test-type + Challenges + HelixQA bar (no integration/e2e/security/DDoS/scaling/chaos/stress/performance/benchmarking/UI/UX suite exists). → **G41**. |
| §11.4.28 submodule decoupling + dependency layout | COMPLIANT | `project/helix-deps.yaml` declares 7 deps (`llms_verifier`, `helix_llm`, `helix_agent`, `embeddings`, `helix_qa`, `challenges`, `docs_chain`) at `grouped` layout, schema-conformant with §11.4.31; no nested own-org submodule chain exists (none of the 7 are yet vendored — dry-run only per the manifest's own header). |
| §11.4.29 lowercase snake_case naming | VIOLATION | The project directory itself is `Agent_AI_Skill_Tree_Development` — TitleCase_With_Underscores, not lowercase snake_case. Not a language-mandated exception (this is a plain research-project directory, not a Java/Android resource path). → **G42**. |
| §11.4.30 .gitignore + no-versioned-build-artifacts | COMPLIANT | `project/.gitignore` covers binaries, test artifacts, `.env`, IDE files, `*.db`, `pgdata/`, config overrides, logs; `git ls-files` scan of `project/` found zero tracked binaries/`.db`/`bin/`/`dist/` entries. |
| §11.4.31 submodule-dependency-manifest | COMPLIANT | `helix-deps.yaml` present, `schema_version: 1`, correct `deps[].{name,ssh_url,ref,why,layout}` shape. |
| §11.4.32 post-constitution-pull validation | N-A | No constitution pull occurred in this project's scope this session. |
| §11.4.33 type-aware closure-status vocabulary | VIOLATION | Tied to §11.4.15/16 — the register never uses `Fixed (→ Fixed.md)` / `Implemented (→ Fixed.md)` / `Completed (→ Fixed.md)` since no Fixed.md exists. → folded into **G39**. |
| §11.4.34 reopened-source attribution mandate | PARTIAL → VIOLATION | G01's STATUS block is a rich narrative history but does not carry the literal `**Reopened-Details:** By/On/Reason/Evidence` line format. → folded into **G39**. |
| §11.4.35 canonical-root inheritance clarity | COMPLIANT | Top-level `CLAUDE.md` opens with the exact sanctioned `## INHERITED FROM constitution/CLAUDE.md` heading + `@constitution/CLAUDE.md` import; constitution submodule files are present at `constitution/{Constitution,CLAUDE,AGENTS}.md`. |
| §11.4.36 mandatory install_upstreams on clone/add | COMPLIANT | `git remote -v` shows 4 configured push remotes (gitflic, github, gitlab, gitverse) matching the "Upstreams (4)" note in CONTINUATION.md. |
| §11.4.37 fetch-before-edit mandate | UNCONFIRMED | This audit's own first git action was `git status`/`git log`, not a `git fetch --all --prune` — I did not verify remote-ahead state before reading. This is an honest gap in *this audit's own* procedure, not a finding about the project's prior sessions (which I cannot verify from static inspection). |
| §11.4.38 installable-asset evidence mandate | N-A | No user-distributable package/installer/container image has been built yet. |
| §11.4.40 full-suite retest before release tag | N-A | No release tag has been cut. |
| §11.4.41 pre-force-push merge-first mandate | N-A | No force-push event occurred (git config shows no force-related settings; §11.4.113 also holds — see below). |
| §11.4.42 iteration-discipline mandate | COMPLIANT | P0.5 explicitly orders CRITICAL (G01–G04) before HIGH/MEDIUM/LOW, and the current lane is the "single serialized Go-mutator... §11.4.209 review before every commit" — textbook priority-ordered batching. |
| §11.4.43 TDD-fix-discipline mandate | PARTIAL | G01/G02 register entries follow RED-defect→fix→VERIFY narrative; the LIVE-ADB-PROBE sub-step is N-A (no device). Whether each fix's test was authored *before* the fix (true RED) vs. after is UNCONFIRMED from static inspection of committed test files alone. |
| §11.4.44 document revision header mandate | VIOLATION | Only `CONTINUATION.md` carries the `**Revision:** N` / `**Last modified:**` header (Rev 5). `REQUIREMENTS.md`, `IMPLEMENTATION_PLAN.md`, `SPEC.md`, `GAPS_AND_RISKS_REGISTER.md`, `plan.md`, and every `research/*.md` file lack it. → **G43** (MUST-FIX-NOW — cheap, mechanical). |
| §11.4.45 integration-status-doc maintenance | N-A | No distinct multi-fix "integration domain" beyond the MVP itself yet; CONTINUATION.md + the register jointly serve this role informally. |
| §11.4.46 validate-recent-work-before-post-flash-tests | N-A | No device-flash concept in this project. |
| §11.4.47 Firebase data review mandate | N-A | No Firebase/Crashlytics integration. |
| §11.4.48/.49/.51/.117/.136/.137/.143/.158-.160/.163/.170/.193 (UI-driven/video/CV-OCR/streaming/vision mandates) | N-A (grouped) | No UI surface exists yet (`cmd/tui` is a stub directory per project structure); no video/audio output; no streaming app. Will re-open as PENDING once the TUI/Web client (R3) lands. |
| §11.4.50 deterministic consistency mandate | PENDING-AT-COMPLETION | No N-iteration (3×/10×) determinism proof exists yet for any test; will be satisfied as the stress/chaos suite (§11.4.85, G41-adjacent) lands. |
| §11.4.52 autonomous-validation mandate | VIOLATION | The 10 existing tests are unit/integration-shaped Go tests, not an end-to-end autonomous validation path with captured evidence per feature; the coverage ledger classification (`AUTONOMOUS_VERIFIED`/`DESIGNED`/`OPERATOR_ATTENDED_ONLY`/`N-A`) does not exist for this project. → folded into **G41**. |
| §11.4.53 Fixed_Summary parity | N-A | No Fixed.md/Fixed_Summary.md constellation adopted (see §11.4.93 discussion). |
| §11.4.54 ATM-NNN ticket identifier mandate | VIOLATION (operator-decision-needed) | See dedicated verdict in §3 below. → **G44**. |
| §11.4.55 reopens-history + per-item Reopens.md | VIOLATION | G01 has a rich in-place reopen-like history (attempt 1 NO-GO → attempt 2 GO) but no separate `docs/issues/<id>/Reopens.md`, and no `Reopens` column exists (no summary doc exists at all). → folded into **G39**. |
| §11.4.56 Status_Summary two-audience parity | N-A | No `Status.md`/`Status_Summary.md` domain doc pattern adopted for this project yet. |
| §11.4.57 README doc-link section + revision metadata | VIOLATION | `project/README.md` has no `Tracked-Items + Status Documents` section (confirmed: `grep -n "Tracked-Items"` → no match). → **G45**. |
| §11.4.58 parallel-development PWU pipeline | PARTIAL | CONTINUATION.md documents an informal "1 Go-mutator + 2–3 design-research streams" pattern consistent in spirit with the 5-stage PWU pipeline, but no formal PWU manifest / HelixQA Challenge-bank entries / lock-hierarchy tooling exists for this project. |
| §11.4.59/.60 README/documentation always-sync composite | VIOLATION | Direct consequence of §11.4.12/§11.4.44/§11.4.57 gaps above — no composite sync wrapper exists. → folded into **G38**. |
| §11.4.63 workable-items procedure docs SSoT | N-A | Tied to the §11.4.93 DB-adoption decision (§3). |
| §11.4.65/§11.4.73 four-format export + main-spec versioning | VIOLATION | `SPEC.md` and every MVP doc lack the §11.4.44 header (feeding into §11.4.73) and lack `.html`/`.pdf`/`.docx` exports (§11.4.65). → folded into **G38/G43**. |
| §11.4.66 blocker-resolution interactive-clarification | COMPLIANT | G14/X1 (submodule policy) and the R19 secondary-facet question are both surfaced as explicit open operator decisions rather than auto-resolved guesses. |
| §11.4.67 shell-script target-shell-parseability | UNCONFIRMED | The 13 scripts under `project/scripts/` were not run through `bash -n`/`sh -n` by this audit; no evidence of prior verification found either. Recommend as a companion check when G40 (script docs) is remediated. |
| §11.4.68 audio-specific sink-side evidence | N-A | No audio feature. |
| §11.4.69 universal sink-side positive-evidence taxonomy | PENDING-AT-COMPLETION | No feature yet maps to a taxonomy class with `ab_pass_with_evidence`-style proof (the REST/MCP surfaces are the closest candidate — `network_connectivity`/`network_throughput` — but are still dead-code-adjacent per G01 O3/G03). |
| §11.4.70 subagent-driven execution default | COMPLIANT | Same evidence as §11.4.20. |
| §11.4.71 pre-push fetch+investigate+integrate mandate | UNCONFIRMED | No visibility into actual past `git push` invocations from static file inspection. |
| §11.4.72 audio top-priority mandate | N-A | No audio feature; priority ordering in this project is security-first (G01) per R1, which is the correct domain-equivalent of "highest-risk-first" (§11.4.132). |
| §11.4.74 submodule-catalogue-first discovery | COMPLIANT | `helix-deps.yaml` reuses 7 existing Helix-family submodules (LLMsVerifier, HelixLLM, HelixAgent, Embeddings, HelixQA, Challenges, docs_chain) instead of reimplementing a jury/embeddings/QA engine from scratch — the "extend, don't reimplement" discipline is directly evidenced. |
| §11.4.75 mechanical enforcement without exception | PENDING-AT-COMPLETION | This project has not yet authored its own project-specific gates (only inherits constitution-level gates); as the doc-sync/DB/container fixes below land, project-level gates should accompany them. |
| §11.4.76 containers-submodule mandate | VIOLATION | `project/Makefile` and `project/docker-compose.yml` orchestrate containers directly (`CONTAINER_RUNTIME ?= docker`, raw `docker compose` invocations), not via the `vasic-digital/containers` submodule's `pkg/boot`/`pkg/compose`/`pkg/health` layer mandated by §11.4.76. → folded into **G46**. |
| §11.4.77 regeneration-mechanism-required | N-A | No gitignored generated artifact beyond ordinary `go build` output exists that needs a bespoke regen script. |
| §11.4.78/§11.4.79 CodeGraph mandate + own-org submodule inclusion | UNCONFIRMED | `.codegraph/codegraph.db` exists at the `helix_skills` repo root (evidence CodeGraph is wired at the parent-project level); whether the MVP subtree specifically is indexed, and whether it is kept current, was not independently verified (would require invoking the CodeGraph CLI, out of this audit's read-only scope). |
| §11.4.80 CodeGraph regular-update automation | UNCONFIRMED | Same basis as above. |
| §11.4.81 cross-platform-parity mandate | PENDING-AT-COMPLETION | `REQUIREMENTS.md` R3 explicitly marks Web/Desktop/Mobile clients as "greenfield" (not yet built) — correctly deferred, not silently dropped. |
| §11.4.82 iteration-speedup discipline | COMPLIANT | Batched fixes (§11.4.9 evidence) + Fable-xhigh review batching are the concrete instance of this discipline. |
| §11.4.83 docs/qa/ end-user evidence mandate | PENDING-AT-COMPLETION | No `docs/qa/<run-id>/` directory exists anywhere under the MVP tree (confirmed by directory search) because no feature has shipped/gone live yet (G03: flagship pipelines still dead code at the time of the last register update). **This MUST exist before any G0x item is marked closed/shipped** — flagged as a gating precondition, not yet a violation since nothing has been declared "done" without it. |
| §11.4.84 working-tree quiescence rule | COMPLIANT | CONTINUATION.md explicitly documents the "never two Go mutators on this tree" + "residue-scan + account for every staged file; narrow-stage" discipline verbatim, and states the mutation-in-flight state at time of writing. |
| §11.4.85 stress+chaos test mandate | VIOLATION | No stress or chaos test files found (none of the 10 `_test.go` files exercise sustained load, concurrent contention, process-death injection, network-fault injection, or resource-exhaustion). → folded into **G41**. |
| §11.4.86 roster/corpus-backed Status-doc auto-sync | N-A | No roster/corpus-backed Status doc pattern in use by this project. |
| §11.4.87 endless-loop autonomous work mandate | UNCONFIRMED | Cannot verify continuous-loop behavior across sessions from a static file snapshot; the git log shows a dense, continuously-advancing commit sequence consistent with it but this is circumstantial, not proof. |
| §11.4.88/§11.4.89 background-push / background-test-execution mandates | UNCONFIRMED | No visibility into whether past commits/tests were run in foreground vs. background from static inspection. |
| §11.4.90 Obsolete status + obsolescence audit | N-A | No item has been classified Obsolete; no register entry cites a superseding change or duplicate that would warrant it. |
| §11.4.91 summary-doc clarity mandate | N-A | No summary doc exists yet (tied to §11.4.15/16/93 cluster) — the anchor cannot be violated by a document that does not exist; it becomes live the moment a summary doc is authored. |
| §11.4.92 multi-pass change-evaluation discipline | COMPLIANT | G01's STATUS block demonstrates all 5 passes in substance: main-task verification (build/vet/test=0), blast-radius (auth scenario-space re-enumerated), cross-feature (both `/api/v1` and `/mcp/v1` groups checked), research (§11.4.8 citations elsewhere in the register), anti-bluff confirmation (6 §1.1 mutations M1–M6 run in a scratchpad copy). |
| §11.4.93 SQLite-backed workable-items SSoT | VIOLATION (with nuance — see §3) | → **G47**. |
| §11.4.94 zero-idle priority-first parallel-by-default | COMPLIANT (in spirit) | CONTINUATION.md's stream-sizing narrative explicitly surveys parallel-work feasibility and justifies the current fleet size rather than idling silently. |
| §11.4.95 workable-items DB tracked in git | N-A | Moot — no DB exists yet (§11.4.93). |
| §11.4.96 safe-parallel-work-with-long-build catalogue | N-A | This project's build is a fast Go build, not the AOSP-class multi-hour build this catalogue targets. |
| §11.4.97 maximum-use-of-idle-time | UNCONFIRMED | Same basis as §11.4.87. |
| §11.4.98 full-automation anti-bluff (no manual test intervention) | PENDING-AT-COMPLETION | No live/e2e/Challenge-class test exists yet to assess for manual-intervention dependence; the 10 existing unit-ish tests do not require manual steps (COMPLIANT at their own layer), but the mandate's full scope (live/e2e) is not yet reachable. |
| §11.4.99 latest-source documentation cross-reference | COMPLIANT | `REQUIREMENTS.md`'s TOON-format correction cites a live `git ls-remote` verification against `github.com/toon-format/toon` with an explicit HEAD SHA and date, rather than relying on training-data assumptions. |
| §11.4.100 (RETIRED) | N-A | Anchor repealed; no action required. |
| §11.4.101 autonomous-decision-over-blocking | COMPLIANT | G14/X1 is escalated (irreversible, high-blast-radius, cannot be determined from evidence alone) while everything else continues in parallel — exactly the decision rule's intended split. |
| §11.4.102 systematic-debugging auto-activation | COMPLIANT | G01/G02/G03 register entries are root-cause investigations (evidence → why-it-matters → decision) before any fix, matching the Iron Law. |
| §11.4.103 continuous parallel-stream routine | COMPLIANT (honest bottleneck disclosed) | CONTINUATION.md explicitly states the fleet is "honestly at 1 running" pending the Fable verdict and explains why spawning a filler stream would violate priority ordering (§11.4.42/§11.4.183) rather than help — this is the mandate's own honest-boundary clause satisfied, not violated. |
| §11.4.104/§11.4.105 participant-attribution / intent-recognition | N-A | No messenger/notification surface in this project. |
| §11.4.106 Docs Chain engine | VIOLATION | `docs_chain` is referenced only in research/incorporation notes (`research/docs_chain_incorporation.md`, `helix-deps.yaml`); it is not actually vendored or invoked to keep any MVP doc in sync. → folded into **G38**. |
| §11.4.107 AV anti-bluff techniques | N-A | No audio/video output. |
| §11.4.108 four-layer fix-verification + runtime-signature | PARTIAL | G01's STATUS block explicitly frames SOURCE (`go build/vet/test`=0), and RUNTIME-layer tests (one listener serves both route groups, unauth calls get 401/503) — strong evidence at those two layers. The ARTIFACT layer (bytes-in-the-built-binary) and RUNTIME-ON-CLEAN-TARGET (a freshly-deployed instance, not the dev tree) are not yet demonstrated — correctly PENDING given no deployment exists yet. |
| §11.4.109 anti-forgetting PreToolUse guard hook | VIOLATION | `.claude/settings.local.json` (top-level, applies to this whole repo including the MVP subtree) has **no hooks configured at all** — no `guard-forbidden-commands.sh` PreToolUse wiring, despite `constitution/docs/AGENT_GUARDRAILS.md` and the hook script both existing in the constitution submodule. → **G48** (MUST-FIX-NOW — cheap, mechanical, closes a real enforcement gap). |
| §11.4.110 pre-build readiness verdict + clash detection | PENDING-AT-COMPLETION | No project-specific READY-FOR-BUILD gate authored yet; reasonable to defer until the P0.5 Go-mutator lane stabilizes. |
| §11.4.111 resolve-by-stable-name | N-A | No enumerated hardware devices in this project. |
| §11.4.112 structural-impossibility classification | N-A | No goal has been proven structurally impossible in this project. |
| §11.4.113 absolute no-force-push mandate | COMPLIANT | `git config --list \| grep -i force` returns nothing; CONTINUATION.md explicitly states "Absolute no force-push (§11.4.113); fast-forward only"; git log shows a clean linear history with no rewritten commits. |
| §11.4.114 last-known-good-tag regression isolation | N-A | No release tag exists yet to serve as a regression oracle. |
| §11.4.115 RED-baseline-on-the-broken-artifact + polarity-switch | UNCONFIRMED | `sandbox_test.go`/`pipeline_test.go` exist and are named consistently with G02/G03 remediation, but whether they implement the literal `RED_MODE` polarity-switch pattern (vs. a plain post-fix regression test) was not verified by reading test source in this audit — a recommended spot-check, not confirmed either way. |
| §11.4.116 real-time conductor↔test-framework sync channel | N-A | No autonomous test-framework of that scale exists for this project yet. |
| §11.4.118 discovery-pressure to confirm known-issue-set completeness | COMPLIANT | `GAPS_AND_RISKS_REGISTER.md` is explicitly framed as "an adversarial audit" going beyond the originally-reported issues (its own header: "Method note... this audit did not re-run [build/vet]... findings are about design, behaviour, wiring, security, and contract fidelity") — exactly the discovery-pressure discipline. |
| §11.4.119 single-resource-owner partitioning | N-A | No shared exclusive hardware resource in this project. |
| §11.4.120 fix-breaks-its-own-gate reconciliation | N-A | No such event is documented. |
| §11.4.121 no-commit-while-build-writes-tracked-artifacts | N-A | No tracked directory receives build-time writes in this project. |
| §11.4.122 no-silent-removal-without-operator-confirmation | COMPLIANT | No end-user capability has been silently removed; G14/X1 (a decoupling/vendoring decision, not a removal) is correctly escalated rather than decided unilaterally. |
| §11.4.123 rock-solid-proof-or-deep-research mandate | COMPLIANT | 20+ `research/*.md` documents exist specifically because validation methods were not obvious (sandbox isolation feasibility, embedding-provider choice, tree-sitter integration) — the deep-research-first discipline is directly evidenced rather than a metadata-only closure. |
| §11.4.124 dead/unwired-code investigate-before-remove | COMPLIANT | G01's investigation explicitly confirms `internal/api` has zero importers via `grep -rln` *before* deciding to wire it in (not delete it) — the exact "investigate, then wire-in rather than delete" discipline the anchor mandates. |
| §11.4.125 code-review-agent gate before build | COMPLIANT | CONTINUATION.md's binding constraints state the Fable-xhigh review is mandatory "before EVERY commit/build". |
| §11.4.126 default autonomous-loop working mode | UNCONFIRMED | Cannot verify from a static snapshot whether the loop was engaged from the operator's first prompt of every session. |
| §11.4.127 session-handoff resumption-prompt mandate | COMPLIANT | CONTINUATION.md carries both a SHORT resume sentence and a FULL detail block, moment-valid (cites exact HEAD-adjacent state, PHASE/NEXT/terminal goal). |
| §11.4.128 always-on device-recording mandate | N-A | No physical test/debug device exists for this project. |
| §11.4.129 huge-blocker release protocol | N-A | No release-validation cycle has begun (pre-release phase). |
| §11.4.130 validate-fix-first-after-redeploy | N-A | No redeploy cycle exists yet. |
| §11.4.131 standing session-resumption file mandate | COMPLIANT | `CONTINUATION.md` exists at the declared path, carries the §11.4.44 header (Revision 5, Last modified 2026-07-15T17:10:00Z), and is current relative to the latest commits inspected. |
| §11.4.132 risk-ordered validation priority mandate | COMPLIANT | P0.5 explicitly orders the 4 CRITICAL findings (G01–G04) ahead of the 11 HIGH / 8 MEDIUM / 4 LOW findings, and within CRITICAL, the security hole (G01) was closed first. |
| §11.4.133 target+hardware safety mandate | N-A | No physical target hardware for a Go backend service. |
| §11.4.134 code-review iterate-until-GO mandate | COMPLIANT | G01: "Attempt 1... NO-GO... Attempt 2 (GO'd, committed)"; G02/G03/G05/G16/G21: "review-1 NO-GO (2 warnings) → FIXED → re-verified... GATED on the §11.4.209 Fable-xhigh RE-review" — explicit iterate-to-zero-finding-GO loops evidenced twice. |
| §11.4.135 standing regression-guard suite | PARTIAL | New tests accompany each fix (e.g. `skill_create_draft_test.go` for the MCP draft-invariant), but no curated "regression-guard registry" cross-referencing which test guards which closed G-item was found. |
| §11.4.138 operator-escape bluff-audit | N-A | No operator-escape (bypassed guard) event is documented. |
| §11.4.139 fresh-process clean-artifact runtime-signature | PENDING-AT-COMPLETION | Same basis as §11.4.108's ARTIFACT/RUNTIME-ON-CLEAN-TARGET gap. |
| §11.4.140 action-prefix system | N-A | Mechanism-level (constitution-side); not a distinguishing project-specific item. |
| §11.4.141 token-efficiency mandate | UNCONFIRMED | Cannot verify prior sessions' token-efficiency from a static snapshot. |
| §11.4.142 universal code-review mandate | COMPLIANT | Same evidence as §11.4.125/§11.4.209 — the Fable-xhigh review is applied to "every change", per CONTINUATION's stated binding constraint, with no stated carve-out. |
| §11.4.145 independent multi-angle impact-research | COMPLIANT | The register itself is explicitly an independent adversarial audit (not self-review) — see §11.4.118 evidence. |
| §11.4.146 reproduce-first + extend-to-all-cases workflow | COMPLIANT | G02/G03/G05/G16/G21 were fixed together as one related batch rather than G02 alone, then G05 separately, etc. — the "extend to all related cases" discipline. |
| §11.4.147 crashed-agent respawn-until-complete registry | UNCONFIRMED | No agent-crash event is documented in the available materials; cannot confirm the registry mechanism was engaged because it was never needed, versus not wired. |
| §11.4.148/§11.4.171 workable-item integrity + comprehensive descriptions | COMPLIANT (description quality) / VIOLATION (mechanical wrapper — see §11.4.93) | Each G-item's description (Category/Severity/Evidence/Why-it-matters/DECISION/Test-coverage/Challenges/HelixQA) is genuinely comprehensive — it exceeds the §11.4.91 "≥6 words, self-contained" bar by a wide margin. The mechanical DB+ATM-id wrapper is the missing piece, tracked separately as **G47/G44**. |
| §11.4.149 per-workable-item testing diary | VIOLATION | No append-only `test_diary` exists (tied to §11.4.93 DB absence). → folded into **G47**. |
| §11.4.150 deep multi-angle web research before declaring fixed | COMPLIANT | Same evidence as §11.4.8/§11.4.123. |
| §11.4.151 project-prefixed release-tag naming | N-A | No release tag has been cut yet. |
| §11.4.152 Crashlytics monitoring mandate | N-A | No Firebase Crashlytics integration. |
| §11.4.153 per-feature Status+Status_Summary+video-confirmation | N-A | No feature has shipped to have a Status doc for yet (correctly early-stage, not silently skipped — the register tracks the gating work). |
| §11.4.154/§11.4.155 window-scoped capture + naming | N-A | No recordings exist (no UI/video feature). |
| §11.4.156 CI/CD automation disabled mandate | COMPLIANT | No `.github/workflows/` or `.gitlab-ci.yml` found anywhere in the repo outside the constitution submodule. |
| §11.4.157 GEMINI.md lockstep mandate | VIOLATION | Top-level repo has `CLAUDE.md` + `AGENTS.md` only — no `QWEN.md`, no `GEMINI.md`. → **G49**. |
| §11.4.161 rootless container runtime mandate | VIOLATION | `project/Makefile:36`: `CONTAINER_RUNTIME ?= docker` (rootful Docker as the default; overridable but not defaulting to Podman). → **G46** (MUST-FIX-NOW — one-line default change). |
| §11.4.162 OpenDesign UI design system mandate | PENDING-AT-COMPLETION (planning COMPLIANT) | No UI exists yet to apply tokens to; `research/opendesign_incorporation.md` documents forward planning for R12/R18 doc-delivery diagrams — the planning discipline is present ahead of need. |
| §11.4.164 constitution auto-propagation hook system | N-A | Inherited mechanism; this project does not need to reimplement it. |
| §11.4.165 independent verification agent mandate | COMPLIANT | The Fable-xhigh reviewer is precisely this: structurally independent, iterates to zero-finding GO. |
| §11.4.166 (REPEALED — Semgrep) | N-A | Anchor repealed; no action required. |
| §11.4.167 big-work-item feature-branch lifecycle | N-A | Work proceeds directly on `main`, operator-authorized (`IMPLEMENTATION_PLAN.md` header: "Repo branch: `main` (operator authorized direct-on-main)") — no `feat/` stream has been opened, correctly not yet applicable. |
| §11.4.168 exported-document visual validation | N-A | No exports exist yet (tied to §11.4.65/§11.4.12). |
| §11.4.169 comprehensive test-type coverage | VIOLATION | Same basis and remediation as §11.4.27/§11.4.85. → folded into **G41**. |
| §11.4.172 production-readiness planning w/ timeline | PARTIAL | `IMPLEMENTATION_PLAN.md` has phases (P0–P13) with tasks/subtasks/evidence gates, but no explicit realistic wall-clock timeline projection was found. |
| §11.4.173 containerized+distributed build mandate | VIOLATION | Same root cause as §11.4.76/§11.4.161 — build/deploy is orchestrated via raw `docker`/`docker-compose`, not the mandated containers submodule. → folded into **G46**. |
| §11.4.174 shared-host process-ownership verification | N-A | Single-operator host, not a shared multi-tenant host for this project's scope. |
| §11.4.176/§11.4.178/§11.4.179/§11.4.187/§11.4.188/§11.4.191/§11.4.192 multi-track engine family | N-A | This project is a single-track effort (one Go-mutator lane); the formal multi-track ruler orchestration, exactly-once claim registry, and device-lock mechanisms are designed for ≥2 concurrent tracks on shared hardware, which does not describe this project's current operating mode. |
| §11.4.180 stale-lock auto-reap | N-A | No lock-based commit/push wrapper specific to this project was found requiring this. |
| §11.4.181 branch-name consistency | N-A | Single `main` branch; no feature-branch naming collision risk yet. |
| §11.4.182 track+branch label on every agent dispatch | COMPLIANT | This very audit task was dispatched with the label prefix `(T1/main - claude1)`, confirming the mechanism is actively wired at the dispatch layer for this project. |
| §11.4.183 maximum multi-agent utilization per work-stream | PARTIAL | CONTINUATION.md's honest fleet-sizing narrative (§11.4.103 evidence) shows genuine engagement with this discipline, bounded correctly by §11.4.42 priority-ordering rather than blind agent-count maximization. |
| §11.4.184 SonarQube CLI installed + PATH-discoverable | COMPLIANT | `command -v sonar-scanner` resolves in a fresh login shell (`bash -lc`) to `/home/milos/Factory/software/sonar-scanner/bin/sonar-scanner`; `~/.bashrc:287-288` durably exports the PATH addition, explicitly tagged `# §11.4.184`. |
| §11.4.185 manual QA final confirmation | PENDING-AT-COMPLETION | Correctly not yet applicable — no scope has reached a release-candidate state requiring manual QA sign-off. |
| §11.4.186 anti-divergence cross-document consistency gate | PENDING-AT-COMPLETION | No mechanical cross-doc consistency gate has been authored for this project's multiple docs (REQUIREMENTS/IMPLEMENTATION_PLAN/SPEC/register); the docs appear internally consistent by manual authorial cross-referencing (e.g. REQUIREMENTS.md explicitly reconciles against IMPLEMENTATION_PLAN.md), but this has not been mechanically verified/gated. |
| §11.4.189 most-reopened extra-depth live-testing scrutiny | N-A | No live-testing phase has begun yet. |
| §11.4.190 website engineering-quality mandate | N-A | No website surface exists yet for this project; will become live once the R18 documentation-delivery website (referenced in CONTINUATION.md) is built. |
| §11.4.194 exhaustive all-scenario code-review mandate | COMPLIANT | G01's attempt-1→attempt-2 cycle is a direct instance: the first review caught a real multi-factor defect (co-bound listener) a shallower review had missed, and the second review re-enumerated the full auth scenario-space including both degenerate cells (`api_keys=[""]`, `["${UNSET}"]`). |
| §11.4.195 branch-taxonomy governance | N-A (operator-authorized exception) | `IMPLEMENTATION_PLAN.md` header explicitly records "Repo branch: `main` (operator authorized direct-on-main)" — an explicit, documented authorization rather than an unremarked taxonomy violation. |
| §11.4.196/§11.4.198 native-alias-first orchestration + default mechanisms | N-A | Conductor/orchestration-layer concern (which alias drives which track), not a property of the audited project's own source tree. |
| §11.4.197 research/kicked-off-work completion mandate | PARTIAL (actively tracked, not yet fully wired — the honest, correct state) | This is the anchor most directly at stake for this project's overall posture: G03 documents that the flagship validation/autoexpand pipelines are still dead code, and G01's O3 (dead `internal/api.Server` consolidation) is explicitly deferred with a cited reason (§11.4.101/§11.4.124) rather than silently dropped. Because every such item is tracked in the register with a stated next-step, this is NOT the forbidden "silently un-wired in the backlog" state — it is the "in progress, tracked" state the anchor permits. It will graduate to full COMPLIANT only once G03/O3 are wired, or to N-A-via-closure if explicitly Obsoleted. |
| §11.4.199 exact-reproduction-sequence mandate | COMPLIANT | G01's re-verification narrative re-drives the exact scenario space (both `http`/`both`/default modes, both degenerate `api_keys` cells) rather than a hand-rolled approximation. |
| §11.4.200 deploy-target-isolated-and-verified mandate | N-A | No deploy/flash tooling exists yet for this project. |
| §11.4.201 guard-asserts-real-condition mandate | N-A | This project has not yet authored any host-safety/resource guard of its own to assess. |
| §11.4.202 reporting directives (ISSUE/BUG/TASK) | PARTIAL | The mechanism exists and is functional at the constitution layer (gate `CM-REPORTING-DIRECTIVES` PASSes — see §2); this project's own G0x items are hand-authored Markdown, not created via the `ISSUE:`/`BUG:`/`TASK:` directive → DB pipeline. The mechanism is available but unused for this project's own tracking. |
| §11.4.207 continuum instant-resume engine | N-A (mechanism present upstream, not adopted downstream) | The engine exists in `constitution/submodules/continuum/` (see §2 gate run); this project uses the simpler §12.10/§11.4.131 `CONTINUATION.md` mechanism instead, which §11.4.207 explicitly states it "EXTENDS — never replaces". Not adopting the continuum snapshot mechanism is therefore not itself a violation of this project. |
| §11.4.208 operator-request-history document | VIOLATION | No `docs/requests/history.md` (or equivalent) exists anywhere under this project or the parent repo. → **G50** (MUST-FIX-NOW per explicit task framing — mechanical, one doc + one helper invocation). |
| §11.4.209 code-review MUST run on Fable-xhigh | PARTIAL (COMPLIANT per project self-report; independent verification limited) | CONTINUATION.md's binding constraints and multiple G-register STATUS blocks explicitly cite "§11.4.209 Fable-xhigh review" / "Fable-xhigh GO" as having gated every commit in the P0.5 lane. I cannot independently verify from static file inspection which model actually executed those past reviews (no session/model metadata is recoverable from committed files) — this is recorded as an **honest audit-coverage gap**, not a disproof. The project's own documentation is internally consistent with compliance. |

### 1.2 Grouped N-A anchors (hardware / device / audio / video / messenger / mobile-flash — none apply to a Go knowledge-graph backend)

§11.4.46 (post-flash validation), §11.4.47 (Firebase), §11.4.48/.49 (UI-driven/dual-approach video testing), §11.4.51 (live-ADB-first), §11.4.68 (audio sink-side), §11.4.72 (audio top-priority), §11.4.104/.105 (messenger identity/intent), §11.4.107 (AV liveness techniques), §11.4.111 (resolve-by-stable-name for enumerated hardware), §11.4.117 (CV/OCR pixel oracle), §11.4.119 (single-resource-owner hardware partitioning), §11.4.128 (device recording), §11.4.133 (target+hardware safety), §11.4.136/.137/.143 (real-content/subtitle/streaming-app journeys), §11.4.144 (tracked-device availability-following), §11.4.152 (Crashlytics), §11.4.154/.155 (window-scoped recording naming), §11.4.158/.159/.160/.163 (video-recording + vision-verification pipelines), §11.4.193 (anti-blind-typing UI mandate) — **reason:** this project has no physical device, no audio/video output, no messenger surface, and no UI to drive at this phase.

### 1.3 Grouped N-A anchors (single-track project / no multi-track ruler orchestration in use)

§11.4.176, §11.4.178, §11.4.179, §11.4.187, §11.4.188, §11.4.191, §11.4.192, §11.4.196, §11.4.198 — **reason:** these anchors govern coordination across ≥2 concurrently-running tracks/agents sharing hardware and a ruler orchestrator; this project currently runs a single serialized Go-mutator lane (CONTINUATION.md's own honest disclosure). §11.4.182 (track+branch label) is the one anchor from this family that DOES apply at the dispatch layer and is COMPLIANT (§1.1).

### 1.4 Grouped N-A anchors (release/deployment lifecycle not yet reached)

§11.4.38, §11.4.40, §11.4.41 (N-A — no force-push occurred), §11.4.46, §11.4.114, §11.4.129, §11.4.130, §11.4.151, §11.4.189, §11.4.200 — **reason:** no release tag, no redeploy cycle, no huge-blocker release-validation event has occurred yet; this is the correct phase-appropriate state for a P0.5 remediation-stage project, not an evasion.

### 1.5 Grouped PENDING-AT-COMPLETION anchors (will be satisfiable once the system is further built/deployed)

§11.4.2, §11.4.25, §11.4.50, §11.4.52 (already individually flagged VIOLATION above for the *current* test shape, but its *full* scope is PENDING), §11.4.69, §11.4.81, §11.4.83, §11.4.98, §11.4.108 (ARTIFACT/RUNTIME-ON-CLEAN-TARGET sub-layers), §11.4.110, §11.4.139, §11.4.153, §11.4.162 (full UI application), §11.4.168, §11.4.170, §11.4.185, §11.4.186, §11.4.190 — what will satisfy each: a live deployed instance, a captured `docs/qa/<run-id>/` transcript per shipped feature, a stress/chaos suite, a coverage ledger, and (for §11.4.185) a QA-team manual sign-off before any release tag.

---

## 2. Gate-execution results (real runs, real exit codes — §11.4.201)

All gates below were executed with `bash <gate-script>`, not grep-simulated. Where a gate's default root resolution swept in unrelated sibling projects under `/home/milos/Factory/projects/tools_and_research/` (a much wider directory than this repo), I re-ran with `CONSUMER_ROOT` explicitly scoped to `helix_skills` to get the project-accurate signal; both runs are reported for transparency.

| Gate | Default-root run | Scoped (`CONSUMER_ROOT=.../helix_skills`) run | Verdict for this project |
|---|---|---|---|
| `cm_reporting_directives.sh` | EXIT 0 — PASS (4/4 invariants: escapes-to-literal, both honest SKIP reasons defined, tracker PASS gated on real exit 0) | not re-run (root-independent — checks the constitution engine itself) | **COMPLIANT** — the §11.4.202 mechanism is functional; this project simply does not use it yet for its own G-items (§1.1). |
| `cm_continuum_resume_engine_present.sh` | EXIT 1 — FAIL. `submodules/continuum/test/e2e/e2e_test.go` contains a project-specific literal (decoupling violation, §11.4.28(B)); `go test -race ./...` PASS; `continuum selfcheck` good=PASS/bad=FAIL/negctrl=PASS all PASS. | not re-run (checks the constitution submodule's own engine, not this project) | **Out of this project's scope** — the failure is inside the shared `continuum` engine (constitution submodule), not something this MVP project owns or can remediate. Reported for completeness, not counted in this project's violation list. |
| `cm_covenant_114_202_propagation.sh` | EXIT 1 — 1167 MISSING across `/home/milos/Factory/projects/tools_and_research/` (sweeps in unrelated sibling repos: `claude_toolkit`, `code_server`, `tmux`, etc.) | EXIT 1 — 4 PRESENT (constitution's own 4 mirror files), **2 MISSING**: `helix_skills/CLAUDE.md`, `helix_skills/AGENTS.md` | **VIOLATION, confirmed at correct scope** — see §5 below for the honest tension this reveals. |
| `cm_covenant_114_207_propagation.sh` | EXIT 1 (same wide-root sweep) | EXIT 1 — same "2 MISSING": `helix_skills/CLAUDE.md`, `helix_skills/AGENTS.md` | Same as above. |
| `cm_covenant_114_199/.200/.201_propagation.sh` (×3) | EXIT 1 (same wide-root sweep) | EXIT 1 — same "2 MISSING" pattern each | Same as above — this is a **uniform pattern across every §11.4.167-through-§11.4.209-era propagation gate**, not per-anchor noise. |
| `cm_covenant_114_196/.191/.187/.176/.167_propagation.sh` (×5) | not separately logged (same pattern confirmed) | EXIT 1 — same "2 MISSING" pattern each | Same as above. |
| `cm_cli_agent_plugins_wired.sh` | not applicable (constitution-root gate) | EXIT 0 — PASS (9/9 invariants) when scoped to `helix_skills` (root resolves to `constitution/` internally) | **COMPLIANT** — constitution-level plugin wiring is intact; not a project-specific signal either way. |

**Honest summary of the gate sweep:** every `CM-COVENANT-114-NNN-PROPAGATION` gate tested (10 of them) FAILs when correctly scoped to this repository, and FAILs for the *same* reason each time: the project's top-level `CLAUDE.md`/`AGENTS.md` use the sanctioned `@import` pointer pattern (per the constitution's own "How inheritance works" section) rather than literal anchor-text restatement, and the gate scripts do a literal `grep` for the anchor string — they do not resolve `@import` directives. This is a genuine, mechanically-reproducible FAIL, not a false positive from mis-scoping (I verified the gate's own root-resolution logic at `cm_covenant_114_202_propagation.sh:68` confirms this). Whether the correct remediation is "restate every anchor literally in the project CLAUDE.md/AGENTS.md" (defeats the purpose of the `@import` pointer) or "teach the gate to resolve `@import`" (a constitution-submodule-level fix, out of this project's remediation scope) is an **operator decision**, not something this audit unilaterally resolves. See §5.

I did not run `cm_opendesign_ui_system.sh`, `cm_subsystem_shortcuts.sh`, or `cm_track_branch_label.sh` — these are lower-signal for a pre-UI, single-track Go backend and were deprioritized under the effort budget for this audit; their absence from this run is an honest coverage gap (§4), not a claimed PASS.

---

## 3. Specific dedicated verdicts (as requested)

### §11.4.93 + §11.4.95 — SQLite workable-items DB

**Verdict: VIOLATION, with earned nuance.** The constitution's inherited Go engine (`constitution/scripts/workable-items/`, with `cmd/workable-items`, `schema.sql`, a built `bin/workable-items` binary) is present and ready to consume — this project could adopt it with **zero new engineering**, only an invocation (`workable-items sync md-to-db` against `GAPS_AND_RISKS_REGISTER.md`). It has not done so: no `docs/workable_items.db` exists anywhere under this project. The project instead tracks its 27+ findings (G01–G28+, growing) in a single well-structured Markdown register with rich per-item narrative (Category/Severity/Evidence/Why-it-matters/DECISION/Test-coverage/Challenges/HelixQA fields) that substantively satisfies the *content* bar of §11.4.148/§11.4.171 even though it lacks the *mechanical* DB/summary/export wrapper. This is a real, no-escape-hatch violation of §11.4.93 ("no escape hatch... a project that ships without this is not delivering a stable product" — applied here to the tracking layer itself), but the fix is genuinely cheap since the engine already exists. → **G47**.

### §11.4.54 — ATM-NNN stable ids

**Verdict: VIOLATION, operator-decision-needed (not unilaterally fixable).** The project uses `G01`–`G28`+ (register), `R1`–`R21` (requirements), `P0`–`P13` (plan phases) identifiers instead of the mandated `[ATM-NNN]` scheme. These ids are already deeply cross-referenced throughout 20+ research documents, `CONTINUATION.md`, `IMPLEMENTATION_PLAN.md`, and `REQUIREMENTS.md` — a mechanical rename to `ATM-NNN` at this stage carries real risk of broken cross-references (§11.4.6: renaming without proof every reference was updated is itself forbidden) for a benefit (uniform id-namespace) whose priority is lower than the CRITICAL/HIGH security and correctness items the register is actively tracking. Two honest remediation paths exist: **(a)** operator-authorizes `G0x`/`R0x`/`P0x` as documented, project-scoped aliases of the ATM-NNN pattern (cheapest, lowest-risk — a one-line addendum to this project's CLAUDE.md under "Project-specific rules" citing this exact tension), or **(b)** run the existing `assign_atm_ticket_ids.sh`-class tooling to mint real `ATM-NNN` ids in parallel (not replacing) the G0x ids inside the future `workable_items.db` (§11.4.93), so both id schemes co-exist without a destructive rename. → **G44**, flagged for §11.4.66 operator clarification, not a unilateral fix.

---

## 4. Prioritized violation-remediation list

MUST-FIX-NOW items are cheap, mechanical, and carry real risk if left (security-adjacent enforcement gap, container-runtime posture, missing audit trail). PENDING items are correctly deferred to feature-completion and are listed for tracking, not urgency.

| Gap id | Anchor(s) | Severity | Fix direction | Class |
|---|---|---|---|---|
| **G38** | §11.4.12, §11.4.22, §11.4.59, §11.4.60, §11.4.65, §11.4.73, §11.4.106 | HIGH | Wire the Docs Chain engine (already declared in `helix-deps.yaml`) for this project's doc set (`REQUIREMENTS.md`, `IMPLEMENTATION_PLAN.md`, `SPEC.md`, `GAPS_AND_RISKS_REGISTER.md`, `CONTINUATION.md`, `plan.md`), producing `.html`/`.pdf` siblings kept in sync via a `--docs-only` lightweight commit path. | MUST-FIX-NOW |
| **G39** | §11.4.15, §11.4.16, §11.4.19, §11.4.33, §11.4.34, §11.4.55 | MEDIUM | Either (a) migrate the register's `Category`/`STATUS`-prose fields to the closed `{Status}`/`{Type}` vocabularies + add a `**Reopened-Details:**` line per reopened item, or (b) explicitly document the register's own vocabulary as a project-specific, operator-approved substitute per §11.4.17/§11.4.35, citing this audit. | MUST-FIX-NOW (documentation-only fix, cheap) |
| **G40** | §11.4.18 | MEDIUM | Author `docs/scripts/<name>.md` companions for the 13 scripts under `project/scripts/`, plus add the in-source documentation block (Purpose/Usage/Inputs/Outputs/Side-effects/Dependencies) to each script header. | MUST-FIX-NOW |
| **G41** | §11.4.27, §11.4.52, §11.4.85, §11.4.169 | HIGH | Extend the 10 existing tests toward the 13-test-type floor: add integration tests against a real Postgres+pgvector instance (not mocked), a security-scenario suite (the auth/CORS scenario-space is already partially covered — extend it), and at least one stress test (N≥100 iterations) + one chaos test (process-death or network-fault injection) per §11.4.85's closed-set families. | MUST-FIX-NOW (start with security+stress, given G01's history) |
| **G42** | §11.4.29 | LOW | Rename `Agent_AI_Skill_Tree_Development` → `agent_ai_skill_tree_development`, updating every reference atomically (CONTINUATION.md, cross-links, any absolute paths in research docs), with a regression check that every reference resolves post-rename (§11.4.29's own test-coverage-of-renames requirement). | MUST-FIX-NOW (mechanical, but touches many files — schedule deliberately) |
| **G43** | §11.4.44, §11.4.73 | MEDIUM | Add the `**Revision:** N` / `**Last modified:** ISO8601` header to `REQUIREMENTS.md`, `IMPLEMENTATION_PLAN.md`, `SPEC.md`, `GAPS_AND_RISKS_REGISTER.md`, `plan.md`, and every `research/*.md` file; wire `doc_revision_bump.sh`-equivalent into the commit path so it auto-bumps. | MUST-FIX-NOW (cheapest fix in this list) |
| **G44** | §11.4.54 | MEDIUM | Operator decision needed (§11.4.66) — see §3 above for the two honest remediation paths (alias G0x as documented, or mint parallel ATM-NNN ids without renaming). | MUST-FIX-NOW to the extent of *surfacing the decision*; the id-scheme change itself follows the operator's choice. |
| **G45** | §11.4.57, §11.4.59 | MEDIUM | Add the `Tracked-Items + Status Documents` section (with the required markers) to `project/README.md`, wired to regenerate via the same sync wrapper as G38. | MUST-FIX-NOW |
| **G46** | §11.4.76, §11.4.161, §11.4.173 | HIGH | Change `project/Makefile:36` default `CONTAINER_RUNTIME ?= docker` → `CONTAINER_RUNTIME ?= podman`; evaluate adopting the `vasic-digital/containers` submodule's `pkg/boot`/`pkg/compose` layer instead of raw `docker-compose.yml` invocations (or document, per §11.4.161's own exception clause, why this platform genuinely has no rootless option — unlikely to be true here). | MUST-FIX-NOW (the one-line default change is trivial; full containers-submodule adoption is a larger follow-up) |
| **G47** | §11.4.93, §11.4.95, §11.4.63, §11.4.149 | HIGH | Invoke the existing `constitution/scripts/workable-items` Go binary against `GAPS_AND_RISKS_REGISTER.md` to populate `docs/workable_items.db`; commit the DB (tracked, never gitignored per §11.4.95). | MUST-FIX-NOW (engine already built; this is adoption, not engineering) |
| **G48** | §11.4.109 | HIGH | Wire `constitution/scripts/hooks/guard-forbidden-commands.sh` as a `PreToolUse` hook in `.claude/settings.local.json` (currently has zero hooks configured), and reference `constitution/docs/AGENT_GUARDRAILS.md` in the project's own `docs/AGENT_GUARDRAILS.md` (or a project-local equivalent). | MUST-FIX-NOW |
| **G49** | §11.4.157 | LOW | Author `QWEN.md` and `GEMINI.md` at the `helix_skills` top level as lockstep mirrors of `CLAUDE.md`/`AGENTS.md` (both currently absent). | MUST-FIX-NOW (mechanical) |
| **G50** | §11.4.208 | MEDIUM | Create `docs/requests/history.md` at the `helix_skills` top level (project-local, per §11.4.35's rule-universal/document-project-local split) and begin appending an entry per operator request going forward, with `UNKNOWN` fields for unrecoverable pre-mandate history rather than fabricated values. | MUST-FIX-NOW (explicitly named in this audit's own task framing) |
| — | §11.4.202 gate literal-vs-`@import` tension (see §2) | MEDIUM | **Operator decision needed**: either restate anchor literals in `helix_skills/CLAUDE.md`/`AGENTS.md` (defeats the `@import` pointer's purpose) or extend the propagation-gate scripts to resolve `@import` as satisfying (a constitution-submodule-level change, out of this project's unilateral remediation authority). Not assigned a project-local gap id since the fix, whichever direction, lands in the constitution submodule or as an explicit CLAUDE.md exception, not in this MVP project's own tree. | Flag-for-operator, not a project gap id |

**PENDING-AT-COMPLETION items** (not assigned gap ids — tracked via existing G03/O3/G14-X1 entries and the grouped list in §1.5): full autonomous-validation coverage ledger (§11.4.25/§11.4.52), `docs/qa/` transcripts per shipped feature (§11.4.83), the ARTIFACT/RUNTIME-ON-CLEAN-TARGET layers of §11.4.108/§11.4.139, manual QA sign-off (§11.4.185), the cross-document consistency gate (§11.4.186), and website quality (§11.4.190) once the R18 doc-delivery site exists.

---

## 5. Honest gaps in this audit (§11.4.6)

1. **§11.4.37 fetch-before-edit**: this audit's own first git action should have been `git fetch --all --prune` before reading anything; it was not. Remote-ahead state was not confirmed. This is a gap in *this audit's own procedure*.
2. **§11.4.67 script-shell-parseability, §11.4.78/§11.4.79/§11.4.80 CodeGraph indexing currency, §11.4.87/§11.4.88/§11.4.89/§11.4.97/§11.4.126/§11.4.141/§11.4.147 (loop/background/token/crash-registry mandates)**: none of these are verifiable from a static, read-only file-tree snapshot — they describe *session behavior* (was work backgrounded, was the loop engaged from the first prompt, was an agent-crash ever handled) that leaves no artifact this audit can inspect after the fact. Marked **UNCONFIRMED**, not COMPLIANT and not VIOLATION — the honest third state.
3. **§11.4.115 RED-polarity-switch**: test file *names* are consistent with the remediation narrative, but I did not open and read the test source to confirm the literal `RED_MODE` polarity-switch pattern versus a plain post-fix-only regression test. Recommended as a follow-up spot-check, not claimed either way.
4. **§11.4.209 (which model actually ran past reviews)**: the project's own documentation consistently and repeatedly claims Fable-xhigh review gated every commit in the P0.5 lane. I have no independent means — from a static file audit — of confirming which model executed those reviews in past sessions; I can only confirm the *project's own record* is internally consistent and non-contradictory. This is reported as the project's self-attested state, not independently re-derived.
5. **The propagation-gate wide-root default (§2)**: `CONSUMER_ROOT`/`--root` defaults to `..` relative to the invoking shell's cwd, not to the script's own location or to any auto-detected repository root. Run naively, this silently scopes the audit to a directory containing dozens of unrelated projects (1167 missing carriers). I caught and corrected this by explicitly setting `CONSUMER_ROOT`, and I flag it here because a less careful gate invocation — including by a future agent — would produce a wildly misleading "1167 violations" headline that has nothing to do with this project.
6. **Deep-dive gates not run**: `cm_opendesign_ui_system.sh`, `cm_subsystem_shortcuts.sh`, `cm_track_branch_label.sh` were not executed against this project (deprioritized under the effort budget — see §2). Their omission is a coverage gap, not a claimed PASS or FAIL.
7. **No live build/test execution performed by this audit**: per the task's own scope ("write NO code, modify NO existing file"), I did not run `go build`/`go vet`/`go test` myself to re-verify the register's claim of "build=0, vet=0, test=0 (all pkgs PASS)" for the currently-uncommitted G02 working-tree state (per CONTINUATION.md, a Go change was uncommitted at time of last update, gated on a pending Fable-xhigh review). This audit's compliance ledger is therefore evaluated against the tree AS COMMITTED plus the documented in-flight state, not a freshly-reproduced build.

---

## 6. Data summary (for the dispatching agent)

- **File path:** `docs/research/mvp/Agent_AI_Skill_Tree_Development/research/constitution_compliance_audit.md`
- **Counts:** COMPLIANT 47 · VIOLATION 34 · PENDING-AT-COMPLETION 26 · N-A 108 · UNCONFIRMED 15
- **MUST-FIX-NOW violations (with severity + gap id):** G38 docs-sync/exports (HIGH) · G39 status/type/reopens vocabulary (MEDIUM) · G40 script docs (MEDIUM) · G41 test-type coverage incl. stress/chaos (HIGH) · G42 project-dir naming (LOW) · G43 revision headers (MEDIUM) · G44 ATM-NNN id decision (MEDIUM, operator-gated) · G45 README doc-link section (MEDIUM) · G46 rootless-container default (HIGH) · G47 SQLite workable-items DB adoption (HIGH) · G48 PreToolUse guard-hook wiring (HIGH) · G49 QWEN/GEMINI mirrors (LOW) · G50 request-history doc (MEDIUM).
- **Gate-execution results:** `CM-REPORTING-DIRECTIVES` PASS (exit 0). `CM-CONTINUUM-RESUME-ENGINE-PRESENT` FAIL (exit 1) — failure is inside the shared constitution-submodule engine (a project-specific literal in its own e2e test), out of this project's remediation scope. Ten `CM-COVENANT-114-NNN-PROPAGATION` gates (167/176/187/191/196/199/200/201/202/207) each FAIL (exit 1) when correctly scoped to `helix_skills`, uniformly because the project's top-level `CLAUDE.md`/`AGENTS.md` use the sanctioned `@import` pointer form rather than literal anchor-text restatement, which the gates' literal-grep implementation does not resolve — flagged for operator decision, not a project-code fix. `CM-CLI-AGENT-PLUGINS-WIRED` PASS (exit 0, constitution-level, not project-specific).
- **§11.4.93/§11.4.95 verdict:** VIOLATION — no SQLite `docs/workable_items.db` exists; the constitution's own Go engine to build one is already present and unused (adoption-only fix, G47).
- **§11.4.54 verdict:** VIOLATION, operator-decision-needed — project uses G0x/R0x/P0x ids, not ATM-NNN; two non-destructive remediation paths identified (G44).
- **Honest audit-coverage gaps:** §11.4.37 (this audit itself skipped a pre-read `git fetch`), 12+ session-behavior anchors unconfirmable from a static snapshot (§11.4.67/.78-.80/.87-.89/.97/.115/.126/.141/.147/.209), 3 gate scripts not run, and no independent re-run of `go build/vet/test` against the current in-flight tree state.
