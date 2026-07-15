# TOON Go Codec — Decision Doc

**Revision:** 1
**Last modified:** 2026-07-15T00:00:00Z
**Scope:** HelixKnowledge Skill Graph System — API wire-format serialization
**Question:** How does the Go service serialize/deserialize TOON (primary wire format) with JSON as fallback?
**Verdict:** **VENDOR the official Go implementation** — `github.com/toon-format/toon-go@v0.0.0-20251202084852-7ca0e27c4e8c` (MIT). Do **not** hand-roll a codec.

> Supersedes the earlier "interpret TOON as TOML" note. TOON is a distinct, real
> serialization format (Token-Oriented Object Notation), unrelated to TOML.

---

## 1. Reachability + provenance evidence (real lookups)

TOON is real and actively maintained. All lines below are captured command output, not recollection.

**Canonical repo reachable (`git ls-remote`):**
```
$ git ls-remote https://github.com/toon-format/toon.git
a19a1179193451fad40f11ef88de5f363ea3684a	HEAD
a19a1179193451fad40f11ef88de5f363ea3684a	refs/heads/main
...
```
Default branch `main`, HEAD `a19a117`. Repo description (GitHub API):
`"🎒 Token-Oriented Object Notation (TOON) – Compact, human-readable, schema-aware JSON for LLM prompts. Spec, benchmarks, TypeScript SDK."`

**The `toon-format` org ships official SDKs across languages** (GitHub API `orgs/toon-format/repos`), including a **Go** one:
```
toon-format/toon         TypeScript   (reference SDK + benchmarks)
toon-format/spec         JavaScript   (canonical SPEC.md)
toon-format/toon-go      Go           ← the Go implementation
toon-format/toon-python  Python
toon-format/toon-rust    Rust
toon-format/toon-java    Java
toon-format/toon-swift   Swift
toon-format/toon-dart    Dart
toon-format/toon-dotnet  C#
toon-format/ToonFormat.jl Julia
```

**Spec versioning:** the main repo carries semver tags (`git ls-remote --tags`) up to at least `v2.3.0`; the normative grammar now lives in `toon-format/spec` (`SPEC.md`). The format is stable and past v2.

---

## 2. TOON grammar summary (from the real SPEC + README examples)

TOON is a whitespace-significant, schema-aware, token-minimizing encoding of the JSON
data model. Its whole reason to exist is fewer LLM tokens than JSON for the same data
(no repeated keys in arrays-of-objects, no braces/brackets clutter).

**Defaults:** indent = **2 spaces** per depth level (configurable); delimiter = **comma** (tab / pipe selectable).

### 2.1 Objects
`key: value` per line. A bare `key:` (no value) opens a nested object; its fields are
indented one level deeper. Indentation is significant (it defines structure).
```
user:
  id: 123
  name: Ada
  contact:
    email: ada@example.com
```
≡ JSON `{"user":{"id":123,"name":"Ada","contact":{"email":"ada@example.com"}}}`

### 2.2 Scalars
- Numbers unquoted (`42`, `3.14`). Canonical form: no exponent for values in `[1e-6, 1e21)`;
  **trailing fractional zeros are not emitted** by a conforming encoder (`14.50` → `14.5`).
- Booleans `true` / `false` and `null` unquoted.
- Strings unquoted when unambiguous; **double-quoted when** they contain a delimiter/colon/
  quote/backslash, equal `true`/`false`/`null`, look numeric, or have leading/trailing space.
- Escapes inside quotes: `\\`, `\"`, `\n`, `\t`, `\uXXXX`.
```
url: "http://example.com:8080"
version: "3.0"
empty: ""
```

### 2.3 Inline primitive arrays — `[N]` length marker
```
tags[3]: admin,ops,dev
```
≡ `{"tags":["admin","ops","dev"]}`. The `[3]` is the element count (self-describing / validatable).

