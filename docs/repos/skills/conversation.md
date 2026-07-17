# conversation

> **Repo:** [vasic-digital/conversation](https://github.com/vasic-digital/conversation.git)
> **Type:** Declared dependency · **Status:** Active · **Org:** vasic-digital

## Overview

A library for managing conversation context in AI agent systems.
Implements infinite context compression techniques and event-sourced
conversation state, enabling agents to maintain coherent long-running
dialogues without hitting context window limits.

## Key capabilities

- Conversation context window management and overflow handling
- Infinite context compression and summarisation
- Event-sourced conversation history with replay support
- Context-aware chunking for efficient LLM consumption

## Architecture

The conversation module operates as a middleware layer:

1. **Event sourcing** — every message/action is recorded as an
   append-only event, enabling full history replay
2. **Compression engine** — summarises older context segments to
   fit within LLM window limits while preserving key information
3. **Context window manager** — tracks token usage and triggers
   compression when approaching limits
4. **State persistence** — conversation state stored durably for
   cross-session continuity

## Integration points

- **token_optimizer** — direct dependency; conversation feeds
  compressed context into the token optimization pipeline
- **Memory** — memory abstractions used for conversation state
  persistence
- **Normalize** — input normalisation applied before conversation
  processing
- **LLMProvider** — provider adapters powering conversation
  generation
- **TOON** — token-efficient serialization for compressed context

## Configuration

Compression thresholds, summarisation strategies, and persistence
backends are configurable. Check the repo for Go module API
documentation.

## Status

**Active.** Go module dependency consumed via `helix-deps.yaml` by
the token_optimizer submodule.
