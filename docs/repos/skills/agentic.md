# Agentic

> **Repo:** [vasic-digital/Agentic](https://github.com/vasic-digital/Agentic)
> **Type:** vasic-digital repo · **Status:** Active

## Overview

Agentic implements a directed-graph execution engine for composing
multi-step AI agent workflows. Nodes represent individual agent tasks
or decision points, and edges encode dependencies and branching logic.
This enables complex, multi-agent pipelines to be defined declaratively
and executed reliably.

## Key capabilities

- Directed acyclic graph (DAG) based workflow definition
- Declarative node and edge configuration for agent pipelines
- Branching and conditional execution paths
- Integration with containerized agent runtimes
- Graph validation and cycle detection

## Architecture

Agentic is built around a graph execution core:

1. **Graph definition** — workflows declared as nodes (tasks/decisions)
   and edges (dependencies/branches)
2. **Execution engine** — traverses the DAG, dispatching nodes when
   dependencies resolve
3. **Runtime adapters** — connects to containerized agent runtimes
   (via AgentWrapper) for actual task execution
4. **State management** — tracks node completion, branch outcomes,
   and error propagation

## Integration points

- **AgentWrapper** — provides Docker container isolation for the
  agent runtimes Agentic dispatches to
- **LLMOrchestrator** — headless agent management that feeds into
  Agentic workflow graphs
- **HelixQA** — QA orchestration that can leverage Agentic workflows
  for test pipelines
- **BackgroundTasks** — offloads long-running graph nodes

## Configuration

Workflow graphs are defined declaratively (format TBD — check repo
for schema). Node types, edge conditions, and runtime adapter
bindings are configurable per graph.

## Status

**Active.** Standalone repository. Referenced in the Helix Skills
ecosystem as a core orchestration primitive.
