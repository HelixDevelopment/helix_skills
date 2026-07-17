# Agentic

- **GitHub URL**: <https://github.com/vasic-digital/Agentic>
- **Description**: Graph-based agentic workflow orchestration
- **Category**: Container + Lifecycle
- **Status**: Active

## Overview

Agentic implements a directed-graph execution engine for composing multi-step AI agent workflows. Nodes represent individual agent tasks or decision points, and edges encode dependencies and branching logic. This enables complex, multi-agent pipelines to be defined declaratively and executed reliably.

## Tech Stack

- Language: Multiple
- Framework: Custom graph execution engine

## Key Features

- Directed acyclic graph (DAG) based workflow definition
- Declarative node and edge configuration for agent pipelines
- Branching and conditional execution paths
- Integration with containerized agent runtimes

## Related Repos

- [LLMOrchestrator](../LLMOrchestrator/README.md) — headless agent management that feeds into Agentic workflow graphs
- [AgentWrapper](../AgentWrapper/README.md) — provides containerized agent runtimes for graph nodes
- [HelixQA](../HelixQA/README.md) — QA orchestration that can leverage Agentic workflows for test pipelines

---
*Part of the [vasic-digital catalogue](../README.md)*
