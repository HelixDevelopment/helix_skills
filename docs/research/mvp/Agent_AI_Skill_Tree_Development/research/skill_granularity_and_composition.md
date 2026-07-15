# Skill Granularity & Composition (R16)

**Revision:** 1
**Last modified:** 2026-07-15T00:00:00Z
**Scope:** Answers operator question **R16** — *"Research if it makes sense to split Skills for whole big technologies into smaller Skills wrapped by one bigger Skill. We MUST HAVE a proper mechanism for Skill granulation and definition of its relationships and use as building blocks for bigger pieces — whole technologies and stacks."*
**Status:** Design proposal, backward-compatible superset of `SPEC.md` §4/§5 and the 8 seed TOMLs. No code changed; this document is intended to fold directly into `SPEC.md`.
**Evidence discipline:** Every external claim carries a source tag `[Sxx]` resolving to §2's table; every URL in that table was reachability-verified on 2026-07-15 and the raw status line is pasted there. Unverifiable claims are dropped or explicitly labelled.

---

## 1. Verdict

**Big technologies MUST be decomposed into small, independently-authored *atomic* skills aggregated under a thin *umbrella* skill — not stored as one monolithic node.** The decisive reasons are four, each grounded in prior art rather than taste:

1. **Token economics / retrieval quality.** Agents fetch skills on demand into a bounded context window. A monolithic `android` node forces the whole technology into context even when the task needs only the build sub-topic, and long-context models measurably *lose* mid-context information — retrieval accuracy is highest when the relevant material is a small, targeted chunk at a context boundary, and degrades when it is buried in a long block [S15][S16]. This is exactly the failure Anthropic's own Agent Skills design avoids with **progressive disclosure**: load a cheap name+description first, load the body only when triggered, load bundled detail only when needed [S1][S2][S3]. A monolith defeats progressive disclosure; an umbrella + atomic leaves *is* progressive disclosure expressed as a graph.
2. **Reuse (DRY).** `gradle`/`kotlin_dsl`, `activity_lifecycle`, `art_runtime` are each depended on by more than one parent. A monolith duplicates them into every consumer; atomic nodes let a single leaf be a shared dependency — the same "a package is a reusable unit other packages depend on" model every package manager uses [S4][S6][S7][S8][S10].
3. **Independent versioning & co-change.** Jetpack Compose versions on its own BOM cadence, decoupled from the platform API level; the AOSP build system evolves independently of the app framework. Units that change on different clocks belong in different nodes (the software-engineering rule that things which change together live together, and things which change apart live apart) [S11][S12][S13]. A monolith couples independently-versioned parts into one `version` field and one review cycle.
4. **Distinct validation.** The system's value proposition is anti-bluff validated skills. `activity_lifecycle` is validated by a different probe than a `gradle` build; forcing them into one node forces one coarse validation verdict over unrelated evidence. Separable validation ⇒ separable skills.

The counter-cost — more nodes, more edges — is precisely what the existing DAG (`skill_dependencies`, recursive-CTE traversal) already exists to manage, so the marginal cost is near zero while the retrieval/reuse/versioning/validation wins are structural. **Decompose-under-umbrella wins unambiguously.**

---

## 2. Prior-art evidence

All URLs verified reachable on **2026-07-15** with `curl -sSIL` (last HTTP status line pasted) or `git ls-remote` (HEAD SHA pasted). Where a server blocked `HEAD`, a `GET`-fallback code is shown.

