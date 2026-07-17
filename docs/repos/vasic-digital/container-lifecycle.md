# Container + Lifecycle

Back to [vasic-digital index](./README.md) | [Main index](../README.md)

This group covers the container runtime, process lifecycle, and orchestration layer that wraps AI agents and long-running workloads in isolated, reproducible environments. It provides Docker-based agent sandboxing, graph-driven workflow execution, background task scheduling, QA automation orchestration, and optimized containerized tooling (tmux). These repos form the execution substrate that every multi-agent, multi-track, and CI/CD pipeline in the Helix ecosystem depends on.

| Repo | Description | Status |
|---|---|---|
| [AgentWrapper](./AgentWrapper/README.md) | Wraps AI CLI coding agents (Claude, Cursor, Aider, etc.) in Docker containers with configurable resource limits, network isolation, and volume mounts. Provides a uniform sandbox envelope so every agent runs in a hermetic, reproducible environment regardless of host configuration. | Active |
| [Agentic](./Agentic/README.md) | Graph-based agentic workflow orchestration engine written in Go. Models agent workflows as directed acyclic graphs with typed edges, supports conditional branching, parallel fan-out/fan-in, and checkpoint-resume semantics for long-running multi-step tasks. | Active |
| [AutoTemp](./AutoTemp/README.md) | Benchmark-driven temperature auto-tuning orchestration scaffold. Aims to automatically discover optimal LLM sampling temperatures by running evaluation sweeps against task-specific benchmarks and selecting the temperature that maximises quality metrics. | SCAFFOLD / WIP |
| [BackgroundTasks](./BackgroundTasks/README.md) | Go module providing a durable background task queue with at-least-once delivery, retry policies, dead-letter handling, and lifecycle hooks. Used for scheduling long-running operations (builds, scans, syncs) that must survive process restarts. | Active |
| [HelixQA](./HelixQA/README.md) | AI-driven QA orchestration framework for multi-platform testing. Manages test-bank registration, autonomous session dispatch, challenge scoring, and captured-evidence collection across Android, web, desktop, and CLI targets. Integrates with the Challenges submodule for anti-bluff validation. | Active |
| [HyperTune](./HyperTune/README.md) | Hyperparameter tuning orchestration scaffold. Intended to run automated search sweeps (grid, random, Bayesian) over model hyperparameters, track experiment results, and select optimal configurations based on defined metrics. | SCAFFOLD / WIP |
| [LLMOrchestrator](./LLMOrchestrator/README.md) | Headless CLI agent management for LLM orchestration written in Go. Spawns, monitors, resumes, and controls multiple concurrent LLM-powered agents with per-agent session tracking, crash detection, and automatic respawn. Core engine behind the multi-track ruler orchestration. | Active |
| [tmux](./tmux/README.md) | Optimized and verified containerized tmux build -- reproducible across hosts, jemalloc-aware, OOM-protected, with an 8-test verification gate. Provides a reliable terminal multiplexer foundation for multi-agent session management and background process isolation on any Linux system. | Active |

**Related skills:** [agentwrapper](../skills/agentwrapper.md), [agentic](../skills/agentic.md), [llmorchestrator](../skills/llmorchestrator.md), [multitrack](../skills/multitrack.md), [session-orchestrator](../skills/session-orchestrator.md)
