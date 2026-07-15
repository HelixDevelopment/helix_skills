# R25 — Canonical `hxs` Project-Key Rename: Airtight Integrity-Gated Design (DESIGN ONLY)

**Revision:** 1
**Last modified:** 2026-07-15T18:54:27Z
**Status:** DESIGN ONLY — execution is a dedicated LATER pass (post doc-churn freeze). This document does NOT execute any rename; it produces the airtight, re-runnable specification so the apply pass has zero risk of orphaned/broken references (§11.4.6).
**Scope:** the HelixKnowledge Skill Graph System MVP doc set + `project/` code references, under `docs/research/mvp/Agent_AI_Skill_Tree_Development/`.
**Composes:** §11.4.29 (lowercase) · §11.4.54 (stable auto-incremental workable-item ids) · §11.4.151 (`<prefix>-<version>` release naming) · §11.4.6 (no-guessing — prove every reference resolves) · §11.4.124 (investigate-before-remove — a rename that drops an id is a violation) · §9.2 (hardlinked backup) · §11.4.113 (absolute no-force-push) · §11.4.66 (surface genuine decisions) · §11.4.84 (working-tree quiescence) · §11.4.86 (fingerprint-driven re-arm).

---

## 1. Goal + operator mandate

### 1.1 Verbatim operator mandate (R25, REQ-025, 2026-07-15)

> "Since the project name is HelixSkill, the key which MUST BE used everywhere MUST BE: **hxs**. All occurrences MUST BE fully updated, all references too! Everywhere!"

Recorded in `REQUIREMENTS.md` R25 (lines 297–306) and `requests/history.md` REQ-025 (lines 50–63). The requirement is explicit that:

- The canonical project key is **`hxs`** (lowercase — §11.4.29).
- It is used **everywhere** — as the ticket / workable-item id prefix (`hxs-NNN`, superseding the `G0x`/`R0x`/`P0x` scheme) AND as the §11.4.151 release/version prefix (`hxs-<version>`).
- ALL occurrences AND ALL references MUST be updated, with **ZERO broken references** (§11.4.6: a rename is forbidden without proving every reference resolves).
- Execution is a **DEDICATED integrity-gated pass** (deterministic id-map + a post-rename gate asserting zero orphan old-ids + every `hxs-NNN` resolvable + no broken links), run **once the current design-doc churn settles** so the FULL doc set renames atomically.

### 1.2 What this document delivers (and does NOT)

Delivers: the id-namespace inventory (real counts + re-runnable regexes), the id-space decision (A vs B with a flagged recommendation), the deterministic id-map SPEC (algorithm, not a frozen hand-list), the integrity gate + §1.1 mutation design, the atomic apply plan, the release-prefix + G40-DB coupling, sequencing, and open decisions.

Does NOT: rename anything, edit any existing doc, touch `project/`, run git. The frozen id-map and the actual rewrite are products of the LATER apply pass, generated from the then-current tree (§4.1 — the map is re-generated at apply time because the doc set is still churning).

---

## 2. Id-namespace inventory (real counts — re-runnable, §11.4.6)

**Doc set (`$DOCSET`), all under the absolute MVP path** `docs/research/mvp/Agent_AI_Skill_Tree_Development/`:

```
REQUIREMENTS.md GAPS_AND_RISKS_REGISTER.md IMPLEMENTATION_PLAN.md \
CONTINUATION.md requests/history.md SPEC.md plan.md research/*.md
```

`project/` code is a secondary, smaller reference surface (§2.4).

Counts below were produced 2026-07-15T18:xx UTC by the exact commands shown; re-run at apply time (§7) because the tree is still churning.

### 2.1 Distinct id namespaces found