| Tag | Source (what it grounds) | URL | Verified status line |
|---|---|---|---|
| S1 | Anthropic Engineering — Agent Skills: composable, progressive-disclosure knowledge units | https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills | `HTTP/2 200` |
| S2 | Claude Docs — Agent Skills: `SKILL.md` frontmatter (name+description) always loaded, body/linked files loaded on demand | https://docs.claude.com/en/docs/agents-and-tools/agent-skills | `HTTP/2 200` |
| S3 | anthropics/skills reference repo (skill layout in practice) | https://github.com/anthropics/skills.git | `git ls-remote HEAD → 9d2f1ae187231d8199c64b5b762e1bdf2244733d` |
| S4 | npm — `package.json` `dependencies` (a package declares the other packages it needs) | https://docs.npmjs.com/cli/v10/configuring-npm/package-json | `HTTP/2 200` |
| S5 | npm — workspaces (one root aggregates many member packages) | https://docs.npmjs.com/cli/v10/using-npm/workspaces | `HTTP/2 200` |
| S6 | Cargo — specifying dependencies (a crate depends on other crates by name+version) | https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html | `HTTP/2 200` |
| S7 | Cargo — workspaces / virtual manifest (a top manifest that only aggregates members) | https://doc.rust-lang.org/cargo/reference/workspaces.html | `HTTP/2 200` |
| S8 | Maven — dependency mechanism (transitive dependency resolution) | https://maven.apache.org/guides/introduction/introduction-to-dependency-mechanism.html | `HTTP/2 200` |
| S9 | Maven — POM: `packaging=pom` + `<modules>` aggregator project (an umbrella POM composed of module POMs) | https://maven.apache.org/guides/introduction/introduction-to-the-pom.html | `HTTP/2 200` |
| S10 | Gradle — declaring dependencies | https://docs.gradle.org/current/userguide/declaring_dependencies.html | `HTTP/1.1 200 OK` |
| S11 | Gradle — multi-project builds (a root project composed of subprojects) | https://docs.gradle.org/current/userguide/multi_project_builds.html | `HTTP/1.1 200 OK` |
| S12 | Go — modules reference (`require` directives; a module is a collection of packages) | https://go.dev/ref/mod | `HTTP/2 200` |
| S13 | Go — "Using Go Modules" (module = versioned unit; deps by path+version) | https://go.dev/blog/using-go-modules | `HTTP/2 200` |
| S17 | W3C SKOS Reference §semantic-relations — hierarchical (`broader`/`narrower`, non-transitive) vs associative (`related`) | https://www.w3.org/TR/skos-reference/#semantic-relations | `HTTP/2 200` |
| S18 | W3C SKOS Primer — modelling `broader`/`narrower`/`related` in practice | https://www.w3.org/TR/skos-primer/ | `HTTP/2 200` |
| S19 | OBO Relations Ontology — `part_of` (mereology) kept distinct from `is_a` (subsumption) | https://oborel.github.io/obo-relations/ | `HTTP/2 200` |
| S20 | Cohesion (software) — a module should hold one closely-related responsibility | https://en.wikipedia.org/wiki/Cohesion_(computer_science) | `HTTP/2 200` |
| S21 | Coupling (software) — minimise inter-module dependency | https://en.wikipedia.org/wiki/Coupling_(computer_programming) | `HTTP/2 200` |
| S22 | Interface Segregation Principle — many small focused interfaces beat one fat one | https://en.wikipedia.org/wiki/Interface_segregation_principle | `HTTP/2 200` |
| S23 | Single-Responsibility Principle — one reason to change per unit | https://en.wikipedia.org/wiki/Single-responsibility_principle | `HTTP/2 200` |
| S24 | Dijkstra EWD447 "On the role of scientific thought" — separation of concerns (primary source) | https://www.cs.utexas.edu/users/EWD/ewd04xx/EWD447.PDF | `HTTP/1.1 200 OK` |
| S15 | "Lost in the Middle" (Liu et al.) — long-context models under-use mid-context info | https://arxiv.org/abs/2307.03172 | `HTTP/2 200` |
| S16 | "Lost in the Middle" — TACL 2024 published version | https://aclanthology.org/2024.tacl-1.9/ | `HTTP/1.1 200 OK` |
| S25 | Anthropic Engineering — effective context engineering (smallest high-signal token set; progressive disclosure) | https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents | `HTTP/2 200` |

**(a) Anthropic Agent Skills / progressive disclosure [S1][S2][S3][S25].** A Skill is a `SKILL.md` whose YAML frontmatter (`name` + `description`) is the always-loaded cheap layer; the Markdown body and any bundled files are pulled only when the skill is triggered. That is a three-tier *granularity* built into the format: metadata → body → linked detail. Our umbrella-node (cheap orientation) + atomic-leaf (deep body, fetched via MCP `skill_get`) reproduces this at graph scale [S25].

**(b) Package-manager dependency graphs [S4]–[S13].** Every mainstream manager models "a unit is composed of / depends on other named, versioned units," and every one provides an explicit *aggregator* construct that is nothing but a manifest over member units: Maven `packaging=pom` with `<modules>` [S9], Cargo virtual-manifest workspaces [S7], npm workspaces [S5], Gradle multi-project root [S11], Go modules-of-packages [S12]. This is the direct precedent for `kind = "umbrella"` (aggregator) vs `kind = "atomic"` (leaf unit), and for a `composes` edge distinct from a `requires`/`depends_on` edge — package managers already distinguish "this project *contains* these modules" (aggregation) from "this project *needs* those artifacts" (dependency).

