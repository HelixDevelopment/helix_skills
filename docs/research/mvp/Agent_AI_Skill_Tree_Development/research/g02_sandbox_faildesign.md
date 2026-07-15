# G02 — Fail-Closed Design Decision: Skill Validation & (Non-)Execution of Untrusted Code

**Revision:** 1
**Last modified:** 2026-07-15T15:38:45Z
**Status:** DECIDED (design-research; drives a later Go fix — no code changed by this doc)
**Scope:** Remediation design for gaps-register **G02** (sandbox RCE), folding in **G16** (fake WASM / conflicting Docker mounts), **G03** (validation package dead code), **G05** (jury auto-approve), **G21** (shallow resource verification).
**Authority:** Constitution §11.4.8 (deep-research-before-fix), §11.4.6 (no-guessing / decision-not-maybe), §11.4.112 (structural-impossibility honesty), §11.4.161 (rootless-Podman / containers-submodule mandate), §11.4.123 (rock-solid-proof-or-research), §11.4.197 (started work is tracked, never un-wired).

> **Deliverable contract:** this is the ONLY file created by this research pass. It produces a *vetted decision*, not code. Every source is cited with its URL + fetch date (§11.4.8/§11.4.99).

---

## 0. TL;DR (the decision, plainly — §11.4.6)

1. **Execution of untrusted skill/POC code is NOT required to validate a skill.** Static validation (parse/syntax/type-check + resource verification + cross-reference + LLM jury) suffices for the MVP. The MVP therefore **executes nothing untrusted** — closing the RCE by *construction* (there is no host-exec code path left to misconfigure).
2. **Delete the host-execution "sandbox" entirely.** `WASMSandbox` (which never uses WASM and runs `go run` / `python -c` / `bash -c` on the host), `executeProcess`, `executeGoSnippet`, the `NoOpSandbox` auto-pass, the fake `LD_PRELOAD` "network restriction," and the `sandbox_type="wasm"` default all go. The word "sandbox" applied to a host-exec path is itself the anti-bluff (R11) violation.
3. **IF and only if a genuine "run a POC and observe its output" feature is ever added, it is opt-in, off-by-default, and runs behind exactly ONE real isolation boundary: rootless Podman (per §11.4.161, via the `vasic-digital/containers` submodule).** The contract is **fail-closed**: if the isolated runtime (or the containers submodule) is unavailable, the execution tier returns **SKIP-with-reason** — it NEVER falls back to host execution and NEVER auto-passes.

---

## 1. Threat model — what untrusted content the validation path handles

### 1.1 Attacker-controllable inputs (once P4 wires `Validate()`)
A skill is fully attacker-controllable via `POST /skills` (REST) and MCP `skill_create`. The validation path (`internal/validation/pipeline.go`) touches three classes of untrusted content:

| # | Untrusted input | Where it enters | Current handling | Real risk |
|---|---|---|---|---|
| T1 | **`skill.Content` markdown → every fenced code block** | `extractAndRunCode` → `extractCodeBlocks` (`pipeline.go:336-420`) extracts **every** ` ``` ` block *regardless of language tag* | Passed to `sandbox.Execute` | **Arbitrary host RCE** — an *illustrative* doc snippet (` ```bash\nrm -rf ~\n``` `) is executed |
| T2 | **`skill.Resources[].URL`** | `verifySingleResource` (`pipeline.go:259-303`) does server-side `HEAD`+`GET` on skill-supplied URLs | HEAD `<400` ⇒ pass; hash check best-effort, fail-open | **SSRF** (e.g. `http://169.254.169.254/…`), + fake "verified" (G21) |
| T3 | **`skill.Dependencies` / naming** | `CrossReference` (`pipeline.go:516-560`) | graph lookup only (no execution) | low — this is the correct static shape |