### 2.4 Tabular array-of-objects (the token-efficiency win)
Header declares length **and** field names once; each row supplies values in field order,
indented one level. Applies only when **all elements are objects with the identical key set
and all values are primitive** (SPEC §9.3).
```
users[2]{id,name}:
  1,Ada
  2,Linus
```
≡ `{"users":[{"id":1,"name":"Ada"},{"id":2,"name":"Linus"}]}`

### 2.5 Non-uniform / nested arrays — expanded "list" form (`- ` marker)
When the tabular precondition is violated (mixed shapes, varying keys, nested structures),
each element is emitted as a `- ` list item at depth+1:
```
items[3]:
  - 1
  - a: 1
  - text
```
Objects as list items keep their fields indented under the dash:
```
items[2]:
  - id: 1
    name: First
  - id: 2
    name: Second
    extra: true
```
Nested arrays nest a header per item:
```
pairs[2]:
  - [2]: 1,2
  - [2]: 3,4
```

### 2.6 Empty + root forms
- Empty array: `key: []` (legacy `key[0]:`). Empty nested object: bare `key:`. Empty doc → `{}`.
- Root-level primitive array: `[3]: a,b,c`. Root tabular: `[2]{id,name}:` + rows. Empty root array: `[]`.

---

## 3. Go implementation availability (verified findings)

**A Go implementation EXISTS and is the org-official one.** Verified source lines:

**Repo metadata (GitHub API `repos/toon-format/toon-go`):**
```
"full_name": "toon-format/toon-go",
"description": "🐹 Community-driven Go implementation of TOON",
"license": { "spdx_id": "MIT" },
"default_branch": "main",
"pushed_at": "2025-12-02T08:48:57Z",
"stargazers_count": 142,
"archived": false
```

**Module path (`go.mod`, fetched raw):**
```
module github.com/toon-format/toon-go
go 1.23
```

**License (`LICENSE`, fetched raw):**
```
MIT License
Copyright (c) 2025-PRESENT Bintang Pradana Erlangga Putra
Copyright (c) 2025-PRESENT Johann Schopplich
```

**Indexed on pkg.go.dev** (`https://pkg.go.dev/github.com/toon-format/toon-go`) — verified:
- Latest version string (exact): **`v0.0.0-20251202084852-7ca0e27c4e8c`** (pseudo-version; **no semver tag published yet** — `git ls-remote --tags` returns empty).
- Last published: **Dec 2, 2025**. License detected: **MIT**.
- Current `main` HEAD: `7ca0e27c4e8c695d99c88ca4a123f409086da91e` (matches the pseudo-version hash).

**Exported API surface (from pkg.go.dev index) — a complete, spec-aligned codec:**
```
Marshal(v any, opts ...EncoderOption) ([]byte, error)
MarshalString(v any, opts ...EncoderOption) (string, error)
Unmarshal(data []byte, v any, opts ...DecoderOption) error
UnmarshalString(s string, v any, opts ...DecoderOption) error
Decode(data []byte, opts ...DecoderOption) (any, error)
DecodeString(s string, opts ...DecoderOption) (any, error)
NewEncoder(opts ...EncoderOption) *Encoder
NewDecoder(opts ...DecoderOption) *Decoder
NewObject(fields ...Field) Object
Types: Encoder, Decoder, EncoderOption, DecoderOption, Delimiter, Field, Object
Options: WithLengthMarkers(bool), WithArrayDelimiter(Delimiter), WithDelimiter/WithDocumentDelimiter,
         WithIndent(int), WithDecoderIndent(int), WithStrictMode(bool), WithTimeFormatter(...)
```
Usage confirmed on the repo page:
```go
encoded, err := toon.Marshal(in, toon.WithLengthMarkers(true))
var out Payload
err = toon.Unmarshal(encoded, &out)
```

The option set maps 1:1 onto the spec knobs we need: `[N]` markers (`WithLengthMarkers`),
delimiter choice (`WithArrayDelimiter`/`Delimiter`), indent width (`WithIndent`), strict
round-trip validation (`WithStrictMode`). It follows `encoding/json` conventions (Marshal/
Unmarshal + struct tags), so it drops into a Go service idiomatically.

---

## 4. DECISION — vendor `toon-go`, wrap it behind `internal/serialization`

