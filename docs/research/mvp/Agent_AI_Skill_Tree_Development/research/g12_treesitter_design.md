# G12 — Real Tree-Sitter Parsing Design (Kotlin/C# Coverage + Explicit Reduced-Fidelity Fallback)

**Revision:** 1
**Last modified:** 2026-07-15T16:30:00Z
**Status:** design-research, no code landed
**Scope:** `internal/codeanalysis/treesitter.go` (+ `internal/config/config.go` defaults)
**Baseline read via:** `git show 255061b:docs/research/mvp/Agent_AI_Skill_Tree_Development/project/<path>` (constitution §11.4.84/§11.4.119 — the live working tree is under concurrent mutation and was never read)

---

## 1. Root cause (proved against committed baseline `255061b`)

All quotes below are FACT, read from the pinned commit `255061b` via `git show` —
never from the concurrently-mutating working tree.

1. **`initNativeParser` always fails, unconditionally.**
   `internal/codeanalysis/treesitter.go:106-131` — the whole body is a comment
   block describing what a real CGO implementation *would* do, followed by:
   ```go
   // Since we cannot rely on CGO being available, we always return an error
   // and use the regex fallback. ...
   return fmt.Errorf("native parser not available for %s (CGO may be disabled)", language)
   ```
   (`treesitter.go:130`). This is unconditional — it does not probe whether
   CGO is actually enabled on the running build; it always returns an error
   regardless of the true toolchain state.

2. **`NewTreeSitterParser`'s native-attempt list omits Kotlin and C# entirely.**
   `treesitter.go:90`:
   ```go
   for _, lang := range []string{"go", "python", "java", "javascript", "c", "cpp", "rust"} {
       if err := p.initNativeParser(lang); err != nil {
           p.fallbackParsers[lang] = newRegexParser(lang)
       } else {
           p.nativeParsers[lang] = true
       }
   }
   ```
   Kotlin and C# are not in this 7-language list, so neither `p.nativeParsers`
   nor `p.fallbackParsers` is ever pre-populated for them at construction time.
   `GetSupportedLanguages()` (`treesitter.go:566-575`) and
   `IsLanguageSupported()` (`treesitter.go:578-581`, body
   `return p.nativeParsers[language] || p.fallbackParsers[language] != nil`)
   therefore report Kotlin/C# as **unsupported**, even though Kotlin is a
   default-enabled analysis language (fact 6 below).

3. **`parseNative` is a placeholder that always errors.**
   `treesitter.go:157-160`:
   ```go
   func (p *TreeSitterParser) parseNative(content []byte, language string) (*Tree, error) {
       // This would call into the tree-sitter C library via CGO.
       // Placeholder implementation - always returns error to trigger fallback.
       return nil, fmt.Errorf("native parser not implemented")
   }
   ```
   Combined with fact 1/2, `p.nativeParsers[language]` is never `true` for
   *any* language, so `parseNative` is never even reached from `Parse()`
   (`treesitter.go:143-154`) in the shipped binary — it is 100% dead code,
   and so are `Tree.Root` / `TSNode` (never populated by any code path).

4. **All three native extraction methods are unimplemented placeholders.**
   `treesitter.go:228-230` (`extractImportsNative`), `233-235`
   (`extractFunctionsNative`), `238-240` (`extractClassesNative`) — each is a
   one-line stub returning `fmt.Errorf("native <kind> extraction not
   implemented")`. They are dead code for the same reason as fact 3.

5. **`compilePatterns` has no `kotlin` and no `csharp` case.**
   `treesitter.go:264-296` — the switch covers exactly `"go"` (266),
   `"python"` (273), `"java"` (277), `"javascript", "typescript"` (282),
   `"c", "cpp"` (286), `"rust"` (290). There is no `case "kotlin":` and no
   `case "csharp":` anywhere in the file. `RegexParser.patterns` therefore
   stays an **empty map** for either language.

6. **Kotlin is a default-enabled analysis language; C# is a normalized,
   operator-reachable one.**
   `internal/config/config.go:204`:
   ```go
   Languages:       []string{"java", "kotlin", "c", "cpp", "python", "go"},
   ```
   Kotlin ships enabled **by default**. `normalizeLanguage`
   (`treesitter.go:541-563`) maps `"kt"` → `"kotlin"` (`treesitter.go:558-559`)
   and `"c#"`, `"csharp"` → `"csharp"` (`treesitter.go:554-555`) — so an
   operator adding `"csharp"` to `CodeAnalysisConfig.Languages` reaches the
   exact same empty-pattern path (`isLanguageEnabled` in `analyzer.go`
   string-compares via `normalizeLanguage`, so `"csharp"` is a valid,
   accepted configuration value).