### 1.2 What the current "sandbox" actually does (forensic FACT — read in full)
- **Default path = host RCE.** Default `sandbox_type="wasm"` → `NewWASMSandbox`. `detectWasmRuntime()` looks for `wasmtime/wasmer/wasmedge`; on a bare host (and on THIS host — only `podman`+`docker` are present, no wasm runtime) it returns `""`. `Execute` then falls to `executeProcess`, which runs `python -c <code>` / `bash -c <code>` / `node -e <code>` **as the service UID on the host** (`sandbox.go:206-290`, dispatch `237-245`); the Go path is `go run` on the host (`sandbox.go:143-203`). The only "security" is a trimmed env + `LD_PRELOAD=` (which does **nothing** for network). No namespace, seccomp, cgroup, chroot, uid-drop, or network isolation.
- **"WASM" never runs WASM (G16).** `isWASMSupported` claims go/rust/c/cpp/ts, but `executeWASM` special-cases only Go → `executeGoSnippet` (host `go run`) and sends everything else to `executeProcess` (host). The detected wasm runtime binary is logged and never invoked. The source comment even admits *"In production, this would use a proper WASM compiler service"* (`sandbox.go:118-120`).
- **Every fallback fails OPEN to host exec or auto-pass:** Docker unavailable → `NewWASMSandbox` → host exec (`sandbox.go:338-342`); wasm runtime absent → host exec; `createDefaultSandbox` default/`gvisor` case → WASM → host exec (`pipeline.go:121-135`; note SPEC advertises `gvisor` but there is no `gvisor` case); `sandbox_type="none"` → `NoOpSandbox` returns `ExitCode:0` = **auto-PASS** (`pipeline.go:592-602`).
- **The DockerSandbox path is rootful + broken (G16, §11.4.161 violation):** uses the `docker` CLI (not rootless Podman); mounts `-v tmpDir:/tmp:ro` **and** `--tmpfs /tmp:noexec…` **and** `--read-only` — three conflicting `/tmp` declarations, so `go run /tmp/main.go` reads an empty/shadowed `/tmp`.

### 1.3 What "validating a skill" actually requires — execution vs static
The pipeline's execution "verdict" is `approved := execResult.ExitCode == 0` (`pipeline.go:322`). **That signal is near-worthless AND maximally dangerous:**
- Exit 0 proves only "the attacker's process didn't crash" — it says **nothing** about whether the skill's *knowledge* is correct.
- Most legitimate doc snippets are fragments / pseudo-code / illustrative and are *not meant to run standalone*; executing them is a category error.
- So the execution stage buys ~zero validation signal while opening a full RCE + SSRF surface. **The correctness oracle for a skill is the LLM jury + source verification + cross-reference — none of which execute untrusted code.**

**Decidable statically (NO execution of untrusted code) — covers the MVP's real validation goals:**
- **Syntax/parse validity** — via memory-safe, *non-executing* standard-library front-ends: Go `go/parser`+`go/format`; Python `ast.parse` / `compile(src, dont_inherit=True)` **without** `exec`/`eval`; `node --check`; JSON/TOML/YAML decoders. These *parse*, they do not *run*.
- **Schema / frontmatter / TOON-or-TOML structural validity.**
- **Dependency resolution** = graph existence lookup (already `CrossReference`) — **not** `pip/npm/go` install (those run arbitrary install scripts and are themselves execution).
- **Source/resource verification** = fetch + hash compare + SSRF-guarded egress (static; see §5, G21).
- **Type-check / lint** where a *non-executing* front-end exists.

**Genuinely requires execution — and is OUT of the MVP default:**
- "Compile and run this POC and assert on its *output*." R2's "*working POCs, not stubs*" is a requirement on the **project's own first-party reference code** (Gin/quic-go/tree-sitter/etc.), which is proven working in the project's own build+test/CI — **not** a requirement to run *untrusted skill-submitted* code. The R6 wizard's "validate via jury" is LLM-based. So no MVP requirement forces execution of untrusted skill code.

> **Decision (Q1):** **Static validation suffices for the MVP; untrusted-code execution is not required and is OFF by default.** Execution becomes an opt-in, fail-closed, isolated tier only if a real "observe POC output" feature is later justified.

---

## 2. Options for the execution-required subset (cited — §11.4.8)

For the *hypothetical* opt-in execution tier, these are the real isolation boundaries and their trade-offs **for this system** (a Go service that may run under `systemctl --user`, R15):

