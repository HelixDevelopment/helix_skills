# TOON

> **Repo:** [vasic-digital/TOON](https://github.com/vasic-digital/TOON.git)
> **Type:** Declared dependency · **Status:** Active · **Org:** vasic-digital

## Overview

A reusable Go module (`digital.vasic.toon`) implementing Token-Oriented
Object Notation (TOON) — a serialisation format optimised for LLM
token efficiency. Provides a wrapper around object serialisation that
minimises token count when passing structured data to and from LLMs.

## Key capabilities

- Token-efficient object serialisation format
- LLM-optimised encoding and decoding
- Drop-in replacement for JSON in LLM prompt contexts
- Compact wire format reducing token overhead

## Architecture

TOON is a serialization library with LLM-specific optimizations:

1. **Encoder** — converts Go objects to TOON format, minimizing
   token-expensive characters and redundant structure
2. **Decoder** — parses TOON format back into Go objects with
   full type fidelity
3. **Wire format** — binary/text hybrid format designed for
   minimal token consumption in LLM prompts
4. **Compatibility layer** — JSON-compatible where possible,
   with TOON-specific extensions for token savings

## Integration points

- **token_optimizer** — direct dependency (WS7 wire-encode);
  provides token-efficient serialization for the optimization
  pipeline
- **LLMProvider** — provider adapters consuming TOON-serialised
  payloads for reduced token costs
- **conversation** — conversation context using token-efficient
  formats for compression
- **Embeddings** — embedding request/response serialization

## Configuration

Encoding options (compact vs readable), fallback behaviour, and
custom type registrations are configurable. Check the repo for
Go module API documentation.

## Status

**Active.** Go module dependency consumed via `helix-deps.yaml` by
the token_optimizer submodule.