| Axis | Prefix | Regex used (re-runnable) | Distinct ids | Total occurrences (doc set) | Kind |
|---|---|---|---|---|---|
| **Gaps** | `G` | `grep -rEoh '\bG[0-9]{2}\b' $DOCSET` | **50** (G01–G50; G50 SUPERSEDED → **49 live** G01–G49) | **1,473** | **Workable items** (findings/defects/tasks) |
| **Requirements** | `R` | `grep -rEoh '\bR[0-9]{1,2}\b' $DOCSET` | **25** (R1–R25; zero-padded `R01`/`R09` variants also appear in one design doc) | **839** (INCLUDES domain-term false positives — §2.3) | Requirement clusters (scope SoT) |
| **Phases** | `P<n>` / `P0.5` | `grep -rEoh '\bP[0-9]+(\.5)?\b' $DOCSET` | 16 (P0, P0.5, P1–P14) | P0×82, P1×65, P8×48, P3×45, … | WBS structure (ordered) |
| **Phase tasks** | `P<n>.T<m>` | `grep -rEoh '\bP[0-9]+(\.[0-9]+)?\.T[0-9]+\b' $DOCSET` | **60** (P0.T1 … P14.T2) | **162** | WBS structure (2-level: phase.task) |
| **Compound remediation** | `P0.5.G<nn>` | `grep -rEoh '\bP[0-9]+(\.[0-9]+)?\.G[0-9]{2}\b' $DOCSET` | 9 (P0.5.G01…G14) | ~9 | **NOT independent** — action-plan restatement of the matching `Gxx` (DEDUP, §3.3) |
| **Cross-cutting tracks** | `X<n>` | `grep -rEoh '\bX[1-4]\b' $DOCSET` | 4 (X1–X4) | many (X1 submodule / X2 git-growth / X3 gates / X4 design) | Structural track axis |
| **Intake-ledger** | `REQ-<nnn>` | `grep -rEoh '\bREQ-[0-9]+\b' $DOCSET` | 7 (REQ-019…REQ-025) | ~9 | §11.4.208 request-ledger records (intake axis, distinct from `R`) |
| **Review follow-ups** | `O<n>` | `grep -rEoh '\bO[1-3]\b' $DOCSET` | 3 (O1/O2/O3 — from the G02/Fable review) | several | Local sub-tracking tokens |
| **Tutorial docs** | `D<n>` | (in `research/r18_documentation_delivery_design.md`) | D1–D5 (also `D8` appears — a **false positive**, §2.3) | several | Doc-internal labels |
| **ATM tickets** | `ATM-NNN` | `grep -rEoh '\bATM-[0-9]+\b' $DOCSET` | **0** | **0** | The constitution's canonical scheme was NEVER adopted here — this IS the G47/§11.4.54 tension R25 resolves |

**Per-file G-id occurrence distribution (top of the churn surface):** `GAPS_AND_RISKS_REGISTER.md` 161 · `CONTINUATION.md` 140 · `research/constitution_compliance_audit.md` 132 · `research/p05_completion_audit_and_discovery.md` 121 · `research/testing_infrastructure_plan.md` 82 · `research/r19_anthropic_api_support_design.md` 81 · `research/p05_high_defect_fix_designs.md` 81 · `research/ops_hardening_design.md` 73 · (…full list is `for f in $DOCSET; do echo "$(grep -Eo '\bG[0-9]{2}\b' "$f"|wc -l) $f"; done | sort -rn`). This spread confirms the operator's "everywhere": the id set is referenced across **≥30 files**, and single ids fan out widely (e.g. `G01` in 16 files, `G14` in 17 files).

### 2.2 Where each axis is DEFINED (definition-anchored — load-bearing for the id-map)

The apply pass MUST distinguish a **definition** (the one place an id is introduced) from a **cross-reference** (every other mention). Definitions:

- **G-ids:** two shapes in `GAPS_AND_RISKS_REGISTER.md` — (a) H3 findings `^### G[0-9]{2} — …` (G01–G38, **38** headings) and (b) compact bullet defs `^- \*\*G[0-9]{2} \(…` (G39–G49, under the "R23 audit" section, lines 404–456). The migration parser MUST handle BOTH shapes (already flagged in `research/g40_workable_items_db_adoption_design.md` §0).
- **R-ids:** `^- \*\*R[0-9]{1,2} — …` bullets in `REQUIREMENTS.md` (lines 52–104 for R1–R14; 198–306 for R15–R25) — **25** definitions.
- **Phases:** `^## P[0-9]+(\.5)? — …` headings in `IMPLEMENTATION_PLAN.md` (lines 106–314).
- **Phase tasks:** `- \*\*P[0-9]+(\.[0-9]+)?\.T[0-9]+ — …` bullets in `IMPLEMENTATION_PLAN.md`.
- **Cross-cutting tracks:** `## X[1-4] —` headings in `IMPLEMENTATION_PLAN.md` (lines 325–341).
- **REQ-NNN:** `^### REQ-[0-9]+ —` headings in `requests/history.md`.

Everything else that matches the id grammar is either a cross-reference (→ must be rewritten in lockstep) or a false positive (§2.3, → must NOT be rewritten).

### 2.3 CRITICAL: false-positive collisions (blind token substitution is FORBIDDEN)

The id grammar collides with real technical vocabulary. Confirmed collisions (captured evidence — `research/skill_granularity_and_composition.md:363-364`):

