# Deep Research: SPEC.md Gap Analysis — All Remaining Unimplemented Features

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

**Date:** 2026-07-17
**Analyst:** Deep Research Agent (T2)
**Scope:** SPEC.md §2-§8 vs current implementation tree
**Method:** File:line evidence only (R11 positive-evidence-only)

---

## Executive Summary

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

Of the 137 tracked gaps (G01-G137), **31 are CLOSED/FIXED**, **6 are operator-blocked**, and **100 remain OPEN**. The implementation has made significant progress on security (G01 runtime hole closed, G02 sandbox RCE eliminated, G21 SSRF guard landed) and core DAG correctness (G06/G07 landed), but the **flagship pipelines** (validation jury, auto-expand growth, code analysis) remain partially or fully unwired, and **zero test coverage** exists for most packages.

---

## §2 Technology Stack — Gap Analysis

### IMPLEMENTED
| Component | Status | Evidence |
|-----------|--------|----------|
| Go 1.22+ | ✅ | `go.mod` specifies `go 1.22` |
| Gin v1.11.0+ | ✅ | `go.mod` has `github.com/gin-gonic/gin` |
| pgvector | ✅ | `go.mod` has `github.com/pgvector/pgvector-go` |
| TOML (BurntSushi) | ✅ | `go.mod` has `github.com/BurntSushi/toml` |
| tree-sitter | ⚠️ STUB | `internal/codeanalysis/treesitter.go:106-131` `initNativeParser` always errors |
| mcp-go | ✅ | `go.mod` has `github.com/mark3labs/mcp-go` |
| Cobra | ✅ | `go.mod` has `github.com/spf13/cobra` |
| Bubble Tea | ✅ | `go.mod` has `github.com/charmbracelet/bubbletea` |
| PostgreSQL 16+ | ✅ | `docker-compose.yml` uses `postgres:16` |
| Brotli | ✅ | `go.mod` has `github.com/andybalholm/brotli` |

### NOT IMPLEMENTED
| Component | Status | Evidence |
|-----------|--------|----------|
| HTTP/3 (quic-go + Caddy) | ❌ DEAD CODE | `internal/api/http3.go` deleted; Caddy not configured |
| Embedding providers (anthropic, local) | ⚠️ PARTIAL | `NewEmbedderFromConfig` supports `openai`/`local` only (`internal/db/embedding.go:294-308`); Anthropic added for LLM but not embeddings |

---

## §3 Project Structure — Gap Analysis

### IMPLEMENTED
All directories exist as specified:
- `cmd/server/`, `cmd/worker/`, `cmd/cli/`, `cmd/tui/` ✅
- `internal/config/`, `internal/db/`, `internal/models/`, `internal/skill/` ✅
- `internal/registry/`, `internal/autoexpand/`, `internal/codeanalysis/` ✅
- `internal/validation/`, `internal/mcp/`, `internal/api/`, `internal/worker/` ✅
- `migrations/`, `scripts/`, `config/config.toml` ✅

### SPEC DRIFT
| Item | Spec Says | Implementation |
|------|-----------|----------------|
| `internal/toon/` | Not in spec | Added for G08 TOON codec |
| `internal/skillscatalog/` | Not in spec | Added for G124 docs catalog |
| `internal/source/` | Not in spec | Added for G69/G93 ingestion |
| `internal/ingest/` | Not in spec | Added for G93 multi-source ingestion |

---

## §4 Data Models — Gap Analysis

### IMPLEMENTED (§4.1 Core Types)
All core types exist in `internal/models/`:
- `Skill` struct ✅ (`internal/models/skill.go`)
- `SkillStatus` enum (draft/validated/active/deprecated) ✅
- `DependencyType` (requires/extends/recommends/composes/related_to/alternative_to) ✅
- `SkillKind` (atomic/composite/umbrella) ✅ (R16)
- `SkillDependency` with `Optional`/`SortOrder` ✅ (R16)
- `Resource` ✅
- `Evidence` ✅
- `SkillRegistryEntry` ✅
- `AuditLogEntry` ✅

