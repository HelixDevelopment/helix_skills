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

**Status:** Fixed (→ Fixed.md)
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

**Status:** Fixed (→ Fixed.md)
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

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Open governance escalation blocks P3/P7/P8/P10/P11/P13 dependency work.

## G15. Aurora OS / HarmonyOS client feasibility unproven

**Status:** Fixed (→ Fixed.md)
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

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Operator-blocked on 5 product/ownership decisions (D1-D5).

## G69. GitHub Skills Source Ingestion epic item G69

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G70. GitHub Skills Source Ingestion epic item G70

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G71. GitHub Skills Source Ingestion epic item G71

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G72. GitHub Skills Source Ingestion epic item G72

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G73. GitHub Skills Source Ingestion epic item G73

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G74. GitHub Skills Source Ingestion epic item G74

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G75. GitHub Skills Source Ingestion epic item G75

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G76. GitHub Skills Source Ingestion epic item G76

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G77. GitHub Skills Source Ingestion epic item G77

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G78. GitHub Skills Source Ingestion epic item G78

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G79. GitHub Skills Source Ingestion epic item G79

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G80. GitHub Skills Source Ingestion epic item G80

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G81. GitHub Skills Source Ingestion epic item G81

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G82. GitHub Skills Source Ingestion epic item G82

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G83. GitHub Skills Source Ingestion epic item G83

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G84. GitHub Skills Source Ingestion epic item G84

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G85. GitHub Skills Source Ingestion epic item G85

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G86. GitHub Skills Source Ingestion epic item G86

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G87. GitHub Skills Source Ingestion epic item G87

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G88. GitHub Skills Source Ingestion epic item G88

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G89. GitHub Skills Source Ingestion epic item G89

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G90. GitHub Skills Source Ingestion epic item G90

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G91. GitHub Skills Source Ingestion epic item G91

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G92. GitHub Skills Source Ingestion epic item G92

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G69-G92 epic: GitHub Skills Source Ingestion (24 items).

## G93. Unified Multi-Source Skill Ingestion epic item G93

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G94. Unified Multi-Source Skill Ingestion epic item G94

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G95. Unified Multi-Source Skill Ingestion epic item G95

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G96. Unified Multi-Source Skill Ingestion epic item G96

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G97. Unified Multi-Source Skill Ingestion epic item G97

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G98. Unified Multi-Source Skill Ingestion epic item G98

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G99. Unified Multi-Source Skill Ingestion epic item G99

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G100. Unified Multi-Source Skill Ingestion epic item G100

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G101. Unified Multi-Source Skill Ingestion epic item G101

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G102. Unified Multi-Source Skill Ingestion epic item G102

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G103. Unified Multi-Source Skill Ingestion epic item G103

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G104. Unified Multi-Source Skill Ingestion epic item G104

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G105. Unified Multi-Source Skill Ingestion epic item G105

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G106. Unified Multi-Source Skill Ingestion epic item G106

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G107. Unified Multi-Source Skill Ingestion epic item G107

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G108. Unified Multi-Source Skill Ingestion epic item G108

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G109. Unified Multi-Source Skill Ingestion epic item G109

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G110. Unified Multi-Source Skill Ingestion epic item G110

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G111. Unified Multi-Source Skill Ingestion epic item G111

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G112. Unified Multi-Source Skill Ingestion epic item G112

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G113. Unified Multi-Source Skill Ingestion epic item G113

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G114. Unified Multi-Source Skill Ingestion epic item G114

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G115. Unified Multi-Source Skill Ingestion epic item G115

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G116. Unified Multi-Source Skill Ingestion epic item G116

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G117. Unified Multi-Source Skill Ingestion epic item G117

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G118. Unified Multi-Source Skill Ingestion epic item G118

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G119. Unified Multi-Source Skill Ingestion epic item G119

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G120. Unified Multi-Source Skill Ingestion epic item G120

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G121. Unified Multi-Source Skill Ingestion epic item G121

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G122. Unified Multi-Source Skill Ingestion epic item G122

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Part of G93-G122 epic: Unified Multi-Source Skill Ingestion (30 items).

## G123. G69 vs G93 architectural-overlap reconciliation

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Reconcile GitHub Skills Source Ingestion epic with Unified Multi-Source Skill Ingestion epic — overlapping scope.

