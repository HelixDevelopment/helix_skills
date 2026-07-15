# R18 Documentation-Delivery Architecture — HelixKnowledge Skill Graph System

**Revision:** 2
**Last modified:** 2026-07-15T16:17:17Z
**Status:** research / design. DESIGNED, not yet delivered — every producing
component is gated (Docs Chain vendoring on G14/X1; OpenDesign vendoring on a
§11.4.66 operator sign-off; the **A3 interactive-API-docs serving feature** — the
`/docs` handler + vendored Scalar/Redoc assets + `/openapi.{json,toon}` endpoints —
is not yet implemented, a P13.T4 deliverable). NB: the Go backend itself **compiles
clean** (`go build ./...`=0 at the committed baseline, P0 build-fix commit `5532e2b`);
the gap is the serving feature, not a compile failure. Nothing below is claimed shipped.
**Authority:** operator mandate **R18** (REQUIREMENTS.md:206-215) — composes
**R10** (Docs Chain, REQUIREMENTS.md:80) + **R12** (OpenDesign, REQUIREMENTS.md:85-89)
+ Constitution §11.4.12/.44/.45/.53/.56/.57/.59/.60/.65/.86/.106/.162/.168/.170/.190.
**Scope:** the FULL documentation surface for this system; the mandatory export
matrix; the always-in-sync wiring through Docs Chain + git hooks; OpenDesign
diagram authoring + host-rendered pixel proof; §11.4.168 anti-bluff render
validation; the phased, tracked plan (§11.4.197) with honest blockers.
**Anti-bluff boundary (§11.4.6):** every capability is grounded in a real
constitution anchor, a real file read this session, or a host tool verified
present this session. Anything unverified is marked UNCONFIRMED / gated. This
document is DESIGN — designed ≠ working.

---

## 0. Honest verdict (read first)

| Question | Answer |
|---|---|
| Is the full doc surface enumerated concretely? | **YES** — 42 deliverables across 8 classes (§1). |
| Are the mandatory export transforms present on host? | **`pandoc` YES, `weasyprint` YES** (docs_chain research §4.3); **`mmdc` YES, `scalar` YES, `convert` YES** (verified this session); **layered-PSD tooling `gimp`/`photopea` ABSENT** (§4 gap G1). |
| Is the sync engine (Docs Chain) usable here today? | **NO — gated.** Engine present only in a *sibling* checkout (`helix_code/submodules/docs_chain`), NOT incorporated into `helix_skills`; vendoring is operator-gated on **G14/X1** (parent-priority vs single-canonical, unresolved); shipped binary is wrong-arch → must `go build`. |
| Is OpenDesign usable here today? | **NO — gated.** Not vendored; a §11.4.66 operator sign-off is owed to add `submodules/open_design/` (third-party). |
| Can the interactive API docs be served today? | **NO — feature not yet implemented.** The Go backend compiles clean (`go build ./...`=0 at the committed baseline); what is missing is the A3 serving feature itself (the `/docs` handler + vendored Scalar/Redoc assets + `/openapi.{json,toon}` endpoints), a P13.T4 deliverable — not a build failure. (The REQUIREMENTS.md:111-120 "FAILS" line is the *pre-P0-build-fix* historical baseline, superseded by commit `5532e2b`.) |
| What IS provable today? | The surface, the matrix, the context designs, and host-tool presence. Nothing is claimed rendered/synced until the gates clear and evidence is captured. |

**Bottom line:** the architecture is complete and grounded; the four things
between DESIGN and WORKING are (a) resolve G14/X1 and vendor Docs Chain
(operator-gated), (b) operator-sign-off + vendor OpenDesign, (c) implement the
A3 interactive-docs serving feature (P13.T4) — the Go backend already compiles,
this adds the `/docs` handler + assets + endpoints, (d) capture the render + sync
evidence (§5). Each is a tracked P-task (§6) so no requirement is left un-wired
(§11.4.197).

---

## 1. The full documentation surface (concrete enumeration)

**42 concrete deliverables across 8 classes** (11 of them OpenDesign diagrams;
the interactive + static API docs cover all **21** OpenAPI paths in
`api/openapi.yaml`). Paths are relative to the future service repo root (the
"consuming repo root", to be declared per §11.4.35 — see §6 blocker B8). The
existing `project/docs/*` set (ARCHITECTURE, API, INSTALL, MCP_INTEGRATION,
DANGER_ZONES — docs_chain research §3.1) is folded in and expanded.

