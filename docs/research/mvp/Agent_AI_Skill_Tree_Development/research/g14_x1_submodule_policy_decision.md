# G14 / X1 — Submodule Dependency-Layout Policy: Operator Decision Package

**Revision:** 2
**Last modified:** 2026-07-15T17:25:00Z
**Status:** decision-package for operator (§11.4.66) — no code written, no file modified, no git operation performed by this document's author.
**Authority cited:** `constitution/CLAUDE.md` §11.4.28 (esp. §11.4.28(C)), §11.4.29, §11.4.31, §11.4.36, §11.4.76, §11.4.161, §11.4.177; project docs `REQUIREMENTS.md` (R7, R9, R12, R19, R20), `IMPLEMENTATION_PLAN.md` (X1, global constraint #2, P3/P7/P8/P9/P10/P11/P12/P13), `GAPS_AND_RISKS_REGISTER.md` (G13, G14, Adjudication row #2, G28).

---

## 1. What this document is

This is a decision package, not a decision. It lays out the two candidate policies the
Gaps & Risks Register names for **G14**, states precisely how each satisfies (or strains)
the cited Constitution clauses, incorporates a new operator mandate (**R20** — mandatory
Containers-submodule routing) and freshly-verified real-world layout precedent from three
sibling `helix_*` projects, enumerates exactly what work is currently blocked on this
decision, and ends with a recommendation plus the exact question(s) the operator must
answer before any implementer touches `helix-deps.yaml` or `scripts/sync_submodules.sh`
in `--apply` mode. No file was modified and no git command was run to produce it — every
sibling-project fact below was gathered read-only (`cat`, `grep`, `git remote -v`,
`git log -1`, `gh repo view`), no `--apply`/`git submodule add`/write executed.

---

## 2. The conflict, precisely, with citations

Three project artifacts currently disagree with each other in wording, even though two of
the three are trying to describe the same intended behavior:

| Source | Exact wording | Reading |
|---|---|---|
| `project/helix-deps.yaml` header | *"§11.4.28C Dependencies resolve to ONE canonical location: `<repo_root>/<name>/` (ungrouped) or `<repo_root>/submodules/<name>/` (grouped). Nested own-org submodule chains are FORBIDDEN."* | Strict single-canonical (Option A). |
| `REQUIREMENTS.md` R9 (line 78-81) | *"All deps and dependency submodules live under `submodules/<snake_case>/`, OR reuse the versions from the parent dir/submodules — **parent-dir versions have PRIORITY**. **BOTH copies must always be in sync** with main/master."* | Two copies, both kept in sync (Option B). |
| `IMPLEMENTATION_PLAN.md` §2 global constraint #2 | *"**Single-canonical dependencies** (§11.4.28C): one location per dependency; parent-ecosystem copy has priority, else `submodules/<snake_case>/`."* | Reads as Option A but borrows Option B's "priority" language — the plan document itself hedges both ways. |
| `IMPLEMENTATION_PLAN.md` X1 (cross-cutting) | *"Flagged discrepancy: Constitution §11.4.28C mandates single-canonical; operator asked for parent-priority+both synced. **Escalation open** — resolve before any `--apply` run that would create a second copy."* | Names the conflict explicitly and blocks `--apply`. |
| `GAPS_AND_RISKS_REGISTER.md` G14 | *"an open governance escalation blocks P3/P7/P8/P10/P11/P13 dependency work."* | States the blast radius directly. |

**Constitutional fact established by this research (not previously stated in the register):**
all seven dependencies declared in `helix-deps.yaml` (`llms_verifier`, `helix_llm`,
`helix_agent`, `embeddings`, `helix_qa`, `challenges`, `docs_chain`) belong to orgs on the
§11.4.28(A) owned-org list (`vasic-digital`, `HelixDevelopment`) — none of them are
third-party. This matters because §11.4.28(C)'s "third-party submodules exempt" carve-out
does **not** apply to any of the seven; the strict single-canonical / no-nested-chain
language binds them in full, and the §11.4.28(C) depth-1 reusable-engine carve-out is
explicitly scoped to `constitution/submodules/<name>/` only — *"This applies to NO other
submodule."* So the G14 decision cannot lean on that carve-out; it has to be resolved on
its own terms. The same is true of the two new dependencies this round adds — **containers**
(`vasic-digital`) and **LLMProvider** (`vasic-digital`/`HelixDevelopment`, see §4) — both
owned-org. **open-design** (`nexu-io`) is the one dependency in this whole set that genuinely
qualifies for the third-party exemption.

**Independently verified, host-side fact (read directly, not assumed):** all seven
originally-declared dependencies are *already* present as real git submodule checkouts of
the sibling ecosystem project `helix_code`, confirmed by `.git` presence in every case:

```
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/llms_verifier   (.git present)
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/helix_llm       (.git present)
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/helix_agent     (.git present)
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/embeddings      (.git present)
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/helix_qa        (.git present)
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/challenges      (.git present)
/home/milos/Factory/projects/tools_and_research/helix_code/submodules/docs_chain      (.git present)
```

`/home/milos/Factory/projects/tools_and_research/helix_code/.gitmodules` registers each
one (e.g. `[submodule "challenges"] path = submodules/challenges`). Zero of the seven
appear in `project/go.mod` as a Go module `require` — confirmed by reading the full
`go.mod`. Combined with the research findings already in `REQUIREMENTS.md` ("HelixAgent
(:7061) and HelixLLM (:8443) are both OpenAI-compatible" over HTTP; "llms-verifier/helixqa
are PATH binaries"), this establishes as **fact** that none of the seven is consumed as a
Go source import — they are consumed as (a) HTTP services (`helix_llm`, `helix_agent`),
(b) PATH binaries built from source (`llms_verifier`, `helix_qa`, `docs_chain`), or (c) a
script/test corpus exercised against the running service (`challenges`). `embeddings` is
the one exception with its own importable Go module (`digital.vasic.embeddings`), read
directly from its `go.mod`, though it is not yet imported by skill-system. This is
material to "what actually needs vendoring vs. what could be reached over the network or
via `go get`" — see §9 (Honest gaps).

**Zero of the seven are vendored inside `skill-system` today** — `helix-deps.yaml`
declares them, `submodules/` does not exist under the project root, confirmed by directory
listing.

---

## 3. New evidence incorporated mid-research: sibling-ecosystem layout precedent

The coordinator supplied real `.gitmodules` layout evidence from three sibling `helix_*`
projects as precedent for the single-canonical-vs-parent-priority decision. Each claim was
independently re-verified by reading the actual `.gitmodules` files in full (not merely
grepped), so the table below states the confirmed reality, including two corrections to
the summary as originally supplied.

| Project | Constitution submodule | Layout convention observed |
|---|---|---|
| `helix_ota` | `HelixConstitution` at root `constitution` (`git@github.com:HelixDevelopment/HelixConstitution.git`) | **Mixed, and internally inconsistent.** `constitution`, `containers`, and `docs_chain` are registered at ROOT (ungrouped); `ota-*`, `http3`, `helixqa`, `challenges`, `doc_processor`, `llm_orchestrator`, `llm_provider`, `security`, `vision_engine`, `llms_verifier`, `website` are registered under `submodules/` (grouped). **Correction to the supplied summary:** `containers` is registered **twice** in this one file — once at root (`[submodule "containers"] path = containers`, line 4-6) and again under `submodules/containers` (line 37-39) — **both pointing at the identical URL** `git@github.com:vasic-digital/containers.git`. This is not "root, full stop" as summarized; it is a live, unresolved rival-registration inside a single project, i.e. direct evidence that this exact class of problem (two registered locations for one dependency) already exists unresolved in a sibling, not merely as a hypothetical this research is inventing. |
| `helix_track` | `constitution` at root, same URL | **Fully root/ungrouped for every single dependency** — all ~45 entries (`containers`, `docs_chain`, `challenges`, `helix_qa`, `helix_agent`, `llm_provider`, `embeddings`, `auth`, and ~35 more) are registered directly at the repo root with no `submodules/` prefix at all. **Correction to the supplied summary:** the client matrix (`web_client`, `harmony_os_client`, `aurora_os_client`, etc.) is real and root-based, but so is *every other* dependency in this project — the root-based convention is not specific to clients, it is this project's uniform house style. |
| `helix_terminator` | `constitution` at root, same URL | **Fully grouped under `submodules/` for every non-constitution dependency** — `open-design` → `submodules/open-design`, `auth` → `submodules/auth`, `llmprovider` → `submodules/llmprovider`, `containers` → `submodules/containers`, `helixllm` → `submodules/helixllm`, `helixtrack-core` → `submodules/helixtrack-core`. Confirms the supplied summary as stated. |
| `helix_code` *(the project this decision's parent-ecosystem-root already points at)* | `constitution` at root, same URL | **Fully grouped under `submodules/` for every dependency checked** — all 7 originally-declared deps plus `containers` (`submodules/containers`, `git@github.com:vasic-digital/Containers.git` — GitHub confirms `containers` lowercase is the canonical name; the capitalized form resolves to the same repo) and `llm_provider` (`submodules/llm_provider`, primary remote `github.com:HelixDevelopment/LLMProvider.git`, with an additional `vasicdigitalgithub` push remote to `github.com:vasic-digital/LLMProvider.git` — confirmed via `git remote -v`, i.e. a §2.1 multi-upstream-push pair, **not** two independent forks). |

**Conclusion from this table:** the operator's own ecosystem does **not** have one uniform
convention — `helix_track` is 100% root/ungrouped, `helix_terminator` is 100% grouped,
`helix_ota` is a live mix with one unresolved duplicate. What settles this for the
skill-graph project is not "which sibling looks most authoritative" but **which layout the
project's own already-resolved parent-ecosystem root (`helix_code`) actually uses** —
because that is the literal, load-bearing convention every one of this project's real
dependency resolutions passes through today. `helix_code` uses `submodules/<name>`
(grouped) uniformly, for all 9 dependencies checked (the original 7 + `containers` +
`llm_provider`). This matches the skill-graph project's own `helix-deps.yaml`, which
already declares `layout: grouped` for all seven current entries. Recommending anything
other than "continue grouped, matching `helix_code`" would mean deliberately diverging
from the one sibling convention that is actually in the resolution path, in favor of a
different sibling (`helix_track`) whose convention is not otherwise connected to this
project at all.

---

## 4. New operator mandate incorporated: R20 — Containers submodule is mandatory, not optional

**R20 (2026-07-15):** all containerization work in this project MUST route through the
`vasic-digital/containers` submodule's `pkg/boot` / `pkg/compose` / `pkg/health` layer,
never ad-hoc `podman`/`docker` invocations, per the pre-existing universal constitution
mandates §11.4.76 (Containers-submodule mandate) and §11.4.161 (rootless-container-runtime
mandate — Podman rootless only, no rootful docker, no sudo).

**Independently verified (not merely accepted from the mandate's own wording):**

- `vasic-digital/containers` exists as a real, public GitHub repository — confirmed live
  via `gh repo view vasic-digital/containers` returning `{"name":"containers","url":
  "https://github.com/vasic-digital/containers"}`.
- It is **already vendored** inside `helix_code` at `helix_code/submodules/containers`,
  confirmed by `.git` presence, and its remotes (`git remote -v`) confirm both GitHub
  (`vasic-digital/Containers.git`) and GitLab (`vasic-digital/containers.git`) mirrors are
  configured — a real, live, multi-upstream-pushed checkout at commit `a432efa8` (2026-07-05).
- Its own `helix-deps.yaml` (read in full) declares itself a **leaf**: *"Containers is a
  leaf Go submodule with ZERO own-org submodule dependencies"* — `deps: []`,
  `transitive_handling.recursive: true`, `transitive_handling.conflict_resolution:
  operator-required`, `language_specific_subtree: false`. This is the **full** §11.4.31
  schema (unlike the skill-graph project's own `helix-deps.yaml`, which is missing the
  `transitive_handling`/`language_specific_subtree` fields — noted in §9). Being a leaf
  means vendoring it introduces zero nested own-org chain risk under §11.4.28(C).
- Its `pkg/` directory (read directly) confirms the three named layers exist for real:
  `pkg/boot`, `pkg/compose`, `pkg/health` — alongside `pkg/discovery`, `pkg/lifecycle`,
  `pkg/orchestrator`, `pkg/runtime`, `pkg/monitor`, and others. The mandate's cited
  integration points are not aspirational — the package surface to route G13's compose
  work through is real and present.

**How R20 changes the G14 decision's dependency set:** `containers` is not a new,
unresolved dependency the way `open-design` or `LLMProvider` are (see §6/§7) — it is the
**tenth** dependency this research finds already vendored at the exact same
parent-ecosystem root (`helix_code`) that resolves the other seven. Under **Option A**,
adding `containers` to `helix-deps.yaml` resolves it immediately, with **zero** new
vendoring, by the exact same `find_ecosystem_copy()` mechanism already exercised for the
other seven. Under **Option B**, it becomes an eighth dependency requiring the not-yet-
built mirror-pin-and-guard mechanism (§6).

**Scope boundary (stated plainly, not glossed over):** *how* `internal/validation/sandbox.go`
(G02, real-isolation-vs-fail-closed-SKIP) and the ops-hardening compose work (G13 —
`project/deploy/docker-compose.yml` canonicalization) are rewritten to call
`containers`' `pkg/boot`/`pkg/compose`/`pkg/health` APIs instead of shelling out to
`podman`/`docker` directly is a **P12/P4-scoped engineering design**, not this G14/X1
layout decision. This document's scope is *where the containers submodule's source tree
lives and how it is kept in sync* (the G14 question, now including this dependency);
*how the codebase calls into it* is a follow-on task flagged in §9 and left to the P12/G13
design that already exists (`research/ops_hardening_design.md`).

---

## 5. Option A — §11.4.28(C) single-canonical

**Statement:** each dependency resolves to **exactly one** location on the whole host at
any given time. If a parent-ecosystem copy already exists, that copy *is* the canonical
location and `skill-system` vendors nothing locally for it. If no ecosystem copy exists,
`skill-system` vendors it once, locally, as a normal owned-org git submodule.

**This is the exact behavior `scripts/sync_submodules.sh` already implements today** (read
in full; `process_dep()` at lines 252-281): it checks `find_ecosystem_copy()` first (search
of `--ecosystem-root`/`HELIX_ECOSYSTEM_ROOTS`-supplied roots for `<root>/submodules/<name>`
or `<root>/<name>`), and only falls through to `add_new_submodule()` when no ecosystem
copy and no local copy exist. The two branches are mutually exclusive by construction —
there is never a code path that creates a second copy alongside a discovered ecosystem
copy. The script already supports multiple `--ecosystem-root` values (an array, looped in
`find_ecosystem_copy()`) or a colon-separated `HELIX_ECOSYSTEM_ROOTS` — relevant because
`open-design` is not present under `helix_code` but **is** present under the sibling
`helix_terminator` (§3, §7), so Option A can still resolve it to a single existing
ecosystem copy *if* `helix_terminator` is added as a second `--ecosystem-root`, without any
script change.

**Directory layout:**

```
# When --ecosystem-root /home/milos/Factory/projects/tools_and_research/helix_code is supplied
# (all 9 owned-org deps identified so far resolve here today; verified above):
skill-system/                       # this project's root
├── helix-deps.yaml                 # declares all deps, unchanged in shape
├── submodules/                     # DOES NOT gain entries for any ecosystem-resolved dep
└── (no local checkout of llms_verifier/helix_llm/helix_agent/embeddings/
     helix_qa/challenges/docs_chain/containers/llm_provider anywhere under this repo)

# The single canonical copy of each remains exactly where it is today:
helix_code/submodules/llms_verifier/   ← canonical
helix_code/submodules/helix_llm/       ← canonical
helix_code/submodules/helix_agent/     ← canonical
helix_code/submodules/embeddings/      ← canonical
helix_code/submodules/helix_qa/        ← canonical
helix_code/submodules/challenges/      ← canonical
helix_code/submodules/docs_chain/      ← canonical
helix_code/submodules/containers/      ← canonical (R20)
helix_code/submodules/llm_provider/    ← canonical (LLMProvider, if adopted — see §7)
```

If a future dependency is declared that has **no** ecosystem copy anywhere the resolver is
pointed at (e.g. `open-design` if `helix_terminator` is not added as a second ecosystem
root — see §7), it vendors exactly once at `skill-system/submodules/<name>/` and *that*
becomes canonical.

**§11.4.28 compliance:**
- **(A) Equal-codebase:** unaffected by this choice — whichever repo holds the canonical
  copy is worked on with the same engineering attention either way; (A) does not prefer
  either option.
- **(B) Decoupling:** satisfied identically under both options — the dependency source
  itself is never edited to know about `skill-system`; only the *resolver* differs.
- **(C) Dependency-layout / nested-chain-forbidden:** satisfied **literally and directly**
  — there is one location, full stop; no second copy exists to diverge or to be
  independently edited.

**§11.4.31 (helix-deps.yaml) compliance:** the manifest's existing header text ("ONE
canonical location") is **already true** under this option — no manifest edit needed
beyond adding the new entries (§7).

**§11.4.36 (install_upstreams on clone/add) compliance:** only fires for the fallback
branch (no ecosystem copy found → `add_new_submodule` clones fresh). For the 9 deps
confirmed at `helix_code` this branch never executes, so `install_upstreams` is never
invoked *by skill-system* for them — it remains the responsibility of whichever project
originally cloned them (`helix_code`, confirmed already holding `upstreams/`-bearing
checkouts for `embeddings` and `containers`, read directly in both cases).

**Sync/pinning mechanism:** `sync_submodules.sh --apply` fetches+fast-forwards the *one*
existing copy (`sync_existing_checkout()`) to the ref declared in `helix-deps.yaml`. There
is no cross-copy pinning problem because there is no second copy to keep in step.

**Failure modes (identified by tracing the actual code path, not assumed):**
1. **Ecosystem-root discovery is a per-invocation CLI flag / env var, not a committed
   value.** If `--ecosystem-root` / `HELIX_ECOSYSTEM_ROOTS` is absent or wrong on a given
   machine or in CI, the *same* dependency set resolves differently on different hosts —
   one host finds the ecosystem copy and vendors nothing, another finds nothing and
   vendors a fresh local copy. Both individual runs are internally consistent (never two
   copies on one host), but the *fleet* can end up with divergent per-host topologies
   unless the ecosystem-root value(s) are captured somewhere committed — this is a
   **documentation/wiring obligation**, not a defect in the option itself, and it grows
   with this round's finding that a *second* root (`helix_terminator`) may be needed for
   `open-design`.
2. **A fresh clone of `skill-system` alone (no sibling checkouts, no `--ecosystem-root`)
   has zero local copies of any declared dependency** until an operator either (a) clones
   the relevant sibling ecosystem repo(s) and supplies their paths, or (b) runs
   `sync_submodules.sh --apply` with no ecosystem root, which then vendors fresh local
   copies. This is real and must be stated plainly in onboarding docs — it is the direct
   cost of "logical, not physical" single-canonical.
3. **No commit-SHA record inside `skill-system`'s own git history for the ecosystem-
   resolved case.** Reproducibility of "what commit of `llms_verifier` (or `containers`)
   was this build tested against" lives entirely in the sibling project's own
   `.gitmodules`/submodule pointer, not in anything `skill-system` commits — a person
   auditing `skill-system` alone cannot answer that question from its own git log.
4. **`helix_ota`'s live duplicate `containers` registration (§3) is a concrete, observed
   instance of exactly the failure Option A is designed to prevent** — it is evidence that
   without the mutually-exclusive resolver logic `sync_submodules.sh` already has, this
   exact ecosystem drifts into rival copies given the chance.

---

## 6. Option B — operator parent-priority + both-synced

**Statement:** the parent-ecosystem copy is the one **logical** canonical (edits, upstream
pulls, and version bumps happen there). `skill-system` **additionally always vendors its
own copy**, as a real, committed git submodule at `skill-system/submodules/<name>/`. That
local copy is a **read-only mechanical mirror**: `sync_submodules.sh` (extended) pins it to
exactly the commit the ecosystem copy is at, and nothing else is permitted to write inside
it directly.

**Directory layout:**

```
skill-system/
├── helix-deps.yaml
├── .gitmodules                      # NEW: one entry per dependency added
└── submodules/
    ├── llms_verifier/                # real git submodule, pointer pinned == ecosystem copy's HEAD
    ├── helix_llm/                    # same
    ├── helix_agent/                  # same
    ├── embeddings/                   # same
    ├── helix_qa/                     # same
    ├── challenges/                   # same
    ├── docs_chain/                   # same
    ├── containers/                   # same (R20)
    └── llm_provider/                 # same, if adopted (§7)

# unchanged, remains the place where real edits/upstream pulls happen:
helix_code/submodules/llms_verifier/   ← logical canonical (edit here)
helix_code/submodules/containers/      ← logical canonical (edit here)
... (same for every other dependency)
```

**§11.4.28 compliance:**
- **(A) Equal-codebase:** unaffected, same as Option A.
- **(B) Decoupling:** unaffected, same as Option A — the *dependency's own source* is
  never made project-aware under either option; only `skill-system`'s local mirror
  existing or not existing changes.
- **(C) Dependency-layout / nested-chain-forbidden:** satisfied **only under an added
  constraint that is not automatic** — the letter of §11.4.28(C) says "ONE canonical
  location"; two independently-existing git checkouts of the same owned-org repository
  is, by the plain reading of that sentence, two locations. The register's own framing
  ("this satisfies §11.4.28C's 'one canonical' *semantically* … while honouring the
  operator's 'parent priority + both in sync' *operationally*") is a **reconciliation
  argument, not a literal reading** — it holds only if the local copy is mechanically
  prevented from ever becoming an independently-editable second source of truth. That
  requires a **new enforcement component** (a guard rejecting any commit that touches
  `submodules/<name>/**` unless it is exactly the pointer-bump produced by the sync
  script) that does not exist today. Without that guard, Option B is not "one canonical
  location with a mirror" — it is two copies that can silently diverge, which is the
  precise anti-pattern §11.4.28(C) exists to forbid, and which `helix_ota`'s own live
  duplicate `containers` registration (§3) shows already happens in this ecosystem absent
  such a guard.

**§11.4.31 (helix-deps.yaml) compliance:** the manifest's own header currently asserts
strict single-canonical wording. Adopting Option B **requires editing that header** so the
manifest stops contradicting the mechanism it drives.

**§11.4.36 (install_upstreams on clone/add) compliance:** Option B's local-mirror vendor
step **is** a `git submodule add` (`add_new_submodule()` runs unconditionally for every
dependency, since none exists locally yet) — this is a **new, real obligation** Option A
does not have today: the first `--apply` run under Option B must invoke `install_upstreams`
against each newly-cloned local mirror (`embeddings` and `containers` are both confirmed
today to carry `upstreams/` directories at their ecosystem-copy locations; the remaining
seven were not individually re-verified for `upstreams/` content in this pass and should be
checked at implementation time, not assumed).

**Sync/pinning mechanism (new — does not exist in the script today):** `sync_submodules.sh`
would need a new mode: fetch/read the ecosystem copy's current HEAD, then fast-forward /
checkout the local mirror to that exact SHA, then commit the resulting submodule-pointer
bump in `skill-system` as its own tracked git object. This is materially more code than
Option A requires — Option A's `process_dep()` needs no changes at all; Option B needs (i)
a "read peer HEAD" step, (ii) a "pin mirror to that SHA" step, (iii) a write-guard against
direct edits inside the mirror, and (iv) a manifest-header rewrite.

**Failure modes:**
1. **Drift-by-omission**, now with a *live, observed precedent* — `helix_ota` (§3) already
   carries two registrations of `containers` pointed at the same URL inside one project;
   this is the concrete shape of what happens when a "should be the same thing, mirrored"
   pair of paths exists without a guard actively preventing divergence.
2. **Doubled review/diff surface:** every submodule-pointer bump in the mirror produces its
   own commit + diff noise in `skill-system`'s history, on top of whatever churn already
   exists in the ecosystem copy's own history.
3. **New engineering surface = new defect surface.** The guard/mirror-pin logic described
   above does not exist yet; every line of it is untested new code before `--apply` can be
   trusted, whereas Option A is provably already correct (it is the code currently in the
   repo, already exercised by `process_dep()`'s existing branching), and is furthermore the
   *same* convention already uniformly used by the actual resolved parent (`helix_code`,
   §3) for every one of the 9 dependencies checked.

---

## 7. Concrete dependency list + canonical path under each option

| Dependency | Owning org | Already at `helix_code` (the currently-configured ecosystem root)? | Canonical path — Option A | Canonical path — Option B |
|---|---|---|---|---|
| `llms_verifier` | vasic-digital | Yes (`.git` confirmed) | `helix_code/submodules/llms_verifier` only | same, **+** `skill-system/submodules/llms_verifier/` mirror |
| `helix_llm` | HelixDevelopment | Yes | `helix_code/submodules/helix_llm` only | same **+** mirror |
| `helix_agent` | HelixDevelopment | Yes | `helix_code/submodules/helix_agent` only | same **+** mirror |
| `embeddings` | vasic-digital | Yes | `helix_code/submodules/embeddings` only | same **+** mirror |
| `helix_qa` | HelixDevelopment | Yes | `helix_code/submodules/helix_qa` only | same **+** mirror |
| `challenges` | vasic-digital | Yes | `helix_code/submodules/challenges` only | same **+** mirror |
| `docs_chain` | vasic-digital | Yes | `helix_code/submodules/docs_chain` only | same **+** mirror |
| `containers` (**R20**, new this round) | vasic-digital | Yes (`.git` confirmed, remotes verified, leaf `helix-deps.yaml` confirmed) | `helix_code/submodules/containers` only — **zero new vendoring**, resolves through the existing mechanism the instant it is added to `helix-deps.yaml` | same **+** mirror at `skill-system/submodules/containers/` |
| `llm_provider` (LLMProvider — R7/R19, **not yet in `helix-deps.yaml`**, flagged "under separate evaluation") | HelixDevelopment (primary remote) **+** vasic-digital (confirmed additional push mirror, same content — verified via `git remote -v`, not two independent repos) | Yes (`helix_code/submodules/llm_provider`, `.git` confirmed) | `helix_code/submodules/llm_provider` only, if/when adopted | same **+** mirror |
| `open-design` (**R12**) | nexu-io (third-party — the one dependency in this set genuinely exempt from §11.4.28(C)'s nested-chain clause) | **No** — confirmed absent from `helix_code/.gitmodules` and `helix_code/submodules/`. **Present at a different sibling**, `helix_terminator/submodules/open-design`. | Resolves only if `helix_terminator` is added as a **second** `--ecosystem-root`; otherwise vendors fresh, once, at `skill-system/submodules/open_design/` per `research/opendesign_incorporation.md` §3.1 | Same choice applies; if mirrored, mirror path is `skill-system/submodules/open_design/` regardless of which ecosystem root supplies the logical canonical |

**Net verified finding:** for the eight owned-org dependencies now identified (seven
original + `containers`), **Option A unblocks all eight immediately with zero new
engineering** (the resolver already implements it, and all eight are already present at
the one ecosystem root already configured); **Option B unblocks all eight only after new
engineering lands** (mirror-pin mode + write-guard + manifest-header edit), described in
§6. `open-design` and (if adopted) `llm_provider` sit slightly outside that clean picture
and are called out individually above rather than folded into the "all resolve identically"
claim, because their ecosystem-presence facts differ from the original seven.

---

## 8. Blast radius — what this decision unblocks

Everything below is currently named, in the cited documents, as blocked or pending on this
exact decision, or is newly identified as blocked by this round's R20/sibling-precedent
research. Nothing here is inferred beyond what the citation states.

| Blocked item | Citation | What it needs from G14/X1 |
|---|---|---|
| **All 7 originally-declared dependencies** | `GAPS_AND_RISKS_REGISTER.md` G14 evidence: *"The manifest declares 7 deps … all `layout: grouped` — none are vendored yet (`submodules/` absent)."* | Cannot be vendored until the policy is fixed — `IMPLEMENTATION_PLAN.md` X1 explicitly states *"no second-copy `--apply` until resolved."* |
| **`containers` (R20)** | New operator mandate this round; §11.4.76 (Containers-submodule mandate) + §11.4.161 (rootless mandate), both pre-existing universal constitution anchors this project has not yet reconciled against its P12 ops work. | Needs the same G14 layout decision as the other 8 before it can be added to `helix-deps.yaml` and resolved — it is not exempt from the conflict merely because it is newly mandated. |
| **P3 — Model access layer** (`ModelProvider` chain: HelixLLM → LLMsVerifier → claude-toolkit alias → OpenAI) | `IMPLEMENTATION_PLAN.md` P3 goal + `GAPS_AND_RISKS_REGISTER.md` G14 severity line: *"blocks P3/P7/P8/P10/P11/P13 dependency work."* | The provider implementations need to know where their PATH binaries (`llms_verifier`, built from source) and HTTP endpoints (`helix_llm`) come from. |
| **P7 — API surface** | Same G14 severity-line citation. | Named directly by the register as blocked. |
| **P8 — Shared-core & clients** | Same G14 severity-line citation. | Named directly by the register as blocked. |
| **P9 — Helix ecosystem & agent interop** | `IMPLEMENTATION_PLAN.md` P9.T1 depends on P3's ModelProvider wiring, itself gated by the same G14 dependency-layout question. | Same underlying blocker as P3. |
| **P10 — OpenDesign design artifacts (R12)** | `REQUIREMENTS.md` R12; `IMPLEMENTATION_PLAN.md` P10.T1; `GAPS_AND_RISKS_REGISTER.md` G14 severity line names P10 directly; `research/opendesign_incorporation.md` §3.1. | Whether `open-design` is registered in `helix-deps.yaml` and resolved by `sync_submodules.sh` inherits whichever resolver policy G14 settles on; this round additionally clarifies it is **not** present at the currently-configured ecosystem root, only at a different sibling (`helix_terminator`). |
| **P11 — Exhaustive testing (R8)** (`Challenges`, `HelixQA`) | `IMPLEMENTATION_PLAN.md` P11.T2/P11.T3; `GAPS_AND_RISKS_REGISTER.md` G14 severity line names P11 directly. | `challenges` and `helix_qa` are two of the seven G14-blocked deps. |
| **P12 — Deployment & ops (R15) — specifically G13's compose canonicalization** | `IMPLEMENTATION_PLAN.md` P12.T4 / `GAPS_AND_RISKS_REGISTER.md` G13 (*"Two rival docker-compose.yml files"*, decision: canonicalize on `project/deploy/`). | R20 requires this work to route through the `containers` submodule's `pkg/boot`/`pkg/compose`/`pkg/health` layer rather than ad-hoc podman/docker — but `containers` itself cannot be vendored/resolved until G14 is settled, so **G13's remediation is now transitively gated on G14** in a way the register did not previously state (the register's G13 entry does not mention `containers` at all — this is a new dependency this research surfaces, not a restatement of an existing citation). |
| **P13 — Docs Chain + packaging (R10, R18)** | `IMPLEMENTATION_PLAN.md` P13.T1; `REQUIREMENTS.md` R10 and R18; `research/docs_chain_incorporation.md`: *"Do not vendor a second copy blindly — resolve X1 first (§11.4.6)."* | `docs_chain` is one of the seven G14-blocked deps; R18's documentation-delivery mandate cannot wire the Docs Chain engine in until G14 resolves. |
| **R19 — Anthropic API / `LLMProvider`** | `GAPS_AND_RISKS_REGISTER.md` G28 (Anthropic `LLMClient` design, DESIGN DONE, Go impl PENDING) does not itself depend on `LLMProvider`; the coordinator's mandate flags `LLMProvider` as a candidate R7/R19 dependency "under separate evaluation," **not yet declared in `helix-deps.yaml`.** | If/when adopted, faces the identical G14 policy question; already confirmed present at the current ecosystem root (`helix_code/submodules/llm_provider`), so it resolves as cleanly as `containers` under Option A. This document does not decide *whether* to adopt `LLMProvider` — only records its canonical-path implications if it is. |

---

## 9. Recommendation

**Recommendation: adopt Option A (§11.4.28(C) strict single-canonical), with the
already-declared parent-ecosystem copy of each dependency used as the canonical location,
`HELIX_ECOSYSTEM_ROOTS`/`--ecosystem-root` pinned to committed, documented value(s) rather
than left to per-invocation discretion, and `containers` (R20) added to `helix-deps.yaml`
alongside the original seven — all eight resolve identically, today, with zero new
vendoring.**

**Rationale (reasoned, not hedged):**

1. **It is already built and already correct.** `scripts/sync_submodules.sh`
   implements Option A's exact resolution order today, mutually exclusive by construction.
   Choosing Option A requires **zero** new code for the eight owned-org dependencies now
   identified — only closing two loose ends: (a) commit an `HELIX_ECOSYSTEM_ROOTS` default
   so the choice is not per-developer, and (b) fix `helix-deps.yaml`'s and
   `IMPLEMENTATION_PLAN.md`'s wording to stop citing R9's "both copies … in sync" framing
   as if it were the active policy.
2. **It satisfies §11.4.28(C) by the plain text of the clause, not by a reconciliation
   argument.** All eight core dependencies (the original seven plus `containers`) are
   confirmed owned-org — the third-party exemption and the constitution-submodule-only
   depth-1 carve-out both do not apply. Option B's compliance depends on a guard component
   that has to be built, tested, and kept correct forever; Option A's compliance is
   structural.
3. **Every one of these eight dependencies already has exactly one real, live, git-tracked
   copy on this host, at the exact ecosystem root this project is already configured to
   resolve against**, verified directly. Vendoring a second copy under Option B would be
   creating a rival copy of something that demonstrably already works as a single copy
   today — and `helix_ota`'s own live duplicate `containers` registration (§3) is direct,
   observed proof of what that rival-copy state looks like once it happens, unguarded, in
   this exact ecosystem.
4. **The sibling-project layout survey (§3) supports matching `helix_code`'s own
   convention, not inventing a new one.** The ecosystem is not uniform across all four
   projects surveyed, but the one convention that is actually load-bearing for
   `skill-system` — because it is the literal path every real dependency resolution
   already passes through — is `helix_code`'s grouped `submodules/<name>` layout, which
   already matches this project's own `helix-deps.yaml` declarations.
5. **None of the eight (save `embeddings`, and now potentially `llm_provider`) is
   Go-imported**, confirmed by reading `project/go.mod` in full. Most are consumed as HTTP
   services, PATH binaries, or an exercised test/tooling corpus — none of these
   consumption modes benefit from a second local git checkout; they need the *artifact*
   (a running service, a built binary, a package layer to call into) to be reachable,
   which the ecosystem copy already provides.
6. **R9's intent ("parent-dir versions have PRIORITY") is fully honored by Option A** — the
   parent-dir copy *is* used, with priority, every time. What Option A does **not** do is
   also keep a second, local, synced copy — which is the part of R9 in direct tension with
   §11.4.28(C)'s literal text for owned-org dependencies. Recommending Option A is
   therefore a recommendation to resolve that specific tension in the constitution's favor,
   not a claim that R9 is fully satisfiable verbatim.

**This recommendation does not resolve `open-design`'s placement** — it is third-party,
not one of the owned-org deps, is not present at the currently-configured ecosystem root,
and is exempt from the nested-chain clause regardless of which G14 option wins; it either
resolves via a *second* `--ecosystem-root` pointed at `helix_terminator`, or vendors fresh
once at `skill-system/submodules/open_design/`, under either G14 option (§7).

**This recommendation does not itself implement R20's routing requirement** — it only
establishes that `containers` resolves cleanly, at zero new vendoring cost, once added to
`helix-deps.yaml` under Option A. *How* `internal/validation/sandbox.go` (G02) and the
ops-hardening compose work (G13) are rewritten to call into `containers`' `pkg/boot`/
`pkg/compose`/`pkg/health` layer instead of shelling out directly is explicitly out of this
document's scope (§4, §9 gaps).

### §11.4.66 question(s) the operator must answer to proceed

1. **Does the operator accept Option A (parent-ecosystem copy is the sole canonical
   location; `skill-system` vendors nothing locally for any already-ecosystem-resolved
   dependency) as the resolution of the R9-vs-§11.4.28(C) conflict — superseding R9's
   "both copies … in sync" clause for these owned-org dependencies specifically?** This is
   the load-bearing question; everything else follows from the answer.
2. **If yes:** should `HELIX_ECOSYSTEM_ROOTS` be committed as a project default pointing
   at `/home/milos/Factory/projects/tools_and_research/helix_code`, and should a
   **second** root, `/home/milos/Factory/projects/tools_and_research/helix_terminator`, be
   committed alongside it specifically to resolve `open-design` without a fresh local
   vendor? Or does the operator require CI/other hosts to resolve ecosystem roots
   differently?
3. **Is `containers` (R20) to be added to `helix-deps.yaml` as part of this same G14
   remediation commit, or does the operator want it staged as a separate, dedicated P12/G13
   commit** (given it also requires the separate, not-yet-scoped engineering work of
   rewriting the sandbox/compose call sites to use its `pkg/*` layer, per §4's scope
   boundary)?
4. **Is `LLMProvider` (R7/R19) adopted as a declared dependency now, or does it remain
   "under separate evaluation" and out of `helix-deps.yaml` until that evaluation
   concludes?** This document only records that, if adopted, it resolves identically to
   `containers` under Option A (already present at `helix_code/submodules/llm_provider`).
5. **If Option B is required instead:** who owns building and testing the new mirror-pin +
   write-guard mechanism before any `--apply` is trusted, and does the operator authorize
   rewriting `helix-deps.yaml`'s header text (currently asserting strict "ONE canonical
   location") to describe the mirror-and-guard framing explicitly?

### Reversibility classification (§11.4.101)

| Choice | Reversible? | Basis |
|---|---|---|
| Adopting Option A today (no code change) | **Fully reversible, negligible cost.** | Nothing is committed or mutated by choosing it; it is the status quo of the existing script. Switching to Option B later is additive new engineering, not an undo. |
| Committing `HELIX_ECOSYSTEM_ROOTS` default value(s), including a second root for `open-design` | **Reversible.** | A config value change; no data at risk; §9.2 backup applies trivially to a one-line edit. |
| Adding `containers` (and, if adopted, `llm_provider`) to `helix-deps.yaml` | **Reversible.** | Plain-text manifest edit; the entries resolve to an already-existing ecosystem copy under Option A, so there is no new git object created by this step alone. |
| Building Option B's mirror-pin + write-guard mechanism, then running `--apply` to vendor local mirror submodules | **Reversible at the git-mechanics level** (submodule add/deinit is a standard, backed-up (§9.2), non-force operation) **but NOT cost-free** — once P3/P9/P10/P11/P12/P13 implementers begin writing build scripts, Go `replace` directives, or CI steps that assume a specific one of the two layouts, reverting the *policy* becomes a bounded-but-real refactor across every phase that consumed it, not merely a `git submodule deinit`. This is the operative argument for deciding **now**, while zero dependents exist (verified: `submodules/` does not exist yet; `go.mod` has zero references to any declared dependency). |
| Editing `helix-deps.yaml`'s header wording | **Fully reversible.** | Plain-text documentation change; git history preserves the prior wording regardless. |
| Rewriting G02's sandbox and G13's compose call sites to route through `containers`' `pkg/*` layer (R20 implementation, out of scope here) | **Not evaluated in this document** — that is new engineering work with its own reversibility profile, to be assessed when that design is actually written (§9). | — |

---

## 10. Honest gaps

- **`embeddings` is the one originally-declared dependency with its own importable Go
  module** (`digital.vasic.embeddings`, confirmed by reading its `go.mod`) not yet imported
  by `skill-system` (P3.T3 not yet implemented). A fully rigorous treatment would also
  weigh a *third* path for this one dependency — resolving it as a normal Go module
  dependency via `go.mod`/GOPROXY rather than any git-submodule vendoring, which would
  sidestep the G14 question for it entirely. Not evaluated here (out of the assigned
  scope); flagged `UNCONFIRMED` — whether `digital.vasic.embeddings` is reachable via a
  public Go module proxy was not tested.
- **`upstreams/` presence (§11.4.36 trigger) was verified directly only for `embeddings`
  and `containers`** among the now nine-dependency set. The remaining seven ecosystem
  copies' directory listings were not individually inspected for an `upstreams/` directory
  in this pass — affects only Option B's "must run `install_upstreams` on first vendor"
  obligation (§6) and should be checked, not assumed, if that path is ever taken.
- **The current `project/helix-deps.yaml` does not include the full §11.4.31 schema** — it
  is missing the `transitive_handling.recursive`, `transitive_handling.conflict_resolution`,
  and `language_specific_subtree` top-level fields. `containers`' own `helix-deps.yaml`
  (read in full, §4) **does** include all three, so a concrete, in-ecosystem example of
  full-schema compliance now exists to model the fix on. This is a pre-existing
  manifest-completeness gap independent of the G14 decision, noted because this research
  read the full file.
- **CI-host behavior under Option A was not tested** — this research verified the
  ecosystem copies exist on *this* development host; it did not verify what a CI runner or
  a second operator's machine currently has checked out, so the "fresh clone with no
  `--ecosystem-root`" failure mode (§5) is stated as a structural risk from reading the
  script's logic, not as an observed failure on a second host.
- **Whether `open-design` should be added as an entry in `helix-deps.yaml` at all** (so it
  is resolved by the same `sync_submodules.sh` mechanism, per `opendesign_incorporation.md`'s
  recommendation) is a P10-scoped follow-on decision, not resolved here — this document
  records only its canonical-path implications under each G14 option.
- **`LLMProvider`'s adoption decision (R7/R19) is explicitly out of scope here** — the
  coordinator flagged it "under separate evaluation." This document verified the
  apparent two-org naming is a multi-upstream-mirror pair, not two divergent repos
  (`git remote -v` on `helix_code/submodules/llm_provider` shows `HelixDevelopment/
  LLMProvider.git` as the primary `github`/`origin`/`upstream` remote and an *additional*
  `vasicdigitalgithub` push remote to `vasic-digital/LLMProvider.git` — both are real,
  live GitHub repositories per `gh repo view`, consistent with §2.1 multi-upstream push,
  not a fork ambiguity) — but does not adjudicate whether the dependency itself should be
  adopted.
- **R20's actual implementation — rewriting `internal/validation/sandbox.go` (G02) and the
  ops-hardening compose work (G13, `research/ops_hardening_design.md`) to call into
  `containers`' `pkg/boot`/`pkg/compose`/`pkg/health` layer instead of ad-hoc
  podman/docker — is not designed in this document.** This package establishes only that
  the `containers` submodule itself resolves cleanly under Option A; the follow-on
  engineering design (which functions to call, how `pkg/compose` maps onto the
  `project/deploy/docker-compose.yml` canonicalization G13 already decided on) is a
  separate, not-yet-written P12-scoped task, explicitly flagged rather than silently
  assumed complete.
- **`helix_ota`'s duplicate `containers` registration was read and cited as precedent
  evidence; it was not investigated for *why* it exists** (e.g., a migration-in-progress
  from root to grouped layout, or a deliberate transitional state) — stated as an observed
  fact, not diagnosed as a root cause, since diagnosing a sibling project's own history is
  outside this document's assigned scope.

---

*This document is a decision package only. Per §11.4.66 / §11.4.101, no `--apply` run
against `sync_submodules.sh`, no edit to `helix-deps.yaml`, and no new submodule was
created while preparing it — including while incorporating the mid-task R20 mandate and
sibling-layout evidence, which was gathered exclusively via read-only commands (`cat`,
`grep`, `git remote -v`, `git log -1`, `gh repo view`).*