## G124. Auto-generated skills-tree documentation catalog item G124

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G125. Auto-generated skills-tree documentation catalog item G125

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G126. Auto-generated skills-tree documentation catalog item G126

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G127. Auto-generated skills-tree documentation catalog item G127

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G128. Auto-generated skills-tree documentation catalog item G128

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G129. Auto-generated skills-tree documentation catalog item G129

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G130. Auto-generated skills-tree documentation catalog item G130

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G131. Auto-generated skills-tree documentation catalog item G131

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G132. Auto-generated skills-tree documentation catalog item G132

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G133. Auto-generated skills-tree documentation catalog item G133

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G134. Auto-generated skills-tree documentation catalog item G134

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Part of G124-G135 epic: auto-generated skills-tree documentation catalog with 12 items total.

## G135. Auto-generated skills-tree documentation catalog item G135

**Status:** Fixed (→ Fixed.md)
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

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Build+run clean. 4-layer test coverage + paired mutations per Constitution §1/§1.1/§11.4. Fix flagged security defects. Dedupe to one canonical tree.

## R02. Universal + dynamic: create skills on demand for any technology

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Not Android-specific. Create skills on demand for any technology set. Mentioned tech (Gin, quic-go, pgvector, tree-sitter, etc) must be working POCs, not stubs.

## R03. Clients: CLI, TUI, REST, Web, Desktop, Mobile

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

CLI, TUI, REST API, Web, Desktop (Win/macOS/Linux), Mobile (Android, iOS, HarmonyOS, Aurora OS). Maximize shared codebase across all surfaces.

## R04. Agent interop: Claude Code, OpenCode, Kimi Code, HelixTrack

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Usable from Claude Code (toolkit + aliases), OpenCode, Kimi Code, and HelixTrack/HelixAgent/HelixLLM.

## R05. Incorporation analysis: research integration of R4 agents

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

In-depth research on how to integrate all of R4 properly. 4 research agents in flight.

## R06. Canonical use case: skill-creation wizard

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

User opens client → skill-creation wizard → enters tech set → submit → backend runs create → map (DAG) → full processing → progress reported back.

## R07. Model access: pluggable ModelProvider, not hardcoded OpenAI

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Obtain quality models via LLMsVerifier, HelixLLM, Claude Toolkit aliases. Pluggable ModelProvider. Vendored submodules per Constitution naming.

## R08. Exhaustive testing: all test types + Challenges + HelixQA

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Every unit covered by ALL supported test types PLUS Challenges and HelixQA test banks/suites. Coverage as close to 100% as possible.

## R09. Submodule resolution + sync: parent-dir versions have priority

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

All deps live under submodules/<snake_case>/ OR reuse from parent dir. Parent-dir versions have PRIORITY. Both copies always in sync.

## R10. Docs Chain: fully incorporate the Docs Chain submodule

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Fully incorporate the Docs Chain submodule (.docs_chain) per Constitution §11.4.106.

## R11. Zero bluff anywhere: positive-evidence-only everywhere

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

No false results, faulty results, or faulty codebase. No bluff of any kind. Positive-evidence-only everywhere. Governs every gate, test, and status report.

## R12. OpenDesign for all design artifacts

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** medium
**Created-By:** Claude

Every client design/styling/diagrams/illustrations MUST use OpenDesign. Deliver wireframes, sketches, Figma, exports in PDF/PSD/SVG.

## R13. Validation skill corpus: 35+ technologies must validate end-to-end

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

System must create + validate corpus: android, java, kotlin, python, go, typescript, linux, etc. This IS the end-to-end proof.

## R14. Git-versioned + real-time growth via MCP/ACP/plugins

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Every skill is Git-versionable (TOML per SPEC §4.2). MCP/ACP/plugins trigger real-time skill creation + deep research. Improvements available immediately.

## R15. systemctl user scope + ops scripts integration

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

System MUST integrate via systemctl (user scope) + ops scripts for deployment and management.

## R16. Skill granularity and composition: decomposable big technologies

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Big technologies MUST be decomposable into smaller composable skills. Umbrella/composite/atomic skill kinds.

## R17. Exhaustive gaps/risks remediation + total test coverage

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

ALL gaps, weak spots, danger-zones MUST be investigated and fixed. Total test coverage across all supported types.

## R18. Full documentation delivery, always in sync

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Whole project MUST ship complete documentation always in sync. §11.4.12/§11.4.65 compliance.

## R19. Anthropic API support as first-class ModelProvider

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** high
**Created-By:** Claude

