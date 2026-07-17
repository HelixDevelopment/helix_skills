# LLMOrchestrator

> **Repo:** [vasic-digital/LLMOrchestrator](https://github.com/vasic-digital/LLMOrchestrator)
> **Type:** vasic-digital repo · **Status:** Active

## Overview

LLMOrchestrator manages the lifecycle of headless CLI-based LLM
agents. It handles agent spawning, health monitoring, task
dispatching, and graceful shutdown across multiple concurrent agent
instances. This serves as the central coordination point for
multi-agent systems that run without a GUI.

## Key capabilities

- Headless CLI agent spawning and lifecycle management
- Concurrent multi-agent dispatch and coordination
- Health monitoring and automatic restart of failed agents
- Task routing across available agent instances
- Graceful shutdown and resource cleanup

## Architecture

LLMOrchestrator is structured as a lifecycle manager:

1. **Agent pool** — maintains a pool of available agent instances
   with health status tracking
2. **Task dispatcher** — routes incoming tasks to available agents
   based on capability and load
3. **Health monitor** — periodic liveness checks with automatic
   restart on failure
4. **Shutdown coordinator** — graceful drain and cleanup of agent
   resources on termination

## Integration points

- **AgentWrapper** — provides Docker container isolation for the
  agents LLMOrchestrator manages
- **Agentic** — graph-based workflows that consume orchestrated
  agent outputs
- **BackgroundTasks** — offloads long-running work from orchestrated
  agents
- **LLMProvider** — provider adapters used by managed agents for
  LLM communication

## Configuration

Agent pool size, health check intervals, task routing strategies,
and restart policies are configurable. Check the repo for
CLI flags and environment variable documentation.

## Status

**Active.** Standalone repository. Referenced in the Helix Skills
ecosystem as a core orchestration primitive.
