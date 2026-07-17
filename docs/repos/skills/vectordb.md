# VectorDB

> **Repo:** [vasic-digital/VectorDB](https://github.com/vasic-digital/VectorDB.git)
> **Type:** Declared dependency · **Status:** Active · **Org:** vasic-digital

## Overview

A reusable Go module (`digital.vasic.vectordb`) providing vector
database abstractions for similarity search. Implements a
provider-agnostic interface for storing and querying high-dimensional
vectors, used as the retrieval backend in RAG and semantic search
pipelines.

## Key capabilities

- Provider-agnostic vector storage and retrieval interface
- Similarity search with configurable distance metrics
- Batch upsert and query operations
- ANN (Approximate Nearest Neighbor) search support

## Architecture

VectorDB follows a storage-abstraction pattern:

1. **Store interface** — abstract `VectorStore` interface for
   upsert, query, and delete operations
2. **Backend adapters** — concrete implementations for different
   vector database backends (in-memory, file-based, remote)
3. **ANN engine** — approximate nearest neighbor indexing for
   efficient high-dimensional search
4. **Batch processor** — handles bulk vector operations with
   configurable concurrency

## Integration points

- **token_optimizer** — direct dependency (WS6-L2 ANN store);
  provides vector search for the token optimization pipeline
- **Embeddings** — embedding generation that feeds vectors into
  this store
- **RAG** — retrieval-augmented generation consuming vector
  search results
- **Storage** — generic storage abstractions complementing
  vector-specific storage
- **Memory** — memory retrieval powered by vector similarity

## Configuration

Distance metrics, index types, batch sizes, and backend selection
are configurable. Check the repo for Go module API documentation.

## Status

**Active.** Go module dependency consumed via `helix-deps.yaml` by
the token_optimizer submodule.
