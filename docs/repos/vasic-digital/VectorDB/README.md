# VectorDB

- **GitHub URL**: <https://github.com/vasic-digital/VectorDB>
- **Description**: Generic reusable Go module: digital.vasic.vectordb
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

A reusable Go module providing vector database abstractions for similarity search. Implements a provider-agnostic interface for storing and querying high-dimensional vectors, used as the retrieval backend in RAG and semantic search pipelines.

## Tech Stack

- Language: Go
- Module: digital.vasic.vectordb

## Key Features

- Provider-agnostic vector storage and retrieval interface
- Similarity search with configurable distance metrics
- Batch upsert and query operations

## Related Repos

- [Embeddings](../Embeddings/README.md) — embedding generation that feeds vectors into this store
- [RAG](../RAG/README.md) — retrieval-augmented generation consuming vector search results
- [Storage](../Storage/README.md) — generic storage abstractions complementing vector-specific storage

---
*Part of the [vasic-digital catalogue](../README.md)*