### Class A — API docs (3)
| ID | Deliverable | Kind |
|---|---|---|
| A1 | `api/openapi.yaml` — the OpenAPI 3.1 spec, 21 paths, TOON-primary + JSON-fallback content negotiation | **source** (single source of truth for A2/A3) |
| A2 | **STATIC** API reference site — rendered from A1 to a self-contained HTML page (+ PDF) | derived |
| A3 | **INTERACTIVE** API docs — self-hosted Scalar/Redoc-class explorer **served by the Go server** at `GET /docs` (assets vendored, **no external CDN**), serving `/openapi.json` + `/openapi.toon` derived from A1 | served surface |

### Class B — User manuals (3)
| ID | Deliverable |
|---|---|
| B1 | `docs/manuals/USER_MANUAL.md` — end-user manual: the wizard, skill search/get/tree, learn-from-codebase, git-versioned growth |
| B2 | `docs/manuals/CLIENT_MANUAL.md` — per-surface usage (CLI, TUI, Web, Desktop, Mobile) over the shared REST/MCP core (R3) |
| B3 | `docs/manuals/OPERATOR_MANUAL.md` — install/run/ops runbook: `systemctl --user`, Compose (pgvector), `scripts/{start,stop,restart,status,install,uninstall,logs}` (R15); supersedes/expands `project/docs/INSTALL.md` |

### Class C — Guides (8)
| ID | Deliverable |
|---|---|
| C1 | `docs/guides/ARCHITECTURE.md` — layered architecture + package map (expands `project/docs/ARCHITECTURE.md`) |
| C2 | `docs/guides/MCP_INTEGRATION.md` — MCP + ACP + plugins interop; Claude Code / OpenCode / Kimi wiring (expands `project/docs/MCP_INTEGRATION.md`) |
| C3 | `docs/guides/MODEL_PROVIDER_GUIDE.md` — the pluggable `ModelProvider` chain (HelixLLM → LLMsVerifier → claude-toolkit alias → OpenAI) (R7) |
| C4 | `docs/guides/DEPS_SUBMODULE_GUIDE.md` — `helix-deps.yaml`, `sync_submodules.sh`, parent-priority policy (R9/G14) |
| C5 | `docs/guides/DANGER_ZONES.md` — expands `project/docs/DANGER_ZONES.md` (Aurora/Harmony, TOON codec, embedding drift, jury false-approve) |
| C6 | `docs/guides/FEASIBILITY_METHODOLOGY.md` — the explicit "Can we do this? Here's exactly how" deliverable (REQUIREMENTS.md:45-46) |
| C7 | `docs/guides/CONTRIBUTING.md` — contributor/dev guide (build, test taxonomy, gates) |
| C8 | `docs/guides/SECURITY.md` — CORS/auth model, secret handling (§11.4.10), sandbox isolation posture |

### Class D — Tutorials (5)
| ID | Deliverable |
|---|---|
| D1 | `docs/tutorials/wizard_walkthrough.md` — enter `android, android_aosp, java, kotlin, c++, cmake` → create → map → process (R6) |
| D2 | `docs/tutorials/learn_from_codebase.md` — point at a real repo → extract evidence skills |
| D3 | `docs/tutorials/add_model_provider.md` — implement + register a new `ModelProvider` |
| D4 | `docs/tutorials/agent_integration.md` — consume skills from Claude Code / OpenCode / Kimi over MCP |
| D5 | `docs/tutorials/author_skill_toml.md` — hand-write a skill TOML (SPEC §4.2) + import |

### Class E — FAQ (1)
| ID | Deliverable |
|---|---|
| E1 | `docs/FAQ.md` |

