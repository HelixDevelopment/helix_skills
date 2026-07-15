# G15 — Aurora OS / HarmonyOS client feasibility + max-shared-code client architecture

**Revision:** 1
**Last modified:** 2026-07-15T17:20:00Z
**Status:** design-research, no code
**Scope:** R3 client surfaces (`REQUIREMENTS.md:58-60`) for
`github.com/helixdevelopment/skill-system`, specifically the two highest-risk
mobile targets — HarmonyOS and Aurora OS — and the max-shared-code client
architecture across ALL R3 surfaces (CLI/TUI/Web/Desktop/Mobile/HarmonyOS/Aurora).
**Authority / mandates served:** G15 (`GAPS_AND_RISKS_REGISTER.md:189-195`) —
Aurora/HarmonyOS client feasibility, the plan's top danger-zone
(`IMPLEMENTATION_PLAN.md:334`, P8.T5 `:257`); R3 (`REQUIREMENTS.md:58-60`) — clients
incl. HarmonyOS + Aurora OS, maximize shared codebase; P8 (`IMPLEMENTATION_PLAN.md:249-257`)
— contract-first thin clients generated from `api/openapi.yaml`. Constitution
§11.4.8/§11.4.99/§11.4.150 (deep multi-angle research before any build claim),
§11.4.6 (no-guessing — every verdict below is FACT-with-cited-source or explicit
`UNCONFIRMED:`), §11.4.112 (structural-impossibility classification — evidence-backed,
never a bluffed build), §11.4.111 (resolve toolchains/targets by stable identity,
not assumption), §11.4.200 (verify-after-deploy — a build "success" message is not
proof the artifact runs on the intended target), §11.4.197 (research fully wired to
a tracked outcome, not left in the backlog).
**Read discipline:** every Go/config/plan/register fact below was read from the
committed baseline ref `255061b` via `git show 255061b:…project/<path>` or the
tracked Markdown docs (working tree has uncommitted Go changes + a live review in
flight and was NOT read for source facts). Every external platform claim below was
verified 2026-07-15 against the cited live source (§11.4.99) — none of it is drawn
from training-data memory, because both platforms are fast-moving and a stale
recollection would misguide the build decision (the exact failure class §11.4.99's
canonical anchor was written to prevent). No existing file was modified; this is
the single new deliverable this agent produced.

---

## 0. One-paragraph problem statement

