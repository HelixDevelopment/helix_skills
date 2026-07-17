# LLMGateway

- **GitHub URL**: <https://github.com/vasic-digital/LLMGateway>
- **Description**: Gateway service for routing, load balancing, and managing LLM provider traffic
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

LLMGateway acts as a centralised gateway for routing LLM requests across multiple providers. It handles load balancing, failover, rate limiting, and request transformation so downstream consumers interact with a single endpoint regardless of which provider serves the request. Designed for production deployments where provider reliability, cost optimisation, and traffic management are critical.

## Tech Stack

- Language: Multiple
- Framework: Gateway / reverse proxy pattern

## Key Features

- Centralised routing across multiple LLM providers
- Load balancing and automatic failover on provider errors
- Rate limiting and quota management per provider
- Request transformation and response normalisation
- Health monitoring and provider status tracking

## Related Repos

- [LLMProvider](../LLMProvider/README.md) — provider adapters that the gateway routes through
- [I-LLM](../I-LLM/README.md) — introspection layer for gateway observability
- [LLMOps](../LLMOps/README.md) — operations platform for gateway metrics and experiment tracking

---
*Part of the [vasic-digital catalogue](../README.md)*
