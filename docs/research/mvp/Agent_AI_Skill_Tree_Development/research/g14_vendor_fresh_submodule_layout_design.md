# G14 — Vendor-Fresh Submodule Layout + Incorporation Plan (HelixKnowledge Skill Graph)

**Revision:** 1
**Last modified:** 2026-07-15T18:36:44Z
**Status:** design-only. No submodule added, no `.gitmodules` touched, no file under
`project/` modified, no git operation performed by this document's author.
**Authority cited:** `constitution/CLAUDE.md` §11.4.28 (esp. §11.4.28(C)), §11.4.29,
§11.4.31, §11.4.36, §11.4.6, §11.4.74, §11.4.28(B); project docs `CONTINUATION.md`
Rev 8 (operator decision), `project/helix-deps.yaml` (existing manifest),
`project/scripts/sync_submodules.sh` (existing tooling), `research/g14_x1_submodule_policy_decision.md`
(the superseded decision package), `research/r22_full_catalogue_incorporation_design.md`,
`research/g05_g11_reuse_debate_dag_design.md`, `research/helix_family_reusable_practices.md`,
`research/g40_workable_items_db_adoption_design.md`, `research/g43_docs_chain_export_wiring_design.md`.

---

## 1. Goal + operator decision

**Verbatim operator decision** (`CONTINUATION.md` Rev 8, 2026-07-15T18:24:57Z):

> "G14/X1 submodule policy = VENDOR FRESH under this project — each dep a fresh git
> submodule under our own `submodules/<snake_case>/` (§11.4.28(C)), REVERSES the
> autonomously-adopted Option-A single-canonical; R22-catalogue ADOPT verdicts stand,
> layout becomes project-local vendoring + `install_upstreams` + recursive helix-deps.
> open-design → `submodules/open_design/`."

**What this reverses, precisely.** `research/g14_x1_submodule_policy_decision.md` (Revision
2) had autonomously adopted **"Option A"**: a dependency already vendored at a
**parent-ecosystem-root sibling project** (e.g. `../helix_code/submodules/llm_provider/`)
is **consumed by reference** from that sibling checkout — the skill-system project vendors
**nothing locally**, and `project/scripts/sync_submodules.sh` implements exactly this
"search `--ecosystem-root` dirs first, only vendor fresh as a fallback" logic (read
directly, `project/scripts/sync_submodules.sh:1-30`). The operator's decision **reverses**
that: every dependency this project needs becomes its **own, independent, normal git
submodule**, added directly under **this project's own** `submodules/<name>/` (or a
documented root-level exception, §6), resolved from its **own canonical upstream SSH URL**
— **never** by pointing at, or depending on the presence/freshness of, `helix_code`,
`helix_ota`, `helix_terminator`, or any other sibling checkout. The skill-system repo must
be **fully self-contained and clonable in isolation** (`git clone --recurse-submodules`
alone reconstructs every dependency) — a genuinely independent Helix-family project, not
one that silently assumes a sibling checkout sits next to it on disk.

**What this does NOT reverse.** Two mechanisms that look superficially similar to Option A
are **orthogonal** to the G14/X1 question and remain unchanged by this decision:

1. **The `constitution/` submodule itself.** Every Helix-family project (all 7 surveyed
   siblings, confirmed `research/helix_family_reusable_practices.md` §9) vendors
   `HelixDevelopment/HelixConstitution.git` at the identical path `constitution/` — this is
   universal governance inheritance (§11.4.35), not a "consume dependency X from sibling
   project Y" pattern. The skill-system project already does this today (it is nested
   inside `helix_skills`, which already vendors `constitution/` at its own root — verified
   directly, `helix_skills/.gitmodules`). The eventual extracted, standalone skill-system
   repo vendors `constitution/` at its own root too, per this universal precedent.
2. **The §11.4.28(C) constitution-submodule depth-1 carve-out.** `token_optimizer`,
   `session_orchestrator`, and `continuum` already live at
   `constitution/submodules/token_optimizer/` etc. — **inside** the one governance
   submodule every project equally vendors, not inside a sibling **project's** tree. A
   project consuming `token_optimizer` from `constitution/submodules/token_optimizer/` is
   not "reaching into another project's copy" (the failure mode G14/X1 was about) — it is
   reading the same shared governance submodule everyone already has. This mechanism is
   **not gated on G14/X1** and is **not reversed** by the vendor-fresh decision.

Concretely: `llm_provider`, `http3`, `containers`, `dag_orchestrator`, `pipeline_runtime`,
`vision_engine`, `panoptic`, `llms_verifier`, `helix_llm`, `helix_agent`, `embeddings`,
`helix_qa`, `challenges`, `open_design`, `design_system` all move from "consume-by-reference
from `helix_code`/`helix_ota`" to **"vendor fresh under this project's own `submodules/`"**.
`token_optimizer` and (pending, §7) `docs_chain` stay **consumed from `constitution/`**,
because that was never the Option-A pattern in the first place.

**Physical-location caveat (honest, §11.4.6).** `project/` today lives at
`docs/research/mvp/Agent_AI_Skill_Tree_Development/project/` **inside** the `helix_skills`
repository — `git -C project rev-parse --show-toplevel` resolves to `helix_skills`'s own
root, confirming `project/` is a plain directory, **not** its own git repository or
submodule, at this design-research stage. This document therefore designs the layout for
the **eventual extracted, standalone skill-system repository** (module
`github.com/helixdevelopment/skill-system`, per `IMPLEMENTATION_PLAN.md` R1 "harden
extracted Go backend") whose root will be wherever `project/` is extracted to (`P13.T3`
packaging). Every path below (`<skill-system-root>/submodules/<name>/`) is relative to that
future root, not to `helix_skills`'s root. No submodule is added and no extraction is
performed by this pass — it is a design document per the task contract.

