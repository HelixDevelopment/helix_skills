## G01. Two rival API servers; hardened internal/api is unwired dead code

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** critical
**Created-By:** Claude

Runtime security hole CLOSED (2026-07-15). Dead server code removed (2026-07-18, -1826 lines). internal/api coverage 41.6% -> 59.8%.

## G02. Sandbox provides no isolation; default path executes arbitrary code on host (RCE)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** critical
**Created-By:** Claude

Latent RCE primitive. Fixed: sandbox now fails closed when no WASM runtime found.

## G03. Flagship pipelines (validation jury + autoexpand) are dead code

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** critical
**Created-By:** Claude

internal/validation and internal/autoexpand never instantiated. Fixed: both pipelines now wired.

## G04. Insufficient test coverage across the codebase

**Status:** In progress
**Type:** Task
**Severity:** critical
**Created-By:** Claude

Tests exist (144 files, 27/27 packages GREEN). Coverage boosted: models 0%->100%, skill 5.5%->8.3%, api 41.6%->59.8%. Remaining: registry (2%), db (12%), worker (0%) — DB-dependent.

## G05. LLMJury auto-approves when no jury is configured (default state)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Zero-bluff gate defaults to approve-everything. Fixed: fail-closed when unconfigured.

## G06. GetDependencyTree returns only depth-1 children (recursive tree truncated)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Core recursive dependency DAG feature broken. Fixed: recursive CTE with depth bound.

## G07. TOML/JSON dependency+resource round-trip broken (edges silently dropped)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Breaks R14 git-versionable round-trip. Fixed: all 6 relation types + components carried through.

## G08. No TOON codec exists; wire format mandate unmet

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

R/founding mandate TOON primary + JSON fallback. Fixed: TOON codec implemented.

## G09. Pervasive OpenAPI-implementation drift; most endpoints unimplemented or differently shaped

**Status:** Queued
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Contract is SSoT for R3 thin clients. Clients generated from it will not work.

## G10. Embedding dimension: no model-column assertion; vector(768) hard-coded

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

768/1536/384 conflict resolved only by unenforced constants. Fixed: dynamic dimension from config.

## G11. Worker does no real work and can panic (unchecked type assertions)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Background auto-growth/validation/review non-functional. Fixed: proper error handling + recover.

## G12. tree-sitter is a stub: native parsing always fails; regex-only

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

R2 requires working tree-sitter POC. Fixed: native parsing works for Kotlin, C#, and more.

## G13. Two rival docker-compose.yml files (rival-copy risk)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

P12.T4 requires one canonical compose. Fixed: consolidated to single file.

## G14. Submodule policy conflict unresolved; all 7 declared deps unvendored

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

Open governance escalation blocks P3/P7/P8/P10/P11/P13 dependency work.

## G15. Aurora OS / HarmonyOS client feasibility unproven

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

R3 hard-requires Aurora + HarmonyOS clients. Plan ranks this top danger.

## G29. Store.Search advertises hybrid vector search but is trigram/ILIKE-only

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Store.VectorSearch had zero callers. Fixed: hybrid search now wired.

## G31. learn_from_project project_path has zero validation (path-traversal/LFI)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Latent path-traversal when G03 wiring lands. Fixed: path validation added.

## G32. registry.ReviewScheduler fully built but has zero callers (dead pipeline)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Dead flagship pipeline. Fixed: scheduler now wired.

## G35. CLI+TUI send Authorization:Bearer but server reads X-API-Key

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Both first-party clients 401 when auth enforced. Fixed: unified auth header.

## G57. MCP ACPAdapter stdio transport was unwired

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Fixed: stdio transport now wired (commit 8fa4e27).

## G59. Embedding ingestion never wired; StoreSkillEmbedding dead code

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Fixed (2026-07-17): Store.Create calls embedWriteThrough, ClearSkillEmbedding on failure branches.

## G63. 4th divergent route-contract surface (registry CLI/TUI)

**Status:** Operator-blocked
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Operator-blocked on 5 product/ownership decisions (D1-D5).

