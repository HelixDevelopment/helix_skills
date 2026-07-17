# LLMProvider

- **GitHub URL**: <https://github.com/HelixDevelopment/LLMProvider>
- **Description**: Unified LLM provider abstraction with 40+ adapters, automatic failover, load balancing, cost optimization, and streaming response support
- **Category**: AI / Infrastructure
- **Status**: Active

## Overview

LLMProvider is the foundational LLM access layer for the entire Helix ecosystem. It provides a single interface for 40+ provider adapters (OpenAI, Anthropic, Google, Mistral, local models, and dozens more) with automatic failover, load balancing, cost optimization, and streaming response support. Every LLM-consuming component in the ecosystem routes through LLMProvider, making it the single point of provider management and optimization.

## Tech Stack

- Language: Go
- Integration: Provider-specific SDKs, Go HTTP client
- Architecture: Adapter pattern with provider-specific implementations behind a common interface
- Key patterns: Strategy pattern for provider selection, circuit breaker for failover, response streaming

## Key Features

- Unified LLM API -- single interface for 40+ provider adapters
- Automatic provider failover with configurable fallback chains
- Load balancing across multiple providers and API keys
- Cost optimization with per-token pricing tracking and provider selection
- Rate-limit handling with exponential backoff and provider rotation

## Related Repos

- [HelixAgent](../HelixAgent/README.md) -- core consumer for multi-model LLM access
- [HelixCode](../HelixCode/README.md) -- multi-model code generation
- [HelixBuilder](../HelixBuilder/README.md) -- intelligent build failure analysis
- [LLMOrchestrator](../LLMOrchestrator/README.md) -- alias resolution and model routing
- [HelixLLM](../HelixLLM/README.md) -- local provider adapter alongside cloud providers
- [DebateOrchestrator](../DebateOrchestrator/README.md) -- multi-model debate participants
- [LLMProvider](../../vasic-digital/LLMProvider/README.md) -- related LLM provider in the vasic-digital ecosystem

---
*Part of the [HelixDevelopment catalogue](../README.md)*
