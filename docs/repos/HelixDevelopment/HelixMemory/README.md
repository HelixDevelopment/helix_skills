# HelixMemory

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixMemory>
- **Description**: Super Memory Provider -- Unified Cognitive Memory Engine fusing Mem0 + Cognee + Letta for HelixAgent
- **Category**: AI / Cognitive Engine
- **Status**: Active

## Overview

HelixMemory provides multi-tier memory management for the Helix agent ecosystem -- short-term working memory, long-term semantic memory, and episodic memory for persistent agent context across sessions. By fusing Mem0 (personal memory), Cognee (knowledge graphs), and Letta (persistent agents), it gives agents the ability to remember, reason over past experiences, and resume with full context across session boundaries.

## Tech Stack

- Language: Go (core engine), Python (ML integration)
- Storage: Vector database integration (Qdrant, ChromaDB), knowledge graph libraries
- Architecture: Multi-store memory system with unified query interface
- Key patterns: Embedding-based retrieval, graph traversal, relevance scoring

## Key Features

- Multi-tier memory architecture -- working memory, semantic memory, episodic memory
- Unified interface fusing Mem0, Cognee, and Letta
- Semantic search across memory stores with vector embeddings
- Knowledge graph construction and traversal for relational memory
- Session persistence -- agents resume with full context across session boundaries

## Related Repos

- [HelixAgent](../HelixAgent/README.md) -- core dependency for persistent cognitive context
- [HelixCode](../HelixCode/README.md) -- project context persistence across coding sessions
- [DebateOrchestrator](../DebateOrchestrator/README.md) -- stores and recalls debate outcomes
- [LLMProvider](../LLMProvider/README.md) -- embedding generation via configured models
- [HelixSpecifier](../HelixSpecifier/README.md) -- specification similarity matching via embeddings

---
*Part of the [HelixDevelopment catalogue](../README.md)*
