# TOON

- **GitHub URL**: <https://github.com/vasic-digital/TOON>
- **Description**: Generic reusable Go module: digital.vasic.toon - Token-Oriented Object Notation wrapper
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

TOON (Token-Oriented Object Notation) is a compact serialisation format optimised for LLM token efficiency. It wraps structured data into a notation that minimises token count when passed as context to language models, reducing cost and latency compared to JSON or YAML for equivalent payloads. Designed as a drop-in replacement for JSON in LLM prompt construction where token budget matters.

## Tech Stack

- Language: Go
- Module: digital.vasic.toon

## Key Features

- Compact token-oriented serialisation format reducing LLM token usage
- Drop-in replacement for JSON in prompt construction contexts
- Bidirectional conversion between TOON and standard JSON
- Optimised for structured data frequently passed in LLM contexts
- Low-overhead encoding and decoding with minimal allocations

## Related Repos

- [LLMProvider](../LLMProvider/README.md) — provider adapters that can use TOON for compact prompt payloads
- [conversation](../conversation/README.md) — conversation context management that benefits from token-efficient serialisation
- [Normalize](../Normalize/README.md) — input canonicalisation that may process TOON-encoded payloads

---
*Part of the [vasic-digital catalogue](../README.md)*