## G69. GitHub Skills Source Ingestion epic item G69

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G70. GitHub Skills Source Ingestion epic item G70

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G71. GitHub Skills Source Ingestion epic item G71

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G72. GitHub Skills Source Ingestion epic item G72

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G73. GitHub Skills Source Ingestion epic item G73

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G74. GitHub Skills Source Ingestion epic item G74

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G75. GitHub Skills Source Ingestion epic item G75

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G76. GitHub Skills Source Ingestion epic item G76

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G77. GitHub Skills Source Ingestion epic item G77

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G78. GitHub Skills Source Ingestion epic item G78

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G79. GitHub Skills Source Ingestion epic item G79

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G80. GitHub Skills Source Ingestion epic item G80

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G81. GitHub Skills Source Ingestion epic item G81

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G82. GitHub Skills Source Ingestion epic item G82

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G83. GitHub Skills Source Ingestion epic item G83

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G84. GitHub Skills Source Ingestion epic item G84

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G85. GitHub Skills Source Ingestion epic item G85

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G86. GitHub Skills Source Ingestion epic item G86

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G87. GitHub Skills Source Ingestion epic item G87

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G88. GitHub Skills Source Ingestion epic item G88

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G89. GitHub Skills Source Ingestion epic item G89

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G90. GitHub Skills Source Ingestion epic item G90

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G91. GitHub Skills Source Ingestion epic item G91

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G92. GitHub Skills Source Ingestion epic item G92

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G93. Unified Multi-Source Skill Ingestion epic item G93

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G94. Unified Multi-Source Skill Ingestion epic item G94

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G95. Unified Multi-Source Skill Ingestion epic item G95

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G96. Unified Multi-Source Skill Ingestion epic item G96

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G97. Unified Multi-Source Skill Ingestion epic item G97

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G98. Unified Multi-Source Skill Ingestion epic item G98

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G99. Unified Multi-Source Skill Ingestion epic item G99

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G100. Unified Multi-Source Skill Ingestion epic item G100

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G101. Unified Multi-Source Skill Ingestion epic item G101

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G102. Unified Multi-Source Skill Ingestion epic item G102

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G103. Unified Multi-Source Skill Ingestion epic item G103

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G104. Unified Multi-Source Skill Ingestion epic item G104

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G105. Unified Multi-Source Skill Ingestion epic item G105

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G106. Unified Multi-Source Skill Ingestion epic item G106

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G107. Unified Multi-Source Skill Ingestion epic item G107

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G108. Unified Multi-Source Skill Ingestion epic item G108

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G109. Unified Multi-Source Skill Ingestion epic item G109

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G110. Unified Multi-Source Skill Ingestion epic item G110

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G111. Unified Multi-Source Skill Ingestion epic item G111

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G112. Unified Multi-Source Skill Ingestion epic item G112

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G113. Unified Multi-Source Skill Ingestion epic item G113

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G114. Unified Multi-Source Skill Ingestion epic item G114

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G115. Unified Multi-Source Skill Ingestion epic item G115

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G116. Unified Multi-Source Skill Ingestion epic item G116

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G117. Unified Multi-Source Skill Ingestion epic item G117

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G118. Unified Multi-Source Skill Ingestion epic item G118

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G119. Unified Multi-Source Skill Ingestion epic item G119

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G120. Unified Multi-Source Skill Ingestion epic item G120

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G121. Unified Multi-Source Skill Ingestion epic item G121

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G122. Unified Multi-Source Skill Ingestion epic item G122

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G123. G69 vs G93 architectural-overlap reconciliation

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Reconcile GitHub Skills Source Ingestion epic with Unified Multi-Source Skill Ingestion epic — overlapping scope.

## G124. Auto-generated skills-tree documentation catalog item G124

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G125. Auto-generated skills-tree documentation catalog item G125

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G126. Auto-generated skills-tree documentation catalog item G126

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G127. Auto-generated skills-tree documentation catalog item G127

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G128. Auto-generated skills-tree documentation catalog item G128

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G129. Auto-generated skills-tree documentation catalog item G129

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G130. Auto-generated skills-tree documentation catalog item G130

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G131. Auto-generated skills-tree documentation catalog item G131

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G132. Auto-generated skills-tree documentation catalog item G132

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G133. Auto-generated skills-tree documentation catalog item G133

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G134. Auto-generated skills-tree documentation catalog item G134

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G135. Auto-generated skills-tree documentation catalog item G135

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G137. Additional enterprise features and integrations

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Enterprise-grade features for the skill system — additional integrations and capabilities.

## R01. Standard: build+run clean, 4-layer test coverage, paired mutations

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

Build+run clean. 4-layer test coverage + paired mutations per Constitution §1/§1.1/§11.4. Fix flagged security defects. Dedupe to one canonical tree.

## R02. Universal + dynamic: create skills on demand for any technology

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Not Android-specific. Create skills on demand for any technology set. Mentioned tech (Gin, quic-go, pgvector, tree-sitter, etc) must be working POCs, not stubs.

