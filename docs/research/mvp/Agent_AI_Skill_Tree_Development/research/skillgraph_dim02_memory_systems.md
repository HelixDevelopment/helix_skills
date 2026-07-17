# SkillGraph Memory Systems Research: Comprehensive Findings

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

**Research Date:** 2026-01-18
**Researcher:** AI Research Agent
**Scope:** LLM Memory Systems for Persistent Agent Knowledge
**Blueprint Context:** PostgreSQL 16 + pgvector, 1536-dim embeddings (OpenAI ada-002 compatible), HNSW indexing, Graph-RAG approach, Evidence table for learned experiences

---

## Table of Contents

**Revision:** 1
**Last modified:** 2026-07-17T23:11:30Z

1. [pgvector Production Performance](#1-pgvector-production-performance)
2. [Memory Architectures for LLM Agents](#2-memory-architectures-for-llm-agents)
3. [Embedding Models Comparison](#3-embedding-models-comparison)
4. [GraphRAG Memory Approaches](#4-graphrag-memory-approaches)
5. [Production Memory Systems](#5-production-memory-systems)
6. [Summary & Recommendations](#6-summary--recommendations)
7. [Danger Zones](#7-danger-zones)

---

## 1. pgvector Production Performance

### 1.1 Performance at Scale

**Claim:** pgvector comfortably handles up to 10 million 1536-dimensional vectors on a single managed PostgreSQL instance with proper HNSW indexing and halfvec quantization. Many teams run tens of millions with careful tuning and partial indexes. The practical rule of thumb: if your HNSW index fits in RAM, performance is good. [^30^]
**Source:** Firecrawl Vector Database Comparison 2026
**URL:** https://www.firecrawl.dev/blog/best-vector-databases
**Date:** 2025-10-09
**Excerpt:** "Beyond 50-100M vectors, extensions hit throughput and latency limits that purpose-built systems avoid."
**Context:** Comparison of all major vector databases including pgvector, Pinecone, Milvus, Qdrant, Weaviate
**Confidence:** High

**Claim:** Real-world deployments of pgvector at 10 million vectors require approximately 20-25 KB per vector including HNSW index overhead. At this scale, memory pressure becomes a central design concern. At 100 million vectors, raw FP32 storage alone exceeds 600 GB before index overhead. [^100^]
**Source:** ClickHouse Engineering Blog - Scaling Vector Search in Postgres
**URL:** https://clickhouse.com/resources/engineering/scale-vector-search-postgres
**Date:** 2026-06-19
**Excerpt:** "Real-world deployments can land around 20 to 25 kilobytes per vector. At ten million vectors, memory pressure becomes a central design concern."
**Context:** Detailed analysis of pgvector scaling thresholds and optimization strategies
**Confidence:** High

**Claim:** pgvector with pgvectorscale achieves 471 QPS at 99% recall on 50M vectors, which is 11.4x better than Qdrant's 41 QPS at the same recall. [^30^]
**Source:** Firecrawl Vector Database Comparison (citing Timescale benchmarks)
**URL:** https://www.firecrawl.dev/blog/best-vector-databases
**Date:** 2025-10-09
**Excerpt:** "pgvectorscale achieves 471 QPS (Queries Per Second) at 99% recall on 50M vectors. That's 11.4x better than Qdrant's 41 QPS at the same recall."
**Context:** May 2025 benchmarks comparing pgvector+pgvectorscale against dedicated vector databases
**Confidence:** Medium (vendor-sponsored benchmark)

### 1.2 HNSW Index Tuning

**Claim:** The HNSW `m` parameter (maximum connections per node) should be set based on dimensionality. For 384-dimensional vectors, m=16 is sufficient. For 1536-dimensional vectors (OpenAI embeddings), m should be 32 or even 64 to maintain graph navigability. [^85^]
**Source:** pgvector HNSW Performance Tuning at Scale
**URL:** https://pradeepbhandari.com/blog/pgvector-hnsw-performance-tuning-scaling-guide
**Date:** 2026-07-15
**Excerpt:** "For 1536-dimensional embeddings, like those from OpenAI, you'll likely need to bump m to 32 or even 64 to keep the graph navigable."
**Context:** Detailed tuning guide for pgvector HNSW indexes at million-row scale
**Confidence:** High

**Claim:** The recommended starting point for HNSW index creation with 1536d vectors is: `m = 32`, `ef_construction = 128`. Increasing `ef_construction` beyond 128 or 256 offers diminishing returns for most RAG applications. [^85^]
**Source:** pgvector HNSW Performance Tuning at Scale
**URL:** https://pradeepbhandari.com/blog/pgvector-hnsw-performance-tuning-scaling-guide
**Date:** 2026-07-15
**Excerpt:** "CREATE INDEX ON document_embeddings USING hnsw (embedding vector_cosine_ops) WITH (m = 32, ef_construction = 128);"
**Context:** Practical SQL examples and tuning recommendations
**Confidence:** High

**Claim:** For fast index builds on high-volume datasets, set `max_parallel_maintenance_workers = 8` and `maintenance_work_mem = '4GB'` before creating the index. HNSW index creation in pgvector is parallel-aware. [^85^]
**Source:** pgvector HNSW Performance Tuning at Scale
**URL:** https://pradeepbhandari.com/blog/pgvector-hnsw-performance-tuning-scaling-guide
**Date:** 2026-07-15
**Excerpt:** "SET max_parallel_maintenance_workers = 8; SET maintenance_work_mem = '4GB'; CREATE INDEX ON document_embeddings USING hnsw..."
**Context:** Memory configuration for large-scale index builds
**Confidence:** High

### 1.3 Quantization and Storage Optimization

**Claim:** Using `halfvec` (16-bit float) instead of `vector` (32-bit float) reduces storage by 50% with minimal recall loss. For 1M vectors at 1536 dimensions, this turns a ~6GB index into ~3GB. This should be the default for production pgvector deployments. [^123^]
**Source:** Neon Blog - Use halfvec Instead and Save 50%
**URL:** https://neon.com/blog/dont-use-vector-use-halvec-instead-and-save-50-of-your-storage-cost
**Date:** 2024-07-31
**Excerpt:** "halfvec not only helps you reduce the memory and storage of your tables and indexes by half but also speeds up index build without compromising recall."
**Context:** Benchmarks on DBpedia 1M 1536-dimension dataset
**Confidence:** High

**Claim:** Binary quantization with `bit` type achieves 32x storage reduction but with significant recall degradation for 1536-dimensional vectors. It works better as a pre-filter (fast candidate retrieval followed by full-precision reranking) rather than standalone search. [^124^]
**Source:** Jonathan Katz Blog - Scalar and Binary Quantization for pgvector
**URL:** https://jkatz.github.io/post/postgres/pgvector-scalar-binary-quantization/
**Date:** 2024-04-09
**Excerpt:** "Binary quantization works better for datasets that produce larger vectors that can be distinguished by their bits."
**Context:** Comprehensive benchmarks comparing vector, halfvec, and bit types across multiple datasets
**Confidence:** High

**Claim:** Matryoshka Representation Learning (MRL) models like OpenAI text-embedding-3-small allow truncating 1536-dim vectors to 512 or 768 dimensions with minimal quality loss. The recommended two-stage retrieval pattern: use 256-dim truncated vectors for initial ANN pass, then re-score top candidates at full dimensions. [^46^]
**Source:** Tensoria Embedding Models 2026 Guide
**URL:** https://tensoria.fr/en/blog/embedding-models-2026-guide
**Date:** 2026-05-15
**Excerpt:** "You can truncate a 1536-dimensional MRL embedding to 256 dimensions and retain approximately 93 to 95% of the retrieval quality."
**Context:** Detailed analysis of dimension trade-offs and Matryoshka patterns
**Confidence:** High

### 1.4 pgvector Limitations

**Claim:** pgvector inherits PostgreSQL's single-node architecture, creating practical limits. Performance degrades noticeably beyond 10-20 million vectors. High-concurrency workloads (thousands of QPS) strain PostgreSQL's connection model. Standard auto-vacuuming won't re-optimize HNSW internal graph connections - manual REINDEX is required when index size balloons beyond 1.5x raw data size. [^130^]
**Source:** VeloDB - What Is pgvector?
**URL:** https://www.velodb.io/glossary/what-is-pgvector
**Date:** 2026-02-04
**Excerpt:** "Performance degrades noticeably beyond 10-20 million vectors, depending on dimensionality and hardware."
**Context:** Comprehensive overview of pgvector capabilities and limitations
**Confidence:** High

**Claim:** A single PostgreSQL instance hits hard walls with pgvector: memory constraints (1M 1536d embeddings require ~6GB raw storage before indexing), I/O bottlenecks (vector searches don't parallelize well), and concurrent query limitations. At 50M+ vectors, distributed SQL or dedicated vector databases become necessary. [^32^]
**Source:** Medium - Outgrowing PostgreSQL: Why Distributed SQL Is the Only Way
**URL:** https://medium.com/@dikhyantkrishnadalai/outgrowing-postgresql-why-distributed-sql-is-the-only-way-to-scale-vector-search-beyond-50m
**Date:** 2025-07-19
**Excerpt:** "A million 1536-dimensional embeddings (OpenAI's ada-002) require ~6GB just for raw storage, before indexing overhead."
**Context:** Analysis of pgvector limitations at enterprise scale
**Confidence:** High

### 1.5 pgvectorscale and StreamingDiskANN

**Claim:** pgvectorscale from Timescale adds StreamingDiskANN, a disk-oriented ANN index that reduces the requirement for the vector graph to remain entirely in RAM. It achieves 28x lower p95 latency and 16x higher query throughput than Pinecone's s1 index at 99% recall. [^83^]
**Source:** DBVis - pgvectorscale Extension for Improved Vector Search
**URL:** https://www.dbvis.com/thetable/pgvectorscale-an-extension-for-improved-vector-search-in-postgres/
**Date:** 2025-09-03
**Excerpt:** "PostgreSQL with pgvector and pgvectorscale delivers 28x lower p95 latency and 16x higher query throughput than Pinecone's storage-optimized (s1) index."
**Context:** Benchmark on 50 million Cohere embeddings of 768 dimensions each
**Confidence:** Medium (vendor-sponsored benchmark)

**Claim:** StreamingDiskANN uses Statistical Binary Quantization (SBQ) keeping vectors on disk efficiently while maintaining high recall. pgvectorscale supports vectors up to 16,000 dimensions. [^88^]
**Source:** DBI Services - pgvector Guide for DBA
**URL:** https://www.dbi-services.com/blog/pgvector-a-guide-for-dba-part-2-indexes-update-march-2026/
**Date:** 2026-03-01
**Excerpt:** "StreamingDiskANN uses SBQ compression to reduce storage and speed up traversal. Supports vectors up to 16,000 dimensions."
**Context:** Technical comparison table of HNSW, IVFFlat, and DiskANN options for pgvector
**Confidence:** High

### 1.6 pgvector Decision Framework

**Claim:** For teams already running PostgreSQL with <50M vectors, pgvector + pgvectorscale is the recommended starting point. Purpose-built databases like Milvus or Pinecone should be considered when scaling to >100M vectors or when vector search is the primary workload. [^30^]
**Source:** Firecrawl Vector Database Comparison
**URL:** https://www.firecrawl.dev/blog/best-vector-databases
**Date:** 2025-10-09
**Excerpt:** "Existing PostgreSQL (<50M vectors): Start with pgvector + pgvectorscale. New project (>100M vectors): Milvus/Zilliz Cloud or Pinecone serverless."
**Context:** Decision framework based on infrastructure and scale
**Confidence:** High

---

## 2. Memory Architectures for LLM Agents

### 2.1 MemGPT Architecture

**Claim:** MemGPT introduces an OS-inspired multi-level memory architecture with two primary memory types: main context (analogous to RAM - the LLM prompt tokens) and external context (analogous to disk storage). Main context consists of three contiguous sections: system instructions (read-only), working context (read/write), and a FIFO queue. [^35^]
**Source:** MemGPT Paper (arXiv:2310.08560)
**URL:** https://arxiv.org/pdf/2310.08560
**Date:** 2023
**Excerpt:** "MemGPT's OS-inspired multi-level memory architecture delineates between two primary memory types: main context (analogous to main memory/physical memory/RAM) and external context (analogous to disk memory/disk storage)."
**Context:** Original MemGPT paper from UC Berkeley
**Confidence:** High

**Claim:** MemGPT uses function calls that the LLM processor can invoke to manage its own memory without user intervention. This includes searching recall storage, managing archival memory, and paging relevant context in and out. The LLM itself manages the memory cycle autonomously. [^37^]
**Source:** Information Matters - MemGPT Engineering Semantic Memory
**URL:** https://informationmatters.org/2025/10/memgpt-engineering-semantic-memory-through-adaptive-retention-and-context-summarization/
**Date:** 2026-05-08
**Excerpt:** "What makes MemGPT particularly innovative is its use of the LLM itself as the memory manager. Through self-directed memory editing via tool calling, the system can actively manage its own memory contents."
**Context:** Analysis of MemGPT's memory management innovations
**Confidence:** High

**Claim:** MemGPT implements strategic forgetting through summarization and targeted deletion. When token usage approaches a threshold (e.g., 70% capacity), the system inserts an internal alert. The LLM then halts current reasoning, reviews working memory, determines least critical content, summarizes it, and writes to external storage. [^41^]
**Source:** Serokell - Design Patterns for Long-Term Memory in LLM-Powered Architectures
**URL:** https://serokell.io/blog/design-patterns-for-long-term-memory-in-llm-powered-architectures
**Date:** 2025-12-09
**Excerpt:** "When token usage in the primary context approaches a defined threshold (e.g., 70% capacity), the system inserts an internal alert... the LLM halts its current reasoning, reviews its working memory, determines which content is least critical, summarizes it, and writes it to the appropriate external tier."
**Context:** Comprehensive overview of memory design patterns including MemGPT
**Confidence:** High

### 2.2 Episodic vs Semantic Memory

**Claim:** The agent-engineering community has converged on a taxonomy borrowed from cognitive science (CoALA paper): Working memory (active context), Episodic memory (records of specific past experiences), Semantic memory (general facts and knowledge), and Procedural memory (knowledge of how to do things). [^36^]
**Source:** AI Agent Memory Design Guide (hidekazu-konishi.com)
**URL:** https://hidekazu-konishi.com/entry/ai_agent_memory_design_guide.html
**Date:** 2026-06-11
**Excerpt:** "The canonical academic source is the CoALA paper (Cognitive Architectures for Language Agents), which maps the memory modules of classic cognitive architectures onto language agents."
**Context:** Detailed taxonomy of agent memory types grounded in research
**Confidence:** High

**Claim:** Episodic memory operates across four stages: encoding (capturing full event context), retrieval (pulling relevant episodes), consolidation (transforming episodes into durable semantic knowledge), and eviction (managing what gets dropped). Consolidation is identified as the most important and least implemented stage. [^39^]
**Source:** Atlan - Episodic Memory for AI Agents
**URL:** https://atlan.com/know/episodic-memory-ai-agents/
**Date:** 2026-04-17
**Excerpt:** "Consolidation, turning raw experience into reusable knowledge, is the most important and the least implemented."
**Context:** Detailed explanation of episodic memory stages for AI agents
**Confidence:** High

**Claim:** The canonical implementation of memory consolidation is the reflection mechanism from Generative Agents (Park et al. 2023): periodically synthesizing recent episodes by recency, relevance, and salience into higher-level insights that are written to semantic memory. [^39^]
**Source:** Atlan - Episodic Memory for AI Agents
**URL:** https://atlan.com/know/episodic-memory-ai-agents/
**Date:** 2026-04-17
**Excerpt:** "The canonical implementation is the reflection mechanism from Generative Agents (Park et al. 2023): periodically synthesizing recent episodes by recency, relevance, and salience into higher-level insights."
**Context:** Reference to the Generative Agents paper and its reflection mechanism
**Confidence:** High

**Claim:** Memory retrieval in generative agents uses a composite scoring function: score(i) = alpha_recency * recency(i) + alpha_importance * importance(i) + alpha_relevance * relevance(i), where importance is LLM-inferred and relevance is embedding-space cosine similarity. [^93^]
**Source:** Emergent Mind - Generative Agents
**URL:** https://www.emergentmind.com/topics/generative-agents
**Date:** 2026-01-21
**Excerpt:** "score(i) = alpha_recency * recency(i) + alpha_importance * importance(i) + alpha_relevance * relevance(i)"
**Context:** Analysis of Park et al. 2023 Generative Agents memory mechanisms
**Confidence:** High

### 2.3 Short-Term vs Long-Term Memory Patterns

**Claim:** Existing research has predominantly treated long-term memory (LTM) and short-term memory (STM) as independent components. STM is commonly enhanced through RAG, while LTM management has progressed along separate lines (trigger-based and agent-based paradigms). This separation leads to fragmented memory construction and suboptimal performance. [^34^]
**Source:** Agentic Memory Paper (arXiv:2601.01885v1)
**URL:** https://arxiv.org/html/2601.01885v1
**Date:** 2025
**Excerpt:** "Existing architectures generally follow two patterns: (a) static STM with trigger-based LTM, or (b) static STM with agent-based LTM. In both settings, the two memory systems are optimized independently and later combined in an ad hoc way."
**Context:** AgeMem framework proposing unified memory management for LLM agents
**Confidence:** High

**Claim:** The AgeMem (Agentic Memory) framework proposes a unified approach where both LTM and STM are jointly managed via explicit tool-based operations. The LLM autonomously invokes memory operations for both types through a three-stage progressive RL strategy: first acquire LTM storage capabilities, then learn STM context management, then coordinate both. [^34^]
**Source:** Agentic Memory Paper (arXiv:2601.01885v1)
**URL:** https://arxiv.org/html/2601.01885v1
**Date:** 2025
**Excerpt:** "We propose Agentic Memory (AgeMem), a unified framework that jointly manages LTM and STM... through a unified tool-based interface, the LLM autonomously invokes and executes memory operations for both LTM and STM."
**Context:** Novel framework for unified memory management in LLM agents
**Confidence:** High

**Claim:** The three memory types most production systems use are: Semantic memory (facts and concepts independent of time), Episodic memory (time-indexed experiences), and Procedural memory (skills and routines). Most systems use a mix, with episodic memory often consolidated into semantic memory over time. [^38^]
**Source:** Redis Blog - Long-Term Memory Architectures for AI Agents
**URL:** https://redis.io/blog/long-term-memory-architectures-ai-agents/
**Date:** 2026-04-28
**Excerpt:** "Semantic memory stores facts and concepts independent of time or context. Episodic memory records time-indexed experiences and events. Procedural memory captures skills and routines for performing tasks."
**Context:** Production-focused memory architecture patterns
**Confidence:** High

### 2.4 Production Memory Frameworks Comparison

**Claim:** In 2026, four memory solutions dominate: Mem0 (48K GitHub stars, general-purpose), Zep (temporal knowledge graphs), LangMem (LangGraph-native), and Letta (OS-inspired tiered memory, formerly MemGPT). Mem0 achieves 92.5 LoCoMo benchmark score vs Letta's 74.0. [^81^]
**Source:** NiteAgent - AI Agent Memory Comparison 2026
**URL:** https://niteagent.com/blog/ai-agent-memory-comparison-2026/
**Date:** 2026-05-17
**Excerpt:** "In 2026, four solutions dominate -- Mem0 (48K stars, general-purpose), Zep (temporal knowledge graphs, 63.8% LongMemEval), LangMem (LangGraph-native agent SDK), and Letta (OS-inspired tiered memory, 83.2% LongMemEval)."
**Context:** Comparison of leading agent memory frameworks with benchmarks
**Confidence:** Medium (some data may be promotional)

**Claim:** Mem0 achieves the lowest search latency among memory methods (p50: 0.148s, p95: 0.200s) compared to LangMem (p50: 17.99s) and A-Mem (p50: 0.668s). Mem0's graph-enhanced variant (Mem0g) achieves highest J score (68.44%) across all methods. [^71^]
**Source:** Mem0 Paper (arXiv:2504.19413)
**URL:** https://arxiv.org/pdf/2504.19413
**Date:** 2025
**Excerpt:** "Our proposed Mem0 approach achieves the lowest search latency among all methods (p50: 0.148s, p95: 0.200s)."
**Context:** Benchmark results on LOCOMO dataset
**Confidence:** High

**Claim:** Letta (formerly MemGPT) creates framework lock-in - it is a full agent runtime, not just a memory layer. Teams using LangChain or LlamaIndex would need a complete architectural rewrite. Additionally, relying on the LLM to explicitly call memory functions creates reliability gaps and higher inference costs. [^95^]
**Source:** Evermind - Best Letta Alternatives for AI Agent Memory
**URL:** https://evermind.ai/blogs/letta-alternative
**Date:** 2026-06-29
**Excerpt:** "Letta is not just a memory layer; it is a full agent runtime. To use Letta's memory, your agents must run inside the Letta environment."
**Context:** Analysis of friction points when using Letta in production
**Confidence:** High

---

## 3. Embedding Models Comparison

### 3.1 Best Models for Code/Technical Knowledge

**Claim:** Voyage AI voyage-code-3 consistently tops retrieval benchmarks for code and technical documentation, outperforming OpenAI text-embedding-3-large by 13.80% and CodeSage-large by 16.81% on a suite of 238 code retrieval datasets. It is specifically optimized for code search and code-to-text retrieval. [^79^]
**Source:** Hugging Face - voyageai/voyage-code-3
**URL:** https://huggingface.co/voyageai/voyage-code-3
**Date:** Current
**Excerpt:** "voyage-code-3 is optimized for code retrieval, outperforming OpenAI-v3-large and CodeSage-large by an average of 13.80% and 16.81% on a suite of 238 code retrieval datasets."
**Context:** Official model card with benchmark results
**Confidence:** High

**Claim:** VoyageCode3 features: 32K token context, supports embeddings of 2048/1024/512/256 dimensions with Matryoshka learning, multiple quantization options (float, int8, uint8, binary), and comprehensive training across 300+ programming languages. [^78^]
**Source:** Modal Blog - 6 Best Code Embedding Models Compared
**URL:** https://modal.com/blog/6-best-code-embedding-models-compared
**Date:** 2025-03-31
**Excerpt:** "VoyageCode3: trained on trillions of tokens with carefully tuned code-to-text ratio, comprehensive dataset with docstring-code and code-code pairs across 300+ programming languages."
**Context:** Comparison guide specifically for code embedding models
**Confidence:** High

**Claim:** BGE-M3 from BAAI is the strongest open-source embedding model available. It supports dense, sparse, and multi-vector retrieval in a single model, enabling hybrid search without separate models. Multilingual support covers 100+ languages. Requires at minimum an A10G GPU for self-hosting. [^48^]
**Source:** PE Collective - Best Embedding Models 2026
**URL:** https://pecollective.com/tools/best-embedding-models/
**Date:** 2026-04-06
**Excerpt:** "BGE-M3 from BAAI is the strongest open-source embedding model available. It supports dense, sparse, and multi-vector retrieval in a single model."
**Context:** Comprehensive benchmark of 6 embedding models on 50K technical documents
**Confidence:** High

### 3.2 OpenAI Models

**Claim:** OpenAI text-embedding-3-large is the safest default choice, scoring near the top of MTEB benchmarks. Matryoshka representation support allows reducing dimensions from 3072 down to 256 with minimal quality loss. However, it is not the absolute best on any single benchmark. [^48^]
**Source:** PE Collective - Best Embedding Models 2026
**URL:** https://pecollective.com/tools/best-embedding-models/
**Date:** 2026-04-06
**Excerpt:** "OpenAI text-embedding-3-large: Best Overall. Scores near the top of MTEB benchmarks across English retrieval, classification, and clustering tasks."
**Context:** Detailed comparison of 6 embedding models
**Confidence:** High

**Claim:** OpenAI text-embedding-3-small at $0.02/1M tokens is the recommended default for most teams starting out. Excellent quality at low cost, Matryoshka-compatible for truncation to 512 or 768 later without re-embedding. [^121^]
**Source:** DanubeData - Build a RAG System with pgvector
**URL:** https://danubedata.ro/blog/pgvector-rag-managed-postgres-2026
**Date:** 2026-07-15
**Excerpt:** "For most teams starting out, text-embedding-3-small at 1536 dimensions is the right default: excellent quality, cheap, and Matryoshka-compatible so you can truncate to 512 or 768 later without re-embedding."
**Context:** Practical guide for building RAG with pgvector on managed PostgreSQL
**Confidence:** High

### 3.3 Local/Open-Source Models

**Claim:** Nomic Embed v2 punches above its weight class at 137M parameters - small enough to run on CPU in production. Supports Matryoshka dimensions (768 down to 64), long context up to 8192 tokens. The quality-to-size ratio is the best in the market for local deployment. [^48^]
**Source:** PE Collective - Best Embedding Models 2026
**URL:** https://pecollective.com/tools/best-embedding-models/
**Date:** 2026-04-06
**Excerpt:** "Nomic Embed v2 punches way above its weight class. At 137M parameters, it's small enough to run on a CPU in production. The quality-to-size ratio is the best in the market."
**Context:** Comparison of embedding models for different deployment scenarios
**Confidence:** High

**Claim:** Jina Embeddings v3 is best for long documents with 8192-token context and 572M parameters. Available via Apache 2.0 license. Well-suited for technical documentation where full context matters. [^48^]
**Source:** PE Collective - Best Embedding Models 2026
**URL:** https://pecollective.com/tools/best-embedding-models/
**Date:** 2026-04-06
**Excerpt:** "Jina Embeddings v3: Best for Long Documents. 8K context, strong on document-level retrieval."
**Context:** Benchmark results on 50K technical documents
**Confidence:** High

### 3.4 Dimension Trade-offs

**Claim:** The retrieval quality curve flattens quickly after 768 dimensions for most tasks. Going from 256 to 768 gives meaningful recall gains. Going from 1536 to 3072 gives marginal gains at roughly 6x the storage cost. For most RAG applications, 768 or 1024 dimensions is the practical sweet spot. [^46^]
**Source:** Tensoria Embedding Models 2026 Guide
**URL:** https://tensoria.fr/en/blog/embedding-models-2026-guide
**Date:** 2026-05-15
**Excerpt:** "The retrieval quality curve flattens quickly after 768 dimensions for most tasks. Going from 256 to 768 gives meaningful recall gains. Going from 1536 to 3072 gives marginal gains at roughly 6x the storage cost."
**Context:** Detailed analysis of dimension trade-offs with concrete storage numbers
**Confidence:** High

**Claim:** For 10M vectors: 3072 dims (float32) = ~117 GB raw + 30-40 GB HNSW index; 1536 dims = ~58 GB raw + 15-20 GB indexed; 768 dims = ~29 GB raw + 8-12 GB indexed; 256 dims = ~10 GB raw + 3-4 GB indexed. [^46^]
**Source:** Tensoria Embedding Models 2026 Guide
**URL:** https://tensoria.fr/en/blog/embedding-models-2026-guide
**Date:** 2026-05-15
**Excerpt:** "Concrete numbers for 10M vectors: 3072 dimensions (float32): ~117 GB raw, ~30 to 40 GB with HNSW index at ef_construction=200"
**Context:** Storage cost analysis by dimension
**Confidence:** High

**Claim:** Computing similarity between 384-dimensional vectors is roughly 4x faster than 1536-dimensional vectors. For real-time search, lower dimensions enable sub-50ms responses that higher dimensions can't match. [^51^]
**Source:** Particula Tech - How Many Dimensions Should Your Embeddings Have?
**URL:** https://particula.tech/blog/embedding-dimensions-rag-vector-search
**Date:** 2025-12-22
**Excerpt:** "Computing similarity between 384-dimensional vectors is roughly 4x faster than 1536-dimensional vectors."
**Context:** Practical guide for choosing embedding dimensions
**Confidence:** High

---

## 4. GraphRAG Memory Approaches

### 4.1 Microsoft GraphRAG

**Claim:** Microsoft Research's GraphRAG (released mid-2024) uses a knowledge graph instead of or alongside vector search to retrieve structured context. The pipeline: Source Documents -> Text Chunks -> Entity/Relationship Extraction -> Knowledge Graph -> Community Detection (Leiden algorithm) -> Community Summaries -> Query-Focused Answers. [^72^]
**Source:** Microsoft GraphRAG Paper (arXiv:2404.16130v2)
**URL:** https://arxiv.org/html/2404.16130v2
**Date:** 2024
**Excerpt:** "Graph indexing with a 600 token window took 281 minutes for the Podcast dataset... We implemented Leiden community detection using the graspologic library."
**Context:** Original Microsoft GraphRAG research paper
**Confidence:** High

**Claim:** GraphRAG improved comprehensiveness of answers by 50-70% over baseline RAG on complex multi-document summarization tasks. A 2024 Stanford study confirmed graph-augmented retrieval reduced factual errors by 35-45% compared to vector-only approaches on multi-hop question answering. [^70^]
**Source:** Atlan - What Is GraphRAG?
**URL:** https://atlan.com/know/what-is-graphrag/
**Date:** 2026-03-30
**Excerpt:** "Microsoft Research found that GraphRAG improved comprehensiveness of answers by 50-70% over baseline RAG on complex, multi-document summarization tasks."
**Context:** Comprehensive overview of GraphRAG architecture and benefits
**Confidence:** High

**Claim:** The Leiden algorithm is used for hierarchical community detection in GraphRAG, allowing examination of communities at multiple granularity levels. The algorithm identifies 5 levels of communities, with higher levels being more specific. [^82^]
**Source:** Bertelsmann Tech Blog - How Microsoft GraphRAG Works
**URL:** https://tech.bertelsmann.com/en/blog/articles/how-microsoft-graphrag-works-step-by-step-part-12
**Date:** 2025-07-28
**Excerpt:** "The Leiden algorithm, a hierarchical clustering method, to identify communities within the graph. One advantage of using a hierarchical community detection algorithm is the ability to examine communities at multiple levels of granularity."
**Context:** Step-by-step explanation of GraphRAG implementation
**Confidence:** High

**Claim:** Uncapped 5-hop queries on enterprise-scale graphs (100M+ nodes) can exceed 10-second response times, compared to sub-second latency for 2-hop queries. Production systems need query-time budgets that cap traversal depth. [^70^]
**Source:** Atlan - What Is GraphRAG?
**URL:** https://atlan.com/know/what-is-graphrag/
**Date:** 2026-03-30
**Excerpt:** "Uncapped 5-hop queries on enterprise-scale graphs (100M+ nodes) can exceed 10-second response times."
**Context:** Performance analysis of graph traversal in production GraphRAG
**Confidence:** High

### 4.2 Hybrid Search (Vector + Keyword + Graph)

**Claim:** Modern agent architectures often use hybrid memory combining the precise, symbolic recall of graphs with the broad, semantic recall of vector embeddings. A common pattern: use the knowledge graph to identify context, then use vector search for details within that context. [^67^]
**Source:** ZBrain - Role of Knowledge Graphs in Building Agentic AI
**URL:** https://zbrain.ai/knowledge-graphs-for-agentic-ai/
**Date:** 2026-05-05
**Excerpt:** "A common integration pattern is to use the knowledge graph to identify context, then use vector search for details within that context."
**Context:** Analysis of hybrid graph+vector memory approaches
**Confidence:** High

**Claim:** The two-stage retrieval pattern: the knowledge graph provides focus (narrowing relevant entities and relationships), and the vector store provides depth (detailed unstructured information). This has been shown to significantly improve both accuracy and efficiency. [^67^]
**Source:** ZBrain - Role of Knowledge Graphs in Building Agentic AI
**URL:** https://zbrain.ai/knowledge-graphs-for-agentic-ai/
**Date:** 2026-05-05
**Excerpt:** "The two-stage retrieval means the knowledge graph provides focus... and the vector store provides depth. The result is more precise and efficient."
**Context:** Detailed explanation of hybrid retrieval patterns
**Confidence:** High

**Claim:** Agent-as-a-Graph retrieval represents tools and parent agents as nodes and edges in a knowledge graph. It achieves 14.9% and 14.6% improvements in Recall@5 and nDCG@5 over prior state-of-the-art retrievers through vector search + weighted reciprocal rank fusion + graph traversal. [^66^]
**Source:** PwC Paper - Knowledge Graph-Based Tool and Agent Retrieval
**URL:** https://arxiv.org/html/2511.18194v1
**Date:** 2025
**Excerpt:** "Agent-as-a-Graph... achieving 14.9% and 14.6% improvements in Recall@5 and nDCG@5 over prior state-of-the-art retrievers."
**Context:** Research paper on KG-based retrieval for multi-agent systems
**Confidence:** High

### 4.3 Context Assembly from Graph Neighborhoods

**Claim:** Graph-RAG can use fewer tokens than pure vector approaches for the same queries because it supplies a concise representation of relevant info. Fewer tokens means faster responses and lower cost. The graph provides an "evidence trail" making outcomes traceable and auditable. [^67^]
**Source:** ZBrain - Role of Knowledge Graphs in Building Agentic AI
**URL:** https://zbrain.ai/knowledge-graphs-for-agentic-ai/
**Date:** 2026-05-05
**Excerpt:** "Graph-RAG can use fewer tokens for the same queries, since it supplies a concise representation of relevant info. Fewer tokens means faster responses and lower cost."
**Context:** Benefits of graph-based retrieval for enterprise applications
**Confidence:** High

---

## 5. Production Memory Systems

### 5.1 Production RAG Architecture Lessons

**Claim:** A reliable production RAG platform has five major layers: Ingestion, Processing/Indexing, Retrieval/Ranking, Generation, and Observability/Operations. The retrieval layer is the heart - if it is weak, the LLM will sound polished but still be wrong. [^68^]
**Source:** Dev.to - Inside a Production RAG System
**URL:** https://dev.to/seasia_infotech_899dc2c59/inside-a-production-rag-system-architecture-stack-and-lessons-learned-28h7
**Date:** 2026-03-19
**Excerpt:** "This retrieval layer is the heart of the production RAG system. If it is weak, the LLM will sound polished but still be wrong."
**Context:** Architecture and lessons from production RAG deployment
**Confidence:** High

**Claim:** The biggest production lesson: retrieval quality matters more than model size. Better chunking, metadata, and reranking usually improve results more than changing the LLM. Chunking is architecture, not preprocessing - structure-aware chunking dramatically improves quality. [^68^]
**Source:** Dev.to - Inside a Production RAG System
**URL:** https://dev.to/seasia_infotech_899dc2c59/inside-a-production-rag-system-architecture-stack-and-lessons-learned-28h7
**Date:** 2026-03-19
**Excerpt:** "Lesson 1: Retrieval quality matters more than model size. Lesson 3: Chunking is architecture, not preprocessing."
**Context:** Five hard-won lessons from production RAG
**Confidence:** High

**Claim:** Reranking is not optional in production RAG - it is a lifeline. Skipping reranking resulted in 40% of answers pulling topically adjacent but irrelevant chunks. A lightweight reranker cut irrelevant retrievals by more than half. Chunking strategy matters more than the embedding model choice. [^71^]
**Source:** Medium - Production-Grade RAG: Architecture, Trade-offs, and Hard-Won Lessons
**URL:** https://medium.com/@pavansaish/production-grade-rag-architecture-trade-offs-hard-won-lessons-bc28fcc6b8b8
**Date:** 2025-11-03
**Excerpt:** "I skipped reranking in my first attempt to save latency. Result? 40% of answers pulled in topically adjacent but irrelevant chunks."
**Context:** First-hand account of production RAG challenges
**Confidence:** High

### 5.2 Agent Memory Production Patterns

**Claim:** The read-before-reasoning, write-after-acting loop is a common production pattern for agent memory: (1) Receive input, (2) Memory read, (3) Reason and plan, (4) Act, (5) Observe, (6) Memory write, (7) Loop. Retrieval quality and write discipline usually decide whether it works in production. [^38^]
**Source:** Redis Blog - Long-Term Memory Architectures for AI Agents
**URL:** https://redis.io/blog/long-term-memory-architectures-ai-agents/
**Date:** 2026-04-28
**Excerpt:** "The hardest part is context assembly: given everything that could go into the context window, what should actually go in?"
**Context:** Production memory architecture patterns from Redis
**Confidence:** High

**Claim:** Memory evaluation in production should measure: retrieval precision and recall, staleness detection, contradiction handling, and user satisfaction. Benchmarks like LoCoMo and LongMemEval are used to test memory quality over long conversations. [^74^]
**Source:** GitHub - Agent Memory Techniques (NirDiamant)
**URL:** https://github.com/NirDiamant/Agent_Memory_Techniques
**Date:** 2026-05-30
**Excerpt:** "Measure memory quality. Check retrieval precision and recall, staleness, contradictions, and user satisfaction."
**Context:** Comprehensive collection of 30 agent memory techniques
**Confidence:** High

### 5.3 Knowledge Graphs for Skill Systems

**Claim:** Knowledge and Skill Graphs (KSG) extend traditional knowledge graphs to process both static and dynamic knowledge. KSG has entity nodes (fact, environment, skill nodes) and attribute nodes (entity description, skill display, pre-train network, offline dataset). Two types of directed edges indicate entity-entity and entity-attribute relationships. [^84^]
**Source:** KSG Paper (arXiv:2209.05698)
**URL:** https://arxiv.org/pdf/2209.05698
**Date:** 2022
**Excerpt:** "KSG has two types of nodes, entities and attributes, as well as two types of directed edges, which indicate entity-entity and entity-attribute relationships."
**Context:** Original paper introducing Knowledge and Skill Graph concept
**Confidence:** High

**Claim:** Graph Neural Networks (GCN, GraphSAGE, GAT) can be applied to skill graphs to construct skill embeddings that model interdependence between skills. The GraphKT architecture uses a GNN to produce node embeddings for skills, concatenated with correctness values, passed through a causal transformer. [^131^]
**Source:** Stanford CS224W - Graph Neural Networks for Knowledge Tracing
**URL:** https://medium.com/stanford-cs224w/graph-neural-networks-for-knowledge-tracing-ef31fdaa5f00
**Date:** 2023-05-15
**Excerpt:** "We propose creating an explicit graph structure describing the interaction between skills in student problem solving sequences, leveraging a graph neural network to construct skill embeddings."
**Context:** Stanford course project on GNNs for skill learning
**Confidence:** High

---

## 6. Summary & Recommendations

### 6.1 pgvector Assessment for SkillGraph

**Verdict: RECOMMENDED for the SkillGraph system, with specific tuning.**

Given the blueprint context (PostgreSQL 16 + pgvector, <10M vectors expected for skill knowledge), pgvector is an excellent choice:

1. **Unified data model**: Vectors live alongside relational skill data, evidence tables, and graph relationships in a single transaction
2. **HNSW index**: Use `m=32, ef_construction=128` for 1536-dimensional embeddings
3. **Quantization**: Implement `halfvec` to halve storage (6GB -> 3GB per 1M vectors) with minimal recall loss
4. **Scaling headroom**: 10M vectors is well within pgvector's comfortable range
5. **Partial indexes**: Create per-skill-type or per-project partial indexes for faster filtered retrieval

**Migration trigger**: Consider dedicated vector DB only if exceeding 50M vectors or needing >1000 QPS sustained throughput.

### 6.2 Memory Architecture Recommendation

**Recommended four-tier memory architecture for SkillGraph:**

| Tier | Type | Storage | Purpose | Implementation |
|------|------|---------|---------|----------------|
| Working Memory | Short-term | Context window | Active skill execution context | LLM prompt management |
| Semantic Memory | Long-term | pgvector + relational | Skill knowledge, facts, patterns | PostgreSQL tables with vector search |
| Episodic Memory | Long-term | PostgreSQL JSONB | Project experiences, outcomes | Evidence table with full context |
| Procedural Memory | Long-term | PostgreSQL + code | Learned workflows, validated patterns | Skill graph + executable templates |

**Key design principles:**
1. **Consolidation pipeline**: Periodically synthesize episodic records (project experiences) into semantic knowledge (skill patterns) using a reflection mechanism
2. **Composite retrieval**: Score memories by `recency * importance * relevance` (following Generative Agents)
3. **Graph awareness**: Use the skill graph to provide structural context for vector retrieval (hybrid graph+vector approach)

### 6.3 Embedding Model Recommendation

**For code/technical skill knowledge:**

| Scenario | Model | Dimensions | Cost |
|----------|-------|------------|------|
| Best code retrieval quality | Voyage voyage-code-3 | 1024-2048 | $0.18/1M tokens |
| Best overall value | OpenAI text-embedding-3-small | 1536 (Matryoshka) | $0.02/1M tokens |
| Best open-source local | BGE-M3 | 1024 | Free (GPU) |
| Best CPU-only local | Nomic Embed v2 | 768 (Matryoshka) | Free |

**Recommendation**: Start with OpenAI text-embedding-3-small at 1536 dimensions (Matryoshka-compatible). This allows truncating to 768 or 512 later without re-embedding if storage/performance needs require it. For GDPR-sensitive deployments, use BGE-M3 self-hosted.

### 6.4 GraphRAG Integration

**Recommended hybrid approach for SkillGraph:**

1. **Knowledge layer**: PostgreSQL tables for skill nodes, experience nodes, and relationship edges
2. **Vector layer**: pgvector HNSW index on skill/experience embeddings
3. **Retrieval pattern**: Two-stage - (a) graph traversal to identify relevant skill neighborhoods, (b) vector search within those neighborhoods for semantic similarity
4. **Community detection**: Use Leiden algorithm for identifying skill clusters at multiple granularity levels
5. **Context assembly**: Serialize retrieved subgraph + vector-matched chunks into structured LLM prompt

---

## 7. Danger Zones

### 7.1 Critical Warnings

1. **HNSW index bloat**: Every vector update in PostgreSQL is a delete+insert due to MVCC. Dead rows accumulate in the HNSW graph, degrading query quality. Plan for periodic `REINDEX CONCURRENTLY` operations. Monitor index size - refresh when it exceeds 1.5x raw data size.

2. **Memory wall**: If your HNSW index exceeds available RAM, query latency degrades from milliseconds to seconds. Size your PostgreSQL `shared_buffers` to keep the working set cache-resident. Use `halfvec` to halve memory requirements.

3. **Consolidation complexity**: The most important AND least implemented memory operation is consolidation (episodic -> semantic). Without active consolidation, episodic memories accumulate endlessly and retrieval quality degrades. Implement a periodic reflection job.

4. **Embedding lock-in**: Switching embedding models requires re-indexing everything. This is weeks of engineering time, not just an API change. Choose Matryoshka-compatible models to preserve flexibility.

5. **Filter overfiltering**: Pre-filtering vector search (applying WHERE clauses before ANN) can silently degrade recall. Use pgvector 0.8+ iterative scans (`hnsw.iterative_scan = relaxed_order`) to fix this.

### 7.2 Anti-Patterns to Avoid

- **Storing only vector similarity**: Semantic search alone misses exact matches. Always combine with keyword/graph search.
- **Ignoring chunk boundaries**: Naive chunking breaks code context. Use AST-aware or semantic chunking for code.
- **Treating memory as append-only**: Without forgetting/consolidation, memory systems become slower and less accurate over time.
- **No observability**: Track retrieval precision, recall, and latency percentiles. Bad retrieval degrades gracefully - the LLM compensates with hallucinations.

### 7.3 Scaling Thresholds

| Scale | Vectors | Action Required |
|-------|---------|----------------|
| Comfortable | <5M | Basic pgvector + HNSW |
| Monitor | 5-10M | Add halfvec, tune memory, partial indexes |
| Optimize | 10-20M | pgvectorscale/DiskANN, read replicas |
| Migrate | >50M | Evaluate dedicated vector DB |

---

## References

[^30^] Firecrawl, "Best Vector Databases in 2026," 2025-10-09
[^32^] Medium/Dalai, "Outgrowing PostgreSQL: Why Distributed SQL Is the Only Way," 2025-07-19
[^34^] arXiv:2601.01885v1, "Agentic Memory: Learning Unified Long-Term and Short-Term Memory Management"
[^35^] arXiv:2310.08560, "MemGPT: Towards LLMs as Operating Systems"
[^36^] hidekazu-konishi.com, "AI Agent Memory Design Guide," 2026-06-11
[^37^] Information Matters, "MemGPT: Engineering Semantic Memory," 2026-05-08
[^38^] Redis Blog, "Long-Term Memory Architectures for AI Agents," 2026-04-28
[^39^] Atlan, "Episodic Memory for AI Agents," 2026-04-17
[^41^] Serokell, "Design Patterns for Long-Term Memory in LLM-Powered Architectures," 2025-12-09
[^46^] Tensoria, "Embedding Models in 2026," 2026-05-15
[^48^] PE Collective, "Best Embedding Models 2026," 2026-04-06
[^51^] Particula Tech, "How Many Dimensions Should Your Embeddings Have?," 2025-12-22
[^66^] arXiv:2511.18194v1, "Knowledge Graph-Based Tool and Agent Retrieval"
[^67^] ZBrain, "Role of Knowledge Graphs in Building Agentic AI," 2026-05-05
[^68^] Dev.to, "Inside a Production RAG System," 2026-03-19
[^69^] Salfati Group, "Graph RAG Guide 2025," 2026-03-12
[^70^] Atlan, "What Is GraphRAG?," 2026-03-30
[^71^] arXiv:2504.19413, "Mem0: Building Production-Ready AI Agents"
[^72^] arXiv:2404.16130v2, "A GraphRAG Approach to Query-Focused Summarization"
[^79^] Hugging Face, "voyageai/voyage-code-3" model card
[^81^] NiteAgent, "AI Agent Memory Comparison 2026," 2026-05-17
[^83^] DBVis, "pgvectorscale: An Extension for Improved Vector Search," 2025-09-03
[^84^] arXiv:2209.05698, "KSG: Knowledge and Skill Graph"
[^85^] pradeepbhandari.com, "pgvector HNSW Performance Tuning," 2026-07-15
[^88^] DBI Services, "pgvector Guide for DBA Part 2," 2026-03-01
[^93^] Emergent Mind, "Generative Agents," 2026-01-21
[^100^] ClickHouse, "How to Scale Vector Search in Postgres," 2026-06-19
[^121^] DanubeData, "Build a RAG System with pgvector," 2026-07-15
[^123^] Neon Blog, "Don't use vector. Use halfvec instead," 2024-07-31
[^124^] Jonathan Katz, "Scalar and binary quantization for pgvector," 2024-04-09
[^130^] VeloDB, "What Is pgvector?," 2026-02-04
[^131^] Stanford CS224W, "Graph Neural Networks for Knowledge Tracing," 2023-05-15
