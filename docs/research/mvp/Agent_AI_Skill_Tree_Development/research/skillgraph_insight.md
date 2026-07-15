# Insight Extraction: Skill Graph System

## Cross-Dimension Insights

### Insight 1: Single-Database Architecture Eliminates an Entire Class of Integration Failures
- **Derived From**: Dim 02 (Memory Systems), Dim 04 (Vector DB), Dim 01 (CodeGraph)
- **Rationale**: By using PostgreSQL + pgvector as both relational store and vector DB, the system avoids distributed transaction complexity, network partitions between relational and vector stores, and schema synchronization issues. This is validated by Discourse running thousands of pgvector databases and Supabase outperforming dedicated vector DBs.
- **Implications**: Simpler deployment, stronger consistency guarantees, easier backup/restore, no need for event sourcing between stores.
- **Confidence**: High

### Insight 2: The Memory Architecture Should Mirror Human Cognitive Memory Types
- **Derived From**: Dim 02 (Memory Systems), Dim 06 (Validation)
- **Rationale**: Research shows four distinct memory types (Working, Episodic, Semantic, Procedural) mapped to different storage patterns. Our skill graph naturally maps to Semantic (facts about technologies), Evidence table maps to Episodic (experiences from projects), the API cache is Working memory, and the dependency graph structure itself is Procedural (how to do things).
- **Implications**: Design should explicitly separate these memory types with different retention policies, consolidation pipelines (Episodic → Semantic), and access patterns.
- **Confidence**: High

### Insight 3: MCP Universality Means One Integration Covers All Major CLI Agents
- **Derived From**: Dim 03 (MCP/ACP), Dim 05 (Go HTTP/3)
- **Rationale**: MCP is supported by Claude Code, OpenCode, Continue.dev, Claude Desktop, and more. A single MCP server implementation provides universal integration. The discovery that ACP also uses stdio-based JSON-RPC means our transport layer design serves both protocols.
- **Implications**: Focus investment on one robust MCP server rather than N protocol adapters. The plugin ecosystem is converging on MCP as the standard.
- **Confidence**: High

### Insight 4: The Zero-Bluff Guarantee Requires Defense-in-Depth, Not a Single Technique
- **Derived From**: Dim 06 (Validation), Dim 01 (CodeGraph), Dim 02 (Memory)
- **Rationale**: No single validation technique is sufficient. Sandboxed execution catches code errors but not conceptual ones. Multi-model jury catches hallucinations but has shared error floors. Source verification catches drift but not original inaccuracies. The research shows that 5+ stacked controls are needed for production reliability.
- **Implications**: Implement all validation layers: source hashing → sandbox execution → 3-model jury → RAGAS evaluation → periodic re-validation. Each layer catches what others miss.
- **Confidence**: High

### Insight 5: Code Analysis Pipeline Should Be Incremental, Not Batch
- **Derived From**: Dim 01 (CodeGraph), Dim 06 (Validation)
- **Rationale**: The `codebase-memory-mcp` project indexes the Linux kernel in 3 minutes using incremental updates. Full re-indexing on every change is impractical for large codebases. Content hashing enables detecting which files changed.
- **Implications**: Design the code analysis worker with incremental processing from day one. Use git hooks or polling with SHA-based change detection.
- **Confidence**: High

### Insight 6: TOML/JSON Duality Creates Unnecessary Complexity — Simplify
- **Derived From**: Dim 05 (Go HTTP/3), Dim 03 (MCP/ACP)
- **Rationale**: TOML is 5-10x slower to parse than JSON and has no registered MIME type. MCP itself uses JSON for tool schemas. However, TOML is excellent for human-editable config and skill definitions.
- **Implications**: Use TOML for skill definition files (human-written) and JSON for API wire format (machine-optimized). Support TOML responses via Accept header but default to JSON. This eliminates the "Toon" ambiguity (TOML is correct).
- **Confidence**: High

### Insight 7: HTTP/3 Adds Complexity Without Clear API Benefit — Use Strategically
- **Derived From**: Dim 05 (Go HTTP/3), Dim 04 (Vector DB)
- **Rationale**: HTTP/3 improves latency for lossy networks (mobile, global) but adds complexity. For a localhost-first system (CLI agents communicating with local API), HTTP/1.1 or HTTP/2 over TCP is simpler and potentially faster. HTTP/3 matters more for remote deployments.
- **Implications**: Support HTTP/3 for remote deployments but make HTTP/2 the default for local usage. Caddy reverse proxy handles the protocol upgrade transparently.
- **Confidence**: Medium

### Insight 8: The Dependency Graph IS the Product Differentiator
- **Derived From**: Dim 01 (CodeGraph), Dim 02 (Memory), Dim 04 (Vector DB)
- **Rationale**: Vector search alone (RAG) is commoditized. The recursive skill dependency tree with auto-growth is the unique value. Graph traversal provides context that pure vector similarity cannot — understanding that Android AOSP depends on Linux Kernel Modules, which depends on C, which depends on memory management.
- **Implications**: Invest heavily in graph traversal performance, visualization, and the auto-growth pipeline. The dependency tree is the moat.
- **Confidence**: High

### Insight 9: Auto-Growth Must Be Bounded to Prevent Explosion
- **Derived From**: Dim 06 (Validation), Dim 01 (CodeGraph), Dim 02 (Memory)
- **Rationale**: A "full Android AOSP" skill recursively references Java, Kotlin, C++, Linux kernel, Python, Bazel, Git, Make, and each of those references dozens more. Unbounded recursion leads to thousands of skills for a single starting point.
- **Implications**: Implement configurable depth limits, relevance scoring, and human approval gates. Track "coverage percentage" rather than aiming for 100% automated completeness.
- **Confidence**: High

### Insight 10: Learning from Codebases Requires Project-Specific Context
- **Derived From**: Dim 01 (CodeGraph), Dim 06 (Validation)
- **Rationale**: Generic pattern extraction produces low-value skills. The real value comes from understanding project-specific patterns — how THIS codebase uses ViewModel, not just what ViewModel is. Evidence nodes must capture project context (repo, branch, file path).
- **Implications**: Design Evidence schema with rich provenance. Enable project-specific skill overlays that don't pollute the global skill graph.
- **Confidence**: Medium

---

## Strategic Implications

1. **Build the graph first, then the AI**: The dependency graph structure is foundational. Get the DAG traversal, cycle detection, and recursive queries solid before adding LLM features.

2. **Local-first, remote-capable**: Most CLI agents run locally. Optimize for localhost deployment with Docker, support remote via HTTP/3.

3. **Pluggable LLM backend**: Don't hard-code OpenAI. Support local models (Ollama), Anthropic, Google, and others. Validation jury requires model diversity anyway.

4. **Human-in-the-loop by default**: Auto-growth is powerful but risky. Require human approval for new skills until the system proves reliability.

5. **Start narrow, expand deep**: Begin with Android AOSP → Java → Kotlin → C++ chain. Perfect this before adding other technology domains.