## R03. Clients: CLI, TUI, REST, Web, Desktop, Mobile

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

CLI, TUI, REST API, Web, Desktop (Win/macOS/Linux), Mobile (Android, iOS, HarmonyOS, Aurora OS). Maximize shared codebase across all surfaces.

## R04. Agent interop: Claude Code, OpenCode, Kimi Code, HelixTrack

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Usable from Claude Code (toolkit + aliases), OpenCode, Kimi Code, and HelixTrack/HelixAgent/HelixLLM.

## R05. Incorporation analysis: research integration of R4 agents

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

In-depth research on how to integrate all of R4 properly. 4 research agents in flight.

## R06. Canonical use case: skill-creation wizard

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

User opens client → skill-creation wizard → enters tech set → submit → backend runs create → map (DAG) → full processing → progress reported back.

## R07. Model access: pluggable ModelProvider, not hardcoded OpenAI

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Obtain quality models via LLMsVerifier, HelixLLM, Claude Toolkit aliases. Pluggable ModelProvider. Vendored submodules per Constitution naming.

## R08. Exhaustive testing: all test types + Challenges + HelixQA

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

Every unit covered by ALL supported test types PLUS Challenges and HelixQA test banks/suites. Coverage as close to 100% as possible.

## R09. Submodule resolution + sync: parent-dir versions have priority

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

All deps live under submodules/<snake_case>/ OR reuse from parent dir. Parent-dir versions have PRIORITY. Both copies always in sync.

## R10. Docs Chain: fully incorporate the Docs Chain submodule

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

Fully incorporate the Docs Chain submodule (.docs_chain) per Constitution §11.4.106.

## R11. Zero bluff anywhere: positive-evidence-only everywhere

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

No false results, faulty results, or faulty codebase. No bluff of any kind. Positive-evidence-only everywhere. Governs every gate, test, and status report.

## R12. OpenDesign for all design artifacts

**Status:** Queued
**Type:** Feature
**Severity:** medium
**Created-By:** Claude

Every client design/styling/diagrams/illustrations MUST use OpenDesign. Deliver wireframes, sketches, Figma, exports in PDF/PSD/SVG.

## R13. Validation skill corpus: 35+ technologies must validate end-to-end

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

System must create + validate corpus: android, java, kotlin, python, go, typescript, linux, etc. This IS the end-to-end proof.

## R14. Git-versioned + real-time growth via MCP/ACP/plugins

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Every skill is Git-versionable (TOML per SPEC §4.2). MCP/ACP/plugins trigger real-time skill creation + deep research. Improvements available immediately.

## R15. systemctl user scope + ops scripts integration

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

System MUST integrate via systemctl (user scope) + ops scripts for deployment and management.

## R16. Skill granularity and composition: decomposable big technologies

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Big technologies MUST be decomposable into smaller composable skills. Umbrella/composite/atomic skill kinds.

## R17. Exhaustive gaps/risks remediation + total test coverage

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

ALL gaps, weak spots, danger-zones MUST be investigated and fixed. Total test coverage across all supported types.

## R18. Full documentation delivery, always in sync

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

Whole project MUST ship complete documentation always in sync. §11.4.12/§11.4.65 compliance.

## R19. Anthropic API support as first-class ModelProvider

**Status:** Queued
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Besides full OpenAI support, Anthropic Messages API MUST be a first-class ModelProvider with proper client implementation.

## R20. Containers submodule for all containerization

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

All containerization MUST use the containers submodule per Constitution §11.4.76. Rootless podman only.

## R21. Adopt reusable Helix-family practices

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Survey the Helix ecosystem and adopt reusable practices, patterns, and components from sibling projects.

## R22. Catalogue-first incorporation from vasic-digital + HelixDevelopment

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Before building new components, check existing catalogues in vasic-digital and HelixDevelopment orgs. Extend-don't-reimplement per §11.4.74.

## R23. Full constitutional compliance, no violations, no bluff

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

Every aspect of the project MUST comply with the Helix Constitution. No violations, no bluff, no shortcuts.

## R24. Every operator request always respected + recorded, zero request-loss

**Status:** Queued
**Type:** Task
**Severity:** high
**Created-By:** Claude

Every operator request/prompt MUST be captured, tracked, and processed. §11.4.197/§11.4.202/§11.4.210 compliance.

## R25. Canonical project key hxs, used everywhere

**Status:** Queued
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Project key 'hxs' MUST be used consistently across all configurations, prefixes, and identifiers.

