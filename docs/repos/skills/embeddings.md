# Embeddings

> **Repo:** [vasic-digital/Embeddings](https://github.com/vasic-digital/Embeddings.git)
> **Type:** Declared dependency · **Status:** Active · **Org:** vasic-digital

## Overview

A reusable Go module (`digital.vasic.embeddings`) for generating,
storing, and querying text and code embeddings. Provides a
provider-agnostic interface for embedding generation with support
for multiple backends, used as a building block across RAG and
search pipelines.

## Key capabilities

- Provider-agnostic embedding generation interface
- Batch embedding support for large document sets
- Caching and deduplication of embedding vectors
- Multiple embedding model backend support

## Architecture

The embeddings module follows a provider-adapter pattern:

1. **Provider interface** — abstract `Embedder` interface for
   generating vectors from text/code
2. **Backend adapters** — concrete implementations for different
   embedding models (local and remote)
3. **Batch processor** — handles bulk embedding requests with
   configurable concurrency
4. **Cache layer** — deduplicates identical inputs to avoid
   redundant embedding calls

## Integration points

- **token_optimizer** — direct dependency (WS6-L2 embed); provides
  embedding generation for the token optimization pipeline
- **VectorDB** — vector storage and retrieval companion; embeddings
  are stored in VectorDB for similarity search
- **RAG** — retrieval-augmented generation pipeline using embeddings
  for document retrieval
- **Memory** — memory retrieval powered by semantic similarity

## Configuration

Embedding model selection, batch sizes, cache TTL, and provider
credentials are configurable. Check the repo for Go module API
documentation.

## Status

**Active.** Go module dependency consumed via `helix-deps.yaml` by
the token_optimizer submodule.
