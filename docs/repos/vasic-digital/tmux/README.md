# tmux

- **GitHub URL**: <https://github.com/vasic-digital/tmux>
- **Description**: Optimized + verified containerized tmux build -- reproducible across hosts, jemalloc-aware, OOM-protected, 8-test verification gate. Reusable on any Linux system.
- **Category**: Container + Lifecycle
- **Status**: Active

## Overview

This repository provides a production-hardened, containerized build of tmux that is reproducible across different Linux hosts. It integrates jemalloc for memory efficiency, applies OOM-killer protections, and runs an 8-test verification gate to ensure build correctness before deployment. The build is designed for reuse in any containerized agent or dev environment.

## Tech Stack

- Language: Shell, C
- Framework: Docker, tmux, jemalloc

## Key Features

- Reproducible containerized tmux build across Linux hosts
- jemalloc integration for optimized memory allocation
- OOM-killer protection for stable long-running sessions
- 8-test verification gate ensuring build integrity

## Related Repos

- [AgentWrapper](../AgentWrapper/README.md) — uses containerized tmux for terminal multiplexing inside agent containers

---
*Part of the [vasic-digital catalogue](../README.md)*
