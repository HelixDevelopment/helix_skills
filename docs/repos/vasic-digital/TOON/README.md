# TOON

- **GitHub URL**: <https://github.com/vasic-digital/TOON>
- **Description**: Generic reusable Go module: digital.vasic.toon - Token-Oriented Object Notation wrapper
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

A reusable Go module implementing Token-Oriented Object Notation (TOON) — a serialisation format optimised for LLM token efficiency. Provides a wrapper around object serialisation that minimises token count when passing structured data to and from LLMs.

## Tech Stack

- Language: Go
- Module: digital.vasic.toon

## Key Features

- Token-efficient object serialisation format
- LLM-optimised encoding and decoding
- Drop-in replacement for JSON in LLM prompt contexts

## Related Repos

- [LLMProvider](../LLMProvider/README.md) — provider adapters consuming TOON-serialised payloads
- [conversation](../conversation/README.md) — conversation context using token-efficient formats

---
*Part of the [vasic-digital catalogue](../README.md)*