**(c) Ontology relationship typing [S17][S18][S19].** SKOS separates **hierarchical** relations (`broader`/`narrower`) from **associative** relations (`related`), and crucially defines `broader`/`narrower` as **non-transitive**, with a *separate explicit* super-property (`broaderTransitive`) for the closure [S17]. That is precisely our design: hierarchical structural edges (`composes`, `extends`) that the resolver walks transitively, versus associative advisory edges (`related_to`) that it does not. The OBO Relations Ontology further keeps **`part_of`** (mereology / whole-part) rigidly distinct from **`is_a`** (subsumption) [S19] — the same distinction we draw between `composes`/`part_of` (umbrella↔component) and `extends` (specialization↔general).

**(d) Module theory — cohesion/coupling, ISP, SRP, separation of concerns [S20]–[S24].** High cohesion (one responsibility per unit [S20][S23]) + low coupling [S21] is the classic justification for splitting a low-cohesion monolith; ISP ("no client should depend on methods it does not use" [S22]) is the direct analogue of "no agent should load technology sub-topics its task does not use"; Dijkstra's separation of concerns [S24] is the primary root of the whole argument. R. C. Martin's package principles — Common Closure ("classes that change together are packaged together") and Common Reuse ("classes used together are packaged together") — are the co-change and reuse split-triggers in §3 (concept attributed; load-bearing citations are the verified cohesion/SRP pages [S20][S23]).

**(e) LLM context-window & token economics [S15][S16][S25].** "Lost in the Middle" empirically shows retrieval/answer quality is highest when relevant content is a small chunk at the start or end of context and drops when it sits mid-monolith [S15][S16]. Anthropic's context-engineering guidance makes the operational rule explicit: put the *smallest set of high-signal tokens* in context and disclose progressively [S25]. Fine-grained atomic retrieval is therefore not a nicety — it is the measured-optimal way to feed skills to an agent.

---

## 3. Granularity model

### 3.1 `kind` levels

A new first-class node attribute `kind` classifies every skill on the aggregation axis (orthogonal to `metadata.complexity`, which is a difficulty axis):

| `kind` | Meaning | `content` role | Typical out-edges |
|---|---|---|---|
| `atomic` | Indivisible building block — one cohesive responsibility, one version cadence, one validation surface. e.g. `gradle_kotlin_dsl`, `android.components.lifecycle` | The deep reference body (the leaf an agent actually loads to *do* the thing) | `requires` (prereqs), maybe `extends` (its general), advisory edges |
| `composite` | Intermediate aggregator that is itself `part_of`/`extends` a higher node. e.g. `android.aosp.build_system` (the platform build, part of the Android world but not an app component) | Orientation for its sub-area + a component manifest | `composes` (its parts), `requires`, `extends`, advisory |
| `umbrella` | Top-level technology/stack aggregator and **wizard entry point** (a whole "technology" a user types). e.g. `android.overview` | Thin map: "what this technology is + which component to load for which task" (progressive-disclosure top tier) | `composes` (its atomic/composite parts), `requires`, advisory |

`composite` and `umbrella` are the same *machinery* (both aggregate via `composes`); they differ only in **role/altitude**: `umbrella` is a user-selectable technology root, `composite` is an internal mid-level grouping. The distinction is a flag so the wizard knows which nodes are top-level selectable (§7) without hard-coding names.

**Default & backward-compatibility:** `kind` defaults to `atomic`. The 8 existing seed TOMLs omit `kind` and therefore remain valid unchanged (they load as `atomic`); the §6 worked example reclassifies `android.overview`→`umbrella` and `android.aosp.build_system`→`composite` as an additive migration.

### 3.2 When to split — testable criteria

**SPLIT** a would-be node into atomic children under an umbrella when **ANY** trigger fires (each is mechanically checkable):

- **S1 — Size.** Rendered `content` would exceed **~2,000 tokens** (≈ one progressive-disclosure page) **or ~350 Markdown lines**. Rationale: [S15][S25] — keep the loadable unit small.
- **S2 — Reuse across ≥2 parents.** The sub-topic is (or is planned to be) a `requires`/`composes` target of **≥2 distinct** parent skills (Common Reuse; DRY) [S20]. e.g. `gradle` is needed by `android.overview`, `kotlin_multiplatform`, and `spring_boot`.
- **S3 — Independent version cadence.** The sub-topic ships/releases on a **different clock** than the parent (Common Closure) [S11][S12]. e.g. Jetpack Compose BOM vs Android API level.
- **S4 — Distinct validation.** The sub-topic has its **own runnable validation / `evidences` signature** separable from its siblings' — separable proof ⇒ separable node.
- **S5 — Low cohesion.** `content` covers **≥2 responsibilities that do not change together** (SRP/ISP violation) [S22][S23].

**KEEP ATOMIC** (do *not* split) when **ALL** hold: **K1** single cohesive responsibility; **K2** one version cadence; **K3** one validation surface; **K4** under the size budget; **K5** not reused by ≥2 parents. An atomic node with `composes` out-edges is a contradiction and MUST fail validation (see §5.4).