- **`R8`** — the **Android R8 shrinker/optimizer** ("D8/R8, DEX, shrink/obfuscate"). A blind `\bR8\b → hxs-…` rewrite would CORRUPT a domain skill description. `R8` the shrinker is inside the 839 `R[0-9]{1,2}` matches but is NOT a requirement id.
- **`D8`** — the **Android DEX compiler** ("D8/R8 lowering"). Collides with the `D<n>` tutorial-label grammar.
- Residual risk class: any two-char `[A-Z][0-9]` technical token (part numbers, register names, protocol versions) that appears in skill content.

**Design consequence (load-bearing, applies to BOTH Model A and Model B):** the id-map + rewrite MUST be **definition-anchored and deny-list-aware**, never a blind `sed s/\bG[0-9][0-9]\b/…/g`. An id is rewritten ONLY if it is in the DEFINED id universe (§2.2) OR is a cross-reference whose token matches a defined id AND is not on the domain-term deny-list. The deny-list (seeded: `R8`, `D8`; extended at apply time by a scan of skill/technology content) is itself captured evidence and reviewed before the pass.

### 2.4 `project/` code reference surface (secondary)

Real refs in `project/` (word-boundary + `§`-prefixed in-code forms):

- **`§G0x` in-code comment citations — 17 occurrences:** `§G03`×9, `§G21`×4, `§G05`×3, plus `G07`×5, `G05`×3, `G21`×2, `G16`×2, `G03`×2, `G02`×2 across `internal/api/server.go`, `internal/api/skills_handler.go`, `internal/api/skills_validation_test.go`, `internal/mcp/server.go`, `internal/mcp/tools.go`, `internal/models/skill.go`, `internal/validation/{pipeline,sandbox}.go` + their tests, `internal/skill/migration_granularity_test.go`. (Command: `grep -rEoh '§?G[0-9]{2}\b' project/ --include='*.go' | sort | uniq -c`.)
- **`R<n>` cites in seed data + compose:** `project/deploy/docker-compose.yml` (R15), `project/seed/CORPUS.yaml` (R13/R6), `project/seed/skills/*.toml` (R6/R13), `internal/models/skill.go` + `internal/skill/graph.go` (`// NEW (R16)`), `migrations/002_granularity.up.sql` (R16). **These R-cites are heavily entangled with domain content** (the seed TOMLs describe Android/Java skills) — the deny-list discipline (§2.3) is essential here.
- **`P<n>.T<m>` cites:** granularity-migration files under `internal/skill/`, `internal/db/`, `migrations/` (`P1.T1` etc.).

**No `hxs` collision:** `grep -rn 'hxs'` across the whole tree returns ONLY the R25 mandate text itself (REQUIREMENTS.md, requests/history.md, CONTINUATION.md) — the new key is collision-free.

---

## 3. Id-space decision — Model A vs Model B (recommendation + flagged sign-off)

### 3.1 The two candidate models

**(A) UNIFIED.** One monotonic `hxs-NNN` sequence covering ALL workable items — gaps + requirements + phases collapse into a single flat id space, with an old→new mapping table.

**(B) SEMANTIC-PRESERVING.** Keep the gap / requirement / phase / track distinction but re-key each axis under `hxs`. Two sub-variants:
- **(B1)** re-prefix each axis: `Gxx → hxs-g-NN`, `Rxx → hxs-r-NN`, `Pn.Tm → hxs-p<n>.t<m>`, `Xn → hxs-x<n>`.
- **(B2)** convert ONLY the workable-item **G-ids** to a clean monotonic `hxs-NNN` ticket sequence (pure Model A on the ticket axis), and keep R / P / X as semantic axes re-prefixed under `hxs` (`hxs-r-NN` / `hxs-p<n>` / `hxs-x<n>`).

### 3.2 Analysis against the mandate + §11.4.54

| Consideration | Model A (flat) | Model B (semantic) |
|---|---|---|
| Operator letter — `hxs-NNN` "superseding G0x/R0x/P0x" | Literal match (one flat sequence) | Met by re-prefix; the *shape* differs from `hxs-NNN` for R/P |
| Operator spirit — "hxs everywhere, one canonical key" | Satisfied | Satisfied — every id becomes `hxs-…`, no bare `G/R/P` survives |
| §11.4.54 scope — targets **workable items** | Over-applies: forces requirement-clusters + WBS structure into a ticket sequence they are not | Correctly scopes the `hxs-NNN` ticket sequence to the workable items (gaps); keeps R/P as their real axes |
| Phase ordering / hierarchy | **Lost** — `P3`→`hxs-071` destroys "after P2"; `Pn.Tm`'s 2-level hierarchy cannot survive a flat number | Preserved — `hxs-p3`, `hxs-p3.t2` keep ordering + hierarchy |
| Apply safety | Requires an AUTHORED merge-order decision ("what number is R13 in the merged sequence?") for 49 G + 25 R + 60 P.T — churn-sensitive, error-prone, needs a human | Mechanical + reversible: re-prefix is a deterministic transform; only the G-ticket sequence needs monotonic assignment |
| Cross-ref blast radius | Same ~2,300 doc occurrences either way | Same |
| Information loss (§11.4.6-adjacent) | Conflates 3 entity kinds (finding / requirement / WBS node) into one namespace | None — each entity keeps its kind |