### 4.1 The choice
```
go get github.com/toon-format/toon-go@v0.0.0-20251202084852-7ca0e27c4e8c
```
Pin the exact pseudo-version in `go.mod` (reproducible builds; §11.4.108). MIT license is
compatible with our distribution. It is the `toon-format`-org Go SDK, indexed on pkg.go.dev,
Go 1.23, with the full Marshal/Unmarshal/streaming API and every spec option we need.

### 4.2 Why trustworthy (and the honest caveats)
**Trust factors:** lives under the official `toon-format` org alongside the spec + reference
TS SDK; MIT; 142 stars; recently pushed (2025-12-02); parses the same conformance model as the
canonical `SPEC.md`; complete encoder+decoder+streaming surface.

**Caveats (recorded, not hidden — §11.4.6):**
1. **No semver tag yet** → pinned by pseudo-version. Watch for a `v1.x` tag and upgrade
   deliberately behind our own golden tests.
2. **Small maintainer base**, self-described "community-driven." Mitigated by our wrapper
   (§4.3) and by our own conformance/golden test suite (§5) so a swap is cheap.
3. **Pre-1.0** → API may shift. The wrapper isolates every direct call site.

### 4.3 Wiring — `internal/serialization` + Gin content negotiation
Thin adapter package (no re-implementation): `internal/serialization` exposes
`Marshal(v)`, `Unmarshal(data, v)`, and a `Codec` interface with a `toonCodec` (delegates to
`toon.Marshal/Unmarshal`, configured `WithLengthMarkers(true)` + comma delimiter + 2-space
indent) and a `jsonCodec` (stdlib `encoding/json`, the fallback).

Content negotiation in the Gin/HTTP layer:
- **Response:** inspect `Accept`. `application/toon` → TOON codec, `Content-Type: application/toon`.
  `application/json` (or `*/*`, or unknown) → JSON fallback.
