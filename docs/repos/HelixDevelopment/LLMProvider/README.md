# LLMProvider

- **GitHub URL**: <https://github.com/HelixDevelopment/LLMProvider>
- **Description**: Shared LLM provider abstraction layer with 40+ provider adapters. Provides a unified interface for accessing any LLM provider (OpenAI, Anthropic, Google, Mistral, local models, and dozens more) through a single API, with automatic failover, load balancing, and cost optimization.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Unified LLM API -- single interface for 40+ provider adapters
- Automatic provider failover with configurable fallback chains
- Load balancing across multiple providers and API keys
- Cost optimization with per-token pricing tracking and provider selection
- Rate-limit handling with exponential backoff and provider rotation
- Streaming response support across all providers
- Embedding generation via provider-specific embedding models
- Model capability discovery -- query available models and their features
- Response caching with configurable TTL for repeated queries
- Token usage metering and budget enforcement

## Technology

- **Language**: Go
- **Frameworks**: Go HTTP client, provider-specific SDKs
- **Architecture**: Adapter pattern with provider-specific implementations behind a common interface
- **Key patterns**: Strategy pattern for provider selection, circuit breaker for failover, response streaming

## Integration

- Core dependency of HelixAgent, HelixCode, HelixBuilder, and every LLM-consuming component
- Used by HelixLLM as a local provider adapter alongside cloud providers
- Consumed by LLMOrchestrator for alias resolution and model routing
- Provides embedding generation for HelixMemory vector stores
- Integrates with DebateOrchestrator for multi-model debate participants
- Token usage metrics feed into cost tracking and budget enforcement systems
- Foundation of the native-alias-first priority system in multi-track orchestration

## Status

Active development. 40+ provider adapters operational. Core API interface stable. Streaming and embedding support across providers. Automatic failover and rate-limit handling production-ready. New provider adapters added regularly.