---

## 2. Dependency inventory table

Verdicts are carried forward unchanged from the already-accepted incorporation designs
(`r22_full_catalogue_incorporation_design.md`, `g05_g11_reuse_debate_dag_design.md`,
`helix_family_reusable_practices.md`) — this document changes only the **layout mechanism**
(vendor-fresh vs. consume-by-reference), never a verdict. Every `ssh_url` below was
independently re-confirmed live via `gh repo view <org>/<name> --json name,sshUrl,...`
(2026-07-15, this pass) except where marked UNCONFIRMED/CONFIRMED-ABSENT.

### 2.A — VENDOR FRESH under `<skill-system-root>/submodules/<name>/` (this project's own copy)

| # | Name | snake_case dir | Org / upstream `ssh_url` | Verdict | Go-replace? | `helix-deps.yaml` needed? |
|---|---|---|---|---|---|---|
| 1 | LLMProvider | `llm_provider` | `git@github.com:HelixDevelopment/LLMProvider.git` | **ADOPT** — redirects R7/R19 away from a hand-rolled 2nd Anthropic client | **Yes** — `replace digital.vasic.llmprovider => ./submodules/llm_provider` | Yes (leaf, `deps: []`) |
| 2 | http3 | `http3` | `git@github.com:vasic-digital/http3.git` | **ADOPT** — escalated from STUDY-only, real diff done vs. target's thin `internal/api/http3.go` | **Yes** — `replace digital.vasic.http3 => ./submodules/http3` | Yes (leaf) |
| 3 | Containers | `containers` (root-level exception, §6) | `git@github.com:vasic-digital/containers.git` | **ADOPT** — ROOT path per family precedent (`helix_ota` live-code) | **Yes** — `replace digital.vasic.containers => ./containers` | Yes (leaf) |
| 4 | DagOrchestrator | `dag_orchestrator` | `git@github.com:HelixDevelopment/DagOrchestrator.git` | **INCORPORATE (scoped)** — G11 worker-scheduling machinery only | **Yes** — `replace dev.helix.dag => ./submodules/dag_orchestrator` | Yes (leaf) |
| 5 | PipelineRuntime | `pipeline_runtime` | `git@github.com:HelixDevelopment/PipelineRuntime.git` | **ADOPT (scoped, deferred)** — R6 wizard progress-streaming; not yet built | **Yes, once wired** — `replace dev.helix.pipeline => ./submodules/pipeline_runtime` | Yes (leaf) |
| 6 | VisionEngine | `vision_engine` | `git@github.com:HelixDevelopment/VisionEngine.git` | **N/A-until-R3-clients-exist** | Deferred | Yes, when actionable (leaf) |
| 7 | Panoptic | `panoptic` | `git@github.com:vasic-digital/Panoptic.git` | **N/A-until-R3-clients-exist** | Deferred | Yes, when actionable (leaf; `module panoptic` — unnamespaced) |
| 8 | LLMsVerifier | `llms_verifier` | `git@github.com:vasic-digital/LLMsVerifier.git` | **ADOPT** (already in `project/helix-deps.yaml`) | **Yes** — `replace llmsverifier => ./submodules/llms_verifier` if the server imports its jury/router API directly; **else** consumed as an external scoring service only (see §5.2) | Yes — **1 transitive dep** (`Challenges`, its own `helix-deps.yaml`) |
| 9 | HelixLLM | `helix_llm` | `git@github.com:HelixDevelopment/HelixLLM.git` | **ADOPT** (already in `project/helix-deps.yaml`) | **No** — confirmed HTTP-only (`OpenAI-compatible local model daemon`, port `8443`/`18434`; consumed via `HELIX_LLM_LOCAL_OPENAI_ENDPOINT`, `research/helix_interop_incorporation.md:31-58` — never Go-imported) | No manifest exists upstream (confirmed 404, §8) — record `deps: []` locally with a note |
| 10 | HelixAgent | `helix_agent` | `git@github.com:HelixDevelopment/HelixAgent.git` | **ADOPT** (already in `project/helix-deps.yaml`) | **No** — confirmed HTTP/gRPC-only (REST `:7061`, gRPC `LLMFacade`, MCP-SSE bridge; the project hand-writes its own thin HTTP client, `research/helix_interop_incorporation.md:97-135` — never Go-imported) | Yes, for documentation only — its own manifest lists a heavy transitive set (`VectorDB`, `VisionEngine`, `conversation`, more) that this project does **not** need to recursively vendor because it consumes HelixAgent as a running service, not a compiled dependency (§5.2) |
| 11 | Embeddings | `embeddings` | `git@github.com:vasic-digital/Embeddings.git` | **ADOPT** (already in `project/helix-deps.yaml`) | **Yes** — `module digital.vasic.embeddings` confirmed via `go.mod` (leaf, only `stretchr/testify`); `replace digital.vasic.embeddings => ./submodules/embeddings` | Yes (leaf, `deps: []`, confirmed) |
| 12 | HelixQA | `helix_qa` | `git@github.com:HelixDevelopment/HelixQA.git` | **ADOPT** (already in `project/helix-deps.yaml`) | **No, by design** — consumed as an external test-bank runner (`tools/helixqa_runner`-class separate binary/`go.mod`), never compiled into the main server, mirroring `helix_ota`'s own precedent (its `server/go.mod` does not require `digital.vasic.helixqa`; only a side tool does) | Yes, for documentation — its own manifest declares **8** own-org transitive deps (`Challenges`, `Containers`, `DocProcessor`, `LLMOrchestrator`, `LLMProvider`, `LLMsVerifier`, `security`, `VisionEngine`); **not all 8** need vendoring here unless HelixQA is later compiled in-tree (§5.2) |
| 13 | Challenges | `challenges` | `git@github.com:vasic-digital/Challenges.git` | **ADOPT** (already in `project/helix-deps.yaml`) | **Only inside a dedicated challenge-runner `go.mod`** (own module, not the main server's), matching the "no fakes beyond unit tests" split (§11.4.27) | Yes (leaf, `module digital.vasic.challenges`) |
| 14 | open-design | `open_design` | `git@github.com:nexu-io/open-design.git` | **ADOPT** — already fully designed (`research/opendesign_incorporation.md`) | **No** — not a Go module, a design/brand-generation corpus | N/A (third-party, no manifest shipped) |
| 15 | design_system | `design_system` | `git@github.com:vasic-digital/design_system.git` | **ADOPT** — alongside open-design | **No** — npm/Tailwind-v4 CSS-token package, no `go.mod` | **CONFIRMED-ABSENT upstream** (`gh api .../contents/helix-deps.yaml` → 404; `gh api .../contents` root listing has no `helix-deps.yaml` among `[.gitignore, LICENSE, README.md, assets, components, docs, fonts, i18n, manifest.json, package.json, scripts, tailwind, tokens, upstreams]`) — treat as leaf, npm-shaped, record `deps: []` locally with the confirmed-absent note (§11.4.31 gap, not a blocker) |

### 2.B — Consumed from `constitution/` (NOT vendored under this project's own `submodules/`; orthogonal to G14/X1, §1)

| # | Name | Real location | Mechanism | Notes |
|---|---|---|---|---|
| 16 | token_optimizer | `constitution/submodules/token_optimizer/` | §11.4.28(C) depth-1 carve-out, consumed by reference from the constitution submodule this project already vendors | **WIRE**, not vendor — this project supplies a `wire.Encoder` TOON adapter (§4.4 of `r22_full_catalogue_incorporation_design.md`), no new submodule |
| 17 | docs_chain | Pending — **recommended** (`g43_docs_chain_export_wiring_design.md`) to land at `constitution/submodules/docs_chain/`, same carve-out as #16 | Same mechanism as #16, **once the constitution submodule maintainers land it there** | Landing `constitution/submodules/docs_chain/` is a **constitution-submodule-level decision**, outside this project's scope to perform unilaterally. **Interim fallback**, if the project needs Docs Chain wiring before that lands: vendor fresh at this project's own `submodules/docs_chain/` (`git@github.com:vasic-digital/docs_chain.git`, module confirmed pure-Go/`fsnotify`/`yaml.v3`, leaf) exactly like any Category-A dependency, then **retire** the interim copy the moment `constitution/submodules/docs_chain/` exists (a tracked, evidence-backed migration per §11.4.197 — never left as a permanent duplicate) |
| 18 | Workable-items engine | `constitution/scripts/workable-items/` | Inherited automatically as part of `constitution/`'s own `scripts/` tree — not a submodule at all | No action needed beyond the standing constitution inheritance every project already has (§11.4.35); `g40_workable_items_db_adoption_design.md` confirms it builds clean and is unused-but-adoptable today |

### 2.C — STUDY-only / deferred / NOT vendored at this time (verdict unchanged, listed for completeness)

| # | Name | Org / upstream | Verdict | Action |
|---|---|---|---|---|
| 19 | DebateOrchestrator | `git@github.com:HelixDevelopment/DebateOrchestrator.git` | **STUDY-only** — the project's own `LLMJury`/`computeOverallVerdict` already implements the correct two-factor consensus; adopting `ConductDebate`'s single blended-float threshold would be a regression | None — no submodule added |
| 20 | Auth | `git@github.com:vasic-digital/Auth.git` | **STUDY-pattern-only** — promote only when a bearer-token/API-key auth requirement is stated (not yet) | None |
| 21 | MCP_Module | `git@github.com:vasic-digital/MCP_Module.git` | **STUDY-only** — target's own `mark3labs/mcp-go`-based MCP layer is already hardened; not redirected | None |
| 22 | KMP suite (`UI-Components-KMP`, `Security-KMP`, `Concurrency-KMP`, `Storage-KMP`) | `git@github.com:vasic-digital/{UI-Components-KMP,Security-KMP,Concurrency-KMP,Storage-KMP}.git` (all confirmed real via `gh repo view`) | **ADOPT-consideration**, deferred weigh-off vs. Flutter to the client-build phase (R3/G15); does not resolve HarmonyOS/Aurora either way | None until R3 client-architecture decision |
| 23 | `config`/`security`/`Storage`/`Filesystem` generic modules | Present under `helix_code/submodules/`, not independently diffed | **FLAGGED**, no verdict given | None — honest gap, not a blocker |

**UNCERTAIN upstreams: none.** Every `ssh_url` in §2.A and §2.C was independently
re-confirmed this pass via live `gh repo view`/`gh api` calls (2026-07-15) — no fabricated
URL appears anywhere in this table. The only genuinely open items are (a) whether
`llms_verifier`'s jury API is directly Go-imported or consumed as an external scoring
service (#8 — a design decision for whoever lands the P3/R7 wiring, not an upstream
uncertainty), and (b) the `docs_chain` landing-location timing (#17 — a scheduling
dependency on the constitution-submodule maintainers, not an upstream uncertainty).

---

## 3. Reference layout from a real `helix_*` project (cited, read directly)

`helix_ota` is the cleanest live precedent for the exact vendor-fresh + `go.mod replace`
pattern this design adopts, confirmed by reading its files directly (not inferred):

**`helix_ota/.gitmodules`** (excerpt, verbatim):
```
[submodule "containers"]
	path = containers
	url = git@github.com:vasic-digital/containers.git
[submodule "submodules/http3"]
	path = submodules/http3
	url = git@github.com:vasic-digital/http3.git
[submodule "submodules/llm_provider"]
	path = submodules/llm_provider
	url = git@github.com:vasic-digital/LLMProvider.git
[submodule "submodules/vision_engine"]
	path = submodules/vision_engine
	url = git@github.com:vasic-digital/VisionEngine.git
```
(Full file also declares `submodules/ota-*`, `submodules/helixqa`, `submodules/challenges`,
`docs_chain`, `submodules/doc_processor`, `submodules/llm_orchestrator`,
`submodules/security`, `submodules/llms_verifier`, `submodules/website` — 19 submodules
total, every one **owned and vendored independently by `helix_ota` itself**, none consumed
by pointing at `helix_code` or any other sibling.)

**`helix_ota/server/go.mod`** (the Go `replace` wiring, verbatim):
```go
require (
	digital.vasic.containers v0.0.0-00010101000000-000000000000
	digital.vasic.http3      v0.0.0-00010101000000-000000000000
	...
)

replace digital.vasic.containers => ../containers
replace digital.vasic.http3 => ../submodules/http3

// ota-protocol is co-developed (§11.4.28); built against the local submodule
// during development. Production pins a tagged version.
replace github.com/HelixDevelopment/ota-protocol => ../submodules/ota-protocol
```
This is the exact **"vendor fresh + `replace` => local checkout"** mechanism this design
generalizes for the skill-system project (§5). Note the **one legacy wrinkle in
`helix_ota` itself** this design deliberately avoids: `helix_ota/.gitmodules` declares
`containers` at **both** `containers/` (root) **and** `submodules/containers/` (identical
`ssh_url`) — a duplicate-copy artifact from an earlier layout migration, not a pattern to
copy. This design's §6 mandates exactly **one** canonical location per dependency, never
two.

**`helix_ota/helix-deps.yaml`** confirms the **recursive resolution** contract this design
adopts for the skill-system project (§7): every dependency is declared with `name`,
`ssh_url`, `ref`, `describe`, `layout: flat|grouped`, `why`; `transitive_handling.recursive:
true`; own-org-only scope (third-party deps declared separately, honestly, never
incorporated as submodules).

**`helix_ota/upstreams/`** confirms the §11.4.36 `install_upstreams` wiring (§5.1):
top-level `upstreams/{GitHub,GitLab,GitFlic,GitVerse}.sh` (one `UPSTREAMABLE_REPOSITORY=`
export each) plus a per-submodule `upstreams/submodules/<name>/<remote>.sh` tree for
submodules that themselves declare additional push mirrors (e.g.
`upstreams/submodules/website/vasic-digital.sh`).

---

## 4. Canonical `submodules/` tree (target state, post-incorporation)

```
<skill-system-root>/                      # today: project/, physically inside helix_skills
├── constitution/                         # universal governance submodule (§1, unchanged)
│   └── submodules/
│       ├── token_optimizer/              # §11.4.28(C) carve-out — consumed, not vendored here
│       ├── session_orchestrator/
│       ├── continuum/
│       └── docs_chain/                   # PENDING — recommended landing site (§2.B #17)
├── containers/                           # ROOT-level exception (§6) — vasic-digital/containers
├── go.mod                                # module github.com/helixdevelopment/skill-system
├── helix-deps.yaml                       # this project's OWN manifest (§7) — extended
├── .helix-manifest.yaml                  # recursive-resolution audit record (§7.3), generated
├── upstreams/                            # §11.4.36 — this project's OWN push-mirror recipes
│   └── GitHub.sh                         # (+ GitLab.sh / GitFlic.sh / GitVerse.sh as configured)
└── submodules/
    ├── llm_provider/                     # HelixDevelopment/LLMProvider
    ├── http3/                            # vasic-digital/http3
    ├── dag_orchestrator/                 # HelixDevelopment/DagOrchestrator
    ├── pipeline_runtime/                 # HelixDevelopment/PipelineRuntime
    ├── vision_engine/                    # HelixDevelopment/VisionEngine (N/A until R3)
    ├── panoptic/                         # vasic-digital/Panoptic (N/A until R3)
    ├── llms_verifier/                    # vasic-digital/LLMsVerifier
    ├── helix_llm/                        # HelixDevelopment/HelixLLM (service, no go.mod replace)
    ├── helix_agent/                      # HelixDevelopment/HelixAgent (service, no go.mod replace)
    ├── embeddings/                       # vasic-digital/Embeddings
    ├── helix_qa/                         # HelixDevelopment/HelixQA (external test-bank runner)
    ├── challenges/                       # vasic-digital/Challenges (own go.mod, not main server)
    ├── open_design/                      # nexu-io/open-design — snake_case per operator (§8)
    └── design_system/                    # vasic-digital/design_system
```

Everything under `submodules/` is a **single, flat, non-nested** list — no dependency's own
`.gitmodules` is ever re-walked into a further own-org submodule chain (§11.4.28(C)'s
"nested own-org submodule chains are FORBIDDEN" is verified for every Category-A dependency
in §2.A: `dag_orchestrator`/`pipeline_runtime`/`vision_engine`/`panoptic`/`llm_provider`/
`http3` are all confirmed **leaves** — `deps: []` in their own `helix-deps.yaml` — and
`llms_verifier`'s one transitive dep (`Challenges`) is itself already declared as a
Category-A sibling at the same flat depth, not nested inside `submodules/llms_verifier/`).

---

## 5. `.gitmodules` + `go.mod` replace plan

### 5.1 `.gitmodules` (target state — additive; nothing here exists yet, no file touched by this document)

```gitconfig
[submodule "constitution"]
	path = constitution
	url = git@github.com:HelixDevelopment/HelixConstitution.git

[submodule "containers"]
	path = containers
	url = git@github.com:vasic-digital/containers.git

[submodule "submodules/llm_provider"]
	path = submodules/llm_provider
	url = git@github.com:HelixDevelopment/LLMProvider.git

[submodule "submodules/http3"]
	path = submodules/http3
	url = git@github.com:vasic-digital/http3.git

[submodule "submodules/dag_orchestrator"]
	path = submodules/dag_orchestrator
	url = git@github.com:HelixDevelopment/DagOrchestrator.git

[submodule "submodules/pipeline_runtime"]
	path = submodules/pipeline_runtime
	url = git@github.com:HelixDevelopment/PipelineRuntime.git

[submodule "submodules/vision_engine"]
	path = submodules/vision_engine
	url = git@github.com:HelixDevelopment/VisionEngine.git

[submodule "submodules/panoptic"]
	path = submodules/panoptic
	url = git@github.com:vasic-digital/Panoptic.git

[submodule "submodules/llms_verifier"]
	path = submodules/llms_verifier
	url = git@github.com:vasic-digital/LLMsVerifier.git

[submodule "submodules/helix_llm"]
	path = submodules/helix_llm
	url = git@github.com:HelixDevelopment/HelixLLM.git

[submodule "submodules/helix_agent"]
	path = submodules/helix_agent
	url = git@github.com:HelixDevelopment/HelixAgent.git

[submodule "submodules/embeddings"]
	path = submodules/embeddings
	url = git@github.com:vasic-digital/Embeddings.git

[submodule "submodules/helix_qa"]
	path = submodules/helix_qa
	url = git@github.com:HelixDevelopment/HelixQA.git

[submodule "submodules/challenges"]
	path = submodules/challenges
	url = git@github.com:vasic-digital/Challenges.git

[submodule "submodules/open_design"]
	path = submodules/open_design
	url = git@github.com:nexu-io/open-design.git

[submodule "submodules/design_system"]
	path = submodules/design_system
	url = git@github.com:vasic-digital/design_system.git
```

`vision_engine` and `panoptic` entries are added to `.gitmodules` at the same time as the
rest (cheap, zero-risk `git submodule add`) but their **incorporation** (§9) is deferred
until R3 clients exist — vendoring the submodule pointer early costs nothing and avoids a
second churn-y `.gitmodules` edit later; this mirrors how `helix_ota` vendors submodules it
does not yet fully wire (e.g. `submodules/website`).

### 5.2 `go.mod` replace plan (main server module, `github.com/helixdevelopment/skill-system`)

Only dependencies that are **genuinely Go-imported by the main server binary** get a
`replace` directive in the top-level `go.mod` — mirroring `helix_ota/server/go.mod`'s own
selective pattern (it vendors 19 submodules but `replace`s only 3 in the server module;
others are consumed by separate tool binaries with their own `go.mod`, §3).

```go
module github.com/helixdevelopment/skill-system

require (
	digital.vasic.containers     v0.0.0-00010101000000-000000000000
	digital.vasic.http3          v0.0.0-00010101000000-000000000000
	digital.vasic.llmprovider    v0.0.0-00010101000000-000000000000
	dev.helix.dag                v0.0.0-00010101000000-000000000000
	digital.vasic.embeddings     v0.0.0-00010101000000-000000000000
	// dev.helix.pipeline added when R6 progress-streaming is actually built (§9)
)

replace digital.vasic.containers  => ./containers
replace digital.vasic.http3       => ./submodules/http3
replace digital.vasic.llmprovider => ./submodules/llm_provider
replace dev.helix.dag             => ./submodules/dag_orchestrator
replace digital.vasic.embeddings  => ./submodules/embeddings
// replace dev.helix.pipeline     => ./submodules/pipeline_runtime   (once wired, §9)
```

**NOT `replace`d in the main server `go.mod`** (confirmed reasons, §2.A):

- `llms_verifier` — pending the §9.5 design decision on direct-import vs. external-service
  consumption; if direct-import is chosen, add `replace llmsverifier =>
  ./submodules/llms_verifier`.
- `helix_llm`, `helix_agent` — confirmed HTTP-only services (`research/helix_interop_incorporation.md`);
  the project hand-writes a thin `HelixLLMProvider`-shaped HTTP client, never a Go import of
  either module. Vendoring the submodule still matters (running the service locally via
  `docker-compose`, referencing its README/`.env.example` defaults, `install_upstreams`
  parity) — it just carries no `replace` line.
- `helix_qa`, `challenges` — consumed via a **separate** tool binary/`go.mod` (e.g.
  `tools/helixqa_runner/go.mod`, `tools/challenges_runner/go.mod`, mirroring
  `helix_ota/tools/helixqa_runner/go.mod` + `helix_ota/tools/device_claim/go.mod`'s own
  precedent of per-tool `go.mod` files distinct from `server/go.mod`) — **never** the main
  server module, per §11.4.27's "mocks/stubs/fakes only in unit tests; every other test
  type interacts with the real system" (Challenges/HelixQA power acceptance-test banks, they
  are not production dependencies of the shipped server).
