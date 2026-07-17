# conversation

- **GitHub URL**: <https://github.com/vasic-digital/conversation>
- **Description**: Conversation context management, infinite context compression, and event sourcing for AI agents
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

conversation provides a comprehensive library for managing conversation context in AI agent systems. It implements infinite context compression techniques that allow agents to maintain coherent long-running dialogues without hitting context window limits, and uses event sourcing to preserve full conversation history for replay and analysis. Designed as the context backbone for any agent that needs to sustain multi-turn interactions across extended sessions.

## Tech Stack

- Language: Go
- Module: digital.vasic.conversation

## Key Features

- Conversation context window management with configurable retention strategies
- Infinite context compression and summarisation for long-running dialogues
- Event-sourced conversation history with full replay capability
- Context segmentation and prioritisation for relevance-aware truncation
- Integration hooks for memory, tool calls, and provider adapters

## Related Repos

- [Memory](../Memory/README.md) — long-term memory abstractions used for conversation state persistence
- [LLMProvider](../LLMProvider/README.md) — provider adapters powering conversation generation and summarisation
- [TOON](../TOON/README.md) — token-efficient serialisation for compact conversation context payloads

---
*Part of the [vasic-digital catalogue](../README.md)*