### 3.3 DEDUP fact (applies to both models)

`P0.5.Gxx` lines in `IMPLEMENTATION_PLAN.md` are the **action-plan restatement of the matching `Gxx` finding** (same fix, same file:line citations — verified `IMPLEMENTATION_PLAN.md:127-133` "P0.5.G01 — Wire internal/api…" ≡ `GAPS_AND_RISKS_REGISTER.md` G01 DECISION line; this is the load-bearing dedup fact from `research/g40_workable_items_db_adoption_design.md` §0). Therefore `P0.5.Gxx` maps to the SAME `hxs` id as its `Gxx` — it is NOT a new id. Likewise `G50` is SUPERSEDED (register lines 410–411) → excluded from the live sequence; its supersession note is preserved as history (§11.4.124 — never dropped).

### 3.4 RECOMMENDATION — **Model B2** (pure `hxs-NNN` for gaps; semantic re-prefix for R/P/X), FLAGGED for operator sign-off

**Recommend B2** because it is the only model that simultaneously honors the operator's letter AND spirit AND §11.4.54's actual scope AND the safe-apply principle:

1. **§11.4.54 targets workable items; gaps ARE the workable items** → they get the clean monotonic `hxs-NNN` ticket sequence the operator literally wrote (`hxs-001`, `hxs-002`, …). This is the exact "hxs-NNN superseding G0x" instruction, applied where it belongs.
2. **Requirements (R) and phases (P) are DIFFERENT axes.** Requirements are a stable scope SoT (`REQUIREMENTS.md`, cited 839× incl. seed code); phases are WBS structure with ordering + a 2-level `phase.task` hierarchy. Collapsing them into one flat sequence (Model A) destroys the ordering/hierarchy and conflates three entity kinds — a real information-loss + readability regression, not an improvement.
3. **"hxs everywhere" is FULLY satisfied** by the re-prefix: `R13 → hxs-r13`, `P3 → hxs-p3`, `P1.T1 → hxs-p1.t1`, `X1 → hxs-x1`. No bare `G`/`R`/`P`/`X` id survives; `hxs` becomes the single canonical key across every axis.
4. **B2 is the SAFER apply.** Re-prefixing is a deterministic, collision-guarded, reversible transform that needs NO human "what number does R13 become" decision. Only the G-ticket sequence requires monotonic assignment, and that is fully mechanical (sort by G-number).

**FLAG (genuine operator/conductor decision — surface via §11.4.66):** the operator's LITERAL phrasing `hxs-NNN` "superseding G0x/**R0x/P0x**" is compatible with a strict Model A reading (fold R and P into the SAME flat sequence). B2 keeps R/P as `hxs-r-`/`hxs-p-` sub-axes instead. This is the ONE genuinely ambiguous point. The design defaults to **B2** with a strong rationale, but the apply pass MUST NOT execute until the operator confirms whether:

- **(D1)** R/P/X stay semantic sub-axes (`hxs-r-`/`hxs-p-`/`hxs-x-`) — **B2, recommended**; OR
- **(D2)** R and P fold into the one flat `hxs-NNN` ticket sequence — **strict Model A**.

Both are fully specified by the id-map algorithm in §4 (the algorithm is parameterized by this choice); only the axis-fold flag differs.

---

## 4. Deterministic id-map SPEC (algorithm, re-runnable — NOT a frozen hand-list)

The doc set is still churning (38 research docs + register + plan actively edited). The frozen `old_id → hxs_id` map is therefore a BUILD ARTIFACT generated at apply time by this algorithm run against the then-current tree — NOT hand-authored now (a hand-list would be stale before the pass runs). This section is the algorithm + invariants.

### 4.1 Algorithm

**Inputs:** the doc set `$DOCSET`; the axis-fold flag (D1/D2 from §3.4); the domain-term deny-list (seed `{R8, D8}`, extended by a content scan).