7. **The failure is SILENT — no error, no log, no signal.**
   `parseFallback` (`treesitter.go:163-175`) never returns an error; it always
   returns `&Tree{..., Parsed: parsed}` where `parsed = parser.Parse(content)`.
   `RegexParser.Parse` (`treesitter.go:298-308`) calls
   `extractImports`/`extractFunctions`/`extractClasses`, each of which
   `switch`es on `p.language` (`treesitter.go:317-356` imports,
   `treesitter.go:359-397` — approximate offsets — functions/classes) with **no
   `kotlin`/`csharp` case**, so each returns `nil`/empty with **no error**.
   `Analyzer.AnalyzeProject`'s per-file goroutine
   (`analyzer.go` — read at baseline, "Parse the file" block) only logs at
   `Debug` level when `err != nil`; here `err == nil`, so a Kotlin or C# file
   with real imports/functions/classes is indistinguishable in the emitted
   `AnalysisResult` from a genuinely-empty file. This is exactly the
   §11.4.69/§11.4.107 "PASS with no evidence the feature works" failure
   shape, applied to code analysis instead of AV.

8. **Zero test coverage of `treesitter.go` exists today.**
   The package tree at `255061b` is exactly `analyzer.go`, `analyzer_test.go`,
   `treesitter.go` — **no `treesitter_test.go` file exists.**
   `analyzer_test.go:14-58` (`TestDetectLanguage`) DOES assert
   `"Main.kt" → "kotlin"` (line 22) and `"Program.cs" → "csharp"` (line 37),
   but that only tests `detectLanguage`'s file-extension mapping in
   `analyzer.go` — it never touches `RegexParser`, `compilePatterns`, or any
   extraction function in `treesitter.go`. The plan's row-1 claim of
   "`detectLanguage`/regex extraction (codeanalysis)" done
   (`research/testing_infrastructure_plan.md:190`) conflates these two —
   `detectLanguage` is tested, the regex *extraction* is not (§7 reconciliation
   below).