| Option | Isolation strength | Host requirements | Fit under `systemctl --user` (R15) | Untrusted arbitrary-lang snippets? | Verdict |
|---|---|---|---|---|---|
| **Host process (current)** | **none** — shared UID, FS, net | none | n/a | yes (this IS the RCE) | **REJECT — this is the defect** |
| **Rootless Podman** (§11.4.161) | good defense-in-depth: userns remap (container-root→unprivileged host UID) + default seccomp + drop-caps + `--network none` + `--read-only` + cgroup-v2 limits; breakout lands as an *unprivileged* user, not host root | rootless Podman + cgroup v2; **Constitution-mandated boundary** | **native fit** — Quadlet units in `~/.config/containers/systemd/`, managed by `systemctl --user` | yes (real container image per language) | **RECOMMEND (the boundary)** |
| **gVisor / runsc** | strong: userspace application kernel intercepts *every* syscall; host kernel attack surface dramatically reduced | KVM **or** ptrace/systrap platform + OCI-runtime wiring | awkward rootless/`--user`; not in containers submodule | yes | defense-in-depth **upgrade later** (run the Podman container *with* runsc), not MVP |
| **Firecracker microVM** | strongest: own guest kernel on KVM hardware virt | `/dev/kvm` + heavier orchestration | not feasible for an arbitrary unprivileged `--user` session | yes | **out of MVP** (future high-volume/high-risk ceiling) |
| **WASM (wasmtime/wazero)** | excellent, capability-based **deny-by-default** (no FS/net/clock unless granted) | none privileged; low overhead | native | **NO** — only runs code *compiled to WASM*; cannot run attacker `python -c`/`bash -c`/`go run` | good for *first-party* POCs compiled to WASM; **not** the boundary for untrusted skill snippets (this is exactly the G16 mismatch) |
| **nsjail (namespaces+cgroups+rlimits+seccomp-bpf)** | strong, purpose-built for "contestants upload+run arbitrary code" (Google CTF); sub-20ms, no daemon | namespace privileges | works but non-container dependency | yes | viable alternative, **subordinate to §11.4.161's Podman mandate** |
| **Refuse to execute (static-only)** | perfect (nothing runs) | none | native | n/a | **the MVP default** |

