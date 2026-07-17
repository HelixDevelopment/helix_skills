# token_optimizer

> **Path:** `constitution/submodules/token_optimizer/`
> **Type:** Go submodule · **Status:** Active

## What it provides

Token efficiency optimization engine — multi-tier prompt caching,
context-aware routing (§11.4.141). Optimizes agent token usage across
sessions and tracks.

## How consumed

Go module. Used by the constitution's prompt optimization pipeline.

## Source paths

- Module: `constitution/submodules/token_optimizer/`

## Dependencies

- TOON — WS7 wire-encode
- Embeddings — WS6-L2 embed
- VectorDB — WS6-L2 ANN store
- Normalize — Request canonicalization
- LLMProvider — Transport-tier provider interface
- conversation — Infinite-context compression

## Constitution references

§11.4.141
