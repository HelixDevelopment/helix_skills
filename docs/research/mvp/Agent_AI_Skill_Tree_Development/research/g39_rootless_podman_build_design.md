# G39 — Rootless-Podman Conversion Design (§11.4.161 / §11.4.76)

**Revision:** 2
**Last modified:** 2026-07-15T21:33:52Z

## As-landed addendum

**Verified against HEAD, 2026-07-15T21:33:52Z (§11.4.6 — read, not guessed).**
§2's touchpoint table below ("Current rootful-Docker touchpoints") predates
commit `b597623`, which closed rows #1-#3 (and shifted row #4's cited line
ranges as a result):

- **Row #1** (`project/Makefile:36` `CONTAINER_RUNTIME ?= docker`) —
  **FIXED.** At HEAD, `project/Makefile:41` reads
  `CONTAINER_RUNTIME ?= podman`, preceded by a `§11.4.161`-citing comment
  block (Makefile:36-40) explaining the rootless-first rationale and the
  operator override path (`make CONTAINER_RUNTIME=docker ...`).
- **Row #2** (`COMPOSE_CMD` "inherits #1's rootful default") — **FIXED** as a
  direct consequence of #1: `COMPOSE_CMD ?= $(CONTAINER_RUNTIME) compose`
  (Makefile:42) now expands podman-first.
- **Row #3** (hardcoded `docker tag`/`docker push` in the `docker-push`
  target, originally cited at `Makefile:254-259`) — **FIXED.** At HEAD the
  target (now at `Makefile:259-264`) reads `$(CONTAINER_RUNTIME) tag
  $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest` /
  `$(CONTAINER_RUNTIME) push $(IMAGE_NAME):$(IMAGE_TAG)` /
  `$(CONTAINER_RUNTIME) push $(IMAGE_NAME):latest` (Makefile:262-264) — no
  literal `docker` invocation remains in this target.
- **Row #4**'s cited line ranges (`Makefile:131-140, 245-277, 318-322`) have
  shifted with the fix (the `docker-build`/`docker-push`/`docker-up`/
  `docker-down`/`docker-logs` targets now sit at approximately
  `Makefile:250-278`); the substance of row #4 (these targets consume
  `$(COMPOSE_CMD)`, now podman-first thanks to #1/#2) is unaffected by the
  line-number shift.
- **`check_container_runtime_default.sh` was added** alongside the fix
  (verified present at `project/scripts/check_container_runtime_default.sh`).