**PROMOTE to `umbrella`/`composite`** when a node **aggregates ≥2 components via `composes`** *and* its own `content` is orientation/manifest rather than deep leaf material. An umbrella whose `content` is itself a 3,000-token deep-dive is a mis-modelled monolith — move the depth into leaves.

---

## 4. Relationship taxonomy

### 4.1 Edge `type` enum (superset of the existing three)

The existing `DependencyType` (`requires | extends | recommends`) is **preserved verbatim** and extended. Canonical stored set:

| `type` | Semantics (direction: **source → target**) | Symmetry | Acyclic? | In hard closure? |
|---|---|---|---|---|
| `requires` | Hard **knowledge prerequisite**: understanding *source* assumes *target*. (existing) | asymmetric | **MUST** be acyclic | **Yes** |
| `composes` | Whole–part **aggregation**: umbrella/composite *source* is built from component *target*. (NEW) | asymmetric | **MUST** be acyclic | **Yes** |
| `extends` | **Specialization / is-a**: *source* is a specialization of the more-general *target*. (existing) | asymmetric | **MUST** be acyclic | **Yes** |
| `recommends` | **Soft, asymmetric** suggestion: *target* complements *source*, no hard need. (existing) | asymmetric | exempt (advisory) | No |
| `related_to` | **Symmetric associative** "see also" — non-hierarchical topical link. (NEW) | symmetric | exempt (may cycle by nature) | No |
| `alternative_to` | **Symmetric substitute / mutual-exclusion** — *source* and *target* are competing choices for the same job (e.g. `maven ↔ gradle`). (NEW) | symmetric | exempt (may cycle by nature) | No |

