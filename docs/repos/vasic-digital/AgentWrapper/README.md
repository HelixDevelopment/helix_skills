# AgentWrapper

- **GitHub URL**: <https://github.com/vasic-digital/AgentWrapper>
- **Description**: Wrap AI CLI Coding Agents in Docker containers
- **Category**: Container + Lifecycle
- **Status**: Active

## Overview

AgentWrapper provides a Docker-based isolation layer for AI CLI coding agents. It packages agent runtimes into reproducible containers, ensuring consistent environments across hosts and preventing agent processes from interfering with the host system. This is a foundational building block for running multiple agents concurrently.

## Tech Stack

- Language: Multiple (Docker, Shell)
- Framework: Docker

## Key Features

- Docker containerization of CLI-based AI coding agents
- Reproducible agent environments across different host machines
- Process isolation and resource boundary enforcement
- Integration with orchestration layer for multi-agent workflows

## Related Repos

- [LLMOrchestrator](../LLMOrchestrator/README.md) — manages headless CLI agent lifecycles that AgentWrapper containers run
- [Agentic](../Agentic/README.md) — graph-based orchestration that can compose wrapped agents into workflows
- [tmux](../tmux/README.md) — optimized containerized tmux build used for terminal multiplexing inside agent containers

---
*Part of the [vasic-digital catalogue](../README.md)*
