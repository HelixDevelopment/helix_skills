# HelixAgent

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixAgent>
- **Description**: LLM-powered autonomous agent framework providing the core cognitive engine for the Helix ecosystem
- **Category**: AI / Cognitive Engine
- **Status**: Active

## Overview

HelixAgent is the cognitive substrate for the entire Helix agent ecosystem. It implements multi-step reasoning, tool-use, memory integration, and task decomposition for building intelligent agents that can plan, execute, and self-correct across complex workflows. Every higher-level agent component (HelixCode, HelixBuilder, HelixSpecifier) is built on HelixAgent as its agent substrate.

## Tech Stack

- Language: Go (core), with bindings for Python and TypeScript
- Concurrency: Go concurrency patterns, plugin system via Go interfaces
- Architecture: Modular cognitive engine with pluggable reasoning, memory, and tool layers
- Key patterns: Strategy pattern for reasoning, observer pattern for tool events

## Key Features

- Multi-step reasoning engine with chain-of-thought and tree-of-thought planning
- Tool-use framework -- agents discover and invoke external tools, APIs, and services
- Task decomposition -- breaks complex goals into executable sub-task graphs
- Self-correction loops with error detection, diagnosis, and retry strategies
- Plugin architecture for custom reasoning strategies and domain-specific behaviors

## Related Repos

- [LLMProvider](../LLMProvider/README.md) -- core consumer for multi-model LLM access
- [HelixMemory](../HelixMemory/README.md) -- provides long-term and working memory management
- [LLMOrchestrator](../LLMOrchestrator/README.md) -- manages headless agent lifecycle
- [DebateOrchestrator](../DebateOrchestrator/README.md) -- multi-agent consensus on complex decisions
- [HelixCode](../HelixCode/README.md), [HelixBuilder](../HelixBuilder/README.md), [HelixSpecifier](../HelixSpecifier/README.md) -- built on HelixAgent as their cognitive substrate
- [VisionEngine](../VisionEngine/README.md) -- visual reasoning and UI understanding capabilities
- [Agentic](../../vasic-digital/Agentic/README.md) -- graph-based workflow orchestration using HelixAgent nodes

---
*Part of the [HelixDevelopment catalogue](../README.md)*