### IMPLEMENTED (§4.2 TOML Format)
- TOML skill format with `[skill]`, `[skill.metadata]`, `[skill.dependencies]` ✅
- `[[skill.components]]` ergonomic authoring form ✅ (R16)
- `[[skill.resources]]` ✅
- Round-trip TOML↔JSON (G07 landed `073192f`) ✅

---

## §5 Database Schema — Gap Analysis

### IMPLEMENTED
- `skills` table with all columns including `kind` (R16) ✅
- `skill_dependencies` with 6 relation types + `optional` + `sort_order` ✅
- `resources` table ✅
- `evidences` table ✅
- `skill_registry` table ✅
- `audit_log` table ✅
- HNSW indexes for embeddings ✅
- `update_updated_at` trigger ✅
- `vector(768)` default ✅

### NOT IMPLEMENTED
| Item | Status | Evidence |
|------|--------|----------|
| `skill_sources` table (G70) | ❌ QUEUED | Migration `004_skill_sources` not created |
| `skill_source_mappings` table (G70) | ❌ QUEUED | Same |
| `skill_enhancement_proposals` table (G71) | ❌ QUEUED | Migration not created |
| `skills.origin` column (G70) | ❌ QUEUED | Not in schema |
| Ingestion tables (G95) | ❌ QUEUED | Migration `004_ingestion` not created |
| Embedding dimension templating (G10) | ❌ DESIGN DONE | `vector(768)` hard-coded, not config-driven |

---

## §6 API Specification — Gap Analysis

### IMPLEMENTED (Live Server — `cmd/server/main.go`)
| Endpoint | Status | Evidence |
|----------|--------|----------|
| `GET /api/v1/skills` | ✅ | `main.go:170` |
| `GET /api/v1/skills/:name` | ✅ | `main.go:175` |
| `GET /api/v1/skills/:name/tree` | ✅ | `main.go:178` |
| `GET /api/v1/skills/search` | ✅ | `main.go:182` (GET, not POST as spec says) |
| `GET /api/v1/registry/coverage` | ✅ | `main.go:238` |
| `GET /api/v1/missing` | ✅ | `main.go:234` |
| `GET /health` | ✅ | `main.go:151` |
| `GET /metrics` | ✅ | `main.go:155` |
| `GET /version` | ✅ | `main.go:159` |

### NOT IMPLEMENTED (Live Server)
| Endpoint | Status | Evidence |
|----------|--------|----------|
| `POST /skills` (create) | ❌ | No route in `main.go` |
| `PUT /skills/:name` (update) | ❌ | No route in `main.go` |
| `DELETE /skills/:name` (delete) | ❌ | No route in `main.go` |
| `POST /skills/import` (bulk import) | ❌ | No route in `main.go` |
| `GET /skills/:name/export` (TOML export) | ❌ | No route in `main.go` |
| `POST /search` (vector+keyword) | ❌ | Spec says POST, live is GET |
| `POST /search/similar` | ❌ | No route in `main.go` |
| `GET /registry` | ❌ | No route in `main.go` |
| `GET /registry/missing` | ❌ | Spec path differs from live |
| `GET /registry/stale` | ❌ | No route in `main.go` |
| `POST /registry/review/:name` | ❌ | No route in `main.go` |
| `POST /expand/:name` | ❌ | No route in `main.go` |
| `GET /expand/status/:id` | ❌ | No route in `main.go` |
| `POST /expand/gap-report` | ❌ | No route in `main.go` |
| `POST /learn` | ❌ | No route in `main.go` |
| `GET /learn/status/:id` | ❌ | No route in `main.go` |
| `GET /evidences/:skill_name` | ❌ | No route in `main.go` |

### DEAD CODE (Internal API — `internal/api/server.go`)
The hardened `internal/api.Server` has full CRUD handlers but is **never instantiated** (G01-O3):
- `skills_handler.go` — full CRUD ✅ (dead)
- `search_handler.go` — search ✅ (dead)
- `expand_handler.go` — expand ✅ (dead)
- `learn_handler.go` — learn ✅ (dead)
- `registry_handler.go` — registry ✅ (dead)