### Class F — Diagrams / schemes / graphs — OpenDesign illustrations (11) (R12/§11.4.162)
| ID | Diagram | Source shape |
|---|---|---|
| DG1 | System architecture (layered: clients → interfaces → domain → model-access → storage) | IMPLEMENTATION_PLAN §3 |
| DG2 | **Skill DAG** — recursive typed-edge dependency graph (atomic ↔ composite/umbrella, R16) | SPEC §4 + g06_g07 design |
| DG3 | ER / DB schema diagram (skills, dependencies, resources, evidences, registry, audit_log) | SPEC §5 SQL |
| DG4 | Wizard sequence — create → map (DAG) → expand closure → process (R6) | REQUIREMENTS R6 |
| DG5 | Auto-growth jury pipeline sequence — draft → resource-verify → sandbox → cross-ref → jury(≥2) → merge (R11) | IMPLEMENTATION_PLAN P4 |
| DG6 | MCP real-time growth sequence — MCP `create_skill` → live pipeline → immediately queryable (R14) | IMPLEMENTATION_PLAN P6 |
| DG7 | Learn-from-codebase flow — repo → tree-sitter → candidate skills → evidence | IMPLEMENTATION_PLAN P5 |
| DG8 | Model-provider fallback chain | IMPLEMENTATION_PLAN P3 |
| DG9 | Deployment topology — `systemctl --user` + Compose (pgvector) + scripts (R15) | IMPLEMENTATION_PLAN P12 |
| DG10 | Docs-sync chain — the Docs Chain wiring itself (this document's §3) | this doc |
| DG11 | Client / shared-core topology — CLI/TUI/Web/Desktop/Mobile over OpenAPI+MCP (R3) | IMPLEMENTATION_PLAN P8 |

### Class G — SQL definitions / schema docs (2)
| ID | Deliverable | Kind |
|---|---|---|
| G1 | `docs/DATABASE_SCHEMA.md` — documents every table/column/index/HNSW/trigger + embeds DG3 | derived doc |
| G2 | `migrations/*.sql` — the authoritative SQL definitions (SPEC §5) | **source SQL** |

### Class H — Templates & other materials (9)
| ID | Deliverable |
|---|---|
| H1 | `config/config.toml` template (SPEC §8) |
| H2 | Skill TOML template (SPEC §4.2) |
| H3 | `.env.example` (§11.4.10/§11.4.77 — secrets never committed, template only) |
| H4 | `systemd/helixknowledge.service` user-unit template (R15) |
| H5 | Compose template (`pgvector/pgvector:pg16`) (R15) |
| H6 | `submodules/open_design/design-systems/helixknowledge/DESIGN.md` + `tokens.css`/`design-tokens.json` — the brand contract (R12/§11.4.162) |
| H7 | `README.md` — top-level, with the §11.4.57 Tracked-Items + Status doc-link section |
| H8 | `CHANGELOG.md` |
| H9 | Packaging materials — reproducible `zip` + `tar.gz` of the deliverable + release script (founding request) |

---

## 2. Export matrix (per §11.4.65; DOCX where the class warrants)

**§11.4.65 floor = `.md` + `.html` + `.pdf`** for every narrative doc. **DOCX**
is added where the class is a formal deliverable / status / workable-item doc
(§11.4.148/.149/.153/.168 — the docs_chain research already produced `.docx` for
the planning-doc context). Diagrams export **SVG + PNG + PDF** (+ PSD, gated).
The interactive API docs (A3) are a **served** surface (validated in §5, not a
file export). "Transform present?" columns are grounded in host checks
(docs_chain research §4.3 + this session's `command -v`).

| Class / item | md | html | pdf | docx | svg | png | psd | Producing transform(s) | Present? |
|---|:--:|:--:|:--:|:--:|:--:|:--:|:--:|---|---|
| B/C/D/E manuals·guides·tutorials·FAQ | ● | ● | ● | ○¹ | – | – | – | `pandoc-html`, `weasyprint-pdf`, (`pandoc-docx`¹) | ✅ pandoc+weasyprint |
| G1 DATABASE_SCHEMA | ● | ● | ● | ○ | – | – | – | same | ✅ |
| A1 openapi.yaml (source) | – | – | – | – | – | – | – | authored; validated by an OpenAPI linter | ✅ (spec.) |
| A2 static API reference | – | ● | ● | – | – | – | – | `exec: scalar`² (spec→HTML) → `weasyprint-pdf` | ✅ `scalar`, `weasyprint` |
| A3 interactive API docs | – | served | – | – | – | – | – | Go server serves vendored Scalar/Redoc + `/openapi.{json,toon}` | ⚠ feature unimplemented (P13.T4); backend compiles |
| F DG1–DG11 diagrams | src | ● | ● | – | ● | ● | ○³ | Mermaid src → `mmdc`(→SVG) → `convert`(→PNG) → `weasyprint`/embed(→PDF); branded via OpenDesign skills | ✅ `mmdc`,`convert`; ⚠ OD gated |
| G2 migrations/*.sql (source) | – | – | – | – | – | – | – | authored; `psql`/CI schema-apply is the proof, not a doc export | ✅ |
| H1–H5,H7,H8 text templates/materials | ● | ●⁴ | ●⁴ | – | – | – | – | `pandoc-html`,`weasyprint-pdf` where a rendered copy is wanted | ✅ |
| H6 OpenDesign brand | src | – | – | – | – | – | – | authored tokens (consumed, not exported) | ⚠ OD gated |
| H9 packaging | – | – | – | – | – | – | – | `zip` + `tar.gz` (checksums) | ✅ |

¹ **DOCX** applies to the formal-deliverable subset (manuals, the feasibility
doc, DATABASE_SCHEMA) per §11.4.153; ordinary tutorials/FAQ take the md+html+pdf
floor unless promoted. ○ = applicable-where-class-warrants.
² Scalar CLI renders a self-contained HTML API reference from the spec; `@redocly/cli`
is a drop-in alternative via `npx` (both self-hostable, no runtime CDN).
³ **PSD** is a §11.4.6-gated gap: OpenDesign has **no PSD path** (opendesign
research G1); flat PSD via `convert` (present) is possible, **layered** PSD needs
`gimp`/`photopea` — both **ABSENT on host** (verified). Confirm the operator's
real PSD need (deliverable vs source) before investing (§4).
⁴ Rendered copies of templates/README produced only where a human-readable
export is wanted; the `.toml`/`.env`/unit files stay as tracked source.

**Transform → engine map (all Docs Chain builtins except the two `exec:`):**
`md→html` = `pandoc-html`; `html→pdf` / `md→pdf` = `weasyprint-pdf`;
`md→docx` = `pandoc-docx`; `html→html` colorize (§11.4.23) = `colorize-html`;
members-fingerprint (§11.4.86) = `members-fingerprint`; `mermaid→svg` = `exec:
mmdc`; `svg→png` / `png→psd(flat)` = `exec: convert`; `openapi→html` = `exec:
scalar` (all verified from docs_chain `CONFIG_SCHEMA.md §5` + host `command -v`).

---

## 3. Always-in-sync wiring (Docs Chain + git hooks + fingerprints)

**One-paragraph mechanism.** Every source doc, diagram source, and the OpenAPI
spec is a node in a **Docs Chain context** (`.docs_chain/contexts/*.yaml`, the
verified schema — docs_chain research §3.3); each derived export (`.html`/`.pdf`/
`.docx`/`.svg`/`.png`) is a downstream node reached by a `derive-from` edge whose
`transform` is one of the §2 engines, so `docs_chain sync` regenerates only what
changed and `docs_chain verify` is a read-only PASS/FAIL/SKIP gate that exits
non-zero when any export is stale; change is detected by **content hash, not
mtime** (`internal/hash/hash.go` sha256 — §11.4.86), re-armed by a
`members-fingerprint` node so a stale verify cannot pass on changed inputs; the
gate is wired at **three write-seams** per §11.4.106(F) — a **commit-seam**
PreToolUse/pre-commit hook that refuses a commit whose staged set touches a chain
source while its exports are stale, a **build-seam** pre-build gate
(`docs_chain verify --all` == 0), and the **constitution-pull seam** (§11.4.164
`post_update_hook.sh`); the **interactive API docs regenerate from the spec** by
making `api/openapi.yaml` a fingerprinted source node whose change re-runs
`exec: scalar` (A2 static) and re-derives the served `/openapi.json` +
`/openapi.toon` the Go server hands to the vendored Scalar/Redoc UI — so a spec
edit and its docs never drift; conflicts (both sides dirty) are exit-2
never-silent-merge, transform failure is exit-3 rollback (no partial write), and
every run drops evidence to `qa-results/docs_chain/<run-id>/`.

**Contexts (extend the docs_chain research §3.2 set):**

| Context | Sources | Derived | Buildable when |
|---|---|---|---|
| `mvp_planning_docs` | 5 top-level `.md` | html/pdf/docx | now (post-vendor) |
| `research_docs` | `research/*.md` | html/pdf | now (post-vendor) |
| `product_docs` (NEW) | B/C/D/E manuals·guides·tutorials·FAQ + G1 | html/pdf(/docx) | after docs authored |
| `api_docs` (NEW) | `api/openapi.yaml` (fingerprint node) | A2 static html+pdf; triggers A3 asset re-embed | after `scalar` wired |
| `diagrams` (NEW) | DG1–DG11 Mermaid/OD sources | svg/png/pdf, fingerprinted | after `mmdc`/OD wired |
| `sql_schema` (NEW) | `migrations/*.sql` + G1 | G1 html/pdf + DG3 | after schema authored |
| `workable_items` | future DB ↔ Issues.md | summary/html/pdf | **blocked** (DB+generators absent — docs_chain research §3.4) |

**Honest boundary (§11.4.106(F)):** the seam hooks give
**eventual-consistency-at-every-write**, not literal real-time. The Docs Chain
`watch` daemon (fsnotify — docs_chain research §2.1) is the optional real-time
add-on; if adopted it is a tracked §11.4.197 upgrade, never claimed as shipped.
`state.json` + `*.docs_chain.tmp` are gitignored (§11.4.28); the SSoT docs
(and the future workable_items DB, §11.4.95) are tracked.

---

## 4. OpenDesign diagrams — authoring, tokenizing, pixel-proof (§11.4.162/.170)

**Authoring pipeline (per diagram DG1–DG11):**
1. **Source of truth = Mermaid** (constitution-native; renders in Markdown/artifacts;
   `mmdc` **verified present** on host). Every diagram's canonical source is a
   `.mmd`/fenced block — versionable, diff-able, regenerable.
2. **Tokenization via OpenDesign.** The diagram theme is driven by our brand
   `design-tokens.json` (`od-design-tokens/v1`, opendesign research §2.2): map
   `--accent`/`--surface`/`--fg`/`--border` into Mermaid `themeVariables` (or,
   for the polished/branded deliverables, author the diagram through the OD
   `hand-drawn-diagrams` / `frame-flowchart-sticky` skills against our
   `helixknowledge/DESIGN.md`). One brand definition, both light and dark tiers
   authored (opendesign research §3.4 — dark is brand-dependent, we author it).
3. **Export.** `mmdc -o DGn.svg` → `convert DGn.svg DGn.png` → PDF (weasyprint/
   embed). SVG is a first-class OD pipeline format (opendesign research §2.4).
4. **Validation = §11.4.170 host-rendered pixel proof** (§5).

**Honest OpenDesign gaps (verbatim from opendesign research §7 — do NOT claim
what OD lacks):**
- **G1 — PSD unsupported.** OpenDesign has **zero** PSD path. Layered PSD needs
  `gimp`/`photopea` — **both ABSENT on this host** (verified). Flat PSD via
  `convert` only. Confirm the operator's real PSD need before investing; track
  as a §11.4.197 item, never bluff a layered export.
- **G2 — Figma is round-trip, not native `.fig`.** OD produces `.od-figma.json`
  re-imported via the OD Figma plugin (desktop only); SVGs **rasterize** on
  import. Editable Figma layers are achievable with fidelity caveats — no native
  `.fig` is generated.
- **G3 — Mobile = framed HTML prototypes, not native kits.** No Material/Android,
  HarmonyOS, or Aurora design system ships; native mobile themes are ours to
  author (§11.4.74 extend-upstream).
- **G4 — `@open-design/components` is `private:true`** — not a consumable lib;
  we build from tokens, not from that package.
- **G5 — dark mode is brand-dependent** (57/151 systems ship a dark tier) — we
  author the dark tier where missing.
- **Tooling gate:** the OD MCP server / `od` CLI presence is **UNCONFIRMED on
  this host**; OpenDesign is **not vendored** (operator-gated, §6 blocker B2).
  Until then, diagrams render via Mermaid+`mmdc`+`convert` (all present) with
  brand tokens applied — the OD-authored/branded polish layer is the gated part.

---

## 5. Anti-bluff render validation (§11.4.168) — what "delivered" requires

A doc/diagram is **NOT "delivered"** on a green `sync`; it is delivered only with
captured proof it **renders correctly AND is in sync**. Per §11.4.168
(exported-doc independent **content + textual + full-visual** validation) +
§11.4.170 (device-independent host-rendered pixel proof), the evidence set is:

1. **In-sync proof.** `docs_chain verify --all` exit 0 + the §11.4.86 fingerprint
   sidecar matches live inputs + the three seam hooks green. Evidence:
   `qa-results/docs_chain/<run-id>/`.
2. **Content parity.** Text-extract the `.html`/`.pdf`/`.docx` (and OCR the PDF)
   and assert it contains the source `.md`'s headings + section text — no
   dropped section, no truncation (the §11.4.108 SOURCE→ARTIFACT check applied to
   docs).
3. **Textual integrity.** No conflict markers, no broken tables, no empty code
   blocks, no missing §11.4.44 revision header.
4. **Full-visual proof (§11.4.170).** Render each `.html`/`.pdf` **and each
   diagram PNG** on the host; run the OCR/vision layout oracle: **no overflow /
   clip / overlap / label-over-label / off-screen / blank page**; every diagram
   node label is present and legible.
5. **Interactive API docs (A3).** The served `GET /docs` returns 200, the vendored
   Scalar/Redoc renders all **21** OpenAPI paths (host-rendered Playwright/
   chrome-devtools screenshot), and a network trace proves **zero external-CDN
   requests** (self-hosted, CSP-clean) — the §11.4.190 web-UI mandate binds this
   surface (responsive across Chromium/Firefox/WebKit, SEO, OpenDesign-branded,
   light+dark, all pixel-proven).
6. **Self-validated oracle (§11.4.107(10)).** The render-validator ships a
   golden-good doc (renders → PASS) + a golden-bad doc (truncated / overflowing
   → MUST FAIL) + a negative control, wired into meta-test; a validator that
   passes its golden-bad is itself the bluff.

All evidence lands under `docs/qa/<run-id>/` (§11.4.83) + `qa-results/docs_chain/`.
A doc reported delivered without artifacts 1–4 (5–6 where applicable) is a §11.4
PASS-bluff at the documentation layer.

---

## 6. Phasing + tracked items (§11.4.197) — extend the docs phase

Maps to **IMPLEMENTATION_PLAN P13 (Docs, Docs Chain, packaging)** + **P10
(OpenDesign)**. Existing tasks are honored; NEW tasks extend the docs phase so
the R18 surface is fully wired, not partially landed.

| Task | Status | Deliverable | Evidence gate |
|---|---|---|---|
| **P13.T1** Docs Chain incorporated | existing, **gated** | vendor engine into `helix_skills`; build binary here | resolves per §11.4.28C; host-built `doctor --all` clean; `sync --all` + `verify --all` exit 0 (docs_chain research §5) |
| **P13.T2** Nano-detail docs + Mermaid | existing, **extend** | author the full §1 surface (Classes B–E, G1) with §11.4.44 headers | docs build; every Mermaid renders (`mmdc`) |
| **P13.T3** Packaging (zip + tar.gz) | existing | H9 | both archives + checksums + extract-build smoke |
| **P13.T4** (NEW) API docs static + interactive | new | A2 (`scalar` render) + A3 (Go-served, no CDN) | A2 html+pdf render-validated; A3 served 200 + 21 paths + zero-CDN trace (§5.5) |
| **P13.T5** (NEW) Docs-sync wiring | new | the 5 new contexts (§3) + 3 seam hooks + fingerprints | `verify --all` == 0 at commit/build/pull seams; paired §1.1 mutation (stale export → gate FAILs) |
| **P13.T6** (NEW) Render-validation harness | new | §11.4.168/.170 content+textual+visual oracle | golden-good PASS / golden-bad FAIL / negative-control (§5.6) wired into meta-test |
| **P10.T1** Vendor OpenDesign | existing, **gated** | `submodules/open_design/` + brand `DESIGN.md`/tokens (H6) | submodule resolves; license/inheritance pointers; dark tier authored |
| **P10.T3/T4** OD diagrams + exports | existing, **extend** | DG1–DG11 branded + SVG/PNG/PDF (PSD gated) | each diagram render + pixel-proven (§5.4); PSD tracked as a gap, not bluffed |

**Honest blockers — the doc-delivery surface is NOT yet deliverable, and why (§11.4.6):**
- **B1 — Docs Chain not vendored (G14/X1 OPEN).** The parent-priority vs
  single-canonical policy is unresolved (GAPS G14; IMPLEMENTATION_PLAN X1 —
  "no second-copy `--apply` until resolved"). Recommended framing: parent copy =
  the one *logical* canonical, `submodules/docs_chain/` = read-only mirror pinned
  by `sync_submodules.sh`; **escalate for operator sign-off before any vendor/
  `--apply`.** Also: shipped binary wrong-arch → `go build ./cmd/docs_chain` here
  (offline module resolution UNCONFIRMED). No live `sync`/`verify` capturable until built.
- **B2 — OpenDesign not vendored.** Adding `submodules/open_design/` (third-party,
  §11.4.28(C)) is a §11.4.66 operator action; the OD MCP/CLI presence is
  UNCONFIRMED on host. Diagrams render via Mermaid today; the OD-branded polish
  layer is gated.
- **B3 — A3 serving feature not yet implemented (NOT a compile failure).** The Go
  backend compiles clean (`go build ./...`=0 at the committed baseline; the
  REQUIREMENTS.md:111-120 "FAILS" line is the pre-P0-build-fix historical baseline,
  superseded by commit `5532e2b`). A3 interactive API docs cannot be served until
  the P13.T4 serving feature (the `/docs` handler + vendored Scalar/Redoc assets +
  `/openapi.{json,toon}` endpoints) is built — a feature gap, not a build gap.
- **B4 — Layered-PSD tooling absent.** `gimp`/`photopea` ABSENT (verified); OD has
  no PSD path. Confirm the real PSD need (§4 G1) before investing; flat-only via `convert`.
- **B5 — §11.4.44 revision headers absent** on the planning docs (docs_chain
  research §3.1) — author-owned, must be added before the chain treats them as
  compliant (Docs Chain does not author them).
- **B6 — `workable_items` chain blocked** — the DB (`docs/workable_items.db`),
  the md↔db binary, and `generate_issues_summary.sh` do not exist yet
  (§11.4.93/.95; docs_chain research §3.4). Do not register that context as working.
- **B7 — `api/openapi.yaml` has no native Docs Chain node kind** — rendered only
  via an `exec: scalar` transform (designed §3), not a builtin.
- **B8 — Consuming repo root undecided** — where `.docs_chain/` lives and what
  `--root` points at is a §11.4.35 operator declaration (docs_chain research §5.6).

Every blocker maps to a tracked task above so no started work stays un-wired
(§11.4.197). This document IS the R18 design artefact for those tasks;
**designed ≠ working** — the tasks stay OPEN until their evidence gates (§5) are captured.

---

## Sources verified (this session)
- Constitution anchors §11.4.12/.44/.45/.53/.56/.57/.59/.60/.65/.86/.106/.107(10)/.162/.164/.168/.170/.190/.197 (project CLAUDE.md, read in context).
- `REQUIREMENTS.md` (R10:80, R12:85-89, R18:206-215, R6, R9), `SPEC.md` (§4-§8), `IMPLEMENTATION_PLAN.md` (P10/P13/X1/X4 + §3 architecture), `api/openapi.yaml` (21 paths; TOON+JSON; tags), `GAPS_AND_RISKS_REGISTER.md` (G14/X1).
- `research/opendesign_incorporation.md` (§2,§4,§7 gaps G1–G6), `research/docs_chain_incorporation.md` (§2 capabilities, §3 contexts, §4.3 host tooling, §5 gaps).
- Host tool presence (`command -v`, this session): `mmdc` ✅, `scalar` ✅, `convert` ✅, `npx`/`node` ✅; `redoc*`/`swagger-ui`/`rsvg-convert`/`inkscape`/`gimp`/`photopea` **ABSENT**. `pandoc`/`weasyprint` ✅ (docs_chain research §4.3).
