# LLMOrchestrator

- **GitHub URL**: <https://github.com/vasic-digital/LLMOrchestrator>
- **Description**: Headless CLI agent management for LLM orchestration
- **Category**: Container + Lifecycle
- **Status**: Active

## Overview

LLMOrchestrator manages the lifecycle of headless CLI-based LLM agents. It handles agent spawning, health monitoring, task dispatching, and graceful shutdown across multiple concurrent agent instances. This serves as the central coordination point for multi-agent systems that run without a GUI.

## Tech Stack

- Language: Multiple
- Framework: Custom agent lifecycle manager

## Key Features

- Headless CLI agent spawning and lifecycle management
- Concurrent multi-agent dispatch and coordination
- Health monitoring and automatic restart of failed agents
- Task routing across available agent instances

## Related Repos

- [AgentWrapper](../AgentWrapper/README.md) — provides Docker container isolation for the agents LLMOrchestrator manages
- [Agentic](../Agentic/README.md) — graph-based workflows that consume orchestrated agent outputs
- [BackgroundTasks](../BackgroundTasks/README.md) — offloads long-running work from orchestrated agents

---
*Part of the [vasic-digital catalogue](../README.md)*