**Step 1 — enumerate the DEFINED id universe (definition-anchored, §2.2), not blind tokens.**
```
G-defs (H3):   grep -rEn '^#{3,4} G[0-9]{2} '                 GAPS_AND_RISKS_REGISTER.md
G-defs (bullet):grep -rEn '^- \*\*G[0-9]{2} \('              GAPS_AND_RISKS_REGISTER.md
R-defs:        grep -rEn '^- \*\*R[0-9]{1,2} '               REQUIREMENTS.md
P-defs:        grep -rEn '^#{2} P[0-9]+(\.5)? '              IMPLEMENTATION_PLAN.md
PT-defs:       grep -rEn '\*\*P[0-9]+(\.[0-9]+)?\.T[0-9]+ '  IMPLEMENTATION_PLAN.md
X-defs:        grep -rEn '^#{2} X[1-4] '                     IMPLEMENTATION_PLAN.md
```
This yields the authoritative id set. Any grammar-matching token NOT in this set is a cross-reference (map it) or a false positive (deny-list — skip).

**Step 2 — assign `hxs` ids.**
- **Gap ticket axis (both models):** sort defined G-ids by numeric value; assign `hxs-001`, `hxs-002`, … monotonically, zero-padded ≥3 digits, in G-number order (append-only). `G50` (superseded) and every `P0.5.Gxx` (dedup, §3.3) are EXCLUDED from new-number assignment — `P0.5.Gxx` inherits its `Gxx` mapping.
- **Under D1 (B2, recommended):** `R<n> → hxs-r<zeropad2>`, `P<n> → hxs-p<n>`, `P<n>.T<m> → hxs-p<n>.t<m>`, `P0.5 → hxs-p0.5`, `X<n> → hxs-x<n>`.
- **Under D2 (strict A):** append R-defs then P/P.T-defs into the SAME monotonic `hxs-NNN` sequence after the gaps (order = operator-confirmed; default gaps→requirements→phases, ascending within each). A published `hxs-NNN → (old_id, kind)` table preserves the axis semantics that the flat number drops.

**Step 3 — persist an APPEND-ONLY, monotonic state file `hxs_id_map.json`** (jsonl, one row per mapping), schema:
```json
{"old_id":"G14","hxs_id":"hxs-014","axis":"gap","defined_in":"GAPS_AND_RISKS_REGISTER.md:189",
 "aliases":["P0.5.G14"],"assigned_at":"<ISO-UTC>","superseded":false}
```
Invariants (§11.4.54): once assigned an `hxs_id` is **never** renumbered, reused, decremented, or gapped; the sequence is monotonic with no holes; the map is append-only. The `defined_in` field is the binding key (survives heading reflow). Aliases (`P0.5.Gxx`) are recorded so the rewrite maps them to the same id.

**Step 4 — rewrite (definition-anchored + deny-list-aware).** For every file in `$DOCSET` (+ the §2.4 `project/` refs if in scope — decision D4), replace each occurrence of every mapped `old_id` — its definition AND every cross-reference (including `§G0x` in-code forms) — with its `hxs_id`, using an anchored word-boundary transform that:
- SKIPS any token on the domain-term deny-list (`R8`, `D8`, extended);
- SKIPS any token not in the mapped universe (unknown two-char tokens are left literal, never guessed — §11.4.6);
- handles the zero-padded variants (`R01`/`R09`) by normalizing to the canonical id before lookup;
- rewrites compound `P0.5.Gxx` to the `Gxx`→`hxs` mapping (dedup), and `Pn.Tm` under D1 to `hxs-p<n>.t<m>`.

**Step 5 — emit the human-readable crosswalk** (`hxs_id_map.md`, exported per §11.4.65) listing every `old_id → hxs_id`, kind, and definition location — the auditable record + the alias source for anyone reading historical git commits / external trackers.

### 4.2 Why the map is generated at apply time (not frozen here)

Between now and the apply pass, new G-ids (the register is actively growing — G39–G49 were appended this session; more may land), new research docs, and new `Pn.Tm` tasks WILL be added. Freezing a hand-list now guarantees drift. The re-runnable algorithm (§4.1) regenerates the authoritative universe from the CURRENT tree at apply time; §11.4.86 fingerprinting (§7) proves the map matches the tree it renames.

---

## 5. Integrity gate `CM-HXS-RENAME-COMPLETE` + §1.1 mutation

A deterministic PASS/FAIL gate that runs AFTER the rewrite and BEFORE the commit (§11.4.6 — a rename is forbidden without proving every reference resolves). GREEN is a hard precondition of the commit.

### 5.1 Gate invariants (all must hold)