**Aliases normalized at import** (accept the operator's vocabulary, store one canonical type — keeps the DB `CHECK` set small and semantics unambiguous):

| Author writes… | Stored as… | Note |
|---|---|---|
| `depends_on` | `requires` | package-manager vocabulary → knowledge-prerequisite edge |
| `prerequisite` | `requires` | plain-language synonym |
| `part_of` (authored on the **child**, pointing at the parent) | `composes` (**inverted** to parent→child) | one canonical whole→part edge regardless of authoring side |

**Direction & dedup rule for `composes`/`part_of`.** Canonical storage is always **umbrella `--composes-->` component** (matching the TOML convention that `[skill.dependencies]` lists the *targets the source points at*). A child may instead author `part_of = ["umbrella"]`; the importer inverts it to a `composes` edge on the umbrella. If both an umbrella's `composes` list and a child's `part_of` name the same pair, they **collapse to one edge** (idempotent import).

### 4.2 What the recursive-CTE resolver walks for "give me everything needed for X"

The **hard closure set = `{requires, composes, extends}`** (asymmetric, acyclic, structural). The resolver for "everything needed for X" transitively follows exactly these three from each root:

- `requires` → pull prerequisites (you can't learn the source without them),
- `composes` → pull the umbrella's components (the parts *are* the technology),
- `extends` → pull the general parent (a specialization presumes its general; Liskov).

The **advisory set = `{recommends, related_to, alternative_to}`** is collected into a *separate* "see-also / alternatives" result but is **never auto-pulled** into the required closure. `alternative_to` additionally powers conflict detection (§7): if the resolved root set contains two nodes linked `alternative_to`, the wizard surfaces a pick-one.

### 4.3 Acyclicity invariant (refines `seed/validate_dag.py`)

The current validator checks that `requires + extends + recommends` are *jointly* acyclic. The new model **scopes the acyclicity invariant to the hard closure set** `{requires, composes, extends}` and **excludes the advisory set**, because `related_to`/`alternative_to` are symmetric by definition (storing both directions is intentional and would false-positive a cycle). This is a *relaxation* — the 8 seed TOMLs (whose `recommends` edges happen to be acyclic) remain valid; nothing that passed before fails now. Concretely, `validate_dag.py` should: (1) also resolve `composes`/`part_of` targets against the closed-world corpus; (2) run the cycle check over `{requires, composes, extends}` only; (3) assert that no `atomic` node has an outgoing `composes` edge.

---

## 5. Data-model addendum (drop-in for SPEC §4/§5)

Every change below is **additive** and preserves the 8 seed TOMLs.

### 5.1 TOML format additions (§4.2)

New optional keys — omitting them yields the pre-existing behaviour:

```toml
[skill]
name    = "android.overview"
version = "0.2.0"
title   = "Android Application Platform (umbrella)"
kind    = "umbrella"          # NEW — atomic (default) | composite | umbrella
description = "..."
content = """..."""

[skill.metadata]
tags = ["android", "mobile"]
domain = "android"
complexity = "beginner"

[skill.dependencies]
requires        = ["java.language", "kotlin.language"]   # existing
extends         = []                                     # existing
recommends      = []                                     # existing
composes        = [                                      # NEW — component leaves this node aggregates
  "android.components.lifecycle",
  "android.ui.compose",
  "android.build.gradle_agp",
  "android.runtime.art",
  "android.distribution.play",
  "android.ndk.native",
]
related_to      = []                                     # NEW — symmetric "see also"
alternative_to  = []                                     # NEW — symmetric substitute
# depends_on / prerequisite / part_of are also ACCEPTED here and normalized per §4.1

# OPTIONAL ergonomic authoring form for umbrellas that need per-component
# ordering/optionality. Each entry materializes as one `composes` edge.
[[skill.components]]
name  = "android.components.lifecycle"
order = 1
optional = false

[[skill.components]]
name  = "android.ndk.native"
order = 6
optional = true          # only pulled when native code is in scope
```

`composes = [...]` (a simple name list, symmetric with `requires`) and `[[skill.components]]` (array-of-tables carrying `order`/`optional`/`note`) are **two authoring forms for the same edge**; the importer emits one `composes` edge per entry. Umbrellas needing ordering use `[[skill.components]]`; simple ones use the list.

### 5.2 Go model additions (§4.1)

```go
// NEW — aggregation-axis classification (orthogonal to complexity)
type SkillKind string
const (
    SkillKindAtomic    SkillKind = "atomic"    // indivisible building block (default)
    SkillKindComposite SkillKind = "composite" // mid-level aggregator
    SkillKindUmbrella  SkillKind = "umbrella"  // technology/stack root; wizard entry point
)

// Skill gains one field (all others unchanged)
type Skill struct {
    // ...existing fields...
    Kind SkillKind `json:"kind" db:"kind"` // default "atomic"
}

// DependencyType gains three canonical values (existing three preserved)
const (
    DepTypeRequires   DependencyType = "requires"   // existing — hard closure
    DepTypeExtends    DependencyType = "extends"    // existing — hard closure
    DepTypeRecommends DependencyType = "recommends" // existing — advisory
    DepTypeComposes   DependencyType = "composes"   // NEW      — hard closure (whole→part)
    DepTypeRelatedTo  DependencyType = "related_to" // NEW      — advisory, symmetric
    DepTypeAlternative DependencyType = "alternative_to" // NEW — advisory, symmetric
)

// HardClosureTypes is the set the "everything needed for X" resolver walks.
var HardClosureTypes = []DependencyType{DepTypeRequires, DepTypeComposes, DepTypeExtends}

// SkillDependency gains two optional edge attributes (for composes ordering)
type SkillDependency struct {
    SkillID      uuid.UUID      `json:"skill_id" db:"skill_id"`
    DependsOn    uuid.UUID      `json:"depends_on" db:"depends_on"`
    RelationType DependencyType `json:"relation_type" db:"relation_type"`
    Optional     bool           `json:"optional" db:"optional"`       // NEW — default false
    SortOrder    *int           `json:"sort_order,omitempty" db:"sort_order"` // NEW — component ordering
}
```

Import-time normalization: `depends_on`,`prerequisite` → `requires`; `part_of` (child→parent) → `composes` (parent→child, inverted). Symmetric `related_to`/`alternative_to` are stored once and treated as undirected by the resolver.

### 5.3 SQL schema additions (§5) — fresh-schema form

```sql
-- skills: add the aggregation-axis column (default keeps existing rows valid)
ALTER TABLE skills
    ADD COLUMN kind TEXT NOT NULL DEFAULT 'atomic'
    CHECK (kind IN ('atomic', 'composite', 'umbrella'));
CREATE INDEX idx_skills_kind ON skills(kind);

-- skill_dependencies: widen the relation_type CHECK and add edge attributes
ALTER TABLE skill_dependencies
    DROP CONSTRAINT skill_dependencies_relation_type_check;      -- name per Postgres default
ALTER TABLE skill_dependencies
    ADD CONSTRAINT skill_dependencies_relation_type_check
    CHECK (relation_type IN (
        'requires', 'extends', 'recommends',      -- existing
        'composes', 'related_to', 'alternative_to' -- NEW
    ));
ALTER TABLE skill_dependencies
    ADD COLUMN optional   BOOLEAN NOT NULL DEFAULT FALSE,  -- NEW
    ADD COLUMN sort_order INT;                             -- NEW (nullable)
```

For a greenfield `001_initial.up.sql` the same lands inline: `kind TEXT NOT NULL DEFAULT 'atomic' CHECK (kind IN ('atomic','composite','umbrella'))` on `skills`; the widened `CHECK (relation_type IN ('requires','extends','recommends','composes','related_to','alternative_to'))` plus `optional BOOLEAN NOT NULL DEFAULT FALSE` and `sort_order INT` on `skill_dependencies`. The `PRIMARY KEY (skill_id, depends_on)` **must widen to `(skill_id, depends_on, relation_type)`** so a pair may carry, e.g., both a `requires` and a historical `recommends` without collision.

### 5.4 New validation invariants (mechanical, anti-bluff)

1. **Closure acyclicity:** the sub-graph over `{requires, composes, extends}` is a DAG (recursive-CTE / DFS cycle check; §4.3).
2. **Atomic-has-no-parts:** `kind='atomic'` ⇒ zero outgoing `composes` edges.
3. **Umbrella-has-parts:** `kind IN ('umbrella','composite')` ⇒ ≥1 outgoing `composes` edge (else it is a mislabelled atomic).
4. **Closed-world targets:** every `composes` target resolves to a declared skill (same rule the corpus already enforces for `requires`/`extends`).
5. **Symmetric-edge canonicalization:** `related_to`/`alternative_to` are stored once; a duplicate reverse row is a redundant-edge warning, not a cycle.

---

## 6. Worked example — decomposing `android`

`android.overview` becomes a thin **umbrella** aggregating six real Android sub-skills (five mandatory + one optional native bridge). `android.aosp.build_system` stays a **composite** that `extends` the umbrella (the *platform* build is a specialization-context of the Android world, **not** an app component — hence `extends`, not `part_of`; this is the key modelling distinction §4 draws).

### 6.1 Umbrella — `android.overview` (post-decomposition)

```toml
[skill]
name    = "android.overview"
version = "0.2.0"
title   = "Android Application Platform (umbrella)"
kind    = "umbrella"
description = "Umbrella entry point for Android app development: orientation + a manifest of the atomic component skills (framework components, Compose UI, Gradle/AGP build, ART runtime, distribution, NDK). Load a component leaf for the specific task rather than this whole node."
content = """
# Android (umbrella)

## What this is
Google's Linux-kernel-based mobile OS and app platform. Apps run on the Android
Runtime (ART), are written in Kotlin/Java, and ship as APKs or App Bundles (AABs).

## Which component skill to load
- Build/package an app -> `android.build.gradle_agp`
- Write UI -> `android.ui.compose`
- Understand app structure & lifecycle -> `android.components.lifecycle`
- Bytecode/runtime questions -> `android.runtime.art`
- Ship to users -> `android.distribution.play`
- Bundle C/C++ native code -> `android.ndk.native` (optional)

## Boundary
This umbrella is app-level. The *platform* source build (compiling AOSP itself)
is a separate specialization: `android.aosp.build_system` (extends this node).
"""

[skill.metadata]
tags = ["android", "mobile", "umbrella"]
domain = "android"
complexity = "beginner"

[skill.dependencies]
requires   = ["java.language", "kotlin.language"]   # unchanged from seed
extends    = []
recommends = []
composes   = [
  "android.components.lifecycle",
  "android.ui.compose",
  "android.build.gradle_agp",
  "android.runtime.art",
  "android.distribution.play",
  "android.ndk.native",
]

[[skill.components]]
name = "android.components.lifecycle"
order = 1
optional = false

[[skill.components]]
name = "android.ui.compose"
order = 2
optional = false

[[skill.components]]
name = "android.build.gradle_agp"
order = 3
optional = false

[[skill.components]]
name = "android.runtime.art"
order = 4
optional = false

[[skill.components]]
name = "android.distribution.play"
order = 5
optional = false

[[skill.components]]
name = "android.ndk.native"
order = 6
optional = true

[[skill.resources]]
url = "https://developer.android.com/guide"
title = "Android Developer Guides"
resource_type = "official-doc"
```

The six components and their real Android substance:

| Atomic child | Responsibility | Own version cadence (S3) | Reuse (S2) |
|---|---|---|---|
| `android.components.lifecycle` | Activities/Fragments, Services, BroadcastReceivers, ContentProviders, `AndroidManifest.xml`, lifecycle callbacks | platform API level | flutter/react-native embedders |
| `android.ui.compose` | Jetpack Compose `@Composable`, recomposition, state hoisting | **Compose BOM** (independent) | compose-multiplatform |
| `android.build.gradle_agp` | Gradle + Android Gradle Plugin, D8/R8, DEX, shrink/obfuscate | **AGP** (independent) | any AGP-built module |
| `android.runtime.art` | ART, DEX bytecode, D8/R8 lowering, ahead-of-time/JIT | platform | kotlin (DEX target) |
| `android.distribution.play` | APK vs AAB packaging, Play Console, signing, sideload | Play policy | any shipped app |
| `android.ndk.native` | NDK, `externalNativeBuild`, JNI bridge to C/C++ | **NDK** (independent) | qt-on-android, game engines |

### 6.2 Atomic child — `android.ui.compose` (full TOML)

Authored with the **inverse `part_of` form** to demonstrate alias normalization (§4.1): the child names its umbrella, and the importer materializes the single canonical `android.overview --composes--> android.ui.compose` edge (idempotent with the umbrella's own `composes` list).

```toml
[skill]
name    = "android.ui.compose"
version = "0.1.0"
title   = "Jetpack Compose — Android Declarative UI"
kind    = "atomic"
description = "Reference for Jetpack Compose: @Composable functions, recomposition, state hoisting, and the Compose BOM version model — Android's recommended Kotlin-first UI toolkit."
content = """
# Jetpack Compose

## Overview
Jetpack Compose is Android's recommended, Kotlin-first declarative UI toolkit.
UI is described as composable functions rather than inflated from XML layouts.

## Key facts
- `@Composable` functions emit UI; the runtime re-invokes them (**recomposition**)
  when the state they read changes.
- State is held in observable holders (`mutableStateOf`, `remember`) and pushed
  down via **state hoisting** so composables stay stateless/testable.
- `Modifier` chains express layout, drawing, and input decoration compositionally.
- Versioning: Compose libraries are aligned via the **Compose BOM**, which moves
  on its own cadence independently of the platform API level — the reason this is
  a separate atomic skill (independent version clock).
- Requires Kotlin: Compose ships a Kotlin compiler plugin and leans on Kotlin
  language features (trailing lambdas, default args); it is not available to Java.
- Interoperates with the legacy View system via `ComposeView` / `AndroidView`.

## Validation surface
Rendered-UI proof (screenshot/host-render of a composable tree) — distinct from a
Gradle build success or a lifecycle-callback probe.
"""

[skill.metadata]
tags = ["android", "compose", "ui", "kotlin", "declarative"]
domain = "android"
complexity = "intermediate"

[skill.dependencies]
requires = ["kotlin.language"]          # Compose is Kotlin-only
extends  = []
recommends = []
part_of  = ["android.overview"]          # normalized -> composes edge on android.overview

[[skill.resources]]
url = "https://developer.android.com/compose"
title = "Jetpack Compose (Android Developers)"
resource_type = "official-doc"
```

*(Resource URLs in §6 illustrate the format; the system's existing resource-validation path is what verifies them at ingest — they are not part of this document's own §2 evidence set.)*

---

## 7. Wizard mapping — create → map → process

User enters a stack, e.g. the R6 canonical set: **`android, android_aosp, java, kotlin, c++, cmake`**. The pipeline turns free-text tokens into a fully-expanded, de-duplicated, topologically-ordered **atomic building-block closure**.

### Stage 1 — CREATE / resolve (token → canonical skill)
Each token is normalized to a canonical `name` via the existing name index + `POST /search` (hybrid vector/keyword) for fuzzy/alias hits:

| Token | Resolved `name` | `kind` |
|---|---|---|
| `android` | `android.overview` | umbrella |
| `android_aosp` | `android.aosp.build_system` | composite |
| `java` | `java.language` | atomic |
| `kotlin` | `kotlin.language` | atomic |
| `c++` | `cpp.language` | atomic |
| `cmake` | `cmake.build_system` | atomic |

Unresolved tokens → gap report (`POST /expand/gap-report`) rather than a silent drop.

### Stage 2 — MAP (classify roots by `kind`)
`umbrella`/`composite` nodes are **expansion roots** (they will fan out via `composes`/`extends`); `atomic` nodes are leaves that still fan out via `requires`. All six become roots of the closure walk.

### Stage 3 — PROCESS / expand (hard-closure over `{requires, composes, extends}`)
Run the recursive CTE from all roots, following only the hard closure set, de-duplicating shared leaves:

```sql
WITH RECURSIVE closure AS (
    -- seed: the wizard's resolved roots
    SELECT s.id, s.name, s.kind, 0 AS depth
    FROM skills s
    WHERE s.name = ANY($1::text[])            -- e.g. {android.overview, android.aosp.build_system, java.language, kotlin.language, cpp.language, cmake.build_system}
  UNION
    -- walk ONLY the hard closure edge set
    SELECT dep.id, dep.name, dep.kind, c.depth + 1
    FROM closure c
    JOIN skill_dependencies e ON e.skill_id = c.id
    JOIN skills dep           ON dep.id = e.depends_on
    WHERE e.relation_type IN ('requires', 'composes', 'extends')
      AND (e.optional = FALSE OR $2::bool)     -- $2 = include_optional_components
)
SELECT DISTINCT name, kind, MAX(depth) AS max_depth
FROM closure
GROUP BY name, kind
ORDER BY max_depth DESC;                       -- deepest prereqs first ≈ topological order
```

`SELECT ... WHERE kind = 'atomic'` yields the **pure building-block set** an agent loads; the `umbrella`/`composite` rows are the cheap "table of contents" surfaced first (progressive disclosure [S25]).

**Closure trace for the R6 stack** (hard edges only; `android.overview` decomposed per §6, native component optional):

- **Roots pulled directly:** `android.overview` (umbrella), `android.aosp.build_system` (composite).
- **`android.overview` `composes`** → `android.components.lifecycle`, `android.ui.compose`, `android.build.gradle_agp`, `android.runtime.art`, `android.distribution.play` (+ `android.ndk.native` iff `include_optional`).
- **`android.overview` `requires`** → `java.language`, `kotlin.language` (already roots → de-duped once).
- **`android.aosp.build_system` `extends`** → `android.overview` (already present); **`requires`** → `linux.os`, `python.language`, `make.build_system`.
- **`kotlin.language` `requires`** → `java.language` (de-duped). **`cpp.language` `requires`** → `c.language`. **`cmake.build_system` `requires`** → `make.build_system`, `cpp.language` (de-duped). **`make.build_system` `requires`** → `c.language` (de-duped).

**Topologically-ordered atomic closure (prereqs first):**
`c.language`, `java.language`, `python.language`, `linux.os` → `make.build_system`, `cpp.language`, `kotlin.language` → `cmake.build_system`, `android.runtime.art`, `android.components.lifecycle`, `android.ui.compose`, `android.build.gradle_agp`, `android.distribution.play` (→ `android.ndk.native` if opted-in) → then the aggregator ToC nodes `android.overview`, `android.aosp.build_system`.

**Advisory / see-also (collected, NOT auto-pulled):** `bazel.build_system` (recommended by AOSP), `design_patterns.overview` (recommended by java/cpp). **Alternatives / conflicts:** if the user's tokens also resolved to two `alternative_to`-linked nodes (e.g. `maven` *and* `gradle`), the wizard surfaces a pick-one instead of loading both.

### Why this shape is correct
The umbrella gives the agent a cheap map; the `composes`/`requires`/`extends` closure gives the exact, de-duplicated, ordered set of atomic skills the stack needs; the advisory edges are offered but never bloat context. That is package-manager transitive resolution [S8][S12] fused with Agent-Skills progressive disclosure [S1][S25], expressed over the DAG the system already has.

---

## 8. Backward-compatibility & migration checklist

- **Seed TOMLs (×8):** load unchanged — `kind` defaults to `atomic`; no new required keys. ✅
- **Existing edges:** `requires`/`extends`/`recommends` semantics and storage unchanged; only new `type` values and two nullable/defaulted edge columns added. ✅
- **PK widening:** `skill_dependencies` PK → `(skill_id, depends_on, relation_type)` (a strict superset that still rejects true duplicates). Migration must de-dup any pre-existing exact-duplicate rows first (the seed has none). ✅
- **Validator:** update `seed/validate_dag.py` to (1) resolve `composes`/`part_of` targets, (2) scope the cycle check to `{requires, composes, extends}`, (3) assert atomic-has-no-`composes` (§5.4). Relaxation-only for the current corpus. ✅
- **Reclassification migration (worked example):** bump `android.overview`→`umbrella`, author the six component leaves, add `composes` edges, keep `android.aosp.build_system extends android.overview`. Additive; no node deleted. ✅
- **API/MCP:** `GET /skills/:name/tree` and `skill_tree` gain a `closure=hard|all` (or `edge_types=`) parameter to select the §4.2 hard set vs. everything; default `hard` reproduces "everything needed for X". Existing callers unaffected if `hard` is the default. ✅

---

## 9. Consolidated verified sources

All reachability-checked 2026-07-15 (status lines in §2). Anthropic Agent Skills / progressive disclosure [S1][S2][S3][S25]; package-manager composition & dependency graphs — npm [S4][S5], Cargo [S6][S7], Maven [S8][S9], Gradle [S10][S11], Go [S12][S13]; ontology relation typing — SKOS [S17][S18], OBO RO [S19]; module theory — cohesion [S20], coupling [S21], ISP [S22], SRP [S23], separation of concerns [S24]; LLM context/token economics — Lost-in-the-Middle [S15][S16], context engineering [S25].