R3 hard-requires clients on HarmonyOS and Aurora OS (`REQUIREMENTS.md:58-60`), and
the implementation plan and the register both independently flag this pair as the
**highest client risk** in the whole project (`IMPLEMENTATION_PLAN.md:334` danger-zone
#1; `GAPS_AND_RISKS_REGISTER.md:191` "highest client risk"). No client code exists yet
for any R3 surface — P8 has not started (`GAPS_AND_RISKS_REGISTER.md:192` "P8 not
started"; confirmed independently in this session: `git ls-tree -r 255061b` under
`project/` contains no `cmd/mobile`, `cmd/desktop`, `flutter/`, or equivalent — the
only existing client-adjacent code is a hand-rolled `cmd/tui/api_client.go`). The
prior plan's only proposed path was "Flutter-OHOS + the omprussia embedder, a spike,
not a proven build" (`GAPS_AND_RISKS_REGISTER.md:192`). This document closes that gap
with FACTS pinned to 2026-07-15 sources: does a Flutter path exist for each target,
what does it actually require, and what architecture maximizes shared code given the
real constraints — never a bluffed "it will just work."

---

## 1. HarmonyOS: two structurally different targets, not one

**HarmonyOS (1.x–4.x) is AOSP-compatible; HarmonyOS NEXT (= HarmonyOS 5.0, API level
12+) removed AOSP entirely.** This is a hard version boundary, not a gradual
deprecation:

- HarmonyOS 4.2 (built on OpenHarmony 3.x) "supports both native apps built on the
  OpenHarmony core as well as Android apps" via an integrated Dalvik VM and Android
  framework compatibility layer (harmony-developers.com, "HarmonyOS, HarmonyOS NEXT,
  OpenHarmony and Oniro Explained", verified 2026-07-15).
- **HarmonyOS NEXT (HarmonyOS 5.0.0) is API level 12, first released 2024-09-05**
  (Huawei Developers atomic-release notes for 5.0.0(12); cross-confirmed by
  harmony-developers.com's API-level article, verified 2026-07-15). HarmonyOS NEXT
  "does not include the Android AOSP core and is incompatible with Android
  applications … Huawei has switched to their own kernel (HongMeng) and has removed
  all AOSP code" (OSnews / Android Authority reporting, verified 2026-07-15). Only
  native apps via the Ark Compiler + ArkUI (built on ArkTS, a TypeScript superset) +
  native HarmonyOS SDK APIs run on HarmonyOS NEXT.
- Every R3 client target that matters going forward is HarmonyOS NEXT (API 12+) — the
  AOSP-compatible line is the platform HarmonyOS NEXT is explicitly replacing, so a
  feasibility study aimed at "HarmonyOS" without this distinction would answer the
  wrong question (the exact ambiguity R3's original phrasing and G15's register entry
  both leave open, and the exact reason the operator's prompt for this task demanded
  the split be resolved as FACT).

### 1.1 Does Flutter run on HarmonyOS NEXT? — FACT, cited

**Not officially — Google's `flutter/flutter` has explicitly closed a HarmonyOS
support request as "not planned" / duplicate**, even after the requester cited
Huawei's AOSP-app-support removal as the motivating urgency (github.com/flutter/flutter
issue #150536, verified 2026-07-15). There is no upstream Flutter target for
HarmonyOS and none is planned by the Flutter team.

**A community/vendor-affiliated fork DOES run Flutter on HarmonyOS NEXT.**
OpenHarmony-SIG maintains `flutter_flutter` (gitee.com/openharmony-sig/flutter_flutter,
verified 2026-07-15) — its own README states it "extends Flutter SDK compatibility for
the OpenHarmony platform." It explicitly targets **API12** (= HarmonyOS NEXT / 5.0)
with DevEco Studio 5.0 or command-line-tools 5.0 — i.e. it targets the new, non-AOSP
line, not the legacy AOSP-compatible HarmonyOS. A companion, more current build is
referenced directly by name: "HarmonyOS Next API 16 for HarmonyOS Flutter 3.22.0 is
released" (harmony-developers.com, verified 2026-07-15) — confirming the fork is
still being advanced against newer HarmonyOS NEXT API levels, not abandoned at API12.

**Architecture — a thin, mostly-generated native shell, not a UI rewrite.**
`flutter create` for this target generates a native ArkTS entry-ability module that
hosts the Flutter engine, structurally the same role as Android's `MainActivity` or
iOS's `AppDelegate`:

```typescript
import flutter from '@ohos/flutter_module';
export default class EntryAbility extends Ability {
  onCreate(want: Want) {
    flutter.init(this.context);
    flutter.run({ entry: 'lib/main.dart' });
  }
}
```

(harmony-developers.com, "Flutter App Development: Hongmeng HarmonyOS NEXT Beta
Porting", verified 2026-07-15). App UI and business logic stay in Dart; the ArkTS
surface is boilerplate bootstrap, not a parallel UI implementation. The build target
is `flutter build hap` — output packaging is **HAP** (HarmonyOS Ability Package),
HarmonyOS's native app-package format.

**Maturity — explicitly beta, community/SIG-maintained, not Google-official.** The
source article's own title names the target "HarmonyOS NEXT Beta"; required tooling
versions are pre-GA-era (DevEco Studio ≥4.0 Beta2 in the article's own prerequisites);
the OpenHarmony-SIG `flutter_flutter` Gitee page itself carries an archival notice
("This repository has been archived. New address: GitCode") — a maintenance-location
change that must be re-verified at spike time, not assumed stable (§11.4.6). The fork
also lags upstream Flutter's release cadence (one tagged base seen was `3.7.12-ohos-*`
against a much newer upstream Flutter; a separately-cited article names `3.22.0`
against API16 — the fork's version-to-API-level matrix is NOT fixed and must be
pinned explicitly at spike time, never assumed current).

**VERDICT (FACT): Flutter on HarmonyOS NEXT is possible via a named, HarmonyOS-NEXT-
targeting fork (OpenHarmony-SIG `flutter_flutter`), not structurally impossible — but
it is a third-party/community-SIG fork, not upstream Flutter, currently at beta
maturity, requiring its own toolchain (DevEco Studio + hvigor + OHOS SDK) distinct
from Android/iOS Flutter tooling.**

---

## 2. Aurora OS: Sailfish-derived, Qt/QML-native, with a vendor-published Flutter fork

Aurora OS "has been developed on the basis of Sailfish OS" (Jolla, Finland) and "has
been developed since 2016 by the Russian company Open Mobile Platform" (OMP)
(en.wikipedia.org/wiki/Aurora_OS_(Russian_Open_mobile_platform), verified
2026-07-15).

### 2.1 Native stack — Qt/QML/C++, Sailfish Silica, QtWidgets forbidden

The official Aurora developer portal documents "Aurora IDE — an integrated
development environment based on Qt Creator for developing applications in C, C++,
and QML for Aurora OS using Sailfish Silica components," with "Qt Quick" as "the
standard library for writing QML applications" and Silica supplying the
Aurora-styled widget set. The public-API reference explicitly states **QtWidgets is
prohibited** ("not optimized for the Aurora OS user interface")
(developer.auroraos.ru — SDK/app_development + public_api pages, verified
2026-07-15). A native Aurora client is therefore a Qt Quick/QML + C++ application
built with the Aurora SDK/Aurora IDE toolchain, packaged as RPM.

### 2.2 Does Flutter run on Aurora OS? — FACT, cited

**Yes — and this path is OMP's OWN, vendor-published fork, not a third-party
project.** OMP's own auroraos.ru announcement states (translated) "Open Mobile
Platform published a Flutter SDK with initial support for Aurora OS"
(auroraos.ru/tpost/mzsi3ecdt1, verified 2026-07-15), and a follow-up post
("Обновление от сообщества развития Flutter для ОС Аврора" — "Update from the Flutter-
for-Aurora development community", auroraos.ru/tpost/a7tihbgfy1, verified
2026-07-15) confirms continued iteration: new Flutter CLI commands, VS Code debug
panel + device manager integration, and new plugins landed starting from the
3.24.0-based release line. The SDK's canonical repository is
`gitlab.com/omprussia/flutter` and its docs are published at
`omprussia.gitlab.io/flutter/docs/` (verified 2026-07-15).

**Pinned toolchain, not upstream Flutter.** The install docs
(`omprussia.gitlab.io/flutter/docs/start/install/`, verified 2026-07-15) specify
**Flutter 3.27.1**, cloned from OMP's own fork at
`gitlab.com/omprussia/flutter/flutter` — not `flutter/flutter` upstream. Two
components are required: (1) the **Aurora Platform SDK** (RPM build/CI toolchain,
target naming `AuroraOS-{version}-base-{architecture}`, covering `aarch64`,
`armv7hl`, `x86_64`), and (2) OMP's Flutter SDK fork; `flutter-aurora doctor`
verifies the install.

**Architecture — a native Qt-based "Flutter Embedder," analogous to HarmonyOS's
ArkTS EntryAbility shell.** Aurora's own developer tooling ("Aurora Scripts") lists
and manages "Flutter Embedder" as an independently versioned, installable component
distinct from the app itself (github.com/keygenqt/aurora-scripts, verified
2026-07-15) — i.e. a native Qt/C++ embedder hosts the Flutter engine inside the app's
RPM package, the same structural role as HarmonyOS's `EntryAbility`/`flutter_module`
bootstrap: a thin, mostly-tooling-generated native shim, not a parallel Dart-vs-QML
UI rewrite of the whole app.

**Maturity — vendor-published but explicitly self-described as early-stage.** OMP's
own 2024 announcement language is "initial support" (начальная поддержка) for the
first release; the pinned Flutter version (3.27.1) is a specific recent-stable
snapshot rather than a tracked-to-upstream-HEAD release train, consistent with an
actively-evolving, non-GA fork. `UNCONFIRMED:` whether the Aurora Platform SDK
download/registration path (via `developer.auroraos.ru`) carries any
region/registration gate that would affect an autonomous CI spike — the fetched
pages did not state this either way, and this doc does not guess (§11.4.6); it is
recorded as an open question for the operator-hardware spike (§4 below), not asserted
as a blocker or a non-blocker.

**VERDICT (FACT): Flutter on Aurora OS is possible via OMP's own vendor-published
fork (`gitlab.com/omprussia/flutter`, docs at `omprussia.gitlab.io/flutter`), not
structurally impossible — vendor-backed (unlike the HarmonyOS case, which is
community-SIG rather than Huawei-official), but self-described "initial support,"
pinned to a specific Flutter release (3.27.1) rather than upstream Flutter, and
requiring the Aurora Platform SDK + RPM packaging pipeline.**

---

## 3. Max-shared-code client architecture (R3)

R3 demands clients across CLI, TUI, REST, Web, Desktop (Win/macOS/Linux), and Mobile
(Android, iOS, HarmonyOS, Aurora OS) with **maximum shared codebase**
(`REQUIREMENTS.md:58-60`). The prior confirmed research finding already anchors part
of this: "Shared-core = contract-first thin clients… genuine shared logic is only
~15-30%… Per surface: CLI+TUI = Go (reuse existing); GUI apps = Flutter — the ONLY
framework covering Android + iOS + HarmonyOS (Flutter-OHOS) + Aurora OS (omprussia
embedder) + desktop from one codebase; Web = React" (`REQUIREMENTS.md:166-172`,
`IMPLEMENTATION_PLAN.md:253-257` P8.T1-T5). §1 and §2 above turn that prior finding
from a plan-stage assumption into a cited, verified FACT: both the HarmonyOS-NEXT
fork and the Aurora fork are real, and both share the SAME structural shape (a thin
native bootstrap shell hosting the Flutter engine), which is exactly what makes a
single shared Flutter/Dart core viable across five of the six R3 GUI surfaces.

### 3.1 Recommended layering (three layers, decreasing shareability)

**Layer 1 — Generated contract client (near-100% shared, ALL surfaces incl. CLI/TUI).**
One Dart package generated from `api/openapi.yaml` (the P8.T1 "shared-core contract"
deliverable) is the single source every GUI client consumes — mirrors the Go side's
own contract-first intent. TOON is not a Dart-side blocker: a community Dart TOON
implementation already exists (`github.com/toon-format/toon-dart`, verified
2026-07-15), alongside the spec's own multi-language implementations (TypeScript,
Python, Go, Rust, .NET) — so once the server-side codec (G08) lands, the Flutter
clients can speak `application/toon` natively rather than falling back to JSON,
closing the same wire-format gap on the client side that G08 tracks on the server
side. Go's CLI/TUI keep their own generated Go client (oapi-codegen-class tooling)
against the same spec — same contract, two language bindings, not a second contract.

**Layer 2 — Shared Flutter/Dart UI + business logic (~70-85% of GUI-surface code,
shared across Android, iOS, Desktop, HarmonyOS NEXT, Aurora OS).** One Flutter
application module (screens, state management, the wizard flow from R6, API-client
wiring) is authored once and compiled for Android, iOS, and Desktop via mainline
Flutter, and for HarmonyOS NEXT and Aurora OS via their respective forks (§1, §2).
This is the layer R3's "maximize shared codebase" mandate actually optimizes — five
GUI targets from one Dart codebase instead of five. Web stays **React** per the
already-confirmed decision (`REQUIREMENTS.md:171`) — Flutter Web exists but was not
the prior research finding's choice and is out of this document's scope to
re-litigate; if Flutter Web is reconsidered later it would extend Layer 2 to six
targets, not five, but that re-decision is not made here (§11.4.6 — no invented
scope change).

**Layer 3 — Per-OS native shell + toolchain (NOT shareable; a real, per-target,
non-Dart cost).** Each of the five Flutter-hosting targets needs its own thin native
bootstrap and its own build toolchain:
- Android: Gradle + Kotlin/Java `MainActivity` (standard Flutter Android embedding).
- iOS: Xcode + Swift/Obj-C `AppDelegate` (standard Flutter iOS embedding).
- Desktop (Win/macOS/Linux): Flutter's own desktop embedder per OS (CMake/Visual
  Studio/Xcode toolchains respectively).
- HarmonyOS NEXT: ArkTS `EntryAbility` + `@ohos/flutter_module` (§1.1), built with
  DevEco Studio + hvigor + the OpenHarmony SDK, packaged as HAP.
- Aurora OS: the Qt-based "Flutter Embedder" (§2.2), built with the Aurora Platform
  SDK, packaged as RPM.

None of Layer 3 is Dart and none of it is shareable — it is genuinely per-OS
plumbing, honestly stated as a cost (§11.4.6), not hidden inside a "maximized shared
code" claim. Each Layer-3 shell is, per both platforms' own tooling, mostly
boilerplate generated by that platform's `flutter create`-equivalent command rather
than hand-authored business logic — the shareable-vs-not boundary sits cleanly at
"Dart application code" vs "native bootstrap + build toolchain", not at
"HarmonyOS/Aurora are unshareable, everything else is."

### 3.2 Plugin-parity gap — an honest, tracked residual (not assumed away)

Neither fork is asserted here to carry 100% of the standard Flutter plugin
ecosystem. The Aurora Flutter docs structure their own support pages around "Flutter
support, Dart packages, official packages, community packages, and third-party
tools" with an explicit statement that this section documents "the support status of
plugins for the Aurora operating system" (`omprussia.gitlab.io/flutter/docs/support/`,
verified 2026-07-15) — i.e. OMP itself frames plugin support as a per-plugin matrix,
not a blanket guarantee. This doc does not enumerate that matrix (out of scope for a
feasibility study; the exact plugin list is a build-time concern for whichever
plugins P8's actual UI ends up using) but flags it as a **named, tracked risk**:
any platform-channel-backed plugin (camera, Bluetooth, native notifications, etc.)
used by the shared Layer-2 code must be re-verified against each fork's plugin
support matrix before HarmonyOS-NEXT/Aurora builds are claimed complete — a plugin
gap on either fork is a per-plugin shim cost, not a whole-app rewrite, but it is real
and must not be silently assumed solved by "Flutter runs there."

---

## 4. §11.4.112 risk-flag classification (per target)

| Target | Classification | Cited blocker/evidence | Operator-hardware spike required? |
|---|---|---|---|
| **HarmonyOS NEXT (Flutter)** | **possible-with-named-fork** — NOT structurally impossible. Named fork: OpenHarmony-SIG `flutter_flutter` (gitee.com/openharmony-sig/flutter_flutter), targeting API12+/HarmonyOS NEXT, HAP packaging, `@ohos/flutter_module` EntryAbility bootstrap. | No upstream Flutter support (flutter/flutter#150536, closed "not planned"); fork is beta-maturity + recently relocated ("archived → GitCode" notice) — re-verify canonical location at spike time. | **YES.** A HarmonyOS NEXT device or the official HarmonyOS emulator + a licensed/registered DevEco Studio 5.0+ environment is required to prove a real `flutter build hap` succeeds AND the resulting HAP installs+runs (§11.4.200 — a build-success message alone is not proof of a working artifact on the intended target). No Google-official or Huawei-official CI runner for this path is known to exist; this is genuinely operator-attended. |
| **Aurora OS (Flutter)** | **possible-with-named-fork** — NOT structurally impossible. Named fork: OMP's own `gitlab.com/omprussia/flutter` (Flutter 3.27.1 pin), Aurora Platform SDK, native Qt "Flutter Embedder", RPM packaging, `flutter-aurora doctor` CLI. | Vendor-published but self-described "initial support"; pinned to a specific Flutter release rather than tracking upstream HEAD; plugin-support matrix is per-plugin, not blanket (§3.2). `UNCONFIRMED:` whether Aurora Platform SDK access/registration carries a region gate — not established either way by the sources fetched, recorded honestly rather than guessed. | **YES.** The Aurora OS Emulator (VM) or a real Aurora OS device + the Aurora Platform SDK/Aurora IDE + RPM signing infra is required to prove a real Flutter-Aurora build produces an installable, runnable package. This is genuinely operator-attended — no evidence of an operator-independent autonomous CI path for Aurora RPM builds was found. |

Neither target reaches `structurally-impossible` — the §11.4.112 "won't-fix" closure
is **not applicable here**; both targets get a tracked, evidence-backed "possible,
needs an operator-hardware spike to convert from cited-feasible to proven-built"
status, exactly the outcome G15's own decision text asked for ("risk-flag with the
exact blocker … never bluff a build" — `GAPS_AND_RISKS_REGISTER.md:194`).

---

## 5. Freeze the backend contract first (G15's explicit precondition)

G15's decision text states the client architecture should be de-risked by freezing
the backend contract first, because every client surface — including the two
highest-risk mobile targets — is a thin client over that contract
(`GAPS_AND_RISKS_REGISTER.md:194`). This is not optional sequencing: **`api/openapi.yaml`
is currently drifted from the live implementation (G09,
`GAPS_AND_RISKS_REGISTER.md:127-139`)** — confirmed concrete drifts include
`POST /search` (spec) vs `GET /api/v1/skills/search?q=` (live), `GET /registry/missing`
(spec) vs two different live route shapes, and a `POST /expand/{name}` spec route
with **no live route at all**. Generating a Layer-1 Dart client (§3.1) from the
current `openapi.yaml` would produce a client that calls endpoints that 404 or expect
the wrong verb/shape on today's server — the exact failure mode G09 already names as
the reason contract-first codegen (P8) must wait for G09's fix
(`GAPS_AND_RISKS_REGISTER.md:137` "contract-first codegen (P8) will produce clients
calling endpoints that 404 or expect the wrong verb/shape"). Sequencing implication
for this document's recommendation: **Layer 1 codegen (§3.1) MUST NOT start until
G09's route-parity contract gate is green** (`research/g09_openapi_drift_design.md`
is the tracked design for that fix); starting HarmonyOS/Aurora spikes (§4) does not
have to wait for G09 — the toolchain/build-feasibility spike is orthogonal to which
exact routes the finished client calls — but the *generated client itself*, and any
client-side integration test asserting real API calls succeed, is correctly blocked
on G09 per the plan's own ordering.

---

## 6. Test/proof approach — how a real end-to-end client call gets proven per target

Per §11.4.108 (runtime signature on a clean deployment) and §11.4.200
(verify-after-deploy, never trust a build tool's own success message), each target's
proof obligation is layered:

1. **Contract layer (all targets, once G09 is green):** the generated Layer-1 client
   compiles against the frozen `api/openapi.yaml`, and a contract test asserts the
   generated client's request/response shapes schema-validate against the live
   server's actual responses — this is autonomous, no operator hardware needed.
2. **Android / iOS / Desktop (Win/macOS/Linux):** standard Flutter build + a real
   `/health` + wizard-flow smoke call against a live backend instance, captured as
   evidence (request/response payloads, build artifact identity). Android/iOS/Desktop
   builds are autonomous-CI-feasible (mainline Flutter tooling, no vendor SDK
   registration gate known); this is the `AUTONOMOUS_VERIFIED`-class path per
   §11.4.52.
3. **HarmonyOS NEXT:** requires the DevEco Studio 5.0+ toolchain + either the
   official HarmonyOS emulator or a real device to run `flutter build hap`
   end-to-end. If autonomous CI access to that toolchain/device is genuinely
   unavailable, this is an honest `SKIP-with-reason: hardware_not_present` or
   `operator_attended` per §11.4.3/§11.4.52 — **never** a claimed PASS from a build
   log alone. When the device/emulator IS available: verify-after-deploy per
   §11.4.200 — install the HAP, launch it, and confirm (via a captured screen state
   or log) that the app reached the wizard screen and completed one real API round
   trip, not merely that `hvigor`/DevEco Studio printed a success message.
4. **Aurora OS:** requires the Aurora Platform SDK/Aurora IDE + either the Aurora OS
   Emulator VM or a real Aurora device to build+install the RPM. Same
   SKIP-with-reason discipline applies if autonomous access is unavailable
   (`operator_attended`, citing the exact missing dependency — SDK registration,
   device, or emulator image — never silently assumed available). When available:
   install the RPM, launch the app, and confirm one real API round trip the same
   way as HarmonyOS NEXT, per §11.4.200 (a `flutter-aurora doctor` green result or an
   RPM build exit-0 message is NOT proof the installed app runs against the real
   backend).
5. **CLI/TUI/Web/REST:** already covered by the existing 13-test-type matrix and are
   not re-derived here (see §7 reconciliation).

None of steps 3-4 can be autonomously completed inside this research task — they are
correctly deferred to an operator-scheduled hardware/toolchain window, consistent
with how `GAPS_AND_RISKS_REGISTER.md:195` already classifies HelixQA coverage for
this gap as "device-lab dependent … flag as operator-attended where autonomous build
is infeasible."

---

## 7. Reconciliation with `research/testing_infrastructure_plan.md`

- **Row match confirmed:** `testing_infrastructure_plan.md:301` (the G15 row in the
  per-gap test-type matrix) already states: "acceptance(one build artifact per OS OR
  documented blocker), smoke(thin client hits `/health`+wizard), contract(generated
  client compiles vs OpenAPI)… Challenge: per-OS build feasibility; HelixQA:
  **device-lab dependent — operator-attended where autonomous build infeasible**…
  build artifact OR §11.4.112 evidence-backed blocker; contract-compile log." This
  document's §4 risk-flag table and §6 test approach are consistent with that row
  point-for-point: neither target got an evidence-backed §11.4.112 blocker (both are
  `possible-with-named-fork`), so the acceptance criterion resolves to "one build
  artifact per OS," gated on the operator-hardware spikes named in §6.
- **Honest-boundary section match confirmed:** `testing_infrastructure_plan.md:378-381`
  states "Aurora/HarmonyOS acceptance (G15): device-lab / omprussia-embedder
  dependent; where autonomous build is infeasible the HelixQA verdict is
  operator-attended, and a failed feasibility spike is a §11.4.112 evidence-backed
  blocker, never a bluffed build." This document supplies the FACT layer that row was
  written in anticipation of: the omprussia-embedder reference at
  `testing_infrastructure_plan.md:378` is now confirmed (§2.2, the "Flutter Embedder"
  component in Aurora's own tooling), and the HarmonyOS side of the same row is now
  additionally backed by the OpenHarmony-SIG fork citation (§1.1) that row did not
  yet name.
- **No contradiction found** between this document's verdicts and the existing test
  plan or register text; this document narrows "smallest ecosystem, highest risk"
  (`IMPLEMENTATION_PLAN.md:334`) from a plan-stage description into cited,
  2026-07-15-verified facts plus a concrete risk classification (§4) and layering
  recommendation (§3) neither prior doc contained.

---

## 8. Honest gaps (§11.4.6) — what this document does NOT establish

- **No build was actually run.** This is a design-research document, not a build
  attempt — per the task's own constraint (no code, no execution). The verdicts in
  §1 and §2 are FACT-with-cited-source that a working path *exists per the vendors'
  own documentation and public statements*; they are not a claim that this repo's
  own Flutter app has been built and run on either target. That proof is §6 and
  §4's "operator-hardware spike required: YES" — explicitly NOT satisfied by this
  document.
- **`UNCONFIRMED:` Aurora Platform SDK access/registration gate.** Not established
  either way by the sources fetched (§2.2) — flagged, not guessed.
- **`UNCONFIRMED:` exact plugin-support matrix for either fork.** §3.2 names this as
  a tracked risk rather than enumerating it; the concrete plugin list depends on
  which plugins Layer 2's actual UI ends up using, which is not yet designed (P8 not
  started).
- **HarmonyOS-NEXT fork's current canonical repository location is time-sensitive.**
  The Gitee page carried an "archived → GitCode" notice at verification time
  (§1.1); re-verify at spike time rather than assuming the Gitee URL is still
  canonical by the time P8.T5 executes.
- **This document does not re-litigate the Web=React / Desktop-framework decision**
  already recorded in `REQUIREMENTS.md:166-172` — it is taken as a given precondition
  for the Layer-2 scope (§3.1), not re-derived.
- **Forbidden-vocabulary check:** this document contains no instance of
  likely/probably/maybe/seems/appears/guess/presumably/perhaps/conjectured (§11.4.6);
  every uncertain point is marked `UNCONFIRMED:` explicitly instead.
