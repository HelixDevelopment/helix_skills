# I-LLM

- **GitHub URL**: <https://github.com/vasic-digital/I-LLM>
- **Description**: Introspection layer for LLM providers
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

I-LLM provides an introspection and observability layer for LLM provider interactions. It captures, structures, and exposes metadata about LLM requests and responses — including token usage, latency, model routing decisions, and error patterns — enabling operators to monitor provider health, debug integration issues, and optimise cost. Designed to sit between application code and LLMProvider adapters as a transparent instrumentation layer.

## Tech Stack

- Language: Multiple
- Framework: LLM observability instrumentation

## Key Features

- Transparent request/response capture for LLM provider calls
- Token usage tracking and cost estimation per request
- Latency measurement and percentile reporting
- Error pattern classification and provider health scoring
- Integration with observability and ops platforms

## Related Repos

- [LLMProvider](../LLMProvider/README.md) — provider adapters that I-LLM instruments
- [LLMOps](../LLMOps/README.md) — operations platform consuming I-LLM's introspection data
- [LLMGateway](../LLMGateway/README.md) — gateway service using I-LLM for routing decision observability

---
*Part of the [vasic-digital catalogue](../README.md)*