1. **Zero orphan old-ids.** No live `G0x`/`R0x`/`P0x`/`Pn.Tm`/`Xn` id token survives anywhere in `$DOCSET` (+ in-scope `project/` refs), EXCEPT: (a) tokens on the domain-term deny-list (`R8`/`D8`/extended), (b) the `hxs_id_map.{json,md}` crosswalk itself (which legitimately records old ids), (c) historical/supersession notes explicitly quoting an old id as history (allow-listed by line). Check:
   `grep -rEn '\b(G[0-9]{2}|R[0-9]{1,2}|P[0-9]+(\.[0-9]+)?(\.T[0-9]+)?|X[1-4])\b' $DOCSET | <filter deny-list + crosswalk + allow-list> | wc -l == 0`.
2. **Every `hxs-NNN` resolves to exactly one definition.** Build a definition index from the rewritten tree; assert every gap `hxs-NNN` has exactly one definition heading/bullet (no dup, no missing).
3. **Sequence integrity (§11.4.54).** The gap `hxs-NNN` sequence is monotonic with NO gaps and NO duplicates; the `hxs_id_map.json` is append-only and internally consistent (every mapped old id appears exactly once).
4. **Every cross-reference resolves.** Every `hxs-…` token used as a reference has a matching definition in the index (no dangling reference).
5. **No broken markdown links / anchors.** Extract every `[text](target)` + intra-doc `#anchor`; resolve each against the rewritten tree; zero unresolved.
6. **Deny-list untouched.** Assert each seed domain term (`R8`/`D8`) still appears verbatim in its original file:line (the rewrite did NOT corrupt a technical term).
7. **Release-prefix present (§6).** `HELIX_RELEASE_PREFIX=hxs` resolvable in `.env.example`; version strings that carry a prefix use `hxs-`.

Each invariant prints its resolved evidence on FAIL (§11.4.201 — a guard reports why it refused), naming the exact offending file:line.

### 5.2 Paired §1.1 mutation (proves the gate is not a bluff)

Run in a scratchpad copy (§11.4.84 — never the real tree). Each mutation MUST flip the gate RED:

- **M1 (orphan):** leave one `G14` un-renamed in one research doc → invariant 1 FAILs, naming the file:line.
- **M2 (dup id):** map two old ids to the same `hxs-014` → invariant 2/3 FAIL.
- **M3 (gap in sequence):** skip `hxs-013` (jump G12→hxs-012, G13→hxs-014) → invariant 3 FAILs.
- **M4 (dangling ref):** insert an `hxs-999` reference with no definition → invariant 4 FAILs.
- **M5 (broken link):** rewrite a `[…](GAPS_AND_RISKS_REGISTER.md#g14)` anchor without updating the target → invariant 5 FAILs.
- **M6 (domain-term corruption):** let the rewrite touch `R8`→`hxs-…` in the Android skill description → invariant 6 FAILs.

A gate that PASSes any of M1–M6 is itself the bluff (§11.4.107(10) self-validation): ship a golden-good fixture (correctly renamed tree → PASS) + these six golden-bad fixtures (→ FAIL) wired into meta-test.

---

## 6. Atomic apply plan (§9.2 backup · §11.4.113 no-force-push)

The apply pass is a SINGLE atomic transaction on a QUIESCENT tree (§11.4.84).

1. **Freeze + fetch (§11.4.37/§11.4.84).** Declare the doc-churn freeze window (§7); `git fetch --all --prune`; integrate any upstream delta FIRST; assert the working tree is clean (no in-flight doc edits, no mutation markers).
2. **§9.2 hardlinked backup (mandatory, near-zero cost).** `cp -al .git ../hxs_rename_backup/repo.git.mirror`; record HEAD, HEAD tree sha256, and a full id-inventory snapshot (the §2 counts) as the pre-op state. Define the expected post-op state (zero orphan old-ids, N `hxs-NNN` defs where N = live-gap count).
3. **Regenerate the id-map (§4.1).** Run the algorithm against the CURRENT tree → `hxs_id_map.json` + `hxs_id_map.md`. Human/operator review of the crosswalk + the deny-list (§2.3) BEFORE the rewrite (§11.4.66 checkpoint if any ambiguous token surfaced).
4. **Single-pass scripted rewrite.** Apply the definition-anchored, deny-list-aware transform (§4.1 step 4) across the FULL doc set (+ in-scope `project/` refs) in ONE pass → all files change together.
5. **Run `CM-HXS-RENAME-COMPLETE` (§5).** GREEN is a HARD gate. On any FAIL → do NOT commit; fix the transform/deny-list; re-run from step 3.
6. **Re-export + revision-bump.** Regenerate `.html`/`.pdf` siblings via the exporter / Docs Chain (couples with G43 export wiring); bump each touched doc's `**Revision:**` + `**Last modified:**` header (§11.4.44); update CONTINUATION (§12.10) + the request ledger.
7. **Commit + push (§11.4.113).** ONE descriptive commit ("rename: adopt canonical `hxs` project key per R25 — id-map + gate evidence") citing the `hxs_id_map.md` crosswalk + the gate PASS. Push **ff-only to ALL upstreams** (§2.1) — **NEVER force-push** (§11.4.113). The commit descends from every mirror tip (merge-onto-latest-main if needed), so every push is a fast-forward.
8. **Rollback path.** Because the rewrite is a single commit on a clean tree, rollback is one operation: `git reset --hard <pre-op HEAD>` (or restore the §9.2 hardlinked mirror) if the gate FAILs post-commit or any downstream check regresses. No history rewrite, no force-push, no lost commits (§9.2).