**Key evidence for the trade-offs:**
- **Containers share the host kernel; a permissive container running untrusted LLM-generated code is easily escaped** — microVMs are the gold standard, gVisor the middle ground, plain containers the minimum ([Northflank, *How to sandbox AI agents*, fetched 2026-07-15](https://northflank.com/blog/how-to-sandbox-ai-agents); [Northflank, *Firecracker vs gVisor*](https://northflank.com/blog/firecracker-vs-gvisor)).
- **Rootless Podman:** user namespaces remap container-root to an unprivileged host user, so *"even if an attacker escapes the container, they only gain access as an unprivileged user on the host"*; **caveat** — *"with user namespaces turned on, podman has access to kernel apis that have not been rigorously tested for non-root users"* ([Red Hat, *rootless Podman user namespace modes*, fetched 2026-07-15](https://www.redhat.com/en/blog/rootless-podman-user-namespace-modes); [Podman rootless tutorial](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md)).
- **Rootless Podman under `systemctl --user`:** Quadlet is the documented rootless path — unit files in `$HOME/.config/containers/systemd/`, managed via `systemctl --user daemon-reload/start/enable` + `journalctl --user`; requires cgroup v2 ([Podman `podman-systemd.unit(5)` docs, fetched 2026-07-15](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html)). This is why Podman — not gVisor/Firecracker — is the R15-feasible boundary.
- **gVisor:** *"an application kernel… it is an application kernel, not a virtual machine hypervisor or a system call filter"*; *"No system call is passed through directly to the host"*; **the sandbox process is itself treated as untrusted**; explicit residual gaps — *"gVisor does not provide protection against hardware side channels"* and *"a sandbox is not a substitute for a secure architecture"* ([gVisor security model, fetched 2026-07-15](https://gvisor.dev/docs/architecture_guide/security/); [gVisor intro](https://gvisor.dev/docs/architecture_guide/intro/)).
- **Firecracker:** each microVM boots its own guest kernel on KVM; *"to reach the host, an attacker would have to escape the guest kernel and break out of the VM"* — KVM-dependent, Linux-only, higher overhead ([Northflank, fetched 2026-07-15](https://northflank.com/blog/how-to-sandbox-ai-agents)).
- **WASM/wasmtime:** *"WebAssembly is inherently sandboxed by design since it must import all functionality"*; **deny-by-default capability model** — no FS/net/clock without an explicitly granted host capability ([Wasmtime security docs, fetched 2026-07-15](https://docs.wasmtime.dev/security.html); [Microsoft, *Wassette: WASM tools for AI agents*](https://opensource.microsoft.com/blog/2025/08/06/introducing-wassette-webassembly-based-tools-for-ai-agents/); [The New Stack, *WebAssembly sandboxing*](https://thenewstack.io/how-webassembly-offers-secure-development-through-sandboxing/)). **Limit:** WASM cannot execute arbitrary attacker `python -c`/`go run`; C-extension support (NumPy/etc.) is incomplete — so it is *not* a drop-in for the current "run any snippet" behaviour.
- **Language-level sandboxing of Python is fundamentally unworkable** (e.g. `__traceback__.tb_frame.f_globals['__builtins__']` escapes deleted builtins) — *"the solution isn't better prompts. It's isolation"* at the infrastructure level ([mavdol, *Notes on sandboxing untrusted code*, fetched 2026-07-15](https://gist.github.com/mavdol/2c68acb408686f1e038bf89e5705b28c)). This kills any "restrict the interpreter" shortcut and validates removing host execution outright.
- **nsjail:** Linux namespaces + cgroups + rlimits + seccomp-bpf; *"Google uses it for hosting CTF challenges, where contestants upload and execute arbitrary code"* ([google/nsjail, fetched 2026-07-15](https://github.com/google/nsjail)).

---

## 3. Recommended fail-closed design

### 3.1 Two-tier pipeline (rename kills the R11 naming bluff)
Replace the single `Sandbox` interface with two honestly-named tiers:

- **Tier A — `StaticValidator` (default, always-on, executes nothing untrusted).**
  Stages: (1) **resource verification** (SSRF-guarded, fail-closed hash — see §5); (2) **static code check** — parse/syntax/type-check via non-executing standard-library front-ends only (`go/parser`, Python `ast.parse`/`compile` w/o exec, `node --check`); **only** check blocks explicitly tagged runnable, and even then *parse-not-run* by default — never execute an illustrative doc block; (3) **cross-reference** (graph existence); (4) **LLM jury** (fail-closed, §4). This tier is buildable now with **zero** container dependency.

- **Tier B — `IsolatedExecutor` (opt-in, OFF by default, gated).**
  Exists only to *observe POC output* if a future feature justifies it. The **only** permitted implementation is a `PodmanIsolatedExecutor` backed by the `vasic-digital/containers` submodule (§11.4.161). Container hardening (all required): rootless Podman; ephemeral `--rm`; non-root in-container user; `--network none`; `--read-only`; distinct **`/work:ro`** mount for the code (NOT `/tmp`, and **no** `--tmpfs /tmp` shadow — this is the G16 fix) or code via stdin; `--cap-drop=ALL`; `--security-opt=no-new-privileges`; default seccomp profile; cgroup-v2 `--memory`/`--cpus`/`--pids-limit`; wall-clock timeout with container kill. Optional later defense-in-depth: run the same container under the `runsc` (gVisor) OCI runtime.

### 3.2 The fail-closed contract (explicit — no fail-open cell)
There is **no** code path from "isolation unavailable" or "verifier can't run" to either host execution or PASS.

| Condition | Old behaviour (the bug) | New behaviour (fail-closed) |
|---|---|---|
| Tier B requested, rootless Podman / containers submodule **absent** | fell back to `WASMSandbox` → **host exec** | **SKIP-with-reason** (`isolation_runtime_absent`); stage = SKIP, never PASS, never host-exec |
| `sandbox_type="none"` | `NoOpSandbox` → `ExitCode:0` = **auto-PASS** | **removed** — no auto-pass sink exists |
| `sandbox_type="wasm"` default, no wasm runtime | host exec | **removed** — default is Tier A static; no host exec anywhere |
| Docker/Podman unavailable | fell back to host exec | **SKIP-with-reason**, never host-exec |
| A parse/type-check tool for a language is missing | n/a | **SKIP-with-reason** for that block (honest coverage gap), never PASS-by-default |
| A required stage cannot produce a *real* verdict | mixed | stage = `BLOCKED`/`SKIP` with machine-readable reason |

**Overall verdict rule (composes G03 + G05):** a skill reaches `validated`/`active` **only** when every *enabled* stage produced a real *positive* verdict (incl. a recorded jury verdict). A `SKIP` in any stage NEVER silently upgrades to overall PASS — a pipeline whose only non-skipped stages passed while a mandatory stage skipped is `BLOCKED`, not `PASSED`. This is the same fail-closed posture the R1/G01 auth fix uses ("no keys ⇒ refuse", not "no keys ⇒ open").

### 3.3 Anti-bluff naming (R11) — deletions
- Delete `WASMSandbox`, `executeWASM`, `executeGoSnippet`, `executeProcess`, `detectWasmRuntime`, `NoOpSandbox`, the `LD_PRELOAD` "restrict network" line, and the `sandbox_type` `"wasm"`/`"none"` enum values + the SPEC `gvisor` claim (unimplemented). Nothing named "sandbox" may denote a host-exec path. Package doc-comment (`sandbox.go:1-4`) rewritten to describe static validation + an optional isolated executor.

---

## 4. G05 composition — fail-closed jury
`LLMJury` currently returns `Consensus:true` ("no jury configured, auto-approved") when `len(p.jury)==0` (`pipeline.go:428-439`) — an auto-pass in the *default* state. Under the same fail-closed principle: an empty jury while `validation.enabled` is a **hard BLOCK** (or forces `require_human_review`), never an auto-pass; require `approval_threshold ≥ 2` **real** votes. The static tier's "SKIP-not-PASS on inability to check" and the jury's "BLOCK-not-PASS on no jurors" are the *same* fail-closed rule applied at two stages, so the overall verdict has no fail-open sink.

---

## 5. G21 note — source verification hardened (static, in the default tier)
`verifySingleResource` is fail-open (HEAD `<400` passes; hash check only when a prior hash exists; GET/read errors return `nil`=pass). Fold into Tier A: require a stored hash for `official-doc`/`code` resources; treat any fetch/read/hash error as **verification FAILURE** (not pass); add an **egress allowlist / block link-local + cloud-metadata IPs** (169.254.0.0/16, ::1, RFC-1918 per policy) to close the SSRF (T2). This needs no execution runtime — it is pure static verification and ships with the MVP.

---

## 6. Honest boundary — now vs submodule-gated vs infeasible (§11.4.6/§11.4.112/§11.4.197)

**Buildable NOW (no submodule; this IS the MVP G02 fix — closes the RCE by construction):**
- Full Tier A static validator (parse/syntax check, SSRF-guarded fail-closed resource verify, cross-ref, fail-closed jury).
- Delete all host-exec paths + `NoOpSandbox` auto-pass + `LD_PRELOAD` line + `wasm`/`none`/`gvisor` sandbox enums; rename to `StaticValidator`.
- Define the `IsolatedExecutor` interface **with a default impl that always SKIPs-with-reason** (`isolation_runtime_absent`). After this, *no host-execution code exists anywhere in the tree* — the RCE is gone regardless of config, before P4 wires `Validate()`.

**Submodule-gated (tracked, NOT done — §11.4.197):**
- The concrete `PodmanIsolatedExecutor` requires the `vasic-digital/containers` submodule vendored at `submodules/containers/` (§11.4.28/§11.4.161/R7/R9) and wired through its `pkg/boot`/`pkg/compose`. **FACT (2026-07-15): the containers submodule is NOT present** — `project/` has no `submodules/` dir, and `constitution/submodules/` holds only `clickup_sync`, `continuum`, `session_orchestrator`, `token_optimizer`. Vendoring it + building `PodmanIsolatedExecutor` behind the fail-closed interface is a **tracked workable item** (must not sit as un-wired research). Until then Tier B is DISABLED (SKIP), which is the correct fail-closed state — so shipping the MVP without it loses no safety, only the (non-MVP) "observe POC output" capability.

**Structurally infeasible / deferred under `systemctl --user` (§11.4.112 — stated, not pretended):**
- **Firecracker microVM** needs `/dev/kvm` + orchestration an arbitrary unprivileged `--user` session generally lacks → NOT the MVP default; documented future ceiling, never claimed as shipped.
- **gVisor/runsc** needs KVM or a ptrace/systrap platform + OCI-runtime wiring that is awkward (not impossible) under pure rootless `systemctl --user`, and is not in the containers submodule → documented defense-in-depth *upgrade* (run the Podman container under runsc later), not MVP.
- **Rootless Podman + Quadlet** IS the documented, cgroup-v2 `systemctl --user` path — hence it is the chosen boundary: the recommended option is the one that is *feasible* under R15; the stronger boundaries are honestly deferred.
- Residual honesty: even Tier A parsers process attacker input — mitigated by non-executing memory-safe front-ends + input-size/time/recursion caps; a parser 0-day is a far smaller surface than deliberate execution but **not zero** (residual, flagged, not silently assumed safe). Rootless Podman shares the host kernel; a kernel-LPE/escape reaches the *unprivileged* service user, not host root (the Red-Hat "untested kernel APIs" caveat applies) — accept for MVP, note gVisor/microVM as the strengthen-later path.

---

## Sources verified (§11.4.8 / §11.4.99) — fetched 2026-07-15
1. Wasmtime — Security. https://docs.wasmtime.dev/security.html (primary; capability deny-by-default)
2. gVisor — Security Model. https://gvisor.dev/docs/architecture_guide/security/ (primary; sandbox-is-untrusted, residual gaps)
3. gVisor — Introduction to security. https://gvisor.dev/docs/architecture_guide/intro/ (primary)
4. google/nsjail — namespaces+cgroups+rlimits+seccomp-bpf; CTF arbitrary-code hosting. https://github.com/google/nsjail (primary)
5. Podman — `podman-systemd.unit(5)` / Quadlet under `systemctl --user`. https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html (primary)
6. Podman — rootless tutorial (user-namespace UID/GID remap). https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md (primary)
7. Red Hat — rootless Podman user-namespace modes ("untested kernel APIs" caveat). https://www.redhat.com/en/blog/rootless-podman-user-namespace-modes (vendor)
8. Microsoft OSS — Wassette: WASM-based tools for AI agents (deny-by-default). https://opensource.microsoft.com/blog/2025/08/06/introducing-wassette-webassembly-based-tools-for-ai-agents/ (vendor)
9. Northflank — How to sandbox AI agents: microVMs, gVisor & isolation. https://northflank.com/blog/how-to-sandbox-ai-agents (secondary)
10. Northflank — Firecracker vs gVisor. https://northflank.com/blog/firecracker-vs-gvisor (secondary)
11. mavdol — Notes on sandboxing untrusted code (why Python can't be sandboxed; F'cracker/gVisor/WASM trade-offs). https://gist.github.com/mavdol/2c68acb408686f1e038bf89e5705b28c (secondary)
12. The New Stack — How WebAssembly offers secure development through sandboxing. https://thenewstack.io/how-webassembly-offers-secure-development-through-sandboxing/ (secondary)

**Negative findings (§11.4.99(B)):** No source supports language-level sandboxing of Python/JS as sufficient — the opposite (source 11). No source supports "container = safe for untrusted code" unqualified; every source frames plain containers as the *minimum*, hardened-rootless-Podman as adequate-with-caveats, gVisor/microVM as stronger. No authoritative source describes running Firecracker/gVisor under an unprivileged `systemctl --user` session as a turnkey path — corroborating the §11.4.112 deferral. Constitution `vasic-digital/containers` submodule content was NOT web-researched (internal); its `pkg/boot`/`pkg/compose` API is asserted from §11.4.76/§11.4.161 references and MUST be confirmed against the submodule when it is vendored (tracked).
