# OpenDesign Incorporation Research (R12)

**Revision:** 1
**Last modified:** 2026-07-15T00:00:00Z
**Scope:** How the HelixKnowledge Skill Graph System incorporates OpenDesign
(`git@github.com:nexu-io/open-design.git`) for ALL client / design / styling
work per operator mandate R12, across the R3 client surfaces (CLI, TUI,
REST-consumers, Web/React, Desktop, Mobile) + the skill-creation wizard flow +
architecture diagrams.
**Author:** research subagent (read-only; zero-bluff — every external claim
carries a verified source line).
**Verification substrate:** repo cloned read-only into the session scratchpad at
commit `94a5bd2e08c5a13c774785907532babcbf9777f7` (default branch `main`) and
inspected directly. Claims below cite the exact file/line in that tree or the
exact command output.

---

## 0. TL;DR

- **Repo reachable: YES.** `git ls-remote` succeeds over BOTH SSH and HTTPS;
  HEAD = `94a5bd2e08c5a13c774785907532babcbf9777f7`. Contents inspected from a
  real clone — nothing below is invented.
- **What it is:** OpenDesign (npm name `open-design`, v0.15.1, **Apache-2.0**) is
  a **local-first, agent-native design product** — self-described "the
  open-source Claude Design alternative" and "the Figma alternative for the agent
  era." It is a pnpm + TypeScript monorepo (Node ~24) shipping as **(a) a desktop
  app, (b) an `od` CLI, (c) an MCP server, (d) 150 brand-grade `DESIGN.md` design
  systems each with machine-readable design tokens, (e) 100+ design "skills",
  (f) plugins, (g) a Figma-import plugin**. It does NOT ship as a simple
  `npm install <tokens>` component library — the consumable design substrate is
  the `DESIGN.md` systems + their `design-tokens.json` / `tokens.css`, driven by a
  coding agent.
- **Incorporation approach:** add as a **third-party (vendor) git submodule at
  `submodules/open_design/`** (§11.4.28 layout, §11.4.29 snake-case dir; the
  upstream repo name `open-design` is preserved per the §11.4.29 vendor-submodule
  exception). `nexu-io` is NOT one of our owned orgs, so the §11.4.28(C)
  nested-own-org prohibition does not apply — it is treated as an external
  dependency, not an equal-codebase owned submodule.
