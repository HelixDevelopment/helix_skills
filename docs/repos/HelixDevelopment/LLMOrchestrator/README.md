# LLMOrchestrator

- **GitHub URL**: <https://github.com/HelixDevelopment/LLMOrchestrator>
- **Description**: Headless CLI agent lifecycle management -- spawning, monitoring, resuming, and controlling agents with crash resilience, rate-limit handling, and multi-alias support
- **Category**: AI / Infrastructure
- **Status**: Active

## Overview

LLMOrchestrator provides the lifecycle management layer for LLM-powered agents. It handles spawning, monitoring, resuming, and controlling headless agent sessions with crash resilience, rate-limit handling, and multi-alias support. The orchestrator is the conductor layer for parallel-development methodology, managing agent pools and ensuring no work is lost to crashes or rate limits.

## Tech Stack

- Language: Go
- Concurrency: Go goroutines and channels, process management
- Architecture: Orchestrator pattern with pluggable agent backends
- Key patterns: Supervisor tree, circuit breaker for rate limits, temp-then-rename state

## Key Features

- Headless agent session lifecycle management (spawn, monitor, resume, terminate)
- Multi-alias support with native-first priority and provider fallback chains
- Rate-limit detection and automatic rebind to healthy aliases
- Crash-resilient session recovery with durable state persistence
- Real-time session monitoring via JSONL event stream consumption

## Related Repos

- [HelixAgent](../HelixAgent/README.md) -- manages HelixAgent instances as headless workers
- [LLMProvider](../LLMProvider/README.md) -- alias resolution and model routing
- [HelixConstitution](../HelixConstitution/README.md) -- multi-track orchestration rules (§11.4.187)
- [HelixGitpx](../HelixGitpx/README.md) -- agent git operations during autonomous workflows
- [HelixCode](../HelixCode/README.md) -- headless coding agent sessions
- [Agentic](../../vasic-digital/Agentic/README.md) -- workflow orchestration that can leverage LLMOrchestrator for agent management

---
*Part of the [HelixDevelopment catalogue](../README.md)*
