# HelixTranslate

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixTranslate>
- **Description**: High-performance, enterprise-grade universal ebook translation toolkit supporting any ebook format and any language pair
- **Category**: Content / Translation
- **Status**: Active

## Overview

HelixTranslate is a universal ebook translation system that supports any ebook format (EPUB, MOBI, PDF, FB2, AZW3, CBZ, and more) and any language pair. It features multiple translation engine backends, a REST API with HTTP/3 (QUIC) support, and real-time WebSocket events for streaming translation progress. Format-preserving translation maintains ebook structure, styling, and metadata throughout the process.

## Tech Stack

- Language: Go
- API: REST with HTTP/3 (QUIC) support, WebSocket for real-time events
- Architecture: Microservice with REST API, WebSocket gateway, and translation engine abstraction
- Key patterns: Pipeline-based format processing, engine adapter pattern, streaming responses

## Key Features

- Universal ebook format support (EPUB, MOBI, PDF, FB2, AZW3, CBZ, and more)
- Any-to-any language pair translation with automatic language detection
- Multiple translation engine backends (cloud APIs, local LLMs, hybrid strategies)
- REST API with HTTP/3 (QUIC) support for high-performance client access
- Format-preserving translation -- maintains ebook structure, styling, and metadata

## Related Repos

- [LLMProvider](../LLMProvider/README.md) -- accessing translation-capable LLM models across providers
- [HelixAgent](../HelixAgent/README.md) -- AI-assisted translation quality review
- [HelixMemory](../HelixMemory/README.md) -- translation memory and glossary persistence
- [HelixLLM](../HelixLLM/README.md) -- local LLM inference for offline translation

---
*Part of the [HelixDevelopment catalogue](../README.md)*