- **Honest gaps (detailed in §7):** **PSD export is NOT supported by OpenDesign**
  (0 occurrences in tree); native `.fig` is not produced (Figma round-trips via a
  JSON capture + the OD Figma plugin, and SVG rasterizes on import); mobile is
  "framed HTML prototypes" (iPhone/Pixel), NOT native iOS / Android / HarmonyOS /
  Aurora OS design kits; the bundled desktop app targets macOS + Windows only
  (Linux via source/Docker/MCP); `@open-design/components` is `private:true`
  (OD's own React app UI, not a publishable consumable lib).

---

## 1. Reachability evidence (Task 1)

Both transports succeed and return the identical HEAD SHA.

```
$ git ls-remote git@github.com:nexu-io/open-design.git HEAD
94a5bd2e08c5a13c774785907532babcbf9777f7	HEAD
EXIT_SSH=0

$ git ls-remote https://github.com/nexu-io/open-design.git HEAD
94a5bd2e08c5a13c774785907532babcbf9777f7	HEAD
EXIT_HTTPS=0
```

The repo is **public and reachable** (HTTPS anonymous ls-remote works → not a
private repo requiring credentials). A read-only shallow clone
(`--depth 1 --filter=blob:limit=200k`) succeeded: **11,373 tracked files**,
default branch `main`, latest commit:

```
94a5bd2e08c5a13c774785907532babcbf9777f7  2026-07-15 14:44:22 +0000  fix BYOK OpenCode permission bypass (#5701)
origin/HEAD -> origin/main
```

No access blocker. (If a future run finds it private/unreachable, the correct
action is to request repo access from `nexu-io`, not to invent contents.)

---

## 2. What OpenDesign actually provides (Task 2 — verified)

### 2.1 Identity, license, stack

| Fact | Value | Source (in cloned tree @ `94a5bd2e`) |
|---|---|---|
| npm name / version | `open-design` / `0.15.1`, `private: true` | `package.json` |
| License | **Apache-2.0** | `package.json` `"license"`; `LICENSE` header "Apache License Version 2.0" |
| Description | "Local-first design product: detects your installed code-agent CLI, runs design skills + design systems, streams artifacts into a sandboxed preview." | `package.json` `"description"` |
| Package manager / runtime | pnpm 10.33.2 workspace, Node ~24, TypeScript 5.9.3, ESM | `package.json` `engines`/`packageManager`; `pnpm-workspace.yaml` |
| CLI binary | `od` → `./apps/daemon/bin/od.mjs` | `package.json` `"bin"` |
| Self-positioning | "The open-source Claude Design alternative"; "the Figma alternative for the agent era" | `README.md` H1 + "What is Open Design" |

**It is NOT** a design-token npm package, a component library you `npm install`,
or a Figma kit. **It IS** a whole agent-native design *product / workspace*:
a coding agent (Claude Code, Codex, Cursor, Gemini, OpenCode, …) reads OpenDesign
skills + a chosen `DESIGN.md` design system off the filesystem and streams design
artifacts (HTML / decks / images / video) into a sandboxed preview, exportable to
HTML / PDF / PPTX / MP4. (`README.md` "What is Open Design"; `CONTEXT.md`.)

### 2.2 Monorepo layout (the consumable pieces)

Verified `ls` of the tree:

- `design-systems/` — **153 directories** (README states "150 brand-grade
  `DESIGN.md` systems"; catalog incl. `apple`, `claude`, `stripe`, `vercel`,
  `notion`, `figma`, `tesla`, `bmw`, `ibm`, `nvidia`, `default`, …). **This is the
  primary consumable design substrate.** Each system directory (e.g.
  `design-systems/apple/`) ships:
  - `DESIGN.md` — the human "brand contract" (9-section schema: visual theme,
    color palette & roles, typography, component stylings, layout, depth &
    elevation, do's/don'ts, …) — verified against `design-systems/claude/DESIGN.md`
    section headers. Multilingual variants (`DESIGN-ar.md`, `DESIGN-ja.md`, …).
  - **`design-tokens.json`** — machine-readable tokens, `"format":
    "od-design-tokens/v1"`, `"contract": "TOKEN_SCHEMA"`, typed token entries
    (`{name:"--bg", value:"#ffffff", type:"color", layer:"A1-identity", …}`).
    Verified in `design-systems/apple/design-tokens.json`.
  - **`tokens.css`** — the `:root { … }` CSS-variable block agents paste verbatim
    into an artifact's first `<style>`; standard token names (`--bg`, `--surface`,
    `--fg`, `--muted`, `--border`, `--accent`, `--radius-*`, `--ease-*`, …).
    Verified in `design-systems/apple/tokens.css` (header comment documents intent
    + lint enforcement `apps/daemon/src/lint-artifact.ts`).
  - **`tailwind-v4.css`** — Tailwind v4 binding of the same tokens.
  - `components.html` + `components.manifest.json` — a rendered component
    inventory + compact manifest for the brand.
  - `USAGE.md`, `manifest.json`, `preview/`, `source/`, `system/`.
- `skills/` — **164 directories** ("100+ skills ship in the box"), each a
  `SKILL.md` folder (Claude Code skill convention) extended with `od:` frontmatter
  (`mode`, `platform`, `scenario`, `design_system.requires`, `fidelity`,
  `example_prompt`). Modes: `prototype` (web/mobile/desktop single-page artifacts),
  `deck`, `image`, `video`, `audio`, `template`, `design-system`, `utility`.
  Relevant skills present: `apple-hig`, `flutter-animating-apps`,
  `hand-drawn-diagrams`, `frame-flowchart-sticky`, `login-flow`, `pptx`, `pdf`,
  `pptx-generator`, `slides`, `figma-generate-design`, `figma-generate-library`,
  `figma-create-design-system-rules`, `figma-code-connect-components`,
  `brand-extract`, `brand-guidelines`, `color-expert`, `canvas-design`. (Verified
  via `ls skills/` + `README.md` "Skills".)