---

## 7. Sequencing — when to run + how to freeze the churn

**Trigger (post doc-churn settle).** Run the pass ONLY once the design-doc wave settles: the 38 `research/*.md` docs + `GAPS_AND_RISKS_REGISTER.md` + `IMPLEMENTATION_PLAN.md` reach a stable state (no more G-ids being appended, no more research docs landing). Empirically: after the current design-doc round completes and BEFORE (a) the G40 workable-items-DB population (§8) and (b) P-phase implementation ids proliferate into `project/` code.

**Freeze mechanism.** Declare a short doc-freeze window; treat the doc set as single-writer (§11.4.84 quiescence) for the duration — no concurrent doc edits land while the rewrite + gate + export run. Ideally a dedicated focused session (the rewrite + gate is minutes; the value is atomicity).

**Re-run the inventory at apply time (do NOT trust this snapshot).** The §2 regexes are the re-runnable spec; the apply pass regenerates the definition universe + counts from the THEN-current tree. A §11.4.86 fingerprint (sha256 of the sorted defined-id list) is captured pre-rewrite and asserted by the gate, so a stale map cannot rename a tree it does not match. If the fingerprint moved since this design, the map is regenerated — never applied from a stale list.

---

## 8. Release-prefix (§11.4.151) + G40-DB dependency

### 8.1 `hxs-<version>` release/version prefix (§11.4.151)

- The §11.4.151 release-tag + version-name prefix becomes **`hxs`** on the main repo AND every owned submodule in one release (e.g. `hxs-1.0.0-dev-0.0.1`).
- Resolution order per §11.4.151: `HELIX_RELEASE_PREFIX=hxs` in `.env` (git-ignored §11.4.30) documented in the tracked `.env.example` (§11.4.77); else lowercased snake_case project-root dir name. Set `HELIX_RELEASE_PREFIX=hxs` explicitly in `.env.example` so resolution is unambiguous.
- Version codes increment monotonically; the same `hxs` prefix is greppable across main + all owned submodules in one release. The gate's invariant 7 (§5.1) asserts the prefix is present + resolvable.

### 8.2 G40 workable-items-DB coupling (SEQUENCE this)

The G40 design (`research/g40_workable_items_db_adoption_design.md`) adopts the constitution's SQLite workable-items engine as the §11.4.93/§11.4.95 single source of truth, keyed on `items.atm_id` (which it found is **TEXT with no format CHECK** — so `hxs-NNN` is a valid key with ZERO schema change).

**Dependency:** if G40 populates the DB BEFORE the hxs rename, it stores `G01`…`G49` as primary keys, and the rename must then migrate `items.atm_id` + `item_history` + every FK/alias row — extra churn + risk across a live DB (destructive-op territory → §9.2 backup + operator auth).