- `vision_engine`, `panoptic` — deferred (§9), no `replace` until R3 clients exist.
- `open_design`, `design_system` — not Go modules; no `go.mod` involvement at all (their
  consumption is a Web-client build-time asset pipeline, R12).

**Per-tool `go.mod` replace plan** (illustrative, exact tool names TBD at implementation
time):

```go
// tools/helixqa_runner/go.mod
module github.com/helixdevelopment/skill-system/tools/helixqa_runner
require digital.vasic.helixqa v0.0.0-00010101000000-000000000000
replace digital.vasic.helixqa => ../../submodules/helix_qa

// tools/challenges_runner/go.mod
module github.com/helixdevelopment/skill-system/tools/challenges_runner
require digital.vasic.challenges v0.0.0-00010101000000-000000000000
replace digital.vasic.challenges => ../../submodules/challenges
```

If (and only if) `helix_qa` is ever compiled in-tree by one of these tool modules, its own
**8** transitive own-org deps (`Challenges`, `Containers`, `DocProcessor`, `LLMOrchestrator`,
`LLMProvider`, `LLMsVerifier`, `security`, `VisionEngine` — per its own `helix-deps.yaml`,
confirmed §2.A #12) would each need their own flat-sibling vendor entry too, per the
recursive-resolution contract (§7). Today, this project only needs 4 of those 8 for other
reasons already (`Challenges`, `LLMProvider`, `LLMsVerifier`, `VisionEngine`) — `Containers`,
`DocProcessor`, `LLMOrchestrator`, `security` would be **new** additions triggered solely by
compiling HelixQA in-tree. Following `helix_ota`'s own documented precedent (§2.A #12,
"HelixQA does not fully compile in-tree here; the build is unaffected — it consumes the bank
via `tools/helixqa/run_bank.sh`, never by compiling HelixQA"), this design recommends
**not** compiling HelixQA in-tree, avoiding that 4-dependency expansion entirely.

---

## 6. Root-level exception: `containers/`

`containers` is the **one** documented exception to the flat `submodules/<name>/` rule,
justified by live family-wide precedent (`helix_ota` mounts it at bare `containers/`,
`research/helix_family_reusable_practices.md` §10 row 2: "ADOPT-as-submodule, ROOT path
`containers/` (per `helix_ota`'s live-code precedent)"). §11.4.28(C) explicitly permits
"ungrouped" root-level mounting (`<repo_root>/<name>/`) as an equally-canonical alternative
to "grouped" (`<repo_root>/submodules/<name>/`) — this is a **documented layout choice**,
not a violation of the single-canonical rule, so long as exactly **one** of the two forms is
used per dependency (§4's tree uses root-level `containers/` only — it is never also
duplicated at `submodules/containers/`, unlike `helix_ota`'s own accidental double-mount,
flagged in §3 as a wrinkle this design avoids). `helix-deps.yaml`'s `layout` field records
this per dependency (`layout: ungrouped` for `containers`, `layout: grouped` for every other
Category-A dependency) — the existing `project/helix-deps.yaml` schema already supports
this field (`layout: grouped` on every current entry), so no schema change is needed, only a
new entry with `layout: ungrouped`.

---

## 7. `install_upstreams` + `helix-deps` recursive resolution

### 7.1 `install_upstreams` wiring (§11.4.36)

Per the constitution's `install_upstreams.sh` contract (read directly,
`constitution/install_upstreams.sh:1-40`): on **every** clone/add of a repository that
carries an `upstreams/` directory populated with `*.sh` recipe files, `install_upstreams`
must be invoked from that repository's root. This applies at **two** levels for the
skill-system project:

1. **This project's own root** — once extracted, the skill-system repo carries its own
   `upstreams/GitHub.sh` (`export UPSTREAMABLE_REPOSITORY="git@github.com:<org>/skill-system.git"`),
   mirroring `helix_ota/upstreams/GitHub.sh`'s exact one-line shape. `install_upstreams` runs
   once at repo setup, configuring every declared push mirror as a local git remote so a
   single push fans out per §2.1.
2. **Any vendored submodule that itself ships an `upstreams/` directory** — confirmed for
   `design_system` (§2.A #15, "carries an `upstreams/` directory" per
   `r22_full_catalogue_incorporation_design.md` §7.1's `gh api contents` listing). Per
   §11.4.36, `install_upstreams` must be run **from inside** `submodules/design_system/`
   immediately after that submodule is added, exactly like `helix_ota`'s own
   `upstreams/submodules/website/vasic-digital.sh` pattern (a per-submodule push-mirror
   recipe nested under the **parent** project's `upstreams/submodules/<name>/`, not inside
   the submodule's own tree — the parent records which of its vendored submodules need the
   invocation and where their own recipes live). No other Category-A dependency in §2.A is
   confirmed to ship its own `upstreams/` directory (their `helix-deps.yaml`/`go.mod`
   listings show no such directory) — this is re-checked at actual incorporation time per
   the standing §11.4.36 mandate (a newly-discovered `upstreams/` dir in any dependency
   triggers the same invocation, whether or not this document names it today).

### 7.2 Recursive `helix-deps.yaml` resolution

The existing `project/helix-deps.yaml` already implements the correct **schema**
(`schema_version: 1`, `name`/`ssh_url`/`ref`/`why`/`layout` per entry) — this design
**extends** it (adds entries for §2.A rows 1–7 and 14–15 which are not yet declared) rather
than replacing it. The **recursion rule** (§11.4.31, confirmed by `helix_ota/helix-deps.yaml`'s
own `transitive_handling.recursive: true` + `conflict_resolution: operator-required`):
for every Category-A dependency, its **own** `helix-deps.yaml` is read, and any own-org
transitive dependency it declares is added as a **sibling** entry in this project's own
manifest (never nested under the dependency's tree, §4) — **unless** that transitive
dependency is not actually needed by this project's own consumption pattern (§5.2's
HelixQA/HelixAgent carve-out: a dependency consumed as a running **service**, or as an
external **test-bank/tool** rather than a compiled library, does not obligate recursively
vendoring its own transitive Go-import graph). This project's own recursion audit (per
dependency, from data already gathered in §2.A/§5.2):

| Dependency | Own transitive own-org deps | Recursively vendored here? |
|---|---|---|
| `llm_provider` | none (leaf) | N/A |
| `http3` | none (leaf) | N/A |
| `dag_orchestrator` | none (leaf) | N/A |
| `pipeline_runtime` | none (leaf) | N/A |
| `vision_engine` | none (leaf) | N/A |
| `panoptic` | none (leaf) | N/A |
| `embeddings` | none (leaf) | N/A |
| `llms_verifier` | `Challenges` | **Yes** — already a Category-A sibling (#13) |
| `helix_agent` | `VectorDB`, `VisionEngine`, `conversation`, + more (heavy) | **No** — consumed as an HTTP/gRPC service only (§5.2); not recursively pulled |
| `helix_qa` | 8 deps (`Challenges`, `Containers`, `DocProcessor`, `LLMOrchestrator`, `LLMProvider`, `LLMsVerifier`, `security`, `VisionEngine`) | **Partial** — 4 already needed for other reasons (`Challenges`, `LLMProvider`, `LLMsVerifier`, `VisionEngine`); the remaining 4 (`Containers` overlaps #3 already; `DocProcessor`, `LLMOrchestrator`, `security` are **not** vendored — HelixQA is consumed via its external bank-runner, never compiled in-tree, §5.2) |
| `challenges` | none confirmed beyond its own leaf declaration | N/A |
| `helix_llm` | no manifest exists (404, confirmed) | N/A — service-only consumption regardless |
| `design_system` | confirmed-absent manifest, npm-shaped | N/A |
| `open_design` | third-party, no manifest | N/A |

### 7.3 `.helix-manifest.yaml` audit record

Following `helix_ota`'s own precedent (`.helix-manifest.yaml`, read directly, §1's format:
`schema_version`, `consumer`, `added_for`, `deps: [{name, path, repo}]`, `build_proof`), this
project's own `.helix-manifest.yaml` is generated (not hand-written) by the incorporation
tooling (§9) the first time each dependency batch is actually vendored, recording exactly
which path each landed at and a build-proof pointer (`go build`/`go vet` exit codes +
evidence path under `docs/qa/<run-id>/`) — never claimed complete until that build proof
exists (§11.4.38, §11.4.108).

### 7.4 `sync_submodules.sh` redesign implication (design note, not a code change)

`project/scripts/sync_submodules.sh` (read directly, §1) currently implements Option A's
**ecosystem-root search-then-fallback** logic (`--ecosystem-root` flags /
`HELIX_ECOSYSTEM_ROOTS` env var). Under the reversed vendor-fresh decision, this script's
ecosystem-root search branch becomes **dead code for this project's own resolution path** —
every dependency now **always** vendors fresh at its own canonical local path (§4),
regardless of whether a parent-ecosystem copy exists anywhere on disk. A follow-on
implementation pass (tracked, not performed by this design document per its own scope
contract) should either (a) delete the ecosystem-root search branch entirely, simplifying
the script to "always vendor at the manifest's declared local path," or (b) keep the search
as an **informational** cross-check only (warn if a parent-ecosystem copy exists and looks
newer, but never substitute it for the local vendor) — recommendation (a), since the whole
point of the operator's decision is independence from sibling-checkout presence/freshness;
keeping dead search logic around is itself a §11.4.124 "investigate before removing"
candidate the moment this lands, not before.

---

## 8. `open_design` + `design_system` naming

The operator's decision text is explicit and final on this point: **"open-design →
`submodules/open_design/`"** — snake_case, per §11.4.29's general naming mandate.

**The tension this resolves (documented, not silently dropped):**
`research/r22_full_catalogue_incorporation_design.md` §7.2 flagged a genuine, previously
unreconciled discrepancy: `research/opendesign_incorporation.md` (Revision 1) independently
recommended the snake_case rename `submodules/open_design/`, while **100% of directly-observed
sibling precedent** (`helix_terminator/.gitmodules:4-6` and `helix_vpn/.gitmodules:63-65`,
both confirmed by directory listing — `helix_terminator/submodules/open-design/` and
`helix_vpn/submodules/open-design/` both exist, hyphenated, matching the upstream repo name
`nexu-io/open-design` verbatim) mounts it **unrenamed**, hyphenated. The operator's verbatim
G14/X1 decision resolves this discrepancy **in favor of the snake_case rename** — the
skill-system project's copy is the outlier relative to 2 of 2 observed siblings, and that is
an accepted, deliberate consequence of the operator's explicit instruction, not an oversight.
This is exactly the "vendor-exception reading" of §11.4.29 `opendesign_incorporation.md` §3.1
originally proposed — the operator's decision confirms that reading as the one to implement.

**`design_system`** carries no such discrepancy — it is a `vasic-digital`-owned repo whose
own name is already snake_case (`design_system`), so `submodules/design_system/` requires no
rename decision at all; every observed convention (§11.4.29's own "every sibling uses
`submodules/<snake_case>/`" finding, `helix_family_reusable_practices.md` §9) already agrees.

**Import path implication (Web/R12 client, not yet built).** Neither module affects Go
`import`/`replace` paths (§5.2 — both are non-Go). The only path implication is at the
**client build tool** layer (whatever bundler/asset pipeline the eventual R12 Web client
uses) resolving `submodules/open_design/` and `submodules/design_system/` as its design-token
source directories — a detail for the client-build phase (§9's step 5), not this document.

---

## 9. Incorporation order (deps-before-dependents)

Ordered so no step ever references a not-yet-vendored path, and so the highest-value,
lowest-risk (leaf, zero transitive deps) dependencies land first, per §11.4.132's
risk-ordering discipline applied to a vendoring batch:

1. **`.gitmodules` scaffold (all 15 Category-A entries + `containers` root exception, §5.1)**
   — a single batch `git submodule add` pass; cheap, reversible, no build-graph impact yet.
   `vision_engine`/`panoptic` pointers land here too (deferred wiring, §5.1's rationale).
2. **Leaf, zero-transitive-dep Go modules first** (`llm_provider`, `http3`, `dag_orchestrator`,
   `pipeline_runtime`, `embeddings`) — add `go.mod replace` + `require`, `go build ./...`,
   `go vet ./...`, confirm green before touching anything with a transitive dependency.
   `pipeline_runtime`'s `replace` stays commented until the R6 capability is actually built
   (§2.A #5) — vendoring the submodule now costs nothing; wiring it into `go.mod` before
   there is code to use it would be dead weight (§11.4.124's own "unwired code" caution
   applies to the reverse direction too — don't wire ahead of need).
3. **Root-level `containers/`** — the one layout exception (§6); wire its `replace` alongside
   step 2 (it is itself a leaf, `digital.vasic.containers`, `deps: []`).
4. **`llms_verifier` (1 transitive dep, already satisfied by step 1's `challenges` pointer)**
   — decide direct-import vs. external-service consumption (§9.5 open item) before adding a
   `go.mod replace`; either way its submodule pointer and `helix-deps.yaml` entry land now.
5. **Service-only dependencies (`helix_llm`, `helix_agent`)** — vendor the submodule (for
   local `docker-compose` runnability + README/env-default reference), write the thin
   `HelixLLMProvider`-shaped HTTP client by hand (no Go import, no `replace`), confirm the
   client's fallback chain (`HelixLLM → LLMsVerifier → HelixAgent`, per
   `research/helix_interop_incorporation.md` §2) against a live local instance before
   claiming this step done (§11.4.108 runtime-signature, not a compile-only claim).
6. **External test-bank/tool dependencies (`helix_qa`, `challenges`)** — vendor the
   submodule, scaffold the separate `tools/helixqa_runner/go.mod` +
   `tools/challenges_runner/go.mod` (§5.2), confirm each tool binary builds independently of
   the main server module (proves the "not compiled into the main server" boundary holds).
7. **Design assets (`open_design`, `design_system`)** — vendor both submodules,
   run `install_upstreams` inside `submodules/design_system/` (§7.1's confirmed
   `upstreams/` directory), leave their client-build-tool wiring for the R12/R3 client phase
   (not yet built — no premature coupling).
8. **Deferred-until-actionable (`vision_engine`, `panoptic`)** — submodule pointer already
   landed in step 1; no further action (no `go.mod` entry, no client wiring) until an R3
   client genuinely exists, per the unchanged N/A-until-R3-clients-exist verdict.
9. **`.helix-manifest.yaml` regeneration + `docs/qa/<run-id>/` build-proof capture** (§7.3) —
   the closing step of the whole incorporation batch, never skipped, never claimed complete
   without it (§11.4.38, §11.4.108).

**Open item carried forward, not resolved by this document (§9.5):** whether `llms_verifier`
is directly Go-imported (its own jury/router API called in-process) or consumed purely as an
external scoring service alongside `helix_llm`/`helix_agent` in the same fallback chain — the
interop design (`research/helix_interop_incorporation.md` §5, "LLMsVerifier is consumed
in-process by HelixLLM/HelixAgent today") suggests the **existing** Helix-family pattern is
in-process-within-those-services, which would argue for this project **also** treating it as
an external service (no `go.mod replace`) rather than a direct import — but this document
does not fabricate that conclusion as settled FACT; it is flagged here as the one design
decision the incorporation implementer must make explicitly before step 4 above, citing
whichever evidence resolves it at that time.

---

## 10. Risks / open items

1. **No UNCERTAIN upstream SSH URLs** — every URL in §2.A/§2.C was independently
   re-confirmed live via `gh repo view`/`gh api` this pass (2026-07-15); none are fabricated
   or carried over unverified from an older document.
2. **`docs_chain` landing-location timing (§2.B #17)** is a genuine open scheduling risk —
   if the constitution-submodule maintainers do not land `constitution/submodules/docs_chain/`
   before this project needs Docs Chain wiring, the interim vendor-fresh fallback (still
   fully compliant with the operator's G14/X1 decision, since it is *this project's own*
   copy) must be tracked as a migration item (§11.4.197) so it does not silently persist as
   a permanent duplicate once the canonical constitution-submodule copy exists.
3. **`llms_verifier` direct-import-vs-external-service decision (§9, item 9.5)** is
   explicitly unresolved by this design — it changes whether a `go.mod replace` line exists
   for it, but does not change its vendor-fresh submodule status (it is vendored either way).
4. **`helix_ota`'s own duplicate `containers` mount** (§3, §6) is flagged as a wrinkle this
   design's tree (§4) deliberately does not reproduce — a live reminder that "vendor fresh"
   still requires the single-canonical-location discipline (§11.4.28(C)) to be enforced
   *within* this project too, not just *between* this project and its siblings.
5. **`sync_submodules.sh`'s ecosystem-root search branch becoming dead code (§7.4)** is a
   real follow-on cleanup this document flags but does not perform — leaving it un-simplified
   past the point this design lands would itself become a §11.4.124-class "investigate before
   removing" candidate, tracked here so it is not lost.
6. **`design_system`'s confirmed-absent `helix-deps.yaml` (§2.A #15)** is a real upstream
   gap against §11.4.31, not a blocker for this project (§11.4.31's own guidance: absence of
   a Go-shaped manifest on an npm-shaped package is different in kind from a hidden own-org
   coupling) — optionally reported upstream per §11.4.74's extend-don't-reimplement path, not
   this project's obligation to fix.
7. **HarmonyOS/Aurora client-platform gap (KMP suite, §2.C #22)** remains unresolved by any
   dependency named in this document — carried forward unchanged from `g15_aurora_harmonyos_client_feasibility.md`,
   out of scope for a submodule-layout design.

**No BLOCKER identified.** Every dependency this project currently needs (§2.A, already
R22/G05-G11-accepted verdicts) has a confirmed real upstream and a concrete vendor-fresh
mount point; the two items with residual ambiguity (§9 item 9.5, §2.B #17 timing) are design
decisions/scheduling risks for the implementer, not blockers to landing this layout plan.