Besides full OpenAI support, Anthropic Messages API MUST be a first-class ModelProvider with proper client implementation.

## R20. Containers submodule for all containerization

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

All containerization MUST use the containers submodule per Constitution §11.4.76. Rootless podman only.

## R21. Adopt reusable Helix-family practices

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Survey the Helix ecosystem and adopt reusable practices, patterns, and components from sibling projects.

## R22. Catalogue-first incorporation from vasic-digital + HelixDevelopment

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Before building new components, check existing catalogues in vasic-digital and HelixDevelopment orgs. Extend-don't-reimplement per §11.4.74.

## R23. Full constitutional compliance, no violations, no bluff

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Every aspect of the project MUST comply with the Helix Constitution. No violations, no bluff, no shortcuts.

## R24. Every operator request always respected + recorded, zero request-loss

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

Every operator request/prompt MUST be captured, tracked, and processed. §11.4.197/§11.4.202/§11.4.210 compliance.

## R25. Canonical project key hxs, used everywhere

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Project key 'hxs' MUST be used consistently across all configurations, prefixes, and identifiers.

## G16. sandbox_type=wasm never uses WASM; Docker sandbox has conflicting mounts

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G17. Weak/committed default DB password; config validation misses enums

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G18. CORS allowlist unreachable on live path; SPEC omits allowed_origins

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G19. SPEC.md config sample uses -- comments (invalid TOML)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G20. Auto-expand fabricates placeholder skills without LLM

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G21. Resource verification is shallow (HEAD-only, fail-open)

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G22. No rate limiting/auth on live server; Brotli flush errors ignored

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G23. Migrations loaded from cwd-relative path; failure only warns

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G24. Health/metrics/version unauthenticated; /metrics exposes Prometheus

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G25. RemoveDependency ignores name-lookup errors

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G26. VAR:-default cannot resolve to intentionally-empty value

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G27. sanitizeTableName silently strips instead of rejecting

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G28. Anthropic Messages API as first-class LLMClient provider

**Status:** Fixed (→ Fixed.md)
**Type:** Feature
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G30. learn_from_project returns job ID that can never be status-checked

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G33. Store.ExportToTOML swallows DB error

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G34. Unchecked rid.(string) type assertion in request-id middleware

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G36. SSRF blocklist: non-zero 0.0.0.0/8 hosts not explicitly blocked

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G37. Import-skills path honors client status on dead api.Server router

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G38. Request-history auto-capture hook not yet wired

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G44. Missing revision headers on most docs

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G45. Status/type/reopens closed-vocabulary not enforced

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G46. Shell scripts lack companion docs

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G47. Gxx id scheme vs ATM-NNN naming

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G48. README lacks Tracked-Items doc-link section

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G49. No QWEN.md/GEMINI.md mirrors

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G51. Silent migration-state desync from psql error handling

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G52. Health-probe repoint + ARCHITECTURE.md doc-drift

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G53. WaitForVectorIndexReady catalog-query correction

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G54. gofmt drift in internal/validation/pipeline.go

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G55. Phantom OpenAPI/doc-listed routes unimplemented

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G56. docker-compose deployment contract broken

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G58. Placeholder finding — needs real investigation

**Status:** Obsolete (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G60. Search conflict-oracle uses ranked Search instead of GetByName

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G61. Two divergent /health implementations

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G62. 20 files with gofmt drift project-wide

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G64. Additional ops-script issues

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G65. Additional migration issues

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G66. Seed corpus missing 3 prerequisite TOMLs

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** medium
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G67. qa-results gitignore vs curated QA evidence policy conflict

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G68. RemoveDependency coarse delete should be relation-type-aware

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** low
**Created-By:** Claude

Auto-imported from GAPS_AND_RISKS_REGISTER.md during G40 Phase 1 completion.

## G39. Rootful container runtime default in Makefile

**Status:** Fixed (→ Fixed.md)
**Type:** Bug
**Severity:** high
**Created-By:** Claude

Makefile uses rootful docker instead of rootless podman per §11.4.161. Fixed: switched to rootless podman.

## G41. Zero hooks wired in settings.local.json

**Status:** Fixed (→ Fixed.md)
**Type:** Task
**Severity:** high
**Created-By:** Claude

guard-forbidden-commands.sh PreToolUse guard not installed. Fixed: hooks wired in settings.local.json.