- `design-templates/` — **115 directories** (web-prototype, saas-landing,
  dashboard, mobile-app [iPhone 15 Pro / Pixel framed], mobile-onboarding,
  social-carousel, email-marketing, deck templates, `hyperframes`, `critique`, …).
- `craft/` — cross-cutting design-quality rule files: `accessibility-baseline.md`
  (WCAG 2.2 AA target, jurisdiction floor, native-mobile parity notes),
  `color.md`, `typography.md`, `typography-hierarchy.md`, `animation-discipline.md`,
  `anti-ai-slop.md`, `state-coverage.md`, `form-validation.md`, `laws-of-ux.md`,
  `rtl-and-bidi.md`. These are the "how to use the tokens well" guardrails.
- `packages/` — engine packages incl. `components` (`@open-design/components`
  v0.8.0 — "Shared Open Design React UI primitives", React 18.3.1 peer, **but
  `private: true`** = workspace-internal, not npm-published), `contracts`, `host`,
  `plugin-runtime`, `registry-protocol`, `metatool`, `sidecar`, `platform`,
  `release`, `download`.
- `apps/` — `daemon` (the `od` CLI + local HTTP/MCP server), `desktop` (Electron),
  `web`, `landing-page`, `packaged`.
- `figma-plugin/` — a Figma **development plugin** that rebuilds an OpenDesign
  page-capture (`.od-figma.json`, produced by the OD Clipper / `od library figma`)
  into editable Figma layers (frames/text/images/fills/strokes/radii/shadows).
- `plugins/` — plugin system (`_official`, `community`, `registry`, `spec`).
- `tools/` — `dev`/`pack`/`release`/`serve` build tooling.
- `clipper/` — the OD Clipper (web page → capture, incl. Figma JSON).
- `mcp` — shipped as skills + CLI + **MCP server**; `od mcp install <agent>` wires
  it into an agent's config (`README.md` "Platform Compatibility"; MCP code in
  `apps/daemon/src/mcp-config.ts`, tests under `apps/daemon/tests/mcp-*`).

### 2.3 Frameworks / targets it addresses

- **Output artifacts are real single-page HTML** using real CSS / fonts /
  components shaped by the chosen design system's tokens — "web · desktop · mobile
  prototypes", live dashboards/artifacts, decks, images, video, HyperFrames
  motion. (`README.md` "What is Open Design" + Demo sections 1–5.)
- **React** is the framework of OD's own app UI (`@open-design/components`, React
  18.3.1) and Tailwind v4 tokens are provided per design system — directly usable
  by our React web client.
