# Cross-Verification Report: Skill Graph System Technologies

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

## Verification Date: 2026-07-15

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

---

## High Confidence Findings (Confirmed by 2+ agents from independent sources)

### 1. pgvector + PostgreSQL 16 is Optimal Choice
- **Confirmed by**: VectorDB agent, Memory Systems agent, CodeGraph agent
- **Evidence**: Discourse uses pgvector in thousands of DBs. Supabase reports 1185% more QPS than Pinecone s1. HNSW handles 100K vectors in 2-3ms query time.
- **Recommendation**: PROCEED with PostgreSQL 16 + pgvector. Use HNSW index with m=32, ef_construction=128.

### 2. tree-sitter Go Bindings are Production-Ready
- **Confirmed by**: CodeGraph agent (primary), Validation agent (supporting)
- **Evidence**: Official `github.com/tree-sitter/go-tree-sitter` actively maintained. Pure-Go alternative (`gotreesitter`) available for CGO-free builds.
- **Recommendation**: Use official bindings. Fallback to pure-Go if CGO issues arise.

### 3. MCP Protocol is Stable and Well-Supported
- **Confirmed by**: MCP/ACP agent (primary), Validation agent (supporting)
- **Evidence**: MCP spec 2025-11-25 stable. mcp-go library at 8.9k stars, actively maintained. Claude Code, OpenCode, Continue.dev all support MCP.
- **Recommendation**: PROCEED with mcp-go for MCP server implementation.

### 4. quic-go HTTP/3 is Production-Ready
- **Confirmed by**: Go HTTP/3 agent (primary), CodeGraph agent (supporting)
- **Evidence**: Used by Caddy, Cloudflare, Traefik in production. Gin v1.11.0 added experimental HTTP/3 support.
- **Recommendation**: Use dual-stack (HTTP/2 + HTTP/3). Consider Caddy as reverse proxy for simplicity.

### 5. 3-Model Jury is Optimal for Validation
- **Confirmed by**: Validation agent (primary), Memory Systems agent (supporting)
- **Evidence**: Research shows diminishing returns after 3-5 models. 3-model jury with diverse architectures gives best accuracy/cost tradeoff.
- **Recommendation**: Use 3 diverse models (e.g., GPT-4o, Claude 3.5, local model) with 2-of-3 approval.

### 6. GraphRAG Significantly Improves Retrieval Quality
- **Confirmed by**: Memory Systems agent (primary), CodeGraph agent (supporting)
- **Evidence**: Microsoft GraphRAG shows 50-70% improvement in comprehensiveness. LightRAG provides lightweight alternative.
- **Recommendation**: Implement LightRAG-style dual retrieval (graph traversal + vector search).

---

## Medium Confidence Findings (Single authoritative source)

### 7. Embedding Dimension: 768 is Sweet Spot
- **Source**: Memory Systems agent research
- **Evidence**: Quality curve flattens after 768. 1536→3072 gives marginal gains at 6x cost.
- **Note**: Conflicts with blueprint's 1536d choice. Either is workable; 768d offers performance gains.

### 8. halfvec Halves Storage with <1% Recall Loss
- **Source**: VectorDB agent research (pgvector docs)
- **Evidence**: 16-bit float quantization available in pgvector. Verified by benchmarks.
- **Recommendation**: Enable halfvec for production deployments.

### 9. Caddy as Reverse Proxy Simplifies HTTP/3
- **Source**: Go HTTP/3 agent research
- **Evidence**: Caddy handles HTTP/3, TLS, Alt-Svc automatically. Go backend serves HTTP/1.1/2.
- **Trade-off**: Adds infrastructure component vs. direct HTTP/3.

---

## Conflict Zones

### C1. TOML as Primary API Format ⚠️
- **Blueprint specifies**: TOML primary, JSON fallback
- **Go HTTP/3 agent finding**: TOML parsing is 5-10x slower than JSON. No registered MIME type. Designed for config, not wire serialization.
- **Resolution**: Keep TOML for configuration files and skill definitions (human-readable). Use JSON for API responses. Keep TOML as supported Accept format but make JSON default.
- **Status**: RESOLVED — modify blueprint

### C2. ACP Uses JSON-RPC over stdio, NOT gRPC ⚠️
- **Blueprint specifies**: "ACP adapter via gRPC for OpenCode"
- **MCP/ACP agent finding**: ACP (Agent Client Protocol by Zed Industries) uses JSON-RPC 2.0 over stdio, not gRPC.
- **Resolution**: ACP adapter must use JSON-RPC over stdio, not gRPC. Remove gRPC dependency for ACP.
- **Status**: RESOLVED — correct blueprint

### C3. Embedding Model Choice
- **Memory Systems agent**: Voyage voyage-code-3 best for code (+13.8% vs OpenAI)
- **VectorDB agent**: OpenAI text-embedding-3-small best value, BGE-M3 best open-source
- **Resolution**: Support pluggable embedding providers. Default to OpenAI text-embedding-3-small for ease, allow local BGE-M3 for privacy.
- **Status**: RESOLVED — design for pluggable embeddings

---

## Low Confidence Findings

### 10. LightRAG Scalability
- **Source**: Limited production usage data
- **Concern**: LightRAG is newer than Microsoft GraphRAG. Long-term stability unproven.
- **Mitigation**: Abstract the GraphRAG interface so either can be swapped.

### 11. auto-growth Pipeline Full Automation
- **Source**: Theoretical design, no production precedent found
- **Concern**: Fully automated skill generation without human review risks quality degradation.
- **Mitigation**: Implement human-in-the-loop for initial deployments. Make full automation opt-in.
