# LLMProvider

- **GitHub URL**: <https://github.com/vasic-digital/LLMProvider>
- **Description**: Shared LLM provider interface, 40+ provider adapters, retry, circuit breaker, health monitoring
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

A unified interface and adapter library for interacting with 40+ LLM providers. Implements resilience patterns including automatic retry, circuit breaker, and health monitoring so consumers can switch providers without changing application code.

## Tech Stack

- Language: Multiple

## Key Features

- Unified provider interface with 40+ adapter implementations
- Automatic retry with exponential backoff
- Circuit breaker pattern for provider fault tolerance
- Health monitoring and provider status reporting

## Related Repos

- [LLMGateway](../LLMGateway/README.md) — gateway service routing through these adapters
- [I-LLM](../I-LLM/README.md) — introspection layer for observability
- [LLMOps](../LLMOps/README.md) — operations platform consuming provider metrics

---
*Part of the [vasic-digital catalogue](../README.md)*