9. **A same-class defect exists for Bash and Dart, discovered as a byproduct
   of this review, out of the assigned Kotlin/C# scope.** `analyzer.go`'s
   `detectLanguage` maps `.sh`/`.bash` → `"bash"` and `.dart` → `"dart"`, but
   `treesitter.go`'s native-attempt list (fact 2) and `compilePatterns`
   switch (fact 5) have no case for either — the identical silent-zero-
   extraction failure mode applies. This is flagged honestly in §8 as a
   tracked follow-up, not designed here (scope discipline — the assignment is
   Kotlin/C#).

---

## 2. Decision

### 2.1 Chosen binding: official `github.com/tree-sitter/go-tree-sitter` core + per-language first-party/community grammar modules, gated by Go's built-in `cgo` build tag

**Chosen:** `github.com/tree-sitter/go-tree-sitter` (the tree-sitter
org's own Go bindings) as the core parser/query engine, plus one grammar
module per R13-corpus language:

| Core / grammar need | Module |
|---|---|
| Core bindings | `github.com/tree-sitter/go-tree-sitter` |
| Go | `github.com/tree-sitter/tree-sitter-go/bindings/go` |
| Python | `github.com/tree-sitter/tree-sitter-python` (bindings/go) |
| Java | `github.com/tree-sitter/tree-sitter-java` (bindings/go) |
| C | `github.com/tree-sitter/tree-sitter-c` (bindings/go) |
| C++ | `github.com/tree-sitter/tree-sitter-cpp/bindings/go` |
| JavaScript | `github.com/tree-sitter/tree-sitter-javascript/bindings/go` |
| TypeScript | `github.com/tree-sitter/tree-sitter-typescript` (bindings/go) |
| C# | `github.com/tree-sitter/tree-sitter-c-sharp` (confirmed on pkg.go.dev: "comprehensive support for C# versions 1 through 13.0") |
| Kotlin | `github.com/fwcd/tree-sitter-kotlin/bindings/go` (community — no first-party `tree-sitter/tree-sitter-kotlin` with Go bindings exists; confirmed via search) |
| Bash (follow-up, §1 fact 9, not designed here) | `github.com/tree-sitter/tree-sitter-bash` (bindings/go) |

**Rejected alternative: `github.com/smacker/go-tree-sitter`.** This
community wrapper vendors many grammars (including a confirmed `kotlin`
subpackage, `github.com/smacker/go-tree-sitter/kotlin`) as subpackages of one
module — a genuine one-`go.mod`-entry convenience. It was rejected because:

- **C# support is unconfirmed.** Search of the smacker repository and its
  `pkg.go.dev` listing did not surface a `csharp` subpackage; the official
  binding's C# grammar (`tree-sitter/tree-sitter-c-sharp`) is confirmed and
  first-party. Kotlin/C# is the assigned scope, and only the official-org
  path confirms both halves.
- **Documented cross-platform build fragility.** The smacker repository's
  own issue tracker (read via search, not asserted from memory) shows open
  reports of exactly the class of silent build fragility this design must
  avoid: "Unable to build for windows" (issue #120), a CGO build failure
  ("undefined: Node", issue #167), and a missing-header build failure on
  Python ("../array.h file not found", issue #175). These are the same
  "looks fine, silently breaks on another platform/toolchain" failure shape
  §11.4.108/§11.4.201 forbid; the official binding's grammars are each
  independently released and tested by the tree-sitter org itself.
- **Monolithic grammar vendoring vs. R7's pluggable-dependency philosophy.**
  REQUIREMENTS.md R7 mandates dependencies be "pluggable, not hardcoded" and
  vendored per-submodule where owned; per-grammar Go modules (the official
  binding's pattern) let each language's grammar be pinned/upgraded
  independently, matching that discipline more directly than one
  community-vendored bundle whose internal grammar versions move together
  on the bundler's own release cadence.

Both binding families require **CGO** — this is not a smacker-specific cost.
Confirmed via search: "The smacker/go-tree-sitter library uses CGO to call
into the C runtime, since Tree-sitter is written in C… CGO must be enabled";
the official binding's own README states "this library requires CGO to
build" and additionally documents an escape hatch — loading a grammar at
runtime from a prebuilt shared library via `purego` (no CGO) — which is
noted here as a **future option for cross-compiled clients (R3)**, not
designed in this pass (§8 honest gap).

### 2.2 Build-tag gating (Go's built-in `cgo` tag, not a bespoke one)

Split `treesitter.go` into two files along the existing `initNativeParser` /
`parseNative` / `extract*Native` seam:

- `treesitter_native.go` — `//go:build cgo` — the REAL implementation:
  imports the core + per-language grammar modules from §2.1, implements
  `initNativeParser` to actually construct a `sitter.Parser`, set its
  language via the matching grammar module's `Language()` function, and
  return success; `parseNative` calls `parser.ParseCtx`/`Parse` and builds a
  real `Tree.Root` from the returned node tree; `extractImportsNative` /
  `extractFunctionsNative` / `extractClassesNative` walk that real tree via
  tree-sitter queries (`.scm` query files per language, the standard
  tree-sitter idiom) instead of regex.
- `treesitter_native_stub.go` — `//go:build !cgo` — **exactly today's
  code** (facts 1/3/4 verbatim), so a CGO-disabled build (hermetic CI,
  cross-compiled binaries without a C toolchain) degrades to the regex
  fallback exactly as it does today — no behavior change for that build
  mode, only a genuine native path added for CGO-enabled builds.

`cgo` is Go's own pre-declared build constraint (true only when
`CGO_ENABLED=1` and a working C compiler is found by the Go toolchain at
build time) — not a custom flag this project invents, so no new "did the
maintainer remember to pass a special flag" failure mode is introduced.

### 2.3 Explicit, labelled fallback — never silent (closes fact 7)

Two additive, non-breaking changes:

1. **`Tree` gains a `Fidelity` field**: `"native"` | `"regex-fallback"`.
   Set by `Parse()` at the point it decides which path served the request
   (`treesitter.go:143-154`'s existing branch). `AnalysisResult` gains
   `FidelityByLanguage map[string]string`, populated by the analyzer the
   first time each language is seen in a run, and surfaced in the JSON
   output — an operator/QA reviewer can see at a glance that, e.g., every
   `python` file in a given run was parsed at `"regex-fallback"` fidelity on
   a CGO-disabled build, per R11's "never silently degrade" mandate.
2. **A configured-but-genuinely-unsupported language now fails loud.**
   `RegexParser.Parse` returns a new sentinel error
   `ErrNoPatternsForLanguage` when `len(p.patterns) == 0` (today it always
   "succeeds" with an empty `FallbackParse{}` — fact 7). `Analyzer`'s
   per-file goroutine promotes this to a `WARN`-level log line
   (`zap.String("language", file.language)`) and appends the file path to a
   new `AnalysisResult.UnsupportedLanguageFiles []string` field — a Kotlin
   file today silently reports "0 imports, 0 functions, 0 classes"
   indistinguishably from a truly-empty file; after this change it is
   explicit and visible, never absorbed.

### 2.4 Interim Kotlin/C# regex patterns (ship immediately, before the CGO work lands)

Per §11.4.197 ("no gap left open in the backlog un-wired") this interim
patch stops the SILENT zero-extraction **now**, independent of the CGO
timeline, by adding two new `compilePatterns` cases modeled directly on the
existing Java/JavaScript shapes (`treesitter.go:277-285`):

- **`case "kotlin":`** — `import\s+([\w.]+)` (imports); `fun\s+(\w+)\s*\(`
  (top-level and member functions); `(?:class|object|interface)\s+(\w+)`
  (classes/objects/interfaces).
- **`case "csharp":`** — `using\s+([\w.]+)\s*;` (usings/imports);
  `(?:public|private|protected|internal|static|\s)+\s*(?:[\w<>\[\],\s]+)\s+(\w+)\s*\(([^)]*)\)\s*\{`
  (methods, same shape as the existing `java` `func` pattern at
  `treesitter.go:279`); `(?:public\s+)?(?:abstract\s+|sealed\s+)?class\s+(\w+)`
  and `interface\s+(\w+)` (classes/interfaces).

These patterns are always tagged `Fidelity: "regex-fallback"` (§2.3) even
after this interim patch lands — the native grammar work in §2.1/§2.2 is
what earns `"native"` for Kotlin/C#, and R11 forbids ever mislabeling
regex-derived results as tree-sitter-derived ones.

### 2.5 Containerized CGO build (§11.4.161 / §11.4.173)

CGO compilation needs a C toolchain (`gcc`/`clang`) at `go build` time to
compile the vendored tree-sitter core + grammar C sources — neither binding
family needs a *system-installed* `libtree-sitter` package (confirmed:
"uses CGO to call into the C runtime" describes compiling the vendored C
sources via cgo, not linking a system library). This toolchain dependency
belongs in the project's existing rootless-podman build container
(§11.4.161, composing §11.4.173's containerized-build mandate) rather than
assumed present on the bare host: the container's build stage needs
`gcc`/`clang` (+ `pkg-config`, the conventional cgo companion) added: a
one-line Containerfile/Dockerfile change, not a new build system. The
`!cgo`-tagged stub path (§2.2) remains the correct behavior for any build
environment where that toolchain is genuinely unavailable — never a build
failure, always a documented fidelity drop (§2.3).

### 2.6 R13-corpus language coverage matrix

R13 (`REQUIREMENTS.md`) lists 38 items; the subset that are directly
*parseable source languages* for `codeanalysis` (the rest — android,
android-aosp, rockchip, orange-pi, postgres, gin-gonic, cmake, make, gcc,
bazel, brotli, quic, http3, http, protocols, design-patterns, algorithms,
security, snyk, sonarqube, maven, gradle, linux, macos, debugging — are
platforms/frameworks/protocols/tools whose artifacts are themselves written
in one of the languages below, or are non-source-code skills entirely) is:

| Language | In R13 corpus | Default-enabled (`config.go:204`) | Today (`255061b`) | After interim patch (§2.4) | After native work (§2.1/§2.2) |
|---|---|---|---|---|---|
| Go | yes | yes | regex-fallback (native attempted, always fails — fact 1) | regex-fallback (unchanged) | **native** |
| Python | yes | yes | regex-fallback | regex-fallback | **native** |
| Java | yes | yes | regex-fallback | regex-fallback | **native** |
| C | yes | yes | regex-fallback | regex-fallback | **native** |
| C++ (cpp) | yes | yes | regex-fallback | regex-fallback | **native** |
| Kotlin | yes | **yes (default)** | **silent zero-extraction** (facts 2/5/7) | regex-fallback (interim patterns) | **native** (community grammar, §2.1) |
| C# (csharp) | not listed by name in R13; register/task title calls it out explicitly, operator-configurable via `normalizeLanguage` (fact 6) | no (operator must add it) | **silent zero-extraction** | regex-fallback (interim patterns) | **native** (official grammar, §2.1) |
| JavaScript | yes | no | regex-fallback (works — on-the-fly `newRegexParser` still hits the shared `case "javascript","typescript"`) | unchanged | native (optional, R13 does not force priority) |
| TypeScript | yes | no | regex-fallback | unchanged | native (optional) |
| Bash | yes | no | **silent zero-extraction** (fact 9 — NOT designed in this pass) | tracked follow-up | tracked follow-up |
| Rust | not in R13 | no | regex-fallback (already correct today) | unchanged | native (optional) |

---

## 3. Why (external precedent, §11.4.8) + Runtime Signature (§11.4.108)

**External precedent (verified this session, not from training-data memory
alone):**

- `github.com/tree-sitter/go-tree-sitter` — official Go bindings, requires
  CGO to build, grammars fetched individually
  (`go get github.com/tree-sitter/tree-sitter-<lang>@latest`), and supports
  an alternative `purego` shared-library loading path with no CGO.
- `github.com/tree-sitter/tree-sitter-c-sharp` — official, first-party C#
  grammar with Go bindings, "comprehensive support for C# versions 1 through
  13.0" (pkg.go.dev listing).
- `github.com/fwcd/tree-sitter-kotlin/bindings/go` — the maintained
  Go-bindable Kotlin grammar; not under the `tree-sitter` GitHub org.
- `github.com/smacker/go-tree-sitter` — community wrapper, vendors a
  `kotlin` subpackage, requires `CGO_ENABLED=1`, with documented
  cross-platform build-failure issues (#120 Windows, #167 CGO build,
  #175 missing header) cited as the reason it was rejected over the
  official binding for this design.
- `github.com/malivvan/tree-sitter` — a CGO-free alternative wrapping a WASM
  build of tree-sitter via `wazero`; noted as a possible future path for
  cross-compiled R3 clients, not adopted here (§8 honest gap — not
  evaluated in depth, out of this pass's scope).

**Runtime signature (definition of done, §11.4.108) — proves native parsing
actually works on a real file on a clean deployment, never a stub error:**

1. Build the `codeanalysis` package with `CGO_ENABLED=1` inside the
   §2.5 container (a "clean deployment" per §11.4.108 — fresh container,
   no leftover state).
2. Parse the project's own committed `internal/codeanalysis/analyzer.go`
   (a real, checked-in Go file, hand-countable ground truth) through the
   `treesitter_native.go` path.
3. Assert: `Tree.Fidelity == "native"`; the extracted import count,
   exported-function-name set, and struct/interface-name set exactly equal
   a hand-verified ground-truth list for that exact file (e.g. it must
   include `Analyzer`, `AnalysisResult`, `Import`, `Pattern`,
   `NewAnalyzer`, `AnalyzeProject`, `detectLanguage`, …) — never merely
   "`err == nil`".
4. For Kotlin: parse a real, checked-in Kotlin source file from the P5.T2
   real-repo corpus (`research/testing_infrastructure_plan.md`'s own G12
   integration-test description: "parse a real Android/Kotlin repo → real
   symbols") through the interim regex path (§2.4) TODAY, and again through
   `treesitter_native.go` once the `fwcd/tree-sitter-kotlin` grammar lands —
   same file, same ground-truth assertions, only `Fidelity` should differ
   (`"regex-fallback"` → `"native"`); any content divergence between the two
   fidelities is itself a finding, not an accepted difference.
5. Same for C#, same pattern.

A PASS that only shows `err == nil` and a non-empty slice, without a
hand-verified ground-truth comparison, is rejected as a §11.4/§11.4.1
bluff — it is exactly the class of defect fact 7 describes.

---

## 4. Test-case count (RED-first, §11.4.115) + reconciliation with `research/testing_infrastructure_plan.md`

**Enumerated cases** (all §1.1-paired where marked):

| # | Type | Case | RED today? |
|---|---|---|---|
| 1 | unit | `compilePatterns` for `kotlin` produces a non-empty pattern map; extraction on a real fixture Kotlin file matches hand-verified ground truth | **RED now** (fact 5 — no `kotlin` case exists) |
| 2 | unit | Same, for `csharp` | **RED now** (fact 5) |
| 3 | unit | `Parse`/`ExtractImports`/`ExtractFunctions`/`ExtractClasses` return `ErrNoPatternsForLanguage` (new sentinel) for a language with zero compiled patterns, instead of a silent empty success | **RED now** (fact 7 — `parseFallback` never errors) |
| 4 | unit | `Tree.Fidelity` is populated (`"native"`/`"regex-fallback"`) for every language on every `Parse` call, never left blank | **RED now** (field doesn't exist yet) |
| 5 | unit | `normalizeLanguage("kt")=="kotlin"`, `normalizeLanguage("c#")==normalizeLanguage("csharp")=="csharp"` (already-passing behavior, made explicit so a future edit to the switch is caught) | GREEN today (behavior exists, only the assertion is new) |
| 6 | integration | Native path (CGO build) parses a real committed Go file; extracted symbols match hand-verified ground truth (§3 runtime signature) | **RED now** (native path doesn't exist — fact 1/3) |
| 7 | integration | Real Kotlin file from the P5.T2 corpus parsed via interim regex (today) and later via native grammar; ground truth matches in both, `Fidelity` differs correctly | **RED now** (fact 2/5/7) |
| 8 | integration | Same, for a real C# file | **RED now** |
| 9 | fuzz | Malformed/truncated source (all languages incl. kotlin/csharp) does not panic the regex engine or the native parser across N iterations (seeded from real files with random truncation/byte-flips) | new coverage |
| 10 | mutation (§1.1) | Remove the `kotlin` case from `compilePatterns` → test #1 must go RED | proves #1 load-bearing |
| 11 | mutation (§1.1) | Remove the `csharp` case → test #2 must go RED | proves #2 load-bearing |
| 12 | mutation (§1.1) | Revert `ErrNoPatternsForLanguage` to the old always-nil-error behavior → test #3 must go RED | proves #3 load-bearing |
| 13 | mutation (§1.1) | Stub out one native grammar's `Language()` binding → the corresponding integration test (#6/#7/#8) must go RED | proves the native-path integration tests exercise the real grammar, not a fixture-independent bluff |

Tests #1, #2, #3, #4, #6, #7, #8 are **RED against the committed `255061b`
baseline right now** — this is the direct, provable demonstration of the
gap (§11.4.115 RED-baseline-on-the-broken-artifact), not a hypothetical.

**Reconciliation with `research/testing_infrastructure_plan.md`:**

- Line 190 (test-type row 1, unit): claims "…`detectLanguage`/regex
  extraction (codeanalysis)…" as done. Per fact 8, this conflates two
  different things: `detectLanguage` (in `analyzer.go`) IS tested
  (`TestDetectLanguage`, `analyzer_test.go:14`); the regex *extraction*
  (`RegexParser`/`compilePatterns` in `treesitter.go`) has **zero** tests
  today. **Extension needed (§11.4.186 cross-doc consistency):** split this
  row's codeanalysis claim into two distinct sub-items — `detectLanguage`
  (DONE) and per-language `RegexParser` extraction incl. kotlin/csharp (NOT
  YET — closed by this design's tests #1/#2/#5).
- Line 298 (per-gap matrix, G12 row): `"unit(per-language extraction incl.
  kotlin), integration(parse real Android/Kotlin repo → real symbols),
  fuzz(malformed source no crash), mutation(remove a grammar → extraction
  FAILs) | Challenge: yes | extracted-symbols JSON from a real repo (no
  lorem); per-language coverage incl. kotlin"` — directly satisfied by
  tests #1/#2/#4/#5 (unit), #6/#7/#8 (integration), #9 (fuzz), #10-#13
  (mutation). **Extension needed:** the row's unit-test wording names only
  "kotlin" even though the register entry (`GAPS_AND_RISKS_REGISTER.md`
  G12 heading) and this design cover "Kotlin/C#" — recommend the plan row
  be reworded "…incl. kotlin, csharp…" for consistency with the register
  (§11.4.186 anti-divergence; a wording gap, not a blocking defect).
- Section 3.2's `HQA-WIZARD-ANDROID` example (`input: { techs: [android,
  android_aosp, java, kotlin, cpp, cmake] }`) already exercises Kotlin at
  the wizard/end-to-end layer — this design's test #7 is the
  package-local counterpart that isolates the codeanalysis defect
  specifically, rather than only proving it transitively through the whole
  wizard pipeline.

---

## 5. Honest gaps (§11.4.6 — never guessed, marked explicitly)

- **UNCONFIRMED:** whether the project's existing rootless-podman build
  container (§11.4.161) already has `gcc`/`clang` + `pkg-config` installed.
  Not verified in this offline design pass (no container introspection was
  available/in-scope here) — must be checked before implementation, and the
  Containerfile change (§2.5) added if absent.
- **UNCONFIRMED:** exact ABI/version compatibility between a single pinned
  `github.com/tree-sitter/go-tree-sitter` core release and each
  independently-released grammar module (`tree-sitter-go`, `-python`,
  `-java`, `-c`, `-cpp`, `-javascript`, `-typescript`, `-c-sharp`,
  `fwcd/tree-sitter-kotlin`) at actual implementation time. Each grammar
  module is versioned separately; this design's read-only pass performed
  web searches confirming each module's *existence* and general shape, but
  did not (and could not, without `go get`/network module-resolution
  access in this pass) build a concrete go.mod and verify a green
  `go build`/`go vet` with a specific version set. This is a build-time
  verification step for the implementation, not assumed solved here.
- **UNCONFIRMED:** `github.com/fwcd/tree-sitter-kotlin`'s grammar
  completeness against current Kotlin language features (e.g. context
  receivers, newer syntax additions) and its release/maintenance cadence —
  it is community-maintained (not under the `tree-sitter` GitHub org), and
  this pass confirmed only that a Go-bindable package exists on pkg.go.dev,
  not its feature completeness.
- **UNCONFIRMED:** whether `github.com/smacker/go-tree-sitter` has a `csharp`
  subpackage. Not found in this pass's searches; recorded as unconfirmed
  absence, not asserted non-existence, since it was not a deciding factor
  (the official-binding path was chosen regardless — §2.1).
- **UNCONFIRMED / not engineered in this pass:** the precise Go mechanics
  for test #13's "stub out one native grammar's binding without breaking
  every other grammar's build" — a build-tag-per-grammar or interface-
  injection seam is implied but needs a short implementation-time spike;
  not asserted as solved here.
- **Out of scope, flagged not designed:** the Bash/Dart same-class silent-
  zero-extraction defect (§1 fact 9) is a real, provable finding discovered
  during this review but is **not** designed in this document — the
  assignment is Kotlin/C#. It MUST be tracked (a G12 follow-up note or a new
  register entry) per §11.4.197/§11.4.186 rather than left un-filed.
- **Not evaluated in depth:** the `purego`/WASM (`malivvan/tree-sitter`)
  no-CGO alternative, relevant to R3's multi-OS/multi-arch client
  cross-compilation story. Noted as a candidate future direction (§2.1,
  §3), not designed or compared in depth here — doing so is future work,
  not asserted as already decided.

---

## Sources verified (2026-07-15)

- https://github.com/tree-sitter/go-tree-sitter
- https://github.com/tree-sitter/tree-sitter-c-sharp (pkg.go.dev listing)
- https://github.com/fwcd/tree-sitter-kotlin (bindings/go, via pkg.go.dev)
- https://github.com/smacker/go-tree-sitter (+ issues #120, #167, #175)
- https://github.com/malivvan/tree-sitter
- `git show 255061b:docs/research/mvp/Agent_AI_Skill_Tree_Development/project/internal/codeanalysis/treesitter.go`
- `git show 255061b:docs/research/mvp/Agent_AI_Skill_Tree_Development/project/internal/codeanalysis/analyzer.go`
- `git show 255061b:docs/research/mvp/Agent_AI_Skill_Tree_Development/project/internal/codeanalysis/analyzer_test.go`
- `git show 255061b:docs/research/mvp/Agent_AI_Skill_Tree_Development/project/internal/config/config.go`
- `docs/research/mvp/Agent_AI_Skill_Tree_Development/GAPS_AND_RISKS_REGISTER.md` (G12 entry)
- `docs/research/mvp/Agent_AI_Skill_Tree_Development/research/testing_infrastructure_plan.md` (lines 190, 298)
- `docs/research/mvp/Agent_AI_Skill_Tree_Development/REQUIREMENTS.md` (R2, R7, R11, R13)
- `docs/research/mvp/Agent_AI_Skill_Tree_Development/SPEC.md:29` ("tree-sitter (official Go bindings)")
