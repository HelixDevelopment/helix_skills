# HelixMemory

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixMemory>
- **Description**: Super Memory Provider -- a unified cognitive memory engine fusing Mem0, Cognee, and Letta for HelixAgent. Provides multi-tier memory management with short-term working memory, long-term semantic memory, and episodic memory for persistent agent context across sessions.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Multi-tier memory architecture -- working memory, semantic memory, episodic memory
- Unified interface fusing Mem0 (personal memory), Cognee (knowledge graphs), and Letta (persistent agents)
- Semantic search across memory stores with vector embeddings
- Knowledge graph construction and traversal for relational memory
- Automatic memory consolidation -- working memory promoted to long-term based on relevance scoring
- Session persistence -- agents resume with full context across session boundaries
- Memory pruning and summarization to manage context window constraints
- Multi-agent memory isolation with shared knowledge pools

## Technology

- **Language**: Go (core engine), Python (ML integration)
- **Frameworks**: Vector database integration (Qdrant, ChromaDB), knowledge graph libraries
- **Architecture**: Multi-store memory system with unified query interface
- **Key patterns**: Embedding-based retrieval, graph traversal, relevance scoring

## Integration

- Core dependency of HelixAgent for persistent cognitive context
- Used by HelixCode for project context persistence across coding sessions
- Integrates with DebateOrchestrator for storing and recalling debate outcomes
- Provides embeddings consumed by HelixSpecifier for specification similarity matching
- Connects to LLMProvider for embedding generation via configured models
- Memory stores feed into the session-resumption system (CONTINUATION.md generation)

## Status

Active development. Core memory tiers (working, semantic) are operational. Episodic memory and knowledge graph construction are in active development. Mem0 and Cognee integration stable; Letta integration advancing. Vector search performance continuously optimized.