### SPEC DRIFT
| Issue | Spec | Live |
|-------|------|------|
| Search verb | `POST /search` | `GET /api/v1/skills/search` |
| Registry missing | `GET /registry/missing` | `GET /api/v1/missing` |
| Expand | `POST /expand/:name` | No route |
| Learn | `POST /learn` | No route |
| Health auth | `ApiKeyAuth` required | Unauthenticated |

---

## §7 MCP Server Tools — Gap Analysis

### IMPLEMENTED
| Tool | Status | Evidence |
|------|--------|----------|
| `skill_search` | ✅ | `internal/mcp/tools.go` |
| `skill_get` | ✅ | `internal/mcp/tools.go` |
| `skill_tree` | ✅ | `internal/mcp/tools.go` (G06 recursive fix landed) |
| `skill_create` | ✅ | `internal/mcp/tools.go` |
| `learn_from_project` | ✅ | `internal/mcp/tools.go` |
| `missing_skills` | ✅ | `internal/mcp/tools.go` |
| `get_coverage` | ✅ | `internal/mcp/tools.go` |

### SPEC DRIFT
| Issue | Spec | Implementation |
|-------|------|----------------|
| `skill_search` | "Vector/hybrid search" | Trigram/ILIKE only (G29 — `VectorSearch` has zero callers) |
| `learn_from_project` | Returns job ID | Returns job ID but no status-check path (G30) |

---

## §8 Configuration — Gap Analysis

### IMPLEMENTED
All config sections exist in `internal/config/config.go`:
- `[server]` ✅
- `[database]` ✅
- `[embedding]` ✅
- `[validation]` ✅
- `[autoexpand]` ✅
- `[codeanalysis]` ✅
- `[mcp]` ✅
- `[registry]` ✅
- `[logging]` ✅

### NOT IMPLEMENTED
| Config Key | Status | Evidence |
|------------|--------|----------|
| `server.enable_http3` | ❌ | HTTP/3 deleted |
| `server.tls_cert`/`tls_key` | ❌ | Not wired |
| `embedding.provider = "anthropic"` | ❌ | Not in factory |
| `validation.sandbox_type = "wasm"` | ❌ | Sandbox deleted (G02) |
| `validation.sandbox_type = "gvisor"` | ❌ | Not implemented |
| `validation.sandbox_type = "docker"` | ❌ | Not implemented |
| `server.allowed_origins` | ⚠️ | Config exists, not in SPEC §8 sample (G18) |
| `server.api_keys` | ⚠️ | Config exists, not in SPEC §8 sample |

---

## Summary: Remaining Work by Priority

### CRITICAL (must fix before any feature work)
1. **G01-O3** — Dead `internal/api.Server` consolidation (all CRUD handlers dead)
2. **G03** — `handleValidate` and `handleCodeAnalysis` still stubs; `runValidationCycle` log-only
3. **G04** — Test coverage still thin (~30 test files vs 53 source files)

### HIGH (core functionality gaps)
4. **G10** — Embedding dimension assertion (DESIGN DONE, impl PENDING)
5. **G12** — tree-sitter native parsing (DESIGN DONE, impl PENDING)
6. **G29** — Hybrid vector search (VectorSearch has zero callers)
7. **G59** — Embedding ingestion (StoreSkillEmbedding unwired)
8. **G14** — Submodule policy conflict
9. **G15** — Aurora/HarmonyOS feasibility spike
10. **G69** — GitHub Skills Source Ingestion (24 sub-items)
11. **G93** — Unified Multi-Source Ingestion (30 sub-items)

### MEDIUM (reliability/ops)
12. **G17** — Weak default DB password
13. **G18** — CORS allowlist SPEC doc update
14. **G20** — Auto-expand placeholder persist (createMinimalDraft)
15. **G60** — Search conflict oracle
16. **G61** — Two divergent /health implementations
17. **G124** — Docs catalog (12 sub-items)

### LOW (cosmetic/deferred)
18. **G54/G62** — gofmt drift
19. **G65/G67/G68** — Ops-script edge cases