- **Rows #7-#12 remain exactly as described below — still open, re-verified
  unchanged.** Re-checked at HEAD: `project/scripts/_lib.sh`'s
  `hx_detect_engine_soft()` still tries `docker compose` before `podman
  compose` / `podman-compose` (still docker-first ordering), and
  `project/scripts/migrate.sh`'s `detect_compose()` still tries `docker
  compose` first, then `docker-compose`, then `podman-compose` (still its
  own independent, still-docker-first, still-divergent detection chain).
  §4 (Minimal fix) and §6.3 (paired mutation, using HEAD-before-fix as the
  golden-bad fixture) below remain the valid, not-yet-applied framing for
  these still-open rows.

§4's golden-bad/golden-good framing and §6.3's paired mutation below already
correctly anticipate this exact partial-fix scenario — as §6.3 step 4 states,
"assert the gate FAILs, citing specifically that file — proving each of the
twelve touchpoints in §2 is independently load-bearing to the gate." Rows
#1-#3 have now individually flipped to golden-good; rows #7-#12 have not.
Nothing in §3, §5, §6, §7, or §8 below is affected by this addendum — those
sections describe the target-state Containers-submodule mapping, the
rootless preflight design, the gate design, the honest §11.4.112 boundary,
and captured anti-bluff evidence, none of which depended on rows #1-#3 being
still-open.

> **Gap-id cross-reference note (honest, §11.4.6 no-guessing):** the dispatching
> task labels this gap **G39** (§11.4.161). The on-disk audit this design was
> pointed at, `research/constitution_compliance_audit.md`, files the identical
> finding — rootful-Docker-first default across `Makefile` + `scripts/*.sh` —
> as **G46** (HIGH, §11.4.76/§11.4.161/§11.4.173; table row: *"Change
> `project/Makefile:36` default `CONTAINER_RUNTIME ?= docker` →
> `CONTAINER_RUNTIME ?= podman`; evaluate adopting the `vasic-digital/containers`
> submodule's `pkg/boot`/`pkg/compose` layer instead of raw `docker-compose.yml`
> invocations"*). This document does not guess which numbering is authoritative
> across audit rounds — it designs against the **same real finding** cited by
> both labels and is filed at the path the dispatching task specified
> (`g39_rootless_podman_build_design.md`).

---

## 1. Goal + §11.4.161

**Mandate (verbatim scope, §11.4.161):** "every project MUST use Podman in
rootless mode (or equivalent rootless container runtime) for ALL containerized
workloads — Docker in rootful mode, sudo, or any escalation to root is
FORBIDDEN unless the target platform has no rootless option AND that constraint
is documented per §11.4.112; the `vasic-digital/containers` submodule
(§11.4.76) MUST be used as the sole container orchestration layer — no ad-hoc
docker/podman commands outside `pkg/boot`/`pkg/compose`/`pkg/health`."

Two independent obligations, both currently unmet by this project:

1. **Rootless-first, not rootful-first.** Every place this project selects a
   container engine/compose provider MUST prefer rootless Podman before any
   Docker fallback, and MUST NEVER shell out to `sudo docker` / `sudo podman`.
2. **Single orchestration layer, not four reimplementations.** All container
   lifecycle logic (boot / compose / health) MUST route through the
   `vasic-digital/containers` submodule's Go API (`pkg/boot`, `pkg/compose`,
   `pkg/health`) — never ad-hoc `docker compose` / `podman-compose` shell
   invocations scattered and re-implemented per script.

Adjacent binding rules this design must stay compatible with:

- **§11.4.184** — the SonarQube local server (when wired) runs via rootless
  Podman, never rootful Docker; any ES-class dependency it or this project
  introduces inherits the `vm.max_map_count` preflight this design specifies.
- **§11.4.173** — containerized + distributed build mandate; the build path
  this design touches (compose-driven dev/test stack, not the Go compiler
  itself) must not introduce a rootful container as a side door around
  §11.4.173's own rootless requirement.

---

## 2. Current rootful-Docker touchpoints (cited, read-only findings)

Twelve distinct touchpoints across two compose files, one Makefile (four
sub-issues), four scripts (each an **independently copy-pasted** detection
chain — a DRY violation on top of the ordering bug), one packaging step, and
one doc comment. None of `project/internal/**` (Go source/tests/migrations,
owned by the concurrent subagent) is touched by any of these — every touchpoint
below lives in `Makefile`, `Dockerfile`, `docker-compose.yml`,
`deploy/docker-compose.yml`, `scripts/*.sh`, or `.env.example`.

| # | File : Lines | Finding |
|---|---|---|
| 1 | `project/Makefile:36` | `CONTAINER_RUNTIME ?= docker` — the root default names the rootful binary first. |
| 2 | `project/Makefile:37` | `COMPOSE_CMD ?= $(CONTAINER_RUNTIME) compose` — inherits #1's default; every `docker-*` target below is downstream of this one line. |
| 3 | `project/Makefile:254-259` (`docker-push` target) | **Worst offender**: `docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest` / `docker push …` (×2) are **hardcoded literal `docker`** invocations that bypass `$(CONTAINER_RUNTIME)` entirely — even flipping #1 to `podman` would NOT fix this target. |
| 4 | `project/Makefile:131-140, 245-277, 318-322` (`dev`, `docker-build`, `docker-up`, `docker-down`, `docker-logs`, `docker-ps`, `clean-all`) | All invoke `$(COMPOSE_CMD)`, inheriting #1/#2's rootful-first default. |
| 5 | `project/docker-compose.yml` (full stack: db + api + worker + optional prometheus/grafana) | The file itself is engine-neutral YAML; it is a touchpoint only because it is *consumed* rootful-first via #4. |
| 6 | `project/deploy/docker-compose.yml` (datastore-only: `postgres` service, app service intentionally commented out per its own header comment lines 5-9) | Same — engine-neutral YAML, consumed rootful-first via #7-#10 below. |
| 7 | `project/scripts/_lib.sh:108-135` (`hx_detect_engine_soft()`) | Tries, **in this literal order**: `docker compose` (line 119) → `podman compose` (line 124) → `podman-compose` (line 129). Rootful-first priority order — the shared library every ops script (`start.sh`/`stop.sh`/`restart.sh`/`status.sh`/`logs.sh`) sources. |
| 8 | `project/scripts/migrate.sh:47-52` | **Independent, divergent** inline detection: `docker compose` (47) → `docker-compose` (49) → `podman-compose` (51). Does not even source `_lib.sh`'s function, and doesn't try bare `podman compose` at all. |
| 9 | `project/scripts/backup.sh:52-57` | Same divergent pattern as #8 (its own copy, not `_lib.sh`'s). |
| 10 | `project/scripts/restore.sh:46-51` | Same divergent pattern as #8/#9 (a third independent copy). |
| 11 | `project/scripts/package.sh:74-75` | Copies root `docker-compose.yml` into the distributable package unconditionally — a packaging-layer touchpoint: the shipped artifact's install docs should steer new installs to podman-first tooling, not merely inherit whatever `docker-compose.yml`'s name implies. |
| 12 | `project/.env.example:138` | Comment example `# Monitoring (optional - enable with: docker compose --profile monitoring up)` — cosmetic, but it is documentation actively modelling the rootful-first invocation as the canonical example. |

**Root cause, stated as fact (§11.4.6):** four different files (`_lib.sh`,
`migrate.sh`, `backup.sh`, `restore.sh`) each reimplement the same
docker-vs-podman detection logic **independently and inconsistently** — three
of the four never try bare `podman compose`, only `_lib.sh` does. This is not
merely a wrong default; it is the exact "ad-hoc docker/podman commands outside
`pkg/boot`/`pkg/compose`/`pkg/health`" pattern §11.4.76 forbids by name.

---

## 3. Containers-submodule rootless API mapping

### 3.1 The real, working reference implementation (helix_ota)

`helix_ota` (a live sibling checkout at `/home/milos/Factory/projects/tools_and_research/helix_ota`,
already vendoring `vasic-digital/containers` at root-level `containers/` per
the family precedent this project's own G14 design adopted — see §3.3 below)
boots a real PostgreSQL for integration tests **exclusively** through the
submodule's Go API, never a manual `podman`/`docker compose` shell step.
Verbatim (`helix_ota/server/internal/store/postgres_integration_test.go:1-77`,
also mirrored in `postgres_fault_integration_test.go`,
`postgres_coverage_integration_test.go`, `postgres_migrations_integration_test.go`,
and the `internal/rollout/postgres_*_integration_test.go` quartet):

```go
import (
    "digital.vasic.containers/pkg/boot"
    "digital.vasic.containers/pkg/compose"
    "digital.vasic.containers/pkg/endpoint"
    "digital.vasic.containers/pkg/health"
    "digital.vasic.containers/pkg/logging"
    "digital.vasic.containers/pkg/runtime"
)

// --- boot PostgreSQL on-demand via the containers submodule ---
rt, err := runtime.AutoDetect(ctx)                    // podman-first priority
orch, err := compose.NewDefaultOrchestrator(projectDir, logging.NopLogger{})

ep := endpoint.NewEndpoint().
    WithHost("localhost").WithPort(pgHostPort).
    WithHealthType("tcp").WithRequired(true).WithEnabled(true).
    WithComposeFile("postgres.compose.yml").WithServiceName("postgres").
    WithTimeout(120 * time.Second).WithRetryCount(60).
    Build()

mgr := boot.NewBootManager(
    map[string]endpoint.ServiceEndpoint{"postgres": ep},
    boot.WithRuntime(rt),
    boot.WithOrchestrator(orch),
    boot.WithHealthChecker(health.NewDefaultChecker()),
    boot.WithProjectDir(projectDir),
    boot.WithLogger(logging.NopLogger{}),
)

summary, err := mgr.BootAll(ctx)
// ... t.Cleanup(func() { _ = mgr.Shutdown(context.Background()) })
```

The file's own header comment states the anti-bluff intent verbatim: *"proves
the pgx/PostgreSQL Repository satisfies the exact same behavioural contract as
the in-memory one, against a REAL PostgreSQL that is booted on-demand through
the containers submodule (digital.vasic.containers) — **never** a manual
`podman`/`compose` step, never a fake."*

### 3.2 `runtime.AutoDetect` is already podman-first by construction

`helix_ota/containers/pkg/runtime/detect.go:38-45`:

```go
RuntimePriority = []string{
    "podman",
    "docker",
    "nerdctl",
    "cri-o",
    "lxd",
    "kubernetes",
}
```

`AutoDetect(ctx)` calls `AutoDetectWithPriority(ctx, GetRuntimePriority())`
(`detect.go:104-106`) — podman is tried **first**, unconditionally, on every
call. This is precisely the ordering this project's shell scripts (touchpoints
#7-#10) get backwards.

### 3.3 The compose layer already solves the "docker-is-actually-a-podman-shim" classification problem

`helix_ota/containers/pkg/compose/orchestrator.go:80-166` implements exactly
the ambiguity this project's naive shell detection cannot resolve: a resolved
`docker` binary may be a **podman-docker compatibility shim** (confirmed live
on this build host, §7 below) rather than real Docker Engine. Verbatim
comments + code:

```go
// isPodmanBackedCmd reports whether the resolved compose command actually
// delegates to podman even when composeCmd is not literally "podman-compose"
// ... a podman-docker compatibility-shim host `docker compose version` exits
// 0 by silently re-execing into podman-compose ...
var podmanBannerMarkers = []string{"podman"}

func isPodmanBackedCmd(composeCmd string, composeArgs []string, timeout time.Duration) bool {
    if composeCmd == "podman" || composeCmd == "podman-compose" {
        return true
    }
    // ... probes the resolved command's own banner output for a podman marker
}
```

This is the **authoritative, already-tested classifier** this project's four
divergent shell copies (`_lib.sh`, `migrate.sh`, `backup.sh`, `restore.sh`)
should delegate to rather than reimplement — reimplementing it in Bash a
fifth time (even correctly) would still violate §11.4.76's "no ad-hoc
docker/podman commands outside `pkg/boot`/`pkg/compose`/`pkg/health`" clause.

### 3.4 Touchpoint → replacement mapping

| Touchpoint(s) | Replacement via Containers submodule |
|---|---|
| #1-#5 (Makefile `docker-*` family + root `docker-compose.yml`, the app+db+worker dev stack) | A new `cmd/devstack` (or `internal/devstack`) Go entrypoint invoked by `make dev` / `make dev-down`, calling `runtime.AutoDetect(ctx)` → `compose.NewDefaultOrchestrator(".", logger)` → `boot.NewBootManager(endpoints, boot.WithRuntime(rt), boot.WithOrchestrator(orch), boot.WithHealthChecker(health.NewDefaultChecker()), boot.WithProjectDir("."), boot.WithLogger(logger))` → `mgr.BootAll(ctx)` / `mgr.Shutdown(ctx)`, exactly mirroring §3.1. `docker-compose.yml` itself is kept as the compose FILE the Orchestrator drives (compose YAML is still the on-disk service definition `pkg/compose` consumes) — only the *selection/invocation* code changes, never a hand-rolled `$(COMPOSE_CMD)` shell-out. |
| #6-#10 (`deploy/docker-compose.yml` + `_lib.sh`/`migrate.sh`/`backup.sh`/`restore.sh` detection chains) | Collapse the four independent detection copies into **one**: either (a) minimal fix — `_lib.sh`'s `hx_detect_engine_soft()` reorders to podman-first and becomes the ONLY implementation (`migrate.sh`/`backup.sh`/`restore.sh` are edited to `source _lib.sh` and call `hx_detect_engine_soft` + `hx_compose` instead of re-declaring their own `COMPOSE_CMD` logic), or (b) full migration — a thin `cmd/opsctl` Go CLI wraps `runtime.AutoDetect` + `compose.NewDefaultOrchestrator` + `boot.NewBootManager` and `start.sh`/`stop.sh`/`restart.sh`/`migrate.sh`/`backup.sh`/`restore.sh`/`status.sh`/`logs.sh` become thin wrappers around `opsctl {up,down,migrate,backup,restore,status,logs}`, so there is exactly ONE selection algorithm in the whole project (living in the vendored submodule, not reimplemented locally at all). Path (b) is the complete §11.4.76 remediation; path (a) is the minimal HIGH-severity-closing fix that can land first. |
| #3 (hardcoded `docker tag`/`docker push`) | Replace with the resolved runtime's own image-registry verb — `$(CONTAINER_RUNTIME) tag` / `$(CONTAINER_RUNTIME) push` (minimal fix) or, in the Go-CLI path, a call through `pkg/runtime`'s `ContainerRuntime` interface's own tag/push method if the submodule exposes one, else shell out to the **resolved** `rt.Name()` binary, never a literal `docker`. |
| #11 (`package.sh` packaging step) | No code change required to the copy itself; the packaged artifact's `README`/`install.sh` (already a Category-A install touchpoint) must document rootless-Podman as the primary supported path, consistent with the corrected `Makefile`/`deploy/docker-compose.yml` defaults it now ships. |
| #12 (`.env.example:138` comment) | Cosmetic doc fix: show `podman compose --profile monitoring up` (or an engine-neutral `<engine> compose ...`) as the primary example, Docker demoted to a parenthetical fallback note. |

Honest scope note (§11.4.6): path (b) (full `pkg/boot`/`pkg/compose`/`pkg/health`
adoption replacing every shell wrapper) is the complete remediation this
design recommends as the target state; it is a larger, separately-tracked
follow-up (consistent with the audit's own framing: *"the one-line default
change is trivial; full containers-submodule adoption is a larger
follow-up"*). This design specifies **both** the minimal fix that closes the
HIGH-severity rootful-default finding immediately and the full-adoption target
state, so neither is silently dropped per §11.4.197.

### 3.5 Prerequisite: `containers` submodule is not yet vendored

`project/helix-deps.yaml` does not currently declare `containers` as a
dependency (confirmed by direct grep — no match). The G14 vendor-fresh design
(`research/g14_vendor_fresh_submodule_layout_design.md` §6, "Root-level
exception: `containers/`") already specifies its incorporation: root-level
`containers/` path (not `submodules/containers/`), `layout: ungrouped` in
`helix-deps.yaml`, `go.mod replace digital.vasic.containers => ./containers`,
mirroring `helix_ota`'s live precedent (minus its own accidental
double-mount, which G14 explicitly avoids). This design's Go-API replacement
(§3.4) is **blocked on** that vendoring step landing first; the minimal
shell-level fix (§4 below) is **not** blocked on it and can land independently
today.

---

## 4. Minimal fix (unblocks immediately, no submodule dependency)

For each touchpoint in §2, the source-level change:

1. **Makefile:36** — `CONTAINER_RUNTIME ?= docker` → `CONTAINER_RUNTIME ?= podman`.
2. **Makefile:254-259** — replace the two literal `docker tag`/`docker push`
   lines with `$(CONTAINER_RUNTIME) tag ...` / `$(CONTAINER_RUNTIME) push ...`.
3. **`_lib.sh:108-135`** — reorder `hx_detect_engine_soft()` to try, in order:
   bare `podman compose` (native, rootless-by-default per §7) → `podman-compose`
   (standalone tool) → `docker compose` (last resort, and only after the
   §3.3-style banner-probe confirms it is not itself a podman shim masquerading
   as extra evidence, OR simply demoted to last-resort regardless of shim
   status — conservative-safe per §11.4.201 clause 4).
4. **`migrate.sh`/`backup.sh`/`restore.sh`** — delete their independent inline
   detection blocks; `source "$(dirname "$0")/_lib.sh"` and call
   `hx_detect_engine` + `hx_compose` instead, collapsing four implementations
   into one.
5. **`.env.example:138`** — swap the example command's engine word.

None of these five edits touch `project/internal/**`, migrations, or Go
tests — they are strictly `Makefile` / `scripts/*.sh` / `.env.example` edits,
matching this design task's read-only-elsewhere scope.

---

## 5. Rootless prerequisites + preflight check design

### 5.1 Prerequisites (closed set)

1. **`podman` binary present** on `PATH` (`command -v podman`).
2. **subuid/subgid ranges allocated** for the invoking user — rootless Podman
   maps container-internal UID/GID ranges through `/etc/subuid` and
   `/etc/subgid`; an entry for the user MUST exist (`grep -q "^${USER}:"
   /etc/subuid /etc/subgid`).
3. **Rootless mode actually active** — `podman info --format
   '{{.Host.Security.Rootless}}'` MUST report `true` (a `podman` binary being
   *present* does not itself prove it is configured/running rootless — this
   is the exact §11.4.201 distinction between a proxy signal and the real
   condition).
4. **A compose provider resolves** — `podman compose version` (the
   docker-compose-compatible provider invoked via the podman CLI) OR
   `podman-compose --version` (the standalone Python tool).
5. **A rootless network backend is available** — `pasta` or `slirp4netns` on
   `PATH` (`command -v pasta || command -v slirp4netns`); rootless Podman
   needs one of these to publish container ports to the host network
   namespace without root.
6. **`vm.max_map_count` floor for any ES-class service** — read
   `/proc/sys/vm/max_map_count` directly (portable; do **not** depend on a
   `sysctl` binary being on `PATH` — this build host does not have one, §7)
   and require it `≥ 262144` (Elasticsearch's own documented floor). This
   project has no ES-class service today, but §11.4.184's SonarQube local
   server does, and the preflight is written generically so it applies the
   moment such a dependency is introduced.
7. **systemd `--user` lingering, for any long-running rootless service** — if
   the ops layer's existing `systemctl --user` wiring (`_lib.sh`'s
   `hx_systemd_unit_path`/`hx_has_systemd_user`, already present) is used to
   keep a rootless-Podman-backed service running after the interactive
   session ends, `loginctl show-user "$USER" --property=Linger` MUST report
   `Linger=yes` (else the container stops at logout) — `loginctl
   enable-linger "$USER"` is the one-time fix, itself a privileged-adjacent
   op (touches system state under the invoking user's own account, no `sudo`
   needed) that the preflight reports as an actionable remediation rather
   than performing silently.

### 5.2 Preflight design (`scripts/rootless_preflight.sh`, new, in the same
`scripts/` directory as the other ops helpers)

```
rootless_preflight() {
    local fail=0
    command -v podman >/dev/null 2>&1 || { echo "MISSING: podman binary"; fail=1; }
    grep -q "^${USER}:" /etc/subuid 2>/dev/null || { echo "MISSING: /etc/subuid entry for $USER"; fail=1; }
    grep -q "^${USER}:" /etc/subgid 2>/dev/null || { echo "MISSING: /etc/subgid entry for $USER"; fail=1; }
    [[ "$(podman info --format '{{.Host.Security.Rootless}}' 2>/dev/null)" == "true" ]] \
        || { echo "FAIL: podman is not running rootless (real condition, not inferred)"; fail=1; }
    { podman compose version >/dev/null 2>&1 || command -v podman-compose >/dev/null 2>&1; } \
        || { echo "MISSING: no podman-backed compose provider resolves"; fail=1; }
    { command -v pasta >/dev/null 2>&1 || command -v slirp4netns >/dev/null 2>&1; } \
        || { echo "MISSING: no rootless network backend (pasta/slirp4netns)"; fail=1; }
    if needs_es_class_service; then   # project-declared, not this script's concern to guess
        local mmc; mmc="$(cat /proc/sys/vm/max_map_count 2>/dev/null || echo 0)"
        (( mmc >= 262144 )) || { echo "FAIL: vm.max_map_count=$mmc < 262144"; fail=1; }
    fi
    return "$fail"
}
```

Every ops script that boots containers (`start.sh`, `dev`/`docker-up` in the
Makefile once migrated) calls `rootless_preflight` before `hx_compose up`, and
refuses to proceed on failure with the specific missing/failing item named —
never a generic "container start failed."

---

## 6. `CM-ROOTLESS-CONTAINER` gate + §1.1 mutation (real-condition per §11.4.201)

### 6.1 Why a bare `grep docker` is the wrong design (the false-positive the task warns against)

A naive `grep -rn docker Makefile scripts/ *.yml` would flag, as false
positives:

- **Prose/comments**, e.g. `_lib.sh:9-18`'s own docstring listing "docker
  compose / podman compose / podman-compose" as the three detected
  implementations, and `_lib.sh:144-147`'s error message that names all three
  by design (it is telling the operator what was checked, not selecting one).
- **`.env.example:138`**'s comment (a real touchpoint, but for a different
  reason — it is a stale *documentation example*, not a code default; conflating
  the two would misclassify the fix).
- **The podman-docker compat shim case** (confirmed live on this host, §7): a
  literal `docker compose` invocation in a script may, on a given host, already
  be running rootless Podman under the hood via `/usr/bin/docker` being a
  shell script that re-execs into `podman`. A bare string match cannot tell
  the difference between "this line runs rootful Docker" and "this line
  happens to run rootless Podman today because of an installed shim, but will
  silently run rootful Docker the moment a real `dockerd` is installed
  instead." Only the **second** case is what makes the current code wrong —
  the shim's presence does not exonerate the code, it merely masks the defect
  on this one host (§7's honest-boundary point).

### 6.2 Two-part design: static (fast, CI-friendly) + dynamic (real-condition, §11.4.201)

**Part A — static, code-vs-prose discriminated:**

Scans tracked files (`Makefile`, `scripts/*.sh`, `.env.example`, compose
YAML) and FAILs on:

1. `sudo\s+(docker|podman)\b` anywhere (comment or code) — privilege
   escalation is forbidden unconditionally, no prose exception.
2. A variable-default assignment whose right-hand side is the literal
   `docker` for a container-runtime-selecting variable (pattern: `^\s*[A-Z_]*RUNTIME\s*[?:]?=\s*docker\b`,
   catching Makefile:36's exact shape) — this is a CODE default, never prose.
3. A detection-chain function/script whose **first** successfully-tested
   branch names `docker compose` / `docker-compose` ahead of any `podman`
   variant (structural check on the ordered list of `command -v` /
   `... version` probes in the function body, not a bag-of-words match) —
   catches `_lib.sh:119-122`, `migrate.sh:47-49`, `backup.sh:52-54`,
   `restore.sh:46-48` precisely, while `_lib.sh`'s own docstring/error-message
   lines (which merely *mention* all three names without ordering a
   preference) do not match this structural pattern and are correctly
   excluded.
4. A literal `docker <verb>` invocation (verb ∈ `run|build|push|tag|login|exec`)
   that is **not** reached through the resolved `$(CONTAINER_RUNTIME))`/
   `hx_compose` abstraction — catches Makefile:257-259's hardcoded
   `docker tag`/`docker push`.

Lines inside `#`-prefixed comment blocks are excluded from checks 2-4 (only
check 1 — the escalation ban — has no prose exception, since even a comment
instructing "run with sudo docker" is itself a governance-violating
instruction, distinct from a comment merely *mentioning* the word "docker").

**Part B — dynamic, real-condition (the §11.4.201 core of this gate):**

Rather than trusting static text, Part B **executes** the project's own
resolution logic in the current environment and asserts the *outcome*:

1. Source `scripts/_lib.sh`, call `hx_detect_engine_soft`, and assert
   `HX_COMPOSE_BIN[0]` resolves to `podman` or `podman-compose` — never
   `docker`. This is real per §11.4.201: it doesn't check what string appears
   in the source file, it runs the actual selection function and inspects
   what it *returns* in this real environment (where both `docker` — the
   shim — and native `podman` are on `PATH`, exactly as confirmed in §7).
2. For whichever binary step 1 resolves, independently confirm rootlessness
   from the authoritative source: `podman info --format
   '{{.Host.Security.Rootless}}'` == `true`. If the resolved binary is
   literally `docker` (shim or real), reuse the Containers submodule's own
   `isPodmanBackedCmd`-style banner probe (§3.3) to classify it — if it does
   NOT delegate to podman, this is an unconditional FAIL (real, unshimmed
   rootful Docker resolved first) regardless of what step 1 or the static
   scan concluded.
3. For the Makefile: `make -n CONTAINER_RUNTIME= docker-build 2>&1` (dry-run,
   `-n` prevents any real build) piped through a check that the printed
   `$(COMPOSE_CMD)` expansion begins with `podman`, never `docker`.
4. Assert `command -v sudo` is never invoked as part of any of the above
   resolution paths' actual subprocess calls (traced via a `PATH`-shadowing
   sentinel `sudo` shim during the gate's own test run that records + fails
   if called — the mechanical instantiation of check A.1's static ban).

A tree PASSES `CM-ROOTLESS-CONTAINER` only when **both** Part A and Part B are
green. Part A alone would be gameable by an ordering trick that still resolves
to Docker in practice; Part B alone would miss a hardcoded rootful invocation
that happens not to fire in the gate's specific test environment (e.g.
Makefile:257-259's `docker push`, which only runs on an actual release push,
not on every `make docker-build`). Both are required, composing exactly as
§11.4.201 requires ("every guard/gate MUST assert the REAL condition ... from
the AUTHORITATIVE source").

### 6.3 Paired §1.1 mutation — using the repo's own pre-fix state as the golden-bad fixture

Unusually, this project does not need to *invent* a synthetic golden-bad
fixture: **the repository's current, unfixed HEAD is itself the golden-bad
case.** The mutation-test design is:

1. **Golden-bad = HEAD before the §4 minimal fix lands.** Running the §6.2
   gate against the current tree (Makefile:36 `CONTAINER_RUNTIME ?= docker`,
   `_lib.sh:119-122` docker-first ordering, `migrate.sh`/`backup.sh`/
   `restore.sh`'s independent docker-first chains, Makefile:257-259's
   hardcoded `docker tag`/`push`) MUST produce FAIL, citing each of the four
   Part-A findings + a Part-B mismatch (native `podman compose` — the correct
   choice — being available but `_lib.sh` still resolving `docker compose`
   first). This is not hypothetical: it is exactly what running the designed
   gate today, on the actual unmodified files read in §2, produces.
2. **Golden-good = HEAD after the §4 minimal fix.** After the five edits in
   §4 land, re-running the identical gate MUST produce PASS — same repo, same
   host, same available binaries, only the source ordering/defaults changed.
3. **Negative control (the false-positive guard, §11.4.201(3)):** a tree
   where the ONLY docker-mentioning lines are `_lib.sh`'s own docstring
   (lines 9-18) and error message (lines 144-147), with the golden-good
   ordering/defaults from step 2 otherwise intact, MUST still PASS — proving
   the gate does not fire on prose, only on the structural code patterns
   defined in §6.2 Part A.
4. **The mutation itself:** starting from golden-good, re-apply the exact
   pre-fix diff for **one** file at a time (e.g. revert only Makefile:36 back
   to `CONTAINER_RUNTIME ?= docker`, leaving the four scripts fixed) and
   assert the gate FAILs, citing specifically that file — proving each of the
   twelve touchpoints in §2 is independently load-bearing to the gate, not
   merely one of them accidentally carrying the whole check.

This design is stronger than an invented fixture because every one of its
three fixture states (golden-bad, golden-good, negative-control) is either
the literal current repository state or a small, precisely-described diff
from it — fully reproducible from this document without guessing.

---

## 7. Honest boundary (§11.4.112)

§11.4.112 permits a documented Docker-rootful exception **only** when deep
research + reproducible captured evidence PROVE rootless is structurally
impossible on the target platform — never on convenience grounds. On the
actual build host this design was researched against, no such exception
applies (confirmed, §8 below: rootless Podman 5.7.1 is fully operational,
subuid/subgid are allocated, a compose provider resolves, and a live rootless
container is running right now). If a *different* target platform for this
project's operators/CI genuinely lacks user namespaces (e.g. a locked-down
corporate host with `CONFIG_USER_NS` disabled and no admin path to enable it,
or a legacy kernel predating rootless support), the exception path is:

1. Attempt `podman info` and `unshare --user` on that specific host; capture
   the actual failure (e.g. `ERRO[0000] cannot ... user namespaces are not
   enabled`) as evidence — never assert impossibility from a template guess.
2. Record the finding next to the touchpoint it affects (e.g. a
   platform-specific note in `Makefile` or `README.md`) citing §11.4.112,
   the captured failure text, and the date verified.
3. The exception is scoped to that ONE platform, never generalised project-wide
   — this project's default (Makefile/`_lib.sh`/`.env.example`) stays
   rootless-Podman-first for every other platform.

No such exception is needed or filed by this design for the current
development host.

---

## 8. Anti-bluff evidence (captured, this session, read-only)

All commands below were run read-only against the live host to ground this
design in fact rather than assumption (§11.4.6), without touching any file
under `project/`:

- **Rootless confirmed at the source:** `podman info --format
  '{{.Host.Security.Rootless}}'` → `true`.
- **The "docker" binary on this host is NOT real Docker — it is the
  podman-docker compatibility shim:** `file /usr/bin/docker` →
  `POSIX shell script, ASCII text executable`; invoking `docker ps` prints the
  shim's own banner first: `Emulate Docker CLI using podman. Create
  /etc/containers/nodocker to quiet msg.`, then delegates to the real podman
  backend (identical container listing to `podman ps`).
- **The project's own isolated test-DB container already runs on rootless
  Podman, not rootful Docker** — `podman ps -a --format '{{.Names}}' | grep
  skillsys` and `docker ps -a --format '{{.Names}}' | grep skillsys` both
  list `skillsys_p1t1_test_pg`; `podman inspect skillsys_p1t1_test_pg` shows
  image `docker.io/pgvector/pgvector:pg16`, networked via `pasta` (the
  rootless-specific port-publishing backend, confirmed in the inspect output:
  `pasta map[5432/tcp:[{0.0.0.0 55433}]]`). **This is exactly the finding this
  design exists to fix stated precisely**: the currently-running test
  database is *accidentally* rootless-compliant on this one host only because
  a compat shim happens to be installed — the actual Makefile/script code
  (§2) still explicitly codifies a rootful-Docker-first default that would
  silently select real rootful `dockerd` the moment this project's build runs
  on any host that has genuine Docker Engine installed instead of (or ahead
  of, in `PATH` order) the shim. Accidental compliance on one developer's
  machine is not compliance — it is the exact PASS-bluff pattern §11.4.6/§11.4.201
  forbids: the code's *default* is wrong regardless of today's *incidental*
  runtime outcome.
- **subuid/subgid allocated:** `/etc/subuid` → `milos:100000:65536`;
  `/etc/subgid` → `milos:100000:65536`.
- **Rootless UID mapping live:** `podman unshare cat /proc/self/uid_map` →
  `0 1000 1` / `1 100000 65536` (host uid 1000 maps to container uid 0; the
  subordinate range maps container uids 1-65535 to host uids
  100000-165535) — direct proof of an active rootless user-namespace mapping,
  not merely a configured-but-unused range.
- **Compose providers present:** `podman --version` → `podman version 5.7.1`;
  `command -v podman-compose` → `/usr/bin/podman-compose`; `podman compose
  version` resolves and reports `podman-compose version 1.5.0` under the hood
  (both the native `podman compose` provider and the standalone
  `podman-compose` tool are available).
- **`vm.max_map_count` already satisfies any future ES-class dependency's
  floor:** `cat /proc/sys/vm/max_map_count` → `2147483642` (≫ Elasticsearch's
  documented 262144 floor); confirmed via `/proc/sys/vm` directly because this
  host has no `sysctl` binary on `PATH` at all (`sysctl vm.max_map_count` →
  `command not found`) — validating §5.2's design choice to read
  `/proc/sys/vm/max_map_count` directly rather than depending on the `sysctl`
  wrapper.

No blocker was hit gathering this evidence; every fact above is a real,
reproducible command output from the actual build host, not an inference.