- **Request:** inspect `Content-Type`. `application/toon` → TOON decode; else JSON decode.
- Register a Gin render/binding shim (or a small middleware) that routes through the `Codec`
  interface so handlers stay format-agnostic. TOON is primary; JSON is always available as
  the safety net (§11.4.6 honest boundary — a client that can't speak TOON still works).

### 4.4 Rejected alternatives
- **Hand-roll a Go codec in `internal/serialization`** — rejected. Reproduces a maintained,
  spec-conformant MIT library; the tabular/list/quoting/number-canonicalization rules are
  fiddly (a from-scratch encoder is a bug farm and a §11.4.124/dead-effort risk). Reserve
  this only if `toon-go` is later abandoned AND no other org SDK is bindable.
- **Treat TOON as TOML** (the superseded note) — rejected/incorrect. TOON ≠ TOML; different
  grammar, different goal (LLM token efficiency vs config files). No shared parser.
- **JSON-only, drop TOON** — rejected. TOON is the mandated primary wire format (token savings
  on the LLM/API path); JSON stays only as fallback.

---

## 5. Test-vector table (from REAL spec/README examples)

Round-trip (`JSON → Marshal → TOON` and `TOON → Unmarshal → JSON`) golden tests use these.
Every pair is drawn from an example actually observed in `toon-format/spec` `SPEC.md` or the
`toon-format/toon` README — none invented. Encoder options: 2-space indent, comma delimiter,
length markers on.

| # | Shape | JSON input | Expected TOON output | Source |
|---|-------|-----------|----------------------|--------|
| 1 | Simple object | `{"id":1,"name":"Alice","active":true}` | `id: 1`<br>`name: Alice`<br>`active: true` | spec SPEC.md §simple object |
| 2 | Inline primitive array + `[N]` | `{"tags":["admin","ops","dev"]}` | `tags[3]: admin,ops,dev` | spec SPEC.md §inline array |
| 3 | Tabular array-of-objects | `{"users":[{"id":1,"name":"Ada"},{"id":2,"name":"Linus"}]}` | `users[2]{id,name}:`<br>`  1,Ada`<br>`  2,Linus` | toon README |
| 4 | Tabular w/ numeric fields | `{"items":[{"id":1,"name":"Alice","price":9.99},{"id":2,"name":"Bob","price":14.5}]}` | `items[2]{id,name,price}:`<br>`  1,Alice,9.99`<br>`  2,Bob,14.5` | spec SPEC.md §tabular (see note ⚠) |
| 5 | Nested objects / indentation | `{"user":{"id":123,"name":"Ada","contact":{"email":"ada@example.com"}}}` | `user:`<br>`  id: 123`<br>`  name: Ada`<br>`  contact:`<br>`    email: ada@example.com` | spec SPEC.md §nested |
| 6 | String quoting (colon/numeric-like) | `{"url":"http://example.com:8080","version":"3.0","empty":""}` | `url: "http://example.com:8080"`<br>`version: "3.0"`<br>`empty: ""` | spec SPEC.md §quoting |
| 7 | Empty array | `{"key":[]}` | `key: []` | spec SPEC.md §empty |
| 8 | Root-level primitive array | `["a","b","c"]` | `[3]: a,b,c` | spec SPEC.md §root array |
| 9 | Mixed / non-uniform array (dash list) | `{"items":[1,{"a":1},"text"]}` | `items[3]:`<br>`  - 1`<br>`  - a: 1`<br>`  - text` | spec SPEC.md Appendix A |
| 10 | Non-tabular objects (varying keys) | `{"items":[{"id":1,"name":"First"},{"id":2,"name":"Second","extra":true}]}` | `items[2]:`<br>`  - id: 1`<br>`    name: First`<br>`  - id: 2`<br>`    name: Second`<br>`    extra: true` | spec SPEC.md Appendix A |
| 11 | Array of arrays (nested headers) | `{"pairs":[[1,2],[3,4]]}` | `pairs[2]:`<br>`  - [2]: 1,2`<br>`  - [2]: 3,4` | spec SPEC.md Appendix A |

**⚠ Note on vector 4 (honesty per §11.4.6):** the SPEC example block literally shows `14.50`,
but TOON's canonical number rule states *trailing fractional zeros are not emitted by a
conforming encoder*. So `14.50` (JSON) canonicalizes to `14.5` (TOON). The golden test MUST
assert `14.5` — assert against the encoder's real output, not the illustrative README digits.
`9.99` is already canonical and stays `9.99`.

**Test discipline:** (a) `JSON → toon.Marshal → assert == expected TOON`; (b) `expected TOON →
toon.Unmarshal → assert deep-equal to JSON input`; (c) full round-trip
`JSON → Marshal → Unmarshal → deep-equal JSON` for value-identity; (d) verify the `[N]` count
in the header equals the actual row/element count (self-describing-length invariant). These
vectors are the golden fixtures the wrapper (§4.3) and its meta-test mutation pair validate.

---

## 6. Sources verified (2026-07-15)

- `git ls-remote https://github.com/toon-format/toon.git` — repo reachable, `main` HEAD `a19a117`.
- GitHub API `orgs/toon-format/repos` — confirmed `toon-format/toon-go` (Go) + spec + multi-language SDKs.
- GitHub API `repos/toon-format/toon-go` — MIT, pushed 2025-12-02, 142★, not archived.
- Raw `toon-format/toon-go/main/go.mod` — `module github.com/toon-format/toon-go`, `go 1.23`.
- Raw `toon-format/toon-go/main/LICENSE` — MIT.
- `https://pkg.go.dev/github.com/toon-format/toon-go` — indexed, `v0.0.0-20251202084852-7ca0e27c4e8c`, API surface.
- Raw `toon-format/spec/main/SPEC.md` — grammar + example blocks (objects, inline `[N]` arrays, tabular `[N]{fields}:`, nested, quoting, numbers/booleans/null, empty, root arrays, Appendix A dash-list forms).
- `toon-format/toon` README — `users[2]{id,name}:` tabular example.

**Negative finding:** no published semver tag on `toon-go` yet (`git ls-remote --tags` empty) — consume via pinned pseudo-version.
