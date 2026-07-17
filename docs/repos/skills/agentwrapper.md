# AgentWrapper

> **Repo:** [vasic-digital/AgentWrapper](https://github.com/vasic-digital/AgentWrapper)
> **Type:** vasic-digital repo · **Status:** Active

## Overview

AgentWrapper provides a Docker-based isolation layer for AI CLI
coding agents. It packages agent runtimes into reproducible containers,
ensuring consistent environments across hosts and preventing agent
processes from interfering with the host system. This is a foundational
building block for running multiple agents concurrently.

## Key capabilities

- Docker containerization of CLI-based AI coding agents
- Reproducible agent environments across different host machines
- Process isolation and resource boundary enforcement
- Integration with orchestration layer for multi-agent workflows
- Configurable resource limits (CPU, memory, network)

## Architecture

AgentWrapper operates as a container lifecycle manager:

1. **Container definition** — Dockerfiles and compose configs for
   wrapping CLI agent runtimes
2. **Lifecycle management** — start, stop, health-check, and cleanup
   of agent containers
3. **Resource enforcement** — cgroup-based limits preventing agent
   processes from consuming excessive host resources
4. **Integration bridge** — exposes container status and agent I/O
   to upstream orchestrators

## Integration points

- **LLMOrchestrator** — manages headless CLI agent lifecycles that
  AgentWrapper containers run
- **Agentic** — graph-based orchestration that can compose wrapped
  agents into workflows
- **tmux** — optimized containerized tmux build used for terminal
  multiplexing inside agent containers

## Configuration

Container resource limits, base images, and agent runtime bindings
are configurable. Check the repo for Docker Compose and environment
variable documentation.

## Status

**Active.** Standalone repository. Referenced in the Helix Skills
ecosystem as a core container primitive for agent isolation.