**Recommendation (SEQUENCE):** **run the hxs rename FIRST**, OR have the G40 DB population key on `hxs-NNN` from row 1 (using this design's id-map). Either way, the DB is keyed on the canonical `hxs-NNN` from the start — the rename never has to touch DB rows. G40's own finding (`atm_id` is TEXT-no-CHECK) means adopting `hxs-NNN` as the key is free. **G40 DB population WAITS for (or co-lands with) this rename.** Record this ordering in the plan so the two passes do not race.

---

## 9. Open decisions (surface via §11.4.66 before the apply pass)

| # | Decision | Recommendation | Why it needs sign-off |
|---|---|---|---|
| **D-A/B** | Model A (flat `hxs-NNN` for everything) vs Model B2 (pure `hxs-NNN` for gaps; `hxs-r-`/`hxs-p-`/`hxs-x-` for the other axes) | **B2** (§3.4) | Operator's literal `hxs-NNN` "superseding …R0x/P0x" is compatible with strict A; B2 honors §11.4.54's workable-item scope + preserves phase ordering/hierarchy. Genuinely ambiguous. |
| **D2 (R-fold)** | Do requirements R1–R25 become flat tickets (fold) or stay a `hxs-r-NN` requirement axis? | Stay `hxs-r-NN` (part of B2) | Requirements are a distinct scope SoT cited 839× incl. seed code; folding is high-churn + conflates axes. |
| **D3 (REQ ledger)** | Re-key intake-ledger `REQ-NNN` → `hxs-req-NNN`, or leave? | Leave (or `hxs-req-`) | §11.4.208 request-ledger is a distinct intake axis, not a workable-item ticket. "hxs everywhere" arguably reaches it — flag. |
| **D4 (code scope)** | Rename the 17 `§G0x` in-code comments + seed-TOML `R`-cites too, or docs-only? | Include code (operator said "everywhere") | Code refs are stable design citations tangled with domain content; the deny-list (§2.3) is essential; confirm scope. |
| **D5 (alias retention)** | Keep a `legacy_id` alias (crosswalk / DB column) so old cross-refs in git history + external trackers still resolve, or hard-supersede? | Keep the `hxs_id_map.md` crosswalk as the durable alias record (G47 option-b spirit) | §11.4.124 — never drop the old id's provenance; historical commits + any external tracker still reference `G14`. |
| **D6 (deny-list)** | Confirm the domain-term exclusions (`R8` Android shrinker, `D8` DEX, + any others found by the apply-time content scan) | Confirm + extend at apply time | A missed deny-list entry corrupts a technical term (§2.3); an over-broad one leaves an orphan id. |

---

## 10. Honest boundary (§11.4.6)

This is a DESIGN. It proves the apply pass CAN be made airtight (definition-anchored id-map + deny-list + self-validated gate + §9.2 backup + ff-only single commit + rollback), and it captures the REAL current inventory (§2) with re-runnable regexes. It does NOT prove the rename is correct on a tree that has not been frozen yet — the frozen id-map, the deny-list extension, and the gate PASS are products of the apply run against the then-current tree (§4.2/§7). The one genuinely ambiguous point (Model A vs B2 / axis-fold) is flagged for operator sign-off (§9) and NOT decided unilaterally.

## Conductor resolutions of the flagged decisions (§11.4.101 / §11.4.66)

**DEFERRED to operator (§11.4.66 — genuine, high-blast-radius ambiguity):**
- **A-vs-B id-space model.** The operator's literal "`hxs-NNN` superseding
  G0x/R0x/P0x" is compatible with BOTH Model A (fold R + P into ONE flat
  `hxs-NNN` sequence) and Model B2 (recommended: `hxs-NNN` for the gap/workable
  axis; `hxs-r-NN` / `hxs-p<n>[.t<m>]` / `hxs-x<n>` semantic-preserved for
  requirements/phases/tracks). Model A loses phase ordering + the phase.task
  hierarchy the IMPLEMENTATION_PLAN depends on; B2 preserves it. Because the
  rewrite touches 2,300+ occurrences and redoing it is costly, this is asked of
  the operator before the apply pass runs — NOT guessed (§11.4.66). The apply
  pass is already deferred (post doc-churn), so this block parks ONE unit and
  stalls nothing (§11.4.101).

**RESOLVED autonomously (§11.4.101 — reversible + evidence-backed):**
1. **REQ-NNN intake-ledger ids = KEEP** (not renamed). They are a request-intake
   channel serial (like ATM-NNN's role), a DISTINCT axis from workable-item ids,
   not a G/R/P prefix the mandate supersedes. `requests/history.md` keeps REQ-NNN.
2. **Code-ref scope = INCLUDE** the in-`project/` `§G0x` comment refs + seed-TOML
   R/P cites in the apply pass ("everywhere / all references" per R25). The rewrite
   is definition-anchored (see finding 1 below), so code refs are re-keyed via the
   same id-map, never blind-substituted.
3. **Legacy-id aliases = NO in-doc aliases.** Bare old-ids are fully removed; the
   old→hxs mapping is captured ONCE in a tracked `MIGRATION_MAP.md` (the §11.4.124
   investigate-before-remove audit evidence) + the §11.4.208 ledger notes the pass.
4. **Deny-list = CONFIRMED exclude** domain terms that collide with the id grammar
   (`R8` Android shrinker, `D8` DEX, any non-definition-anchored token). Only
   definition-anchored ids are re-keyed.

**Load-bearing findings ACCEPTED (bind the apply pass):**
- **No blind token substitution** — the rewrite MUST be definition-anchored +
  deny-list-aware (the `R8`/`D8` collision proves a naive sed would corrupt prose).
- **G40-DB sequencing** — run the hxs rename BEFORE G40 DB population, OR key the
  G40 SQLite `items.atm_id` (TEXT, no CHECK) on `hxs-NNN` from row 1, so the DB
  never needs a keyed-row migration. G40 population therefore waits on this pass.
