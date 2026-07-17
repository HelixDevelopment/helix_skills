# HelixTranslate

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixTranslate>
- **Description**: A high-performance, enterprise-grade universal ebook translation toolkit supporting any ebook format and any language pair. Features multiple translation engines, REST API with HTTP/3 support, and real-time WebSocket events for streaming translation progress.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Universal ebook format support (EPUB, MOBI, PDF, FB2, AZW3, CBZ, and more)
- Any-to-any language pair translation with automatic language detection
- Multiple translation engine backends (cloud APIs, local LLMs, hybrid strategies)
- REST API with HTTP/3 (QUIC) support for high-performance client access
- Real-time WebSocket events for streaming translation progress and status updates
- Batch translation with queue management and parallel processing
- Format-preserving translation -- maintains ebook structure, styling, and metadata
- Translation memory and glossary support for consistent terminology
- Quality scoring and post-translation review workflows

## Technology

- **Language**: Go
- **Frameworks**: Go HTTP/3 libraries, WebSocket implementation, ebook parsing libraries
- **Architecture**: Microservice with REST API, WebSocket gateway, and translation engine abstraction
- **Key patterns**: Pipeline-based format processing, engine adapter pattern, streaming responses

## Integration

- Uses LLMProvider for accessing translation-capable LLM models across providers
- Integrates with HelixAgent for AI-assisted translation quality review
- Connects to HelixMemory for translation memory and glossary persistence
- Provides API endpoints consumed by web frontends and CLI tools
- Translation quality metrics feed into the helixqa testing framework
- Supports integration with external CAT (Computer-Assisted Translation) tools

## Status

Active development. Core translation pipeline with EPUB and MOBI support operational. HTTP/3 REST API and WebSocket streaming stable. Multiple translation engine backends integrated. Ongoing work on additional format support and translation quality improvements.
