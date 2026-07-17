# LLMProvider

> **Repo:** [HelixDevelopment/LLMProvider](https://github.com/HelixDevelopment/LLMProvider.git)
> **Type:** Declared dependency · **Status:** Active · **Org:** HelixDevelopment

## Overview

A unified interface and adapter library for interacting with 40+ LLM
providers. Implements resilience patterns including automatic retry,
circuit breaker, and health monitoring so consumers can switch
providers without changing application code.

## Key capabilities

- Unified provider interface with 40+ adapter implementations
- Automatic retry with exponential backoff
- Circuit breaker pattern for provider fault tolerance
- Health monitoring and provider status reporting
- Provider-agnostic request/response abstraction

## Architecture

LLMProvider follows a layered adapter architecture:

1. **Provider interface** — abstract `Provider` interface defining
   chat/completion/embedding operations
2. **Adapter registry** — 40+ concrete provider adapters (OpenAI,
   Anthropic, Google, local models, etc.)
3. **Resilience layer** — retry, circuit breaker, and timeout
   wrappers applied transparently to all providers
4. **Health monitor** — tracks provider availability and latency,
   feeds into routing decisions

## Integration points

- **token_optimizer** — direct dependency (transport-tier provider
  interface); provides LLM access for the optimization pipeline
- **LLMGateway** — gateway service routing through these adapters
- **I-LLM** — introspection layer for provider observability
- **LLMOps** — operations platform consuming provider metrics
- **conversation** — provider adapters powering conversation
  generation

## Configuration

Provider credentials, retry policies, circuit breaker thresholds,
timeout values, and health check intervals are configurable per
provider. Check the repo for adapter-specific configuration
documentation.

## Status

**Active.** Go module dependency consumed via `helix-deps.yaml` by
the token_optimizer submodule. Hosted under the HelixDevelopment
organisation.