- **Consumed BY coding agents natively** (Claude Code, Codex, Cursor, Gemini CLI,
  Copilot, OpenCode, Kimi, and 21 CLIs) via one-line `od mcp install <agent>`, or
  via a BYOK proxy for any OpenAI-compatible endpoint. (`README.md` "Platform
  Compatibility".)

### 2.4 Export formats — verified from daemon/tools/host source strings

Counts of format-token occurrences in `apps/daemon/src`, `tools`, `packages/host`:

```
219 svg   195 png   129 pdf   85 webp   81 pptx   75 jpeg   46 jpg   42 mp4
```

README states the artifact export surface explicitly as **HTML / PDF / PPTX /
MP4** (plus image formats PNG/JPEG/WebP/SVG present in the pipeline). Dedicated
skills exist for `pdf`, `pptx`, `pptx-generator`, `slides`, `minimax-pdf`,
`export-download-debugging`.

**PSD:** `grep -rniE '\bpsd\b'` over `README.md`, `docs/`, `skills/`, `packages/`
returns **0** — OpenDesign has **no native PSD export**. (See §7 gap G1.)

### 2.5 Light / dark theming — verified

Theming is **DESIGN.md/token-driven, brand-by-brand**, not a single universal
auto dark mode. `tokens.css` files carry the `:root` token block; **57 of 151**
`tokens.css` files reference a dark/`data-theme`/`prefers-color-scheme` hook
(`grep -rlE 'prefers-color-scheme|data-theme|\bdark\b' design-systems/*/tokens.css
| wc -l` → 57 of 151). So: some brands ship an explicit dark token tier, others
express a single-mode brand. Any client requiring guaranteed light **and** dark
(our §11.4.162) must author/extend the dark token tier where the chosen brand
lacks one (§11.4.74 extend-upstream pattern) — do not assume every design system
is dual-mode.

---

## 3. Incorporation plan per Constitution (Task 3)

### 3.1 Submodule placement

Add OpenDesign as a **third-party/vendor submodule**:

```
<project_root>/submodules/open_design/     # §11.4.28(C) layout, §11.4.29 snake_case dir
```

- **Directory name** `open_design` (snake_case, §11.4.29). The upstream repo is
  named `open-design`; §11.4.29's vendor/upstream exception preserves the upstream
  identity — the local mount dir is snake_case, the tracked remote keeps its name.
- **Ownership class:** `nexu-io` is **not** one of our owned orgs
  (vasic-digital / HelixDevelopment / red-elf / ATMOSphere1234321 / Bear-Suite /
  BoatOS123456 / Helix-Flow / Helix-Track / Server-Factory). Therefore this is a
  **third-party dependency submodule**, NOT an "equal-codebase" owned submodule
  (§11.4.28(A) does not bind us to co-engineer it), and the §11.4.28(C)
  nested-own-org-chain prohibition does not apply (it exempts third-party
  submodules). We consume it read-only and extend via the §11.4.74
  extend-upstream-don't-reimplement pattern (fork + PR upstream) when a gap needs
  filling.
- **License compatibility:** Apache-2.0 is permissive → vendoring + redistributing
  design tokens/skills/CSS is compatible; preserve `LICENSE` + NOTICE attribution.
- **Install hygiene:** on add, run `install_upstreams` only if the repo carries an
  `upstreams/` dir (§11.4.36); OpenDesign does not, so this is a no-op — record the
  honest skip.
- **`.gitignore` (§11.4.30):** do NOT track OpenDesign build artifacts if we ever
  build it; we only need the *source design assets* (DESIGN.md, design-tokens.json,
  tokens.css, tailwind-v4.css, components.*, skills/*), so a sparse/partial usage
  is fine — we never version its `node_modules`/`dist`.

### 3.2 What we actually consume (the decoupled contract)

Per §11.4.28(B) decoupling, we consume **data assets**, not OD's app internals:

1. **Design tokens (canonical):** `submodules/open_design/design-systems/<brand>/design-tokens.json`
   (`od-design-tokens/v1`, `TOKEN_SCHEMA`) as the single source of truth for our
   brand; `tokens.css` and `tailwind-v4.css` as the ready CSS-variable / Tailwind
   bindings.
2. **Brand contract:** the chosen `DESIGN.md` (color/type/spacing/component/motion
   rules) as the human + agent spec.
3. **Skills:** relevant `skills/*/SKILL.md` (e.g. `figma-generate-design`,
   `hand-drawn-diagrams`, `pptx-generator`, `apple-hig`, `brand-extract`) invoked
   through the agent, and OD's `od mcp install` MCP server for generate/export.
4. **Craft guardrails:** `craft/accessibility-baseline.md`, `color.md`,
   `anti-ai-slop.md`, `state-coverage.md` as review rules for every UI change.
5. We do **NOT** depend on `@open-design/components` (it is `private:true`, OD's
   own app UI — not a general component lib). Our React client builds its own
   components *from the tokens*, not from that package.

### 3.3 How each client stack consumes the tokens

The token single-source is `design-tokens.json` (od-design-tokens/v1). Transpile
it once, per target, via **Style Dictionary** (or an equivalent token transformer)
so every client shares one brand definition (our brand values override the
starter tokens per §11.4.162):

| Client (R3) | Consumption path |
|---|---|
| **Web (React)** | Import `tokens.css` (`:root` CSS vars) directly, or `tailwind-v4.css` for Tailwind v4; build React components referencing `var(--accent)` etc. OD's own stack is React 18.3.1 → idiomatic fit. |
| **Desktop (Tauri/Flutter, all OS)** | Tauri (web-view) reuses the Web `tokens.css`. Flutter: transpile `design-tokens.json` → a Dart `ThemeData` / token class via Style Dictionary (Dart formatter). OD informs the *design*; it does not build the Tauri/Flutter binary. |
| **Mobile (Android / iOS / HarmonyOS / Aurora)** | Transpile `design-tokens.json` → platform token files (Android `colors.xml`/Compose theme, iOS SwiftUI/`.xcassets`, HarmonyOS ArkTS theme, Aurora/Qt QML palette). `apple-hig` skill guides iOS; **no Material/Android/Harmony/Aurora design system ships** — extend upstream (§7 gap G3). |
| **TUI** | Map the token palette to a terminal 256-color / truecolor theme (a small generator reads `design-tokens.json` colors → TUI theme). OD has no TUI target — our generator, OD's tokens. |
| **CLI** | Minimal styling: derive an ANSI accent/neutral scheme from the same tokens for help/output banners. OD tokens as the color source. |
| **REST-consumers** | No visual surface; brand applies only to any generated docs/OpenAPI HTML (reuse Web `tokens.css`). |
| **Go services** | Transpile `design-tokens.json` → a Go constants file (Style Dictionary custom format) when a Go-rendered surface (e.g. server-rendered report) needs brand colors. |

### 3.4 Brand colors + light/dark (§11.4.162)

- Pick a base design system (e.g. `default` neutral-modern, or author
  `design-systems/helixknowledge/DESIGN.md` following the 9-section schema) and set
  **our brand colors** into `--accent`/surface/neutral tokens (§11.4.162 "project
  brand colors from canonical assets").
- Ship **both light and dark** token tiers. Where a chosen OD brand lacks a dark
  tier (94/151 don't ship one), author the dark `:root[data-theme="dark"]` /
  `@media (prefers-color-scheme: dark)` block ourselves (§11.4.74 extend-upstream;
  do not assume dual-mode). Every UI component ships light + dark variants.

### 3.5 "Elements must not overlap / overlay labels" (§11.4.162)

- OpenDesign contributes **partial** enforcement: `apps/daemon/src/lint-artifact.ts`
  lints token discipline (raw-hex outside `:root` = P1, non-token accent = P0), and
  `craft/accessibility-baseline.md` + `craft/state-coverage.md` +
  `craft/anti-ai-slop.md` codify layout/contrast/state rules. This is design-system
  compliance, **not** pixel-overlap detection.
- The overlap/overlay/label-collision guarantee is enforced on **our** side by the
  constitution's own oracle: §11.4.170 device-independent host-rendered pixel proof
  + the OCR/vision layout oracle (no overlap / label-over-label / clip / off-screen)
  — the same gate §11.4.190(E) mandates for the Web client. OpenDesign feeds the
  design; our §11.4.170/§11.4.190 gates prove it renders cleanly.

---

## 4. Design-artifact production plan per client (Task 4)

R12 mandates: **wireframes, sketches, Figma design, and exported PDF / PSD / SVG
(+ other mandatory design file types)** for every client surface + the wizard flow
+ architecture diagrams. Concrete pipeline per artifact type (tool named; ✅ =
native OpenDesign capability verified in tree, ⚠️ = external tool required, OD
lacks it):

| Artifact type | Pipeline / tool | OD-native? |
|---|---|---|
| **Wireframes / low-fi** | OD `prototype`-mode skills (`web-prototype`, `dashboard`, `mobile-app`, `mobile-onboarding`) → single-page HTML at wireframe fidelity; `hand-drawn-diagrams` / `frame-flowchart-sticky` skills for sketchy wires. | ✅ |
| **Sketches / hand-drawn** | `skills/hand-drawn-diagrams`, `craft/*` guidance; OD Clipper for capturing references. | ✅ |
| **Hi-fi mockups (per client)** | OD prototype skills + chosen `DESIGN.md` tokens → real-CSS HTML mockups for Web/Desktop/Mobile-framed surfaces, sandboxed preview. | ✅ |
| **Figma design** | Generate HTML artifact → OD Clipper → `.od-figma.json` (`od library figma <assetId> --out page.od-figma.json`) → import via the **OD Figma plugin** (`figma-plugin/manifest.json`, Figma **desktop** app, one-time install) → editable Figma layers. Also skills `figma-generate-design`, `figma-generate-library`, `figma-create-design-system-rules`. NOTE: native `.fig` is NOT produced; SVG rasterizes on import (fidelity note in `figma-plugin/README.md`). | ✅ (via capture+plugin) |
| **PDF export** | OD export (HTML→PDF), skills `pdf` / `minimax-pdf`; README export surface HTML/PDF/PPTX/MP4. | ✅ |
| **PPTX / decks** | OD `deck`-mode skills (`guizang-ppt`, `html-ppt-*`, `pptx-generator`) → PPTX; `pptx-html-fidelity-audit` skill. | ✅ |
| **SVG export** | SVG is a first-class pipeline format (219 source-string hits in daemon/tools/host); vector diagrams/illustrations export to SVG. For diagram sources, render Mermaid/HTML→SVG. | ✅ (asset/diagram SVG) |
| **PSD export** | **NOT supported by OpenDesign** (0 occurrences). Produce via external pipeline: SVG/PNG → PSD using **Photopea** (headless/CLI batch), **ImageMagick** (`convert in.png out.psd`, flat), or **Gimp** batch script for layered PSD. This is a HelixKnowledge-side tool, not OD. | ⚠️ external |
| **Images (PNG/JPEG/WebP)** | OD image-mode skills (`gpt-image-2`, ImageRouter, custom API) + export pipeline. | ✅ |
| **Video / motion (MP4)** | OD `video` mode + `hyperframes` (HTML→MP4 motion graphics). | ✅ |

### 4.1 Per client surface

- **CLI / TUI:** design deliverable = a token-derived ANSI/256-color theme sheet +
  a wireframe of help/output layout (OD prototype HTML mock of the terminal
  layout, exported to PDF/SVG for the design doc). OD has no TUI renderer — theme
  generation is ours, tokens are OD's.
- **REST-consumers:** minimal visual surface → brand applies to generated API-doc
  HTML (reuse Web `tokens.css`); OpenAPI/redoc styled with the tokens.
- **Web (React):** full OD prototype pipeline → hi-fi HTML mocks → Figma
  round-trip → PDF/PNG/SVG exports; React implementation consumes `tokens.css` /
  `tailwind-v4.css` directly. **Subject to §11.4.190** (see §5).
- **Desktop (Tauri/Flutter, all OS):** OD produces the *design* (framed
  desktop-window prototypes, exports); the Tauri/Flutter *build* is ours consuming
  transpiled tokens.
- **Mobile (Android/iOS/HarmonyOS/Aurora):** OD `mobile-app`/`mobile-onboarding`
  templates give framed iPhone/Pixel HTML prototypes + `apple-hig` guidance;
  platform-native UI kits are NOT provided → author native token themes + design
  docs ourselves, guided by OD craft rules (§7 gap G3).
- **Skill-creation WIZARD flow (tech set → create → map → process):** design the
  multi-step wizard as an OD `prototype` skill run — one HTML artifact per step
  (tech-set entry, create, map, process/progress), consistent tokens across steps,
  exported to PDF for the design doc + Figma for hand-off. Use `login-flow` /
  `mobile-onboarding` templates as structural precedents for the step sequence.
- **Architecture diagrams:** author as **Mermaid** (constitution-native, renders in
  our Markdown/artifacts) as the source of truth; for polished/branded diagram
  deliverables use OD `hand-drawn-diagrams` / `frame-flowchart-sticky` skills and
  export Mermaid → **SVG** (mermaid-cli `mmdc -o diagram.svg`) → **PDF/PNG**; PSD
  via the §4 external SVG→PSD path if a layered source is demanded.

---

## 5. §11.4.190 cross-reference — Web (React) client

The Web/React client is a website surface → it is bound by §11.4.190 (website
engineering-quality mandate). OpenDesign satisfies the **design/uniqueness** inputs
but NOT the **proof** obligations — those remain our gates:

- **(A) Full responsiveness** — OD prototypes are single-page HTML; responsiveness
  across Chromium/Firefox/WebKit × device classes must be **PROVEN** by our
  device-independent host-rendered screenshots across the breakpoint × engine
  matrix + the §11.4.170 layout oracle (no overlap/clip/overflow/off-screen). OD
  does not ship this proof.
- **(B) Complete SEO** — semantic HTML + per-page title/meta, OG/Twitter cards,
  canonical, schema.org/JSON-LD, robots + sitemap, WCAG AA (OD `craft/
  accessibility-baseline.md` targets WCAG 2.2 AA) + Core Web Vitals — **PROVEN** by
  automated Lighthouse SEO + structured-data validation meeting a score floor. OD
  gives accessibility craft rules, not the audit.
- **(C) Unique OpenDesign-authored templates** — **this is exactly what R12 +
  §11.4.190(C) + §11.4.162 require**: every layout/template authored from the
  OpenDesign design-token system (our brand `DESIGN.md` + `tokens.css`), NOT a
  generic off-the-shelf template. Provenance = the OD design-token binding.
- **(D) Bleeding-edge enterprise visual quality, light + dark** — from the brand
  `DESIGN.md` tokens with both light + dark tiers authored (§3.4), **PROVEN** by
  §11.4.170 host-rendered pixel proof per screen × state × {light,dark}.
- **(E) Anti-bluff proof** — none of (A)–(D) may be claimed without the captured
  evidence above. OpenDesign is the design **source**; the §11.4.170/§11.4.190
  gates are the correctness **oracle**. A website reported responsive/SEO/unique
  without those captures is a §11.4 PASS-bluff.

---

## 6. Recommended incorporation steps (actionable, Constitution-ordered)

1. `git submodule add git@github.com:nexu-io/open-design.git submodules/open_design`
   (third-party; preserve upstream name; record in `.gitmodules`). Pin to a tag/SHA
   for reproducibility (current HEAD `94a5bd2e…`).
2. Author `design-systems/helixknowledge/DESIGN.md` (9-section schema) OR select a
   base OD system; set our brand colors into `tokens.css`/`design-tokens.json`;
   author the **dark** tier (§3.4).
3. Add a **Style Dictionary** transform step: `design-tokens.json` → Web CSS vars
   (reuse OD `tokens.css`), Tailwind v4, Flutter Dart theme, Android/iOS/Harmony/
   Aurora token files, Go constants, TUI/CLI ANSI theme.
4. `od mcp install <our agent>` to wire OD's MCP server for generate/export; adopt
   the relevant skills (prototype/deck/figma/diagram/pdf/pptx).
5. Stand up the **external PSD pipeline** (SVG/PNG → PSD via Photopea/Gimp/
   ImageMagick) since OD lacks PSD (gap G1).
6. Wire the Web client through §11.4.190 gates (responsive + SEO + light/dark pixel
   proof) and §11.4.170 layout oracle for the overlap/label guarantee.
7. Track every OD gap (G1–G5 below) as a §11.4.197 workable item so no requirement
   is left un-wired.

---

## 7. Honest gaps & risks (zero-bluff)

- **G1 — PSD export unsupported.** OpenDesign has **no** PSD path (0 occurrences).
  R12 lists PSD as mandatory → must be produced by an external tool (Photopea /
  Gimp batch for layered PSD; ImageMagick for flat). Flag: layered-PSD fidelity
  from HTML/SVG is imperfect; confirm the operator's real PSD need (source-of-truth
  vs deliverable) before investing.
- **G2 — Figma is round-trip, not native `.fig`.** OD produces `.od-figma.json`
  captures re-imported via the OD Figma **plugin** (Figma desktop only), and SVGs
  **rasterize** on import (`figma-plugin/README.md` fidelity notes). No native
  `.fig` file is generated. Editable Figma layers are achievable but with fidelity
  caveats (gradients beyond first layer, blend modes, transforms, SVG internals
  simplified).
- **G3 — Mobile = framed HTML prototypes, not native platform kits.** OD ships
  `mobile-app`/`mobile-onboarding` (iPhone 15 Pro / Pixel framed HTML) + an
  `apple-hig` skill, but **no** Android/Material, HarmonyOS, or Aurora OS design
  system/skill. Native iOS/Android/HarmonyOS/Aurora design for R3 mobile requires
  authoring those token themes + design docs ourselves (extend-upstream §11.4.74).
- **G4 — `@open-design/components` is not a consumable lib.** It is `private:true`
  React primitives for OD's OWN app UI. Do not depend on it as our component
  library; build from tokens.
- **G5 — Dark mode is brand-dependent.** Only 57/151 `tokens.css` reference a
  dark/theme hook. §11.4.162's light+dark requirement means we author the dark tier
  where the chosen brand lacks one — do not assume dual-mode.
- **G6 — Desktop-app OS coverage of OD itself.** OD's bundled desktop app targets
  macOS + Windows (Electron; AMR CLI distribution slice initially mac arm64 only per
  `CONTEXT.md`). This does not block us (we consume design assets + MCP on Linux via
  source/Docker), but OD is not itself a "runs everywhere native desktop" product.
- **Boundary note (§11.4.6):** everything above about OpenDesign is verified from
  the cloned tree at `94a5bd2e`. Anything the operator adds beyond this
  (e.g., OpenDesign Cloud paid model routing, private OD features) is NOT verified
  here and must not be assumed.

---

## 8. Verification appendix (commands run, all read-only)

- `git ls-remote git@github.com:nexu-io/open-design.git HEAD` → `94a5bd2e…` (exit 0)
- `git ls-remote https://github.com/nexu-io/open-design.git HEAD` → `94a5bd2e…` (exit 0)
- `git clone --depth 1 --filter=blob:limit=200k …` → 11,373 files, branch `main`,
  HEAD commit "fix BYOK OpenCode permission bypass (#5701)" 2026-07-15.
- `cat package.json` → name `open-design` v0.15.1, license `Apache-2.0`, bin `od`.
- `ls design-systems | wc -l` → 153; `ls skills | wc -l` → 164; `ls design-templates
  | wc -l` → 115.
- `sed -n design-systems/apple/design-tokens.json` → `od-design-tokens/v1`,
  `TOKEN_SCHEMA`, 56 typed tokens.
- `cat design-systems/apple/tokens.css` / `USAGE.md` → `:root` token block, lint at
  `apps/daemon/src/lint-artifact.ts`.
- `grep -rhoiE '(pptx|svg|png|pdf|mp4|jpeg|jpg|webp)' apps/daemon/src tools
  packages/host` → svg 219, png 195, pdf 129, webp 85, pptx 81, jpeg 75, jpg 46,
  mp4 42.
- `grep -rniE '\bpsd\b' README.md docs skills packages` → 0.
- `grep -rlE 'prefers-color-scheme|data-theme|dark' design-systems/*/tokens.css |
  wc -l` → 57 (of 151).
- `cat figma-plugin/README.md` → capture (`.od-figma.json`) + plugin import model,
  SVG rasterized on import, no native `.fig`.
- README "Platform Compatibility" → `od mcp install <agent>` for 21 agents; MCP
  server shipped.

*(Clone lives only in the session scratchpad; nothing was written into the project
tree except this file.)*
